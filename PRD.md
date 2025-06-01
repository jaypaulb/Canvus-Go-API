# Canvus Go SDK – Product Requirements Document (PRD)

## 1. Project Overview

- **Goal:**  
  Build a robust, idiomatic Go SDK for the Canvus API, providing full API coverage, developer-friendly abstractions, and a seamless experience for Go developers on Windows/PowerShell.
- **Scope:**  
  - Expand the existing Go API library to cover all Canvus API endpoints and features.
  - Package the library as a full SDK, including utilities, documentation, and developer tools.
  - Ensure all development is Windows/PowerShell compatible.
  - Maintain high code quality, test coverage, and clear documentation.

---

## 2. Stakeholders

- **Primary:**  
  - Go developers integrating with Canvus.
  - Internal Canvus engineering and QA teams.
- **Secondary:**  
  - Technical writers (for docs).
  - DevOps (for packaging and release).

---

## 3. Functional Requirements

### 3.1. API Coverage

- The SDK must provide Go methods for every public Canvus API endpoint, including:
  - Canvas management (CRUD)
  - Folder management
  - Widget operations (notes, images, browsers, videos, PDFs, connectors, anchors)
  - User and token management
  - Client and workspace operations
  - Server info/config endpoints
  - Asset management (images, videos, PDFs, uploads, backgrounds, mipmaps, etc.)
  - Parenting (patching parent ID) must be implemented but is not to be tested due to a known server bug.
- All request/response models must be idiomatic Go structs, with proper error handling.

### 3.2. SDK Features

- **Authentication:**  
  - Support API key authentication (from env or config).
- **Utilities:**  
  - CLI tools for common tasks (e.g., token management, canvas export/import).
  - Helper functions for common workflows (e.g., batch operations, event subscriptions).
- **Documentation:**  
  - Godoc comments for all exported types and functions.
  - Comprehensive README with setup, usage, and examples.
  - Code samples for common use cases.
- **Testing:**  
  - Unit tests for all major features, following Go testing conventions.
  - Test coverage for happy path, edge cases, and error handling.

### 3.3. Developer Experience

- Easy installation via Go modules.
- Windows/PowerShell compatibility for all scripts and commands.
- Clear error messages and logging.
- Consistent, idiomatic Go API design.

### 3.4. Endpoint Implementation & Testing Order

The following order must be followed for endpoint implementation and testing, to ensure safe, isolated, and reliable test environments:

1. **System Management Endpoints**: Users, Access Tokens, Groups, Canvas Folders, Server Config, License (no activation tests), Audit Log, Server Info
2. **Canvas Endpoints**: All canvas actions (CRUD, move, copy, permissions, etc.)
3. **Client & Workspace Endpoints**: All client and workspace actions, including launching the MT-Canvus-Client with a canvas URL and verifying client/workspace state
4. **Widget & Asset Endpoints**: All widget and asset actions, in dependency order (simple elements first, then those requiring files or references, then read-only endpoints)
5. **Parenting**: Implement but do not test parenting (patching parent ID)

### 3.5. Subscription & Buffering Requirements

- For endpoints supporting subscriptions (with `?subscription`), the SDK must:
  - Support initial GET and real-time update streaming (one JSON per line, CR as keep-alive)
  - Provide a function to filter updates for specific elements
  - Provide a buffered subscription handler: only emit updates after a configurable period of inactivity, to reduce noise from rapid, sequential updates

### 3.6. Testing & Cleanup Policy

- All tests must clean up (permanently delete) resources they create, even on failure. Moving to trash is not sufficient.
- Each test must use unique resource names/IDs to avoid collisions and ensure safe cleanup.

## API/SDK Design Notes
- All JSON field names sent to the Canvus API must be lowercase (e.g., 'name', 'mode'), matching the API's requirements. This is critical for PATCH/POST requests to work as expected.

---

## 4. Non-Functional Requirements

- **Performance:**  
  - Efficient HTTP requests, minimal allocations, and fast response parsing.
- **Reliability:**  
  - Robust error handling and retries for transient failures.
- **Maintainability:**  
  - Modular code structure, clear separation of concerns.
  - Adherence to Go best practices and project coding standards.
- **Security:**  
  - Secure handling of API keys and sensitive data.
- **Versioning:**  
  - Semantic versioning for releases.

---

## 5. Constraints

- All development and documentation must be Windows/PowerShell compatible.
- No Linux shell commands or assumptions.
- Use only well-maintained, widely adopted Go modules.
- All code and documentation must be committed to git, following the project's branching and commit message conventions.
- The `tasks.md` file must be updated with every new plan, approach, or after each major prompt cycle, and before every git commit.

---

## 6. Milestones & Deliverables

1. **Project Setup**
   - GitHub repo, `.gitignore`, initial planning docs.
2. **API Coverage Analysis**
   - Coverage matrix (Python vs Go), gap analysis.
3. **Go API Library Expansion**
   - Full endpoint coverage, models, and tests.
4. **SDK Utilities & CLI**
   - Helper tools, code samples, and documentation.
5. **Release**
   - Go module packaging, README, and tagged release.

---

## 7. Success Criteria

- 100% coverage of public Canvus API endpoints.
- All code passes unit tests and is reviewed for idiomatic Go style.
- SDK is installable and usable by Go developers on Windows.
- Documentation is clear, complete, and up-to-date.
- No Linux-specific commands or scripts in the codebase.

---

## 8. Open Questions / To Be Determined

- Are there any private/internal Canvus endpoints that should be included?
- **Minimum Go version to support:** Go 1.24.1 (per developer environment)
- **Frameworks/Libraries:**
  - Use standard `net/http` for HTTP operations
  - Use `github.com/go-playground/validator/v10` for data validation if needed
  - Use `GORM` or `sqlc` for ORM/database interaction if required
  - Follow all conventions in @10-golang-coding-standards.mdc
- **Release cadence:** SDK will be updated twice a year, aligned with MT Canvus releases

---

## SDK Abstractions & Utilities

See [Abstractions.md](./Abstractions.md) for a detailed description of planned SDK abstractions, utilities, and advanced features.

## MCS API Feature Requests

See [MCS-Feature-Requests.md](./MCS-Feature-Requests.md) for proposed improvements and suggestions for the MCS API.
