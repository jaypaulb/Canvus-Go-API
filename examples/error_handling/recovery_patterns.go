// Package main demonstrates advanced error handling patterns with the Canvus Go SDK.
//
// This example shows:
// - Typed error handling with *canvus.APIError
// - Error code checking for specific error conditions
// - Custom retry patterns beyond built-in retries
// - Circuit breaker pattern implementation
// - Graceful degradation strategies
// - Error recovery and cleanup patterns
//
// To run this example:
//   export CANVUS_API_URL="https://your-canvus-server.example.com/api/public/v1"
//   export CANVUS_API_KEY="your-api-key-here"
//   go run recovery_patterns.go
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/jaypaulb/Canvus-Go-API/canvus"
)

func main() {
	// Step 1: Get configuration from environment variables
	apiURL := os.Getenv("CANVUS_API_URL")
	apiKey := os.Getenv("CANVUS_API_KEY")

	if apiURL == "" {
		log.Fatal("CANVUS_API_URL environment variable is required")
	}
	if apiKey == "" {
		log.Fatal("CANVUS_API_KEY environment variable is required")
	}

	fmt.Println("Canvus Go SDK - Error Handling and Recovery Patterns")
	fmt.Println("====================================================")
	fmt.Printf("Server: %s\n\n", apiURL)

	// Step 2: Create session with API key authentication
	cfg := canvus.DefaultSessionConfig()
	cfg.BaseURL = apiURL
	session := canvus.NewSession(cfg, canvus.WithAPIKey(apiKey))

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Step 3: Demonstrate basic error type checking
	fmt.Println("Pattern 1: Basic Error Type Checking")
	fmt.Println("=====================================")

	// Try to get a non-existent canvas
	_, err := session.GetCanvas(ctx, "non-existent-canvas-id")
	if err != nil {
		// Check if it's a typed API error
		var apiErr *canvus.APIError
		if errors.As(err, &apiErr) {
			fmt.Printf("API Error detected:\n")
			fmt.Printf("  Status Code: %d\n", apiErr.StatusCode)
			fmt.Printf("  Error Code: %s\n", apiErr.Code)
			fmt.Printf("  Message: %s\n", apiErr.Message)
			if apiErr.RequestID != "" {
				fmt.Printf("  Request ID: %s\n", apiErr.RequestID)
			}
		} else {
			// It's some other type of error (network, timeout, etc.)
			fmt.Printf("Non-API error: %v\n", err)
		}
	}
	fmt.Println()

	// Step 4: Error code checking patterns
	fmt.Println("Pattern 2: Error Code Checking")
	fmt.Println("==============================")

	demonstrateErrorCodeChecking(ctx, session)
	fmt.Println()

	// Step 5: Custom retry patterns
	fmt.Println("Pattern 3: Custom Retry Patterns")
	fmt.Println("================================")

	demonstrateCustomRetry(ctx, session)
	fmt.Println()

	// Step 6: Circuit breaker pattern
	fmt.Println("Pattern 4: Circuit Breaker Pattern")
	fmt.Println("==================================")

	demonstrateCircuitBreaker(ctx, session)
	fmt.Println()

	// Step 7: Graceful degradation
	fmt.Println("Pattern 5: Graceful Degradation")
	fmt.Println("===============================")

	demonstrateGracefulDegradation(ctx, session)
	fmt.Println()

	// Step 8: Error recovery and cleanup
	fmt.Println("Pattern 6: Error Recovery and Cleanup")
	fmt.Println("=====================================")

	demonstrateErrorRecovery(ctx, session)
	fmt.Println()

	// Step 9: Using SDK error utilities
	fmt.Println("Pattern 7: SDK Error Utilities")
	fmt.Println("==============================")

	demonstrateErrorUtilities(ctx, session)
	fmt.Println()

	// Summary
	fmt.Println("=" + repeatStr("=", 59))
	fmt.Println("Error Handling Best Practices Summary")
	fmt.Println("=" + repeatStr("=", 59))
	fmt.Println()
	fmt.Println("1. Always use errors.As() to check for *canvus.APIError")
	fmt.Println("2. Check error codes for specific handling (401, 403, 404, 429)")
	fmt.Println("3. Use canvus.IsRetryableError() to determine retry eligibility")
	fmt.Println("4. Implement circuit breakers for external service protection")
	fmt.Println("5. Use graceful degradation for non-critical operations")
	fmt.Println("6. Always clean up resources in defer statements")
	fmt.Println("7. Log request IDs for debugging with server logs")
	fmt.Println("8. Handle context cancellation separately from API errors")

	fmt.Println("\nError handling example completed successfully!")
}

