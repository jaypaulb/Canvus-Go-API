# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**Canvus Go SDK**: A complete, production-grade Go SDK for the Canvus API. Provides full access to all Canvus endpoints with strong typing, authentication, context support, centralized error handling, geometry utilities, and import/export functionality.

**Repository**: `Canvus-Go-API` - A Go module (`github.com/jaypaulb/Canvus-Go-API`) containing the SDK in the `canvus/` package.

**Note**: The CLI tool has been moved to a separate repository: [github.com/jaypaulb/canvus-cli](https://github.com/jaypaulb/canvus-cli)

## Important Notes for Claude Code

**Do NOT use the Read tool with directory paths.** The Read tool only works with files. When exploring directories, use:
- `Glob` tool for finding files by pattern (preferred, faster)
- `Bash` tool with `ls` or similar for directory exploration

This prevents "EISDIR: illegal operation on a directory, read" errors.

## Essential Development Commands

### Building
```bash
go build ./canvus/          # Build SDK package
```

### Running Tests
```bash
go test ./canvus/...              # Run all tests
go test ./canvus/... -v          # Run all tests with verbose output
go test -run TestName ./canvus/  # Run a specific test
go test -run TestName -v ./canvus/ # Run specific test verbosely
```

### Linting & Formatting
```bash
gofmt -s -w .              # Format all Go files (Go standard)
goimports -w .             # Organize imports
go vet ./...               # Run Go vet checker
```

### Test Configuration
Tests require `settings.json` in the repository root with:
```json
{
  "api_base_url": "https://your-server/api/v1/",
  "api_key": "your-api-key",
  "timeout_seconds": 5,
  "test_user": {
    "username": "test@example.com",
    "password": "password"
  },
  "test_canvas_id": "canvas-id-for-tests"
}
```

## Architecture Overview

### SDK Core (`canvus/` package)

The SDK is organized by API resource domain:

#### Session Management
- **session.go**: Core `Session` struct, HTTP request handling, authentication, response validation, and retry logic. All API methods are Session methods.
- **options.go**: Session configuration options (`WithAPIKey`, `WithToken`, etc.)
- **errors.go**: Centralized error handling with `APIError` type for server errors

#### Resource API Groups
- **canvases.go**: Canvas CRUD, move, copy, permissions
- **folders.go**: Folder operations
- **widgets.go**: Widget listing and dispatch logic. Single `CreateWidget` method dispatches to type-specific creators.
- **notes.go, anchors.go, images.go, pdfs.go, videos.go, connectors.go**: Type-specific widget creation and CRUD
- **clients.go, workspaces.go**: Client and workspace management
- **users.go, accesstokens.go, groups.go**: System user/auth management
- **serverconfig.go, serverinfo.go, license.go, auditlog.go**: Server administration
- **colorpresets.go, backgrounds.go, uploads.go, mipmaps.go, videoinputs.go, videooutputs.go**: Additional widget/asset endpoints

#### Advanced Features
- **types.go**: Shared type definitions, `Filter` struct for client-side filtering with wildcards and JSONPath support, `Filterable` interface
- **geometry.go**: Geometry utilities (`WidgetsContainId`, `WidgetsTouchId`, `contains`, `touches`) for widget spatial relationships
- **export.go**: Export widgets and assets to folders with asset file handling
- **import.go**: Import widgets from exported folders, restoring asset references and spatial relationships
- **batch.go**: Batch operations with concurrency and partial failure handling

#### Testing Support
- **test_helpers.go**: Utilities for creating unique test resources and cleanup
- **\*_test.go**: Integration tests for all endpoints

## Key Design Patterns

### Session Pattern
All API methods are methods on `*Session`. Create a session once per application:
```go
session := canvus.NewSession("https://api.example.com/v1", canvus.WithAPIKey("key"))
// Use for multiple requests
canvases, _ := session.ListCanvases(ctx)
widgets, _ := session.ListWidgets(ctx, canvasID, nil)
```

### Authentication
- **API Key**: Long-lived static tokens, recommended for most use cases
- **Login**: Username/password to obtain temporary token
- **Token Refresh**: POST token to `/users/login` to extend lifetime

Implemented via `Authenticator` interface and `transportWithAPIKey` round-tripper.

### Error Handling
Server errors wrapped in `APIError` with status code and message. Client errors return Go errors.
```go
if err != nil {
  if apiErr, ok := err.(*canvus.APIError); ok {
    fmt.Println("Server error:", apiErr.StatusCode, apiErr.Message)
  }
}
```

### Request/Response Handling
- **HTTP Method Dispatch**: `doRequest(ctx, method, path, req, resp, ...)` handles all HTTP communication
- **Response Validation**: Centralized validation with relaxed numeric comparison (float64 equality)
- **Automatic Retry**: Transient failures (network errors, 5xx) retried with exponential backoff
- **Context Support**: All methods accept `context.Context` for cancellation and timeouts

### Widget Dispatch Pattern
`CreateWidget` accepts a map and dispatches to type-specific creators (CreateNote, CreateAnchor, etc.) based on `widget_type` field:
```go
req := map[string]interface{}{
  "widget_type": "note",
  "name": "My Note",
  "location": map[string]interface{}{"x": 0, "y": 0},
}
widget, _ := session.CreateWidget(ctx, canvasID, req)
```

### Filtering
`Filter` struct supports client-side filtering with:
- **Wildcards**: `"*"` matches any value
- **Prefix/Suffix**: `"abc*"`, `"*123"`, `"*mid*"`
- **JSONPath**: `"$.location.x"` for nested field matching

```go
filter := &canvus.Filter{Criteria: map[string]interface{}{
  "widget_type": "note",
  "$.location.x": 100.0,
}}
widgets, _ := session.ListWidgets(ctx, canvasID, filter)
```

### Import/Export
- **Export**: `ExportWidgetsToFolder` exports widgets and asset files to a directory
- **Import**: `ImportWidgetsFromFolder` reads exported widgets and assets, restoring relationships
- **Round-Trip Safe**: Exported and re-imported widgets maintain full fidelity
- **Asset Handling**: Images, PDFs, videos exported as files, referenced in JSON

### Batch Operations
`BatchOperation` struct with type dispatch for concurrent bulk operations (move, copy, delete, pin, unpin) with partial failure tracking.

## Important Implementation Notes

### Testing Requirements
- **Resource Cleanup**: All tests must permanently delete created resources, not move to trash
- **Unique Identifiers**: Each test uses unique names/IDs to avoid collisions
- **Test Isolation**: Tests should be independent and runnable in any order
- **Integration Tests**: Use live Canvus server (configured in settings.json), not mocks

### Field Naming
- **JSON tags**: All field names in request/response models must use lowercase JSON tags (e.g., `widget_type`, `location`)
- **Canvus API Requirement**: API expects lowercase JSON field names for PATCH/POST operations

### Response Validation
- **Numeric Comparison**: Relaxed floating-point comparison (not strict equality) for cross-type validation
- **Widget Type Normalization**: Widget types handled case-insensitively in import/export logic
- **Partial Response Handling**: Server may not return complete response body; use polling/retry where needed

### Context & Cancellation
All API methods accept `context.Context` for cancellation and timeout support:
```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()
canvases, _ := session.ListCanvases(ctx)
```

## Code Organization Guidelines

### Adding New Endpoints
1. Define request/response types in `types.go` (if not already defined)
2. Create resource file (e.g., `newresource.go`) or extend existing
3. Implement methods on `*Session` following existing patterns
4. Add context, error wrapping, and documentation
5. Create `*_test.go` with integration tests using unique names
6. Ensure test cleanup (permanent deletion)
7. Update README.md with examples if it's a major feature

### Error Messages
- **Wrap Errors**: Use `fmt.Errorf("MethodName: %w", err)` pattern for context
- **API Errors**: Let `APIError` bubble up; wrap if adding context
- **Validation**: Return early with descriptive messages for invalid inputs

### Documentation
- **Godoc**: All public functions must have godoc comments (starts with function name)
- **Examples**: Include usage examples in godoc for complex methods
- **README**: Update README.md for new major features or API changes

## Common Task Patterns

### To add a new API endpoint:
1. Check if the resource type exists in types.go
2. Add method to appropriate file (canvases.go, widgets.go, etc.)
3. Use `s.doRequest()` with appropriate HTTP method
4. Wrap errors with method name context
5. Write integration test using unique test data

### To run a single test:
```bash
go test -run TestSpecificName -v ./canvus/
# Example: go test -run TestCreateNote -v ./canvus/
```

### To debug API issues:
1. Check settings.json has correct api_base_url and api_key
2. Enable verbose output in tests: `go test -v`
3. Review doRequest() logic in session.go for request/response handling
4. Check error types in errors.go for proper error classification
