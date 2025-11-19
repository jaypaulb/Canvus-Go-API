# Specification: Canvus Go SDK Library Distribution

## Goal

Transform the Canvus-Go-API from a developer-cloned repository into a professionally distributed, well-documented Go library that developers can install via `go get`, learn through comprehensive documentation, and leverage as a foundation for their own projects including an MCP server.

## User Stories

- As a Go developer, I want to install the Canvus SDK via `go get github.com/jaypaulb/Canvus-Go-API/canvus` so that I can integrate Canvus into my project using standard Go practices
- As a system integrator, I want comprehensive documentation and examples so that I can understand best practices for authentication, error handling, and session management
- As a tool builder, I want an OpenAPI specification so that I can build automation tools and coding agents that understand SDK capabilities

## Specific Requirements

**Go Module and Import Path Standardization**
- Update go.mod module name from `canvus-go-api` to `github.com/jaypaulb/Canvus-Go-API`
- Primary import path: `github.com/jaypaulb/Canvus-Go-API/canvus`
- Ensure all internal imports are updated to use the full module path
- Validate go.mod/go.sum compatibility with go.dev proxy requirements
- Add required LICENSE file at repository root for pkg.go.dev compliance

**Semantic Versioning Implementation**
- Version scheme: `MTCS_MAJOR.MTCS_MINOR.SDK_PATCH` (e.g., 1.2.3)
- Start at `v0.1.0` as pre-release until CLI separation and documentation are complete
- Create initial GitHub release with v0.1.0 tag
- Document version scheme in `/docs/VERSIONING.md`
- Reset SDK_PATCH to 0 when MTCS API version changes
- Tag releases as `v0.1.0`, `v0.2.0`, etc. following Go module versioning conventions

**Documentation Suite Creation**
- Create `/docs/` directory with: GETTING_STARTED.md, BEST_PRACTICES.md, TROUBLESHOOTING.md, COMPATIBILITY.md, VERSIONING.md, API_REFERENCE.md, EXAMPLES.md
- GETTING_STARTED.md: Prerequisites, installation, first API call, authentication overview (target: 30 minutes to first successful call)
- BEST_PRACTICES.md: Error handling, authentication patterns, session lifecycle, concurrency, context usage, rate limiting
- TROUBLESHOOTING.md: Top 10 common issues with solutions, debugging tips, custom RoundTripper middleware patterns
- COMPATIBILITY.md: Supported MTCS API versions, known quirks/workarounds, version compatibility matrix

**Runnable Examples Directory**
- Create `/examples/` directory structure with organized subdirectories
- Required examples: getting_started/main.go, authentication/api_key.go, authentication/login_password.go, canvases/create_and_manage.go, widgets/create_and_search.go, users/provision_and_manage.go, batch/bulk_operations.go, import_export/round_trip.go, error_handling/recovery_patterns.go, context/cancellation_and_timeouts.go
- Each example must be runnable with `go run` and include comments explaining the code
- Include README.md in examples directory linking to each example with descriptions

**Starter Templates Creation**
- Create `/templates/` directory with reusable project starters
- Templates: minimal_cli.go (command-line tool structure), web_service.go (HTTP service using SDK), batch_job.go (background job for admin tasks), integration_service.go (microservice pattern)
- Each template demonstrates production patterns: graceful shutdown, configuration management, logging, error handling
- Include inline documentation explaining customization points

**OpenAPI Specification Development**
- Create `/openapi.yaml` or `/openapi/openapi.yaml` at repository root
- Document all 109+ SDK methods in OpenAPI 3.0.0 format
- Include schemas for all Go types (Canvas, Widget, User, etc.) matching SDK types
- Document authentication schemes (API key via Private-Token header, token-based)
- Include operation tags for categorization (Canvases, Widgets, Users, System, etc.)
- Add request/response examples for each operation
- Purpose: Foundation for MCP server, code generation, and coding agents

