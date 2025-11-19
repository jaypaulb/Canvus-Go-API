# Tech Stack

## Language & Runtime

- **Language**: Go 1.24.1
- **Module System**: Go modules (canvus-go-api)
- **Minimum Supported Version**: Go 1.24.1

## Core Dependencies

### HTTP & Networking
- **net/http** (Go Standard Library): Core HTTP client for all API requests
- **http.RoundTripper**: Custom middleware pattern for request/response interception, authentication, and retry logic
- **context** (Go Standard Library): Request cancellation, timeouts, and tracing context propagation

### Testing
- **testify** (github.com/stretchr/testify v1.11.1): Assertion library for unit and integration tests; provides assert and require packages for readable test code
- **testing** (Go Standard Library): Go's built-in testing framework for all test files

### CLI Framework
- **Cobra** (github.com/spf13/cobra v1.9.1): Modern command-line application framework for canvus-cli with subcommands, flags, and help generation
- **pflag** (github.com/spf13/pflag v1.0.6): POSIX-compatible flag parsing (indirect dependency of Cobra)

## Architecture Patterns

### Session Management
- **Session Struct**: Central authenticated session managing API base URL, HTTP client, and authentication state
- **Authenticator Interface**: Abstraction for authentication strategies (API key, token refresh)
- **APIKeyAuthenticator**: Implements header-based API key authentication
- **RoundTripper Middleware**: Custom HTTP transport for injecting authentication, retries, and logging without modifying individual requests

### Error Handling
- **APIError**: Typed error struct capturing HTTP status codes, error messages, and response bodies
- **Error Type Assertions**: Go idiomatic error handling using type assertions and interface{} for precise error discrimination
- **Structured Error Responses**: Server error responses parsed into structured Go types for programmatic handling

### Request/Response Model
- **Typed Structs**: All API requests and responses are strongly typed Go structs with JSON field tags matching the Canvus API specification
- **Lowercase JSON Fields**: All JSON field names in request bodies are lowercase to match Canvus API requirements
- **Flexible Response Handling**: Response unmarshaling with relaxed numeric validation and case-insensitive widget type handling for cross-type compatibility

### Batch Processing
- **BatchProcessor**: Framework for efficiently processing multiple operations in configurable batch sizes
- **Automatic Retries**: Exponential backoff retry logic for transient failures within batch operations
- **Concurrent Request Handling**: Goroutine-based concurrent processing with configurable concurrency limits

### Import/Export
- **Round-trip Safety**: Bidirectional import/export ensuring fidelity when transferring canvases between workspaces
- **Asset File Handling**: Export writes asset files (images, PDFs, videos) to disk with JSON references; import reads files and recreates widgets with correct spatial relationships
- **Type Normalization**: Case-insensitive widget type handling and relaxed numeric validation for cross-version compatibility

### Filtering & Search
- **Filter Abstraction**: Client-side filtering with wildcard, prefix/suffix, and partial match support
- **JSONPath-like Selectors**: Support for nested field selectors (e.g., "$.location.x") in filter criteria
- **Cross-Canvas Search**: FindWidgetsAcrossCanvases utility for querying widgets across multiple canvases with flexible matching

### Geometry Utilities
- **Spatial Algorithms**: Widget containment (WidgetsContainId) and overlap detection (WidgetsTouchId) based on location and size
- **Zone Detection**: Utility functions for determining which widgets are contained within or touch a given widget

### Configuration Management
- **settings.json**: Single source of truth for configuration including API base URL, API key, test user credentials
- **File-based Configuration**: JSON-based configuration for Windows/PowerShell compatibility and simplicity

## Code Organization

