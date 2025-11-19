# Troubleshooting Guide

This guide covers common issues you may encounter when using the Canvus Go SDK and provides solutions for each.

## Table of Contents

- [Top 10 Common Issues](#top-10-common-issues)
  1. [Authentication Failures](#1-authentication-failures)
  2. [Connection Refused Errors](#2-connection-refused-errors)
  3. [Network Timeouts](#3-network-timeouts)
  4. [Certificate Verification Failures](#4-certificate-verification-failures)
  5. [Rate Limiting Responses](#5-rate-limiting-responses)
  6. [Permission Denied Errors](#6-permission-denied-errors)
  7. [Invalid Request Errors](#7-invalid-request-errors)
  8. [Resource Not Found Errors](#8-resource-not-found-errors)
  9. [Batch Operation Failures](#9-batch-operation-failures)
  10. [Import/Export Issues](#10-importexport-issues)
- [Debugging with Custom RoundTripper](#debugging-with-custom-roundtripper)
- [Performance Tuning Tips](#performance-tuning-tips)
- [When to Contact Support](#when-to-contact-support)

---

## Top 10 Common Issues

### 1. Authentication Failures

**Symptoms:**
- `401 Unauthorized` error
- Error code: `unauthorized`
- Message: "Invalid token" or "Authentication required"

**Common Causes:**

**Invalid or expired API key:**
```go
// Check that your API key is correct
apiKey := os.Getenv("CANVUS_API_KEY")
if apiKey == "" {
    log.Fatal("CANVUS_API_KEY environment variable not set")
}

// Ensure you're using the correct option
session := canvus.NewSession(config, canvus.WithAPIKey(apiKey))
```

**Expired login token:**
```go
// Tokens from Login() expire after a period
// Re-authenticate when you get 401 errors
err := session.Login(ctx, email, password)
if err != nil {
    if apiErr, ok := err.(*canvus.APIError); ok && apiErr.StatusCode == 401 {
        // Token expired, try logging in again
        err = session.Login(ctx, email, password)
    }
}
```

**Missing authentication:**
```go
// Ensure you've configured authentication
config := canvus.DefaultSessionConfig()
config.BaseURL = "https://server/api/v1"

// BAD - No authentication configured
session := canvus.NewSession(config)

// GOOD - Authentication configured
session := canvus.NewSession(config, canvus.WithAPIKey(apiKey))
```

**Solutions:**
1. Verify your API key is correct and not expired
2. Check that the API key has the required permissions
3. Ensure the `Private-Token` header is being sent
4. Try generating a new API key in Canvus

---

### 2. Connection Refused Errors

**Symptoms:**
- `connection refused` error
- `dial tcp: connect: connection refused`

**Common Causes:**

**Incorrect server URL:**
```go
// Ensure the URL is correct
config.BaseURL = "https://your-canvus-server/api/v1"  // Include /api/v1

// Common mistakes:
// - Wrong port
// - Missing protocol (https://)
// - Typo in hostname
// - Missing /api/v1 path
```

**Server not running:**
```bash
# Verify the server is accessible
curl -I https://your-canvus-server/api/v1/license
```

**Firewall blocking connection:**
```bash
# Test network connectivity
telnet your-canvus-server 443
```

**Solutions:**
1. Verify the server URL and port
2. Ensure the Canvus server is running
3. Check firewall rules
4. Verify network connectivity to the server
5. Check if VPN is required

---

### 3. Network Timeouts

**Symptoms:**
- `context deadline exceeded`
- `net/http: request canceled (Client.Timeout exceeded)`
- Operations hanging

**Common Causes:**

**Timeout too short:**
```go
// Increase request timeout
config := canvus.DefaultSessionConfig()
config.RequestTimeout = 60 * time.Second  // Increase from default 30s

// Or use context timeout
ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
defer cancel()
```

**Slow network connection:**
```go
// Increase retry settings
config.MaxRetries = 5
config.RetryWaitMax = 5 * time.Second
```

**Server overloaded:**
```go
// Add backoff for server load
if apiErr, ok := err.(*canvus.APIError); ok {
    if apiErr.StatusCode >= 500 {
        time.Sleep(5 * time.Second)
        // Retry...
    }
}
```

**Solutions:**
1. Increase `RequestTimeout` in session config
2. Use longer context timeouts for large operations
3. Increase retry settings
4. Check server health
5. Test network latency to server

---

### 4. Certificate Verification Failures

**Symptoms:**
- `x509: certificate signed by unknown authority`
- `x509: certificate has expired`
- TLS handshake errors

**Common Causes:**

**Self-signed certificates (development):**
```go
import (
    "crypto/tls"
    "net/http"
)

// WARNING: Only use in development!
config := canvus.DefaultSessionConfig()
config.HTTPClient = &http.Client{
    Transport: &http.Transport{
        TLSClientConfig: &tls.Config{
            InsecureSkipVerify: true,  // Disables certificate verification
        },
    },
}
```

**Custom CA certificate:**
```go
import (
    "crypto/tls"
    "crypto/x509"
    "io/ioutil"
    "net/http"
)

// Load custom CA certificate
caCert, err := ioutil.ReadFile("/path/to/ca-cert.pem")
if err != nil {
    log.Fatal(err)
}

caCertPool := x509.NewCertPool()
caCertPool.AppendCertsFromPEM(caCert)

config.HTTPClient = &http.Client{
    Transport: &http.Transport{
        TLSClientConfig: &tls.Config{
            RootCAs: caCertPool,
        },
    },
}
```

**Solutions:**
1. For development: temporarily skip verification (not recommended for production)
2. Install the CA certificate on your system
3. Load custom CA certificates in your application
4. Check that server certificates are valid and not expired

---

### 5. Rate Limiting Responses

**Symptoms:**
- `429 Too Many Requests` error
- Error code: `too_many_requests`
- Message: "Rate limit exceeded"

**Common Causes:**

**Too many requests in short time:**
```go
// The SDK automatically retries on 429, but you can also handle it
for {
    result, err := session.CreateWidget(ctx, canvasID, widget)
    if err != nil {
        if apiErr, ok := err.(*canvus.APIError); ok {
            if apiErr.Code == canvus.ErrTooManyRequests {
                // Wait before retrying
                time.Sleep(5 * time.Second)
                continue
            }
        }
        return err
    }
    break
}
```

**Implement rate limiting:**
```go
// Add rate limiting to your application
limiter := time.NewTicker(100 * time.Millisecond)  // 10 requests/second
defer limiter.Stop()

for _, item := range items {
    <-limiter.C  // Wait for rate limiter
    _, err := session.CreateWidget(ctx, canvasID, item)
    if err != nil {
        return err
    }
}
```

**Use batch operations:**
```go
// Use BatchProcessor for bulk operations
batchConfig := canvus.DefaultBatchConfig()
batchConfig.MaxConcurrency = 5  // Limit concurrent requests

processor := canvus.NewBatchProcessor(session, batchConfig)
```

**Solutions:**
1. Implement client-side rate limiting
2. Use batch operations for bulk updates
3. Add delays between requests
4. Increase retry wait times
5. Contact support if legitimate use case needs higher limits

---

### 6. Permission Denied Errors

**Symptoms:**
- `403 Forbidden` error
- Error code: `forbidden`
- Message: "Permission denied" or "Access denied"

**Common Causes:**

**Insufficient API key permissions:**
```go
// Check if it's a permission error
if apiErr, ok := err.(*canvus.APIError); ok && apiErr.Code == canvus.ErrForbidden {
    log.Printf("Permission denied: %s", apiErr.Message)
    log.Println("Check that your API key has the required permissions")
}
```

**User lacks access to resource:**
```go
// The authenticated user may not have access to the canvas/widget
canvas, err := session.GetCanvas(ctx, canvasID)
if err != nil {
    if apiErr, ok := err.(*canvus.APIError); ok && apiErr.StatusCode == 403 {
        return fmt.Errorf("you don't have access to canvas %s", canvasID)
    }
    return err
}
```

**Solutions:**
1. Verify the API key has required permissions
2. Check user permissions in Canvus admin panel
3. Ensure user has access to the specific resource
4. Try with an administrator API key to confirm

---

### 7. Invalid Request Errors

**Symptoms:**
- `400 Bad Request` error
- Error code: `invalid_request`
- Message: "Invalid parameter" or validation errors

**Common Causes:**

**Missing required fields:**
```go
// Ensure required fields are provided
canvas, err := session.CreateCanvas(ctx, canvus.CreateCanvasRequest{
    Name: "My Canvas",  // Name is required
    // FolderID may also be required depending on server config
})
```

**Invalid field values:**
```go
// Check field value types and formats
widget, err := session.CreateWidget(ctx, canvasID, canvus.CreateWidgetRequest{
    WidgetType: "note",  // Must be a valid widget type
    Location: canvus.Location{
        X: 100,
        Y: 100,
    },
    Size: canvus.Size{
        Width:  200,
        Height: 150,
    },
})
```

**Invalid JSON:**
```go
// Check error details for validation issues
if apiErr, ok := err.(*canvus.APIError); ok && apiErr.Code == canvus.ErrInvalidRequest {
    log.Printf("Invalid request: %s", apiErr.Message)
    if len(apiErr.Details) > 0 {
        log.Printf("Details: %v", apiErr.Details)
    }
}
```

**Solutions:**
1. Check that all required fields are provided
2. Verify field values are valid types
3. Check the API documentation for field requirements
4. Look at the error details for specific validation failures

---

### 8. Resource Not Found Errors

**Symptoms:**
- `404 Not Found` error
- Error code: `not_found`
- Message: "Resource not found"

**Common Causes:**

**Incorrect resource ID:**
```go
// Verify the ID format and value
canvas, err := session.GetCanvas(ctx, canvasID)
if err != nil {
    if apiErr, ok := err.(*canvus.APIError); ok && apiErr.Code == canvus.ErrNotFound {
        log.Printf("Canvas %s not found - check the ID", canvasID)
    }
}
```

**Resource was deleted:**
```go
// Handle the case where resource may have been deleted
canvas, err := session.GetCanvas(ctx, canvasID)
if err != nil {
    if apiErr, ok := err.(*canvus.APIError); ok && apiErr.Code == canvus.ErrNotFound {
        // Canvas may have been deleted
        return nil, fmt.Errorf("canvas not found - it may have been deleted")
    }
}
```

**Incorrect endpoint:**
```go
// Ensure you're using the correct method and parameters
// For widgets, you need both canvas ID and widget ID
widget, err := session.GetWidget(ctx, canvasID, widgetID)
```

**Solutions:**
1. Verify the resource ID is correct
2. Confirm the resource exists (list resources first)
3. Check if the resource was recently deleted
4. Ensure you're using the correct endpoint/method

---

### 9. Batch Operation Failures

**Symptoms:**
- Partial batch failures
- `BatchResult` with errors
- Timeout on large batches

**Common Causes:**

**Individual operation failures:**
```go
// Check individual results
results, err := processor.ExecuteBatch(ctx, operations)
if err != nil {
    log.Printf("Batch error: %v", err)
}

// Check each result
summary := canvus.Summarize(results)
log.Printf("Successful: %d, Failed: %d", summary.Successful, summary.Failed)

for _, failed := range summary.FailedOperations {
    log.Printf("Operation %s failed: %v", failed.OperationID, failed.Error)
}
```

**Timeout on large batches:**
```go
// Increase batch timeout
batchConfig := canvus.DefaultBatchConfig()
batchConfig.Timeout = 10 * time.Minute  // Increase for large batches

// Also use context with timeout
ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
defer cancel()
```

**Too high concurrency:**
```go
// Reduce concurrency if server is overloaded
batchConfig := canvus.DefaultBatchConfig()
batchConfig.MaxConcurrency = 5  // Lower from default 10
```

**Missing metadata:**
```go
// Some operations require metadata
deleteOp := &canvus.BatchOperation{
    ID:       "delete-1",
    Type:     canvus.BatchOperationDelete,
    Resource: widget,
    Metadata: map[string]interface{}{
        "canvas_id":   canvasID,    // Required for widget delete
        "widget_type": "note",       // Required for widget delete
    },
}
```

**Solutions:**
1. Check individual `BatchResult` for errors
2. Increase timeout for large batches
3. Reduce concurrency if rate limited
4. Ensure all required metadata is provided
5. Use `ContinueOnError: true` to complete successful operations

---

### 10. Import/Export Issues

**Symptoms:**
- Import fails with validation errors
- Missing widgets after import
- Asset files not found

**Common Causes:**

**Invalid export format:**
```go
// Ensure export data is valid
exportData, err := session.ExportCanvas(ctx, canvasID)
if err != nil {
    return err
}

// Validate export data before saving
if len(exportData.Widgets) == 0 {
    log.Println("Warning: no widgets in export")
}
```

**Missing asset files:**
```go
// Ensure asset files exist before import
for _, asset := range importData.Assets {
    if _, err := os.Stat(asset.FilePath); os.IsNotExist(err) {
        return fmt.Errorf("asset file not found: %s", asset.FilePath)
    }
}
```

**Widget type mismatches:**
```go
// The SDK normalizes widget types, but watch for case sensitivity
// "Note" vs "note" - SDK handles this automatically
```

**Parent-child relationships:**
```go
// Widgets with parent_id must be imported after their parents
// The SDK handles ordering, but verify your data
for _, widget := range importData.Widgets {
    if widget.ParentID != "" {
        // Ensure parent exists in import data
        found := false
        for _, w := range importData.Widgets {
            if w.ID == widget.ParentID {
                found = true
                break
            }
        }
        if !found {
            log.Printf("Warning: parent %s not found for widget %s",
                widget.ParentID, widget.ID)
        }
    }
}
```

**Solutions:**
1. Validate export data before saving
2. Ensure all asset files are accessible
3. Check widget types are valid
4. Verify parent-child relationships are intact
5. Export to a test canvas first to verify

---

## Debugging with Custom RoundTripper

Use a custom HTTP RoundTripper to inspect all requests and responses:

```go
package main

import (
    "bytes"
    "io"
    "log"
    "net/http"
    "net/http/httputil"
    "os"
    "time"

    "github.com/jaypaulb/Canvus-Go-API/canvus"
)

// debugTransport logs all HTTP requests and responses
type debugTransport struct {
    transport http.RoundTripper
    logger    *log.Logger
}

func (t *debugTransport) RoundTrip(req *http.Request) (*http.Response, error) {
    start := time.Now()

    // Log request
    reqDump, _ := httputil.DumpRequestOut(req, true)
    t.logger.Printf("REQUEST:\n%s\n", redactSensitive(reqDump))

    // Execute request
    resp, err := t.transport.RoundTrip(req)
    if err != nil {
        t.logger.Printf("ERROR: %v (after %v)\n", err, time.Since(start))
        return nil, err
    }

    // Log response (save body for re-reading)
    bodyBytes, _ := io.ReadAll(resp.Body)
    resp.Body.Close()
    resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

    t.logger.Printf("RESPONSE [%d] (after %v):\n%s\n",
        resp.StatusCode, time.Since(start), bodyBytes)

    return resp, nil
}

func redactSensitive(data []byte) []byte {
    // Redact API keys and tokens from logs
    // Implementation depends on your security requirements
    return data
}

func main() {
    // Create debug logger
    logger := log.New(os.Stderr, "[DEBUG] ", log.LstdFlags|log.Lmicroseconds)

    // Configure session with debug transport
    config := canvus.DefaultSessionConfig()
    config.BaseURL = "https://your-server/api/v1"
    config.HTTPClient = &http.Client{
        Transport: &debugTransport{
            transport: http.DefaultTransport,
            logger:    logger,
        },
    }

    session := canvus.NewSession(config, canvus.WithAPIKey("your-api-key"))

    // Now all requests will be logged
    _, _ = session.ListCanvases(context.Background())
}
```

---

## Performance Tuning Tips

### 1. Optimize Connection Pooling

```go
config.HTTPClient = &http.Client{
    Transport: &http.Transport{
        MaxIdleConns:        100,
        MaxIdleConnsPerHost: 10,
        IdleConnTimeout:     90 * time.Second,
    },
}
```

### 2. Use Appropriate Timeouts

```go
// Short timeouts for simple operations
config.RequestTimeout = 10 * time.Second

// Long timeouts for imports/exports
ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
defer cancel()
```

### 3. Batch Operations

```go
// Use BatchProcessor for bulk operations
batchConfig := canvus.DefaultBatchConfig()
batchConfig.MaxConcurrency = 10  // Adjust based on server capacity
```

### 4. Paginate Large Results

```go
// Fetch in pages instead of loading everything
widgets, err := session.ListWidgets(ctx, canvasID, &canvus.ListOptions{
    Limit:  100,
    Offset: 0,
})
```

### 5. Reuse Sessions

```go
// Create session once, reuse for all operations
var session *canvus.Session

func init() {
    session = canvus.NewSession(config, canvus.WithAPIKey(apiKey))
}
```

### 6. Monitor Circuit Breaker

```go
// Configure circuit breaker for resilience
config.CircuitBreaker.MaxFailures = 5
config.CircuitBreaker.ResetTimeout = 30 * time.Second
```

---

## When to Contact Support

Contact support if you experience:

1. **Consistent 5xx errors** - Server-side issues that need investigation
2. **Unexplained behavior** - API returning unexpected results
3. **Performance degradation** - Consistent slow responses
4. **Data inconsistency** - Missing or corrupted data
5. **Authentication issues** - After verifying your credentials are correct
6. **Rate limit increases** - If your legitimate use case requires higher limits

### Information to Include

When contacting support, include:

1. **Request ID** - From `APIError.RequestID`
2. **Timestamp** - When the error occurred
3. **Error details** - Status code, error code, message
4. **SDK version** - From your go.mod
5. **Server version** - If known
6. **Steps to reproduce** - Minimal code example
7. **Expected vs actual behavior**

Example:

```go
// Capture error details for support
if apiErr, ok := err.(*canvus.APIError); ok {
    fmt.Printf("Support ticket info:\n")
    fmt.Printf("  Request ID: %s\n", apiErr.RequestID)
    fmt.Printf("  Status Code: %d\n", apiErr.StatusCode)
    fmt.Printf("  Error Code: %s\n", apiErr.Code)
    fmt.Printf("  Message: %s\n", apiErr.Message)
    fmt.Printf("  Details: %v\n", apiErr.Details)
}
```

---

## Quick Reference: Error Code to Solution

| Error Code | First Steps |
|------------|-------------|
| `unauthorized` (401) | Check API key, re-authenticate |
| `forbidden` (403) | Verify permissions |
| `not_found` (404) | Check resource ID |
| `invalid_request` (400) | Validate request fields |
| `too_many_requests` (429) | Add rate limiting |
| `internal_server_error` (500) | Retry, then contact support |
| `service_unavailable` (503) | Wait and retry |

For more detailed information, see:
- [Best Practices](BEST_PRACTICES.md)
- [Getting Started](GETTING_STARTED.md)
- [API Reference](API_REFERENCE.md)
