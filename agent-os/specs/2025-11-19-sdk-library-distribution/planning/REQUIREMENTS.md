# Canvus Go SDK Library Distribution and Developer Experience - Requirements

## Executive Summary

The Canvus-Go-API SDK must be distributed and packaged as a production-ready, reusable library that enables Go developers to integrate Canvus into their projects through standard Go practices. This requires comprehensive documentation, a robust versioning strategy, separation of CLI tooling, and an agent-consumable specification for future ecosystem expansion. The SDK serves as the foundation for all Canvus-Go integrations, including a future MCP server and custom developer projects.

## Problem Statement and Vision

### Current State
- SDK exists as a Go module with full API coverage (109 methods)
- Advanced features implemented (geometry utilities, batch processing, import/export, filtering)
- Comprehensive integration tests and godoc documentation in place
- No formal distribution/release process
- CLI tools currently embedded in the SDK repository
- Developers unfamiliar with best practices for library usage
- No agent-consumable specifications for future automation tooling

### Vision
Position the Canvus-Go-API as the authoritative, reusable Go library for Canvus integrations. Developers should:
- Install via standard `go get github.com/jaypaulb/Canvus-Go-API/canvus`
- Follow idiomatic patterns through comprehensive documentation and examples
- Build custom integrations with confidence using well-documented APIs
- Understand compatibility expectations as the underlying MTCS API evolves
- Integrate with future automation tools (MCP servers, coding agents) through formal specifications

The SDK is not an end product; it is a foundation for an ecosystem of projects, integrations, and tooling.

## Deliverables

### 1. SDK Library Distribution
- **Primary Import Path**: `github.com/jaypaulb/Canvus-Go-API/canvus` (maintain current standard path)
- **Release Artifacts**: Versioned releases with semantic versioning
- **Go Module Support**: Full go.mod compatibility for dependency management
- **No Breaking Changes Within Patch Versions**: Follow semantic versioning strictly

### 2. Documentation Suite
- **Getting Started Guide**: Quick-start walkthrough for new developers
  - Installation instructions (`go get`)
  - First API call example
  - Authentication methods overview
  - Expected output and error handling

- **Common Workflow Examples**: 5-7 practical examples covering:
  - Canvas creation and management workflows
  - User provisioning and authentication
  - Widget operations and asset management
  - Batch operations and bulk processing
  - Error handling and recovery patterns
  - Context usage for cancellation and timeouts

- **Best Practices Documentation**:
  - Error handling strategies (typed errors, retry logic, fallbacks)
  - Authentication patterns (API keys, login/password, token refresh)
  - Session management (creation, refresh, validation)
  - Concurrency and goroutine safety
  - Context and timeout management
  - Rate limiting and pagination patterns
  - Memory efficiency for large operations

- **Troubleshooting and FAQ**:
  - Common authentication issues and solutions
  - Network and timeout problems
  - Compatibility considerations with different MTCS API versions
  - Performance tuning tips
  - When to use batch operations vs. sequential calls
  - Debugging with custom RoundTripper middleware

- **API Reference Documentation**:
  - Comprehensive godoc (already in place, verify completeness)
  - Method-by-method documentation with usage examples
  - Type documentation for all request/response models
  - Error types and their meanings

### 3. OpenAPI/Swagger Specification
- **Purpose**: Machine-readable specification of SDK capabilities for future automation tools
- **Content**:
  - All 109 methods documented in OpenAPI 3.0 format
  - Request/response schemas matching Go types
  - Authentication schemes (API key, token-based)
  - Error responses and status codes
  - Resource models and relationships
  - Operation tags and categorization

- **Deliverable Location**: `openapi.yaml` or `openapi/` directory at repository root
- **Usage**: Foundation for MCP server, code generation, and coding agents to understand SDK capabilities
- **Format**: Standard OpenAPI 3.0.0 specification, AI-friendly structure

