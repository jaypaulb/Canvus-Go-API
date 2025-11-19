# Changelog

All notable changes to the Canvus Go SDK will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to the [Versioning](docs/VERSIONING.md) scheme documented for this SDK.

## [Unreleased]

### Added
- Nothing yet

### Changed
- Nothing yet

### Deprecated
- Nothing yet

### Removed
- Nothing yet

### Fixed
- Nothing yet

### Security
- Nothing yet

---

## [0.1.0] - 2025-11-19

Initial pre-release of the Canvus Go SDK.

### Added

#### Core Infrastructure
- Go module with proper import path: `github.com/jaypaulb/Canvus-Go-API/canvus`
- Session management with functional options pattern
- Context support for all API operations
- Centralized error handling with typed `APIError` struct
- HTTP client with configurable transport
- Request/response validation with automatic retry logic

#### Authentication
- API key authentication via `WithAPIKey()` option
- Username/password login via `Login()` method
- Token-based authentication with refresh support
- Custom `RoundTripper` support for middleware patterns

#### Canvas Operations
- `ListCanvases()` - List all canvases with optional filtering
- `GetCanvas()` - Get canvas by ID
- `CreateCanvas()` - Create new canvas
- `UpdateCanvas()` - Update canvas properties
- `DeleteCanvas()` - Delete canvas
- `CopyCanvas()` - Copy canvas with all contents
- Background image management (get/set)

#### Widget Operations
- `ListWidgets()` - List widgets on a canvas
- `GetWidget()` - Get widget by ID
- `CreateWidget()` - Create widgets of any type
- `UpdateWidget()` - Update widget properties
- `DeleteWidget()` - Delete widget
- Support for all widget types:
  - Notes
  - Browsers
  - Images
  - PDFs
  - Videos
  - Anchors
  - Connectors

#### User Management
- `ListUsers()` - List all users
- `GetUser()` - Get user by ID
- `CreateUser()` - Create new user
- `UpdateUser()` - Update user properties
- `DeleteUser()` - Delete user
- Access token management (create, list, delete)
- Cleanup utilities for test users

#### Group Operations
- `ListGroups()` - List all groups
- `GetGroup()` - Get group by ID
- `CreateGroup()` - Create new group
- `UpdateGroup()` - Update group properties
- `DeleteGroup()` - Delete group
- Group membership management

#### Folder Operations
- `ListFolders()` - List all folders
- `GetFolder()` - Get folder by ID
- `CreateFolder()` - Create new folder
- `UpdateFolder()` - Update folder properties
- `DeleteFolder()` - Delete folder

#### Workspace Operations
- `ListWorkspaces()` - List workspaces on a canvas
- `GetWorkspace()` - Get workspace by ID
- `CreateWorkspace()` - Create new workspace
- `UpdateWorkspace()` - Update workspace properties
- `DeleteWorkspace()` - Delete workspace

#### Asset Handling
- Image upload and management
- PDF upload and management
- Video upload and management
- Asset file handling for import/export
- Mipmap support

#### System Operations
- `GetLicense()` - Get license information
- `GetSystemInfo()` - Get server system information
- `GetServerConfig()` - Get server configuration
- Audit log access

#### Advanced Features
- **Batch Processing**: `BatchProcessor` with automatic retry logic for bulk operations
- **Import/Export**: Round-trip safe import and export for all widget and asset types
- **Client-side Filtering**: `Filter` abstraction with wildcard and JSONPath-like selectors
- **Cross-canvas Search**: `FindWidgetsAcrossCanvases()` for searching widgets across all canvases
- **Geometry Utilities**: `WidgetsContainId()` and `WidgetsTouchId()` for spatial containment and overlap detection

#### Error Handling
- `APIError` struct with detailed error information:
  - HTTP status code
  - Error code constants
  - Human-readable message
  - Request ID for debugging
  - Error details
  - Wrapped underlying error
- `Unwrap()` and `Is()` methods for error chain handling
- Error code constants for all HTTP status codes and SDK-specific errors

#### Types and Models
- Strongly typed request/response models
- Canvas, Widget, User, Group, Folder types
- Location and Size geometry types
- Color preset types
- Video input/output types
- Comprehensive type coverage for all API entities

#### Testing
- Comprehensive integration tests for all endpoints
- Test helpers and utilities
- Test file fixtures

#### Documentation
- Complete documentation suite in `/docs/`
- Runnable examples in `/examples/`
- Starter templates in `/templates/`
- OpenAPI 3.0 specification
- README with usage examples
- CONTRIBUTING.md guidelines
- Inline godoc comments for all public APIs
- Code examples for common patterns

### Changed
- Nothing (initial release)

### Deprecated
- Nothing (initial release)

### Removed
- Nothing (initial release)

### Fixed
- Nothing (initial release)

### Security
- API key and token handling follow security best practices
- No credentials logged or exposed in error messages

---

## Version History

| Version | Date | Description |
|---------|------|-------------|
| 0.1.0 | 2025-11-19 | Initial pre-release |

---

## Links

- [GitHub Repository](https://github.com/jaypaulb/Canvus-Go-API)
- [Versioning Documentation](docs/VERSIONING.md)
- [Go Package Documentation](https://pkg.go.dev/github.com/jaypaulb/Canvus-Go-API/canvus)

[Unreleased]: https://github.com/jaypaulb/Canvus-Go-API/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/jaypaulb/Canvus-Go-API/releases/tag/v0.1.0
