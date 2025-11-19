# Canvus Go SDK Examples

This directory contains runnable examples demonstrating how to use the Canvus Go SDK.

## Prerequisites

Before running these examples, ensure you have:

1. **Go 1.16 or later** installed
2. **Access to a Canvus server** (MTCS API endpoint)
3. **Valid credentials** (API key or username/password)

## Configuration

All examples use environment variables for configuration:

```bash
# Required: Your Canvus server URL
export CANVUS_API_URL="https://your-canvus-server.example.com/api/public/v1"

# For API key authentication
export CANVUS_API_KEY="your-api-key-here"

# For login/password authentication
export CANVUS_USERNAME="your-email@example.com"
export CANVUS_PASSWORD="your-password"
```

## Running Examples

Each example is a standalone Go program that can be run with `go run`:

```bash
# Run the getting started example
cd examples/getting_started
go run main.go

# Run authentication examples
cd examples/authentication
go run api_key.go
go run login_password.go

# Run canvas management example
cd examples/canvases
go run create_and_manage.go

# Run widget operations example
cd examples/widgets
go run create_and_search.go
```

## Example Directory Structure

```
examples/
  README.md                          # This file
  getting_started/
    main.go                          # Simple first example - list canvases
  authentication/
    api_key.go                       # API key authentication
    login_password.go                # Login/password authentication
  canvases/
    create_and_manage.go             # Canvas lifecycle operations
  widgets/
    create_and_search.go             # Widget operations and search
  users/
    (future examples)
```

## Examples Overview

### Getting Started (Beginner)

- **getting_started/main.go** - Your first API call. Shows basic initialization, authentication, and listing canvases with error handling.

### Authentication (Beginner)

- **authentication/api_key.go** - Demonstrates API key authentication using the `WithAPIKey` functional option.
- **authentication/login_password.go** - Shows login/password authentication flow with proper session cleanup.

### Canvas Operations (Intermediate)

- **canvases/create_and_manage.go** - Complete canvas lifecycle: create, list, get, update, copy, and delete canvases with pagination and filtering.

### Widget Operations (Intermediate)

- **widgets/create_and_search.go** - Widget CRUD operations, cross-canvas search, and using geometry utilities like `WidgetsContainId`.

## Best Practices Demonstrated

All examples follow these best practices:

1. **Environment-based configuration** - No hardcoded credentials
2. **Proper error handling** - Check and handle all errors
3. **Context usage** - Use context for request cancellation and timeouts
4. **Resource cleanup** - Logout sessions when using login/password auth
5. **Comprehensive comments** - Explain what each section does

## Troubleshooting

### Common Issues

**"connection refused" error**
- Verify `CANVUS_API_URL` is correct and the server is accessible
- Check if the URL includes the API path (e.g., `/api/public/v1`)

**"401 Unauthorized" error**
- Verify your API key or credentials are correct
- Check if the API key has sufficient permissions

**"certificate verify failed" error**
- The SDK uses `InsecureSkipVerify: true` by default for self-signed certificates
- For production, configure proper TLS certificates

### Getting Help

- Review the [SDK documentation](../docs/GETTING_STARTED.md)
- Check [troubleshooting guide](../docs/TROUBLESHOOTING.md)
- Review the [API reference](../docs/API_REFERENCE.md)

## Contributing

To add a new example:

1. Create a new `.go` file in the appropriate subdirectory
2. Follow the existing example patterns
3. Include comprehensive comments
4. Ensure it compiles with `go build`
5. Test it runs successfully with `go run`
6. Update this README with the new example
