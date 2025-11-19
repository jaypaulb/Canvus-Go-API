# Getting Started with the Canvus Go SDK

This guide will help you make your first successful API call in under 30 minutes.

## Prerequisites

Before you begin, ensure you have:

1. **Go 1.16 or later** installed on your machine
   ```bash
   go version  # Should output go1.16 or higher
   ```

2. **Access to a Canvus server** with:
   - The server URL (e.g., `https://your-canvus-server`)
   - Either an API key or user credentials (email/password)

3. **A Go module** initialized for your project
   ```bash
   mkdir my-canvus-project
   cd my-canvus-project
   go mod init my-canvus-project
   ```

## Installation

Install the Canvus Go SDK using `go get`:

```bash
go get github.com/jaypaulb/Canvus-Go-API/canvus
```

This will download the SDK and add it to your `go.mod` file.

## Creating Your First Client

The SDK uses a `Session` struct as the main entry point. Here's how to create one:

```go
package main

import (
    "github.com/jaypaulb/Canvus-Go-API/canvus"
)

func main() {
    // Create a session configuration with default settings
    config := canvus.DefaultSessionConfig()
    config.BaseURL = "https://your-canvus-server/api/v1"

    // Create a new session with API key authentication
    session := canvus.NewSession(config, canvus.WithAPIKey("your-api-key-here"))

    // Your session is now ready to use!
    _ = session
}
```

## Authentication Options

The SDK supports two primary authentication methods:

### Option 1: API Key Authentication (Recommended)

Use an API key for server-to-server integrations, automation scripts, and CLI tools. This is the simplest and most secure option for most use cases.

```go
config := canvus.DefaultSessionConfig()
config.BaseURL = "https://your-canvus-server/api/v1"

session := canvus.NewSession(config, canvus.WithAPIKey("your-api-key"))
```

The API key is passed in the `Private-Token` header with every request.

### Option 2: Login with Email/Password

Use the `Login()` method when you need to authenticate with user credentials. This returns a temporary token that's stored in the session.

```go
config := canvus.DefaultSessionConfig()
config.BaseURL = "https://your-canvus-server/api/v1"

session := canvus.NewSession(config)

// Login to obtain a token
err := session.Login(ctx, "user@example.com", "password")
if err != nil {
    log.Fatalf("Login failed: %v", err)
}

// The session now has authentication configured
// Remember to logout when done
defer session.Logout(ctx)
```

## Making Your First API Call

Here's a complete example that lists all canvases:

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/jaypaulb/Canvus-Go-API/canvus"
)

func main() {
    // Create context for the request
    ctx := context.Background()

    // Configure the session
    config := canvus.DefaultSessionConfig()
    config.BaseURL = "https://your-canvus-server/api/v1"

    // Create session with API key authentication
    session := canvus.NewSession(config, canvus.WithAPIKey("your-api-key"))

    // List all canvases
    canvases, err := session.ListCanvases(ctx)
    if err != nil {
        log.Fatalf("Failed to list canvases: %v", err)
    }

    // Print the results
    fmt.Printf("Found %d canvases:\n", len(canvases))
    for _, canvas := range canvases {
        fmt.Printf("  - %s (ID: %s)\n", canvas.Name, canvas.ID)
    }
}
```

Run this example:

```bash
go run main.go
```

## Handling Errors

The SDK returns `*canvus.APIError` for API-specific errors. Always check and handle errors appropriately:

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/jaypaulb/Canvus-Go-API/canvus"
)

func main() {
    ctx := context.Background()

    config := canvus.DefaultSessionConfig()
    config.BaseURL = "https://your-canvus-server/api/v1"
    session := canvus.NewSession(config, canvus.WithAPIKey("your-api-key"))

    // Try to get a non-existent canvas
    canvas, err := session.GetCanvas(ctx, "non-existent-id")
    if err != nil {
        // Check if it's an API error
        if apiErr, ok := err.(*canvus.APIError); ok {
            fmt.Printf("API Error:\n")
            fmt.Printf("  Status Code: %d\n", apiErr.StatusCode)
            fmt.Printf("  Error Code:  %s\n", apiErr.Code)
            fmt.Printf("  Message:     %s\n", apiErr.Message)

            // Handle specific error codes
            switch apiErr.Code {
            case canvus.ErrNotFound:
                fmt.Println("  -> Canvas not found")
            case canvus.ErrUnauthorized:
                fmt.Println("  -> Check your API key")
            case canvus.ErrForbidden:
                fmt.Println("  -> You don't have permission")
            default:
                fmt.Println("  -> Unexpected error")
            }
        } else {
            // Network or other errors
            log.Fatalf("Unexpected error: %v", err)
        }
        return
    }

    fmt.Printf("Canvas: %s\n", canvas.Name)
}
```

