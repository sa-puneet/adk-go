# AGENTS.md - server/adkrest/controllers

> REST API controllers for the ADK (Agent Development Kit) server, implementing HTTP handlers for sessions, artifacts, debug endpoints, and agent interactions. Controllers follow a dependency injection pattern with service layer separation and consistent error handling through HTTP status codes.

## Tech Stack

- **Go** 1.x - Primary language for HTTP controllers and REST API implementation
- **net/http** stdlib - HTTP server and request/response handling
- **encoding/json** stdlib - JSON serialization for API responses and requests

## Key Files

| File | Purpose | Read When |
|------|---------|------------|
| `server/adkrest/controllers/sessions_test.go` (reference) | Test patterns for session management endpoints - reference for testing HTTP handlers | Writing tests for any controller or understanding session lifecycle |
| `server/adkrest/controllers/debug.go` (reference) | Debug endpoints implementation - shows error handling and response patterns | Understanding controller structure or adding new debug endpoints |
| `server/adkrest/controllers/artifacts.go` (reference) | Artifact management endpoints - demonstrates file handling and streaming responses | Implementing file upload/download or artifact-related features |

## Architecture

```
HTTP Request → Controller (validates/parses) → Service Layer → Business Logic → Response

Controller Responsibilities:
- Parse HTTP request (body, params, headers)
- Validate input data
- Call service layer methods
- Transform service responses to HTTP responses
- Set appropriate status codes and headers
- Handle errors with proper HTTP semantics

Pattern: Controllers are thin wrappers that delegate to services
```

## Patterns

### Dependency Injection Controllers

Controllers receive service dependencies via constructor, not global state. Each controller struct holds references to required services.

See `server/adkrest/controllers/artifacts.go` for reference.

### HTTP Error Response Pattern

Errors are translated to HTTP status codes (400 for validation, 404 for not found, 500 for internal). Use http.Error() or json.NewEncoder for consistent responses.

See `server/adkrest/controllers/debug.go` for reference.

### JSON Request/Response Handling

Use json.NewDecoder(r.Body).Decode() for requests and json.NewEncoder(w).Encode() for responses. Always set Content-Type: application/json header.

See `server/adkrest/controllers/sessions_test.go` for reference.

### Path Parameter Extraction

Extract URL path parameters using request context or mux variables. Validate parameters before passing to service layer.

See `server/adkrest/controllers/artifacts.go` for reference.

### Service Layer Delegation

Controllers never contain business logic. All operations are delegated to service layer methods that return domain objects or errors.

See `server/adkrest/controllers/debug.go` for reference.

## Rules

### Always (Do These)

| Rule | Rationale |
|------|------------|
| CRITICAL: Validate all input parameters before calling service layer methods | Controllers are the security boundary. Invalid input must be rejected with 400 Bad Request before reaching business logic. |
| HIGH: Set Content-Type header to 'application/json' for all JSON responses | Clients expect proper content type headers. Missing headers cause parsing errors in client libraries. |
| HIGH: Use http.StatusXXX constants instead of numeric codes | Makes code self-documenting and prevents typos in status codes (e.g., http.StatusNotFound vs 404). |
| HIGH: Close request body after reading with defer r.Body.Close() | Prevents resource leaks in long-running servers. Even though http.Server closes bodies, explicit closure is best practice. |
| MEDIUM: Return early on errors to avoid nested if-else blocks | Improves readability and follows Go idioms. Handle error case immediately and return. |
| MEDIUM: Log errors before returning HTTP error responses | Enables debugging and monitoring. Include request context (session ID, user ID) in logs. |

### Never (Don't Do These)

