# AGENTS.md - server/adkrest/controllers

> REST API controllers for the ADK (Agent Development Kit) server, implementing HTTP handlers for sessions, artifacts, debug endpoints, and agent interactions. Controllers follow a dependency injection pattern with service layer separation and consistent error handling using context-aware responses.

## Tech Stack

- **Go** 1.x - Primary language for HTTP controllers and REST API implementation
- **net/http** stdlib - HTTP server and request/response handling
- **encoding/json** stdlib - JSON serialization for API responses

## Key Files

| File | Purpose | Read When |
|------|---------|------------|
| `server/adkrest/controllers/sessions_test.go` (reference) | Test patterns for session management endpoints - reference for testing controller logic | Writing tests for any controller or understanding session lifecycle |
| `server/adkrest/controllers/debug.go` (reference) | Debug endpoints implementation - shows error handling and response patterns | Understanding controller structure, error handling, or adding debug endpoints |
| `server/adkrest/controllers/artifacts.go` (reference) | Artifact management endpoints - demonstrates file handling and streaming responses | Implementing file upload/download or streaming responses |
| `server/adkrest/controllers/sessions.go` (reference) | Session CRUD operations - primary example of controller pattern with service layer | Creating new controllers or understanding session management flow |

## Architecture

```
Request Flow:
  HTTP Request → Controller Handler
    ↓
  Extract/Validate Parameters (path, query, body)
    ↓
  Call Service Layer (business logic)
    ↓
  Format Response (JSON/stream)
    ↓
  Write HTTP Response

Controller Pattern:
  type Controller struct {
    service ServiceInterface
    logger  Logger
  }
  
  func (c *Controller) HandleEndpoint(w http.ResponseWriter, r *http.Request) {
    // 1. Extract params
    // 2. Validate input
    // 3. Call service
    // 4. Handle errors
    // 5. Write response
  }
```

## Patterns

### Dependency Injection Controllers

Controllers receive service dependencies via constructor, enabling testability and separation of concerns. Never instantiate services inside handlers.

See `server/adkrest/controllers/debug.go` for reference.

### Context-Aware Error Handling

Use context.Context throughout request lifecycle. Errors should be wrapped with context and returned to caller, not logged directly in controller.

See `server/adkrest/controllers/sessions.go` for reference.

### JSON Response Helpers

Standardized response writing with proper status codes. Use helper functions for consistent error responses and success payloads.

See `server/adkrest/controllers/artifacts.go` for reference.

### Path Parameter Extraction

Extract URL path parameters early in handler, validate immediately, and fail fast with 400 Bad Request for invalid input.

See `server/adkrest/controllers/sessions.go` for reference.

### Service Layer Delegation

Controllers are thin wrappers - all business logic lives in service layer. Controllers only handle HTTP concerns (parsing, response formatting).

See `server/adkrest/controllers/debug.go` for reference.

## Rules

### Always (Do These)

| Rule | Rationale |
|------|------------|
| Pass context.Context as first parameter to all service calls | Enables request cancellation, timeout propagation, and distributed tracing across service boundaries |
| Validate all input parameters before calling service layer | Fail fast with clear error messages. Service layer should receive validated data only. |
| Use dependency injection for all external dependencies (services, loggers) | Enables unit testing with mocks and maintains loose coupling between layers |
| Return appropriate HTTP status codes: 200 (success), 201 (created), 400 (bad request), 404 (not found), 500 (server error) | RESTful API conventions - clients depend on correct status codes for error handling |
| Write response headers before writing body, set Content-Type explicitly | HTTP protocol requirement - headers must be written before body, explicit Content-Type prevents ambiguity |
| Use table-driven tests for controller endpoints with multiple scenarios | Ensures comprehensive coverage of success/error paths with maintainable test structure |

### Never (Don't Do These)

