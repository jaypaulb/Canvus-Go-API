# Canvus MCP Server Launchpad

**Briefing Document for Development Team**

## Executive Summary

This document outlines how to build a Model Context Protocol (MCP) server for Canvus using the existing Go SDK. The SDK provides 95% of the required functionality - the MCP server is primarily a thin translation layer.

**Estimated Development Time**: 2-3 weeks (with existing SDK)
**Language**: Go
**Primary Dependencies**: Canvus Go SDK, MCP Go SDK

---

## What is an MCP Server?

**MCP (Model Context Protocol)** is Anthropic's standard for connecting AI assistants to external tools and data sources.

An MCP server exposes functionality as "tools" that Claude can use during conversations.

### Example User Experience

**Without MCP:**
```
User: "Create a canvas and add some notes"
Claude: "Here's code you can run to do that..."
User: [copies code, runs it manually]
```

**With MCP:**
```
User: "Create a canvas called 'Sprint Planning' and add notes for tasks"
Claude: [directly executes via MCP]
       "✓ Created canvas 'Sprint Planning' (ID: abc-123)
        ✓ Added note 'Review requirements'
        ✓ Added note 'Design architecture'
        Done! Here's the canvas: [link]"
```

---

## SDK Assets That Accelerate Development

### 1. Complete API Coverage (109+ Methods)

**What we have:**
- All Canvus operations implemented in Go
- Strong typing with request/response structs
- Comprehensive error handling
- Context support throughout

**What this means:**
- **No API integration work needed** - just call SDK methods
- **Type safety** - compile-time checking of tool parameters
- **Error handling** - SDK errors map directly to MCP error responses

### 2. OpenAPI Specification (4,800+ lines)

**File**: `openapi.yaml`

**What we have:**
- Complete API documentation in machine-readable format
- All 109 operations with parameters and schemas
- Request/response examples
- Authentication schemes

**What this means:**
- **Auto-generate tool definitions** from OpenAPI spec
- **Consistent documentation** - tool descriptions come from spec
- **Validation schemas** - parameter validation already defined
- **Reduced boilerplate** - don't manually define each tool

### 3. Production-Ready Patterns

**What we have:**
- Session management with context
- Retry logic and error recovery
- Batch operations with concurrency
- Import/export functionality
- Geometry utilities

**What this means:**
- **Advanced MCP tools** can be built immediately
- **Reliable execution** - SDK handles retries, timeouts
- **Complex operations** like batch processing already implemented

### 4. Comprehensive Examples

**What we have:**
- 10 runnable examples covering all major operations
- 4 starter templates (CLI, web service, batch job, microservice)
- Production-ready code patterns

**What this means:**
- **Copy-paste foundation** - examples show exactly how to use SDK
- **Testing data** - examples provide test cases for MCP tools

---

## MCP Server Architecture (Go-based)

### High-Level Design

```
┌─────────────────────────────────────────────────┐
│  Claude / AI Assistant                          │
└────────────────┬────────────────────────────────┘
                 │ MCP Protocol (stdio/HTTP)
┌────────────────▼────────────────────────────────┐
│  MCP Server (Go)                                │
│  ┌──────────────────────────────────────────┐  │
│  │  Tool Registry                           │  │
│  │  - create_canvas                         │  │
│  │  - create_note                           │  │
│  │  - list_widgets                          │  │
│  │  - batch_move                            │  │
│  │  - export_region                         │  │
│  │  └─ (109+ tools auto-generated)         │  │
│  └──────────────────────────────────────────┘  │
│  ┌──────────────────────────────────────────┐  │
│  │  Tool Executor                           │  │
│  │  - Parameter validation                  │  │
│  │  - SDK method dispatch                   │  │
│  │  - Response formatting                   │  │
│  └──────────────────────────────────────────┘  │
└────────────────┬────────────────────────────────┘
                 │ Go SDK Method Calls
┌────────────────▼────────────────────────────────┐
│  Canvus Go SDK                                  │
│  - Session management                           │
│  - 109+ API methods                             │
│  - Error handling                               │
│  - Retry logic                                  │
└────────────────┬────────────────────────────────┘
                 │ HTTPS/JSON
┌────────────────▼────────────────────────────────┐
│  Canvus API Server                              │
└─────────────────────────────────────────────────┘
```

### Core Components

