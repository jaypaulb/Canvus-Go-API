# CLI Repository Setup Guide

This guide provides step-by-step instructions for creating the `canvus-cli` repository as a separate project from the Canvus Go SDK.

## Overview

The CLI is being separated from the SDK to allow:
- Independent versioning and release cycles
- Binary distribution without requiring Go installation
- Cleaner separation of concerns (library vs. application)
- Easier maintenance and contribution workflows

## Prerequisites

Before starting, ensure you have:
- GitHub account with permissions to create repositories under `jaypaulb`
- Git installed locally
- Go 1.21+ installed
- Template files from `/home/jaypaulb/Documents/gh/Canvus-Go-API/cli-repo-template/`

---

## Step 1: Create GitHub Repository

### 1.1 Create the Repository on GitHub

1. Go to https://github.com/new
2. Configure the repository:
   - **Repository name**: `canvus-cli`
   - **Description**: "Command-line interface for the Canvus collaborative workspace API"
   - **Visibility**: Public
   - **Initialize**: Do NOT add README, .gitignore, or license (we'll add these from templates)
3. Click "Create repository"

### 1.2 Configure Repository Settings

After creation, configure these settings:

#### General Settings (Settings > General)
- **Features**:
  - [x] Issues
  - [x] Projects (optional)
  - [x] Discussions (optional)
  - [ ] Wiki (not needed)
- **Pull Requests**:
  - [x] Allow merge commits
  - [x] Allow squash merging (default)
  - [ ] Allow rebase merging
  - [x] Automatically delete head branches

#### Branch Protection (Settings > Branches)
Add rule for `main` branch:
- [x] Require a pull request before merging
- [x] Require status checks to pass before merging
  - Required checks: `build`, `test`
- [x] Require branches to be up to date before merging
- [ ] Do not require approvals for single developer

#### Actions (Settings > Actions > General)
- Allow all actions and reusable workflows
- Workflow permissions: Read and write permissions

---

## Step 2: Clone and Initialize

### 2.1 Clone the Empty Repository

```bash
cd ~/Documents/gh
git clone https://github.com/jaypaulb/canvus-cli.git
cd canvus-cli
```

### 2.2 Copy Template Files

```bash
# Copy all template files from the SDK repository
cp -r /home/jaypaulb/Documents/gh/Canvus-Go-API/cli-repo-template/* .
cp -r /home/jaypaulb/Documents/gh/Canvus-Go-API/cli-repo-template/.github .
```

### 2.3 Verify Directory Structure

After copying, your structure should be:
```
canvus-cli/
├── .github/
│   └── workflows/
│       └── ci.yml
├── cmd/
│   └── canvus/
│       └── main.go
├── docs/
│   └── .gitkeep
├── internal/
│   └── .gitkeep
├── go.mod
├── LICENSE
└── README.md
```

---

## Step 3: Initialize Go Module

### 3.1 Tidy Dependencies

```bash
# Ensure Go module is properly initialized
go mod tidy
```

### 3.2 Verify SDK Dependency

The `go.mod` file should show:
```
module github.com/jaypaulb/canvus-cli

go 1.21

require github.com/jaypaulb/Canvus-Go-API v0.1.0
```

**Note**: If v0.1.0 is not yet released, you can use a replace directive temporarily:
```go
require github.com/jaypaulb/Canvus-Go-API v0.0.0

replace github.com/jaypaulb/Canvus-Go-API => ../Canvus-Go-API
```

Remove the replace directive after v0.1.0 is published.

### 3.3 Verify Build

```bash
go build ./...
```

---

## Step 4: Commit and Push

### 4.1 Initial Commit

```bash
git add .
git commit -m "Initial commit: CLI repository setup

- Add go.mod with SDK dependency
- Add README with installation instructions
- Add MIT LICENSE
- Add CI workflow for testing and building
- Add placeholder main.go entry point"
```

### 4.2 Push to GitHub

```bash
git push -u origin main
```

---

## Step 5: Migrate CLI Code

After the initial setup is pushed, migrate the actual CLI code:

### 5.1 Copy CLI Source Files

```bash
# Copy CLI code from SDK repository
cp /home/jaypaulb/Documents/gh/Canvus-Go-API/cmd/canvus-cli/*.go cmd/canvus/
```

### 5.2 Update Imports

In each copied file, update the import path:

**From:**
```go
import "canvus-go-api/canvus"
```

**To:**
```go
import "github.com/jaypaulb/Canvus-Go-API/canvus"
```

### 5.3 Update Package Declaration

In `cmd/canvus/main.go`, ensure:
```go
package main
```

### 5.4 Build and Test

```bash
go build ./cmd/canvus
./canvus --help
```

---

## Step 6: Set Up Releases

### 6.1 Create First Release Tag

After code is migrated and tested:

```bash
git tag -a v0.1.0 -m "Initial CLI release"
git push origin v0.1.0
```

### 6.2 GitHub Release

1. Go to https://github.com/jaypaulb/canvus-cli/releases
2. Click "Draft a new release"
3. Choose tag: `v0.1.0`
4. Release title: `v0.1.0 - Initial Release`
5. Description:
```markdown
## Canvus CLI v0.1.0

Initial release of the Canvus command-line interface.

### Features
- Canvas management (list, get, create, delete)
- Widget operations (create, update, delete)
- User management
- Import/export functionality
- API key and token authentication

### Installation

**Using Go:**
```bash
go install github.com/jaypaulb/canvus-cli/cmd/canvus@latest
```

**Binary Download:**
Download the appropriate binary for your platform from the Assets below.

### Requirements
- Canvus server with API access
- API key or user credentials

### Documentation
See the [SDK documentation](https://github.com/jaypaulb/Canvus-Go-API) for complete API reference.
```

6. Attach binaries (built by CI workflow)
7. Click "Publish release"

---

## File Contents Reference

All template files are located in:
```
/home/jaypaulb/Documents/gh/Canvus-Go-API/cli-repo-template/
```

### go.mod
Module configuration with SDK dependency.

### README.md
Complete CLI documentation with:
- Installation methods
- Quick start guide
- Command reference
- Configuration options
- Links to SDK documentation

### LICENSE
MIT License matching the SDK.

### cmd/canvus/main.go
Placeholder entry point (replace with actual CLI code).

### .github/workflows/ci.yml
GitHub Actions workflow for:
- Running tests
- Building binaries
- Creating releases

---

## Post-Setup Tasks

After completing the setup:

1. **Verify CI Workflow**: Check that GitHub Actions runs successfully
2. **Test Installation**: Verify `go install` works
3. **Update SDK README**: Add link to CLI repository
4. **Remove CLI from SDK**: Delete `/cmd/canvus-cli/` from SDK after migration is complete
5. **Create GoReleaser Config**: Set up `.goreleaser.yaml` for binary distribution

---

## Troubleshooting

### go mod tidy fails with "module not found"

If the SDK v0.1.0 is not yet published:
1. Add a replace directive in go.mod (see Step 3.2)
2. Or wait until SDK v0.1.0 is tagged and pushed

### CI workflow fails

Check:
- Workflow file syntax (validate with `yamllint`)
- Go version matches requirements (1.21+)
- SDK dependency is accessible

### Binary builds fail

Ensure:
- All imports use the full module path
- No internal SDK packages are referenced
- Platform-specific code is properly guarded with build tags

---

## Repository URLs

After setup, these URLs will be available:

- **Repository**: https://github.com/jaypaulb/canvus-cli
- **Releases**: https://github.com/jaypaulb/canvus-cli/releases
- **Go Package**: https://pkg.go.dev/github.com/jaypaulb/canvus-cli
- **Issues**: https://github.com/jaypaulb/canvus-cli/issues

---

## Related Documentation

- [SDK Documentation](https://github.com/jaypaulb/Canvus-Go-API)
- [Getting Started Guide](./GETTING_STARTED.md)
- [API Reference](./API_REFERENCE.md)
- [Release Process](./RELEASING.md)
