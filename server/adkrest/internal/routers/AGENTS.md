# AGENTS.md - server/adkrest/internal/routers

> HTTP router layer for the ADK REST API server. Implements Chi-based routing with middleware for sessions, artifacts, runtime operations, and debug endpoints. Each router file handles a specific domain (sessions, artifacts, runtime, debug) and registers routes with the main router.

## Tech Stack

- **go-chi/chi** v5 - HTTP router and middleware framework for RESTful API routing
- **net/http** stdlib - Standard HTTP server and handler interfaces

## Key Files

| File | Purpose | Read When |
|------|---------|------------|
| `server/adkrest/internal/routers/routers.go` (reference) | Main router registration and middleware setup. Entry point for understanding routing architecture. | Adding new routes or understanding overall API structure |
| `server/adkrest/internal/routers/sessions.go` (reference) | Session management endpoints (create, list, get, delete sessions) | Working with session lifecycle or session-related APIs |
| `server/adkrest/internal/routers/artifacts.go` (reference) | Artifact handling endpoints (upload, download, list artifacts) | Implementing artifact storage or retrieval features |
| `server/adkrest/internal/routers/runtime.go` (reference) | Runtime operation endpoints (execute, stream, health checks) | Working with agent execution or streaming responses |
| `server/adkrest/internal/routers/debug.go` (reference) | Debug and introspection endpoints (metrics, profiling, state inspection) | Adding observability or debugging capabilities |

## Architecture

```
Request Flow:

[HTTP Request] → [routers.go: Main Router]
                      |
                      ├─→ /sessions/* → [sessions.go]
                      ├─→ /artifacts/* → [artifacts.go]
                      ├─→ /runtime/* → [runtime.go]
                      └─→ /debug/* → [debug.go]

Pattern:
1. routers.go creates chi.Router and applies global middleware
2. Each domain router (sessions, artifacts, etc.) registers sub-routes
3. Handlers delegate to service layer (not in this package)
4. Middleware handles cross-cutting concerns (auth, logging, CORS)

See server/adkrest/internal/routers/routers.go for router initialization
```

## Patterns

### Domain-Based Router Separation

Each domain (sessions, artifacts, runtime, debug) has its own router file that registers routes under a specific path prefix. Keeps routing logic organized by business domain.

See `server/adkrest/internal/routers/sessions.go` for reference.

### Chi Router Composition

Uses chi.Router interface for composable routing. Main router mounts sub-routers using r.Mount() pattern. Allows independent testing and organization of route groups.

See `server/adkrest/internal/routers/routers.go` for reference.

### Handler Dependency Injection

Router functions accept dependencies (services, configs) as parameters and return configured chi.Router. Enables testability and loose coupling.

See `server/adkrest/internal/routers/artifacts.go` for reference.

### RESTful Resource Routing

Follows REST conventions: GET for retrieval, POST for creation, DELETE for removal. Uses path parameters for resource IDs (e.g., /sessions/{sessionID}).

See `server/adkrest/internal/routers/sessions.go` for reference.

### Middleware Chain Application

Applies middleware at router level using r.Use(). Common middleware includes authentication, logging, CORS, and request validation.

See `server/adkrest/internal/routers/routers.go` for reference.

## Rules

### Always (Do These)

| Rule | Rationale |
|------|------------|
| CRITICAL: Register all routes through domain-specific router functions, never directly in main | Maintains separation of concerns and testability. Each domain router is independently testable and composable. |
| HIGH: Use chi.URLParam() to extract path parameters, not manual string parsing | Chi provides type-safe parameter extraction. Manual parsing is error-prone and bypasses router validation. |
| HIGH: Return chi.Router from router setup functions for composability | Enables mounting sub-routers and independent testing. Follows chi's composition pattern. |
| MEDIUM: Apply middleware at the appropriate scope (global vs route-specific) | Global middleware in routers.go, route-specific in domain routers. Prevents unnecessary middleware execution. |
| MEDIUM: Use consistent HTTP status codes (200 OK, 201 Created, 204 No Content, 404 Not Found) | RESTful API consistency. Clients depend on standard status codes for error handling. |
| HIGH: Delegate business logic to service layer, keep handlers thin | Routers handle HTTP concerns only. Business logic belongs in services for testability and reusability. |

### Never (Don't Do These)

