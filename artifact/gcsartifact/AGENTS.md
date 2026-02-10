# AGENTS.md - artifact/gcsartifact

> GCS artifact storage implementation providing Google Cloud Storage integration for the ADK artifact system. Implements the artifact.Service interface with GCS-backed storage, handling uploads, downloads, and metadata management with proper error handling and context propagation.

## Tech Stack

- **cloud.google.com/go/storage** - Google Cloud Storage client library for artifact storage operations
- **google.golang.org/api/option** - GCS client configuration and authentication options
- **github.com/googleapis/gax-go/v2** v2 - Google API extensions for retry logic and error handling

## Key Files

| File | Purpose | Read When |
|------|---------|------------|
| `artifact/gcsartifact/service.go` (reference) | Main Service implementation with Upload/Download/Delete/List methods implementing artifact.Service interface | Understanding the public API and how GCS operations are orchestrated |
| `artifact/gcsartifact/gcs_client.go` (reference) | GCS client wrapper abstracting storage operations with proper error handling and context management | Understanding low-level GCS interactions and client lifecycle |
| `artifact/gcsartifact/gcs_test.go` (reference) | Test suite with mock GCS client demonstrating proper testing patterns for storage operations | Writing tests or understanding expected behavior and error cases |

## Architecture

```
Client Request → Service (artifact.Service interface)
    ↓
  Service validates input (bucket, key)
    ↓
  GCSClient wrapper (abstraction layer)
    ↓
  GCS Storage API (cloud.google.com/go/storage)
    ↓
  Google Cloud Storage

Key Components:
- Service: Public API implementing artifact.Service
- GCSClient: Abstraction over GCS SDK with error handling
- Mock interfaces: gcsClientInterface for testing

Data Flow:
- Upload: io.Reader → GCS Writer → Bucket/Key
- Download: Bucket/Key → GCS Reader → io.Writer
- Metadata: Stored as GCS object attributes
```

## Patterns

### Interface Abstraction Pattern

gcsClientInterface abstracts GCS operations enabling mock testing without real GCS dependencies

See `artifact/gcsartifact/gcs_client.go` for reference.

### Context Propagation

All operations accept context.Context as first parameter, enabling cancellation and timeout control throughout the call chain

See `artifact/gcsartifact/service.go` for reference.

### Error Wrapping with Context

Errors wrapped with fmt.Errorf and %w verb to preserve error chains while adding context about bucket/key operations

See `artifact/gcsartifact/service.go` for reference.

### Dependency Injection via Constructor

NewService accepts GCS client options allowing custom credentials, endpoints, or mock clients for testing

See `artifact/gcsartifact/service.go` for reference.

### Deferred Resource Cleanup

defer statements ensure readers/writers are closed even on error paths, preventing resource leaks

See `artifact/gcsartifact/gcs_client.go` for reference.

## Rules

### Always (Do These)

| Rule | Rationale |
|------|------------|
| CRITICAL: Pass context.Context as first parameter to all Service methods | Enables proper cancellation, timeout control, and request tracing through GCS operations. GCS SDK requires context for all operations. |
| CRITICAL: Validate bucket and key parameters are non-empty before GCS operations | Empty bucket/key causes cryptic GCS errors. Early validation provides clear error messages and prevents unnecessary API calls. |
| HIGH: Use defer for closing GCS readers and writers | Ensures resources are released even on error paths. GCS connections must be properly closed to avoid leaks and ensure data is flushed. |
| HIGH: Wrap errors with fmt.Errorf and %w to preserve error chains | Maintains error context for debugging while allowing errors.Is/As checks. Critical for distinguishing GCS errors (not found, permission, etc). |
| HIGH: Use gcsClientInterface abstraction for all GCS operations | Enables testing with mock clients without real GCS dependencies. Maintains clean separation between business logic and GCS SDK. |
| MEDIUM: Include bucket and key in error messages | Provides context for debugging multi-bucket/key operations. Essential for production troubleshooting. |
| MEDIUM: Use io.Copy for streaming data between readers/writers | Efficient memory usage for large artifacts. Avoids loading entire files into memory. |

### Never (Don't Do These)

| Rule | Severity | Rationale |
|------|----------|------------|
| NEVER call GCS SDK directly from Service methods | CRITICAL | Breaks testability and abstraction. All GCS operations must go through gcsClientInterface to enable mocking. |
| NEVER ignore errors from Close() on GCS writers | CRITICAL | GCS writers buffer data; Close() flushes and returns upload errors. Ignoring Close() errors means data may not be persisted. |
| NEVER use empty strings for bucket or key parameters | HIGH | Causes GCS API errors that are harder to debug. Validate at service boundary before making API calls. |
| NEVER create new GCS clients per operation | HIGH | GCS clients are expensive to create and maintain connection pools. Reuse single client instance across operations. |
| NEVER load entire artifact into memory before upload/download | HIGH | Artifacts can be large (GBs). Use streaming io.Reader/Writer to handle arbitrary sizes efficiently. |
| NEVER return raw GCS SDK errors without wrapping | MEDIUM | Loses operation context (which bucket/key failed). Wrap errors with context for debugging. |

