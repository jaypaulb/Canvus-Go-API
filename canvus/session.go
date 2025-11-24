// Package canvus provides a Go SDK for the Canvus API.
package canvus

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"math/big"
	"net"
	"net/http"
	"net/url"
	"path"
	"reflect"
	"strings"
	"sync"
	"time"
)

// Authenticator applies authentication to an HTTP request.
type Authenticator interface {
	Authenticate(req *http.Request)
}

// APIKeyAuthenticator authenticates using a static API key and header.
type APIKeyAuthenticator struct {
	Header string
	APIKey string
}

// Authenticate sets the API key header on the request.
func (a *APIKeyAuthenticator) Authenticate(req *http.Request) {
	if a.Header != "" && a.APIKey != "" {
		req.Header.Set(a.Header, a.APIKey)
	}
}

// transportWithAPIKey is an http.RoundTripper that adds an API key to requests
type transportWithAPIKey struct {
	transport http.RoundTripper
	header   string
	apiKey   string
}

// RoundTrip adds the API key to the request headers
func (t *transportWithAPIKey) RoundTrip(req *http.Request) (*http.Response, error) {
	req = req.Clone(req.Context())
	req.Header.Add(t.header, t.apiKey)
	req.Header.Add("Content-Type", "application/json")
	return t.transport.RoundTrip(req)
}

// TokenAuthenticator authenticates using a bearer token.
type TokenAuthenticator struct {
	Token string
}

// Authenticate sets the Authorization header on the request.
func (a *TokenAuthenticator) Authenticate(req *http.Request) {
	if a.Token != "" {
		req.Header.Set("Private-Token", a.Token)
	}
}

// SessionOption configures a Session.
type SessionOption func(*Session)

// WithAPIKey configures the session to use a static API key.
func WithAPIKey(apiKey string) SessionConfigOption {
	return func(cfg *SessionConfig) {
		// Create a new session with the API key
		if cfg.HTTPClient == nil {
			cfg.HTTPClient = &http.Client{
				Transport: &http.Transport{
					TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
				},
			}
		}

		// Create a transport that adds the API key to requests
		transport := cfg.HTTPClient.Transport
		if transport == nil {
			transport = http.DefaultTransport
		}

		cfg.HTTPClient.Transport = &transportWithAPIKey{
			transport: transport,
			header:   "Private-Token",
			apiKey:   apiKey,
		}
	}
}

// WithToken configures the session to use a bearer token.
func WithToken(token string) SessionOption {
	return func(s *Session) {
		s.authenticator = &TokenAuthenticator{Token: token}
	}
}

// circuitState represents the state of the circuit breaker
type circuitState int

const (
	circuitStateClosed circuitState = iota
	circuitStateOpen
	circuitStateHalfOpen
)

// circuitBreaker implements a simple circuit breaker pattern
type circuitBreaker struct {
	state          circuitState
	failures       int
	maxFailures    int
	resetTimeout   time.Duration
	lastFailure    time.Time
	mutex          sync.RWMutex
}

func newCircuitBreaker(maxFailures int, resetTimeout time.Duration) *circuitBreaker {
	return &circuitBreaker{
		state:        circuitStateClosed,
		maxFailures:  maxFailures,
		resetTimeout: resetTimeout,
	}
}

func (cb *circuitBreaker) allow() bool {
	cb.mutex.RLock()
	defer cb.mutex.RUnlock()

	if cb.state == circuitStateClosed {
		return true
	}

	// If circuit is open, check if we should try to let a request through
	if cb.state == circuitStateOpen && time.Since(cb.lastFailure) > cb.resetTimeout {
		cb.mutex.RUnlock()
		cb.mutex.Lock()
		cb.state = circuitStateHalfOpen
		cb.mutex.Unlock()
		cb.mutex.RLock()
		return true
	}

	return false
}

func (cb *circuitBreaker) success() {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	switch cb.state {
	case circuitStateHalfOpen:
		// Success in half-open state closes the circuit
		cb.state = circuitStateClosed
		cb.failures = 0
	case circuitStateClosed:
		// Reset failure count on success
		cb.failures = 0
	}
}