// demonstrateErrorCodeChecking shows how to check for specific error codes
func demonstrateErrorCodeChecking(ctx context.Context, session *canvus.Session) {
	// Try an operation that will likely fail
	_, err := session.GetCanvas(ctx, "invalid-id")
	if err == nil {
		fmt.Println("Unexpected success - no error to demonstrate")
		return
	}

	var apiErr *canvus.APIError
	if !errors.As(err, &apiErr) {
		fmt.Printf("Non-API error: %v\n", err)
		return
	}

	// Switch on status code for specific handling
	switch apiErr.StatusCode {
	case 400:
		fmt.Println("Bad Request: Check input parameters")
		fmt.Printf("  Details: %s\n", apiErr.Message)

	case 401:
		fmt.Println("Unauthorized: Invalid or expired credentials")
		fmt.Println("  Action: Re-authenticate or refresh token")

	case 403:
		fmt.Println("Forbidden: Insufficient permissions")
		fmt.Println("  Action: Check user roles and canvas access")

	case 404:
		fmt.Println("Not Found: Resource does not exist")
		fmt.Println("  Action: Verify ID or handle missing resource gracefully")

	case 429:
		fmt.Println("Rate Limited: Too many requests")
		fmt.Println("  Action: Implement backoff and reduce request rate")

	case 500, 502, 503:
		fmt.Println("Server Error: Temporary server issue")
		fmt.Println("  Action: Retry with exponential backoff")

	default:
		fmt.Printf("Other error (status %d): %s\n", apiErr.StatusCode, apiErr.Message)
	}

	// Check error code using the ErrorCode constants
	fmt.Println("\nChecking error code constants:")
	switch apiErr.Code {
	case canvus.ErrNotFound:
		fmt.Println("  Error code: ErrNotFound")
	case canvus.ErrUnauthorized:
		fmt.Println("  Error code: ErrUnauthorized")
	case canvus.ErrForbidden:
		fmt.Println("  Error code: ErrForbidden")
	case canvus.ErrTooManyRequests:
		fmt.Println("  Error code: ErrTooManyRequests")
	default:
		fmt.Printf("  Error code: %s\n", apiErr.Code)
	}
}

// demonstrateCustomRetry shows custom retry patterns beyond built-in retries
func demonstrateCustomRetry(ctx context.Context, session *canvus.Session) {
	fmt.Println("Implementing custom retry with exponential backoff:")

	maxRetries := 3
	baseDelay := 500 * time.Millisecond

	var lastErr error
	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			// Calculate exponential backoff
			delay := baseDelay * time.Duration(1<<(attempt-1))
			fmt.Printf("  Retry attempt %d after %v delay...\n", attempt, delay)
			time.Sleep(delay)
		}

		// Try the operation
		_, err := session.GetCanvas(ctx, "retry-test-id")
		if err == nil {
			fmt.Println("  Success!")
			return
		}

		lastErr = err

		// Check if error is retryable
		if !canvus.IsRetryableError(err) {
			fmt.Printf("  Non-retryable error: %v\n", err)
			break
		}

		fmt.Printf("  Attempt %d failed: %v\n", attempt+1, err)
	}

	fmt.Printf("  All retries exhausted. Last error: %v\n", lastErr)

	// Custom retry with context awareness
	fmt.Println("\nCustom retry with context awareness:")
	fmt.Println("```go")
	fmt.Println("func retryWithContext(ctx context.Context, maxRetries int, fn func() error) error {")
	fmt.Println("    for attempt := 0; attempt <= maxRetries; attempt++ {")
	fmt.Println("        err := fn()")
	fmt.Println("        if err == nil {")
	fmt.Println("            return nil")
	fmt.Println("        }")
	fmt.Println("        if ctx.Err() != nil {")
	fmt.Println("            return ctx.Err() // Context cancelled")
	fmt.Println("        }")
	fmt.Println("        if !canvus.IsRetryableError(err) {")
	fmt.Println("            return err // Non-retryable")
	fmt.Println("        }")
	fmt.Println("        time.Sleep(backoff(attempt))")
	fmt.Println("    }")
	fmt.Println("    return errors.New(\"max retries exceeded\")")
	fmt.Println("}")
	fmt.Println("```")
}

