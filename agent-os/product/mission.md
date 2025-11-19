# Product Mission

## Pitch

Canvus-Go-API is an idiomatic Go SDK and CLI toolkit that enables Go developers to programmatically interact with Canvus servers for canvas automation, asset management, system administration, and custom integrations.

## Users

### Primary Customers
- **Go Backend Developers**: Building integrations with Canvus systems, automating canvas operations, and developing custom workflows
- **Systems Administrators**: Managing Canvus server instances, user provisioning, auditing, and configuration through programmatic interfaces
- **DevOps Engineers**: Automating deployment pipelines, server setup, and infrastructure management using Canvus APIs

### User Personas

**Integration Developer** (Mid-Level Backend Engineer)
- **Role:** Backend developer responsible for integrating Canvus into larger systems
- **Context:** Building microservices or applications that interact with collaborative canvas workspaces
- **Pain Points:** Existing solutions require manual API interaction or browser automation; lack of idiomatic Go support; difficulty managing authentication across multiple requests; complex error handling
- **Goals:** Quickly build reliable integrations with minimal boilerplate; leverage Go's concurrency and performance; reuse patterns across projects

**System Administrator** (DevOps/SRE)
- **Role:** Managing Canvus server infrastructure and user lifecycle
- **Context:** Enterprise deployment requiring automation of user creation, workspace provisioning, and audit logging
- **Pain Points:** Manual operations are error-prone and time-consuming; need scriptable, reproducible workflows; limited visibility into system health and audit trails
- **Goals:** Automate common administrative tasks; reliably manage users, permissions, and server configuration; integrate with existing infrastructure-as-code

**Enterprise Architect** (Platform Engineering)
- **Role:** Designing system integrations and collaborative workflows at scale
- **Context:** Large organization needing to connect Canvus with custom applications and workflows
- **Pain Points:** Building SDKs or client libraries requires deep API knowledge; need to support multiple languages and platforms
- **Goals:** Provide a production-ready, battle-tested library for internal teams; reduce development time for integrations; ensure consistency across deployments

## The Problem

### Lack of Idiomatic Go Support for Canvus Integration

Go developers currently lack a comprehensive, idiomatic SDK for the Canvus API. Existing solutions require manual HTTP handling, inconsistent error management, and boilerplate authentication logic. This friction increases development time, introduces bugs, and makes it difficult to build reliable, maintainable integrations. Additionally, administrative tasks (user management, system configuration) lack automation tooling, forcing manual operations or custom scripting.

**Our Solution:** A production-ready Go SDK with full API coverage, strong typing, context support, robust error handling, and built-in utilities for common workflows. The accompanying CLI toolkit enables administrators to script and automate server operations without writing Go code.

## Differentiators

### Comprehensive API Coverage with Idiomatic Go Design

Unlike generic REST clients or incomplete third-party libraries, the Canvus-Go-API provides 109 methods covering all Canvus endpoints (system management, canvas, workspace, widget, and asset operations) designed following Go community best practices. This results in code that is intuitive, performant, and maintainable for Go developers.

### Advanced Features Out of the Box

The SDK includes production-ready features beyond basic API wrappers: geometry utilities for spatial operations, batch processing framework, import/export with round-trip fidelity, comprehensive filtering and search, and multiple authentication strategies. Developers can build sophisticated integrations without reinventing these capabilities.

### Strong Type Safety and Error Handling

Request/response models are strongly typed Go structs with semantic validation. Errors are typed and contextual, enabling precise error handling and recovery. Context support throughout ensures proper request cancellation and timeout handling.

### Integration Test Coverage

All endpoints are covered by comprehensive integration tests against a live server, ensuring the SDK works correctly against real Canvus instances. This gives developers confidence in reliability and compatibility.

### Developer-Friendly Abstractions

The SDK includes utilities like FindWidgetsAcrossCanvases, RoundTripper middleware for custom request handling, geometry functions, and batch operations. These reduce boilerplate and enable common patterns without custom implementation.

## Key Features

### Core SDK Features

- **Full API Coverage**: 109 methods covering all Canvus endpoints (users, tokens, groups, folders, canvases, workspaces, clients, widgets, assets, server config, audit logs)
- **Strong Typing**: Idiomatic Go structs for all request/response models with proper validation and error handling
- **Authentication Methods**: API key, login/password, and token refresh authentication with flexible session management
- **Context Support**: Native context.Context support throughout for cancellation, timeouts, and request tracing
- **Centralized Error Handling**: Typed errors with detailed status codes, messages, and structured responses for precise error handling

### Advanced Features

- **Geometry Utilities**: Widget containment and overlap detection (WidgetsContainId, WidgetsTouchId) for spatial operations
- **Batch Operations**: Framework for efficiently processing multiple operations in batches with automatic retry logic
- **Import/Export**: Robust round-trip import/export for all widget and asset types with case-insensitive widget type handling
- **Client-Side Filtering**: Flexible filtering abstraction with support for wildcards, partial matches, and nested field selectors
- **Widget Search Across Canvases**: FindWidgetsAcrossCanvases utility for querying widgets across multiple canvases
- **Pagination and Streaming**: Helper functions for handling paginated responses and subscription-based streaming
- **RoundTripper Middleware**: Custom http.RoundTripper support for intercepting, logging, and modifying requests
- **Response Validation**: Centralized response validation and retry logic for transient failures

### CLI Tools

- **Canvas Management**: Create, list, rename, move, copy, and delete canvases with CLI commands
- **Widget Operations**: Create and manage widgets programmatically through CLI
- **User Management**: Provision users, manage tokens, and handle group assignments
- **System Administration**: Configure server settings, manage audit logs, and monitor system health
- **Authentication Management**: Interactive login, token creation, and session management
- **Batch Operations**: Automate multi-step workflows and bulk operations
- **Import/Export**: CLI tools for exporting canvases and importing into other workspaces

## Product Goals

1. **Reduce Development Friction**: Enable Go developers to build Canvus integrations quickly without boilerplate or custom utilities
2. **Enable System Automation**: Provide administrators with scriptable tools to manage Canvus servers at scale
3. **Ensure Production Reliability**: Deliver a battle-tested SDK with comprehensive testing, error handling, and documentation
4. **Foster Ecosystem Growth**: Support a vibrant community of developers and integrations built on a solid foundation
