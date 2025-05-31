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

**2024-06-10 Summary:**

- Confirmed existing git repository.
- Added Go-specific rules to `.gitignore` (Windows/PowerShell compatible).
- Created `CONTRIBUTING.md` with branching, commit, and Windows compatibility guidelines.

### 2. MCS REST API Analysis

- [x] Collect and review the official MT Canvus Server (MCS) REST API documentation/spec
- [x] List all available API endpoints and their parameters
- [ ] Identify required abstractions/utilities to improve developer experience (e.g., batching, streaming, error handling)
- [ ] Map API endpoints to planned SDK features and abstractions
- [ ] Document findings and update the PRD if needed

**2024-06-10 Summary:**
- Collected and reviewed the official Canvus API documentation in `Canvus API Docs/`.
- Created a complete, explicit, and fully enumerated list of all API endpoints in `Canvus_API_Endpoint_List.md`.
- Created a comprehensive markdown table of all endpoints in `Canvus_API_Endpoint_Table.md`.

### 3. API Coverage Analysis

- [ ] Extract all endpoints and features from the Python client
- [ ] List all public methods in the Go library
- [ ] Create a coverage matrix (Python vs Go)
- [ ] Identify missing endpoints/features in Go

### 4. Go API Library Expansion

For each missing endpoint/feature:

- [ ] Define Go method signature and data structures
- [ ] Implement the method
- [ ] Add error handling and authentication
- [ ] Write unit tests
- [ ] Update documentation

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

## Ongoing Practices

- **Update this file with every new plan or after each major prompt cycle.**
- **Before every git commit, summarize what was done in this file.**
- Use feature branches for new work; merge via pull requests.
- All development and commands must be Windows/PowerShell compatible.
- No Linux commands or assumptions.
