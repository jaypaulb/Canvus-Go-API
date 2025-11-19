# Canvus Go SDK Examples Index

This document provides an index of all runnable examples included with the SDK. Each example demonstrates specific SDK features and can be run with `go run`.

## Quick Navigation

- [Getting Started](#getting-started) - First steps with the SDK
- [Authentication](#authentication) - API key and login authentication
- [Canvas Operations](#canvas-operations) - Canvas lifecycle management
- [Widget Operations](#widget-operations) - Creating and managing widgets
- [User Management](#user-management) - User provisioning and tokens
- [Batch Operations](#batch-operations) - Bulk operations with retry
- [Import/Export](#importexport) - Canvas data transfer
- [Error Handling](#error-handling) - Recovery patterns
- [Context Management](#context-management) - Cancellation and timeouts

## Difficulty Levels

- **Beginner** - Simple operations, minimal setup
- **Intermediate** - Multiple operations, basic error handling
- **Advanced** - Complex patterns, production considerations

---

## Getting Started

### Basic Usage

**File**: `examples/getting_started/main.go`
**Difficulty**: Beginner
**Topics**: Session creation, API key auth, listing canvases

The simplest example to verify your SDK setup is working. Creates a session with an API key and lists all available canvases.

**What you'll learn**:
- Creating a session with `NewSessionFromConfig`
- Using `WithAPIKey` for authentication
- Making your first API call with `ListCanvases`
- Basic error handling

**Prerequisites**:
- Go 1.16 or later
- Canvus server access
- Valid API key

```bash
export CANVUS_SERVER="https://your-server/api/v1"
export CANVUS_API_KEY="your-api-key"
go run examples/getting_started/main.go
```

---

## Authentication

### API Key Authentication

**File**: `examples/authentication/api_key.go`
**Difficulty**: Beginner
**Topics**: API key setup, Private-Token header, session configuration

Demonstrates the recommended authentication method using a static API key.

**What you'll learn**:
- Configuring API key authentication
- Understanding the Private-Token header
- Session configuration options
- Making authenticated requests

### Username/Password Login

**File**: `examples/authentication/login_password.go`
**Difficulty**: Beginner
**Topics**: Login flow, token management, logout

Shows how to authenticate with username and password, manage the session token, and properly logout.

**What you'll learn**:
- Using the `Login()` method
- Token-based authentication flow
- Session cleanup with `Logout()`
- When to use login vs API key

---

## Canvas Operations

### Create and Manage Canvases

**File**: `examples/canvases/create_and_manage.go`
**Difficulty**: Intermediate
**Topics**: CRUD operations, filtering, pagination, copying

Complete canvas lifecycle example covering creation, listing with filters, updates, copying, and deletion.

**What you'll learn**:
- Creating canvases with `CreateCanvas`
- Listing with client-side `Filter`
- Updating canvas properties
- Copying canvases between folders
- Proper cleanup/deletion

**Key patterns demonstrated**:
- Request/response types
- Filter criteria syntax
- Handling canvas permissions

---

## Widget Operations

### Create and Search Widgets

**File**: `examples/widgets/create_and_search.go`
**Difficulty**: Intermediate
**Topics**: Widget types, cross-canvas search, geometry utilities

Demonstrates working with different widget types, searching across canvases, and using geometry utilities.

**What you'll learn**:
- Creating notes, browsers, and image widgets
- Using `ListWidgets` with filters
- Cross-canvas search with `FindWidgetsAcrossCanvases`
- Spatial queries with `WidgetsContainId`
- Widget parent/child relationships

**Widget types covered**:
- Note widgets (text content)
- Browser widgets (embedded web content)
- Image widgets (with asset upload)

**Advanced features**:
- Wildcard pattern matching
- JSONPath-like nested field queries
- Geometry containment checks

---

## User Management

### Provision and Manage Users

**File**: `examples/users/provision_and_manage.go`
**Difficulty**: Intermediate
**Topics**: User CRUD, access tokens, group membership

Shows how to programmatically manage users, including creation, token generation, and group assignment.

**What you'll learn**:
- Creating users with `CreateUser`
- Generating access tokens for users
- Managing group membership
- Permission considerations
- Cleanup procedures

**Use cases**:
- Automated user provisioning
- Service account management
- Team onboarding automation

---

## Batch Operations

### Bulk Operations with Retry

**File**: `examples/batch/bulk_operations.go`
**Difficulty**: Advanced
**Topics**: BatchProcessor, automatic retry, concurrent operations

Demonstrates the batch processing framework for performing bulk operations efficiently with built-in retry logic.

**What you'll learn**:
- Creating a `BatchProcessor`
- Building operations with `BatchOperationBuilder`
- Executing batches with automatic retry
- Understanding batch results and summary
- Progress tracking patterns

**Operations covered**:
- Bulk delete
- Bulk move
- Bulk copy
- Pin/unpin widgets

**Production patterns**:
- Configuring retry behavior
- Handling partial failures
- Concurrent operation limits

---

## Import/Export

### Round-Trip Import/Export

**File**: `examples/import_export/round_trip.go`
**Difficulty**: Advanced
**Topics**: Export, import, asset handling, fidelity guarantees

Complete example of exporting canvas content with assets and importing to another canvas while maintaining full fidelity.

**What you'll learn**:
- Exporting widgets with `ExportWidgetsToFolder`
- Understanding the export JSON format
- Asset file handling (images, PDFs, videos)
- Importing with `ImportWidgetsToRegion`
- ID remapping for relationships

**Fidelity guarantees**:
- Widget properties preserved
- Spatial positions maintained
- Parent/child relationships intact
- Connector endpoints remapped

**Use cases**:
- Canvas backup and restore
- Template distribution
- Content migration between servers

---

## Error Handling

### Recovery Patterns

**File**: `examples/error_handling/recovery_patterns.go`
**Difficulty**: Advanced
**Topics**: Typed errors, error codes, retry patterns, circuit breaker

Comprehensive error handling example showing production-ready patterns for dealing with failures.

**What you'll learn**:
- Type asserting `*canvus.APIError`
- Checking specific error codes
- Implementing retry logic
- Using the built-in circuit breaker
- Graceful degradation strategies

**Patterns covered**:
- Error type checking
- Retryable vs non-retryable errors
- Exponential backoff
- Circuit breaker state management
- Context error handling

**Error scenarios**:
- Authentication failures
- Network timeouts
- Rate limiting
- Resource not found
- Permission denied

---

## Context Management

### Cancellation and Timeouts

**File**: `examples/context/cancellation_and_timeouts.go`
**Difficulty**: Intermediate
**Topics**: Context cancellation, request timeouts, deadline management

Shows how to properly use Go contexts for cancellation and timeout control.

**What you'll learn**:
- Setting request timeouts with `context.WithTimeout`
- Implementing cancellation with `context.WithCancel`
- Deadline propagation
- Graceful shutdown patterns
- Context best practices

**Patterns covered**:
- Per-request timeouts
- Operation cancellation
- Signal handling for shutdown
- Cleanup on cancellation

---

## Running Examples

### Prerequisites

1. **Install the SDK**:
   ```bash
   go get github.com/jaypaulb/Canvus-Go-API/canvus
   ```

2. **Set environment variables**:
   ```bash
   export CANVUS_SERVER="https://your-canvus-server/api/v1"
   export CANVUS_API_KEY="your-api-key"
   # Or for login-based auth:
   export CANVUS_EMAIL="user@example.com"
   export CANVUS_PASSWORD="your-password"
   ```

3. **Run an example**:
   ```bash
   go run examples/getting_started/main.go
   ```

### Common Environment Variables

| Variable | Description | Required |
|----------|-------------|----------|
| `CANVUS_SERVER` | Server URL (e.g., `https://server/api/v1`) | Yes |
| `CANVUS_API_KEY` | API key for authentication | For API key auth |
| `CANVUS_EMAIL` | User email for login | For login auth |
| `CANVUS_PASSWORD` | User password for login | For login auth |
| `CANVUS_CANVAS_ID` | Default canvas ID for operations | Some examples |

---

## Example Structure

Each example follows a consistent structure:

```go
// Package documentation
package main

import (
    // Standard library imports
    "context"
    "fmt"
    "os"

    // SDK import
    "github.com/jaypaulb/Canvus-Go-API/canvus"
)

func main() {
    // 1. Load configuration from environment
    server := os.Getenv("CANVUS_SERVER")
    apiKey := os.Getenv("CANVUS_API_KEY")

    // 2. Create session
    session := canvus.NewSessionFromConfig(server, apiKey)
    ctx := context.Background()

    // 3. Perform operations with error handling
    result, err := session.SomeOperation(ctx, params)
    if err != nil {
        // Handle error appropriately
        if apiErr, ok := err.(*canvus.APIError); ok {
            fmt.Printf("API Error %d: %s\n", apiErr.StatusCode, apiErr.Message)
        }
        os.Exit(1)
    }

    // 4. Process results
    fmt.Printf("Success: %+v\n", result)

    // 5. Cleanup (if needed)
}
```

---

## Suggested Learning Path

### For Beginners

1. Start with `getting_started/main.go` to verify your setup
2. Try `authentication/api_key.go` to understand auth
3. Explore `canvases/create_and_manage.go` for basic CRUD

### For Intermediate Users

1. Work through `widgets/create_and_search.go` for widget operations
2. Study `users/provision_and_manage.go` for admin tasks
3. Learn `context/cancellation_and_timeouts.go` for production patterns

### For Advanced Users

1. Master `batch/bulk_operations.go` for efficiency
2. Understand `import_export/round_trip.go` for data transfer
3. Implement `error_handling/recovery_patterns.go` for resilience

---

## Contributing Examples

We welcome new examples! If you'd like to contribute:

1. Follow the existing structure and naming conventions
2. Include comprehensive comments explaining each step
3. Add entries to this index with appropriate difficulty level
4. Ensure examples are runnable with `go run`
5. Test with a real Canvus server when possible

---

## Related Documentation

- [GETTING_STARTED.md](GETTING_STARTED.md) - Quick start guide
- [API_REFERENCE.md](API_REFERENCE.md) - Complete method reference
- [BEST_PRACTICES.md](BEST_PRACTICES.md) - Recommended patterns
- [TROUBLESHOOTING.md](TROUBLESHOOTING.md) - Common issues