### Common Error Codes

| Error Code | HTTP Status | Description |
|------------|-------------|-------------|
| `ErrInvalidRequest` | 400 | Bad request, invalid parameters |
| `ErrUnauthorized` | 401 | Invalid or missing authentication |
| `ErrForbidden` | 403 | Not authorized for this action |
| `ErrNotFound` | 404 | Resource not found |
| `ErrConflict` | 409 | Resource conflict |
| `ErrTooManyRequests` | 429 | Rate limited |
| `ErrInternalServer` | 500 | Server error |
| `ErrServiceUnavailable` | 503 | Service unavailable |

## Complete Working Example

Here's a complete example that demonstrates session creation, authentication, API calls, and proper error handling:

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"
    "time"

    "github.com/jaypaulb/Canvus-Go-API/canvus"
)

func main() {
    // Get configuration from environment variables
    serverURL := os.Getenv("CANVUS_URL")
    apiKey := os.Getenv("CANVUS_API_KEY")

    if serverURL == "" || apiKey == "" {
        log.Fatal("Please set CANVUS_URL and CANVUS_API_KEY environment variables")
    }

    // Create context with timeout
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    // Configure the session
    config := canvus.DefaultSessionConfig()
    config.BaseURL = serverURL + "/api/v1"
    config.RequestTimeout = 30 * time.Second
    config.MaxRetries = 3

    // Create session
    session := canvus.NewSession(config, canvus.WithAPIKey(apiKey))

    // List canvases
    canvases, err := session.ListCanvases(ctx)
    if err != nil {
        handleError(err)
        return
    }

    fmt.Printf("Successfully connected! Found %d canvases.\n\n", len(canvases))

    // Display canvas information
    for i, canvas := range canvases {
        fmt.Printf("%d. %s\n", i+1, canvas.Name)
        fmt.Printf("   ID: %s\n", canvas.ID)
        fmt.Printf("   Created: %s\n", canvas.CreatedAt)
        fmt.Println()
    }

    // Get system information
    license, err := session.GetLicense(ctx)
    if err != nil {
        handleError(err)
        return
    }

    fmt.Printf("Server License: %s\n", license.LicenseType)
}

func handleError(err error) {
    if apiErr, ok := err.(*canvus.APIError); ok {
        fmt.Fprintf(os.Stderr, "API Error [%d]: %s\n", apiErr.StatusCode, apiErr.Message)
        if apiErr.RequestID != "" {
            fmt.Fprintf(os.Stderr, "Request ID: %s\n", apiErr.RequestID)
        }
    } else {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
    }
    os.Exit(1)
}
```

## Next Steps

Now that you've made your first API call, explore these resources:

- **[Best Practices](BEST_PRACTICES.md)** - Learn recommended patterns for error handling, session management, and more
- **[Troubleshooting](TROUBLESHOOTING.md)** - Solutions for common issues
- **[API Reference](API_REFERENCE.md)** - Complete API documentation
- **[Examples](EXAMPLES.md)** - More code examples for specific use cases

## Quick Reference

### Session Configuration Options

```go
config := canvus.DefaultSessionConfig()
config.BaseURL = "https://your-server/api/v1"
config.MaxRetries = 3                        // Retry failed requests
config.RetryWaitMin = 100 * time.Millisecond // Min retry delay
config.RetryWaitMax = time.Second            // Max retry delay
config.RequestTimeout = 30 * time.Second     // Request timeout
```

### Common Operations

```go
// List canvases
canvases, err := session.ListCanvases(ctx)

// Get a specific canvas
canvas, err := session.GetCanvas(ctx, "canvas-id")

// Create a canvas
newCanvas, err := session.CreateCanvas(ctx, canvus.CreateCanvasRequest{
    Name: "My Canvas",
})

// List widgets in a canvas
widgets, err := session.ListWidgets(ctx, "canvas-id", nil)

// Get system license info
license, err := session.GetLicense(ctx)
```

### Context Usage

Always use `context.Context` for cancellation and timeouts:

```go
// With timeout
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

// With cancellation
ctx, cancel := context.WithCancel(context.Background())
// Call cancel() when you want to abort the operation
```

## Getting Help

If you encounter issues:

1. Check the [Troubleshooting Guide](TROUBLESHOOTING.md)
2. Review the [API Reference](API_REFERENCE.md)
3. Open an issue on GitHub

---

You're now ready to start building with the Canvus Go SDK!
