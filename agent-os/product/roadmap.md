# Product Roadmap

## Phase 1: Foundation & Core SDK (COMPLETED)

1. [x] **SDK Infrastructure** — Established session management with multiple authentication strategies (API key, login/password, token refresh). Implemented centralized HTTP client, request/response marshaling, and typed error handling.

2. [x] **System Management Endpoints** — Users, access tokens, groups, canvas folders, server config, license, audit logs, and server info fully implemented with comprehensive integration tests.

3. [x] **Canvas Operations** — Complete CRUD operations, move, copy, rename, trash, and permission management for canvases with full test coverage.

4. [x] **Client & Workspace Management** — All client and workspace operations including client info, workspace creation/deletion, and canvas URL launching with verification.

5. [x] **Widget Framework** — Full widget CRUD for all types (notes, anchors, images, connectors, PDFs, videos) with location/size management and spatial relationships.

6. [x] **Asset Management** — Complete asset lifecycle including uploads, image/PDF/video processing, backgrounds, mipmaps, and asset versioning.

7. [x] **Geometry Utilities** — Widget containment and overlap detection (WidgetsContainId, WidgetsTouchId) enabling spatial queries and analysis.

8. [x] **Batch Operations Framework** — Efficient batch processing with retry logic for multi-operation workflows and large-scale automations.

9. [x] **Import/Export System** — Robust round-trip import/export for all widget and asset types with case-insensitive type handling and numeric validation relaxation.

10. [x] **Filtering & Search** — Client-side filtering with wildcards, partial matches, and nested field selectors; FindWidgetsAcrossCanvases for cross-canvas queries.

11. [x] **Documentation** — Complete godoc comments, comprehensive README with examples, code samples, and integration test documentation.

---

## Phase 2: Session Management & CLI Implementation (CURRENT)

1. [ ] **Session Management Enhancements** — Implement session pooling for concurrent request handling, automatic token refresh with background renewal, session timeout/inactivity handling, and connection reuse across multiple SDK instances. `L`

2. [ ] **CLI Core Framework** — Implement root command structure, configuration management via settings.json, authentication handling, and command output formatting (JSON, table, text). `M`

3. [ ] **Canvas CLI Commands** — Implement canvus-cli canvas subcommands (create, list, get, rename, move, copy, delete, permissions, export, import) with interactive prompts where appropriate. `L`

4. [ ] **Widget CLI Commands** — Implement canvus-cli widget subcommands (create, list, get, update, delete) for all widget types with parameterized creation. `L`

5. [ ] **User & Token CLI Commands** — Implement user management commands (create, list, delete, activate) and token management (create, list, revoke, refresh) with permission-aware operations. `M`

6. [ ] **System Admin CLI Commands** — Implement admin-only commands (server-config, audit-log, license-info, user-cleanup) for server management and diagnostics. `M`

7. [ ] **Batch Operations CLI** — Implement CLI support for batch operations including bulk canvas creation, widget import from CSV, bulk asset uploads, and workflow automation. `L`

8. [ ] **Interactive Authentication** — Enhance CLI with interactive login, token creation workflow, and secure credential storage for repeated commands. `M`

---

## Phase 3: Advanced CLI & Automation (PLANNED)

1. [ ] **Workspace Management CLI** — CLI commands for workspace creation, client launching, joining, and leaving; workspace configuration management. `M`

2. [ ] **Template & Snapshot System** — Export canvas templates, create canvas from template, snapshot versioning, and rollback capabilities. `L`

3. [ ] **Event Subscription Framework** — Implement subscription support for real-time updates on canvas changes, widget updates, and user events; streaming JSON handlers. `L`

4. [ ] **Workflow Automation DSL** — Simple configuration format for defining multi-step workflows (canvas creation, bulk imports, user provisioning, cleanup). `L`

5. [ ] **Configuration Profiles** — Support multiple named Canvus server configurations, easy switching between development/staging/production environments. `M`

6. [ ] **Audit & Logging** — Enhanced logging output (structured JSON, log levels), audit trail for CLI operations, and sensitive data masking. `M`

---

## Phase 4: Performance & Ecosystem (FUTURE)

1. [ ] **Connection Pooling & Performance** — Implement HTTP connection pooling, request pipelining, and performance monitoring utilities for large-scale operations. `M`

2. [ ] **Caching Layer** — Optional client-side caching for read-heavy workloads, cache invalidation strategies, and TTL configuration. `M`

3. [ ] **Middleware Ecosystem** — Extensible middleware system for request/response hooks, custom logging, metrics collection, and distributed tracing integration. `L`

4. [ ] **SDK Documentation Site** — Comprehensive API reference, tutorial guides, best practices, migration guides, and architectural overview. `L`

5. [ ] **Go Module Stability** — Semantic versioning enforcement, backward compatibility guarantees, deprecation policies, and changelog maintenance. `M`

6. [ ] **Community Contributions** — Code contribution guidelines, development environment setup, issue templates, PR review process, and contributor recognition. `L`

7. [ ] **Integration Examples** — Reference implementations showing Canvus integration with popular Go frameworks (gin, echo, gRPC), databases (PostgreSQL, MongoDB), and cloud platforms (AWS, GCP, Azure). `L`

8. [ ] **Monitoring & Observability** — Structured logging, metrics export (Prometheus), distributed tracing support (OpenTelemetry), and health check utilities. `L`

---

## Notes

- Order items by technical dependencies and product architecture
- Each item represents an end-to-end (frontend SDK + CLI) functional and testable feature
- Effort scale:
  - XS: 1 day
  - S: 2-3 days
  - M: 1 week
  - L: 2 weeks
  - XL: 3+ weeks
- Phase 2 focuses on completing the productization of the SDK through CLI tools and session management improvements
- Phases 3 and 4 represent future enhancements based on community feedback and ecosystem growth
- All work assumes existing SDK infrastructure is stable; focus shifts to developer tooling and automation