**Release Management Process**
- Create CHANGELOG.md following Keep a Changelog format
- Establish release checklist: run tests, update docs, bump version, update CHANGELOG, create GitHub release
- GitHub releases include: version tag, release notes with SDK changes, MTCS API compatibility notes
- Documentation versioned alongside code (update docs as part of release process)
- Support policy: latest version only (no parallel version maintenance until v1.0.0)

**CLI Separation to Independent Repository**
- Create new repository: `github.com/jaypaulb/canvus-cli`
- Migrate code from `/cmd/canvus-cli/` to new repository
- CLI depends on SDK via go.mod: `require github.com/jaypaulb/Canvus-Go-API v0.1.0`
- CLI distribution: pre-built binaries for macOS (Intel/ARM), Linux (x86_64/ARM), Windows
- Installation methods: binary download from GitHub Releases, `go install github.com/jaypaulb/canvus-cli/cmd/canvus@latest`
- Future: package manager integration (Homebrew formula, Scoop manifest)
- Remove `/cmd/canvus-cli/` from SDK repository after migration

**README Enhancement**
- Restructure README.md as main entry point with feature highlights and use cases
- Include: problem statement, quick feature overview, installation instructions, simple example
- Add clear links to comprehensive docs in /docs/
- List use cases: integrations, automation, admin tools, MCP server foundation
- Include badges: Go version, license, godoc link, latest release

**Testing and Quality Validation**
- Validate all examples compile and run successfully
- Ensure OpenAPI spec matches actual SDK implementation
- Integration test requirements documented
- Example validation as part of CI/release process
- Documentation consistency checking against code

## Existing Code to Leverage

**Session and Authentication Pattern (canvus/session.go)**
- Existing NewSession() with functional options pattern (WithAPIKey, WithToken)
- APIKeyAuthenticator and TokenAuthenticator interfaces already implemented
- Custom RoundTripper pattern (transportWithAPIKey) can be documented as extension point
- Login() method for username/password authentication flow
- Use as basis for authentication documentation and examples

**Typed Error Handling (canvus/errors.go)**
- APIError struct with StatusCode, Code, Message, RequestID, Details, Wrapped
- ErrorCode constants for all HTTP status codes and SDK-specific errors
- Unwrap() and Is() methods for error chain handling
- Use as basis for error handling best practices documentation
- Document patterns for type assertions: `if apiErr, ok := err.(*canvus.APIError); ok {...}`

**Batch Processing Framework (canvus/batch.go)**
- BatchProcessor with automatic retry logic already implemented
- Demonstrates concurrent operation patterns
- Use as example for batch/bulk_operations.go runnable example
- Document as advanced usage pattern in BEST_PRACTICES.md

**Import/Export Functionality (canvus/import.go, canvus/export.go)**
- Round-trip safe import/export for all widget and asset types
- Asset file handling (images, PDFs, videos)
- Use as basis for import_export/round_trip.go example
- Document fidelity guarantees in COMPATIBILITY.md

**Geometry Utilities (canvus/geometry.go)**
- WidgetsContainId, WidgetsTouchId functions
- Spatial containment and overlap detection
- Document as advanced feature in examples and API reference

## Out of Scope

- Creating Python, JavaScript, or other language SDKs (Go only)
- Modifying the underlying Canvus/MTCS API server
- Managed hosting or SaaS offering deployment
- Building a web UI/GUI for administration
- Framework-specific plugins (Rails, FastAPI integration not in scope)
- Database ORM or persistence layer beyond SDK patterns
- Supporting multiple parallel SDK versions (latest only)
- Automated deployment pipelines for the SDK itself
- User authentication/authorization beyond API key and token flows
- Real-time event streaming or WebSocket support (unless already in SDK)
- Internationalization or localization of documentation
- Video tutorials or multimedia documentation
- Community forum or discussion platform setup
- Paid support tiers or enterprise features

