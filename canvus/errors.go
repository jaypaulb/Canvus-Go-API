// Package canvus provides error types and utilities for the Canvus SDK.
package canvus

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// ErrorCode represents a machine-readable error code.
type ErrorCode string

// Common error codes.
const (
	// Client errors (4xx)
	ErrInvalidRequest   ErrorCode = "invalid_request"        // 400
	ErrUnauthorized    ErrorCode = "unauthorized"          // 401
	ErrForbidden       ErrorCode = "forbidden"             // 403
	ErrNotFound        ErrorCode = "not_found"             // 404
	ErrConflict        ErrorCode = "conflict"              // 409
	ErrTooManyRequests ErrorCode = "too_many_requests"     // 429

	// Server errors (5xx)
	ErrInternalServer  ErrorCode = "internal_server_error" // 500
	ErrNotImplemented  ErrorCode = "not_implemented"       // 501
	ErrServiceUnavailable ErrorCode = "service_unavailable" // 503

	// SDK errors
	ErrValidation     ErrorCode = "validation_error"
	ErrRateLimited    ErrorCode = "rate_limited"
	ErrTimeout        ErrorCode = "timeout"
	ErrNetwork        ErrorCode = "network_error"
	ErrUnexpected     ErrorCode = "unexpected_error"
)

// APIError represents an error returned by the Canvus API.
type APIError struct {
	// StatusCode is the HTTP status code from the API response.
	StatusCode int `json:"status_code"`
	
	// Code is a machine-readable error code.
	Code ErrorCode `json:"code"`
	
	// Message is a human-readable error message.
	Message string `json:"message"`
	
	// RequestID is a unique identifier for the request, if available.
	RequestID string `json:"request_id,omitempty"`
	
	// Details contains additional error details, if any.
	Details map[string]interface{} `json:"details,omitempty"`
	
	// Wrapped is the underlying error that triggered this one, if any.
	Wrapped error `json:"-"`
}

// Error implements the error interface.
func (e *APIError) Error() string {
	if e == nil {
		return ""
	}
	
	var b strings.Builder
	fmt.Fprintf(&b, "API error %d", e.StatusCode)
	
	if e.Code != "" {
		fmt.Fprintf(&b, " (%s)", e.Code)
	}
	
	if e.Message != "" {
		if e.Code != "" {
			b.WriteString(": ")
		} else {
			b.WriteString(" ")
		}
		b.WriteString(e.Message)
	}
	
	if e.RequestID != "" {
		fmt.Fprintf(&b, " (request_id: %s)", e.RequestID)
	}
	
	if e.Wrapped != nil {
		fmt.Fprintf(&b, ": %v", e.Wrapped)
	}
	
	return b.String()
}

// Unwrap returns the underlying error, if any.
func (e *APIError) Unwrap() error {
	return e.Wrapped
}

// Is reports whether this error matches the target error.
func (e *APIError) Is(target error) bool {
	t, ok := target.(*APIError)
	if !ok {
		return false
	}
	
	// If the target has a status code, it must match
	if t.StatusCode != 0 && e.StatusCode != t.StatusCode {
		return false
	}
	
	// If the target has a code, it must match
	if t.Code != "" && e.Code != t.Code {
		return false
	}
	
	return true
}

// WithDetails adds additional details to the error.
func (e *APIError) WithDetails(details map[string]interface{}) *APIError {
	e.Details = details
	return e
}

// WithRequestID sets the request ID on the error.
func (e *APIError) WithRequestID(requestID string) *APIError {
	e.RequestID = requestID
	return e
}

// Wrap returns a new error that wraps the current error with additional context.
func (e *APIError) Wrap(err error) *APIError {
	return &APIError{
		StatusCode: e.StatusCode,
		Code:       e.Code,
		Message:    e.Message,
		RequestID:  e.RequestID,
		Details:    e.Details,
		Wrapped:    err,
	}
}

// NewAPIError creates a new APIError with the given status code and message.
func NewAPIError(statusCode int, code ErrorCode, message string) *APIError {
	return &APIError{
		StatusCode: statusCode,
		Code:       code,
		Message:    message,
	}
}

