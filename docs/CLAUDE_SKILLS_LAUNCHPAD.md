# Canvus Claude Skills Launchpad

**Briefing Document for Development Team**

## Executive Summary

This document outlines how to build Claude Skills for Canvus using the existing Go SDK. Skills provide a more context-efficient alternative to MCP for large APIs, loading capabilities on-demand rather than all 109 tools upfront.

**Estimated Development Time**: 2-3 weeks (with existing SDK)
**Language**: Go
**Platform**: Claude (Anthropic-specific)
**Primary Dependencies**: Canvus Go SDK

---

## What are Claude Skills?

**Claude Skills** are Anthropic's approach to giving Claude specific capabilities that can be loaded on-demand and mixed-and-matched as needed.

### Key Difference from MCP

**MCP Approach** (Model Context Protocol):
```
Context Window:
├─ Tool 1: create_canvas (definition: 50 tokens)
├─ Tool 2: get_canvas (definition: 50 tokens)
├─ Tool 3: delete_canvas (definition: 50 tokens)
├─ ... (106 more tools)
└─ Total: ~5,000-10,000 tokens for tool definitions

User: "Create a canvas"
Claude: [Already has all 109 tools loaded, uses 1]
```

**Skills Approach**:
```
Context Window:
└─ (empty, waiting for skill request)

User: "Create a canvas"
Claude: "I need the Canvas Management skill"
System: [Loads Canvas Management skill: 15 operations, ~500 tokens]
Claude: [Uses create_canvas from loaded skill]
```

### Why Skills Are Better for Large APIs

| Aspect | MCP (109 tools) | Skills (8 grouped) |
|--------|-----------------|-------------------|
| **Context Cost** | ~8,000 tokens upfront | ~500-800 tokens per skill |
| **Loading** | All tools always loaded | Load on-demand |
| **Discoverability** | 109 flat tools | 8 logical categories |
| **Updates** | Replace entire toolset | Update individual skills |
| **Versioning** | Single version | Per-skill versioning |

**For Canvus**: With 109 operations, Skills save ~85% context vs MCP

---

## Canvus Skills Architecture

### Proposed Skill Groupings

Based on logical operation clusters in the SDK:

#### 1. Canvas Management (15 operations)
**Purpose**: Create, organize, and manage canvases

**Operations**:
- Create/get/update/delete canvas
- List canvases (with filtering)
- Move canvas to folder
- Copy canvas
- Get/set canvas permissions
- Get/set canvas background
- Manage canvas preview

**When loaded**: User needs to work with canvases

**Example prompts**:
- "Create a new canvas called 'Project Planning'"
- "List all canvases in the Marketing folder"
- "Copy this canvas to another folder"

---

#### 2. Widget Creation (18 operations)
**Purpose**: Add content to canvases

**Operations**:
- Create note widget
- Create image widget (with upload)
- Create PDF widget (with upload)
- Create video widget (with upload)
- Create anchor widget
- Create connector widget
- Create browser widget

**When loaded**: User wants to add content

**Example prompts**:
- "Add a note with text 'Meeting notes' to the canvas"
- "Upload this image to the canvas"
- "Create a connector between these two widgets"

---

#### 3. Widget Organization (12 operations)
**Purpose**: Arrange and manage widgets on canvases

**Operations**:
- List widgets (with filtering)
- Get widget details
- Update widget properties
- Delete widget
- Move widget
- Resize widget
- Set widget parent
- Pin/unpin widget
- Copy widget

**When loaded**: User needs to organize content

**Example prompts**:
- "Move all notes to the left side of the canvas"
- "Pin this widget to the top"
- "Delete all widgets in this region"

---

#### 4. User Administration (12 operations)
**Purpose**: Manage users, groups, and access tokens

**Operations**:
- List/create/update/delete users
- Create/delete access tokens
- List/create/delete groups
- Add/remove users from groups
- Set user properties (admin, blocked)

**When loaded**: User needs admin capabilities