### Package Structure
- **canvus/**: Main SDK package containing all API methods and types
  - **session.go**: Session management and HTTP client setup
  - **clients.go**: Client resource operations
  - **users.go**: User management and authentication
  - **accesstokens.go**: Token creation and management
  - **canvases.go**: Canvas CRUD and operations
  - **widgets.go**: Widget management for all types
  - **assets.go**: Asset lifecycle management
  - **folders.go**: Folder organization
  - **workspaces.go**: Workspace operations
  - **batch.go**: Batch processing framework
  - **import.go / export.go**: Import/export functionality
  - **geometry.go**: Spatial utilities
  - **errors.go**: Error types and handling
  - **types.go**: Shared type definitions
  - **test_helpers.go**: Testing utilities
  - **tests/**: Integration test fixtures and configurations
  - **export/**: Export-related utilities and helper functions

- **cmd/canvus-cli/**: CLI application
  - **main.go**: Entry point and root command setup
  - **widget.go**: Widget subcommands
  - **widget_create.go**: Widget creation implementation
  - Additional command implementations (to be developed)

### Testing Approach

#### Test Organization
- **Unit Tests**: File-level tests (e.g., users_test.go) for isolated functionality
- **Integration Tests**: canvus_test.go for end-to-end testing against live Canvus server
- **Test Fixtures**: tests/ directory containing test server URLs, credentials, and helper functions

#### Testing Standards
- **Test Isolation**: Each test creates unique resources to avoid collisions
- **Complete Cleanup**: All created resources are permanently deleted (not moved to trash) even on test failure
- **Server Configuration**: settings.json provides API base URL, admin API key, and test user credentials
- **Comprehensive Coverage**: All major endpoints and error paths are tested
- **Test Execution Order**: Tests follow dependency order (system management → canvas → client/workspace → widgets/assets)

### Documentation

#### Code Documentation
- **Godoc Comments**: Comprehensive comments for all exported types, functions, and methods
- **Usage Examples**: Code examples in godoc comments showing common usage patterns
- **Error Documentation**: Clear documentation of error conditions and recovery strategies

#### External Documentation
- **README.md**: Overview, feature list, installation, usage examples, and troubleshooting
- **PRD.md**: Detailed product requirements, API coverage map, and design notes
- **CONTRIBUTING.md**: Contribution guidelines, Windows/PowerShell compatibility, and branching conventions
- **tasks.md**: Development task tracking and progress notes
- **Abstractions.md**: Documentation of SDK abstractions, utilities, and advanced features

## Data Serialization

- **JSON**: Standard library json package for all request/response serialization
- **Struct Tags**: JSON field tags for mapping Go types to API contract
- **Custom Marshaling**: Selective use of custom MarshalJSON/UnmarshalJSON for special types (time, geometry, enums)

## Deployment & Distribution

- **Go Module Distribution**: Installable via `go get github.com/jaypaulb/Canvus-Go-API/canvus`
- **Binary Distribution**: CLI distributable as standalone binary (cross-platform compilation)
- **Semantic Versioning**: Git tags for release versioning (v0.x.x, v1.x.x)

## Platform Compatibility

- **Development Environment**: Windows with PowerShell
- **Cross-Platform Support**: SDK and CLI work on Windows, macOS, and Linux
- **Shell Compatibility**: No Linux-specific shell scripts; all scripts use POSIX-compatible patterns
- **API Compatibility**: Compatible with any Canvus server (MT-Canvus-Server, custom deployments)

## Performance Considerations

- **HTTP Connection Reuse**: http.Client maintains connection pools for efficient requests
- **Concurrent Safe**: Session struct is safe for concurrent use across goroutines
- **Memory Efficient**: Streaming response handling for large assets and paginated results
- **Retry Strategy**: Exponential backoff with jitter for resilience against transient failures

## Security

- **API Key Management**: Support for environment-based and file-based API key configuration
- **Secure Authentication**: Private-Token header for authentication; no credentials in URLs
- **Error Message Handling**: Server error messages carefully handled without exposing internal implementation details
- **Context Cancellation**: Proper context cancellation prevents resource leaks and hanging requests

## Future Technology Considerations

- **gRPC**: Potential for gRPC bindings for higher performance inter-service communication
- **OpenTelemetry**: Structured logging and distributed tracing integration
- **Protobuf**: Alternative to JSON for future high-performance implementations
- **WebSocket**: Real-time subscription support for event streaming