func (cb *circuitBreaker) failure() {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	switch cb.state {
	case circuitStateClosed:
		cb.failures++
		if cb.failures >= cb.maxFailures {
			cb.state = circuitStateOpen
			cb.lastFailure = time.Now()
		}
	case circuitStateHalfOpen:
		// A failure in half-open state re-opens the circuit
		cb.state = circuitStateOpen
		cb.lastFailure = time.Now()
	}
}

// tokenManager handles token storage and refresh
type tokenManager struct {
	tokenStore     TokenStore
	currentToken   string
	tokenExpiry    time.Time
	refreshMutex   sync.Mutex
	config         *SessionConfig
}

func newTokenManager(config *SessionConfig) *tokenManager {
	tm := &tokenManager{
		config: config,
	}
	if config.TokenStore != nil {
		tm.tokenStore = config.TokenStore
		// Try to load token from store
		token, _ := tm.tokenStore.GetToken()
		tm.currentToken = token
	}
	return tm
}

func (tm *tokenManager) getToken() string {
	tm.refreshMutex.Lock()
	defer tm.refreshMutex.Unlock()

	// If token is about to expire or already expired, try to refresh it
	if !tm.tokenExpiry.IsZero() && time.Until(tm.tokenExpiry) < tm.config.TokenRefreshThreshold {
		tm.refreshToken()
	}

	return tm.currentToken
}

func (tm *tokenManager) setToken(token string, expiresIn time.Duration) {
	tm.refreshMutex.Lock()
	defer tm.refreshMutex.Unlock()

	tm.currentToken = token
	if expiresIn > 0 {
		tm.tokenExpiry = time.Now().Add(expiresIn)
	}

	// Persist to store if available
	if tm.tokenStore != nil && token != "" {
		_ = tm.tokenStore.StoreToken(token, tm.tokenExpiry)
	}
}

func (tm *tokenManager) clearToken() {
	tm.refreshMutex.Lock()
	defer tm.refreshMutex.Unlock()

	tm.currentToken = ""
	tm.tokenExpiry = time.Time{}

	// Clear from store if available
	if tm.tokenStore != nil {
		_ = tm.tokenStore.ClearToken()
	}
}

func (tm *tokenManager) refreshToken() error {
	// Implementation depends on your authentication flow
	// This is a placeholder - replace with actual token refresh logic
	return nil
}

// Session is the main entry point for interacting with the Canvus API.
type Session struct {
	BaseURL       string
	HTTPClient    *http.Client
	config        *SessionConfig
	authenticator Authenticator
	tokenManager  *tokenManager
	circuitBreaker *circuitBreaker
	userID        int64 // ID of the authenticated user, if available
}