### 4. Examples and Templates
- **Runnable Examples** (in `examples/` directory):
  - `getting_started/main.go` - Simple initialization and first call
  - `authentication/api_key.go` - API key authentication example
  - `authentication/login_password.go` - Login/password flow
  - `canvases/create_and_manage.go` - Canvas lifecycle operations
  - `widgets/create_and_search.go` - Widget operations and cross-canvas search
  - `users/provision_and_manage.go` - User creation and group management
  - `batch/bulk_operations.go` - Batch processing framework usage
  - `import_export/round_trip.go` - Import/export with fidelity
  - `error_handling/recovery_patterns.go` - Typed error handling and retries
  - `context/cancellation_and_timeouts.go` - Context usage patterns

- **Starter Templates**:
  - Minimal CLI application template using the SDK
  - Web service template demonstrating session management
  - Batch job template for administrative tasks
  - Integration service template for system automation

- **Framework Patterns**:
  - Custom RoundTripper middleware example (logging, metrics, auth)
  - Batch operation patterns with retry logic
  - Custom filtering implementations
  - Streaming and pagination patterns

### 5. Release Management
- **Version Numbering Scheme**: `MTCS_MAJOR.MTCS_MINOR.SDK_PATCH`
  - Format: `1.2.3` where `1.2` = MTCS API version, `3` = SDK patch version
  - Example: `0.1.0` (starting version until all features and CLI complete)
  - When MTCS API updates: increment `MTCS_MAJOR.MTCS_MINOR`, reset `SDK_PATCH` to `0`
  - When SDK-only changes: increment `SDK_PATCH`

- **Versioning Strategy**:
  - Start at `v0.1.0` until CLI is complete and all features are finalized
  - Support only the latest version (no parallel version maintenance)
  - Version tags in GitHub for release tracking
  - Release notes documenting SDK changes and MTCS API compatibility

- **Release Artifacts**:
  - Tagged releases in GitHub matching semantic versions
  - Release notes with changelog, SDK improvements, and MTCS API compatibility notes
  - Documentation versioned alongside code

- **Compatibility Documentation**:
  - Current MTCS API version supported
  - Known quirks or oddities inherited from MTCS API
  - Workarounds for known issues
  - Plan for updating as MTCS API evolves
  - Breaking change notice process for major versions

### 6. CLI Separation and Independence
- **New Repository**: `github.com/jaypaulb/canvus-cli`
  - Separate from SDK repository
  - Independent release cycle
  - Can be updated without affecting SDK versions

- **Distribution**:
  - Pre-built binaries for macOS (Intel/ARM), Linux (x86_64/ARM), Windows
  - Available for download from GitHub Releases
  - Installation: binary download or `go install` for developers

- **Relationship to SDK**:
  - CLI depends on canvus SDK via normal go.mod dependency
  - CLI demonstrates real-world SDK usage through example commands
  - Can update CLI without updating SDK and vice versa
  - Current cmd/canvus-cli/ code will be migrated to new repository

- **Future Expansion**:
  - CLI is a reference implementation of SDK usage
  - Developers can use CLI as pattern for their own integrations
  - CLI serves as proof that SDK supports complete Canvus workflows

### 7. Library-First Design Principles
- **Primary Use Case**: Importable library for Go projects
- **Design Decisions**:
  - All functionality exported and documented for library use
  - No assumptions about runtime environment (CLI, web service, batch job)
  - Context support throughout enables any execution model
  - Flexible authentication to support any credential management approach
  - No global state; all operations through client instances

- **Foundation for Future Ecosystem**:
  - MCP Server Project: Wraps SDK to expose capabilities to AI models
  - Coding Agents: Use OpenAPI spec to understand and invoke SDK methods
  - Custom Integrations: Developers build on SDK foundation
  - Community Extensions: Plugins and utilities built using public SDK APIs

- **Guarantees**:
  - Public APIs remain stable within major versions
  - Clear deprecation path for changes
  - No undocumented dependencies or hidden behaviors
  - Full type safety and explicit error handling

## Architecture and Distribution Approach

### Import Path Strategy
- **Import Path**: `github.com/jaypaulb/Canvus-Go-API/canvus`
- **Rationale**: Follows Go standards; differentiates module name from repository name
- **Go Module Declaration**: Module name: `github.com/jaypaulb/Canvus-Go-API`
- **Package Name**: `canvus` at repository root
- **No Changes**: Current structure is correct; maintain as-is