// CircuitBreaker implements a simple circuit breaker pattern
type CircuitBreaker struct {
	maxFailures    int
	resetTimeout   time.Duration
	failures       int
	lastFailure    time.Time
	state          string // "closed", "open", "half-open"
	mu             sync.Mutex
}

// NewCircuitBreaker creates a new circuit breaker
func NewCircuitBreaker(maxFailures int, resetTimeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		maxFailures:  maxFailures,
		resetTimeout: resetTimeout,
		state:        "closed",
	}
}

// Execute runs a function through the circuit breaker
func (cb *CircuitBreaker) Execute(fn func() error) error {
	cb.mu.Lock()

	// Check if circuit is open
	if cb.state == "open" {
		if time.Since(cb.lastFailure) > cb.resetTimeout {
			cb.state = "half-open"
		} else {
			cb.mu.Unlock()
			return errors.New("circuit breaker is open")
		}
	}

	cb.mu.Unlock()

	// Execute the function
	err := fn()

	cb.mu.Lock()
	defer cb.mu.Unlock()

	if err != nil {
		cb.failures++
		cb.lastFailure = time.Now()

		if cb.state == "half-open" || cb.failures >= cb.maxFailures {
			cb.state = "open"
		}
		return err
	}

	// Success - reset on half-open, decrement failures otherwise
	if cb.state == "half-open" {
		cb.state = "closed"
		cb.failures = 0
	} else if cb.failures > 0 {
		cb.failures--
	}

	return nil
}

// State returns the current state of the circuit breaker
func (cb *CircuitBreaker) State() string {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	return cb.state
}

// demonstrateCircuitBreaker shows circuit breaker pattern implementation
func demonstrateCircuitBreaker(ctx context.Context, session *canvus.Session) {
	// Create a circuit breaker
	// Opens after 3 failures, resets after 5 seconds
	cb := NewCircuitBreaker(3, 5*time.Second)

	fmt.Println("Circuit breaker configured: 3 failures to open, 5s reset timeout")
	fmt.Println()

	// Simulate multiple failures
	for i := 1; i <= 5; i++ {
		err := cb.Execute(func() error {
			_, err := session.GetCanvas(ctx, "circuit-breaker-test")
			return err
		})

		fmt.Printf("Request %d: state=%s, error=%v\n", i, cb.State(), err != nil)

		if cb.State() == "open" {
			fmt.Println("  Circuit is OPEN - requests are being rejected")
		}
	}

	fmt.Println("\nCircuit breaker pattern code example:")
	fmt.Println("```go")
	fmt.Println("cb := NewCircuitBreaker(5, 30*time.Second)")
	fmt.Println("err := cb.Execute(func() error {")
	fmt.Println("    return session.GetCanvas(ctx, canvasID)")
	fmt.Println("})")
	fmt.Println("if err != nil && cb.State() == \"open\" {")
	fmt.Println("    // Use cached data or show degraded experience")
	fmt.Println("}")
	fmt.Println("```")
}