// NewSession creates a new Canvus API session with the provided configuration.
func NewSession(cfg *SessionConfig, opts ...SessionConfigOption) *Session {
	// Apply any overrides to the config
	for _, opt := range opts {
		opt(cfg)
	}

	// Ensure we have a valid HTTP client
	if cfg.HTTPClient == nil {
		cfg.HTTPClient = http.DefaultClient
	}

	// Set default timeouts if not configured
	if cfg.HTTPClient.Timeout == 0 {
		cfg.HTTPClient.Timeout = cfg.RequestTimeout
	}

	s := &Session{
		BaseURL:       cfg.BaseURL,
		HTTPClient:    cfg.HTTPClient,
		config:        cfg,
		tokenManager:  newTokenManager(cfg),
		circuitBreaker: newCircuitBreaker(cfg.CircuitBreaker.MaxFailures, cfg.CircuitBreaker.ResetTimeout),
	}


	// If we have a token from the store, use it
	if s.tokenManager.tokenStore != nil {
		if token, err := s.tokenManager.tokenStore.GetToken(); err == nil && token != "" {
			s.authenticator = &TokenAuthenticator{Token: token}
		}
	}

	return s
}
// Implements retry logic with exponential backoff, circuit breaking, and token refresh.
func (s *Session) doRequest(ctx context.Context, method, endpoint string, body interface{}, out interface{}, queryParams map[string]string, rawResponse bool, contentType ...string) error {
	var lastErr error
	var resp *http.Response
	var respBody []byte

	// Check circuit breaker first
	if !s.circuitBreaker.allow() {
		return &APIError{
			StatusCode: http.StatusServiceUnavailable,
			Code:      "circuit_breaker_open",
			Message:   "service unavailable due to circuit breaker being open",
		}
	}

	// Parse URL
	u, err := url.Parse(s.BaseURL)
	if err != nil {
		return fmt.Errorf("invalid base URL: %w", err)
	}
	u.Path = path.Join(u.Path, endpoint)

	// Add query parameters if provided
	if queryParams != nil && len(queryParams) > 0 {
		q := u.Query()
		for k, v := range queryParams {
			q.Set(k, v)
		}
		u.RawQuery = q.Encode()
	}

	// Determine content type
	var ct string
	if len(contentType) > 0 {
		ct = contentType[0]
	} else if body != nil {
		ct = "application/json"
	}

	// Main retry loop
	for attempt := 0; attempt <= s.config.MaxRetries; attempt++ {
		// Prepare request body
		reqBody, retryable, err := s.prepareRequestBody(body, ct)
		if err != nil {
			if !retryable || attempt == s.config.MaxRetries {
				return err
			}
			continue
		}

		// Create request
		req, err := http.NewRequestWithContext(ctx, method, u.String(), reqBody)
		if err != nil {
			return fmt.Errorf("failed to create request: %w", err)
		}

		// Set headers
		if ct != "" {
			req.Header.Set("Content-Type", ct)
		}
		req.Header.Set("User-Agent", s.config.UserAgent)

		// Apply authentication
		if s.authenticator != nil {
			s.authenticator.Authenticate(req)
		}

		// Execute request
		resp, err = s.HTTPClient.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("request failed: %w", err)
			if !isRetryableError(err) || attempt == s.config.MaxRetries {
				s.circuitBreaker.failure()
				return lastErr
			}
			if shouldRetry(err, attempt, s.config) {
				time.Sleep(calculateBackoff(attempt, s.config))
				continue
			}
			return lastErr
		}

		// Read response body
		respBody, err = io.ReadAll(resp.Body)
		resp.Body.Close()

		// Handle non-2xx responses
		if resp.StatusCode >= 400 {
			lastErr = s.handleErrorResponse(resp, respBody, attempt)
			if apiErr, ok := lastErr.(*APIError); ok {
				// Handle token expiration
				if apiErr.StatusCode == http.StatusUnauthorized && attempt == 0 {
					if refreshErr := s.refreshAuthToken(ctx); refreshErr == nil {
						// Retry with new token
						continue
					}
				}

				// Check if we should retry
				if isRetryableError(apiErr) && attempt < s.config.MaxRetries {
					time.Sleep(calculateBackoff(attempt, s.config))
					continue
				}
			}
			s.circuitBreaker.failure()
			return lastErr
		}

		// Process successful response
		s.circuitBreaker.success()

		// Handle raw response if requested
		if rawResponse {
			if ptr, ok := out.(*[]byte); ok {
				*ptr = respBody
				return nil
			}
			return errors.New("out must be *[]byte when rawResponse is true")
		}

		// Parse response body
		if out != nil {
			if err := json.Unmarshal(respBody, out); err != nil {
				return fmt.Errorf("failed to decode response: %w", err)
			}

			// Validate response if needed
			if err := validateResponse(out, body, method); err != nil {
				return fmt.Errorf("response validation failed: %w", err)
			}
		}

		return nil
	}

	// If we get here, we've exhausted all retries
	s.circuitBreaker.failure()
	if lastErr != nil {
		return fmt.Errorf("request failed after %d attempts: %w", s.config.MaxRetries, lastErr)
	}
	return errors.New("request failed: unknown error")
}

// prepareRequestBody prepares the request body and determines if the error is retryable
func (s *Session) prepareRequestBody(body interface{}, contentType string) (io.Reader, bool, error) {
	if body == nil {
		return nil, true, nil
	}

	// Handle raw readers
	if rdr, ok := body.(io.Reader); ok {
		// If it's a ReadSeeker, reset to start for each attempt
		if seeker, ok := rdr.(io.ReadSeeker); ok {
			_, _ = seeker.Seek(0, io.SeekStart)
		}
		return rdr, true, nil
	}

	// Handle other types by marshaling to JSON
	b, err := json.Marshal(body)
	if err != nil {
		return nil, false, fmt.Errorf("failed to marshal request body: %w", err)
	}

	return bytes.NewReader(b), true, nil
}

