# AGENTS.md - cmd/launcher/web/api

> This package provides the HTTP API layer for the launcher web interface, handling REST endpoints for agent management and control. It serves as the bridge between the web UI and the underlying launcher service, exposing operations like starting, stopping, and querying agent status.

## Tech Stack

- **Go net/http** stdlib - HTTP server and routing for REST API endpoints
- **encoding/json** stdlib - JSON serialization for API request/response payloads
- **context** stdlib - Request lifecycle management and cancellation propagation

## Key Files

| File | Purpose | Read When |
|------|---------|------------|
| `cmd/launcher/web/api/api.go` (reference) | Main API handler implementation with route definitions and HTTP endpoint logic | Understanding API contract, adding new endpoints, or debugging HTTP layer issues |

## Architecture

```
HTTP Request → API Handler (api.go) → Launcher Service → Agent Runtime
                                    ↓
                              JSON Response

Flow:
1. Web UI sends HTTP request to /api/* endpoint
2. api.go routes to appropriate handler function
3. Handler validates request, extracts parameters
4. Calls underlying launcher service methods
5. Formats response as JSON and returns to client
6. Error conditions return appropriate HTTP status codes
```

## Patterns

### HTTP Handler Pattern

Standard Go http.HandlerFunc pattern with explicit error handling and status code management. Each endpoint is a separate function that writes JSON responses.

See `cmd/launcher/web/api/api.go` for reference.

### JSON Response Envelope

Consistent response structure with status, data, and error fields. All API responses follow the same envelope format for predictable client parsing.

See `cmd/launcher/web/api/api.go` for reference.

### Context Propagation

Request context flows from HTTP handler through service layer to enable cancellation and timeout handling across the call stack.

See `cmd/launcher/web/api/api.go` for reference.

## Rules

### Always (Do These)

| Rule | Rationale |
|------|------------|
| Return proper HTTP status codes: 200 for success, 400 for client errors, 500 for server errors | REST API clients depend on status codes for error handling logic. Incorrect codes break client error detection and retry logic. |
| Set Content-Type: application/json header for all JSON responses | Clients need proper content type to parse responses correctly. Missing headers cause parsing failures in strict HTTP clients. |
| Validate all input parameters before processing requests | Invalid input can crash the service or cause undefined behavior. Early validation provides clear error messages to API consumers. |
| Use request context for all downstream calls to enable cancellation | Long-running operations must respect client disconnection. Context cancellation prevents resource leaks when clients abort requests. |
| Log all API errors with sufficient context (endpoint, parameters, error details) | API layer is the entry point for debugging production issues. Comprehensive logging enables root cause analysis without reproducing issues. |

### Never (Don't Do These)

| Rule | Severity | Rationale |
|------|----------|------------|
| Never expose internal error details or stack traces in API responses | CRITICAL | Internal errors leak implementation details and security information. API responses should return user-friendly messages while logging full details server-side. |
| Never block indefinitely without timeout in HTTP handlers | HIGH | Blocking handlers exhaust server resources and cause cascading failures. All operations must have reasonable timeouts to prevent resource exhaustion. |
| Never mutate shared state without proper synchronization | CRITICAL | HTTP handlers run concurrently. Race conditions cause data corruption and unpredictable behavior. Use proper locking or pass requests to synchronized service layer. |
| Never ignore errors from json.Encode/Decode operations | HIGH | JSON marshaling errors indicate malformed data or type mismatches. Ignoring these errors sends corrupt responses or processes invalid input. |
| Never use panic() for error handling in HTTP handlers | HIGH | Panics crash the entire server process. HTTP handlers must catch all errors and return appropriate status codes instead of panicking. |

### Ask First

- **Adding new API endpoints that modify agent state** - State-changing operations need careful design for idempotency, error recovery, and consistency. Coordinate with service layer design to ensure proper transaction boundaries.
- **Changing existing API response formats or field names** - Breaking changes impact all API clients including web UI. Requires versioning strategy or backward compatibility plan to avoid breaking existing integrations.
- **Adding authentication or authorization middleware** - Security changes affect all endpoints and require consistent implementation. Must align with overall security architecture and not introduce vulnerabilities.
- **Implementing long-polling or streaming endpoints** - Different concurrency patterns than standard request-response. Requires careful resource management and may need architectural changes to support efficiently.

## Commands

### run-api-server

Start the launcher with web API enabled on specified port

```bash
go run cmd/launcher/main.go --web-port=8080
```

### test-api

Run unit tests for API handlers

```bash
go test ./cmd/launcher/web/api/...
```

### test-api-verbose

Run API tests with verbose output showing each test case

```bash
go test -v ./cmd/launcher/web/api/...
```

### curl-health

Check API health endpoint

```bash
curl http://localhost:8080/api/health
```

## Testing

Unit test HTTP handlers using httptest.ResponseRecorder to verify status codes, response bodies, and error handling without starting a real server. Mock the launcher service layer to isolate API logic.

```bash
go test ./cmd/launcher/web/api/...
go test -race ./cmd/launcher/web/api/...
go test -cover ./cmd/launcher/web/api/...
```

Test directory: `cmd/launcher/web/api`

## Common Tasks

### Add new API endpoint

1. 1. Define handler function: func handlerName(w http.ResponseWriter, r *http.Request)
2. 2. Add route registration in API setup/init function
3. 3. Implement request validation and parameter extraction
4. 4. Call underlying launcher service method with request context
5. 5. Format response as JSON with appropriate status code
6. 6. Add error handling for all failure cases
7. 7. Write unit tests covering success and error paths
8. 8. Update API documentation if maintained separately

### Debug API request failure

1. 1. Check server logs for error messages and stack traces
2. 2. Verify request format matches expected JSON schema
3. 3. Test endpoint with curl to isolate client vs server issues
4. 4. Add debug logging to handler to trace execution flow
5. 5. Check if error originates in API layer or service layer
6. 6. Verify context is not cancelled prematurely
7. 7. Use httptest to reproduce issue in unit test

### Handle API versioning

1. 1. Decide on versioning strategy: URL path (/v1/), header, or query param
2. 2. Create new handler functions for changed endpoints
3. 3. Keep old handlers for backward compatibility period
4. 4. Add version routing logic in API setup
5. 5. Document deprecation timeline for old versions
6. 6. Update client code to use new version
7. 7. Remove old version after deprecation period

### Optimize slow API endpoint

1. 1. Add timing metrics to identify bottleneck (handler vs service layer)
2. 2. Check if operation can be made asynchronous with status polling
3. 3. Add caching for expensive read operations if appropriate
4. 4. Implement pagination for large result sets
5. 5. Consider adding timeout to prevent indefinite blocking
6. 6. Profile with pprof if CPU or memory bound
7. 7. Load test with realistic concurrency to verify improvement