#### 1. MCP Protocol Handler
- Handles stdio or HTTP transport
- Parses MCP requests
- Returns MCP responses

#### 2. Tool Registry
- Auto-generated from OpenAPI spec
- Maps MCP tool names to SDK methods
- Provides tool schemas for Claude

#### 3. Tool Executor
- Validates tool parameters
- Creates SDK session
- Calls appropriate SDK method
- Formats response for MCP

#### 4. Configuration
- API URL and credentials
- Timeout settings
- Rate limiting

---

## Implementation Roadmap

### Phase 1: Foundation (Week 1)

**Goal**: Basic MCP server with 5-10 core tools

**Tasks**:
1. Set up Go MCP server scaffold
   - Use community MCP Go SDK or implement stdio transport
   - Basic request/response handling

2. Implement tool registry
   - Manually define 5-10 core tools (create_canvas, create_note, list_canvases, etc.)
   - Map to SDK methods

3. Implement tool executor
   - Parameter validation
   - SDK session creation with config
   - Error handling and formatting

4. End-to-end test
   - Test with Claude Desktop or MCP Inspector
   - Verify tools execute correctly

**Deliverables**:
- Working MCP server with 5-10 tools
- Basic configuration (API key, URL)
- Test suite

### Phase 2: Auto-Generation (Week 2)

**Goal**: Auto-generate all 109+ tools from OpenAPI spec

**Tasks**:
1. OpenAPI parser
   - Parse `openapi.yaml`
   - Extract operations, parameters, schemas

2. Tool generator
   - Generate tool definitions from OpenAPI operations
   - Generate parameter schemas
   - Generate descriptions

3. Dynamic dispatcher
   - Route tool calls to SDK methods by name
   - Handle different parameter patterns (path params, body, query)

4. Testing
   - Validate all generated tools
   - Test parameter validation
   - Test error responses

**Deliverables**:
- Auto-generated tool registry (109+ tools)
- Tool testing framework
- Updated documentation

### Phase 3: Advanced Features (Week 3)

**Goal**: Production-ready with advanced capabilities

**Tasks**:
1. Resource support
   - Expose canvases as MCP resources
   - Expose widgets as MCP resources
   - Resource templates

2. Advanced tools
   - Batch operations (move, copy, delete multiple widgets)
   - Import/export operations
   - Geometry queries (find widgets in region)

3. Performance optimization
   - Connection pooling
   - Caching (canvas/widget metadata)
   - Rate limiting

4. Documentation and deployment
   - User guide
   - Installation instructions
   - Docker containerization

**Deliverables**:
- Production-ready MCP server
- Docker image
- Complete documentation
- Deployment guide

---

## Code Examples

### Example 1: Manual Tool Definition

```go
package main

import (
    "context"
    "encoding/json"
    "github.com/jaypaulb/Canvus-Go-API/canvus"
)

// Tool definition for create_canvas
type CreateCanvasTool struct {
    session *canvus.Session
}

// Schema returns the MCP tool schema
func (t *CreateCanvasTool) Schema() map[string]interface{} {
    return map[string]interface{}{
        "name": "create_canvas",
        "description": "Create a new canvas",
        "inputSchema": map[string]interface{}{
            "type": "object",
            "properties": map[string]interface{}{
                "name": map[string]interface{}{
                    "type": "string",
                    "description": "Canvas name",
                },
                "folder_id": map[string]interface{}{
                    "type": "string",
                    "description": "Folder ID to create canvas in",
                },
            },
            "required": []string{"name"},
        },
    }
}

// Execute runs the tool
func (t *CreateCanvasTool) Execute(ctx context.Context, params map[string]interface{}) (interface{}, error) {
    // Extract parameters
    name, ok := params["name"].(string)
    if !ok {
        return nil, fmt.Errorf("name is required")
    }

    folderID, _ := params["folder_id"].(string)

    // Call SDK method
    req := canvus.CreateCanvasRequest{
        Name:     name,
        FolderID: folderID,
    }

    canvas, err := t.session.CreateCanvas(ctx, req)
    if err != nil {
        return nil, fmt.Errorf("failed to create canvas: %w", err)
    }

    // Return result
    return map[string]interface{}{
        "id":   canvas.ID,
        "name": canvas.Name,
    }, nil
}
```