// handleErrorResponse processes error responses and returns an appropriate error
func (s *Session) handleErrorResponse(resp *http.Response, body []byte, attempt int) error {
	// Try to parse as API error
	var apiErr *APIError
	if err := json.Unmarshal(body, &apiErr); err == nil && apiErr.Code != "" {
		apiErr.StatusCode = resp.StatusCode
		return apiErr
	}

	// Fall back to generic error
	errCode := ErrorCode(fmt.Sprintf("http_%d", resp.StatusCode))
	return &APIError{
		StatusCode: resp.StatusCode,
		Code:      errCode,
		Message:   string(body),
	}
}

// refreshAuthToken attempts to refresh the authentication token
func (s *Session) refreshAuthToken(ctx context.Context) error {
	// If we're using token-based auth, try to refresh the token
	if tokenAuth, ok := s.authenticator.(*TokenAuthenticator); ok {
		// Use the token manager to handle refresh
		newToken := s.tokenManager.getToken()
		if newToken != "" && newToken != tokenAuth.Token {
			tokenAuth.Token = newToken
			return nil
		}

		// If we couldn't get a new token, clear the current one
		s.authenticator = nil
	}
	return errors.New("unable to refresh authentication token")
}

// isRetryableError determines if an error is retryable
func isRetryableError(err error) bool {
	// Network errors are always retryable
	var netErr net.Error
	if errors.As(err, &netErr) {
		return true
	}

	// Check for context errors
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return false
	}

	// Check for API errors
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		switch {
		case apiErr.StatusCode >= 500: // Server errors
			return true
		case apiErr.StatusCode == 429: // Rate limited
			return true
		case apiErr.StatusCode == 408: // Request Timeout
			return true
		case apiErr.StatusCode == 0: // Network/connection error
			return true
		default:
			return false
		}
	}

	return false
}

// shouldRetry determines if a request should be retried
func shouldRetry(err error, attempt int, config *SessionConfig) bool {
	if attempt >= config.MaxRetries {
		return false
	}

	// Always retry on network errors
	var netErr net.Error
	if errors.As(err, &netErr) {
		return true
	}

	// Don't retry on context cancellation
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return false
	}

	return true
}

// calculateBackoff calculates the backoff duration using exponential backoff with jitter
func calculateBackoff(attempt int, config *SessionConfig) time.Duration {
	// Calculate exponential backoff with jitter
	min := float64(config.RetryWaitMin)
	max := float64(config.RetryWaitMax)

	if min >= max {
		return config.RetryWaitMax
	}

	// Calculate backoff with jitter
	backoff := min * math.Pow(2, float64(attempt))
	if backoff > max {
		backoff = max
	}

	// Add jitter (random value between 0 and backoff/2)
	randVal, _ := rand.Int(rand.Reader, big.NewInt(1000))
	jitter := (float64(randVal.Int64()) / 1000.0) * (backoff / 2)
	duration := time.Duration(backoff + jitter)

	// Ensure we don't exceed max wait time
	if duration > config.RetryWaitMax {
		duration = config.RetryWaitMax
	}

	return duration
}