### Ask First

- **Adding new GCS client options or configuration** - May affect authentication, retry behavior, or connection pooling across all operations. Coordinate with team on defaults.
- **Changing error handling or error types returned** - Service implements artifact.Service interface; error contract changes affect all callers and may break compatibility.
- **Adding metadata or custom attributes to GCS objects** - Affects storage format and may impact other systems reading artifacts. Ensure backward compatibility.
- **Implementing caching or local buffering** - Changes performance characteristics and may introduce consistency issues. Discuss trade-offs with team.
- **Modifying bucket/key naming conventions or path structure** - Affects artifact discoverability and may break existing references. Requires migration strategy.

## Commands

### run-tests

Run all tests for GCS artifact implementation

```bash
go test ./artifact/gcsartifact/...
```

### run-tests-verbose

Run tests with verbose output showing individual test results

```bash
go test -v ./artifact/gcsartifact/...
```

### run-tests-coverage

Run tests with coverage report

```bash
go test -cover ./artifact/gcsartifact/...
```

### check-imports

Check for import formatting issues

```bash
goimports -l artifact/gcsartifact/
```

### lint

Run linter on GCS artifact code

```bash
golangci-lint run ./artifact/gcsartifact/...
```

## Testing

Mock-based unit testing using gcsClientInterface abstraction. Tests verify Service behavior without real GCS dependencies. Mock client simulates success/error cases for upload, download, delete, and list operations.

```bash
go test ./artifact/gcsartifact/...
go test -v -run TestService_Upload ./artifact/gcsartifact/
go test -v -run TestService_Download ./artifact/gcsartifact/
```

Test directory: `artifact/gcsartifact/`

## Common Tasks

### Add new Service method implementing artifact.Service interface

1. 1. Add method signature to Service struct in service.go matching artifact.Service interface
2. 2. Accept context.Context as first parameter
3. 3. Validate input parameters (bucket, key non-empty)
4. 4. Call appropriate gcsClientInterface method
5. 5. Wrap any errors with fmt.Errorf including bucket/key context
6. 6. Add corresponding method to gcsClientInterface in gcs_client.go
7. 7. Implement method in gcsClient struct with proper error handling
8. 8. Add mock implementation in gcs_test.go
9. 9. Write table-driven tests covering success and error cases

### Debug GCS operation failure

1. 1. Check error message for bucket/key context (should be wrapped)
2. 2. Verify bucket name and key are non-empty and valid
3. 3. Check GCS client initialization in NewService - verify options
4. 4. Examine context cancellation/timeout - may be too short for large artifacts
5. 5. Check GCS permissions for service account (read/write/delete)
6. 6. Use GCS console to verify object exists/doesn't exist as expected
7. 7. Add logging around gcsClient methods to trace operation flow
8. 8. Check if error is storage.ErrObjectNotExist for not found cases

### Add new GCS client configuration option

1. 1. Identify required option.ClientOption from GCS SDK
2. 2. Add parameter to NewService function signature in service.go
3. 3. Pass option to storage.NewClient in NewService
4. 4. Update all NewService call sites to provide option
5. 5. Add test case in gcs_test.go verifying option is used
6. 6. Document option purpose and default behavior in godoc
7. 7. Consider if option should be optional (variadic) or required

### Implement new artifact metadata field

1. 1. Determine if metadata should be GCS object attributes or separate object
2. 2. Update Upload method to set metadata via writer.Metadata or writer.ContentType
3. 3. Update Download/List to read metadata from object.Attrs()
4. 4. Add metadata to gcsClientInterface methods if needed
5. 5. Update mock client in tests to return metadata
6. 6. Write tests verifying metadata round-trips correctly
7. 7. Document metadata format and any size/type limitations

### Test with real GCS (integration test)

1. 1. Create test GCS bucket with appropriate permissions
2. 2. Set GOOGLE_APPLICATION_CREDENTIALS or use default credentials
3. 3. Create NewService with real GCS client (no mocks)
4. 4. Use unique key prefix to avoid conflicts (e.g., test-{timestamp}-)
5. 5. Test Upload → Download → verify content matches
6. 6. Test Delete → verify object no longer exists
7. 7. Clean up test objects in defer or test cleanup
8. 8. Consider using testcontainers or GCS emulator for CI

