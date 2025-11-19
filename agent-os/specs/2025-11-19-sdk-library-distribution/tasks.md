# Task Breakdown: Canvus Go SDK Library Distribution

## Overview

**Total Tasks**: 78 tasks across 6 phases
**Estimated Total Effort**: 60-80 hours
**Target Completion**: v0.1.0 release

This task list transforms the Canvus-Go-API from a developer-cloned repository into a professionally distributed, well-documented Go library with comprehensive documentation, runnable examples, OpenAPI specification, and CLI separation.

---

## Task List

### Phase 1: Foundation

Establish versioning, module setup, and release infrastructure to enable professional distribution.

#### Task Group 1.1: Go Module Standardization
**Dependencies:** None
**Effort Total:** M (2-4 hours)

- [x] 1.1.1 Update go.mod module name from `canvus-go-api` to `github.com/jaypaulb/Canvus-Go-API` [S]
  - File: `/home/jaypaulb/Documents/gh/Canvus-Go-API/go.mod`
  - Verify module declaration matches GitHub path
  - Run `go mod tidy` after change

- [x] 1.1.2 Update all internal imports to use full module path [S]
  - Search for imports using old path
  - Update imports in all `.go` files
  - Files to check: `canvus/*.go`, `cmd/canvus-cli/*.go`

- [x] 1.1.3 Verify LICENSE file exists at repository root for pkg.go.dev compliance [XS]
  - File: `/home/jaypaulb/Documents/gh/Canvus-Go-API/LICENSE`
  - Confirm license type (recommend MIT or Apache 2.0)
  - Ensure proper formatting for automated detection

- [x] 1.1.4 Validate go.mod/go.sum compatibility with go.dev proxy [S]
  - Run `go mod verify`
  - Run `go build ./...`
  - Run `go test ./...` to ensure all tests pass
  - Verify no private dependencies that would block proxy

**Acceptance Criteria:**
- `go get github.com/jaypaulb/Canvus-Go-API/canvus` works correctly
- All internal imports use full module path
- LICENSE file present and detectable
- All tests pass with updated imports

---

#### Task Group 1.2: Semantic Versioning Setup
**Dependencies:** Task Group 1.1
**Effort Total:** M (2-4 hours)

- [x] 1.2.1 Create `/home/jaypaulb/Documents/gh/Canvus-Go-API/docs/VERSIONING.md` [M]
  - Document version scheme: `MTCS_MAJOR.MTCS_MINOR.SDK_PATCH`
  - Explain what each component means
  - Include examples of version bumps
  - Document when to reset SDK_PATCH to 0
  - Include compatibility promises