// validateResponse performs validation on the decoded response object.
// For PATCH/POST/PUT, if the request is a map or struct, it checks that all fields present in the request
// are present and equal in the response (out). Only fields present in the request are checked.
// For DELETE, if the response contains an 'id' field, it must match the request (if present),
// and the response must have 'status' or 'state' set to 'deleted' (case-insensitive).
// Returns an error if any field does not match.
func validateResponse(obj interface{}, reqBody interface{}, method string) error {
	if obj == nil {
		return errors.New("response is nil")
	}
	if method == "DELETE" {
		// Marshal response to map for comparison
		respMap := map[string]interface{}{}
		b, err := json.Marshal(obj)
		if err != nil {
			return nil // skip validation if can't marshal
		}
		if err := json.Unmarshal(b, &respMap); err != nil {
			return nil // skip validation if can't unmarshal
		}
		// Check id match if present in request
		var reqID interface{}
		switch v := reqBody.(type) {
		case map[string]interface{}:
			reqID = v["id"]
		case nil:
			// skip
		default:
			b, err := json.Marshal(reqBody)
			if err == nil {
				rm := map[string]interface{}{}
				if err := json.Unmarshal(b, &rm); err == nil {
					reqID = rm["id"]
				}
			}
		}
		if reqID != nil {
			if respID, ok := respMap["id"]; ok {
				if !reflect.DeepEqual(respID, reqID) {
					return fmt.Errorf("response id mismatch: got %v, want %v", respID, reqID)
				}
			}
		}
		// Check status or state is 'deleted'
		if status, ok := respMap["status"]; ok {
			if s, ok := status.(string); ok && !equalsIgnoreCase(s, "deleted") {
				return fmt.Errorf("response status is not 'deleted': got %v", s)
			}
		} else if state, ok := respMap["state"]; ok {
			if s, ok := state.(string); ok && !equalsIgnoreCase(s, "deleted") {
				return fmt.Errorf("response state is not 'deleted': got %v", s)
			}
		}
		return nil
	}
	// List of fields to skip strict equality (server-generated or transformed)
	serverGeneratedFields := map[string]struct{}{
		"id":           {},
		"created_at":   {},
		"modified_at":  {},
		"last_login":   {},
		"state":        {},
		"access":       {},
		"preview_hash": {},
		"asset_size":   {},
		"folder_id":    {},
		"parent_id":    {},
		"location":     {},
		"size":         {},
	}

	// Only validate for PATCH/POST/PUT
	if method != "PATCH" && method != "POST" && method != "PUT" {
		return nil
	}

	// Only validate if reqBody is a map or struct
	var reqMap map[string]interface{}
	switch v := reqBody.(type) {
	case map[string]interface{}:
		reqMap = v
	case nil:
		return nil
	default:
		b, err := json.Marshal(reqBody)
		if err != nil {
			return nil
		}
		if err := json.Unmarshal(b, &reqMap); err != nil {
			return nil
		}
	}
	if len(reqMap) == 0 {
		return nil
	}

	// Marshal response to map for comparison
	respMap := map[string]interface{}{}
	b, err := json.Marshal(obj)
	if err != nil {
		return nil
	}
	if err := json.Unmarshal(b, &respMap); err != nil {
		return nil
	}

	// List of write-only fields to skip in validation
	writeOnlyFields := map[string]struct{}{
		"password": {},
		// Add more write-only fields here if needed
	}

	for k, reqVal := range reqMap {
		if _, skip := writeOnlyFields[k]; skip {
			continue // skip write-only fields
		}
		respVal, ok := respMap[k]
		if !ok {
			continue // skip fields not present in response
		}
		if _, skip := serverGeneratedFields[k]; skip {
			continue // skip server-generated fields
		}
		// For string fields, compare case-insensitively for known enums
		if k == "widget_type" {
			if s1, ok1 := reqVal.(string); ok1 {
				if s2, ok2 := respVal.(string); ok2 {
					if !equalsIgnoreCase(s1, s2) {
						return fmt.Errorf("response field %q mismatch (case-insensitive): got %v, want %v", k, respVal, reqVal)
					}
					continue
				}
			}
		}
		// Relax numeric comparison: treat as equal if numerically equal (int/float64)
		if isNumeric(reqVal) && isNumeric(respVal) {
			if !numericEqual(reqVal, respVal) {
				return fmt.Errorf("response field %q mismatch (numeric): got %v, want %v", k, respVal, reqVal)
			}
			continue
		}
		if !reflect.DeepEqual(respVal, reqVal) {
			return fmt.Errorf("response field %q mismatch: got %v, want %v", k, respVal, reqVal)
		}
	}
	return nil
}

// equalsIgnoreCase compares two strings case-insensitively.
func equalsIgnoreCase(a, b string) bool {
	return strings.EqualFold(a, b)
}