**Example prompts**:
- "Create a new user account for john@company.com"
- "Generate an API token for the CI system"
- "Add Sarah to the Editors group"

---

#### 5. Content Search (8 operations)
**Purpose**: Find and query content

**Operations**:
- Search widgets across canvases
- Filter widgets by type/properties
- Find widgets in geometric region
- Find widgets containing specific widget
- Find widgets touching specific widget
- List widgets by criteria

**When loaded**: User needs to find content

**Example prompts**:
- "Find all note widgets containing 'TODO'"
- "Show me all widgets in the top-left corner"
- "Which widgets are inside this anchor?"

---

#### 6. Batch Operations (6 operations)
**Purpose**: Perform bulk operations efficiently

**Operations**:
- Batch move widgets
- Batch copy widgets
- Batch delete widgets
- Batch pin widgets
- Batch unpin widgets
- Process batch with progress tracking

**When loaded**: User needs to operate on multiple items

**Example prompts**:
- "Move all these widgets to another canvas"
- "Delete all notes created yesterday"
- "Copy these 50 widgets to a new canvas"

---

#### 7. Import/Export (4 operations)
**Purpose**: Migrate content between canvases/systems

**Operations**:
- Export widgets from region
- Export widgets by ID list
- Import widgets to canvas
- Verify import/export fidelity

**When loaded**: User needs to move content

**Example prompts**:
- "Export this region of the canvas to a folder"
- "Import widgets from this export folder"
- "Copy all widgets from Canvas A to Canvas B"

---

#### 8. System Administration (8 operations)
**Purpose**: Configure and monitor the Canvus server

**Operations**:
- Get server info
- Get/update server config
- Get license info
- List audit events
- Send test email
- Get folder permissions
- Set folder permissions

**When loaded**: User needs system admin access

**Example prompts**:
- "What version is the Canvus server running?"
- "Show me recent audit log entries"
- "Update the SMTP settings"

---

## Skills vs MCP: Direct Comparison

### Scenario: "Create a canvas and add 3 notes"

**MCP Approach**:
```
1. Load all 109 tools (~8,000 tokens)
2. User request
3. Claude uses: create_canvas, create_note (3x)
4. Total tools used: 4 out of 109
5. Wasted context: 105 unused tool definitions
```

**Skills Approach**:
```
1. User request
2. Claude determines needed skills: Canvas Management + Widget Creation
3. Load 2 skills (~1,000 tokens for 33 operations)
4. Claude uses: create_canvas, create_note (3x)
5. Context efficiency: 87% savings vs MCP
```

### When to Use Each

**Use MCP when**:
- Simple APIs (< 20 operations)
- All tools frequently used
- Need cross-platform compatibility (non-Claude LLMs)
- MCP infrastructure already exists

**Use Claude Skills when**:
- Large APIs (50+ operations)
- Operations cluster into logical groups
- Most conversations use subset of operations
- Claude-only deployment is acceptable
- Context efficiency matters (long conversations)

**For Canvus**: Skills are clearly better (109 ops, clustered groups, context-sensitive)

---

## Implementation Architecture

### High-Level Design

```
┌─────────────────────────────────────────────────┐
│  Claude Assistant                               │
└────────────────┬────────────────────────────────┘
                 │ Determines needed skill
┌────────────────▼────────────────────────────────┐
│  Skill Loader                                   │
│  ┌──────────────────────────────────────────┐  │
│  │  Skill Registry                          │  │
│  │  - Canvas Management                     │  │
│  │  - Widget Creation                       │  │
│  │  - Widget Organization                   │  │
│  │  - User Administration                   │  │
│  │  - Content Search                        │  │
│  │  - Batch Operations                      │  │
│  │  - Import/Export                         │  │
│  │  - System Administration                 │  │
│  └──────────────────────────────────────────┘  │
└────────────────┬────────────────────────────────┘
                 │ Load requested skill
┌────────────────▼────────────────────────────────┐
│  Skill Executor (Go)                            │
│  - Operation dispatcher                         │
│  - Parameter validation                         │
│  - SDK session management                       │
│  - Result formatting                            │
└────────────────┬────────────────────────────────┘
                 │ SDK method calls
┌────────────────▼────────────────────────────────┐
│  Canvus Go SDK (109+ methods)                   │
└────────────────┬────────────────────────────────┘
                 │ HTTPS/JSON
┌────────────────▼────────────────────────────────┐
│  Canvus API Server                              │
└─────────────────────────────────────────────────┘
```

