# Release Template

Use this template when creating a new GitHub release.

---

## Canvus Go SDK vX.Y.Z

### Summary

Brief description of this release (2-3 sentences highlighting the most important changes).

### Highlights

- Highlight 1: Most important change
- Highlight 2: Second most important change
- Highlight 3: Third most important change (if applicable)

### MTCS API Compatibility

- **Compatible with**: MTCS API vX.Y
- **Tested against**: MTCS Server version X.Y.Z
- **Minimum Go version**: 1.24+

| SDK Version | MTCS API Version | Status |
|-------------|------------------|--------|
| vX.Y.Z | vX.Y | Full support |
| vX.Y.Z | vX.(Y-1) | Backward compatible |

### Installation

Install or update the SDK using Go modules:

```bash
go get github.com/jaypaulb/Canvus-Go-API/canvus@vX.Y.Z
```

Or add to your `go.mod`:

```go
require github.com/jaypaulb/Canvus-Go-API vX.Y.Z
```

Then run:

```bash
go mod tidy
```

### What's Changed

#### Added
- List of new features

#### Changed
- List of changes to existing features

#### Fixed
- List of bug fixes

#### Deprecated
- List of deprecated features (if any)

#### Security
- Security-related changes (if any)

See [CHANGELOG.md](https://github.com/jaypaulb/Canvus-Go-API/blob/master/CHANGELOG.md) for full details.

### Upgrading

**From vA.B.C to vX.Y.Z:**

Notes on upgrading from previous versions. Include:
- Breaking changes (if any)
- Migration steps (if needed)
- Deprecated features to be aware of

If no special steps are needed:
> No breaking changes in this release. Simply update your dependency version.

### Known Issues

- List any known issues or limitations
- Include workarounds where available

If none:
> No known issues in this release.

### Documentation

- [Getting Started Guide](https://github.com/jaypaulb/Canvus-Go-API/blob/master/docs/GETTING_STARTED.md)
- [Best Practices](https://github.com/jaypaulb/Canvus-Go-API/blob/master/docs/BEST_PRACTICES.md)
- [API Reference](https://pkg.go.dev/github.com/jaypaulb/Canvus-Go-API/canvus)
- [Examples](https://github.com/jaypaulb/Canvus-Go-API/tree/master/examples)

### Acknowledgments

Thanks to everyone who contributed to this release!

---

## Release Checklist (for maintainers)

Before publishing this release, ensure:

- [ ] All tests pass in CI
- [ ] CHANGELOG.md is updated with all changes
- [ ] Documentation is updated as needed
- [ ] Version numbers are consistent across docs
- [ ] Tag follows semantic versioning (vMAJOR.MINOR.PATCH)
- [ ] Release notes accurately describe changes
- [ ] MTCS API compatibility is documented

After publishing:

- [ ] Verify installation works: `go get github.com/jaypaulb/Canvus-Go-API/canvus@vX.Y.Z`
- [ ] Check pkg.go.dev is updated
- [ ] Announce release (if appropriate)