### Distribution Model
- **Primary Channel**: GitHub Releases with semantic version tags
- **Go Module Support**: Full go.mod/go.sum support for standard `go get`
- **Discoverability Multi-Pronged Approach**:
  1. **Clarity**: Comprehensive README with feature highlights, use cases, and examples
  2. **Documentation**: Getting started guide, API reference, best practices docs
  3. **Tooling**: Runnable examples, starter templates, framework patterns
  4. **Ecosystem**: OpenAPI spec for automation tools; CLI as reference implementation

### Release Cycle
- **Stable Releases**: Semantic version tags when features are complete
- **Pre-Releases**: v0.x.x used until CLI separation and all documentation complete
- **Patch Releases**: SDK improvements and bug fixes
- **Feature Releases**: When MTCS API updates or significant SDK enhancements

### Integration Points
- **Dependency Specification**: Users specify version in go.mod:
  ```
  require github.com/jaypaulb/Canvus-Go-API v0.1.0
  ```
- **Upgrade Path**: Users manage SDK updates via go get -u or direct version updates
- **Backwards Compatibility**: Maintain compatibility within major versions

## Version Strategy Details

### Semantic Versioning Scheme

**Format**: `MTCS_MAJOR.MTCS_MINOR.SDK_PATCH` (e.g., 1.2.3)

**Component Meanings**:
- **MTCS_MAJOR.MTCS_MINOR**: Reflects the version of the Canvus/MTCS API that the SDK wraps
  - Incremented when the underlying MTCS API has updates
  - Indicates which Canvus API version this SDK release is compatible with
  - Example: If MTCS API is version 1.2, SDK versions start at 1.2.0

- **SDK_PATCH**: Reflects SDK-specific improvements, bug fixes, and enhancements
  - Incremented for SDK-only changes that don't reflect API updates
  - Reset to 0 whenever MTCS_MAJOR or MTCS_MINOR changes
  - Examples:
    - 1.2.0 → 1.2.1: SDK bug fix or improvement
    - 1.2.3 → 1.3.0: MTCS API updated to 1.3, SDK reset patch to 0

**Current State**:
- Start at `v0.1.0` (pre-release)
- Increment patch for improvements until CLI is complete
- Move to stable `v1.0.0` or higher once ready for production

**Future Evolution**:
- Monitor MTCS API for updates
- When MTCS API version changes (e.g., 1.2 → 1.3):
  1. Update SDK code for new/changed endpoints
  2. Bump version to 1.3.0
  3. Document changes in release notes
  4. Update compatibility documentation

**No Parallel Versions**:
- Support only latest version
- No backporting to older versions
- Users must upgrade SDK to get patches
- Encourages staying current with API

### Version Compatibility and Support

**Compatibility Documentation** (in `docs/COMPATIBILITY.md` or similar):
- Current MTCS API version supported (e.g., 1.2.x)
- SDK behavior with older MTCS API versions
- Known oddities or quirks inherited from MTCS API
- Workarounds for known MTCS API issues

**Breaking Changes**:
- Reserve major version increments for breaking changes
- Clear deprecation period before removing features
- Release notes explicitly document breaking changes
- Migration guide provided for updates

**Stability Promise**:
- Accept current MTCS API oddities and document them
- Maintain SDK compatibility with stated MTCS versions
- Plan for updates as MTCS API evolves
- Communicate changes clearly to users

## Documentation Artifacts

### File Structure (in repository root)
```
/
├── README.md (main entry point, feature highlights)
├── docs/
│   ├── GETTING_STARTED.md (quick start guide)
│   ├── EXAMPLES.md (links to runnable examples)
│   ├── BEST_PRACTICES.md (auth, errors, sessions, concurrency)
│   ├── TROUBLESHOOTING.md (FAQ and common issues)
│   ├── COMPATIBILITY.md (MTCS API version support)
│   ├── VERSIONING.md (version scheme explanation)
│   └── API_REFERENCE.md (links to godoc or inline reference)
├── examples/
│   ├── getting_started/main.go
│   ├── authentication/
│   ├── canvases/
│   ├── widgets/
│   ├── users/
│   ├── batch/
│   ├── import_export/
│   ├── error_handling/
│   └── context/
├── templates/
│   ├── minimal_cli.go
│   ├── web_service.go
│   ├── batch_job.go
│   └── integration_service.go
├── openapi.yaml (or openapi/openapi.yaml)
├── canvus/ (main library package)
└── cmd/canvus-cli/ (will be removed after migration)
```

