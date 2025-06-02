# Canvus Go SDK Development Tasks

This document is the authoritative project plan. **Update this file with every new plan, approach, or after each major prompt cycle, especially before committing to git.**

## Project Objective

- Expand the Go API library to support all Canvus API endpoints and features.
- Build a full-featured Go SDK for Canvus, including utilities, documentation, and developer tools.

---

## Task List

### 1. Project Setup

- [x] Initialize a new git repository (GitHub, using CLI)
- [x] Set up standard branching and commit practices
- [x] Confirm Windows/PowerShell as the development environment
- [x] Add `.gitignore` for Go and project-specific files
- [x] Confirmed: All JSON field names in requests must be lowercase to match the Canvus API. This is required for PATCH/POST to work (e.g., canvas renaming).

**2024-06-10 Summary:**

- Confirmed existing git repository.
- Added Go-specific rules to `.gitignore` (Windows/PowerShell compatible).
- Created `CONTRIBUTING.md` with branching, commit, and Windows compatibility guidelines.

### 2. MCS REST API Analysis

- [x] Collect and review the official MT Canvus Server (MCS) REST API documentation/spec
- [x] List all available API endpoints and their parameters
- [x] Identify required abstractions/utilities to improve developer experience (see [Abstractions.md](./Abstractions.md))
- [x] Map API endpoints to planned SDK features and abstractions (see [Abstractions.md](./Abstractions.md))
- [x] Document findings and update the PRD if needed (see [PRD.md](./PRD.md))

**2024-06-13 Summary:**

- Completed identification, mapping, and documentation of required SDK abstractions and utilities. See [Abstractions.md](./Abstractions.md) for details and [PRD.md](./PRD.md) for cross-references.

### 3. API Coverage Analysis

- [x] Extract all endpoints and features from the Docs

### 4. Go API Library Expansion

For each missing endpoint/feature:

- [ ] Define Go method signature and data structures
- [ ] Implement the method
- [ ] Add error handling and authentication
- [ ] Write unit tests
- [ ] Update documentation

#### Endpoint Implementation & Testing Order (2024-06-13)

1. **System Management Endpoints**
    - [x] Users: implement and test all actions
    - [x] Access Tokens: implement and test all actions
    - [x] Groups: implement and test all actions
    - [x] Canvas Folders: implement and test all actions
    - [x] Server Config: implement and test all actions
    - [x] License: implement (do not test activation)
    - [x] Audit Log: implement and test all actions
    - [x] Server Info: implement and test all actions
2. **Canvas Endpoints**
    - [/] Implement and test all Canvas actions (CRUD, move, copy, permissions, etc.)
3. **Client & Workspace Endpoints**
    - [ ] Implement and test all Client actions
    - [ ] Implement and test all Workspace actions
    - [ ] Implement logic to launch MT-Canvus-Client with canvas URL for integration tests [ABROGATED: Not part of Go SDK; for integration/E2E testing only. For API tests, ensure a client is active and running. See documentation.]
    - [ ] Refactor open-canvas code to poll get workspaces and verify canvas ID update (server does not send full response; must check for matching canvas ID)
4. **Widget & Asset Endpoints**
    - [ ] Implement and test all Widget actions (Notes, Anchors, VideoInputs, VideoOutputs, Color Presets, etc.)
    - [ ] Implement and test all Asset actions (Images, Videos, PDFs, Uploads, Connectors, Backgrounds, MipMaps, Assets, etc.)
    - [ ] Implement Parenting (patching parent ID) functions (do not test due to known bug)
    - [ ] Implement and test all read-only endpoints (Widgets, Annotations)

#### Testing & Cleanup Policy

- [ ] All tests must clean up (permanently delete) resources they create, even on failure. Moving to trash is not sufficient.
- [ ] Each test must use unique resource names/IDs to avoid collisions and ensure safe cleanup.

### Required Abstractions/Utilities

- **Authentication:** API key (from env/config), secure handling
- **Context support:** All requests accept `context.Context`
- **Error handling:** Centralized, idiomatic Go error types
- **Pagination/streaming:** Helpers for paginated and streaming endpoints
- **Strong typing:** Request/response models as Go structs
- **Modular structure:** Packages by resource (canvases, folders, widgets, users, etc.)
- **CLI utilities:** (Optional) for common workflows
- **Documentation:** Godoc, README, code samples
- **Testing:** Unit/integration tests for all features

### 5. Build the Go SDK

- [ ] Identify and design SDK utilities (CLI tools, helpers, etc.)
- [ ] Implement SDK features
- [ ] Add code samples and templates
- [ ] Write comprehensive documentation

### 6. Developer Experience & Release

- [ ] Package as a Go module
- [ ] Ensure easy installation and usage
- [ ] Finalize documentation and examples
- [ ] Tag and release initial version

---

## 2024-06-11: Greenfield Go SDK Approach & Planning

### New Approach

- The Canvus Go SDK will be built from scratch, following modern Go community best practices.
- Old Python/Go libraries are not used as a roadmap; they serve only as historical context.
- The SDK will be idiomatic, modular, and designed for the broader Go developer community.
- Focus: developer experience, full API coverage, strong documentation, and extensibility.

### Required Abstractions/Utilities