// demonstrateGracefulDegradation shows graceful degradation strategies
func demonstrateGracefulDegradation(ctx context.Context, session *canvus.Session) {
	fmt.Println("Strategy 1: Fallback to cached data")
	fmt.Println("-----------------------------------")

	// Simulate cached data
	cachedCanvases := []canvus.Canvas{
		{ID: "cached-1", Name: "Cached Canvas 1"},
		{ID: "cached-2", Name: "Cached Canvas 2"},
	}

	// Try to get fresh data
	canvases, err := session.ListCanvases(ctx, nil)
	if err != nil {
		fmt.Printf("Failed to get fresh data: %v\n", err)
		fmt.Printf("Using cached data (%d canvases)\n", len(cachedCanvases))
		canvases = cachedCanvases
	}
	fmt.Printf("Returning %d canvases\n\n", len(canvases))

	fmt.Println("Strategy 2: Partial results")
	fmt.Println("---------------------------")
	fmt.Println("```go")
	fmt.Println("results := make([]Widget, 0)")
	fmt.Println("for _, id := range widgetIDs {")
	fmt.Println("    widget, err := session.GetWidget(ctx, canvasID, id)")
	fmt.Println("    if err != nil {")
	fmt.Println("        log.Printf(\"Failed to get widget %s: %v\", id, err)")
	fmt.Println("        continue // Skip failed widgets, return partial results")
	fmt.Println("    }")
	fmt.Println("    results = append(results, widget)")
	fmt.Println("}")
	fmt.Println("```")
	fmt.Println()

	fmt.Println("Strategy 3: Feature flags for degradation")
	fmt.Println("-----------------------------------------")
	fmt.Println("```go")
	fmt.Println("func getCanvasWithFeatures(ctx context.Context, id string) (*Canvas, error) {")
	fmt.Println("    canvas, err := session.GetCanvas(ctx, id)")
	fmt.Println("    if err != nil {")
	fmt.Println("        return nil, err")
	fmt.Println("    }")
	fmt.Println("    ")
	fmt.Println("    // Try to get widgets, but don't fail if unavailable")
	fmt.Println("    widgets, err := session.ListWidgets(ctx, id, nil)")
	fmt.Println("    if err != nil {")
	fmt.Println("        log.Printf(\"Widgets unavailable: %v\", err)")
	fmt.Println("        // Return canvas without widgets")
	fmt.Println("        return canvas, nil")
	fmt.Println("    }")
	fmt.Println("    ")
	fmt.Println("    canvas.Widgets = widgets")
	fmt.Println("    return canvas, nil")
	fmt.Println("}")
	fmt.Println("```")
}

// demonstrateErrorRecovery shows error recovery and cleanup patterns
func demonstrateErrorRecovery(ctx context.Context, session *canvus.Session) {
	fmt.Println("Pattern: Transactional cleanup on error")
	fmt.Println("---------------------------------------")

	// This demonstrates creating multiple resources and cleaning up on failure
	fmt.Println("```go")
	fmt.Println("func createCanvasWithWidgets(ctx context.Context, name string) (*Canvas, error) {")
	fmt.Println("    // Track resources for cleanup")
	fmt.Println("    var createdWidgets []string")
	fmt.Println("    var canvas *Canvas")
	fmt.Println("    ")
	fmt.Println("    // Cleanup function for rollback")
	fmt.Println("    cleanup := func() {")
	fmt.Println("        if canvas != nil {")
	fmt.Println("            _ = session.DeleteCanvas(ctx, canvas.ID)")
	fmt.Println("        }")
	fmt.Println("    }")
	fmt.Println("    ")
	fmt.Println("    // Create canvas")
	fmt.Println("    canvas, err := session.CreateCanvas(ctx, req)")
	fmt.Println("    if err != nil {")
	fmt.Println("        return nil, err")
	fmt.Println("    }")
	fmt.Println("    ")
	fmt.Println("    // Create widgets")
	fmt.Println("    for i := 0; i < 10; i++ {")
	fmt.Println("        widget, err := session.CreateWidget(ctx, canvas.ID, data)")
	fmt.Println("        if err != nil {")
	fmt.Println("            cleanup() // Rollback on error")
	fmt.Println("            return nil, fmt.Errorf(\"failed at widget %d: %w\", i, err)")
	fmt.Println("        }")
	fmt.Println("        createdWidgets = append(createdWidgets, widget.ID)")
	fmt.Println("    }")
	fmt.Println("    ")
	fmt.Println("    return canvas, nil")
	fmt.Println("}")
	fmt.Println("```")
	fmt.Println()

	fmt.Println("Pattern: Defer cleanup with error capture")
	fmt.Println("-----------------------------------------")
	fmt.Println("```go")
	fmt.Println("func processCanvas(ctx context.Context, id string) (err error) {")
	fmt.Println("    // Acquire resources")
	fmt.Println("    tempCanvas, err := session.CopyCanvas(ctx, id, req)")
	fmt.Println("    if err != nil {")
	fmt.Println("        return err")
	fmt.Println("    }")
	fmt.Println("    ")
	fmt.Println("    // Cleanup using named return for error capture")
	fmt.Println("    defer func() {")
	fmt.Println("        cleanupErr := session.DeleteCanvas(ctx, tempCanvas.ID)")
	fmt.Println("        if err == nil && cleanupErr != nil {")
	fmt.Println("            err = fmt.Errorf(\"cleanup failed: %w\", cleanupErr)")
	fmt.Println("        }")
	fmt.Println("    }()")
	fmt.Println("    ")
	fmt.Println("    // Process...")
	fmt.Println("    return processWidgets(ctx, tempCanvas.ID)")
	fmt.Println("}")
	fmt.Println("```")
}