### Documentation Content Requirements

**README.md** (Enhanced):
- Problem statement (why Canvus-Go-API exists)
- Quick feature overview
- Installation instructions
- Simple example
- Links to comprehensive docs
- Use cases (integrations, automation, admin tools)
- Key differentiators

**GETTING_STARTED.md**:
- Prerequisites (Go 1.16+, Canvus server access)
- Installation via go get
- Creating a client instance
- Authentication options overview
- Making your first API call
- Handling errors
- Common next steps

**BEST_PRACTICES.md**:
- Error handling (typed errors, retries, fallbacks)
- Authentication security
- Session lifecycle management
- Concurrent API calls and goroutine safety
- Context usage and timeouts
- Pagination and streaming patterns
- Resource cleanup (session logout)
- Logging and debugging

**TROUBLESHOOTING.md**:
- Authentication failures (causes and solutions)
- Network timeouts and retries
- Connection refused errors
- Rate limiting responses
- Compatibility issues across MTCS versions
- Performance tips
- When to contact support

**COMPATIBILITY.md**:
- Supported MTCS API versions
- Version compatibility matrix
- Known issues and workarounds
- Deprecated features
- Roadmap for future updates

**VERSIONING.md**:
- Semantic versioning scheme explanation
- What each version component means
- When versions change
- Compatibility promises
- How to stay current

**API_REFERENCE.md**:
- Link to godoc documentation
- Overview of API categories (users, canvases, widgets, etc.)
- Method grouping and organization
- Common patterns across similar methods

### Examples and Templates

**Runnable Examples** (all in `examples/` with README):
1. Getting Started: Basic client setup and first API call
2. Authentication (3 examples): API key, login/password, token refresh
3. Canvas Management: Create, list, update, delete, copy canvases
4. Widget Operations: Create, search, update, delete widgets across canvases
5. User Management: Create users, manage tokens, group membership
6. Batch Processing: Bulk operations with automatic retry
7. Import/Export: Round-trip import/export maintaining fidelity
8. Error Handling: Typed error handling, retry patterns, circuit breakers
9. Context Usage: Cancellation, timeouts, request deadlines

**Starter Templates** (in `templates/`):
1. Minimal CLI: Command-line tool structure
2. Web Service: HTTP service using SDK for business logic
3. Batch Job: Background job for administrative tasks
4. Integration Service: Microservice integrating Canvus with other systems

### OpenAPI/Swagger Specification

**File Location**: `/openapi.yaml` or `/openapi/openapi.yaml`

**Content**:
```yaml
openapi: 3.0.0
info:
  title: Canvus Go SDK
  version: 0.1.0
  description: OpenAPI specification for Canvus-Go-API SDK capabilities

servers:
  - url: https://canvus-server
    description: Canvus API Server

paths:
  /canvases:
    get:
      operationId: ListCanvases
      summary: List all canvases
      parameters: [...]
      responses: {...}
      tags:
        - Canvases

  # ... all 109+ methods documented

components:
  schemas:
    Canvas: {...}
    Widget: {...}
    User: {...}
    # ... all types from SDK

  securitySchemes:
    ApiKeyAuth:
      type: apiKey
      in: header
      name: X-API-Key
    TokenAuth:
      type: http
      scheme: bearer
```

**Specification Purpose**:
- Machine-readable documentation of SDK capabilities
- Foundation for code generation tools
- Input for MCP server generator
- Reference for coding agents (Claude, other AI models)
- Developer reference for API structure

**Maintenance**:
- Update whenever SDK APIs change
- Validate against actual SDK implementation
- Include in release process
- Use for documentation consistency checking

## SDK and CLI Separation Strategy