### Example 2: Auto-Generated from OpenAPI

```go
package main

import (
    "context"
    "fmt"
    "gopkg.in/yaml.v3"
)

// OpenAPISpec represents the parsed OpenAPI specification
type OpenAPISpec struct {
    Paths map[string]map[string]Operation `yaml:"paths"`
}

type Operation struct {
    OperationID string                 `yaml:"operationId"`
    Summary     string                 `yaml:"summary"`
    Parameters  []Parameter            `yaml:"parameters"`
    RequestBody *RequestBody           `yaml:"requestBody"`
}

// GenerateTools generates MCP tools from OpenAPI spec
func GenerateTools(spec *OpenAPISpec) []Tool {
    var tools []Tool

    for path, methods := range spec.Paths {
        for method, op := range methods {
            tool := Tool{
                Name:        op.OperationID,
                Description: op.Summary,
                InputSchema: generateSchema(op),
                Execute: func(ctx context.Context, params map[string]interface{}) (interface{}, error) {
                    return executeSDKMethod(ctx, method, path, params)
                },
            }
            tools = append(tools, tool)
        }
    }

    return tools
}

// executeSDKMethod dispatches to the appropriate SDK method
func executeSDKMethod(ctx context.Context, method, path string, params map[string]interface{}) (interface{}, error) {
    // Parse path and method to determine SDK function
    // Example: POST /canvases -> session.CreateCanvas()
    //          GET /canvases/{id} -> session.GetCanvas(id)

    switch {
    case method == "POST" && path == "/canvases":
        return executeCreateCanvas(ctx, params)
    case method == "GET" && path == "/canvases":
        return executeListCanvases(ctx, params)
    // ... auto-generate all 109 cases
    default:
        return nil, fmt.Errorf("unknown operation: %s %s", method, path)
    }
}
```

### Example 3: MCP Server Main

```go
package main

import (
    "context"
    "encoding/json"
    "fmt"
    "os"
    "github.com/jaypaulb/Canvus-Go-API/canvus"
)

func main() {
    // Load configuration
    apiURL := os.Getenv("CANVUS_API_URL")
    apiKey := os.Getenv("CANVUS_API_KEY")

    // Create SDK session
    session := canvus.NewSession(apiURL, canvus.WithAPIKey(apiKey))

    // Load OpenAPI spec and generate tools
    tools := loadAndGenerateTools("openapi.yaml", session)

    // Start MCP server (stdio transport)
    server := NewMCPServer(tools)
    server.Start()
}

// MCPServer handles MCP protocol
type MCPServer struct {
    tools map[string]Tool
}

func (s *MCPServer) Start() {
    // Read JSON-RPC requests from stdin
    // Process tool calls
    // Write responses to stdout
    decoder := json.NewDecoder(os.Stdin)
    encoder := json.NewEncoder(os.Stdout)

    for {
        var req MCPRequest
        if err := decoder.Decode(&req); err != nil {
            break
        }

        resp := s.handleRequest(req)
        encoder.Encode(resp)
    }
}

func (s *MCPServer) handleRequest(req MCPRequest) MCPResponse {
    switch req.Method {
    case "tools/list":
        return s.listTools()
    case "tools/call":
        return s.callTool(req.Params)
    default:
        return MCPResponse{Error: "unknown method"}
    }
}
```

---

## SDK-to-MCP Mapping Patterns

### Pattern 1: Simple CRUD Operations

**SDK Method**:
```go
canvas, err := session.GetCanvas(ctx, canvasID)
```

**MCP Tool**:
```json
{
  "name": "get_canvas",
  "description": "Get canvas by ID",
  "inputSchema": {
    "type": "object",
    "properties": {
      "canvas_id": {"type": "string"}
    },
    "required": ["canvas_id"]
  }
}
```

### Pattern 2: List/Search Operations

**SDK Method**:
```go
filter := &canvus.Filter{Criteria: map[string]interface{}{"name": "Project*"}}
canvases, err := session.ListCanvases(ctx, filter)
```

**MCP Tool**:
```json
{
  "name": "list_canvases",
  "description": "List canvases with optional filtering",
  "inputSchema": {
    "type": "object",
    "properties": {
      "name_filter": {"type": "string", "description": "Wildcard filter for name"}
    }
  }
}
```

### Pattern 3: Complex Operations