// demonstrateErrorUtilities shows SDK error utilities
func demonstrateErrorUtilities(ctx context.Context, session *canvus.Session) {
	// Test IsRetryableError
	fmt.Println("Using canvus.IsRetryableError():")

	testErrors := []error{
		canvus.NewAPIError(500, canvus.ErrInternalServer, "Server error"),
		canvus.NewAPIError(429, canvus.ErrTooManyRequests, "Rate limited"),
		canvus.NewAPIError(404, canvus.ErrNotFound, "Not found"),
		canvus.NewAPIError(401, canvus.ErrUnauthorized, "Unauthorized"),
		context.Canceled,
		context.DeadlineExceeded,
	}

	for _, err := range testErrors {
		retryable := canvus.IsRetryableError(err)
		fmt.Printf("  %v -> retryable: %v\n", err, retryable)
	}
	fmt.Println()

	// Test IsContextError
	fmt.Println("Using canvus.IsContextError():")

	for _, err := range testErrors {
		contextErr := canvus.IsContextError(err)
		if contextErr {
			fmt.Printf("  %v -> is context error\n", err)
		}
	}
	fmt.Println()

	// Test errors.Is with APIError
	fmt.Println("Using errors.Is() with APIError:")

	err := canvus.NewAPIError(404, canvus.ErrNotFound, "Canvas not found")
	targetErr := &canvus.APIError{StatusCode: 404}

	if errors.Is(err, targetErr) {
		fmt.Printf("  Error matches 404 status code\n")
	}

	targetErr2 := &canvus.APIError{Code: canvus.ErrNotFound}
	if errors.Is(err, targetErr2) {
		fmt.Printf("  Error matches ErrNotFound code\n")
	}
	fmt.Println()

	// Test WrapError
	fmt.Println("Using canvus.WrapError():")
	originalErr := errors.New("original error")
	wrappedErr := canvus.WrapError(originalErr, "failed to process canvas")
	fmt.Printf("  Original: %v\n", originalErr)
	fmt.Printf("  Wrapped: %v\n", wrappedErr)
	fmt.Println()

	// Test WrapErrorf
	fmt.Println("Using canvus.WrapErrorf():")
	canvasID := "canvas-123"
	wrappedErrf := canvus.WrapErrorf(originalErr, "failed to process canvas %s", canvasID)
	fmt.Printf("  Formatted: %v\n", wrappedErrf)
}

// Helper function to repeat a string
func repeatStr(s string, n int) string {
	result := ""
	for i := 0; i < n; i++ {
		result += s
	}
	return result
}