### Core Components

#### 1. Skill Definitions (YAML/JSON)
Each skill is a file describing:
- Skill name and description
- Operations included
- Parameter schemas
- Examples

#### 2. Skill Loader
- Reads skill definitions
- Loads requested skills into Claude's context
- Manages skill versions

#### 3. Skill Executor (Go Server)
- Receives operation requests from Claude
- Validates parameters
- Calls SDK methods
- Returns results

#### 4. SDK Integration
- Same Canvus Go SDK used for MCP
- No API integration work needed
- All 109 operations already implemented

---

## Implementation Roadmap

### Phase 1: Foundation (Week 1)

**Goal**: Define skills and implement 2 core skills

**Tasks**:
1. **Skill Definition Format**
   - Choose format (YAML or JSON)
   - Define schema for skill definitions
   - Create template

2. **Implement Canvas Management Skill**
   - Define skill YAML
   - Map 15 operations to SDK methods
   - Create skill executor logic
   - Test with Claude

3. **Implement Widget Creation Skill**
   - Define skill YAML
   - Map 18 operations to SDK methods
   - Test interaction with Canvas Management

4. **Skill Loader**
   - Load skills on-demand
   - Inject into Claude's context
   - Handle skill dependencies

**Deliverables**:
- Skill definition format documented
- 2 working skills (33 operations)
- Skill loader functional
- Test suite

---

### Phase 2: Complete Skill Set (Week 2)

**Goal**: Implement remaining 6 skills

**Tasks**:
1. **Implement Widget Organization Skill** (12 ops)
2. **Implement User Administration Skill** (12 ops)
3. **Implement Content Search Skill** (8 ops)
4. **Implement Batch Operations Skill** (6 ops)
5. **Implement Import/Export Skill** (4 ops)
6. **Implement System Administration Skill** (8 ops)

**Testing**:
- Test each skill independently
- Test skill combinations
- Test skill loading/unloading
- Verify context efficiency

**Deliverables**:
- All 8 skills implemented (109 operations)
- Comprehensive test suite
- Skill interaction tests

---

### Phase 3: Production Polish (Week 3)

**Goal**: Production-ready skills system

**Tasks**:
1. **Performance Optimization**
   - Skill caching
   - Lazy loading
   - Connection pooling

2. **Error Handling**
   - Graceful skill failures
   - Fallback mechanisms
   - Clear error messages

3. **Documentation**
   - User guide for each skill
   - Skill capability matrix
   - Example conversations
   - Troubleshooting guide

4. **Deployment**
   - Docker containerization
   - Configuration management
   - Monitoring and logging

**Deliverables**:
- Production-ready skills system
- Complete documentation
- Deployment guide
- Performance benchmarks

---

## Code Examples

### Example 1: Skill Definition (YAML)

**File**: `skills/canvas-management.yaml`

```yaml
name: Canvas Management
description: Create, organize, and manage Canvus canvases
version: 1.0.0

operations:
  - name: create_canvas
    description: Create a new canvas
    sdk_method: CreateCanvas
    parameters:
      - name: name
        type: string
        required: true
        description: Canvas name
      - name: folder_id
        type: string
        required: false
        description: Folder to create canvas in
    example:
      input:
        name: "Project Planning"
        folder_id: "folder-123"
      output:
        id: "canvas-456"
        name: "Project Planning"

  - name: list_canvases
    description: List all canvases with optional filtering
    sdk_method: ListCanvases
    parameters:
      - name: filter
        type: object
        required: false
        description: Filter criteria (name, folder_id, etc.)
    example:
      input:
        filter:
          name: "Project*"
      output:
        - id: "canvas-456"
          name: "Project Planning"
        - id: "canvas-789"
          name: "Project Review"

  # ... 13 more operations
```

