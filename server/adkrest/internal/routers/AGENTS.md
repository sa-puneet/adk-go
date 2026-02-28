# AGENTS.md - server/adkrest/internal/routers

> HTTP router layer for the ADK REST API server. Implements Chi-based routing with middleware for sessions, artifacts, runtime operations, and debug endpoints. Each router file handles a specific domain (sessions, artifacts, runtime, debug) and registers routes with appropriate handlers.

## Tech Stack

- **go-chi/chi** v5 - HTTP router and middleware framework for building REST APIs with composable route handlers
- **go** 1.x - Primary language for router implementation with standard net/http integration

## Key Files

| File | Purpose | Read When |
|------|---------|------------|
| `server/adkrest/internal/routers/routers.go` (reference) | Main router setup and initialization. Defines Router interface and creates the root Chi router with all sub-routers mounted | Understanding overall routing architecture or adding new router domains |
| `server/adkrest/internal/routers/sessions.go` | Session management routes (create, list, get, delete sessions). Handles agent session lifecycle endpoints | Working with session-related endpoints or understanding session API patterns |
| `server/adkrest/internal/routers/artifacts.go` | Artifact storage and retrieval routes. Manages file/data artifacts associated with sessions | Implementing artifact upload/download or storage features |
| `server/adkrest/internal/routers/runtime.go` | Runtime operation routes for agent execution, tool calls, and workflow management | Adding runtime execution endpoints or modifying agent execution flow |
| `server/adkrest/internal/routers/debug.go` | Debug and diagnostic endpoints for development and troubleshooting | Adding debugging capabilities or diagnostic endpoints |

## Architecture

```
Request Flow:
1. HTTP Request → routers.go (root Chi router)
2. Route matching → domain-specific router (sessions/artifacts/runtime/debug)
3. Middleware chain execution (auth, logging, context)
4. Handler function invocation
5. Handler calls service layer (not in this package)
6. Response serialization and return

Router Organization:
  routers.go (root)
    ├── /sessions/* → sessions.go
    ├── /artifacts/* → artifacts.go
    ├── /runtime/* → runtime.go
    └── /debug/* → debug.go

Each domain router:
- Implements route registration function
- Returns chi.Router for mounting
- Accepts handler dependencies via constructor
- Uses Chi's r.Get/Post/Delete methods
```

## Patterns

### Domain-Specific Router Files

Each router file handles one domain (sessions, artifacts, runtime, debug). Keeps routing logic organized and testable. Each file exports a function that returns chi.Router

See `server/adkrest/internal/routers/sessions.go` for reference.

### Handler Dependency Injection

Router functions accept handler interfaces/structs as parameters rather than creating dependencies internally. Enables testing and loose coupling

See `server/adkrest/internal/routers/routers.go` for reference.

### Chi Sub-Router Mounting

Use chi.NewRouter() for each domain, then mount to main router with r.Mount(). Allows independent middleware chains per domain

See `server/adkrest/internal/routers/routers.go` for reference.

### RESTful Route Naming

Follow REST conventions: GET for retrieval, POST for creation, DELETE for removal. Use URL parameters for resource IDs (e.g., /sessions/{id})

See `server/adkrest/internal/routers/sessions.go` for reference.

## Rules

### Always (Do These)

| Rule | Rationale |
|------|------------|
| CRITICAL: Each domain router must return chi.Router, not register routes directly on a passed router | Maintains clean separation of concerns and allows independent middleware configuration per domain. Prevents tight coupling between router initialization and route registration |
| HIGH: Use Chi's URL parameter syntax {param} for resource identifiers in routes | Chi's built-in parameter extraction (chi.URLParam) provides type-safe, consistent access to route parameters across all handlers |
| HIGH: Pass handler dependencies as parameters to router constructor functions, never create them inside router files | Router layer should only handle routing logic, not dependency creation. Enables testing with mock handlers and maintains single responsibility |
| MEDIUM: Group related routes under the same domain router file (max 10-15 routes per file) | Keeps files focused and maintainable. If a domain grows beyond 15 routes, consider splitting into sub-domains |
| MEDIUM: Use HTTP method-specific Chi functions (r.Get, r.Post, r.Delete) rather than r.Method() | More readable and idiomatic Chi usage. Makes HTTP method immediately visible in code |