| Rule | Severity | Rationale |
|------|----------|------------|
| CRITICAL: Never put business logic in controllers | CRITICAL | Controllers are HTTP adapters only. Business logic belongs in service layer for testability and reusability. |
| CRITICAL: Never ignore errors from json.Decode() or json.Encode() | CRITICAL | Silent failures lead to incorrect responses or panics. Always check and return appropriate HTTP errors. |
| HIGH: Never use global variables for service dependencies | HIGH | Makes testing impossible and creates hidden dependencies. Use dependency injection via constructor. |
| HIGH: Never return 200 OK with error messages in response body | HIGH | HTTP status codes must reflect actual status. Use 4xx for client errors, 5xx for server errors. |
| HIGH: Never expose internal error details to clients in production | HIGH | Internal errors may leak sensitive information. Return generic messages and log details server-side. |
| MEDIUM: Never write to response after calling WriteHeader() | MEDIUM | Headers cannot be modified after first write. Set all headers before writing body. |

### Ask First

- **Adding new HTTP endpoints or routes** - Must ensure consistency with existing API patterns, versioning strategy, and OpenAPI spec updates
- **Changing error response format or status codes** - Breaking change for API clients. Requires versioning strategy and client migration plan
- **Adding authentication or authorization checks** - Security-critical changes that affect all endpoints. Must align with overall auth strategy
- **Modifying request/response schemas** - API contract changes require coordination with clients and may need backward compatibility
- **Adding middleware or request interceptors** - Affects all requests. Must consider performance impact and interaction with existing middleware

## Commands

### run-tests

Run all controller tests with coverage

```bash
go test ./server/adkrest/controllers/...
```

### test-verbose

Run tests with verbose output to see individual test results

```bash
go test -v ./server/adkrest/controllers/...
```

### test-coverage

Run tests and display coverage percentage

```bash
go test -cover ./server/adkrest/controllers/...
```

### lint

Run linter to check code quality and style

```bash
golangci-lint run ./server/adkrest/controllers/...
```

## Testing

Table-driven tests with httptest.ResponseRecorder for HTTP handler testing. Mock service layer dependencies. Test happy path, error cases, and edge cases separately.

```bash
go test ./server/adkrest/controllers/... -run TestSessionsController
go test ./server/adkrest/controllers/... -cover -coverprofile=coverage.out
```

Test directory: `server/adkrest/controllers/`

## Common Tasks

### Add new REST endpoint

1. 1. Define handler method on controller struct with signature: func (c *Controller) HandlerName(w http.ResponseWriter, r *http.Request)
2. 2. Parse and validate request parameters/body using json.Decoder or URL params
3. 3. Call service layer method with validated inputs
4. 4. Handle service errors and map to appropriate HTTP status codes
5. 5. Encode response using json.NewEncoder(w).Encode() with proper status code
6. 6. Add tests in *_test.go file using httptest package
7. 7. Register route in router configuration (outside controllers package)

### Handle file upload in controller

1. 1. Parse multipart form with r.ParseMultipartForm(maxMemory)
2. 2. Get file from r.FormFile(fieldName)
3. 3. Validate file type, size, and name
4. 4. Pass file reader to service layer (don't read entire file in controller)
5. 5. Return appropriate response with file metadata
6. 6. See server/adkrest/controllers/artifacts.go for reference implementation

### Implement error handling in controller

1. 1. Check error returned from service layer
2. 2. Determine appropriate HTTP status code (400 for validation, 404 for not found, 500 for internal)
3. 3. Log error with context (request ID, user ID, etc.)
4. 4. Return sanitized error message to client (no internal details)
5. 5. Use http.Error() for simple text responses or json.Encoder for structured errors
6. 6. See server/adkrest/controllers/debug.go for error handling patterns

### Write controller unit test

1. 1. Create test file *_test.go in same package
2. 2. Define table-driven test cases with input and expected output
3. 3. Create mock service layer using interface
4. 4. Use httptest.NewRequest() to create test request
5. 5. Use httptest.NewRecorder() to capture response
6. 6. Call controller handler method
7. 7. Assert status code, headers, and response body
8. 8. See server/adkrest/controllers/sessions_test.go for examples