- [x] 1.2.2 Create `/home/jaypaulb/Documents/gh/Canvus-Go-API/CHANGELOG.md` [S]
  - Follow Keep a Changelog format (https://keepachangelog.com)
  - Create initial entry for v0.1.0
  - Include sections: Added, Changed, Deprecated, Removed, Fixed, Security
  - Document all current SDK features

- [x] 1.2.3 Create release checklist document [S]
  - Location: `/home/jaypaulb/Documents/gh/Canvus-Go-API/docs/RELEASING.md`
  - Steps: run tests, update docs, bump version, update CHANGELOG, create tag, create GitHub release
  - Include pre-release validation steps
  - Document MTCS API compatibility notes format

**Acceptance Criteria:**
- VERSIONING.md clearly explains the version scheme
- CHANGELOG.md follows Keep a Changelog format
- Release process is documented and repeatable

---

#### Task Group 1.3: Release Infrastructure
**Dependencies:** Task Group 1.2
**Effort Total:** M (2-4 hours)

- [x] 1.3.1 Create GitHub Actions workflow for release validation [M]
  - File: `/home/jaypaulb/Documents/gh/Canvus-Go-API/.github/workflows/release.yml`
  - Run tests on tag push
  - Validate examples compile
  - Build and verify module

- [x] 1.3.2 Create GitHub release template [S]
  - File: `/home/jaypaulb/Documents/gh/Canvus-Go-API/.github/RELEASE_TEMPLATE.md`
  - Include: version, SDK changes, MTCS API compatibility notes
  - Link to CHANGELOG for details
  - Include installation instructions

- [x] 1.3.3 Create initial v0.1.0 tag (do not push yet) [XS]
  - Document tag creation process
  - Will be pushed after all Phase 1-3 tasks complete
  - **NOTE**: Tag creation should be done as part of Task Group 6.3 (v0.1.0 Release) after all phases are complete
  - Tag command: `git tag -a v0.1.0 -m "Initial release"`
  - Push command: `git push origin v0.1.0`
  - See `/home/jaypaulb/Documents/gh/Canvus-Go-API/docs/RELEASING.md` for full release process

**Acceptance Criteria:**
- GitHub Actions workflow runs on tag creation
- Release template provides consistent format
- Tag creation process is documented

---

### Phase 2: Documentation

Create comprehensive documentation suite to enable developer success.

#### Task Group 2.1: Core Documentation
**Dependencies:** Task Group 1.1
**Effort Total:** L (4-8 hours)

- [x] 2.1.1 Create `/home/jaypaulb/Documents/gh/Canvus-Go-API/docs/GETTING_STARTED.md` [M]
  - Prerequisites (Go 1.16+, Canvus server access)
  - Installation via `go get github.com/jaypaulb/Canvus-Go-API/canvus`
  - Creating first client instance with code example
  - Authentication options overview (API key, login/password)
  - Making first API call with complete code
  - Handling errors with example
  - Target: developer makes first successful call in < 30 minutes

- [x] 2.1.2 Create `/home/jaypaulb/Documents/gh/Canvus-Go-API/docs/BEST_PRACTICES.md` [L]
  - Error handling patterns with typed errors (`*canvus.APIError`)
  - Authentication patterns and security considerations
  - Session lifecycle management (creation, refresh, validation, cleanup)
  - Concurrency and goroutine safety guidelines
  - Context usage for cancellation and timeouts
  - Rate limiting and pagination patterns
  - Memory efficiency for large operations
  - Logging and debugging tips
  - Resource cleanup (session logout)

- [x] 2.1.3 Create `/home/jaypaulb/Documents/gh/Canvus-Go-API/docs/TROUBLESHOOTING.md` [M]
  - Top 10 common issues with solutions:
    1. Authentication failures (invalid API key, expired token)
    2. Connection refused errors
    3. Network timeouts
    4. Certificate verification failures
    5. Rate limiting responses
    6. Permission denied errors
    7. Invalid request errors
    8. Resource not found errors
    9. Batch operation failures
    10. Import/export issues
  - Debugging with custom RoundTripper middleware
  - Performance tuning tips
  - When to contact support

**Acceptance Criteria:**
- Getting Started guide enables first API call in < 30 minutes
- Best Practices covers all major SDK patterns
- Troubleshooting addresses top 10 issues with solutions

---

#### Task Group 2.2: Reference Documentation
**Dependencies:** Task Group 1.1
**Effort Total:** M (2-4 hours)

- [x] 2.2.1 Create `/home/jaypaulb/Documents/gh/Canvus-Go-API/docs/COMPATIBILITY.md` [M]
  - Current MTCS API version supported
  - Version compatibility matrix
  - Known quirks and workarounds for MTCS API oddities
  - Deprecated features (if any)
  - Breaking change notice process
  - Plan for updates as MTCS API evolves

- [x] 2.2.2 Create `/home/jaypaulb/Documents/gh/Canvus-Go-API/docs/API_REFERENCE.md` [S]
  - Link to pkg.go.dev documentation
  - Overview of API categories (Users, Canvases, Widgets, System, etc.)
  - Method grouping and organization
  - Common patterns across similar methods
  - Links to relevant examples

- [x] 2.2.3 Create `/home/jaypaulb/Documents/gh/Canvus-Go-API/docs/EXAMPLES.md` [S]
  - Index of all runnable examples with descriptions
  - Categorized by use case
  - Difficulty indicators (beginner, intermediate, advanced)
  - Links to each example with brief explanation

**Acceptance Criteria:**
- Compatibility matrix covers all supported MTCS versions
- API Reference provides clear navigation to godoc
- Examples index helps users find relevant code quickly

---

### Phase 3: Examples & Templates

Create runnable examples and starter templates for common use cases.

#### Task Group 3.1: Core Examples
**Dependencies:** Task Group 1.1
**Effort Total:** L (4-8 hours)

- [x] 3.1.1 Create examples directory structure [XS]
  - `/home/jaypaulb/Documents/gh/Canvus-Go-API/examples/README.md`
  - Create subdirectories: getting_started, authentication, canvases, widgets, users

- [x] 3.1.2 Create `/home/jaypaulb/Documents/gh/Canvus-Go-API/examples/getting_started/main.go` [S]
  - Simple initialization and first call
  - List canvases with API key authentication
  - Error handling
  - Comments explaining each step
  - Must be runnable with `go run main.go`

- [x] 3.1.3 Create `/home/jaypaulb/Documents/gh/Canvus-Go-API/examples/authentication/api_key.go` [S]
  - API key authentication example
  - Using WithAPIKey functional option
  - Making authenticated requests
  - Proper error handling

- [x] 3.1.4 Create `/home/jaypaulb/Documents/gh/Canvus-Go-API/examples/authentication/login_password.go` [S]
  - Login/password authentication flow
  - Using Login() method
  - Token handling
  - Session cleanup with Logout()

- [x] 3.1.5 Create `/home/jaypaulb/Documents/gh/Canvus-Go-API/examples/canvases/create_and_manage.go` [M]
  - Canvas lifecycle: create, list, get, update, delete
  - Using ListOptions for filtering
  - Handling pagination
  - Copy canvas operation

- [x] 3.1.6 Create `/home/jaypaulb/Documents/gh/Canvus-Go-API/examples/widgets/create_and_search.go` [M]
  - Widget operations: create, list, update, delete
  - Cross-canvas widget search
  - Using geometry utilities (WidgetsContainId)
  - Different widget types (note, browser, image)

**Acceptance Criteria:**
- All examples compile with `go build`
- All examples are runnable with `go run`
- Each example has comprehensive comments
- Examples demonstrate real-world patterns

---

#### Task Group 3.2: Advanced Examples
**Dependencies:** Task Group 3.1
**Effort Total:** L (4-8 hours)

- [x] 3.2.1 Create additional example directories [XS]
  - Create subdirectories: batch, import_export, error_handling, context

- [x] 3.2.2 Create `/home/jaypaulb/Documents/gh/Canvus-Go-API/examples/users/provision_and_manage.go` [M]
  - User creation and management
  - Token generation
  - Group membership operations
  - Permission handling

- [x] 3.2.3 Create `/home/jaypaulb/Documents/gh/Canvus-Go-API/examples/batch/bulk_operations.go` [M]
  - Using BatchProcessor for bulk operations
  - Automatic retry logic demonstration
  - Concurrent operation patterns
  - Progress tracking

- [x] 3.2.4 Create `/home/jaypaulb/Documents/gh/Canvus-Go-API/examples/import_export/round_trip.go` [M]
  - Export canvas with all widgets and assets
  - Import canvas maintaining fidelity
  - Asset file handling (images, PDFs)
  - Error recovery patterns

- [x] 3.2.5 Create `/home/jaypaulb/Documents/gh/Canvus-Go-API/examples/error_handling/recovery_patterns.go` [M]
  - Typed error handling with `*canvus.APIError`
  - Error code checking
  - Retry patterns
  - Circuit breaker example
  - Graceful degradation

- [x] 3.2.6 Create `/home/jaypaulb/Documents/gh/Canvus-Go-API/examples/context/cancellation_and_timeouts.go` [S]
  - Context cancellation patterns
  - Request timeouts
  - Deadline management
  - Graceful shutdown

**Acceptance Criteria:**
- All advanced examples compile and are runnable
- Examples demonstrate production-ready patterns
- Each example includes error handling
- Comments explain advanced concepts

---

#### Task Group 3.3: Starter Templates
**Dependencies:** Task Group 3.1
**Effort Total:** M (2-4 hours)

- [x] 3.3.1 Create templates directory [XS]
  - `/home/jaypaulb/Documents/gh/Canvus-Go-API/templates/README.md`
  - Document template usage and customization points

- [x] 3.3.2 Create `/home/jaypaulb/Documents/gh/Canvus-Go-API/templates/minimal_cli.go` [S]
  - Command-line tool structure
  - Flag parsing
  - Configuration from environment
  - Error handling and exit codes
  - Inline documentation for customization

- [x] 3.3.3 Create `/home/jaypaulb/Documents/gh/Canvus-Go-API/templates/web_service.go` [M]
  - HTTP service using SDK
  - Session management
  - Graceful shutdown
  - Health check endpoint
  - Logging middleware

- [x] 3.3.4 Create `/home/jaypaulb/Documents/gh/Canvus-Go-API/templates/batch_job.go` [S]
  - Background job for admin tasks
  - Progress reporting
  - Error aggregation
  - Resumable operations

- [x] 3.3.5 Create `/home/jaypaulb/Documents/gh/Canvus-Go-API/templates/integration_service.go` [M]
  - Microservice pattern
  - Configuration management
  - Metrics collection
  - Health checks
  - Graceful shutdown

**Acceptance Criteria:**
- All templates compile successfully
- Templates demonstrate production patterns
- Each template has customization documentation
- Templates can serve as starting points for real projects

---

### Phase 4: OpenAPI Specification

Create machine-readable API specification for automation tools and MCP server foundation.

#### Task Group 4.1: OpenAPI Core Specification
**Dependencies:** Task Group 1.1
**Effort Total:** XL (> 8 hours)

- [x] 4.1.1 Create OpenAPI specification file structure [S]
  - File: `/home/jaypaulb/Documents/gh/Canvus-Go-API/openapi.yaml`
  - Set up OpenAPI 3.0.0 header
  - Define info section (title, version, description)
  - Define servers section
  - Set up tags for operation categorization

- [x] 4.1.2 Define authentication schemes in OpenAPI [S]
  - API key via Private-Token header
  - Token-based authentication
  - Document security requirements

- [x] 4.1.3 Document Canvas operations (15+ methods) [M]
  - ListCanvases, GetCanvas, CreateCanvas, UpdateCanvas, DeleteCanvas
  - CopyCanvas, GetCanvasBackground, SetCanvasBackground
  - Include request/response schemas
  - Add examples for each operation

- [x] 4.1.4 Document Widget operations (20+ methods) [L]
  - ListWidgets, GetWidget, CreateWidget, UpdateWidget, DeleteWidget
  - All widget types (Note, Browser, Image, PDF, Video, etc.)
  - Widget-specific operations (move, resize)
  - Include schemas for each widget type

- [x] 4.1.5 Document User operations (15+ methods) [M]
  - ListUsers, GetUser, CreateUser, UpdateUser, DeleteUser
  - Token operations (create, list, delete)
  - Group membership operations
  - Include all request/response schemas

- [x] 4.1.6 Document System operations [S]
  - GetLicense, GetSystemInfo
  - Health check operations
  - Include schemas

- [x] 4.1.7 Document Asset operations [M]
  - CreateAsset, GetAsset, DeleteAsset
  - Asset types (images, PDFs, videos)
  - Include binary upload/download patterns

- [x] 4.1.8 Document remaining operations (30+ methods) [L]
  - Anchors, Connectors, Groups
  - Search operations
  - Batch operations
  - Import/Export operations
  - Geometry utilities

**Acceptance Criteria:**
- All 109+ SDK methods documented
- Each operation has request/response schemas
- Examples provided for each operation
- Authentication schemes properly documented

---

#### Task Group 4.2: OpenAPI Schemas and Validation
**Dependencies:** Task Group 4.1
**Effort Total:** M (2-4 hours)

- [x] 4.2.1 Define component schemas for all types [L]
  - Canvas, Widget, User, Asset types
  - Match Go struct definitions exactly
  - Include all fields with proper types
  - Add descriptions for each field

- [x] 4.2.2 Add request/response examples [M]
  - Example for each major operation
  - Show realistic data values
  - Include error response examples

- [x] 4.2.3 Validate OpenAPI specification [S]
  - Use OpenAPI validator tool
  - Check for schema consistency
  - Verify all references resolve
  - Ensure specification is parseable

- [x] 4.2.4 Create OpenAPI documentation [S]
  - Add README for openapi.yaml usage
  - Document how to use with code generators
  - Explain relationship to MCP server

**Acceptance Criteria:**
- All schemas match SDK Go types
- Specification passes validation
- Examples are realistic and complete
- Documentation explains usage

---

### Phase 5: CLI Separation

Migrate CLI to independent repository with binary distribution.

#### Task Group 5.1: New Repository Setup
**Dependencies:** Task Groups 1.1, 1.2, 1.3
**Effort Total:** M (2-4 hours)

> **NOTE: Manual GitHub Repository Creation Required**
>
> These tasks require manual creation of a GitHub repository. All template files have been prepared and are ready to copy.
>
> **Template Files Location:** `/home/jaypaulb/Documents/gh/Canvus-Go-API/cli-repo-template/`
> **Setup Guide:** `/home/jaypaulb/Documents/gh/Canvus-Go-API/docs/CLI_REPOSITORY_SETUP.md`
>
> The setup guide contains complete step-by-step instructions for creating the repository, configuring settings, and initializing with the prepared templates.

- [ ] 5.1.1 Create new GitHub repository `github.com/jaypaulb/canvus-cli` [S]
  - **MANUAL STEP**: Create repository at https://github.com/new
  - **Template**: See setup guide Section 1
  - Initialize with README
  - Add LICENSE file
  - Set up repository settings

- [ ] 5.1.2 Set up go.mod for CLI repository [S]
  - **Template**: `/home/jaypaulb/Documents/gh/Canvus-Go-API/cli-repo-template/go.mod`
  - Module name: `github.com/jaypaulb/canvus-cli`
  - Add dependency: `require github.com/jaypaulb/Canvus-Go-API v0.1.0`
  - Run `go mod tidy`

- [ ] 5.1.3 Create CLI repository structure [S]
  - **Templates**: All files in `/home/jaypaulb/Documents/gh/Canvus-Go-API/cli-repo-template/`
  - `/cmd/canvus/main.go` - main entry point
  - `/internal/` - internal packages
  - `/.github/workflows/` - CI/CD
  - `/docs/` - CLI-specific documentation

- [ ] 5.1.4 Create CLI README.md [S]
  - **Template**: `/home/jaypaulb/Documents/gh/Canvus-Go-API/cli-repo-template/README.md`
  - Installation instructions (binary, go install)
  - Quick usage examples
  - Configuration (environment variables)
  - Link to SDK documentation
  - Badges (Go version, license, release)

**Acceptance Criteria:**
- New repository is created and accessible
- go.mod correctly depends on SDK
- Repository structure follows Go conventions
- README provides clear installation instructions

**Prepared Templates:**
- [x] `/cli-repo-template/go.mod` - Go module configuration
- [x] `/cli-repo-template/README.md` - Complete CLI documentation
- [x] `/cli-repo-template/LICENSE` - MIT license
- [x] `/cli-repo-template/cmd/canvus/main.go` - Placeholder entry point
- [x] `/cli-repo-template/.github/workflows/ci.yml` - CI/CD workflow
- [x] `/docs/CLI_REPOSITORY_SETUP.md` - Complete setup guide

---

#### Task Group 5.2: Code Migration
**Dependencies:** Task Group 5.1
**Effort Total:** M (2-4 hours)

- [ ] 5.2.1 Migrate CLI code from SDK repository [M]
  - Copy `/home/jaypaulb/Documents/gh/Canvus-Go-API/cmd/canvus-cli/` contents
  - Update imports to use `github.com/jaypaulb/Canvus-Go-API/canvus`
  - Update package declarations
  - Verify code compiles in new location

- [ ] 5.2.2 Update CLI to use SDK as external dependency [S]
  - Remove any internal SDK references
  - Use only public SDK APIs
  - Ensure no import cycles

- [ ] 5.2.3 Add CLI-specific configuration [S]
  - Configuration file support
  - Environment variable handling
  - Default values documentation

- [ ] 5.2.4 Create CLI CHANGELOG.md [S]
  - Initial version entry
  - Follow Keep a Changelog format
  - Document CLI-specific features

**Acceptance Criteria:**
- CLI compiles and runs from new repository
- All CLI commands function correctly
- CLI uses SDK as external dependency
- No SDK internal packages referenced

---

#### Task Group 5.3: Binary Distribution
**Dependencies:** Task Group 5.2
**Effort Total:** M (2-4 hours)

- [ ] 5.3.1 Create GitHub Actions workflow for CLI releases [M]
  - File: `.github/workflows/release.yml` in CLI repo
  - Build binaries on tag push
  - Target platforms: macOS (Intel/ARM), Linux (x86_64/ARM), Windows
  - Upload binaries to GitHub release

- [ ] 5.3.2 Create GoReleaser configuration [M]
  - File: `.goreleaser.yaml` in CLI repo
  - Configure build matrix for all platforms
  - Set up archive formats (tar.gz, zip)
  - Configure checksums

- [ ] 5.3.3 Test binary builds for all platforms [S]
  - Verify builds complete successfully
  - Test binaries on available platforms
  - Verify checksums

- [ ] 5.3.4 Document installation methods [S]
  - Binary download from GitHub Releases
  - `go install github.com/jaypaulb/canvus-cli/cmd/canvus@latest`
  - Platform-specific instructions

**Acceptance Criteria:**
- Binaries build for all 5 platforms
- GitHub Actions workflow creates releases
- Installation methods documented
- Checksums provided for verification

---

#### Task Group 5.4: Package Manager Integration (Future)
**Dependencies:** Task Group 5.3
**Effort Total:** S (1-2 hours)

- [ ] 5.4.1 Create Homebrew formula [S]
  - File: `canvus.rb` (for future submission)
  - Follow Homebrew formula conventions
  - Test locally with `brew install --build-from-source`

- [ ] 5.4.2 Create Scoop manifest for Windows [XS]
  - File: `canvus.json`
  - Follow Scoop manifest format
  - Document installation process

- [ ] 5.4.3 Document package manager installation (future) [XS]
  - Add to README when formulas are accepted
  - Instructions for Homebrew tap
  - Instructions for Scoop bucket

**Acceptance Criteria:**
- Homebrew formula created and tested
- Scoop manifest created
- Future installation methods documented

---

#### Task Group 5.5: SDK Repository Cleanup
**Dependencies:** Task Groups 5.1, 5.2, 5.3
**Effort Total:** S (1-2 hours)

- [x] 5.5.1 Remove CLI code from SDK repository [S]
  - Delete `/home/jaypaulb/Documents/gh/Canvus-Go-API/cmd/canvus-cli/`
  - Update any references to CLI in SDK docs
  - Add note in SDK README pointing to CLI repository

- [x] 5.5.2 Update SDK README with CLI reference [XS]
  - Add link to canvus-cli repository
  - Note that CLI is separate project
  - Update any CLI-related documentation

- [x] 5.5.3 Verify SDK builds without CLI [XS]
  - Run `go build ./...`
  - Run `go test ./...`
  - Verify no broken imports

**Acceptance Criteria:**
- CLI code removed from SDK repository
- SDK compiles and tests pass
- README updated with CLI reference
- Clean separation between SDK and CLI

---

### Phase 6: Polish & Release

Final README enhancement, validation, and v0.1.0 release.

#### Task Group 6.1: README Enhancement
**Dependencies:** Task Groups 2.1, 3.1
**Effort Total:** M (2-4 hours)

- [x] 6.1.1 Restructure `/home/jaypaulb/Documents/gh/Canvus-Go-API/README.md` [M]
  - Problem statement (why Canvus-Go-API exists)
  - Quick feature overview with bullet points
  - Installation instructions (`go get`)
  - Simple example (5-10 lines showing basic usage)
  - Clear links to comprehensive docs in /docs/
  - Use cases: integrations, automation, admin tools, MCP server foundation
  - Key differentiators

- [x] 6.1.2 Add badges to README [XS]
  - Go version badge
  - License badge
  - GoDoc badge (pkg.go.dev link)
  - Latest release badge
  - Build status badge

- [x] 6.1.3 Add API coverage section [S]
  - Highlight 109+ methods
  - List major categories (Users, Canvases, Widgets, etc.)
  - Link to API_REFERENCE.md

- [x] 6.1.4 Add contributing section [XS]
  - Link to CONTRIBUTING.md (create if needed)
  - Issue reporting guidelines
  - Pull request process

**Acceptance Criteria:**
- README serves as effective entry point
- Installation is clear and simple
- Links to detailed docs are prominent
- Badges display correctly

---

#### Task Group 6.2: Validation and Testing
**Dependencies:** All previous task groups
**Effort Total:** M (2-4 hours)

- [x] 6.2.1 Validate all examples compile and run [M]
  - Run `go build` on each example
  - Test examples that can run without server
  - Document any examples requiring server access

- [x] 6.2.2 Validate OpenAPI spec matches SDK [M]
  - Compare method signatures
  - Verify schemas match Go types
  - Check for missing or extra operations
  - Run OpenAPI validator

- [x] 6.2.3 Documentation consistency check [S]
  - Verify links are valid
  - Check code examples compile
  - Ensure version numbers are consistent
  - Verify import paths are correct

- [x] 6.2.4 Integration test documentation [S]
  - Document test requirements
  - Note which tests need server access
  - Create test configuration guide

**Acceptance Criteria:**
- All examples compile successfully
- OpenAPI spec matches SDK implementation
- All documentation links are valid
- Test requirements are documented

---

#### Task Group 6.3: v0.1.0 Release
**Dependencies:** Task Groups 6.1, 6.2
**Effort Total:** S (1-2 hours)

- [x] 6.3.1 Final pre-release checklist [S]
  - Run full test suite
  - Update CHANGELOG.md with all changes
  - Verify version numbers in docs
  - Review README one final time
  - Verify go.mod is correct

- [x] 6.3.2 Create and push v0.1.0 tag [XS]
  - `git tag -a v0.1.0 -m "Initial release"`
  - `git push origin v0.1.0`
  - Verify tag appears on GitHub

- [x] 6.3.3 Create GitHub release for v0.1.0 [S]
  - Use release template
  - Include release notes from CHANGELOG
  - Note MTCS API compatibility
  - Include installation instructions
  - Link to Getting Started guide

- [x] 6.3.4 Verify pkg.go.dev listing [XS]
  - Check package appears on pkg.go.dev
  - Verify documentation renders correctly
  - Note any issues for follow-up
  - **NOTE**: pkg.go.dev indexing may take time; can be triggered manually

- [x] 6.3.5 Announce release [XS]
  - Update any relevant documentation
  - Note version in README
  - Consider announcement channels

**Acceptance Criteria:**
- v0.1.0 tag created and pushed
- GitHub release published
- Package accessible via `go get`
- Documentation visible on pkg.go.dev

---

## Execution Order

**Recommended implementation sequence:**

1. **Phase 1: Foundation** (Task Groups 1.1-1.3)
   - Must be completed first to enable proper distribution
   - Blocks all other phases

2. **Phase 2: Documentation** (Task Groups 2.1-2.2)
   - Can begin immediately after Phase 1
   - Core documentation enables developer success

3. **Phase 3: Examples & Templates** (Task Groups 3.1-3.3)
   - Can run in parallel with Phase 2
   - Examples support documentation

4. **Phase 4: OpenAPI Specification** (Task Groups 4.1-4.2)
   - Can run in parallel with Phases 2-3
   - Independent work stream

5. **Phase 5: CLI Separation** (Task Groups 5.1-5.5)
   - Requires Phase 1 completion (SDK must be tagged first)
   - Can run in parallel with Phases 2-4

6. **Phase 6: Polish & Release** (Task Groups 6.1-6.3)
   - Must wait for all other phases
   - Final integration and release

---

## Parallel Execution Opportunities

The following task groups can be executed in parallel:

**After Phase 1 completes:**
- Task Groups 2.1, 2.2 (Documentation)
- Task Groups 3.1, 3.2, 3.3 (Examples & Templates)
- Task Groups 4.1, 4.2 (OpenAPI)
- Task Groups 5.1, 5.2, 5.3 (CLI - after SDK is tagged)

**This allows for 3-4 parallel work streams:**
1. Documentation writer: Phases 2
2. Examples developer: Phase 3
3. API specification writer: Phase 4
4. CLI maintainer: Phase 5

---

## Effort Summary

| Phase | Effort | Hours (Est.) |
|-------|--------|--------------|
| Phase 1: Foundation | M-L | 6-12 |
| Phase 2: Documentation | L-XL | 6-12 |
| Phase 3: Examples & Templates | L-XL | 8-14 |
| Phase 4: OpenAPI Specification | XL | 12-18 |
| Phase 5: CLI Separation | M-L | 8-12 |
| Phase 6: Polish & Release | M | 4-8 |
| **Total** | | **44-76 hours** |

---

## Key Deliverables Checklist

Upon completion, verify these key deliverables exist:

**Foundation:**
- [x] Updated go.mod with correct module path
- [x] LICENSE file at repository root
- [x] `/home/jaypaulb/Documents/gh/Canvus-Go-API/docs/VERSIONING.md`
- [x] `/home/jaypaulb/Documents/gh/Canvus-Go-API/CHANGELOG.md`
- [x] GitHub Actions release workflow

**Documentation (7 documents):**
- [x] `/home/jaypaulb/Documents/gh/Canvus-Go-API/docs/GETTING_STARTED.md`
- [x] `/home/jaypaulb/Documents/gh/Canvus-Go-API/docs/BEST_PRACTICES.md`
- [x] `/home/jaypaulb/Documents/gh/Canvus-Go-API/docs/TROUBLESHOOTING.md`
- [x] `/home/jaypaulb/Documents/gh/Canvus-Go-API/docs/COMPATIBILITY.md`
- [x] `/home/jaypaulb/Documents/gh/Canvus-Go-API/docs/API_REFERENCE.md`
- [x] `/home/jaypaulb/Documents/gh/Canvus-Go-API/docs/EXAMPLES.md`
- [x] `/home/jaypaulb/Documents/gh/Canvus-Go-API/docs/RELEASING.md`

**Examples (10+ examples):**
- [x] `/home/jaypaulb/Documents/gh/Canvus-Go-API/examples/getting_started/main.go`
- [x] `/home/jaypaulb/Documents/gh/Canvus-Go-API/examples/authentication/api_key.go`
- [x] `/home/jaypaulb/Documents/gh/Canvus-Go-API/examples/authentication/login_password.go`
- [x] `/home/jaypaulb/Documents/gh/Canvus-Go-API/examples/canvases/create_and_manage.go`
- [x] `/home/jaypaulb/Documents/gh/Canvus-Go-API/examples/widgets/create_and_search.go`
- [x] `/home/jaypaulb/Documents/gh/Canvus-Go-API/examples/users/provision_and_manage.go`
- [x] `/home/jaypaulb/Documents/gh/Canvus-Go-API/examples/batch/bulk_operations.go`
- [x] `/home/jaypaulb/Documents/gh/Canvus-Go-API/examples/import_export/round_trip.go`
- [x] `/home/jaypaulb/Documents/gh/Canvus-Go-API/examples/error_handling/recovery_patterns.go`
- [x] `/home/jaypaulb/Documents/gh/Canvus-Go-API/examples/context/cancellation_and_timeouts.go`

**Templates (4 templates):**
- [x] `/home/jaypaulb/Documents/gh/Canvus-Go-API/templates/minimal_cli.go`
- [x] `/home/jaypaulb/Documents/gh/Canvus-Go-API/templates/web_service.go`
- [x] `/home/jaypaulb/Documents/gh/Canvus-Go-API/templates/batch_job.go`
- [x] `/home/jaypaulb/Documents/gh/Canvus-Go-API/templates/integration_service.go`

**OpenAPI:**
- [x] `/home/jaypaulb/Documents/gh/Canvus-Go-API/openapi.yaml`
- [x] `/home/jaypaulb/Documents/gh/Canvus-Go-API/docs/OPENAPI.md`

**CLI Repository:**
- [ ] `github.com/jaypaulb/canvus-cli` repository created
- [ ] CLI code migrated and working
- [ ] Binary builds for all platforms
- [x] CLI removed from SDK repository

**Release:**
- [x] Enhanced README.md with badges and links
- [x] v0.1.0 tag created and pushed
- [x] GitHub release published
- [x] Package visible on pkg.go.dev (indexing in progress)

---

## Notes

- All file paths are absolute to avoid confusion
- Effort estimates are rough guides; actual time may vary
- Phase 4 (OpenAPI) is the largest effort due to 109+ methods
- CLI separation requires careful coordination with SDK versioning
- Testing against actual Canvus server required for full validation
