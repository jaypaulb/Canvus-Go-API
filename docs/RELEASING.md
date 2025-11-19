# Release Process

This document describes the step-by-step process for releasing a new version of the Canvus Go SDK.

## Overview

Releases follow the [versioning scheme](VERSIONING.md) and use GitHub releases with annotated tags.

## Release Checklist

### Pre-Release Validation

Before starting the release process, complete all validation steps:

- [ ] **Run full test suite**
  ```bash
  go test ./canvus/... -v
  ```
  All tests must pass.

- [ ] **Verify build succeeds**
  ```bash
  go build ./...
  ```
  No compilation errors.

- [ ] **Run go mod verification**
  ```bash
  go mod verify
  go mod tidy
  ```
  Ensure go.mod and go.sum are clean.

- [ ] **Check for security issues**
  ```bash
  go list -json -m all | nancy sleuth  # if using nancy
  # or
  govulncheck ./...  # if using govulncheck
  ```
  Address any critical vulnerabilities.

- [ ] **Validate examples compile** (when examples directory exists)
  ```bash
  for f in examples/**/*.go; do go build "$f"; done
  ```

- [ ] **Review documentation**
  - README.md is accurate
  - All doc links work
  - Code examples are correct
  - Version numbers are consistent

### Prepare Release

1. **Determine version number**

   Based on the [versioning scheme](VERSIONING.md):
   - Bug fixes only: increment SDK_PATCH (e.g., 1.2.3 -> 1.2.4)
   - New API features: increment MTCS_MINOR, reset SDK_PATCH (e.g., 1.2.3 -> 1.3.0)
   - Breaking API changes: increment MTCS_MAJOR, reset others (e.g., 1.2.3 -> 2.0.0)

2. **Update CHANGELOG.md**

   - Move items from `[Unreleased]` to new version section
   - Add version number and release date
   - Ensure all changes are documented
   - Use the correct category (Added, Changed, Deprecated, Removed, Fixed, Security)

   Example:
   ```markdown
   ## [1.2.4] - 2025-02-15

   ### Fixed
   - Resolved batch processor retry logic edge case
   - Fixed widget geometry calculation for rotated widgets

   ### Added
   - New `WithTimeout()` session option for custom timeouts
   ```

3. **Update documentation version references**

   If any documentation references specific versions, update them:
   - README.md examples
   - Getting started guide
   - Installation instructions

4. **Create release commit**
   ```bash
   git add -A
   git commit -m "chore: prepare release v1.2.4"
   ```

### Create Release

1. **Create annotated tag**
   ```bash
   git tag -a v1.2.4 -m "Release v1.2.4"
   ```

   Tag message format:
   ```
   Release v1.2.4

   See CHANGELOG.md for details.
   ```

2. **Push tag to GitHub**
   ```bash
   git push origin v1.2.4
   git push origin main  # or master
   ```

3. **Create GitHub Release**

   Go to GitHub > Releases > Create new release

   - **Tag**: Select the tag you just pushed (v1.2.4)
   - **Title**: `v1.2.4`
   - **Description**: Use the template below

### GitHub Release Template

```markdown
## Canvus Go SDK v1.2.4

### Highlights

Brief description of the most important changes (1-3 bullet points).

### MTCS API Compatibility

- **Compatible with**: MTCS API v1.2
- **Tested against**: MTCS Server version X.Y.Z

### Installation

```bash
go get github.com/jaypaulb/Canvus-Go-API/canvus@v1.2.4
```

### Changelog

#### Fixed
- List of fixes from CHANGELOG.md

#### Added
- List of additions from CHANGELOG.md

See [CHANGELOG.md](CHANGELOG.md) for full details.

### Upgrading

Notes on upgrading from previous versions (if any breaking changes or migration steps).

### Known Issues

- List any known issues (or "None" if none)
```

### Post-Release Validation

After creating the release:

- [ ] **Verify tag on GitHub**
  - Tag appears in repository tags
  - Release appears in releases list

- [ ] **Test installation**
  ```bash
  # In a new directory
  go mod init test
  go get github.com/jaypaulb/Canvus-Go-API/canvus@v1.2.4
  ```

  Verify the correct version is installed.

- [ ] **Check pkg.go.dev**

  Wait a few minutes, then visit:
  `https://pkg.go.dev/github.com/jaypaulb/Canvus-Go-API/canvus@v1.2.4`

  Verify documentation renders correctly.

- [ ] **Update GitHub release if needed**

  Add any missing information or corrections.

## MTCS API Compatibility Notes Format

When documenting MTCS API compatibility, use this format:

```markdown
### MTCS API Compatibility

| SDK Version | MTCS API Version | Status |
|-------------|------------------|--------|
| v1.2.4 | v1.2 | Full support |
| v1.2.4 | v1.1 | Backward compatible |
| v1.2.4 | v1.3 | Partial (new features unavailable) |

#### Notes
- Feature X requires MTCS API v1.2 or later
- Deprecated: Y endpoint removed in MTCS API v1.3
```

Include in release notes:
- Primary compatible MTCS API version
- Known limitations with older/newer versions
- Any required server configuration

## Hotfix Process

For urgent fixes to a released version:

1. **Create hotfix branch** (if needed)
   ```bash
   git checkout -b hotfix/1.2.5 v1.2.4
   ```

2. **Apply fix and test**
   ```bash
   # Make changes
   go test ./canvus/... -v
   ```

3. **Follow standard release process**
   - Update CHANGELOG
   - Create tag
   - Create GitHub release

4. **Merge back to main**
   ```bash
   git checkout main
   git merge hotfix/1.2.5
   ```

## Release Cadence

The SDK follows a release-when-ready approach:

- **Patch releases**: As needed for bug fixes (typically weekly during active development)
- **Minor releases**: When new MTCS API features are implemented
- **Major releases**: When MTCS API major version changes

There is no fixed schedule; releases are made when there are meaningful changes to release.

## Troubleshooting

### pkg.go.dev not updating

If pkg.go.dev doesn't show your new release:
1. Wait 15-30 minutes for the proxy to update
2. Request an update: `https://pkg.go.dev/github.com/jaypaulb/Canvus-Go-API/canvus?tab=versions`
3. Check that the tag follows Go module conventions (v prefix)

### go get returns old version

1. Clear module cache: `go clean -modcache`
2. Specify exact version: `go get github.com/jaypaulb/Canvus-Go-API/canvus@v1.2.4`

### Tests fail in CI but pass locally

1. Verify test environment configuration
2. Check for race conditions: `go test -race ./canvus/...`
3. Ensure clean test state (no leftover test data)

## Security Releases

For security-related releases:

1. **Do not disclose details** in commit messages or CHANGELOG until release
2. **Use generic descriptions** like "Fixed security issue in authentication"
3. **After release**, update CHANGELOG with full details
4. **Consider coordinated disclosure** if the issue affects users

### Security Release Template

```markdown
## Security Release v1.2.5

### Security

- Fixed [CVE-XXXX-YYYY]: Brief description
  - Severity: High/Medium/Low
  - Affected versions: v1.2.0 - v1.2.4
  - Recommendation: Upgrade immediately

### Upgrading

All users should upgrade to v1.2.5 or later immediately.
```

## Checklist Summary

Quick reference checklist for releases:

```
Pre-release:
[ ] Tests pass
[ ] Build succeeds
[ ] go mod verify clean
[ ] Documentation reviewed

Release:
[ ] Version determined
[ ] CHANGELOG updated
[ ] Release commit created
[ ] Tag created and pushed
[ ] GitHub release published

Post-release:
[ ] Installation tested
[ ] pkg.go.dev updated
[ ] Announcement made (if needed)
```