---

### Example 2: Skill Loader (Go)

```go
package skills

import (
    "context"
    "fmt"
    "gopkg.in/yaml.v3"
    "os"
)

// Skill represents a Claude skill
type Skill struct {
    Name        string      `yaml:"name"`
    Description string      `yaml:"description"`
    Version     string      `yaml:"version"`
    Operations  []Operation `yaml:"operations"`
}

type Operation struct {
    Name        string      `yaml:"name"`
    Description string      `yaml:"description"`
    SDKMethod   string      `yaml:"sdk_method"`
    Parameters  []Parameter `yaml:"parameters"`
}

// LoadSkill loads a skill definition from file
func LoadSkill(path string) (*Skill, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, fmt.Errorf("failed to read skill file: %w", err)
    }

    var skill Skill
    if err := yaml.Unmarshal(data, &skill); err != nil {
        return nil, fmt.Errorf("failed to parse skill: %w", err)
    }

    return &skill, nil
}

// SkillRegistry manages available skills
type SkillRegistry struct {
    skills map[string]*Skill
}

func NewSkillRegistry() *SkillRegistry {
    return &SkillRegistry{
        skills: make(map[string]*Skill),
    }
}

// Register adds a skill to the registry
func (r *SkillRegistry) Register(skill *Skill) {
    r.skills[skill.Name] = skill
}

// GetSkill retrieves a skill by name
func (r *SkillRegistry) GetSkill(name string) (*Skill, error) {
    skill, ok := r.skills[name]
    if !ok {
        return nil, fmt.Errorf("skill not found: %s", name)
    }
    return skill, nil
}

// ListSkills returns all available skills
func (r *SkillRegistry) ListSkills() []*Skill {
    skills := make([]*Skill, 0, len(r.skills))
    for _, skill := range r.skills {
        skills = append(skills, skill)
    }
    return skills
}
```

---

### Example 3: Skill Executor (Go)

```go
package skills

import (
    "context"
    "fmt"
    "github.com/jaypaulb/Canvus-Go-API/canvus"
)

// SkillExecutor executes skill operations
type SkillExecutor struct {
    session  *canvus.Session
    registry *SkillRegistry
}

func NewSkillExecutor(session *canvus.Session, registry *SkillRegistry) *SkillExecutor {
    return &SkillExecutor{
        session:  session,
        registry: registry,
    }
}

// Execute runs an operation from a skill
func (e *SkillExecutor) Execute(ctx context.Context, skillName, opName string, params map[string]interface{}) (interface{}, error) {
    // Get skill
    skill, err := e.registry.GetSkill(skillName)
    if err != nil {
        return nil, err
    }

    // Find operation
    var op *Operation
    for i := range skill.Operations {
        if skill.Operations[i].Name == opName {
            op = &skill.Operations[i]
            break
        }
    }
    if op == nil {
        return nil, fmt.Errorf("operation not found: %s", opName)
    }

    // Dispatch to SDK method
    return e.dispatch(ctx, op.SDKMethod, params)
}

// dispatch maps SDK method names to actual SDK calls
func (e *SkillExecutor) dispatch(ctx context.Context, method string, params map[string]interface{}) (interface{}, error) {
    switch method {
    case "CreateCanvas":
        return e.createCanvas(ctx, params)
    case "ListCanvases":
        return e.listCanvases(ctx, params)
    case "GetCanvas":
        return e.getCanvas(ctx, params)
    // ... 106 more cases
    default:
        return nil, fmt.Errorf("unknown SDK method: %s", method)
    }
}

// createCanvas executes the CreateCanvas SDK method
func (e *SkillExecutor) createCanvas(ctx context.Context, params map[string]interface{}) (interface{}, error) {
    name, ok := params["name"].(string)
    if !ok {
        return nil, fmt.Errorf("name is required")
    }

    folderID, _ := params["folder_id"].(string)

    req := canvus.CreateCanvasRequest{
        Name:     name,
        FolderID: folderID,
    }

    canvas, err := e.session.CreateCanvas(ctx, req)
    if err != nil {
        return nil, fmt.Errorf("failed to create canvas: %w", err)
    }

    return map[string]interface{}{
        "id":        canvas.ID,
        "name":      canvas.Name,
        "folder_id": canvas.FolderID,
        "created":   canvas.CreatedAt,
    }, nil
}
```

