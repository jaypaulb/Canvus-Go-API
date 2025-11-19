# Canvus Go SDK Starter Templates

This directory contains production-ready starter templates for building applications with the Canvus Go SDK. Each template demonstrates best practices and can be used as a foundation for real projects.

## Overview

These templates provide complete, working Go programs that you can copy and customize for your specific needs. They include proper error handling, configuration management, logging, and graceful shutdown patterns.

## Templates

| Template | Description | Use Case |
|----------|-------------|----------|
| `minimal_cli.go` | Command-line tool structure | One-off scripts, admin utilities |
| `web_service.go` | HTTP service with SDK | REST APIs, webhooks, dashboards |
| `batch_job.go` | Background job processor | Bulk operations, scheduled tasks |
| `integration_service.go` | Microservice pattern | Production services, integrations |

## Getting Started

### 1. Choose a Template

Select the template that best matches your use case:

- **CLI tool**: Use `minimal_cli.go` for command-line utilities that run once and exit
- **Web service**: Use `web_service.go` for long-running HTTP services
- **Batch job**: Use `batch_job.go` for background processing tasks
- **Microservice**: Use `integration_service.go` for production-grade services

### 2. Copy the Template

```bash
# Create your project directory
mkdir my-canvus-tool
cd my-canvus-tool

# Initialize a Go module
go mod init github.com/yourusername/my-canvus-tool

# Copy the template
cp /path/to/canvus-go-api/templates/minimal_cli.go main.go

# Add the SDK dependency
go get github.com/jaypaulb/Canvus-Go-API/canvus
```

### 3. Customize

Each template contains `TODO:` comments marking customization points. Search for these markers and update them for your needs:

```bash
# Find all customization points
grep -n "TODO:" main.go
```

### 4. Configure

Set the required environment variables:

```bash
export CANVUS_API_URL="https://your-canvus-server.example.com/api/public/v1"
export CANVUS_API_KEY="your-api-key"
```

### 5. Build and Run

```bash
# Build
go build -o my-tool

# Run
./my-tool
```

## Customization Points

All templates use `TODO:` comments to mark places that need customization:

- **TODO: Add your imports here** - Add additional imports
- **TODO: Define your command-line flags** - Add CLI flags
- **TODO: Add your business logic** - Main application logic
- **TODO: Add your HTTP handlers** - REST API endpoints
- **TODO: Add your configuration** - Configuration options

## Common Patterns

### Configuration from Environment

All templates use environment variables for configuration:

```go
apiURL := os.Getenv("CANVUS_API_URL")
if apiURL == "" {
    log.Fatal("CANVUS_API_URL environment variable is required")
}
```

### Error Handling

Templates demonstrate proper error handling with typed errors:

```go
if apiErr, ok := err.(*canvus.APIError); ok {
    log.Printf("API Error %d: %s", apiErr.StatusCode, apiErr.Message)
}
```

### Graceful Shutdown

Long-running services implement graceful shutdown:

```go
quit := make(chan os.Signal, 1)
signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
<-quit

ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()
server.Shutdown(ctx)
```

### Session Management

All templates show proper session lifecycle:

```go
// Create session
session := canvus.NewSession(cfg, canvus.WithAPIKey(apiKey))

// Use session
canvases, err := session.ListCanvases(ctx, nil)

// Clean up (if using login/password)
defer session.Logout(ctx)
```

## Template Details

### minimal_cli.go

A command-line tool template with:
- Flag parsing with `flag` package
- Environment variable configuration
- Exit codes for success/failure
- Structured error output

**Best for**: Admin scripts, one-off operations, data exports

### web_service.go

An HTTP service template with:
- RESTful endpoints
- Request/response logging middleware
- Health check endpoint
- Graceful shutdown
- Session management per request

**Best for**: REST APIs, webhook receivers, admin dashboards

### batch_job.go

A batch processing template with:
- Progress reporting
- Error aggregation
- Resumable operations
- Concurrent processing
- Checkpointing

**Best for**: Bulk updates, scheduled maintenance, data migrations

### integration_service.go

A production microservice template with:
- Configuration management
- Prometheus-style metrics
- Health checks (readiness/liveness)
- Graceful shutdown
- Structured logging

**Best for**: Production services, integrations with other systems

## Dependencies

All templates use only:
- Go standard library
- Canvus Go SDK (`github.com/jaypaulb/Canvus-Go-API/canvus`)

No additional dependencies are required, making these templates easy to build and deploy.

## Best Practices

These templates follow best practices from the SDK documentation:

1. **Environment-based configuration** - No hardcoded credentials
2. **Context usage** - All SDK calls use context for cancellation
3. **Proper error handling** - Typed errors, logging, and recovery
4. **Resource cleanup** - Sessions are properly closed
5. **Graceful shutdown** - Services handle termination signals

## Documentation

For more information, see:
- [Getting Started Guide](../docs/GETTING_STARTED.md)
- [Best Practices](../docs/BEST_PRACTICES.md)
- [API Reference](../docs/API_REFERENCE.md)
- [Troubleshooting](../docs/TROUBLESHOOTING.md)

## Contributing

To add a new template:

1. Create a new `.go` file following the existing patterns
2. Include comprehensive `TODO:` comments for customization
3. Ensure it compiles with `go build`
4. Add documentation to this README
5. Test the template as a real project