// ErrorFromStatus creates an appropriate error based on the HTTP status code.
func ErrorFromStatus(statusCode int, message string) error {
	switch statusCode {
	case http.StatusBadRequest:
		return NewAPIError(statusCode, ErrInvalidRequest, message)
	case http.StatusUnauthorized:
		return NewAPIError(statusCode, ErrUnauthorized, message)
	case http.StatusForbidden:
		return NewAPIError(statusCode, ErrForbidden, message)
	case http.StatusNotFound:
		return NewAPIError(statusCode, ErrNotFound, message)
	case http.StatusConflict:
		return NewAPIError(statusCode, ErrConflict, message)
	case http.StatusTooManyRequests:
		return NewAPIError(statusCode, ErrTooManyRequests, message)
	case http.StatusInternalServerError:
		return NewAPIError(statusCode, ErrInternalServer, message)
	case http.StatusNotImplemented:
		return NewAPIError(statusCode, ErrNotImplemented, message)
	case http.StatusServiceUnavailable:
		return NewAPIError(statusCode, ErrServiceUnavailable, message)
	default:
		if statusCode >= 400 && statusCode < 500 {
			return NewAPIError(statusCode, ErrInvalidRequest, message)
		} else if statusCode >= 500 {
			return NewAPIError(statusCode, ErrInternalServer, message)
		}
		return errors.New(message)
	}
}

// IsContextError checks if the error is a context-related error.
func IsContextError(err error) bool {
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return true
	}
	
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.Code == ErrTimeout
	}
	
	return false
}

// IsRetryableError checks if the error is retryable.
func IsRetryableError(err error) bool {
	if err == nil {
		return false
	}
	
	// Context errors are not retryable
	if IsContextError(err) {
		return false
	}
	
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		// Non-API errors are considered retryable by default
		return true
	}
	
	// Retry on server errors and rate limits
	switch apiErr.Code {
	case ErrTooManyRequests, ErrServiceUnavailable, ErrInternalServer:
		return true
	}
	
	// Retry on 5xx errors
	if apiErr.StatusCode >= 500 {
		return true
	}
	
	return false
}

// ErrorResponse represents a standard error response from the API.
type ErrorResponse struct {
	Error            string                 `json:"error,omitempty"`
	ErrorDescription string                 `json:"error_description,omitempty"`
	ErrorURI         string                 `json:"error_uri,omitempty"`
	RequestID        string                 `json:"request_id,omitempty"`
	Details          map[string]interface{} `json:"details,omitempty"`
}

// ParseErrorResponse parses an error response from the API.
func ParseErrorResponse(statusCode int, body []byte) *APIError {
	var resp ErrorResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return NewAPIError(statusCode, "", string(body))
	}
	
	err := NewAPIError(statusCode, "", resp.ErrorDescription)
	if resp.RequestID != "" {
		err.RequestID = resp.RequestID
	}
	if len(resp.Details) > 0 {
		err.Details = resp.Details
	}
	
	return err
}

// ValidationError represents a validation error.
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// Error implements the error interface.
func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error: %s: %s", e.Field, e.Message)
}

// ValidationErrors is a collection of validation errors.
type ValidationErrors []*ValidationError

// Error implements the error interface.
func (e ValidationErrors) Error() string {
	if len(e) == 0 {
		return "no validation errors"
	}
	
	var b strings.Builder
	fmt.Fprintf(&b, "%d validation errors:", len(e))
	for _, err := range e {
		fmt.Fprintf(&b, "\n- %s: %s", err.Field, err.Message)
	}
	return b.String()
}

// Add adds a new validation error.
func (e *ValidationErrors) Add(field, message string) {
	*e = append(*e, &ValidationError{
		Field:   field,
		Message: message,
	})
}

// HasErrors returns true if there are any validation errors.
func (e ValidationErrors) HasErrors() bool {
	return len(e) > 0
}

// WrapError wraps an error with additional context.
// If err is nil, returns nil.
// If msg is empty, returns err as-is.
func WrapError(err error, msg string) error {
	if err == nil {
		return nil
	}
	if msg == "" {
		return err
	}
	return fmt.Errorf("%s: %w", msg, err)
}

// WrapErrorf wraps an error with additional formatted context.
// If err is nil, returns nil.
// If format is empty, returns err as-is.
func WrapErrorf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	if format == "" {
		return err
	}
	return fmt.Errorf("%s: %w", fmt.Sprintf(format, args...), err)
}