### Current State
- CLI commands in `cmd/canvus-cli/` directory
- CLI tied to SDK release cycle
- Users must clone entire repository for SDK

### Separation Plan

**Phase 1: Create New Repository**
- Repository: `github.com/jaypaulb/canvus-cli`
- Migrate code from `cmd/canvus-cli/`
- Set up release infrastructure
- Create distribution binaries

**Phase 2: Update SDK Repository**
- Remove CLI code from `cmd/canvus-cli/`
- Keep SDK as clean library
- Remove CLI build artifacts from releases
- Focus SDK on library usage

**Phase 3: Independence**
- CLI is separate project with separate versioning
- CLI depends on SDK via go.mod (like any other project)
- Can release CLI without SDK changes
- Can update SDK without CLI changes
- CLI serves as reference implementation

### Distribution

**SDK Distribution**:
- `go get github.com/jaypaulb/Canvus-Go-API/canvus`
- GitHub Releases with source code
- GoDoc documentation hosted automatically

**CLI Distribution**:
- Standalone repository: `github.com/jaypaulb/canvus-cli`
- Pre-built binaries for each platform:
  - macOS (Intel/ARM)
  - Linux (x86_64/ARM)
  - Windows
- Installation methods:
  - Download binary from GitHub Releases
  - `go install github.com/jaypaulb/canvus-cli/cmd/canvus@latest`
  - Package managers (Homebrew, etc., future)

### Benefits

**For SDK Users**:
- Cleaner repository focused on library
- Faster clones and downloads
- Clear library vs. tooling distinction
- No CLI dependencies required

**For CLI Users**:
- Standalone tool independent of SDK releases
- Faster iterations without SDK process
- Natural home for command-line specific logic
- Can use SDK as direct dependency with version pinning

**For Maintainers**:
- Cleaner responsibility separation
- Independent release cycles reduce coordination
- CLI can be stable while SDK evolves, or vice versa
- Easier for community contributions to either project

## Key Design Principles

### 1. Library-First Design
- Everything is exposed and documented for library use
- No hidden runtime assumptions
- Flexible enough for any project type (CLI, web, batch)
- Clean public APIs with no global state

### 2. Idiomatic Go
- Follow Go community conventions
- Context support throughout
- Proper error handling with typed errors
- Standard library patterns (io.Reader, json.Marshaler, etc.)
- No framework dependencies for core functionality

### 3. Stability and Compatibility
- Semantic versioning with clear meaning
- Maintain compatibility within major versions
- Document breaking changes clearly
- No surprise API changes
- Long-term support for documented APIs

### 4. Developer Experience
- Comprehensive documentation with examples
- Clear error messages and troubleshooting
- Easy onboarding for new developers
- Best practices documented and exemplified
- Active support and communication

### 5. Foundation for Ecosystem
- Designed for reuse in other projects
- Clear extension points for customization
- Foundation for MCP server and agents
- Enables community contributions and extensions
- Documentation sufficient for framework builders

## Agent-Consumable Documentation

### Purpose
Enable future automation tools (MCP servers, coding agents, code generators) to understand and use the SDK.

### Specification Format
- **OpenAPI 3.0**: Machine-readable API specification
- **Schema Definition**: All types and their relationships
- **Operation Categorization**: Tags for grouping related operations
- **Examples in Spec**: Request/response examples for each operation
- **Error Codes**: Documented error responses for each operation

### Content Specification
1. All 109+ SDK methods defined
2. Request parameters with types and requirements
3. Response schemas with examples
4. Error responses with status codes
5. Authentication methods documented
6. Resource relationships and dependencies
7. Pagination and filtering patterns
8. Batch operation patterns
9. Rate limiting (if applicable)
10. Required fields vs. optional parameters

### Integration Examples
- MCP Server uses spec to expose SDK via Claude protocol
- Code generation tools create client libraries in other languages
- Coding agents understand available operations and call them correctly
- IDE extensions provide autocomplete for SDK methods

## Non-Goals / Out of Scope

The following are explicitly not part of this feature:

1. **Language Bindings**: Python, JavaScript, or other language SDKs
2. **API Server Implementation**: We don't modify the Canvus/MTCS API itself
3. **Managed Hosting**: No SaaS offering or cloud deployment
4. **GUI Administration Tool**: No web UI for admin operations
5. **Framework Integration**: Specific plugins for Rails, FastAPI, etc. (can be community-driven)
6. **Database ORM**: Persistence layer beyond SDK usage patterns
7. **Multiple Version Support**: Support only latest version

## Success Criteria

### Documentation Completeness
- Getting Started guide enables new user to make first API call in < 30 minutes
- All 109 methods have documented godoc with examples
- At least 5 runnable examples demonstrate common workflows
- Best practices guide covers error handling, auth, sessions, concurrency
- Troubleshooting FAQ addresses top 10 common issues
- OpenAPI spec matches SDK implementation

### Distribution Readiness
- SDK installable via `go get github.com/jaypaulb/Canvus-Go-API/canvus`
- Version tags follow semantic versioning scheme
- Release notes document SDK changes and MTCS compatibility
- CLI separated to independent repository
- Pre-built CLI binaries available for major platforms

### Developer Experience
- Developers can build custom integrations without reference to CLI code
- Examples demonstrate best practices for error handling, auth, sessions
- Import path is clear and follows Go standards
- Documentation is discoverable from main README
- Troubleshooting guide reduces support burden

### Foundation for Ecosystem
- OpenAPI spec enables MCP server development
- Library-first design enables custom integration projects
- Public APIs are stable and well-documented
- Extension points clear for community contributions

## Questions Addressed Through Requirements Research

**Q1: What's the biggest distribution challenge?**
- Answer: Discoverability (clarity, tooling, examples) - Addressed through multi-pronged approach (docs, examples, templates, OpenAPI)

**Q2: What versioning strategy?**
- Answer: v0.1.0 start, MTCS_MAJOR.MTCS_MINOR.SDK_PATCH scheme, latest version only - Documented in versioning strategy section

**Q3: Documentation needs?**
- Answer: Getting started, workflows, best practices, troubleshooting, OpenAPI spec - Complete documentation suite defined

**Q4: CLI strategy?**
- Answer: Separate repository, independent cycle - CLI separation strategy documented

**Q5: Import path?**
- Answer: Keep current path - github.com/jaypaulb/Canvus-Go-API/canvus confirmed

**Q6: Examples and templates?**
- Answer: Yes to all (runnable, templates, patterns) - Complete examples and templates defined

**Q7: Breaking changes?**
- Answer: Accept current oddities, document, plan for updates - Compatibility strategy documented

**Q8: Overarching goal?**
- Answer: Reusable foundation for future projects - Library-first design and ecosystem principles documented

## Implementation Priorities

### Phase 1: Foundation (Critical)
1. Documentation suite (Getting Started, Best Practices, Troubleshooting)
2. Runnable examples (getting_started, authentication, canvas_ops, widgets)
3. Version strategy implementation and tagging

### Phase 2: Completeness (High Priority)
1. OpenAPI specification
2. Additional runnable examples (batch, import/export, error handling)
3. CLI separation to new repository
4. Release infrastructure setup

### Phase 3: Polish (Important)
1. Starter templates
2. Framework patterns
3. Advanced examples (custom middleware, filtering)
4. Community documentation setup

### Phase 4: Ecosystem (Future)
1. MCP server project using SDK and OpenAPI spec
2. Documentation for extending and building on SDK
3. Community contribution guidelines
4. Integration showcase

## Stakeholders and Communication

**Internal**:
- Repository Maintainer: Overall project direction
- Documentation Lead: Guides comprehensive documentation creation
- CLI Maintainers: Plan separation and migration

**External**:
- Go Developers: Primary users installing and using SDK
- System Administrators: Using CLI for automation
- DevOps Engineers: Building integrations
- Future Tool Builders: Using OpenAPI spec for automation

**Communication**:
- Release notes for each version with changelog
- Getting Started guide for new users
- FAQ and troubleshooting for common issues
- OpenAPI spec for tool builders

---

**Document Version**: 1.0
**Last Updated**: 2025-11-19
**Status**: Requirements Complete - Ready for Specification Phase
