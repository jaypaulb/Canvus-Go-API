# Canvus Go SDK API Reference

This document provides an organized reference to all SDK methods, types, and utilities. For detailed parameter documentation and examples, see the [pkg.go.dev documentation](https://pkg.go.dev/github.com/jaypaulb/Canvus-Go-API/canvus).

## Package Documentation

The complete API documentation is available on pkg.go.dev:

**[pkg.go.dev/github.com/jaypaulb/Canvus-Go-API/canvus](https://pkg.go.dev/github.com/jaypaulb/Canvus-Go-API/canvus)**

## API Overview

The SDK provides 109+ methods organized into logical categories. All methods:
- Accept `context.Context` as the first parameter for cancellation and timeout support
- Return typed results and `error`
- Use strongly-typed request/response structs where applicable

## Session Creation and Authentication

### Session Constructors

| Function | Description |
|----------|-------------|
| `NewSession(cfg *SessionConfig, opts ...SessionConfigOption) *Session` | Create session with full configuration |
| `NewDefaultSession(baseURL string) *Session` | Create session with defaults (no auth) |
| `NewSessionFromConfig(baseURL, apiKey string) *Session` | Create session with API key |
| `DefaultSessionConfig() *SessionConfig` | Get default configuration |

### Configuration Options

| Option | Description |
|--------|-------------|
| `WithAPIKey(apiKey string)` | Configure API key authentication |
| `WithHTTPClient(client *http.Client)` | Set custom HTTP client |
| `WithMaxRetries(n int)` | Set maximum retry attempts |
| `WithRetryWait(min, max time.Duration)` | Set retry backoff bounds |
| `WithRequestTimeout(timeout time.Duration)` | Set request timeout |
| `WithUserAgent(ua string)` | Set custom User-Agent |
| `WithCircuitBreaker(maxFailures int, resetTimeout time.Duration)` | Configure circuit breaker |
| `WithTokenStore(store TokenStore)` | Set token persistence store |
| `WithTokenRefreshThreshold(threshold time.Duration)` | Set token refresh timing |

### Authentication Methods

| Method | Description |
|--------|-------------|
| `Login(ctx, email, password string) error` | Authenticate with username/password |
| `Logout(ctx) error` | Invalidate current token |
| `UserID() int64` | Get authenticated user's ID |

---

## Users

User management including CRUD operations, access tokens, and group membership.

### User CRUD

| Method | Description |
|--------|-------------|
| `ListUsers(ctx) ([]User, error)` | List all users |
| `GetUser(ctx, id int64) (*User, error)` | Get user by ID |
| `CreateUser(ctx, req CreateUserRequest) (*User, error)` | Create new user |
| `UpdateUser(ctx, id int64, req UpdateUserRequest) (*User, error)` | Update user |
| `DeleteUser(ctx, id int64) error` | Delete user |

### Access Tokens

| Method | Description |
|--------|-------------|
| `ListAccessTokens(ctx, userID int64) ([]AccessToken, error)` | List user's access tokens |
| `GetAccessToken(ctx, userID int64, tokenID string) (*AccessToken, error)` | Get specific token |
| `CreateAccessToken(ctx, userID int64, req CreateAccessTokenRequest) (*AccessToken, error)` | Create new token |
| `DeleteAccessToken(ctx, userID int64, tokenID string) error` | Delete token |

### Groups

| Method | Description |
|--------|-------------|
| `ListGroups(ctx) ([]Group, error)` | List all groups |
| `GetGroup(ctx, id int) (*Group, error)` | Get group by ID |
| `CreateGroup(ctx, req CreateGroupRequest) (*Group, error)` | Create new group |
| `DeleteGroup(ctx, id int) error` | Delete group |
| `ListGroupMembers(ctx, groupID int) ([]GroupMember, error)` | List group members |
| `AddUserToGroup(ctx, groupID int, userID int) error` | Add user to group |
| `RemoveUserFromGroup(ctx, groupID int, userID int) error` | Remove user from group |

---

## Canvases

Canvas management including CRUD, permissions, and special operations.

### Canvas CRUD

| Method | Description |
|--------|-------------|
| `ListCanvases(ctx, filter *Filter) ([]Canvas, error)` | List canvases with optional filter |
| `GetCanvas(ctx, id string) (*Canvas, error)` | Get canvas by ID |
| `CreateCanvas(ctx, req CreateCanvasRequest) (*Canvas, error)` | Create new canvas |
| `UpdateCanvas(ctx, id string, req UpdateCanvasRequest) (*Canvas, error)` | Update canvas |
| `DeleteCanvas(ctx, id string) error` | Permanently delete canvas |

### Canvas Operations

| Method | Description |
|--------|-------------|
| `MoveCanvas(ctx, id string, req MoveOrCopyCanvasRequest) (*Canvas, error)` | Move to folder |
| `CopyCanvas(ctx, id string, req MoveOrCopyCanvasRequest) (*Canvas, error)` | Copy to folder |
| `TrashCanvas(ctx, id string, _ string) (*Canvas, error)` | Move to trash |
| `GetCanvasPreview(ctx, id string) ([]byte, error)` | Download preview image |
| `RestoreDemoCanvas(ctx, id string) error` | Restore demo state |
| `SaveDemoState(ctx, id string) error` | Save current as demo state |

### Canvas Permissions

| Method | Description |
|--------|-------------|
| `GetCanvasPermissions(ctx, id string) (*CanvasPermissions, error)` | Get permission overrides |
| `SetCanvasPermissions(ctx, id string, perms CanvasPermissions) (*CanvasPermissions, error)` | Set permissions |

### Canvas Background

| Method | Description |
|--------|-------------|
| `GetCanvasBackground(ctx, canvasID string) (*CanvasBackground, error)` | Get background settings |
| `PatchCanvasBackground(ctx, canvasID string, req interface{}) error` | Update background |
| `PostCanvasBackground(ctx, canvasID string, multipartBody interface{}) error` | Upload background image |

---

## Folders

Folder management for organizing canvases.

### Folder CRUD

| Method | Description |
|--------|-------------|
| `ListFolders(ctx) ([]Folder, error)` | List all folders |
| `GetFolder(ctx, id string) (*Folder, error)` | Get folder by ID |
| `CreateFolder(ctx, req CreateFolderRequest) (*Folder, error)` | Create new folder |
| `RenameFolder(ctx, id string, name string) (*Folder, error)` | Rename folder |
| `DeleteFolder(ctx, id string) error` | Delete folder |
| `DeleteFolderContents(ctx, id string) error` | Delete folder contents only |

### Folder Operations

| Method | Description |
|--------|-------------|
| `MoveFolder(ctx, id string, parentID string, conflicts string) (*Folder, error)` | Move folder |
| `CopyFolder(ctx, id string, parentID string, conflicts string) (*Folder, error)` | Copy folder |
| `TrashFolder(ctx, id string, _ string) (*Folder, error)` | Move to trash |

### Folder Permissions

| Method | Description |
|--------|-------------|
| `GetFolderPermissions(ctx, id string) (*FolderPermissions, error)` | Get permissions |
| `SetFolderPermissions(ctx, id string, perms FolderPermissions) (*FolderPermissions, error)` | Set permissions |

---

## Widgets

Widget management covering all widget types and operations.

### Generic Widget Operations

| Method | Description |
|--------|-------------|
| `ListWidgets(ctx, canvasID string, filter *Filter) ([]Widget, error)` | List widgets with optional filter |
| `GetWidget(ctx, canvasID, widgetID string) (*Widget, error)` | Get widget by ID |
| `CreateWidget(ctx, canvasID string, req interface{}, contentType ...string) (*Widget, error)` | Create generic widget |
| `UpdateWidget(ctx, canvasID, widgetID string, req map[string]interface{}) (*Widget, error)` | Update widget properties |
| `DeleteWidget(ctx, canvasID, widgetID, widgetType string) error` | Delete widget |
| `PatchParentID(ctx, canvasID, widgetID, parentID string) (*Widget, error)` | Change widget parent |
| `CopyWidget(ctx, widgetID, targetCanvasID string) error` | Copy to another canvas |
| `MoveWidget(ctx, widgetID, targetCanvasID string) error` | Move to another canvas |
| `PinWidget(ctx, widgetID string) error` | Pin widget |
| `UnpinWidget(ctx, widgetID string) error` | Unpin widget |

### Notes

| Method | Description |
|--------|-------------|
| `ListNotes(ctx, canvasID string) ([]Note, error)` | List notes |
| `GetNote(ctx, canvasID, noteID string) (*Note, error)` | Get note |
| `CreateNote(ctx, canvasID string, req interface{}) (*Note, error)` | Create note |
| `UpdateNote(ctx, canvasID, noteID string, req interface{}) (*Note, error)` | Update note |
| `DeleteNote(ctx, canvasID, noteID string) error` | Delete note |
| `UploadNote(ctx, canvasID string, multipartBody interface{}) (*Note, error)` | Upload note with file |

### Images

| Method | Description |
|--------|-------------|
| `ListImages(ctx, canvasID string) ([]Image, error)` | List images |
| `GetImage(ctx, canvasID, imageID string) (*Image, error)` | Get image metadata |
| `CreateImage(ctx, canvasID string, multipartBody io.Reader, contentType string) (*Image, error)` | Upload image |
| `UpdateImage(ctx, canvasID, imageID string, req interface{}) (*Image, error)` | Update image |
| `DeleteImage(ctx, canvasID, imageID string) error` | Delete image |
| `DownloadImage(ctx, canvasID, imageID string) ([]byte, error)` | Download image data |

### PDFs

| Method | Description |
|--------|-------------|
| `ListPDFs(ctx, canvasID string) ([]PDF, error)` | List PDFs |
| `GetPDF(ctx, canvasID, pdfID string) (*PDF, error)` | Get PDF metadata |
| `CreatePDF(ctx, canvasID string, multipartBody interface{}, contentType string) (*PDF, error)` | Upload PDF |
| `UpdatePDF(ctx, canvasID, pdfID string, req interface{}) (*PDF, error)` | Update PDF |
| `DeletePDF(ctx, canvasID, pdfID string) error` | Delete PDF |
| `DownloadPDF(ctx, canvasID, pdfID string) ([]byte, error)` | Download PDF data |

### Videos

| Method | Description |
|--------|-------------|
| `ListVideos(ctx, canvasID string) ([]Video, error)` | List videos |
| `GetVideo(ctx, canvasID, videoID string) (*Video, error)` | Get video metadata |
| `CreateVideo(ctx, canvasID string, multipartBody interface{}, contentType string) (*Video, error)` | Upload video |
| `UpdateVideo(ctx, canvasID, videoID string, req interface{}) (*Video, error)` | Update video |
| `DeleteVideo(ctx, canvasID, videoID string) error` | Delete video |
| `DownloadVideo(ctx, canvasID, videoID string) ([]byte, error)` | Download video data |

### Anchors

| Method | Description |
|--------|-------------|
| `ListAnchors(ctx, canvasID string) ([]Anchor, error)` | List anchors |
| `GetAnchor(ctx, canvasID, anchorID string) (*Anchor, error)` | Get anchor |
| `CreateAnchor(ctx, canvasID string, req interface{}) (*Anchor, error)` | Create anchor |
| `UpdateAnchor(ctx, canvasID, anchorID string, req interface{}) (*Anchor, error)` | Update anchor |
| `DeleteAnchor(ctx, canvasID, anchorID string) error` | Delete anchor |

### Connectors

| Method | Description |
|--------|-------------|
| `ListConnectors(ctx, canvasID string) ([]Connector, error)` | List connectors |
| `GetConnector(ctx, canvasID, connectorID string) (*Connector, error)` | Get connector |
| `CreateConnector(ctx, canvasID string, req interface{}) (*Connector, error)` | Create connector |
| `UpdateConnector(ctx, canvasID, connectorID string, req interface{}) (*Connector, error)` | Update connector |
| `DeleteConnector(ctx, canvasID, connectorID string) error` | Delete connector |

### Video Inputs

| Method | Description |
|--------|-------------|
| `ListVideoInputs(ctx, canvasID string) ([]VideoInput, error)` | List video inputs on canvas |
| `CreateVideoInput(ctx, canvasID string, req interface{}) (*VideoInput, error)` | Create video input widget |
| `DeleteVideoInput(ctx, canvasID, inputID string) error` | Delete video input |

### Color Presets

| Method | Description |
|--------|-------------|
| `ListColorPresets(ctx, canvasID string) ([]ColorPreset, error)` | List color presets |
| `GetColorPresets(ctx, canvasID string) (*ColorPresets, error)` | Get all presets |
| `GetColorPreset(ctx, canvasID, name string) (*ColorPreset, error)` | Get preset by name |
| `CreateColorPreset(ctx, canvasID string, req interface{}) (*ColorPreset, error)` | Create preset |
| `UpdateColorPreset(ctx, canvasID, name string, req interface{}) (*ColorPreset, error)` | Update preset |
| `DeleteColorPreset(ctx, canvasID, name string) error` | Delete preset |
| `PatchColorPresets(ctx, canvasID string, req *ColorPresets) (*ColorPresets, error)` | Bulk update presets |

---

## Assets and Mipmaps

Asset management for binary content and image mipmaps.

| Method | Description |
|--------|-------------|
| `UploadAsset(ctx, canvasID string, multipartBody interface{}) (*Asset, error)` | Upload asset file |
| `GetAssetByHash(ctx, canvasID, publicHashHex string) ([]byte, error)` | Download asset by hash |
| `GetMipmapInfo(ctx, canvasID, publicHashHex string, page *int) (*MipmapInfo, error)` | Get mipmap metadata |
| `GetMipmapLevel(ctx, canvasID, publicHashHex string, level int, page *int) ([]byte, error)` | Download mipmap level |

---

## System

Server configuration, information, and administration.

### Server Info

| Method | Description |
|--------|-------------|
| `GetServerInfo(ctx) (*ServerInfo, error)` | Get server version and status |
| `GetLicenseInfo(ctx) (*LicenseInfo, error)` | Get license information |

### Server Configuration

| Method | Description |
|--------|-------------|
| `GetServerConfig(ctx) (*ServerConfig, error)` | Get server configuration |
| `UpdateServerConfig(ctx, req ServerConfig) (*ServerConfig, error)` | Update configuration |
| `SendTestEmail(ctx) error` | Send test email |

### Audit Log

| Method | Description |
|--------|-------------|
| `ListAuditEvents(ctx, opts *AuditLogOptions) ([]AuditEvent, error)` | Query audit events |

---

## Clients and Workspaces

Client device and workspace management.

### Clients

| Method | Description |
|--------|-------------|
| `ListClients(ctx) ([]ClientInfo, error)` | List connected clients |
| `GetClient(ctx, id string) (*ClientInfo, error)` | Get client info |
| `CreateClient(ctx, req CreateClientRequest) (*ClientInfo, error)` | Create client |
| `UpdateClient(ctx, id string, req UpdateClientRequest) (*ClientInfo, error)` | Update client |
| `DeleteClient(ctx, id string) error` | Delete client |

### Workspaces

| Method | Description |
|--------|-------------|
| `ListWorkspaces(ctx, clientID string) ([]Workspace, error)` | List client's workspaces |
| `GetWorkspace(ctx, clientID string, selector WorkspaceSelector) (*Workspace, error)` | Get workspace |
| `UpdateWorkspace(ctx, clientID string, selector WorkspaceSelector, req UpdateWorkspaceRequest) (*Workspace, error)` | Update workspace |
| `OpenCanvasOnWorkspace(ctx, clientID string, selector WorkspaceSelector, opts OpenCanvasOptions) error` | Open canvas |
| `ToggleWorkspacePinned(ctx, clientID string, selector WorkspaceSelector) error` | Toggle pinned state |
| `ToggleWorkspaceInfoPanel(ctx, clientID string, selector WorkspaceSelector) error` | Toggle info panel |

### Video I/O

| Method | Description |
|--------|-------------|
| `ListClientVideoInputs(ctx, clientID string) ([]VideoInputSource, error)` | List video inputs |
| `ListVideoOutputs(ctx, clientID string) ([]VideoOutput, error)` | List video outputs |
| `SetVideoOutputSource(ctx, clientID string, index int, req interface{}) error` | Set output source |
| `UpdateVideoOutput(ctx, canvasID, outputID string, req interface{}) (*VideoOutput, error)` | Update output |

---

## Batch Operations

Bulk operations with automatic retry and progress tracking.

### Batch Processor

```go
// Create batch processor
bp := canvus.NewBatchProcessor(session, canvus.DefaultBatchConfig())

// Build operations
builder := canvus.NewBatchOperationBuilder()
builder.Delete("widget-1", &widget1)
builder.Move("widget-2", &widget2, targetFolder)
operations := builder.Build()

// Execute batch
results, err := bp.ExecuteBatch(ctx, operations)
summary := canvus.Summarize(results)
```

| Function | Description |
|----------|-------------|
| `NewBatchProcessor(session *Session, config *BatchConfig) *BatchProcessor` | Create processor |
| `DefaultBatchConfig() *BatchConfig` | Get default config |
| `NewBatchOperationBuilder() *BatchOperationBuilder` | Create operation builder |
| `Summarize(results []*BatchResult) *BatchSummary` | Get results summary |

### BatchOperationBuilder Methods

| Method | Description |
|--------|-------------|
| `Delete(id string, resource interface{})` | Add delete operation |
| `Move(id string, resource interface{}, targetFolderID string)` | Add move operation |
| `Copy(id string, resource interface{}, targetCanvasID string)` | Add copy operation |
| `Pin(id string, widget *Widget)` | Add pin operation |
| `Unpin(id string, widget *Widget)` | Add unpin operation |
| `Build() []*BatchOperation` | Get operation list |

---

## Import/Export

Round-trip safe widget and asset export/import.

| Method | Description |
|--------|-------------|
| `ExportWidgetsToFolder(ctx, canvasID string, widgetIDs []string, region Rectangle, sharedCanvasID string, baseFolder string) (string, error)` | Export widgets to folder |
| `ImportWidgetsToRegion(ctx, canvasID string, exported *ExportedWidgetSet, targetRegion Rectangle) ([]string, error)` | Import widgets to canvas |

---

## Geometry Utilities

Spatial analysis and widget relationship utilities.

| Function | Description |
|----------|-------------|
| `WidgetsContainId(ctx, session *Session, canvasID string, widgetID string, widget *Widget, tolerance float64) (WidgetZone, error)` | Find widgets contained within another |
| `Contains(a, b Rectangle) bool` | Check if rectangle a contains b |
| `Touches(a, b Rectangle) bool` | Check if rectangles overlap |
| `WidgetContains(a, b Widget) bool` | Check if widget a contains b |
| `WidgetsTouch(a, b Widget) bool` | Check if widgets overlap |
| `WidgetBoundingBox(w Widget) Rectangle` | Get widget bounding box |

---

## Search Utilities

Cross-canvas widget search with pattern matching.

| Function | Description |
|----------|-------------|
| `FindWidgetsAcrossCanvases(ctx, lister WidgetsLister, query map[string]interface{}) ([]WidgetMatch, error)` | Search widgets across all canvases |

### Query Patterns

- Exact match: `"field": "value"`
- Wildcard any: `"field": "*"`
- Prefix match: `"field": "abc*"`
- Suffix match: `"field": "*xyz"`
- Contains: `"field": "*mid*"`
- Nested field: `"$.location.x": 100`

---

## Filtering

Client-side filtering for list results.

| Function | Description |
|----------|-------------|
| `FilterSlice[T Filterable](elems []T, filter *Filter) []T` | Filter slice of filterable items |

### Filter Usage

```go
filter := &canvus.Filter{
    Criteria: map[string]interface{}{
        "name": "Project*",           // Prefix match
        "widget_type": "note",        // Exact match
        "$.location.x": 100.0,        // Nested field
    },
}
widgets, _ := session.ListWidgets(ctx, canvasID, filter)
```

---

## Error Handling

Typed errors with error codes and wrapping support.

### Error Types

| Type | Description |
|------|-------------|
| `APIError` | HTTP API error with status code, code, and message |
| `ValidationError` | Single field validation error |
| `ValidationErrors` | Collection of validation errors |

### APIError Methods

| Method | Description |
|--------|-------------|
| `Error() string` | Get error message |
| `Is(target error) bool` | Check error equality |
| `Unwrap() error` | Get wrapped error |
| `Wrap(err error) *APIError` | Wrap another error |
| `WithDetails(details map[string]interface{}) *APIError` | Add details |
| `WithRequestID(requestID string) *APIError` | Add request ID |

### Error Functions

| Function | Description |
|----------|-------------|
| `NewAPIError(statusCode int, code ErrorCode, message string) *APIError` | Create API error |
| `ErrorFromStatus(statusCode int, message string) error` | Create from HTTP status |
| `ParseErrorResponse(statusCode int, body []byte) *APIError` | Parse error response |
| `IsRetryableError(err error) bool` | Check if retryable |
| `IsContextError(err error) bool` | Check if context error |
| `WrapError(err error, msg string) error` | Wrap with message |
| `WrapErrorf(err error, format string, args ...interface{}) error` | Wrap with format |

---

## Common Patterns

### CRUD Pattern

All resource types follow a consistent pattern:

```go
// List
items, err := session.ListWidgets(ctx, canvasID, nil)

// Get
item, err := session.GetWidget(ctx, canvasID, widgetID)

// Create
item, err := session.CreateNote(ctx, canvasID, req)

// Update
item, err := session.UpdateNote(ctx, canvasID, noteID, req)

// Delete
err := session.DeleteNote(ctx, canvasID, noteID)
```

### Request Types

Create and Update operations use typed request structs:

```go
req := canvus.CreateCanvasRequest{
    Name:     "My Canvas",
    FolderID: "folder-id",
}
canvas, err := session.CreateCanvas(ctx, req)
```

Or `interface{}` for flexible widget properties:

```go
req := map[string]interface{}{
    "text":     "Note content",
    "location": map[string]interface{}{"x": 100, "y": 200},
    "size":     map[string]interface{}{"width": 300, "height": 200},
}
note, err := session.CreateNote(ctx, canvasID, req)
```

### Context Usage

All methods accept context for cancellation and timeout:

```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

canvases, err := session.ListCanvases(ctx, nil)
```

---

## Related Documentation

- [GETTING_STARTED.md](GETTING_STARTED.md) - Quick start guide
- [BEST_PRACTICES.md](BEST_PRACTICES.md) - Recommended patterns
- [TROUBLESHOOTING.md](TROUBLESHOOTING.md) - Common issues
- [EXAMPLES.md](EXAMPLES.md) - Runnable examples index
- [COMPATIBILITY.md](COMPATIBILITY.md) - Version compatibility
