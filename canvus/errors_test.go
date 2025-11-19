package canvus

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAPIError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *APIError
		expected string
	}{
		{
			name:     "basic error",
			err:      NewAPIError(400, ErrInvalidRequest, "invalid request"),
			expected: "API error 400 (invalid_request): invalid request",
		},
		{
			name: "with request ID",
			err:  NewAPIError(404, ErrNotFound, "not found").WithRequestID("req-123"),
			expected: "API error 404 (not_found): not found (request_id: req-123)",
		},
		{
			name: "with wrapped error",
			err:  NewAPIError(500, ErrInternalServer, "internal error").Wrap(errors.New("underlying error")),
			expected: "API error 500 (internal_server_error): internal error: underlying error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.err.Error())
		})
	}
}

func TestErrorFromStatus(t *testing.T) {
	tests := []struct {
		statusCode int
		expected   ErrorCode
	}{
		{http.StatusBadRequest, ErrInvalidRequest},
		{http.StatusUnauthorized, ErrUnauthorized},
		{http.StatusForbidden, ErrForbidden},
		{http.StatusNotFound, ErrNotFound},
		{http.StatusConflict, ErrConflict},
		{http.StatusTooManyRequests, ErrTooManyRequests},
		{http.StatusInternalServerError, ErrInternalServer},
		{http.StatusNotImplemented, ErrNotImplemented},
		{http.StatusServiceUnavailable, ErrServiceUnavailable},
		{http.StatusTeapot, ErrInvalidRequest}, // 4xx
		{http.StatusBadGateway, ErrInternalServer}, // 5xx
	}

	for _, tt := range tests {
		t.Run(http.StatusText(tt.statusCode), func(t *testing.T) {
			err := ErrorFromStatus(tt.statusCode, "test error")
			apiErr, ok := err.(*APIError)
			require.True(t, ok, "expected *APIError")
			assert.Equal(t, tt.statusCode, apiErr.StatusCode)
			assert.Equal(t, tt.expected, apiErr.Code)
		})
	}
}

func TestIsContextError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "context canceled",
			err:      context.Canceled,
			expected: true,
		},
		{
			name:     "context deadline exceeded",
			err:      context.DeadlineExceeded,
			expected: true,
		},
		{
			name:     "custom timeout error",
			err:      NewAPIError(0, ErrTimeout, "timeout"),
			expected: true,
		},
		{
			name:     "other error",
			err:      errors.New("some error"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, IsContextError(tt.err))
		})
	}
}

func TestIsRetryableError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "retryable server error",
			err:      NewAPIError(500, ErrInternalServer, "internal error"),
			expected: true,
		},
		{
			name:     "retryable rate limit",
			err:      NewAPIError(429, ErrTooManyRequests, "too many requests"),
			expected: true,
		},
		{
			name:     "not retryable client error",
			err:      NewAPIError(400, ErrInvalidRequest, "bad request"),
			expected: false,
		},
		{
			name:     "context canceled",
			err:      context.Canceled,
			expected: false,
		},
		{
			name:     "generic error",
			err:      errors.New("some error"),
			expected: true, // Generic errors are retryable by default
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, IsRetryableError(tt.err))
		})
	}
}

func TestParseErrorResponse(t *testing.T) {
	tests := []struct {
		name     string
		status   int
		body     string
		expected *APIError
	}{
		{
			name:   "standard error response",
			status: http.StatusBadRequest,
			body:   `{"error": "invalid_request", "error_description": "missing required field"}`,
			expected: &APIError{
				StatusCode: http.StatusBadRequest,
				Code:       "",
				Message:    "missing required field",
			},
		},
		{
			name:   "with request ID",
			status: http.StatusNotFound,
			body:   `{"error": "not_found", "error_description": "resource not found", "request_id": "req-123"}`,
			expected: &APIError{
				StatusCode: http.StatusNotFound,
				Code:       "",
				Message:    "resource not found",
				RequestID:  "req-123",
			},
		},
		{
			name:   "invalid JSON",
			status: http.StatusInternalServerError,
			body:   "internal server error",
			expected: &APIError{
				StatusCode: http.StatusInternalServerError,
				Message:    "internal server error",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ParseErrorResponse(tt.status, []byte(tt.body))
			assert.Equal(t, tt.expected.StatusCode, err.StatusCode)
			assert.Equal(t, tt.expected.Code, err.Code)
			assert.Equal(t, tt.expected.Message, err.Message)
			assert.Equal(t, tt.expected.RequestID, err.RequestID)
		})
	}
}

func TestValidationErrors(t *testing.T) {
	t.Run("empty validation errors", func(t *testing.T) {
		var errs ValidationErrors
		assert.False(t, errs.HasErrors())
		assert.Equal(t, "no validation errors", errs.Error())
	})

	t.Run("with validation errors", func(t *testing.T) {
		var errs ValidationErrors
		errs.Add("name", "is required")
		errs.Add("email", "must be a valid email")

		true := true
		assert.True(t, true, "expected validation errors")

		errStr := errs.Error()
		assert.Contains(t, errStr, "2 validation errors:")
		assert.Contains(t, errStr, "name: is required")
		assert.Contains(t, errStr, "email: must be a valid email")
	})
}

func TestWrapError(t *testing.T) {
	t.Run("wrap nil error", func(t *testing.T) {
		err := WrapError(nil, "some context")
		assert.NoError(t, err)
	})

	t.Run("wrap with empty message", func(t *testing.T) {
		original := errors.New("original error")
		err := WrapError(original, "")
		assert.Equal(t, original, err)
	})

	t.Run("wrap with context", func(t *testing.T) {
		original := errors.New("original error")
		err := WrapError(original, "failed to process request")
		assert.EqualError(t, err, "failed to process request: original error")
		assert.True(t, errors.Is(err, original))
	})
}

func TestWrapErrorf(t *testing.T) {
	t.Run("format error message", func(t *testing.T) {
		original := errors.New("connection refused")
		err := WrapErrorf(original, "failed to connect to %s:%d", "example.com", 8080)
		assert.EqualError(t, err, "failed to connect to example.com:8080: connection refused")
	})
}
