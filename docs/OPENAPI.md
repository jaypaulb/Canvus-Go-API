# OpenAPI Specification Guide

This document explains how to use the Canvus API OpenAPI specification (`openapi.yaml`) for code generation, API exploration, and building automation tools.

## Overview

The Canvus API OpenAPI specification provides a complete, machine-readable description of the Canvus Server API. It documents:

- **100+ API operations** across Users, Canvases, Widgets, System, and more
- **50+ data schemas** matching the Go SDK types exactly
- **Authentication schemes** including API key and token-based auth
- **Request/response examples** for each major operation

## File Location

```
/openapi.yaml
```

## Specification Version

- **OpenAPI Version**: 3.0.0
- **API Version**: 1.2.0 (matches MTCS API version)

## Use Cases

### 1. Code Generation

Generate client libraries in any language using OpenAPI code generators:

```bash
# Generate Go client
openapi-generator generate -i openapi.yaml -g go -o ./generated/go-client

# Generate Python client
openapi-generator generate -i openapi.yaml -g python -o ./generated/python-client

# Generate TypeScript client
openapi-generator generate -i openapi.yaml -g typescript-axios -o ./generated/ts-client
```

Popular code generators:
- [OpenAPI Generator](https://openapi-generator.tech/) - Multi-language support
- [swagger-codegen](https://github.com/swagger-api/swagger-codegen) - Swagger ecosystem
- [oapi-codegen](https://github.com/deepmap/oapi-codegen) - Go-specific generator

### 2. API Documentation

Generate interactive API documentation:

```bash
# Using Redoc
npx redoc-cli bundle openapi.yaml -o api-docs.html

# Using Swagger UI
docker run -p 8080:8080 -e SWAGGER_JSON=/openapi.yaml \
  -v $(pwd)/openapi.yaml:/openapi.yaml swaggerapi/swagger-ui
```

### 3. API Testing

Import into API testing tools:

- **Postman**: Import > OpenAPI 3.0
- **Insomnia**: Import/Export > OpenAPI
- **Bruno**: Import OpenAPI specification

### 4. MCP Server Foundation

The OpenAPI specification serves as the foundation for building an MCP (Model Context Protocol) server. Use it to:

- Define available tools for AI assistants
- Generate tool schemas from operation definitions
- Map SDK methods to MCP tool invocations

Example MCP tool definition derived from OpenAPI:

```json
{
  "name": "list_canvases",
  "description": "Retrieves all canvases accessible to the authenticated user",
  "inputSchema": {
    "type": "object",
    "properties": {},
    "required": []
  }
}
```

### 5. Validation and Linting

Validate the specification and your API requests:

```bash
# Validate OpenAPI specification
npx @redocly/cli lint openapi.yaml

# Using spectral
npx @stoplight/spectral lint openapi.yaml

# Online validators
# - https://editor.swagger.io/
# - https://apitools.dev/swagger-parser/online/
```

## Schema Overview

The specification defines schemas for all SDK types. Key schemas include:

### Core Resources

| Schema | Description | Go Type |
|--------|-------------|---------|
| `User` | User account information | `canvus.User` |
| `Canvas` | Canvas resource | `canvus.Canvas` |
| `Widget` | Generic widget | `canvus.Widget` |
| `Folder` | Canvas folder | `canvus.Folder` |
| `Group` | User group | `canvus.Group` |

### Widget Types

| Schema | Description | Go Type |
|--------|-------------|---------|
| `Note` | Note widget | `canvus.Note` |
| `Image` | Image widget | `canvus.Image` |
| `PDF` | PDF document widget | `canvus.PDF` |
| `Video` | Video widget | `canvus.Video` |
| `Anchor` | Anchor widget | `canvus.Anchor` |
| `Connector` | Connector widget | `canvus.Connector` |

### Request/Response Types

| Schema | Description |
|--------|-------------|
| `LoginRequest` | Login credentials |
| `LoginResponse` | Login result with token |
| `CreateUserRequest` | New user payload |
| `CreateCanvasRequest` | New canvas payload |
| `CreateNoteRequest` | New note widget payload |
| `APIError` | Error response |

### Geometry Types

| Schema | Description | Go Type |
|--------|-------------|---------|
| `Point` | X/Y coordinates | `canvus.Point` |
| `Size` | Width/Height dimensions | `canvus.Size` |
| `Rectangle` | Position and size | `canvus.Rectangle` |

## Authentication

The API supports two authentication methods:

### API Key Authentication

```yaml
securitySchemes:
  ApiKeyAuth:
    type: apiKey
    in: header
    name: Private-Token
```

Usage:
```bash
curl -H "Private-Token: your-api-key" https://server/api/v1/canvases
```

### Session Token Authentication

1. Login with email/password to get a token
2. Use the token in the `Private-Token` header

## Operation Tags

Operations are organized by tags for easy navigation:

| Tag | Description |
|-----|-------------|
| `Authentication` | Login/logout operations |
| `Users` | User management |
| `Access Tokens` | API token management |
| `Groups` | User group management |
| `Canvases` | Canvas CRUD operations |
| `Canvas Background` | Background settings |
| `Canvas Permissions` | Permission management |
| `Folders` | Folder management |
| `Widgets` | Generic widget operations |
| `Notes` | Note widget operations |
| `Images` | Image widget operations |
| `PDFs` | PDF widget operations |
| `Videos` | Video widget operations |
| `Anchors` | Anchor widget operations |
| `Connectors` | Connector widget operations |
| `Clients` | Client device management |
| `Video I/O` | Video input/output |
| `System` | Server info and config |
| `Audit Log` | Audit events |
| `Assets` | Asset and mipmap operations |

## Examples

The specification includes examples for major operations. Here are some key examples:

### List Canvases Response

```json
[
  {
    "id": "canvas-123",
    "name": "Project Canvas",
    "access": "owner",
    "asset_size": 1024000,
    "created_at": "2024-01-01T00:00:00Z",
    "folder_id": "folder-1",
    "in_trash": false,
    "mode": "normal",
    "modified_at": "2024-01-15T10:00:00Z",
    "preview_hash": "abc123",
    "state": "active"
  }
]
```

### Create Note Request

```json
{
  "widget_type": "note",
  "text": "This is a note",
  "title": "My Note",
  "background_color": "#ffeb3b",
  "location": {
    "x": 100,
    "y": 200
  },
  "size": {
    "width": 300,
    "height": 200
  }
}
```

### Error Response

```json
{
  "code": "not_found",
  "message": "Resource not found",
  "request_id": "req-abc123"
}
```

## Relationship to SDK

The OpenAPI specification mirrors the Canvus Go SDK:

| OpenAPI Operation | SDK Method |
|-------------------|------------|
| `listCanvases` | `Session.ListCanvases()` |
| `getCanvas` | `Session.GetCanvas()` |
| `createCanvas` | `Session.CreateCanvas()` |
| `updateCanvas` | `Session.UpdateCanvas()` |
| `deleteCanvas` | `Session.DeleteCanvas()` |
| `listWidgets` | `Session.ListWidgets()` |
| `createNote` | `Session.CreateNote()` |

Each operation's `operationId` corresponds to the SDK method name (camelCase).

## Validation

### Validating the Specification

The OpenAPI specification has been validated to ensure:

1. **Schema Consistency**: All `$ref` references resolve correctly
2. **Type Correctness**: Schema types match Go struct definitions
3. **Required Fields**: Required fields are properly marked
4. **Examples**: Examples conform to their schemas

To validate locally:

```bash
# Using Redocly CLI
npx @redocly/cli lint openapi.yaml

# Expected output: No errors
```

### Schema to Go Type Mapping

| OpenAPI Type | Go Type |
|--------------|---------|
| `string` | `string` |
| `integer` | `int` |
| `integer (format: int64)` | `int64` |
| `number` | `float64` |
| `boolean` | `bool` |
| `array` | `[]T` |
| `object` | `struct` |

## Extending the Specification

When adding new SDK methods, update the OpenAPI specification:

1. Add the new path under `paths`
2. Create any new schemas under `components/schemas`
3. Add examples for the operation
4. Validate with a linter

Example new operation:

```yaml
/canvases/{canvasId}/widgets/{widgetId}/duplicate:
  post:
    tags:
      - Widgets
    summary: Duplicate widget
    description: Creates a copy of a widget on the same canvas.
    operationId: duplicateWidget
    parameters:
      - $ref: '#/components/parameters/CanvasIdPath'
      - $ref: '#/components/parameters/WidgetIdPath'
    responses:
      '201':
        description: Widget duplicated
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/Widget'
```

## Best Practices

1. **Keep in Sync**: Update OpenAPI when SDK changes
2. **Add Examples**: Include realistic examples for all operations
3. **Document Errors**: Document all possible error responses
4. **Use References**: Use `$ref` for reusable schemas
5. **Validate Regularly**: Run linters before committing changes

## Resources

- [OpenAPI Specification](https://spec.openapis.org/oas/v3.0.0)
- [OpenAPI Generator](https://openapi-generator.tech/)
- [Swagger Editor](https://editor.swagger.io/)
- [Redocly CLI](https://redocly.com/docs/cli/)

## Troubleshooting

### Common Issues

**Reference not found**
```
$ref '#/components/schemas/NonExistent' not found
```
Solution: Check that the referenced schema exists in `components/schemas`.

**Invalid schema type**
```
Invalid type: 'str' is not a valid type
```
Solution: Use valid OpenAPI types: `string`, `integer`, `number`, `boolean`, `array`, `object`.

**Missing required field**
```
Missing required field 'responses' in operation
```
Solution: Every operation must have a `responses` object with at least one response.

## Version History

| Version | Changes |
|---------|---------|
| 1.2.0 | Initial complete specification with all 100+ operations |

---

For questions or issues with the OpenAPI specification, please open an issue on the [GitHub repository](https://github.com/jaypaulb/Canvus-Go-API).
