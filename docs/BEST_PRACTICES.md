# Best Practices for the Canvus Go SDK

This guide covers recommended patterns and practices for building robust applications with the Canvus Go SDK.

## Table of Contents

- [Error Handling Patterns](#error-handling-patterns)
- [Authentication Patterns](#authentication-patterns)
- [Session Lifecycle Management](#session-lifecycle-management)
- [Concurrency and Goroutine Safety](#concurrency-and-goroutine-safety)
- [Context Usage](#context-usage)
- [Rate Limiting and Pagination](#rate-limiting-and-pagination)
- [Memory Efficiency](#memory-efficiency)
- [Logging and Debugging](#logging-and-debugging)
- [Resource Cleanup](#resource-cleanup)

---

## Error Handling Patterns

### Use Type Assertions for API Errors

The SDK returns `*canvus.APIError` for API-specific errors. Use type assertions to access detailed error information:

```go
result, err := session.GetCanvas(ctx, canvasID)
if err != nil {
    if apiErr, ok := err.(*canvus.APIError); ok {
        // Handle API-specific errors
        log.Printf("API error: status=%d, code=%s, message=%s",
            apiErr.StatusCode, apiErr.Code, apiErr.Message)

        // Access additional details
        if apiErr.RequestID != "" {
            log.Printf("Request ID: %s", apiErr.RequestID)
        }
        if len(apiErr.Details) > 0 {
            log.Printf("Details: %v", apiErr.Details)
        }
    } else {
        // Handle network/other errors
        log.Printf("Unexpected error: %v", err)
    }
    return err
}
```

### Use Error Codes for Control Flow

Switch on error codes to handle different scenarios:

```go
canvas, err := session.GetCanvas(ctx, canvasID)
if err != nil {
    if apiErr, ok := err.(*canvus.APIError); ok {
        switch apiErr.Code {
        case canvus.ErrNotFound:
            // Canvas doesn't exist - maybe create it
            return createDefaultCanvas(ctx, session)
        case canvus.ErrUnauthorized:
            // Invalid credentials - re-authenticate
            return session.Login(ctx, email, password)
        case canvus.ErrForbidden:
            // No permission - inform user
            return fmt.Errorf("you don't have access to this canvas")
        case canvus.ErrTooManyRequests:
            // Rate limited - wait and retry
            time.Sleep(time.Second)
            return session.GetCanvas(ctx, canvasID)
        default:
            return err
        }
    }
    return err
}
```

### Use errors.Is and errors.As for Error Chains

The SDK supports Go's error wrapping conventions:

```go
if apiErr, ok := err.(*canvus.APIError); ok {
    // Check if this specific error matches
    if errors.Is(err, &canvus.APIError{StatusCode: 404}) {
        // Handle not found
    }
}

// Or use errors.As for type assertions
var apiErr *canvus.APIError
if errors.As(err, &apiErr) {
    log.Printf("API error: %s", apiErr.Message)
}
```

### Check for Retryable Errors

Use the `IsRetryableError` helper function:

```go
if canvus.IsRetryableError(err) {
    // This error can be safely retried
    time.Sleep(time.Second)
    return retry(operation)
}
```

---

## Authentication Patterns

### API Key Authentication (Recommended)

For most use cases, use API key authentication. It's simpler and more secure:

```go
config := canvus.DefaultSessionConfig()
config.BaseURL = os.Getenv("CANVUS_URL") + "/api/v1"

session := canvus.NewSession(config, canvus.WithAPIKey(os.Getenv("CANVUS_API_KEY")))
```

**Security Considerations:**
- Never hardcode API keys in source code
- Use environment variables or secure secret management
- Rotate API keys periodically
- Use the minimum required permissions

### Token-Based Authentication

For user-facing applications where you need to authenticate end users:

```go
session := canvus.NewSession(config)

err := session.Login(ctx, email, password)
if err != nil {
    return fmt.Errorf("authentication failed: %w", err)
}

// Store the user ID if needed
userID := session.UserID()

// Remember to logout
defer session.Logout(ctx)
```

### Secure Credential Storage

Never store credentials in plain text:

```go
// BAD - Don't do this
apiKey := "hardcoded-api-key"

// GOOD - Use environment variables
apiKey := os.Getenv("CANVUS_API_KEY")

// BETTER - Use a secret manager
apiKey, err := secretManager.GetSecret("canvus-api-key")
if err != nil {
    log.Fatal("Failed to retrieve API key")
}
```

---

## Session Lifecycle Management

### Session Creation

Create sessions with appropriate configuration:

```go
func newCanvusSession() *canvus.Session {
    config := canvus.DefaultSessionConfig()
    config.BaseURL = os.Getenv("CANVUS_URL") + "/api/v1"

    // Configure timeouts based on your needs
    config.RequestTimeout = 30 * time.Second
    config.MaxRetries = 3
    config.RetryWaitMin = 100 * time.Millisecond
    config.RetryWaitMax = 2 * time.Second

    // Configure circuit breaker
    config.CircuitBreaker.MaxFailures = 5
    config.CircuitBreaker.ResetTimeout = 30 * time.Second

    return canvus.NewSession(config, canvus.WithAPIKey(os.Getenv("CANVUS_API_KEY")))
}
```

### Session Reuse

Reuse sessions across requests - don't create a new session for each operation:

```go
// BAD - Creating new session for each request
func getCanvas(canvasID string) (*canvus.Canvas, error) {
    session := canvus.NewSession(config, canvus.WithAPIKey(apiKey)) // Don't do this!
    return session.GetCanvas(ctx, canvasID)
}

// GOOD - Reuse session
type CanvusService struct {
    session *canvus.Session
}

func NewCanvusService() *CanvusService {
    return &CanvusService{
        session: canvus.NewSession(config, canvus.WithAPIKey(apiKey)),
    }
}

func (s *CanvusService) GetCanvas(ctx context.Context, canvasID string) (*canvus.Canvas, error) {
    return s.session.GetCanvas(ctx, canvasID)
}
```

### Session Cleanup

Always clean up sessions when done, especially with token authentication:

```go
session := canvus.NewSession(config)
err := session.Login(ctx, email, password)
if err != nil {
    return err
}

// Ensure logout on function exit
defer func() {
    if logoutErr := session.Logout(ctx); logoutErr != nil {
        log.Printf("Warning: logout failed: %v", logoutErr)
    }
}()

// Use the session...
```

---

## Concurrency and Goroutine Safety

### Session Safety

The `Session` struct is safe for concurrent use from multiple goroutines:

```go
session := canvus.NewSession(config, canvus.WithAPIKey(apiKey))

var wg sync.WaitGroup
for _, canvasID := range canvasIDs {
    wg.Add(1)
    go func(id string) {
        defer wg.Done()
        canvas, err := session.GetCanvas(ctx, id)
        if err != nil {
            log.Printf("Error getting canvas %s: %v", id, err)
            return
        }
        processCanvas(canvas)
    }(canvasID)
}
wg.Wait()
```

### Use BatchProcessor for Bulk Operations

For concurrent operations with controlled parallelism, use the `BatchProcessor`:

```go
// Configure batch processing
batchConfig := canvus.DefaultBatchConfig()
batchConfig.MaxConcurrency = 10   // Limit concurrent operations
batchConfig.ContinueOnError = true // Continue on individual failures
batchConfig.RetryAttempts = 3      // Retry failed operations
batchConfig.ProgressCallback = func(completed, total int, results []*canvus.BatchResult) {
    log.Printf("Progress: %d/%d", completed, total)
}

processor := canvus.NewBatchProcessor(session, batchConfig)

// Build operations
builder := canvus.NewBatchOperationBuilder()
for _, widget := range widgets {
    builder.Delete(widget.ID, widget)
}

// Execute batch
results, err := processor.ExecuteBatch(ctx, builder.Build())
if err != nil {
    log.Printf("Batch failed: %v", err)
}

// Check results
summary := canvus.Summarize(results)
log.Printf("Completed: %d/%d successful", summary.Successful, summary.TotalOperations)
for _, failed := range summary.FailedOperations {
    log.Printf("Failed operation %s: %v", failed.OperationID, failed.Error)
}
```

### Avoid Race Conditions

When sharing results across goroutines, use proper synchronization:

```go
var (
    mu      sync.Mutex
    results = make(map[string]*canvus.Canvas)
)

var wg sync.WaitGroup
for _, id := range canvasIDs {
    wg.Add(1)
    go func(canvasID string) {
        defer wg.Done()
        canvas, err := session.GetCanvas(ctx, canvasID)
        if err != nil {
            return
        }

        mu.Lock()
        results[canvasID] = canvas
        mu.Unlock()
    }(id)
}
wg.Wait()
```

---

## Context Usage

### Always Pass Context

Every SDK method accepts a `context.Context` - always pass one:

```go
// GOOD - Use context
ctx := context.Background()
canvas, err := session.GetCanvas(ctx, canvasID)

// BAD - Don't create your own context inside functions
func getCanvas(session *canvus.Session, id string) (*canvus.Canvas, error) {
    ctx := context.Background() // Don't do this!
    return session.GetCanvas(ctx, id)
}

// GOOD - Accept context as parameter
func getCanvas(ctx context.Context, session *canvus.Session, id string) (*canvus.Canvas, error) {
    return session.GetCanvas(ctx, id)
}
```

### Use Timeouts

Set appropriate timeouts for operations:

```go
// Timeout for individual operation
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

canvas, err := session.GetCanvas(ctx, canvasID)
if err != nil {
    if errors.Is(err, context.DeadlineExceeded) {
        return nil, fmt.Errorf("operation timed out")
    }
    return nil, err
}
```

### Support Cancellation

Respect cancellation signals from parent contexts:

```go
func processCanvases(ctx context.Context, session *canvus.Session, ids []string) error {
    for _, id := range ids {
        // Check if context is cancelled
        select {
        case <-ctx.Done():
            return ctx.Err()
        default:
        }

        canvas, err := session.GetCanvas(ctx, id)
        if err != nil {
            return err
        }

        if err := processCanvas(canvas); err != nil {
            return err
        }
    }
    return nil
}
```

### Use Deadlines for Long Operations

For operations that might take a while, set deadlines:

```go
// 5 minute deadline for batch import
ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(5*time.Minute))
defer cancel()

err := importCanvasData(ctx, session, data)
```

---

## Rate Limiting and Pagination

### Handle Rate Limits Gracefully

The SDK automatically retries on 429 errors, but you can also handle them explicitly:

```go
for {
    result, err := session.CreateWidget(ctx, canvasID, widget)
    if err != nil {
        if apiErr, ok := err.(*canvus.APIError); ok && apiErr.Code == canvus.ErrTooManyRequests {
            // Wait and retry
            log.Println("Rate limited, waiting...")
            time.Sleep(5 * time.Second)
            continue
        }
        return err
    }
    return nil
}
```

### Implement Rate Limiting for Bulk Operations

When processing many items, add delays:

```go
const requestsPerSecond = 10
limiter := time.NewTicker(time.Second / requestsPerSecond)
defer limiter.Stop()

for _, item := range items {
    <-limiter.C // Wait for rate limiter

    _, err := session.CreateWidget(ctx, canvasID, item)
    if err != nil {
        return err
    }
}
```

### Use Pagination for Large Lists

The SDK supports pagination through `ListOptions`:

```go
func getAllWidgets(ctx context.Context, session *canvus.Session, canvasID string) ([]*canvus.Widget, error) {
    var allWidgets []*canvus.Widget
    offset := 0
    limit := 100 // Fetch 100 at a time

    for {
        widgets, err := session.ListWidgets(ctx, canvasID, &canvus.ListOptions{
            Limit:  limit,
            Offset: offset,
        })
        if err != nil {
            return nil, err
        }

        allWidgets = append(allWidgets, widgets...)

        // Check if we've fetched all items
        if len(widgets) < limit {
            break
        }

        offset += limit
    }

    return allWidgets, nil
}
```

---

## Memory Efficiency

### Process Large Datasets in Batches

Don't load everything into memory at once:

```go
// BAD - Loads all widgets into memory
allWidgets, _ := session.ListWidgets(ctx, canvasID, nil)
for _, widget := range allWidgets {
    processWidget(widget)
}

// GOOD - Process in batches
func processWidgetsInBatches(ctx context.Context, session *canvus.Session, canvasID string) error {
    offset := 0
    batchSize := 100

    for {
        widgets, err := session.ListWidgets(ctx, canvasID, &canvus.ListOptions{
            Limit:  batchSize,
            Offset: offset,
        })
        if err != nil {
            return err
        }

        for _, widget := range widgets {
            if err := processWidget(widget); err != nil {
                return err
            }
        }

        if len(widgets) < batchSize {
            break
        }
        offset += batchSize
    }

    return nil
}
```

### Release Resources Early

Release large objects when no longer needed:

```go
data, err := session.ExportCanvas(ctx, canvasID)
if err != nil {
    return err
}

// Process the data
result := processExportData(data)

// Release the large data object
data = nil

// Continue with result...
```

### Use Streaming for Large Files

For large assets, consider streaming instead of loading into memory:

```go
// For large file uploads, use readers
file, err := os.Open("large-file.pdf")
if err != nil {
    return err
}
defer file.Close()

// The SDK will stream the file content
asset, err := session.CreateAsset(ctx, canvasID, file, "application/pdf")
```

---

## Logging and Debugging

### Custom HTTP Transport for Debugging

Use a custom `http.RoundTripper` to log requests and responses:

```go
type loggingTransport struct {
    transport http.RoundTripper
    logger    *log.Logger
}

func (t *loggingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
    // Log request
    t.logger.Printf("-> %s %s", req.Method, req.URL)
    for key, values := range req.Header {
        if key == "Private-Token" {
            t.logger.Printf("   %s: [REDACTED]", key)
        } else {
            t.logger.Printf("   %s: %s", key, values)
        }
    }

    // Execute request
    resp, err := t.transport.RoundTrip(req)
    if err != nil {
        t.logger.Printf("<- Error: %v", err)
        return nil, err
    }

    // Log response
    t.logger.Printf("<- %s %s", resp.Status, req.URL.Path)

    return resp, nil
}

// Use it
config := canvus.DefaultSessionConfig()
config.HTTPClient = &http.Client{
    Transport: &loggingTransport{
        transport: http.DefaultTransport,
        logger:    log.New(os.Stderr, "[HTTP] ", log.LstdFlags),
    },
}
```

### Structured Logging

Use structured logging for production:

```go
type CanvusService struct {
    session *canvus.Session
    logger  *slog.Logger
}

func (s *CanvusService) GetCanvas(ctx context.Context, canvasID string) (*canvus.Canvas, error) {
    s.logger.Info("getting canvas", "canvas_id", canvasID)

    canvas, err := s.session.GetCanvas(ctx, canvasID)
    if err != nil {
        s.logger.Error("failed to get canvas",
            "canvas_id", canvasID,
            "error", err,
        )
        return nil, err
    }

    s.logger.Info("got canvas",
        "canvas_id", canvasID,
        "canvas_name", canvas.Name,
    )

    return canvas, nil
}
```

### Include Request IDs in Logs

Always log the request ID for debugging:

```go
canvas, err := session.GetCanvas(ctx, canvasID)
if err != nil {
    if apiErr, ok := err.(*canvus.APIError); ok {
        log.Printf("Error getting canvas: %v (request_id: %s)",
            apiErr.Message, apiErr.RequestID)
    }
    return err
}
```

---

## Resource Cleanup

### Use Defer for Cleanup

Always clean up resources using defer:

```go
func processWithSession(ctx context.Context) error {
    session := canvus.NewSession(config)

    err := session.Login(ctx, email, password)
    if err != nil {
        return err
    }
    defer session.Logout(ctx)

    // Use the session...
    return nil
}
```

### Clean Up Test Resources

In tests, ensure created resources are cleaned up:

```go
func TestCreateCanvas(t *testing.T) {
    ctx := context.Background()

    // Create canvas
    canvas, err := session.CreateCanvas(ctx, canvus.CreateCanvasRequest{
        Name: "Test Canvas",
    })
    require.NoError(t, err)

    // Ensure cleanup
    defer func() {
        err := session.DeleteCanvas(ctx, canvas.ID)
        if err != nil {
            t.Logf("Failed to cleanup canvas: %v", err)
        }
    }()

    // Run tests...
}
```

### Handle Cleanup Errors

Log cleanup errors but don't let them mask the original error:

```go
func doOperation(ctx context.Context) (err error) {
    resource, err := createResource(ctx)
    if err != nil {
        return err
    }

    defer func() {
        if cleanupErr := deleteResource(ctx, resource); cleanupErr != nil {
            if err == nil {
                err = cleanupErr
            } else {
                // Log but don't mask original error
                log.Printf("cleanup error: %v", cleanupErr)
            }
        }
    }()

    // Use resource...
    return nil
}
```

---

## Summary

Following these best practices will help you build reliable, maintainable, and efficient applications with the Canvus Go SDK:

1. **Always handle errors properly** - Use type assertions to access `*canvus.APIError` details
2. **Secure your credentials** - Use environment variables or secret managers
3. **Reuse sessions** - Don't create new sessions for each request
4. **Use context for cancellation** - Pass context through all operations
5. **Handle concurrency safely** - Use the BatchProcessor for bulk operations
6. **Paginate large results** - Don't load everything into memory
7. **Clean up resources** - Always logout and delete temporary resources
8. **Log effectively** - Include request IDs and use structured logging

For more information, see the [Troubleshooting Guide](TROUBLESHOOTING.md) and [API Reference](API_REFERENCE.md).