---

### Example 4: Skill Auto-Generation from OpenAPI

```go
package skills

import (
    "fmt"
    "gopkg.in/yaml.v3"
    "os"
)

// OperationGroup represents a logical grouping of operations
type OperationGroup struct {
    Name       string
    PathPrefix string
    Methods    []string
}

var skillGroups = []OperationGroup{
    {
        Name:       "Canvas Management",
        PathPrefix: "/canvases",
        Methods:    []string{"GET", "POST", "PATCH", "DELETE"},
    },
    {
        Name:       "Widget Creation",
        PathPrefix: "/canvases/{id}/widgets",
        Methods:    []string{"POST"},
    },
    // ... more groups
}

// GenerateSkillsFromOpenAPI creates skill definitions from OpenAPI spec
func GenerateSkillsFromOpenAPI(openapiPath string) ([]*Skill, error) {
    // Parse OpenAPI spec
    spec, err := parseOpenAPISpec(openapiPath)
    if err != nil {
        return nil, err
    }

    var skills []*Skill

    // For each skill group
    for _, group := range skillGroups {
        skill := &Skill{
            Name:       group.Name,
            Version:    "1.0.0",
            Operations: []Operation{},
        }

        // Find matching operations in OpenAPI spec
        for path, pathItem := range spec.Paths {
            if matchesGroup(path, group) {
                for method, op := range pathItem {
                    if contains(group.Methods, method) {
                        operation := Operation{
                            Name:        op.OperationID,
                            Description: op.Summary,
                            SDKMethod:   operationIDToSDKMethod(op.OperationID),
                            Parameters:  extractParameters(op),
                        }
                        skill.Operations = append(skill.Operations, operation)
                    }
                }
            }
        }

        skills = append(skills, skill)
    }

    return skills, nil
}

// SaveSkill writes a skill definition to YAML
func SaveSkill(skill *Skill, path string) error {
    data, err := yaml.Marshal(skill)
    if err != nil {
        return err
    }
    return os.WriteFile(path, data, 0644)
}
```

---

## Context Efficiency Analysis

### Token Breakdown

**MCP (All 109 Tools Loaded)**:
```
Tool definitions:
- create_canvas: 80 tokens
- get_canvas: 75 tokens
- list_canvases: 90 tokens
- ... (106 more)
Total: ~8,500 tokens

User conversation: 1,000 tokens
Tools used: 3 operations
Context wasted: 106 unused tool definitions (~8,000 tokens)
```

**Skills (Load 2 Skills)**:
```
Canvas Management skill:
- Skill header: 50 tokens
- 15 operations × 75 tokens = 1,125 tokens
- Total: 1,175 tokens

Widget Creation skill:
- Skill header: 50 tokens
- 18 operations × 75 tokens = 1,350 tokens
- Total: 1,400 tokens

Total loaded: 2,575 tokens (33 operations)
User conversation: 1,000 tokens
Context efficiency: 70% savings vs MCP
```

### Real-World Scenarios

#### Scenario 1: Canvas Creation Workflow
**User**: "Create a project canvas and add meeting notes"

**MCP**: Load 109 tools (8,500 tokens)
**Skills**: Load Canvas Management + Widget Creation (2,575 tokens)
**Savings**: 69% context

---

#### Scenario 2: User Administration
**User**: "Create 5 new user accounts and assign them to groups"

**MCP**: Load 109 tools (8,500 tokens)
**Skills**: Load User Administration (900 tokens)
**Savings**: 89% context

