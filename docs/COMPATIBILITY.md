# Canvus Go SDK Compatibility Guide

This document describes version compatibility between the Canvus Go SDK and the MTCS (Multi-Touch Collaboration System) API, known quirks and workarounds, and the process for handling breaking changes.

## Version Compatibility Matrix

The SDK version scheme follows the pattern `MTCS_MAJOR.MTCS_MINOR.SDK_PATCH`. See [VERSIONING.md](VERSIONING.md) for details.

| SDK Version | MTCS API Version | Status | Notes |
|-------------|------------------|--------|-------|
| v0.1.x | v1.x | **Current** | Initial release, full API coverage |

### Understanding the Matrix

- **SDK Version**: The version of this Go SDK
- **MTCS API Version**: The version of the Canvus/MTCS server API
- **Status**: Current, Supported, Deprecated, or EOL (End of Life)
- **Notes**: Additional compatibility information

## Current MTCS API Version Support

**SDK v0.1.0** targets **MTCS API v1.x**

### Supported Features

The SDK provides full coverage of the MTCS API v1.x, including:

- **Authentication**: API key (Private-Token header), username/password login, token refresh
- **Users**: Full CRUD, access token management, group membership
- **Canvases**: Full CRUD, permissions, move/copy, preview, demo state
- **Folders**: Full CRUD, permissions, move/copy, trash
- **Widgets**: All widget types (Notes, Images, PDFs, Videos, Anchors, Connectors, Browsers, Video Inputs)
- **Assets**: Upload, download, mipmap support
- **System**: Server info, config, license, audit log
- **Clients & Workspaces**: Client management, workspace operations, video I/O
- **Batch Operations**: Bulk delete, move, copy, pin/unpin with retry logic
- **Import/Export**: Round-trip safe widget and asset transfer

## Known Quirks and Workarounds

The MTCS API has some behaviors that require special handling. This section documents known quirks and the SDK's approach to handling them.

### 1. Widget Type Case Sensitivity

**Quirk**: The MTCS API returns widget types in different cases depending on the endpoint (e.g., `"Note"` vs `"note"`).

**Workaround**: The SDK normalizes widget type comparisons using case-insensitive matching. When creating widgets, use lowercase types (`"note"`, `"browser"`, `"image"`).

```go
// SDK handles case-insensitive widget type comparison internally
widgets, _ := session.ListWidgets(ctx, canvasID, nil)
for _, w := range widgets {
    // WidgetType may be "Note" or "note" from API
    // Use strings.EqualFold for comparisons
    if strings.EqualFold(w.WidgetType, "note") {
        // Process note widget
    }
}
```

### 2. Numeric Field Type Variations

**Quirk**: The API sometimes returns numeric values as integers and sometimes as floats (e.g., coordinates in `location` and `size` fields).

**Workaround**: The SDK's response validation uses numeric equality comparison that treats `int` and `float64` values as equal if they represent the same number.

```go
// The SDK handles this internally - you can work with float64 consistently
location := widget.Location
x := location["x"].(float64) // Safe - SDK normalizes numeric types
```

### 3. Server-Generated Fields

**Quirk**: Several fields are generated or modified by the server and may differ from request values:

- `id` - Always server-generated
- `created_at`, `modified_at` - Timestamps
- `state`, `access` - Server-controlled states
- `preview_hash`, `asset_size` - Computed values
- `folder_id`, `parent_id` - May be transformed
- `location`, `size` - May be adjusted for constraints

**Workaround**: The SDK's response validation skips these fields when comparing request/response. Don't rely on exact matches for these fields.

### 4. Empty Response Bodies on DELETE

**Quirk**: DELETE operations may return empty bodies or bodies without the expected `status: "deleted"` field.

**Workaround**: The SDK validates DELETE responses flexibly - it accepts empty bodies and looks for either `status` or `state` fields containing "deleted" (case-insensitive).

### 5. Token Expiration and Refresh

**Quirk**: Tokens have limited lifetimes and the server returns 401 without explicit "token expired" messaging.

**Workaround**: The SDK automatically attempts token refresh on first 401 response. For long-running operations, use the `TokenStore` interface to persist tokens:

```go
type FileTokenStore struct {
    path string
}

func (s *FileTokenStore) GetToken() (string, error) {
    data, err := os.ReadFile(s.path)
    return string(data), err
}

func (s *FileTokenStore) StoreToken(token string, expiry time.Time) error {
    return os.WriteFile(s.path, []byte(token), 0600)
}

func (s *FileTokenStore) ClearToken() error {
    return os.Remove(s.path)
}

cfg := canvus.DefaultSessionConfig()
cfg.BaseURL = "https://server/api/v1"
session := canvus.NewSession(cfg, canvus.WithTokenStore(&FileTokenStore{"/tmp/token"}))
```

### 6. Certificate Verification

**Quirk**: Many Canvus deployments use self-signed certificates.

**Workaround**: The SDK's `WithAPIKey` option creates an HTTP client with `InsecureSkipVerify: true` by default. For production, configure a custom HTTP client with proper certificate handling:

```go
cfg := canvus.DefaultSessionConfig()
cfg.BaseURL = "https://server/api/v1"

// Load custom CA certificate
caCert, _ := os.ReadFile("/path/to/ca.crt")
caCertPool := x509.NewCertPool()
caCertPool.AppendCertsFromPEM(caCert)

cfg.HTTPClient = &http.Client{
    Transport: &http.Transport{
        TLSClientConfig: &tls.Config{
            RootCAs: caCertPool,
        },
    },
}
session := canvus.NewSession(cfg, canvus.WithAPIKey("your-key"))
```

### 7. Pagination Not Implemented

**Quirk**: The MTCS API supports pagination parameters, but many endpoints return all results regardless.

**Workaround**: Treat all List operations as returning complete result sets. Use client-side filtering for large datasets:

```go
// Fetch all canvases
canvases, _ := session.ListCanvases(ctx, nil)

// Client-side filtering
filter := &canvus.Filter{Criteria: map[string]interface{}{"name": "Project*"}}
filtered := canvus.FilterSlice(canvases, filter)
```

### 8. Rate Limiting Responses

**Quirk**: The server may return 429 (Too Many Requests) under heavy load.

**Workaround**: The SDK has built-in retry logic with exponential backoff for 429 responses. Configure retry behavior:

```go
cfg := canvus.DefaultSessionConfig()
cfg.MaxRetries = 5
cfg.RetryWaitMin = 500 * time.Millisecond
cfg.RetryWaitMax = 30 * time.Second
session := canvus.NewSession(cfg, canvus.WithAPIKey("your-key"))
```

### 9. Asset Upload Content Types

**Quirk**: Asset uploads require multipart/form-data with specific field names for different widget types.

**Workaround**: Use the SDK's specialized upload methods which handle content type and field naming:

```go
// Use CreateImage, CreatePDF, CreateVideo instead of generic CreateWidget
image, err := session.CreateImage(ctx, canvasID, multipartBody, contentType)
```

### 10. Workspace Selector Types

**Quirk**: Workspaces can be selected by index (0-based) or ID, requiring different API endpoints.

**Workaround**: Use the `WorkspaceSelector` type which handles both cases:

```go
// By index
selector := canvus.WorkspaceSelector{Index: intPtr(0)}

// By ID
selector := canvus.WorkspaceSelector{ID: "workspace-uuid"}
```

## Import/Export Fidelity Guarantees

The SDK's import/export functionality maintains full fidelity for:

- **Widget Types**: Notes, Images, PDFs, Videos, Anchors, Connectors, Browsers, Video Inputs
- **Spatial Data**: Location and size preserved exactly
- **Parent/Child Relationships**: Parent IDs remapped correctly
- **Assets**: Binary files (images, PDFs, videos) exported and re-uploaded
- **Widget Properties**: All editable properties preserved

### Limitations

- **Server-generated fields** (IDs, timestamps) will differ after import
- **Canvas-specific settings** (background, permissions) are not included in widget export
- **Connector endpoints** are remapped based on widget ID mapping

## Deprecated Features

Currently, there are no deprecated features in SDK v0.1.x.

As the SDK evolves, deprecated features will be listed here with:
- Deprecation version
- Removal target version
- Migration path

## Breaking Change Notice Process

When breaking changes are necessary, we follow this process:

### 1. Advance Notice

Breaking changes are announced:
- In the CHANGELOG.md at least one minor version before removal
- In godoc deprecation comments on affected functions
- In GitHub release notes

### 2. Deprecation Period

Deprecated features:
- Continue to work for at least one minor version
- Log warnings when used (where practical)
- Have documented migration paths

### 3. Removal

After the deprecation period:
- Features are removed in a new minor or major version
- CHANGELOG.md documents all breaking changes
- Migration guide provided in release notes

### Example Timeline

```
v0.2.0 - OldMethod marked deprecated, NewMethod added
v0.3.0 - OldMethod removed
```

## Plan for MTCS API Updates

### Monitoring for API Changes

We monitor for MTCS API changes through:
- Canvus release notes
- API response changes in integration tests
- Community issue reports

### SDK Update Process

When MTCS API changes:

1. **Patch Updates** (v0.1.x) for:
   - Bug fixes
   - New optional fields
   - Documentation improvements

2. **Minor Updates** (v0.x.0) for:
   - New API endpoints
   - New required parameters
   - Deprecated endpoint alternatives

3. **Major Updates** (vX.0.0) for:
   - Breaking API changes
   - Removed endpoints
   - Incompatible authentication changes

### Compatibility Testing

Each SDK release is tested against:
- Target MTCS API version
- Common deployment configurations
- All documented use cases

## Reporting Compatibility Issues

If you encounter compatibility issues:

1. Check this document for known quirks
2. Check the [TROUBLESHOOTING.md](TROUBLESHOOTING.md) guide
3. Open a GitHub issue with:
   - SDK version
   - MTCS server version
   - Minimal reproduction code
   - Expected vs actual behavior

## Related Documentation

- [VERSIONING.md](VERSIONING.md) - Version scheme explanation
- [TROUBLESHOOTING.md](TROUBLESHOOTING.md) - Common issues and solutions
- [BEST_PRACTICES.md](BEST_PRACTICES.md) - Recommended patterns
- [API_REFERENCE.md](API_REFERENCE.md) - Complete method reference