### Never (Don't Do These)

| Rule | Severity | Rationale |
|------|----------|------------|
| CRITICAL: Never implement business logic in router files | CRITICAL | Routers should only handle HTTP routing and parameter extraction. Business logic belongs in handlers/services. Violating this creates untestable, tightly coupled code |
| CRITICAL: Never create database connections, service clients, or external dependencies inside router functions | CRITICAL | Router initialization should be pure routing logic. Dependencies must be injected to enable testing and proper lifecycle management |
| HIGH: Never use global variables for handler storage or router state | HIGH | Prevents testing, creates hidden dependencies, and makes concurrent testing impossible. Always pass dependencies explicitly |
| HIGH: Never register the same route path with multiple HTTP methods in different files | HIGH | Creates confusion about which handler handles which method. Keep all methods for a path in the same domain router file |
| MEDIUM: Never hardcode URL prefixes inside domain router files | MEDIUM | The mounting point is determined by routers.go. Domain routers should define relative paths only |

### Ask First

- **Adding a new domain router file (e.g., users.go, config.go)** - Ensure the domain is distinct enough to warrant a separate file and doesn't overlap with existing routers. Discuss mounting point and URL structure
- **Changing the URL structure of existing routes (e.g., /sessions → /agent-sessions)** - Breaking change for API clients. Requires versioning strategy discussion and backward compatibility plan
- **Adding middleware that affects all routes (auth, rate limiting, etc.)** - Global middleware impacts all endpoints. Discuss performance implications, error handling, and whether it should be selective per domain
- **Implementing WebSocket or SSE endpoints in router files** - Long-lived connections require different patterns than REST. Discuss connection management, cleanup, and whether Chi is appropriate

## Commands

### test-routers

Run all router tests to verify route registration and handler wiring

```bash
go test ./server/adkrest/internal/routers/...
```

### lint-routers

Check router code for style issues and potential bugs

```bash
golangci-lint run ./server/adkrest/internal/routers/
```

### list-routes

Print all registered routes (if implemented) to verify routing table

```bash
go run ./server/adkrest/cmd/server --print-routes
```

## Testing

Router testing focuses on route registration and parameter extraction, not handler logic. Use httptest.NewRecorder and httptest.NewRequest to verify routes are registered correctly. Mock handler dependencies to isolate routing logic. Test that correct handlers are called for each route and HTTP method

```bash
go test ./server/adkrest/internal/routers/... -v
go test ./server/adkrest/internal/routers/... -cover
```

Test directory: `server/adkrest/internal/routers/`

## Common Tasks

### Add a new REST endpoint to existing domain

1. 1. Open the appropriate domain router file (e.g., sessions.go for session-related endpoint)
2. 2. Add route using r.Get/Post/Delete with path and handler: r.Get("/path/{id}", handler.MethodName)
3. 3. Ensure handler method exists in the handler interface/struct passed to router constructor
4. 4. Add test in corresponding _test.go file to verify route registration
5. 5. Update API documentation if maintained separately

### Create a new domain router

1. 1. Create new file: server/adkrest/internal/routers/{domain}.go
2. 2. Define constructor function: func New{Domain}Router(handlers {Domain}Handlers) chi.Router
3. 3. Inside constructor: r := chi.NewRouter(), then register routes with r.Get/Post/Delete
4. 4. Return the chi.Router from constructor
5. 5. In routers.go, import new router and mount: r.Mount("/{domain}", New{Domain}Router(handlers))
6. 6. Create corresponding test file with route registration tests

### Add middleware to specific domain

1. 1. Open domain router file (e.g., sessions.go)
2. 2. After chi.NewRouter(), add middleware: r.Use(middlewareFunc)
3. 3. Middleware only affects routes in that domain router
4. 4. For global middleware, add in routers.go before mounting sub-routers
5. 5. Test middleware execution with httptest

### Debug route not matching

1. 1. Verify route is registered in correct domain router file
2. 2. Check HTTP method matches (GET vs POST vs DELETE)
3. 3. Verify URL parameter syntax uses {param} not :param
4. 4. Check mount point in routers.go matches expected URL prefix
5. 5. Use Chi's Walk function or --print-routes to list all registered routes
6. 6. Test with curl or httptest to isolate issue

