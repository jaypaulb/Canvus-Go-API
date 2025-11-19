# Testing Guide

This document describes how to run tests for the Canvus Go SDK and explains the test requirements.

## Test Types

### Unit Tests
Currently, most tests in the SDK are **integration tests** that require a live Canvus server. There are no mock-based unit tests at this time.

### Integration Tests
Integration tests make actual API calls to a Canvus server. These tests:
- Create, read, update, and delete real resources
- Verify response validation and error handling
- Test authentication flows
- Validate import/export functionality

## Requirements

### Server Access
All integration tests require access to a Canvus server with:
- Valid API endpoint (e.g., `https://your-server/api/v1/`)
- Admin-level API key with full permissions
- Test user credentials (email/password) for authentication tests

### Configuration File
Create a `settings.json` file in the repository root:

```json
{
  "api_base_url": "https://your-canvus-server/api/v1/",
  "api_key": "your-api-key-here",
  "timeout_seconds": 5,
  "test_user": {
    "username": "test@example.com",
    "password": "password"
  },
  "test_canvas_id": "existing-canvas-uuid-for-widget-tests"
}
```

#### Configuration Fields

| Field | Required | Description |
|-------|----------|-------------|
| `api_base_url` | Yes | Full URL to the Canvus API endpoint |
| `api_key` | Yes | API key with admin permissions |
| `timeout_seconds` | No | Request timeout (default: 5s) |
| `test_user.username` | Yes | Email for login authentication tests |
| `test_user.password` | Yes | Password for login authentication tests |
| `test_canvas_id` | Yes | UUID of an existing canvas for widget tests |

### Permissions Required
The API key needs permissions for:
- User management (create, list, delete users)
- Canvas management (create, list, delete canvases)
- Widget management (create, list, delete widgets)
- System information access
- Audit log access

## Running Tests

### Run All Tests
```bash
go test ./canvus/... -v
```

### Run Specific Test
```bash
go test -run TestFunctionName ./canvus/ -v
```

### Run Tests with Coverage
```bash
go test ./canvus/... -cover
```

### Run Tests with Race Detection
```bash
go test ./canvus/... -race
```

## Test Files

The SDK includes the following test files:

| File | Description |
|------|-------------|
| `users_test.go` | User CRUD operations, login/logout |
| `canvases_test.go` | Canvas CRUD operations |
| `widgets_test.go` | Widget creation and management |
| `notes_test.go` | Note widget operations |
| `images_test.go` | Image widget operations |
| `pdfs_test.go` | PDF widget operations |
| `videos_test.go` | Video widget operations |
| `anchors_test.go` | Anchor widget operations |
| `connectors_test.go` | Connector widget operations |
| `folders_test.go` | Folder operations |
| `groups_test.go` | Group management |
| `batch_test.go` | Batch operations |
| `errors_test.go` | Error handling |
| `export_test.go` | Export functionality |
| `import_test.go` | Import functionality |
| `geometry_test.go` | Geometry utilities |
| `session_test.go` | Session management |
| `clients_test.go` | Client device operations |

## Tests Requiring Server Access

**All tests require server access.** There are currently no offline/mock tests.

### Authentication Tests
- `TestLogin` - Tests login flow with username/password
- `TestLogout` - Tests session logout
- `TestLoginFailure` - Tests invalid credentials handling

### User Tests
- `TestListUsers` - Lists all users
- `TestCreateUser` - Creates and deletes a test user
- `TestGetUser` - Retrieves user by ID
- `TestUpdateUser` - Updates user properties

### Canvas Tests
- `TestListCanvases` - Lists all canvases
- `TestCreateCanvas` - Creates and deletes a test canvas
- `TestGetCanvas` - Retrieves canvas by ID
- `TestCopyCanvas` - Copies a canvas

### Widget Tests
- Require `test_canvas_id` to be set to a valid canvas
- `TestListWidgets` - Lists widgets on test canvas
- `TestCreateNote` - Creates and deletes a note widget
- `TestCreateImage` - Creates image widget (requires asset upload)
- `TestCreatePDF` - Creates PDF widget (requires asset upload)

### Batch Tests
- `TestBatchCreate` - Tests batch widget creation
- `TestBatchDelete` - Tests batch deletion

### Import/Export Tests
- `TestExportCanvas` - Exports canvas to folder
- `TestImportCanvas` - Imports widgets from export

## Test Cleanup

Tests are designed to clean up after themselves:
- Created users are deleted at test end
- Created canvases are **permanently deleted** (not trashed)
- Created widgets are deleted
- Temporary files are removed

### Unique Resource Names
Tests use unique names with timestamps to avoid conflicts:
```go
uniqueName := fmt.Sprintf("test-%s-%d", name, time.Now().UnixNano())
```

## Troubleshooting Tests

### Connection Errors
- Verify `api_base_url` is correct and includes `/api/v1/`
- Check server is accessible from test machine
- Verify SSL certificates are valid

### Authentication Errors
- Verify `api_key` is valid and has admin permissions
- Check `test_user` credentials are correct
- Ensure test user account is not blocked

### Canvas Not Found
- Verify `test_canvas_id` contains a valid canvas UUID
- Ensure the canvas exists and is accessible

### Timeout Errors
- Increase `timeout_seconds` in settings.json
- Check network latency to server

## CI/CD Integration

### GitHub Actions Example
```yaml
name: Tests
on: [push, pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      - name: Create settings.json
        run: |
          echo '${{ secrets.TEST_SETTINGS }}' > settings.json
      - name: Run tests
        run: go test ./canvus/... -v
```

### Required Secrets
For CI/CD, store `settings.json` as a repository secret:
- `TEST_SETTINGS` - Complete JSON configuration

## Future Improvements

Planned testing enhancements:
1. Mock server for offline unit tests
2. Test coverage improvements
3. Performance benchmarks
4. Load testing suite

## Related Documentation

- [Getting Started](GETTING_STARTED.md) - SDK installation and setup
- [Troubleshooting](TROUBLESHOOTING.md) - Common issues and solutions
- [Best Practices](BEST_PRACTICES.md) - Error handling patterns