---

#### Scenario 3: Complex Multi-Operation
**User**: "Find all notes, move them to a new canvas, and export the result"

**MCP**: Load 109 tools (8,500 tokens)
**Skills**: Load Content Search + Widget Organization + Canvas Management + Import/Export (4,200 tokens)
**Savings**: 51% context

**Average savings across scenarios**: ~70% context efficiency

---

## SDK Acceleration for Skills

### Why the SDK Makes This Fast

#### 1. Zero API Integration Work
**Without SDK**: ~200 hours to implement 109 operations
**With SDK**: 0 hours (already done)

#### 2. Type Safety Built-In
**Without SDK**: Define and validate 109 parameter schemas
**With SDK**: Use existing Go structs

#### 3. Error Handling Complete
**Without SDK**: Implement retry logic, error classification
**With SDK**: Already implemented with typed errors

#### 4. OpenAPI Spec Ready
**Without SDK**: Document 109 operations manually
**With SDK**: Parse existing `openapi.yaml`

### Development Time Comparison

| Task | Without SDK | With SDK | Time Saved |
|------|-------------|----------|------------|
| API Integration | 200 hours | 0 hours | 200 hours |
| Error Handling | 40 hours | 0 hours | 40 hours |
| Type Definitions | 60 hours | 0 hours | 60 hours |
| Testing | 80 hours | 20 hours | 60 hours |
| Skill Definitions | 40 hours | 20 hours | 20 hours |
| Skill Executor | 20 hours | 10 hours | 10 hours |
| Documentation | 30 hours | 10 hours | 20 hours |
| **Total** | **470 hours** | **60 hours** | **410 hours** |

**Time savings**: 87% reduction

---

## Testing Strategy

### Skill-Level Tests

```go
func TestCanvasManagementSkill(t *testing.T) {
    // Load skill
    skill, err := LoadSkill("skills/canvas-management.yaml")
    assert.NoError(t, err)

    // Verify operations
    assert.Equal(t, 15, len(skill.Operations))
    assert.Contains(t, operationNames(skill), "create_canvas")

    // Test skill execution
    executor := NewSkillExecutor(testSession, registry)
    result, err := executor.Execute(
        context.Background(),
        "Canvas Management",
        "create_canvas",
        map[string]interface{}{
            "name": "Test Canvas",
        },
    )

    assert.NoError(t, err)
    assert.NotNil(t, result["id"])
}
```

### Integration Tests

```go
func TestSkillCombinations(t *testing.T) {
    // Test Canvas Management + Widget Creation
    // 1. Create canvas
    canvas := executeOp(t, "Canvas Management", "create_canvas", ...)

    // 2. Add widgets to canvas
    note := executeOp(t, "Widget Creation", "create_note", ...)

    // 3. Verify
    assert.Equal(t, canvas["id"], note["canvas_id"])
}
```

### Context Efficiency Tests

```go
func TestContextUsage(t *testing.T) {
    // Measure context tokens for MCP
    mcpTokens := measureMCPContextTokens()

    // Measure context tokens for Skills
    skillsTokens := measureSkillsContextTokens([]string{
        "Canvas Management",
        "Widget Creation",
    })

    // Verify efficiency
    savings := (mcpTokens - skillsTokens) / mcpTokens
    assert.Greater(t, savings, 0.6) // At least 60% savings
}
```

---

## Deployment

### Local Development

**Directory Structure**:
```
canvus-skills/
├── skills/
│   ├── canvas-management.yaml
│   ├── widget-creation.yaml
│   ├── widget-organization.yaml
│   ├── user-administration.yaml
│   ├── content-search.yaml
│   ├── batch-operations.yaml
│   ├── import-export.yaml
│   └── system-administration.yaml
├── executor/
│   └── main.go
├── config.yaml
└── README.md
```

**Configuration** (`config.yaml`):
```yaml
canvus:
  api_url: https://your-server/api/v1
  api_key: ${CANVUS_API_KEY}
  timeout: 30s

skills:
  directory: ./skills
  auto_load: false
  cache_duration: 5m

logging:
  level: info
  format: json
```

