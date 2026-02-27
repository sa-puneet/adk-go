# AGENTS.md - artifact/gcsartifact

> GCS artifact storage implementation providing Google Cloud Storage integration for artifact management. Implements the artifact.Service interface with GCS-specific client operations for uploading, downloading, and managing artifacts in GCS buckets.

## Tech Stack

- **cloud.google.com/go/storage** - Google Cloud Storage client library for GCS operations
- **google.golang.org/api/option** - GCP client options and authentication configuration
- **context** stdlib - Context propagation for GCS operations and cancellation

## Key Files

| File | Purpose | Read When |
|------|---------|------------|
| `artifact/gcsartifact/service.go` (reference) | Main service implementation of artifact.Service interface for GCS | Understanding how GCS artifact operations are structured and exposed |
| `artifact/gcsartifact/gcs_client.go` (reference) | GCS client wrapper with bucket and object operations | Implementing or modifying GCS-specific storage operations |
| `artifact/gcsartifact/gcs_test.go` (reference) | Test suite with mock GCS client patterns | Writing tests or understanding mocking strategy for GCS operations |

## Architecture

```
artifact.Service Interface
         |
         v
    GCS Service (service.go)
         |
         +-- Uses GCSClient interface
         |
         v
    GCS Client (gcs_client.go)
         |
         +-- Wraps storage.Client
         +-- Bucket operations
         +-- Object read/write
         |
         v
    Google Cloud Storage API

Flow: Service methods -> GCSClient interface -> storage.Client -> GCS API
Testing: Mock GCSClient interface for unit tests
```

## Patterns

### Interface Abstraction for External Dependencies

GCSClient interface wraps storage.Client to enable mocking and testing without real GCS connections

See `artifact/gcsartifact/gcs_client.go` for reference.

### Context-First API Design

All GCS operations accept context.Context as first parameter for cancellation and timeout control

See `artifact/gcsartifact/service.go` for reference.

### Bucket-Object Path Separation

GCS paths are split into bucket and object components for proper GCS API usage

See `artifact/gcsartifact/gcs_client.go` for reference.

### Reader/Writer Streaming

Uses io.Reader and io.Writer for efficient streaming of artifact data without loading entire files into memory

See `artifact/gcsartifact/service.go` for reference.

## Rules

### Always (Do These)

| Rule | Rationale |
|------|------------|
| Pass context.Context as first parameter to all GCS operations | Enables proper timeout, cancellation, and request tracing through GCS operations |
| Use GCSClient interface instead of direct storage.Client references in service code | Maintains testability and allows mocking GCS operations without real cloud connections |
| Close readers and writers explicitly with defer statements | Prevents resource leaks and ensures proper cleanup of GCS connections |
| Validate bucket and object paths before GCS operations | Prevents invalid GCS API calls and provides clear error messages |
| Return wrapped errors with context about the operation that failed | Helps debugging by providing operation context (bucket, object name) in error messages |

### Never (Don't Do These)

| Rule | Severity | Rationale |
|------|----------|------------|
| Never use storage.Client directly in service methods | CRITICAL | Breaks testability and creates tight coupling to GCS implementation. Always use GCSClient interface |
| Never ignore errors from Close() on readers/writers | HIGH | Close errors may indicate incomplete writes or connection issues that could corrupt artifacts |
| Never load entire artifact content into memory | HIGH | Artifacts can be large; always use streaming with io.Reader/io.Writer to prevent OOM |
| Never hardcode bucket names or GCS paths in service code | MEDIUM | Paths should come from configuration or parameters to support different environments |
| Never create GCS operations without context timeout | HIGH | GCS operations can hang indefinitely; always use contexts with appropriate timeouts |

### Ask First

- **Adding new GCS client options or authentication methods** - Authentication changes affect security and deployment; coordinate with infrastructure team
- **Changing bucket naming conventions or path structures** - May break existing artifact references and require migration strategy
- **Modifying GCSClient interface methods** - Interface changes require updating all implementations and mocks across test suite
- **Adding caching or buffering layers** - Can affect consistency guarantees and memory usage patterns

## Commands

### run-tests

Run all GCS artifact tests with mocked GCS client

```bash
go test ./artifact/gcsartifact/...
```

### run-tests-verbose

Run tests with verbose output showing individual test cases

```bash
go test -v ./artifact/gcsartifact/...
```

### run-tests-coverage

Run tests with coverage report

```bash
go test -cover ./artifact/gcsartifact/...
```

### lint

Run linter on GCS artifact code

```bash
golangci-lint run ./artifact/gcsartifact/...
```

## Testing

Mock-based unit testing using GCSClient interface. Tests use mock implementations to simulate GCS operations without real cloud connections. Focus on error handling, context propagation, and proper resource cleanup.

```bash
go test ./artifact/gcsartifact/...
go test -race ./artifact/gcsartifact/...
```

Test directory: `artifact/gcsartifact/`

## Common Tasks

### Add new GCS operation to service

1. 1. Add method to GCSClient interface in gcs_client.go
2. 2. Implement method in gcsClient struct wrapping storage.Client
3. 3. Add service method in service.go using GCSClient interface
4. 4. Ensure context.Context is first parameter
5. 5. Add proper error wrapping with operation context
6. 6. Create mock implementation in gcs_test.go
7. 7. Write unit tests covering success and error cases

### Debug GCS operation failure

1. 1. Check error message for bucket and object path details
2. 2. Verify context timeout is appropriate for operation
3. 3. Confirm GCS credentials and permissions are correct
4. 4. Check if bucket exists and is accessible
5. 5. Review GCS client initialization in service.go
6. 6. Add logging around GCSClient calls if needed
7. 7. Test with mock client to isolate service logic from GCS issues

### Implement new artifact storage backend

1. 1. Review artifact.Service interface contract
2. 2. Study GCS implementation pattern in service.go
3. 3. Create client interface similar to GCSClient
4. 4. Implement service methods with context-first design
5. 5. Use io.Reader/Writer for streaming operations
6. 6. Add comprehensive error wrapping
7. 7. Create mock client for testing
8. 8. Write unit tests following gcs_test.go patterns

### Update GCS client configuration

1. 1. Locate client initialization code in service.go
2. 2. Review option.ClientOption usage for authentication
3. 3. Update client options as needed
4. 4. Ensure changes don't break GCSClient interface
5. 5. Update tests to reflect new configuration
6. 6. Test with real GCS connection if possible
7. 7. Document configuration changes