**SDK Method**:
```go
results := session.BatchProcessor.ProcessCanvases(ctx, canvasIDs, "move", targetFolderID)
```

**MCP Tool**:
```json
{
  "name": "batch_move_canvases",
  "description": "Move multiple canvases to a folder",
  "inputSchema": {
    "type": "object",
    "properties": {
      "canvas_ids": {"type": "array", "items": {"type": "string"}},
      "target_folder_id": {"type": "string"}
    },
    "required": ["canvas_ids", "target_folder_id"]
  }
}
```

### Pattern 4: Asset Operations

**SDK Method**:
```go
data, err := session.DownloadImage(ctx, canvasID, imageID)
```

**MCP Tool** (returns base64):
```json
{
  "name": "download_image",
  "description": "Download image widget as base64",
  "inputSchema": {
    "type": "object",
    "properties": {
      "canvas_id": {"type": "string"},
      "image_id": {"type": "string"}
    },
    "required": ["canvas_id", "image_id"]
  }
}
```

---

## Resource Mapping (Optional)

MCP also supports "resources" - read-only data that Claude can access.

### Example: Canvas Resource

```go
// Resource definition
type CanvasResource struct {
    URI         string // canvus://canvas/{id}
    Name        string
    Description string
    MimeType    string // application/json
}

// Fetch canvas data
func (r *CanvasResource) Read(ctx context.Context) ([]byte, error) {
    canvas, err := session.GetCanvas(ctx, r.ID)
    if err != nil {
        return nil, err
    }
    return json.Marshal(canvas)
}
```

**Use Case**: Claude can read canvas metadata without explicit tool calls:
- "Read the canvus://canvas/abc-123 resource and summarize it"

---

## Testing Strategy

### Unit Tests
- Test tool parameter validation
- Test SDK method calls with mocks
- Test error handling

### Integration Tests
- Test against live Canvus server
- Verify all 109 tools execute correctly
- Test error scenarios (invalid params, auth failures, etc.)

### MCP Protocol Tests
- Use MCP Inspector tool
- Test with Claude Desktop
- Verify tool schemas are valid

### Example Test

```go
func TestCreateCanvasTool(t *testing.T) {
    // Setup
    mockSession := &MockCanvusSession{}
    tool := &CreateCanvasTool{session: mockSession}

    // Execute
    params := map[string]interface{}{
        "name": "Test Canvas",
        "folder_id": "folder-123",
    }

    result, err := tool.Execute(context.Background(), params)

    // Assert
    assert.NoError(t, err)
    assert.Equal(t, "Test Canvas", result["name"])

    // Verify SDK was called correctly
    assert.Equal(t, "Test Canvas", mockSession.LastCreateRequest.Name)
}
```

---

## Configuration

### Environment Variables

```bash
# Required
CANVUS_API_URL=https://your-server/api/v1
CANVUS_API_KEY=your-api-key

# Optional
CANVUS_TIMEOUT=30s
CANVUS_MAX_RETRIES=3
CANVUS_LOG_LEVEL=info
```

### config.json (alternative)

```json
{
  "canvus": {
    "api_url": "https://your-server/api/v1",
    "api_key": "your-api-key",
    "timeout": "30s",
    "max_retries": 3
  },
  "mcp": {
    "transport": "stdio",
    "log_level": "info"
  }
}
```

---

## Deployment Options

### 1. Local Development (stdio)
```bash
# Add to Claude Desktop config
{
  "mcpServers": {
    "canvus": {
      "command": "/path/to/canvus-mcp-server",
      "env": {
        "CANVUS_API_URL": "...",
        "CANVUS_API_KEY": "..."
      }
    }
  }
}
```

### 2. Docker Container
```dockerfile
FROM golang:1.21-alpine
WORKDIR /app
COPY . .
RUN go build -o canvus-mcp-server
CMD ["./canvus-mcp-server"]
```

### 3. HTTP Transport (for remote clients)
```go
// Serve MCP over HTTP instead of stdio
http.HandleFunc("/mcp", mcpHandler)
http.ListenAndServe(":8080", nil)
```

---

## Key Advantages of This Approach

### 1. Minimal Boilerplate
- **SDK already implements 109 operations** - no API integration needed
- **OpenAPI spec provides schemas** - auto-generate tool definitions
- **Examples provide test cases** - copy-paste testing approach