### Docker Deployment

```dockerfile
FROM golang:1.21-alpine

WORKDIR /app

# Copy skills definitions
COPY skills/ /app/skills/

# Copy executor
COPY executor/ /app/
RUN go build -o canvus-skills-executor

# Configuration
ENV CANVUS_API_URL=""
ENV CANVUS_API_KEY=""

CMD ["./canvus-skills-executor"]
```

### Claude Desktop Integration

**Add to Claude Desktop config**:
```json
{
  "skills": {
    "canvus": {
      "command": "/path/to/canvus-skills-executor",
      "skills_directory": "/path/to/skills",
      "env": {
        "CANVUS_API_URL": "https://your-server/api/v1",
        "CANVUS_API_KEY": "your-api-key"
      }
    }
  }
}
```

---

## Skills vs MCP: Decision Matrix

| Factor | MCP | Claude Skills |
|--------|-----|---------------|
| **Context Efficiency** | Low (all tools loaded) | High (on-demand loading) |
| **Platform Support** | Cross-platform (any LLM) | Claude only |
| **Discoverability** | Flat list of 109 tools | Logical groups (8 skills) |
| **Maintenance** | Update entire toolset | Update individual skills |
| **Versioning** | Single version | Per-skill versions |
| **Development Time** | ~3 weeks | ~2-3 weeks |
| **Best For** | Small APIs, cross-platform | Large APIs, Claude-only |
| **Canvus Fit** | Workable but wasteful | Excellent fit |

**Recommendation for Canvus**: Claude Skills (context efficiency + logical grouping)

---

## Migration Path: MCP → Skills

If you build MCP first, here's how to migrate:

### Step 1: Map MCP Tools to Skills
```
MCP Tools (109) → Skills (8 groups)
- create_canvas → Canvas Management skill
- create_note → Widget Creation skill
- list_widgets → Widget Organization skill
... etc
```

### Step 2: Generate Skill Definitions
```bash
# Use same OpenAPI spec
./generate-skills --from openapi.yaml --output skills/
```

### Step 3: Reuse SDK Integration
```go
// Same executor logic, different wrapper
func (e *SkillExecutor) createCanvas() {
    // Same SDK call as MCP tool
    canvas, err := e.session.CreateCanvas(ctx, req)
    // ... same as before
}
```

### Step 4: Test Side-by-Side
- Run both MCP and Skills
- Compare context usage
- Verify feature parity
- Migrate users gradually

**Migration time**: ~1 week (mostly testing and docs)

---

## Success Metrics

### Week 1
- 2 skills implemented and tested
- Skill loader functional
- Context efficiency verified (> 50% savings vs MCP)

### Week 2
- All 8 skills implemented
- Full SDK coverage (109 operations)
- Integration tests passing

### Week 3
- Production deployment ready
- Documentation complete
- Performance benchmarks met:
  - Skill load time: < 100ms
  - Operation execution: < 200ms overhead
  - Context savings: > 60% vs MCP

---

## Conclusion

Claude Skills provide a superior approach for exposing the Canvus API compared to MCP:

✅ **70% context savings** - Load 33 operations instead of 109
✅ **Logical organization** - 8 skills vs 109 flat tools
✅ **On-demand loading** - Better user experience
✅ **Independent evolution** - Update skills separately
✅ **SDK acceleration** - 87% development time savings

**Expected timeline**: 2-3 weeks from start to production

**Key success factor**: Leverage SDK maximally, group operations logically

---

## Next Steps

1. **Review this document** - Discuss skill groupings with team
2. **Choose approach**: Skills-first or MCP→Skills migration
3. **Set up development environment** - Clone SDK, prepare tooling
4. **Week 1 kickoff** - Implement Canvas Management + Widget Creation skills

---

**Document Version**: 1.0
**Date**: 2025-11-19
**Author**: Development Team
**Status**: Ready for implementation
**Platform**: Claude (Anthropic)
