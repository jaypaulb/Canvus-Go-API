# Canvus Go SDK

[![Go Version](https://img.shields.io/badge/Go-1.16+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![GoDoc](https://pkg.go.dev/badge/github.com/jaypaulb/Canvus-Go-API/canvus.svg)](https://pkg.go.dev/github.com/jaypaulb/Canvus-Go-API/canvus)
[![Latest Release](https://img.shields.io/github/v/release/jaypaulb/Canvus-Go-API?include_prereleases)](https://github.com/jaypaulb/Canvus-Go-API/releases)
[![Build Status](https://github.com/jaypaulb/Canvus-Go-API/actions/workflows/release.yml/badge.svg)](https://github.com/jaypaulb/Canvus-Go-API/actions)

## Why Canvus Go SDK?

The Canvus API is powerful but complex. The Canvus Go SDK eliminates boilerplate and provides a production-ready foundation for building integrations, automation tools, and services that interact with Canvus.

**Stop wrestling with raw HTTP calls.** Get strongly-typed responses, automatic retries, proper authentication handling, and comprehensive error management out of the box.

## Features

- **Complete API Coverage** - 109+ methods covering all Canvus endpoints
- **Strongly Typed** - Full Go structs for all requests and responses
- **Production Ready** - Automatic retries, circuit breakers, context support
- **Multiple Auth Flows** - API key, username/password, token refresh
- **Batch Operations** - Concurrent bulk operations with progress tracking
- **Import/Export** - Round-trip safe widget and asset migration
- **Geometry Utilities** - Spatial queries for widget containment and overlap
- **Flexible Filtering** - Client-side filtering with wildcards and JSONPath

## Installation

```bash
go get github.com/jaypaulb/Canvus-Go-API/canvus
```

Requires Go 1.16 or later.

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "github.com/jaypaulb/Canvus-Go-API/canvus"
)

func main() {
    // Create session with API key
    session := canvus.NewSession("https://your-server/api/v1", canvus.WithAPIKey("YOUR_API_KEY"))

    // List all canvases
    canvases, err := session.ListCanvases(context.Background(), nil)
    if err != nil {
        panic(err)
    }

    for _, c := range canvases {
        fmt.Printf("%s: %s\n", c.ID, c.Name)
    }
}
```

For a complete walkthrough, see the [Getting Started Guide](docs/GETTING_STARTED.md).

## Use Cases

- **System Integrations** - Connect Canvus with your existing tools and workflows
- **Automation Scripts** - Bulk operations, scheduled tasks, CI/CD pipelines
- **Admin Tools** - User provisioning, permission management, audit reporting
- **MCP Server Foundation** - Build AI coding agents with Canvus capabilities
- **Custom Applications** - Embed Canvus functionality in your own products

## API Coverage

The SDK provides **109+ methods** organized into these categories:

| Category | Description | Methods |
|----------|-------------|---------|
| **Users** | User CRUD, access tokens, group membership | 12 |
| **Canvases** | Canvas CRUD, permissions, backgrounds | 15 |
| **Widgets** | All widget types (notes, images, PDFs, videos, etc.) | 40+ |
| **Folders** | Folder organization and permissions | 10 |
| **System** | Server config, license, audit logs | 6 |
| **Clients** | Client devices and workspaces | 12 |
| **Batch** | Bulk operations with retry logic | 6 |
| **Import/Export** | Widget and asset migration | 2 |

See the complete [API Reference](docs/API_REFERENCE.md) for all methods.

## Documentation

- **[Getting Started](docs/GETTING_STARTED.md)** - Installation and first API call
- **[Best Practices](docs/BEST_PRACTICES.md)** - Error handling, auth patterns, concurrency
- **[Troubleshooting](docs/TROUBLESHOOTING.md)** - Common issues and solutions
- **[API Reference](docs/API_REFERENCE.md)** - Complete method reference
- **[Examples Index](docs/EXAMPLES.md)** - Runnable code examples
- **[Compatibility](docs/COMPATIBILITY.md)** - MTCS API version support

## Examples

The [/examples](examples/) directory contains runnable examples for common tasks:

- [Getting Started](examples/getting_started/) - Basic setup and first call
- [Authentication](examples/authentication/) - API key and login flows
- [Canvas Operations](examples/canvases/) - CRUD and permissions
- [Widget Management](examples/widgets/) - Create, search, spatial queries
- [Batch Processing](examples/batch/) - Bulk operations with retries
- [Import/Export](examples/import_export/) - Round-trip widget migration
- [Error Handling](examples/error_handling/) - Recovery patterns
- [Context Usage](examples/context/) - Cancellation and timeouts

## Key Differentiators

| Feature | Canvus Go SDK | Raw HTTP |
|---------|---------------|----------|
| Type Safety | Full Go structs | Manual JSON parsing |
| Error Handling | Typed errors with codes | Raw status codes |
| Retries | Automatic exponential backoff | Manual implementation |
| Auth | Built-in with refresh | Manual header management |
| Batch Ops | Concurrent with progress | Sequential only |
| Testing | Mock-friendly interfaces | HTTP stubs required |

## Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

### Reporting Issues

- Search existing issues before creating new ones
- Include Go version, SDK version, and error messages
- Provide minimal reproduction steps

### Pull Requests

- Fork the repository and create a feature branch
- Follow [Conventional Commits](https://www.conventionalcommits.org/) for commit messages
- Add tests for new functionality
- Update documentation as needed
- Run `go test ./canvus/...` before submitting

## License

MIT License. See [LICENSE](LICENSE) file.

## Related Projects

- **[canvus-cli](https://github.com/jaypaulb/canvus-cli)** - Command-line interface for Canvus (uses this SDK)
- **[OpenAPI Specification](openapi.yaml)** - Machine-readable API specification
