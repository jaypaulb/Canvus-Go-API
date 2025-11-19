// Package canvus provides option types for SDK methods (pagination, filtering, etc.).
package canvus

import (
	"net/http"
	"time"
)

// SessionConfig holds configuration for the API session.
type SessionConfig struct {
	// BaseURL is the base URL for all API requests.
	BaseURL string

	// HTTPClient is the HTTP client to use for requests.
	// If nil, http.DefaultClient is used.
	HTTPClient *http.Client

	// MaxRetries is the maximum number of retries for failed requests.
	// Default: 3
	MaxRetries int

	// RetryWaitMin is the minimum time to wait between retries.
	// Default: 100ms
	RetryWaitMin time.Duration

	// RetryWaitMax is the maximum time to wait between retries.
	// Default: 1s
	RetryWaitMax time.Duration

	// RequestTimeout is the timeout for each HTTP request.
	// Default: 30s
	RequestTimeout time.Duration

	// UserAgent is the User-Agent header to send with requests.
	// Default: "canvus-go-sdk/<version>"
	UserAgent string

	// TokenRefreshThreshold is the duration before token expiration when a refresh should be attempted.
	// Default: 5 minutes
	TokenRefreshThreshold time.Duration

	// CircuitBreaker configures the circuit breaker behavior.
	CircuitBreaker CircuitBreakerConfig

	// TokenStore is used to store and retrieve authentication tokens.
	// If nil, tokens are not persisted between sessions.
	TokenStore TokenStore
}

// CircuitBreakerConfig holds configuration for the circuit breaker.
type CircuitBreakerConfig struct {
	// MaxFailures is the number of consecutive failures before opening the circuit.
	// Default: 5
	MaxFailures int

	// ResetTimeout is the time after which an open circuit will attempt to close.
	// Default: 30s
	ResetTimeout time.Duration
}

// TokenStore defines the interface for storing and retrieving authentication tokens.
type TokenStore interface {
	// GetToken returns the stored token or an error if not found.
	GetToken() (string, error)

	// StoreToken stores the token.
	StoreToken(token string, expiresAt time.Time) error

	// ClearToken removes the stored token.
	ClearToken() error
}

// DefaultSessionConfig returns a default session configuration.
func DefaultSessionConfig() *SessionConfig {
	return &SessionConfig{
		MaxRetries:           3,
		RetryWaitMin:         100 * time.Millisecond,
		RetryWaitMax:         time.Second,
		RequestTimeout:       30 * time.Second,
		UserAgent:            "canvus-go-sdk/v1.0.0",
		TokenRefreshThreshold: 5 * time.Minute,
		CircuitBreaker: CircuitBreakerConfig{
			MaxFailures:  5,
			ResetTimeout: 30 * time.Second,
		},
	}
}

// SessionConfigOption is a function that configures a SessionConfig.
type SessionConfigOption func(*SessionConfig)

// WithHTTPClient sets the HTTP client for the session.
func WithHTTPClient(client *http.Client) SessionConfigOption {
	return func(c *SessionConfig) {
		c.HTTPClient = client
	}
}

// WithMaxRetries sets the maximum number of retries for failed requests.
func WithMaxRetries(maxRetries int) SessionConfigOption {
	return func(c *SessionConfig) {
		c.MaxRetries = maxRetries
	}
}

// WithRetryWait sets the minimum and maximum wait time between retries.
func WithRetryWait(min, max time.Duration) SessionConfigOption {
	return func(c *SessionConfig) {
		c.RetryWaitMin = min
		c.RetryWaitMax = max
	}
}

// WithRequestTimeout sets the timeout for each HTTP request.
func WithRequestTimeout(timeout time.Duration) SessionConfigOption {
	return func(c *SessionConfig) {
		c.RequestTimeout = timeout
	}
}

// WithUserAgent sets the User-Agent header for requests.
func WithUserAgent(ua string) SessionConfigOption {
	return func(c *SessionConfig) {
		c.UserAgent = ua
	}
}

// WithTokenStore sets the token store for the session.
func WithTokenStore(store TokenStore) SessionConfigOption {
	return func(c *SessionConfig) {
		c.TokenStore = store
	}
}

// WithCircuitBreaker sets the circuit breaker configuration.
func WithCircuitBreaker(maxFailures int, resetTimeout time.Duration) SessionConfigOption {
	return func(c *SessionConfig) {
		c.CircuitBreaker = CircuitBreakerConfig{
			MaxFailures:  maxFailures,
			ResetTimeout: resetTimeout,
		}
	}
}

// WithTokenRefreshThreshold sets the token refresh threshold.
func WithTokenRefreshThreshold(threshold time.Duration) SessionConfigOption {
	return func(c *SessionConfig) {
		c.TokenRefreshThreshold = threshold
	}
}

// ListOptions specifies options for list endpoints (pagination, filtering, etc.).
type ListOptions struct {
	Limit  int    // Maximum number of items to return
	Offset int    // Offset for pagination
	Filter string // Optional filter string
}

// GetOptions specifies options for get endpoints (e.g., subscribe to updates).
type GetOptions struct {
	Subscribe bool // Whether to subscribe to updates (if supported)
}

// SubscribeOptions specifies options for streaming/subscription endpoints.
type SubscribeOptions struct {
	Annotations bool // Whether to include annotations
}

// AuditLogOptions specifies options for querying the audit log.
type AuditLogOptions struct {
	Page    int    // Page number
	PerPage int    // Items per page
	Filter  string // Filter string
}

// RequestOption is an option for API requests.
type RequestOption func(*requestOptions)

type requestOptions struct {
	retryable    bool
	headers      map[string]string
	queryParams  map[string]string
	contentType  string
	noAuth       bool
	expectedCode int
}

// WithRetryable sets whether the request is retryable.
// Default: true
func WithRetryable(retryable bool) RequestOption {
	return func(opts *requestOptions) {
		opts.retryable = retryable
	}
}

// WithHeader sets a request header.
func WithHeader(key, value string) RequestOption {
	return func(opts *requestOptions) {
		if opts.headers == nil {
			opts.headers = make(map[string]string)
		}
		opts.headers[key] = value
	}
}

// WithQueryParam sets a query parameter.
func WithQueryParam(key, value string) RequestOption {
	return func(opts *requestOptions) {
		if opts.queryParams == nil {
			opts.queryParams = make(map[string]string)
		}
		opts.queryParams[key] = value
	}
}

// WithContentType sets the Content-Type header.
func WithContentType(contentType string) RequestOption {
	return func(opts *requestOptions) {
		opts.contentType = contentType
	}
}

// WithoutAuth disables authentication for the request.
func WithoutAuth() RequestOption {
	return func(opts *requestOptions) {
		opts.noAuth = true
	}
}

// WithExpectedCode sets the expected HTTP status code.
// If the response status code doesn't match, an error is returned.
func WithExpectedCode(code int) RequestOption {
	return func(opts *requestOptions) {
		opts.expectedCode = code
	}
}

// defaultRequestOptions returns the default request options.
func defaultRequestOptions() *requestOptions {
	return &requestOptions{
		retryable:    true,
		headers:      make(map[string]string),
		queryParams:  make(map[string]string),
		contentType:  "application/json",
		expectedCode: 0, // No specific code expected
	}
}