### 2. Type Safety
- **Go's strong typing** - compile-time parameter checking
- **SDK request/response structs** - validated types throughout

### 3. Robustness
- **SDK handles retries** - transient failures handled automatically
- **Context support** - timeouts and cancellation built-in
- **Error handling** - SDK errors map cleanly to MCP errors

### 4. Maintainability
- **Single source of truth** - OpenAPI spec drives both SDK and MCP
- **Auto-generation** - new SDK methods automatically become tools
- **Consistency** - MCP tools match SDK exactly

---

## Estimated Effort Breakdown

| Phase | Tasks | Effort |
|-------|-------|--------|
| **Week 1** | MCP server scaffold, 5-10 manual tools, testing | 40 hours |
| **Week 2** | OpenAPI parser, auto-generation, full tool set | 40 hours |
| **Week 3** | Resources, advanced tools, docs, deployment | 40 hours |
| **Total** | End-to-end MCP server with 109+ tools | **120 hours** |

**Without SDK**: Estimated 400-500 hours (implement all API operations + MCP layer)

**Time Saved**: ~75% reduction in development time

---

## Next Steps

### Immediate Actions
1. **Review this document** - Discuss with team, clarify questions
2. **Set up development environment** - Install Go, clone SDK repo
3. **Explore MCP protocol** - Review Anthropic's MCP documentation
4. **Test SDK locally** - Run examples, verify API access

### Week 1 Goals
1. **Choose MCP Go library** - Evaluate options or implement stdio transport
2. **Scaffold MCP server** - Basic request/response handling
3. **Implement 5-10 tools** - Start with simple operations (list, get, create)
4. **Test with Claude Desktop** - Verify end-to-end functionality

### Technical Decisions Needed
1. **MCP Go library** - Use community SDK or implement custom?
2. **Transport** - stdio (simpler) or HTTP (more flexible)?
3. **Configuration** - Environment variables or config file?
4. **Error handling** - How detailed should MCP error responses be?
5. **Rate limiting** - Implement in MCP layer or rely on SDK?

---

## Resources

### Canvus Go SDK
- **Repository**: https://github.com/jaypaulb/Canvus-Go-API
- **Documentation**: `/docs/` directory
- **Examples**: `/examples/` directory
- **OpenAPI Spec**: `openapi.yaml`
- **Installation**: `go get github.com/jaypaulb/Canvus-Go-API/canvus@v0.1.0`

### MCP Protocol
- **Specification**: https://modelcontextprotocol.io/
- **Go Examples**: Search GitHub for "mcp server go"
- **Testing Tools**: MCP Inspector, Claude Desktop

### Development Tools
- **Go**: 1.21+
- **Git**: For version control
- **Docker**: For containerization
- **Claude Desktop**: For testing MCP integration

---

## Questions for Team Discussion

1. **MCP Transport**: stdio or HTTP? (stdio simpler for local, HTTP better for remote)
2. **Tool Granularity**: Expose all 109 operations or curate a smaller set?
3. **Resource Support**: Implement MCP resources for canvases/widgets?
4. **Authentication**: How should MCP server get Canvus credentials? (env vars, config file, user input?)
5. **Caching**: Should MCP server cache canvas/widget metadata?
6. **Versioning**: How to handle SDK updates? (rebuild MCP server, automated?)

---

## Success Metrics

### Week 1
- MCP server responds to tool list request
- 5-10 tools execute successfully
- Test suite covers core functionality

### Week 2
- All 109 SDK operations exposed as tools
- Auto-generation from OpenAPI working
- Full integration test suite passing

### Week 3
- Production deployment ready
- Documentation complete
- Performance benchmarks met (< 100ms overhead per tool call)

---

## Conclusion

The Canvus Go SDK provides an exceptional foundation for rapid MCP server development:

✅ **109 operations implemented** - no API integration work
✅ **OpenAPI spec complete** - auto-generate tool definitions
✅ **Production patterns ready** - retry, error handling, batching
✅ **Comprehensive examples** - copy-paste test cases

**Expected timeline**: 2-3 weeks from start to production-ready MCP server.

**Key success factor**: Leverage SDK maximally, minimize custom code.

---

**Document Version**: 1.0
**Date**: 2025-11-19
**Author**: Development Team
**Status**: Ready for implementation