- **Authentication:** API key (from env/config), secure handling
- **Context support:** All requests accept `context.Context`
- **Error handling:** Centralized, idiomatic Go error types
- **Pagination/streaming:** Helpers for paginated and streaming endpoints
- **Strong typing:** Request/response models as Go structs
- **Modular structure:** Packages by resource (canvases, folders, widgets, users, etc.)
- **CLI utilities:** (Optional) for common workflows
- **Documentation:** Godoc, README, code samples
- **Testing:** Unit/integration tests for all features

### API Endpoint → Planned Go SDK Feature Mapping (WIP)

| API Resource         | HTTP Method & Path                              | Planned Go SDK Method Signature                |
|---------------------|-------------------------------------------------|------------------------------------------------|
| Canvases            | GET    /canvases                                | func (c *Client) ListCanvases(ctx context.Context) ([]Canvas, error) |
| Canvases            | GET    /canvases/:id                            | func (c *Client) GetCanvas(ctx context.Context, id string) (Canvas, error) |
| Canvases            | POST   /canvases                                | func (c *Client) CreateCanvas(ctx context.Context, req CreateCanvasRequest) (Canvas, error) |
| Canvases            | PATCH  /canvases/:id                            | func (c *Client) UpdateCanvas(ctx context.Context, id string, req UpdateCanvasRequest) (Canvas, error) |
| Canvases            | DELETE /canvases/:id                            | func (c *Client) DeleteCanvas(ctx context.Context, id string) error |
| ...                 | ...                                             | ...                                            |

> This table will be expanded to cover all endpoints in the official API list. Each resource group will have its own Go type(s) and methods, following idiomatic Go SDK design.

---

## Authentication Methods (SDK & API)

There are two primary ways to authenticate to the Canvus server:

1. **Username/Password Login**
   - Endpoint: `POST /users/login`
   - The client sends a username (email) and password to the endpoint.
   - If password authentication is enabled, the server issues a temporary access token (valid for 24 hours).
   - This token is used for subsequent authenticated requests.
   - Example:

     ```json
     POST /users/login
     { "email": "alice@example.com", "password": "BBBB" }
     ```

   - Response includes a `token` and user info.

2. **Access Token**
   - Endpoint: `POST /access-tokens` (or via Canvus web UI)
   - An access token is created via the API or UI.
   - This token does **not expire** and can be used directly for authentication by including it in the `Private-Token` header.
   - You can also POST an existing token to `/users/login` to validate and prolong its lifetime.

3. **Sign Out**
   - Endpoint: `POST /users/logout`
   - Invalidates the provided access token. If no token is provided in the body, the `Private-Token` header is used.

### SDK Client Authentication Options

The Go SDK client supports three authentication options:

1. Username/password login (temporary token)
2. Static access token (long-lived)
3. Token validation/refresh (prolongs token lifetime)

All authentication logic must be tested for both login and static token flows.

---

## Client Instantiation & Authentication Patterns

The SDK must support three primary client instantiation patterns:

1. **test_client**
   - Uses the main client (from settings) to create and activate a new test user.
   - Logs in as the test user to obtain a temporary token (API key/PrivateToken) via `/users/login`.
   - All actions in the session use this token (sent as `Private-Token` header).
   - On completion, logs out (`/users/logout`, invalidates token) and deletes the test user.

2. **user_client**
   - Logs in as an existing user (by email and password) to obtain a temporary token via `/users/login`.
   - All actions use this token for the session (sent as `Private-Token` header).
   - On completion, logs out to invalidate the token.

3. **client**
   - Uses credentials from the settings/config file for all actions.
   - No automatic cleanup or user/token creation.

**Notes:**

- The terms "API key" and "PrivateToken" are used interchangeably; all authentication headers use `Private-Token`.
- Session cleanup for temporary tokens is handled by calling `/users/logout`.

---

## Ongoing Practices

- **Update this file with every new plan or after each major prompt cycle.**
- **Before every git commit, summarize what was done in this file.**
- Use feature branches for new work; merge via pull requests.
- All development and commands must be Windows/PowerShell compatible.
- No Linux commands or assumptions.

**2024-06-12 Progress Update:**

- Refactored `doRequest` to support query parameters for HTTP requests.
- Updated `ListCanvases` and `GetCanvas` to accept and encode options as query parameters.
- Implemented binary response handling for `GetCanvasPreview`.
- Next: Add or update unit tests for Canvas methods, especially `GetCanvasPreview`.

**2024-06-14 Progress Update:**

- Implemented DeleteFolder method in Go SDK.
- Added integration test for folder creation and deletion.
- Verified folder is removed from ListFolders after deletion.

**2024-06-15 Progress Update:**

- Expanded `Canvas` struct and added all request/response types for Canvas API.
- Implemented all Canvas API methods in `canvases.go` (CRUD, move, copy, permissions, preview, demo state, etc.).
- Next: Add and run comprehensive tests for Canvas lifecycle and error cases in `canvases_test.go`.

- [x] Implement all workspace abstractions and functions:
    - [x] Flexible workspace selection (index, name, user, default)
    - [x] List/Get/Update workspace
    - [x] Toggle info panel and pinned state
    - [x] Set viewport (by coords)
    - [ ] (Revisit: Set viewport by widget, Open canvas with viewport centering) [Return to after widget endpoints are implemented]
- [x] Refactor: Rename Client struct and all related code to Session to avoid confusion with API clients (BREAKING CHANGE)

**2024-06-15 Summary:**
- Completed Client-to-Session refactor across all code and tests.
- All tests pass. Ready to proceed to the next endpoints.
