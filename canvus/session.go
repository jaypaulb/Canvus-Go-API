// Package canvus provides a Go SDK for the Canvus API.
package canvus

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"reflect"
	"strings"
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
func WithAPIKey(apiKey string) SessionOption {
	return func(s *Session) {
		s.authenticator = &APIKeyAuthenticator{Header: "Private-Token", APIKey: apiKey}
	}
}

// WithToken configures the session to use a bearer token.
func WithToken(token string) SessionOption {
	return func(s *Session) {
		s.authenticator = &TokenAuthenticator{Token: token}
	}
}

// Session is the main entry point for interacting with the Canvus API.
type Session struct {
	BaseURL       string
	HTTPClient    *http.Client
	authenticator Authenticator
	userID        int64 // ID of the authenticated user, if available
}

// NewSession creates a new Canvus API session.
// If httpClient is nil, http.DefaultClient is used.
func NewSession(baseURL string, opts ...SessionOption) *Session {
	s := &Session{
		BaseURL:    baseURL,
		HTTPClient: http.DefaultClient,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

// doRequest is a helper for making HTTP requests to the Canvus API.
// Adds centralized response validation and retry logic:
// - Validates the response body (if out != nil) using validateResponse.
// - Retries up to 3 times on transient errors (5xx, network errors, validation failures) with exponential backoff.
func (s *Session) doRequest(ctx context.Context, method, endpoint string, body interface{}, out interface{}, queryParams map[string]string, rawResponse bool, contentType ...string) error {
	u, err := url.Parse(s.BaseURL)
	if err != nil {
		return err
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

	var reqBody io.Reader
	var ct string
	if body != nil {
		if rdr, ok := body.(io.Reader); ok {
			reqBody = rdr
			if len(contentType) > 0 {
				ct = contentType[0]
			} else {
				ct = "application/octet-stream"
			}
		} else {
			b, err := json.Marshal(body)
			if err != nil {
				return err
			}
			reqBody = bytes.NewReader(b)
			ct = "application/json"
		}
	}

	var lastErr error
	maxRetries := 3
	for attempt := 0; attempt < maxRetries; attempt++ {
		req, err := http.NewRequestWithContext(ctx, method, u.String(), reqBody)
		if err != nil {
			return err
		}
		if s.authenticator != nil {
			s.authenticator.Authenticate(req)
		}
		if body != nil {
			req.Header.Set("Content-Type", ct)
		}

		// DEBUG: Log headers and URL for troubleshooting
		fmt.Printf("[DEBUG] %s %s\n", method, u.String())
		for k, v := range req.Header {
			fmt.Printf("[DEBUG] Header: %s: %v\n", k, v)
		}

		resp, err := s.HTTPClient.Do(req)
		if err != nil {
			lastErr = err
			if attempt < maxRetries-1 {
				time.Sleep(time.Duration(1<<attempt) * 100 * time.Millisecond)
				continue
			}
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			b, _ := io.ReadAll(resp.Body)
			lastErr = &APIError{StatusCode: resp.StatusCode, Message: string(b)}
			// Retry on 5xx
			if resp.StatusCode >= 500 && attempt < maxRetries-1 {
				time.Sleep(time.Duration(1<<attempt) * 100 * time.Millisecond)
				continue
			}
			return lastErr
		}

		if out != nil {
			if rawResponse {
				b, err := io.ReadAll(resp.Body)
				if err != nil {
					lastErr = err
					if attempt < maxRetries-1 {
						time.Sleep(time.Duration(1<<attempt) * 100 * time.Millisecond)
						continue
					}
					return err
				}
				if ptr, ok := out.(*[]byte); ok {
					*ptr = b
				} else {
					return errors.New("out must be *[]byte when rawResponse is true")
				}
			} else {
				if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
					lastErr = err
					if attempt < maxRetries-1 {
						time.Sleep(time.Duration(1<<attempt) * 100 * time.Millisecond)
						continue
					}
					return err
				}
				// Centralized response validation
				if err := validateResponse(out, body, method); err != nil {
					lastErr = err
					if attempt < maxRetries-1 {
						time.Sleep(time.Duration(1<<attempt) * 100 * time.Millisecond)
						continue
					}
					return fmt.Errorf("response validation failed: %w", err)
				}
			}
		}
		return nil
	}
	return lastErr
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
		// Try to marshal and unmarshal to map[string]interface{}
		b, err := json.Marshal(reqBody)
		if err != nil {
			return nil // skip validation if can't marshal
		}
		if err := json.Unmarshal(b, &reqMap); err != nil {
			return nil // skip validation if can't unmarshal
		}
	}
	if len(reqMap) == 0 {
		return nil
	}
	// Marshal response to map for comparison
	respMap := map[string]interface{}{}
	b, err := json.Marshal(obj)
	if err != nil {
		return nil // skip validation if can't marshal
	}
	if err := json.Unmarshal(b, &respMap); err != nil {
		return nil // skip validation if can't unmarshal
	}
	for k, v := range reqMap {
		if respVal, ok := respMap[k]; ok {
			if !reflect.DeepEqual(respVal, v) {
				return fmt.Errorf("response field %q mismatch: got %v, want %v", k, respVal, v)
			}
		} else {
			return fmt.Errorf("response missing field %q", k)
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

	// DEBUG: Log headers and URL for troubleshooting
	fmt.Printf("[DEBUG] %s %s\n", method, u.String())
	for k, v := range req.Header {
		fmt.Printf("[DEBUG] Header: %s: %v\n", k, v)
	}

	resp, err := s.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return &APIError{StatusCode: resp.StatusCode, Message: string(b)}
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
		"email":    email,
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
