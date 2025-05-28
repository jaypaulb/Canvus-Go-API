# Canvus Go SDK Development Tasks

This document is the authoritative project plan. **Update this file with every new plan, approach, or after each major prompt cycle, especially before committing to git.**

## Project Objective
- Expand the Go API library to support all Canvus API endpoints and features.
- Build a full-featured Go SDK for Canvus, including utilities, documentation, and developer tools.

---

## Task List

### 1. Project Setup
- [ ] Initialize a new git repository (GitHub, using CLI)
- [ ] Set up standard branching and commit practices
- [ ] Confirm Windows/PowerShell as the development environment
- [ ] Add `.gitignore` for Go and project-specific files

### 2. API Coverage Analysis
- [ ] Extract all endpoints and features from the Python client
- [ ] List all public methods in the Go library
- [ ] Create a coverage matrix (Python vs Go)
- [ ] Identify missing endpoints/features in Go

### 3. Go API Library Expansion
For each missing endpoint/feature:
- [ ] Define Go method signature and data structures
- [ ] Implement the method
- [ ] Add error handling and authentication
- [ ] Write unit tests
- [ ] Update documentation

### 4. Build the Go SDK
- [ ] Identify and design SDK utilities (CLI tools, helpers, etc.)
- [ ] Implement SDK features
- [ ] Add code samples and templates
- [ ] Write comprehensive documentation

### 5. Developer Experience & Release
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