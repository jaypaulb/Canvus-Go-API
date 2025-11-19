# Versioning

This document explains the versioning scheme used by the Canvus Go SDK and how version numbers relate to the MT Canvus Server (MTCS) API.

## Version Scheme

The Canvus Go SDK uses a three-component version format:

```
MTCS_MAJOR.MTCS_MINOR.SDK_PATCH
```

For example: `1.2.3`

### Component Definitions

| Component | Description | Example |
|-----------|-------------|---------|
| **MTCS_MAJOR** | Major version of the MT Canvus Server API | `1` in `1.2.3` |
| **MTCS_MINOR** | Minor version of the MT Canvus Server API | `2` in `1.2.3` |
| **SDK_PATCH** | SDK-specific patch/release number | `3` in `1.2.3` |

### What Each Component Means

**MTCS_MAJOR**: Indicates compatibility with a major version of the Canvus API. Major version changes typically indicate breaking changes in the underlying API that require significant SDK updates.

**MTCS_MINOR**: Indicates compatibility with a minor version of the Canvus API. Minor version changes typically add new API endpoints or capabilities that the SDK can leverage.

**SDK_PATCH**: Represents SDK-only changes such as bug fixes, documentation improvements, performance optimizations, or developer experience enhancements that don't require MTCS API changes.

## Version Bump Examples

### SDK Patch Bump (e.g., 1.2.3 -> 1.2.4)

Increment SDK_PATCH for:
- Bug fixes in the SDK
- Documentation updates
- Performance improvements
- Code refactoring
- New SDK utilities that don't require API changes
- Test improvements

```
1.2.3 -> 1.2.4  # Bug fix in batch processing
1.2.4 -> 1.2.5  # Documentation improvements
1.2.5 -> 1.2.6  # Performance optimization
```

### Minor Version Bump (e.g., 1.2.x -> 1.3.0)

Increment MTCS_MINOR (and reset SDK_PATCH to 0) when:
- The SDK adds support for new MTCS API endpoints
- The MTCS API adds new features that the SDK implements

```
1.2.6 -> 1.3.0  # SDK adds support for new MTCS 1.3 API endpoints
```

### Major Version Bump (e.g., 1.x.x -> 2.0.0)

Increment MTCS_MAJOR (and reset both MTCS_MINOR and SDK_PATCH to 0) when:
- The SDK updates to support a new major MTCS API version
- Breaking changes in the underlying API require SDK restructuring

```
1.3.2 -> 2.0.0  # SDK updated for MTCS API v2.0
```

## Resetting SDK_PATCH to 0

The SDK_PATCH component is reset to 0 whenever:

1. **MTCS_MINOR is incremented**: When new API functionality is added
   ```
   1.2.5 -> 1.3.0  (not 1.3.5)
   ```

2. **MTCS_MAJOR is incremented**: When there are breaking API changes
   ```
   1.3.2 -> 2.0.0  (not 2.3.2)
   ```

This reset ensures that the version number accurately reflects the compatibility relationship between the SDK and the MTCS API.

## Pre-Release Versions

Before v1.0.0, the SDK is considered pre-release:

- **v0.x.y**: Pre-release versions indicate the SDK is under active development
- The API may change between minor versions during pre-release
- Production use is supported but users should expect more frequent updates

Starting version: **v0.1.0**

Transition to v1.0.0 occurs when:
- Core documentation is complete
- All major features are implemented and stable
- CLI has been separated to its own repository
- The SDK has been validated in production use cases

## Compatibility Promises

### Within the Same MTCS Version

SDK versions with the same MTCS_MAJOR.MTCS_MINOR are guaranteed to be compatible:

- `1.2.0`, `1.2.1`, `1.2.2` all work with MTCS API v1.2
- Bug fixes and improvements don't affect compatibility
- You can safely upgrade SDK_PATCH versions

### Backward Compatibility

The SDK aims to maintain backward compatibility:

- SDK v1.2.x should work with MTCS API v1.1 (older minor versions)
- SDK v1.2.x may work with MTCS API v1.3 (newer minor versions, with reduced functionality)
- SDK v1.x.x is not guaranteed to work with MTCS API v2.x (different major versions)

### Forward Compatibility

The SDK provides limited forward compatibility:

- New API features not yet in the SDK can be accessed via raw HTTP methods
- The Session type exposes methods for custom API calls

## Determining Compatible Versions

To determine which SDK version to use:

1. Identify your MTCS API version (check your server's `/api/v1/system/info`)
2. Choose an SDK version where MTCS_MAJOR.MTCS_MINOR matches your API version
3. Use the highest SDK_PATCH for that MTCS version to get the latest bug fixes

**Example**:
- Your server runs MTCS API v1.2
- Available SDK versions: 1.1.0, 1.1.1, 1.2.0, 1.2.1, 1.2.2, 1.3.0
- Choose SDK v1.2.2 (matches your API version, highest patch)

## Release Tags

All releases are tagged in Git following Go module versioning conventions:

```
v0.1.0  # First pre-release
v0.2.0  # Second pre-release
v1.0.0  # First stable release
v1.0.1  # Patch release
v1.1.0  # Minor release
```

Tags are annotated and include release notes.

## Checking Version

You can check the SDK version programmatically:

```go
import "github.com/jaypaulb/Canvus-Go-API/canvus"

// Version constants will be available in future releases
// For now, reference the installed module version in go.mod
```

## Support Policy

Currently, the SDK supports:

- **Latest version only**: Only the most recent release receives updates
- **No parallel version maintenance**: Older versions are not maintained with patches

This policy will be reconsidered after v1.0.0 when the SDK reaches stable status.

## Questions

For questions about versioning or compatibility, please:
1. Check the [COMPATIBILITY.md](COMPATIBILITY.md) document
2. Open an issue on GitHub
3. Review the CHANGELOG for version-specific changes
