# AGENTS.md - cmd/launcher/web/api

> This package provides the HTTP API layer for the launcher web interface, handling REST endpoints for agent management and control. It serves as the bridge between the web UI and the underlying launcher functionality, exposing operations like starting, stopping, and monitoring agents.

## Tech Stack

- **Go** 1.x - Primary language for HTTP API implementation
- **net/http** stdlib - HTTP server and routing for REST API endpoints
- **encoding/json** stdlib - JSON serialization for API request/response handling

## Key Files

| File | Purpose | Read When |
|------|---------|------------|
| `cmd/launcher/web/api/api.go` (reference) | Main API handler implementation with route definitions and endpoint logic | Understanding API structure, adding new endpoints, or debugging web requests |

## Architecture

```
HTTP Request → API Handler (api.go) → Launcher Core Logic → Response

The API layer acts as a thin REST interface:
1. Client sends HTTP request to /api/* endpoint
2. api.go routes to appropriate handler function
3. Handler validates request, calls launcher operations
4. Response serialized to JSON and returned

This is a presentation layer - business logic lives in parent launcher package.
```

## Patterns

### HTTP Handler Functions

Standard Go HTTP handler pattern with http.HandlerFunc signature, using ResponseWriter and Request parameters

See `cmd/launcher/web/api/api.go` for reference.

### JSON Request/Response

All API endpoints use JSON for data exchange, with json.NewDecoder for requests and json.NewEncoder for responses

See `cmd/launcher/web/api/api.go` for reference.

### Error Response Structure

Consistent error response format with HTTP status codes and JSON error messages

See `cmd/launcher/web/api/api.go` for reference.

## Rules

### Always (Do These)

| Rule | Rationale |
|------|------------|
| Return proper HTTP status codes (200 OK, 400 Bad Request, 500 Internal Server Error) based on operation outcome | REST API clients depend on status codes for error handling and flow control |
| Validate all incoming request data before processing | API is external-facing and must protect against malformed or malicious input |
| Use json.NewEncoder(w).Encode() for response serialization rather than json.Marshal | Direct encoding to ResponseWriter is more efficient and handles streaming properly |
| Set Content-Type: application/json header for all JSON responses | Clients need proper content type to parse responses correctly |
| Keep API handlers thin - delegate business logic to launcher core packages | API layer should only handle HTTP concerns, not implement agent logic |

### Never (Don't Do These)

| Rule | Severity | Rationale |
|------|----------|------------|
| Never implement agent lifecycle logic directly in API handlers | CRITICAL | API is presentation layer only - business logic belongs in launcher core to maintain separation of concerns and testability |
| Never expose internal error details or stack traces in API responses | HIGH | Internal errors can leak sensitive information; return user-friendly messages instead |
| Never block API handlers with long-running operations | HIGH | HTTP requests should return quickly; use async patterns for long operations and provide status endpoints |
| Never ignore errors from json.Decode or json.Encode operations | MEDIUM | Serialization errors indicate malformed data or system issues that must be handled |
| Never use panic() in API handlers | CRITICAL | Panics crash the entire server; use proper error returns and recovery middleware if needed |

### Ask First

- **Adding authentication or authorization middleware** - Security model needs to be consistent across the launcher; coordinate with overall architecture
- **Changing API endpoint paths or request/response formats** - Breaking changes affect web UI and any external clients; requires versioning strategy
- **Adding WebSocket or streaming endpoints** - Different protocol handling requires architectural decision about connection management
- **Implementing rate limiting or request throttling** - Should be coordinated with overall launcher resource management strategy

## Commands

### test-api

Run unit tests for API handlers

```bash
go test ./cmd/launcher/web/api/...
```

### build-launcher

Build launcher binary including API server

```bash
go build ./cmd/launcher
```

### test-endpoint

Manual API endpoint testing with curl

```bash
curl -X POST http://localhost:8080/api/endpoint -H 'Content-Type: application/json' -d '{}'
```

## Testing

Test API handlers using httptest.ResponseRecorder to verify HTTP responses without starting actual server. Mock launcher core dependencies to isolate API layer testing.

```bash
go test ./cmd/launcher/web/api/... -v
go test ./cmd/launcher/web/api/... -cover
```

Test directory: `cmd/launcher/web/api`

## Common Tasks

### Add new API endpoint

1. 1. Define handler function with signature: func(w http.ResponseWriter, r *http.Request)
2. 2. Add request struct type if endpoint accepts JSON body
3. 3. Add response struct type for JSON response
4. 4. Implement handler: decode request, validate, call launcher logic, encode response
5. 5. Register route in main API setup (likely in parent web package)
6. 6. Add unit tests using httptest package
7. 7. Document endpoint in API documentation

### Debug API request handling

1. 1. Check cmd/launcher/web/api/api.go for handler implementation
2. 2. Add logging at handler entry to verify request reaches handler
3. 3. Log request body and headers to verify client sends correct data
4. 4. Verify JSON decoding succeeds and check for validation errors
5. 5. Check launcher core operation returns expected results
6. 6. Verify response encoding and HTTP status code setting
7. 7. Use curl or Postman to test endpoint independently

### Handle API errors consistently

1. 1. Define error response struct with 'error' field for message
2. 2. Create helper function: writeError(w http.ResponseWriter, status int, message string)
3. 3. Use appropriate HTTP status: 400 for client errors, 500 for server errors
4. 4. Log internal error details server-side for debugging
5. 5. Return user-friendly error message in response (no stack traces)
6. 6. Set Content-Type: application/json even for error responses
7. 7. Test error paths in unit tests to verify proper handling

### Validate API request data

1. 1. Define request struct with json tags for field mapping
2. 2. Use json.NewDecoder(r.Body).Decode(&req) to parse request
3. 3. Check for decode errors and return 400 Bad Request if malformed
4. 4. Implement validation logic: check required fields, value ranges, formats
5. 5. Return 400 with specific validation error message if invalid
6. 6. Only proceed to business logic after validation passes
7. 7. Consider using validation library for complex validation rules