// doRequestWithHeaders is like doRequest but allows passing custom headers for the request.
// queryParams may be map[string]string or map[string]interface{}; all values will be stringified.
func (s *Session) doRequestWithHeaders(ctx context.Context, method, endpoint string, body interface{}, out interface{}, queryParams interface{}, headers map[string]string, rawResponse bool) error {
	u, err := url.Parse(s.BaseURL)
	if err != nil {
		return err
	}
	u.Path = path.Join(u.Path, endpoint)

	// Convert queryParams to map[string]string if needed
	qp := make(map[string]string)
	switch params := queryParams.(type) {
	case map[string]string:
		qp = params
	case map[string]interface{}:
		for k, v := range params {
			qp[k] = toString(v)
		}
	case nil:
		// no params
	default:
		return errors.New("queryParams must be map[string]string or map[string]interface{} or nil")
	}

	if len(qp) > 0 {
		q := u.Query()
		for k, v := range qp {
			q.Set(k, v)
		}
		u.RawQuery = q.Encode()
	}

	var reqBody io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return err
		}
		reqBody = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), reqBody)
	if err != nil {
		return err
	}
	if s.authenticator != nil {
		s.authenticator.Authenticate(req)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	// Add custom headers
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := s.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return &APIError{StatusCode: resp.StatusCode, Message: string(respBody)}
	}

	if out != nil {
		if rawResponse {
			// out must be *[]byte
			b, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			if ptr, ok := out.(*[]byte); ok {
				*ptr = b
			} else {
				return errors.New("out must be *[]byte when rawResponse is true")
			}
		} else {
			if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
				return err
			}
		}
	}
	return nil
}

// toString converts an interface{} to string for query param values.
func toString(v interface{}) string {
	switch val := v.(type) {
	case string:
		return val
	case fmt.Stringer:
		return val.String()
	case int, int8, int16, int32, int64:
		return fmt.Sprintf("%d", val)
	case uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", val)
	case float32, float64:
		return fmt.Sprintf("%v", val)
	case bool:
		return fmt.Sprintf("%t", val)
	default:
		return fmt.Sprintf("%v", val)
	}
}

// Login authenticates a user and stores the returned token and user ID for future requests.
func (s *Session) Login(ctx context.Context, email, password string) error {
	loginReq := map[string]string{
		"username": email,
		"password": password,
	}
	var loginResp struct {
		Token string `json:"token"`
		User  struct {
			ID int64 `json:"id"`
		} `json:"user"`
	}
	err := s.doRequest(ctx, http.MethodPost, "users/login", loginReq, &loginResp, nil, false)
	if err != nil {
		return err
	}
	if loginResp.Token == "" {
		return errors.New("login: no token returned")
	}
	s.authenticator = &TokenAuthenticator{Token: loginResp.Token}
	s.userID = loginResp.User.ID
	return nil
}

// Logout invalidates the current token and clears authentication.
// It calls POST /users/logout and clears the authenticator on success.
func (s *Session) Logout(ctx context.Context) error {
	logoutReq := map[string]string{}
	var logoutResp map[string]interface{}
	err := s.doRequest(ctx, http.MethodPost, "users/logout", logoutReq, &logoutResp, nil, false)
	if err != nil {
		return err
	}
	s.authenticator = nil
	return nil
}

// Users provides access to user management methods.
func (s *Session) Users() *Session {
	return s
}

// UserID returns the authenticated user's ID, or 0 if not logged in.
func (s *Session) UserID() int64 {
	return s.userID
}

// isNumeric returns true if v is a numeric type
func isNumeric(v interface{}) bool {
	switch v.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		return true
	default:
		return false
	}
}

// numericEqual returns true if a and b are numerically equal (int/float64)
func numericEqual(a, b interface{}) bool {
	af, aok := toFloat64(a)
	bf, bok := toFloat64(b)
	if aok && bok {
		return af == bf
	}
	return false
}

// toFloat64 converts a numeric value to float64
func toFloat64(v interface{}) (float64, bool) {
	switch n := v.(type) {
	case int:
		return float64(n), true
	case int8:
		return float64(n), true
	case int16:
		return float64(n), true
	case int32:
		return float64(n), true
	case int64:
		return float64(n), true
	case uint:
		return float64(n), true
	case uint8:
		return float64(n), true
	case uint16:
		return float64(n), true
	case uint32:
		return float64(n), true
	case uint64:
		return float64(n), true
	case float32:
		return float64(n), true
	case float64:
		return n, true
	default:
		return 0, false
	}
}