| Rule | Severity | Rationale |
|------|----------|------------|
| CRITICAL: Never implement business logic in router handlers | CRITICAL | Violates separation of concerns. Business logic must be in service layer for testability, reusability, and maintainability. |
| CRITICAL: Never bypass chi's routing system with manual URL parsing | CRITICAL | Chi provides validated routing. Manual parsing bypasses security checks and breaks middleware chain. |
| HIGH: Never register routes without considering middleware order | HIGH | Middleware executes in registration order. Authentication must precede authorization, logging should wrap everything. |
| HIGH: Never hardcode paths or URLs in handlers | HIGH | Use chi.URLParam() and path constants. Hardcoded paths break when routes change and are untestable. |
| MEDIUM: Never mix HTTP concerns with domain logic | MEDIUM | Routers handle HTTP (request/response), services handle domain logic. Mixing makes testing and reuse difficult. |
| HIGH: Never return raw errors to clients without sanitization | HIGH | Raw errors may expose internal implementation details or security information. Use error handling middleware. |

### Ask First

- **Adding new top-level route prefix (e.g., /v2/, /admin/)** - Affects API versioning strategy and client contracts. Requires coordination with API documentation and client teams.
- **Changing existing route paths or HTTP methods** - Breaking change for existing clients. Requires deprecation strategy and version management.
- **Adding global middleware that affects all routes** - Performance and security implications. May break existing clients or introduce latency.
- **Implementing custom authentication or authorization in routers** - Security-critical. Should use established middleware patterns and be reviewed by security team.
- **Adding streaming or WebSocket endpoints** - Different lifecycle and resource management than REST. Requires infrastructure and monitoring considerations.

## Commands

### test-routers

Run all router tests

```bash
go test ./server/adkrest/internal/routers/...
```

### test-router-coverage

Run router tests with coverage report

```bash
go test -cover ./server/adkrest/internal/routers/...
```

### lint-routers

Lint router code for common issues

```bash
golangci-lint run ./server/adkrest/internal/routers/...
```

## Testing

Test routers using httptest.NewRecorder and httptest.NewRequest. Mock service layer dependencies. Verify HTTP status codes, response bodies, and header handling. Test middleware application and error cases.

```bash
go test ./server/adkrest/internal/routers/... -v
go test ./server/adkrest/internal/routers/... -race
```

Test directory: `server/adkrest/internal/routers/`

## Common Tasks

### Add new REST endpoint

1. 1. Identify domain (sessions, artifacts, runtime, debug) or create new router file
2. 2. Add route registration in appropriate router function (e.g., NewSessionsRouter)
3. 3. Define handler function that accepts http.ResponseWriter and *http.Request
4. 4. Extract parameters using chi.URLParam() for path params, r.URL.Query() for query params
5. 5. Delegate to service layer for business logic
6. 6. Write response using json.NewEncoder(w).Encode() or w.Write()
7. 7. Set appropriate HTTP status code with w.WriteHeader()
8. 8. Add tests using httptest package
9. 9. Update API documentation

### Add middleware to specific routes

1. 1. Determine scope: global (routers.go) or domain-specific (sessions.go, etc.)
2. 2. Create middleware function: func(next http.Handler) http.Handler
3. 3. For global: add to main router in routers.go using r.Use(middleware)
4. 4. For domain: add to domain router using r.Use(middleware) or r.With(middleware).Get(...)
5. 5. Ensure middleware calls next.ServeHTTP(w, r) to continue chain
6. 6. Test middleware execution order and behavior

### Debug routing issues

1. 1. Check route registration in appropriate router file (sessions.go, artifacts.go, etc.)
2. 2. Verify path pattern matches request (check for trailing slashes, parameter names)
3. 3. Confirm HTTP method matches (GET, POST, DELETE)
4. 4. Check middleware chain isn't blocking request
5. 5. Use chi's built-in debugging: chi.Walk(router, func(method, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler))
6. 6. Add logging middleware to trace request flow
7. 7. Verify routers.go mounts sub-router at correct path

### Refactor existing endpoint

1. 1. Locate handler in appropriate router file (sessions.go, artifacts.go, etc.)
2. 2. Extract business logic to service layer if mixed in handler
3. 3. Update handler to call service methods
4. 4. Ensure error handling uses consistent pattern
5. 5. Update tests to mock service layer
6. 6. Verify HTTP status codes and response format unchanged (or document breaking changes)
7. 7. Run integration tests to verify end-to-end behavior