| Rule | Severity | Rationale |
|------|----------|------------|
| Never instantiate service layer objects inside controller handlers | CRITICAL | Breaks testability and creates tight coupling. Services must be injected via constructor. |
| Never log errors directly in controllers - return them to caller | HIGH | Logging should happen at application boundary (middleware/main). Controllers should propagate errors with context. |
| Never ignore errors from json.Encode/Decode or response writing | HIGH | Silent failures lead to incomplete responses and difficult debugging. Always check and handle encoding errors. |
| Never put business logic in controllers | CRITICAL | Controllers are HTTP adapters only. Business logic belongs in service layer for reusability and testability. |
| Never write to response body after WriteHeader has been called with error status | MEDIUM | Can cause malformed HTTP responses. Write complete error response in single operation. |
| Never use panic for error handling in controllers | HIGH | Panics crash the server. Use proper error returns and let middleware handle recovery if needed. |

### Ask First

- **Adding new HTTP endpoints or changing existing endpoint signatures** - API changes affect clients and require coordination. Verify backward compatibility and update API documentation.
- **Changing error response format or status codes** - Clients may depend on specific error formats. Breaking changes require versioning or migration strategy.
- **Adding authentication or authorization checks to controllers** - Security changes should be reviewed and may require middleware implementation rather than per-controller logic.
- **Implementing streaming or long-polling endpoints** - Requires different patterns than standard request-response. Discuss timeout handling and resource management.

## Commands

### run-controller-tests

Run all controller tests with verbose output

```bash
go test ./server/adkrest/controllers/... -v
```

### test-specific-controller

Run tests for specific controller (replace TestSessionController with target)

```bash
go test ./server/adkrest/controllers -run TestSessionController -v
```

### test-with-coverage

Generate and view test coverage report for controllers

```bash
go test ./server/adkrest/controllers/... -coverprofile=coverage.out && go tool cover -html=coverage.out
```

### lint-controllers

Run linter on controller code

```bash
golangci-lint run ./server/adkrest/controllers/...
```

## Testing

Table-driven tests with httptest.ResponseRecorder for HTTP response validation. Mock service layer dependencies using interfaces. Test both success paths and error conditions (400, 404, 500 responses).

```bash
go test ./server/adkrest/controllers/... -v
go test -race ./server/adkrest/controllers/...
go test -coverprofile=coverage.out ./server/adkrest/controllers/...
```

Test directory: `server/adkrest/controllers/*_test.go`

## Common Tasks

### Add new REST endpoint

1. 1. Define handler method on controller struct: func (c *Controller) HandleNewEndpoint(w http.ResponseWriter, r *http.Request)
2. 2. Extract and validate request parameters (path, query, body)
3. 3. Call service layer method with context: result, err := c.service.DoSomething(r.Context(), params)
4. 4. Handle errors with appropriate status codes (400 for validation, 404 for not found, 500 for internal)
5. 5. Write JSON response with proper Content-Type header
6. 6. Add table-driven tests covering success and error cases
7. 7. Register route in router configuration (see server/adkrest/router.go or similar)

### Add error handling to existing endpoint

1. 1. Identify service call that needs error handling
2. 2. Check error return: if err != nil { /* handle */ }
3. 3. Determine appropriate HTTP status code based on error type
4. 4. Write error response with descriptive message: http.Error(w, err.Error(), statusCode)
5. 5. Add test case for error scenario in *_test.go file
6. 6. Verify error propagates context information for debugging

### Mock service for controller testing

1. 1. Create mock struct implementing service interface in *_test.go
2. 2. Add fields to mock for controlling return values: mockResult, mockError
3. 3. Implement interface methods to return mock values
4. 4. Instantiate controller with mock service in test setup
5. 5. Set mock return values for each test case
6. 6. Verify controller behavior with different mock responses

### Handle file upload in controller

1. 1. Parse multipart form: r.ParseMultipartForm(maxMemory)
2. 2. Extract file from form: file, header, err := r.FormFile('fieldname')
3. 3. Validate file type and size constraints
4. 4. Pass file reader to service layer for processing
5. 5. Return appropriate response (201 Created with resource location)
6. 6. See server/adkrest/controllers/artifacts.go for reference implementation

