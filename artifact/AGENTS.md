# AGENTS.md - artifact

> The artifact package provides storage and retrieval mechanisms for agent-generated artifacts (files, data, content). It defines the core Store interface and includes an in-memory implementation for testing/development, along with comprehensive request validation for artifact operations.

## Tech Stack

- **Go** 1.x - Core implementation language for artifact storage abstraction
- **context** stdlib - Context propagation for cancellation and timeouts in storage operations

## Key Files

| File | Purpose | Read When |
|------|---------|------------|
| `artifact/store.go` (reference) | Defines the Store interface - the contract all artifact storage implementations must follow | Understanding how artifacts are stored/retrieved or implementing a new storage backend |
| `artifact/inmemory.go` (reference) | In-memory Store implementation using sync.Map for concurrent access - reference implementation | Implementing a new storage backend or understanding the Store interface semantics |
| `artifact/request_validation.go` | Validation logic for artifact requests - ensures IDs, names, and operations are valid | Adding new validation rules or understanding what makes a valid artifact request |
| `artifact/request_validation_test.go` (reference) | Comprehensive test suite for validation rules - documents all edge cases and requirements | Understanding validation requirements or adding new validation logic |
| `artifact/types.go` (reference) | Core data structures: Artifact, CreateRequest, GetRequest, ListRequest, etc. | Working with artifact data structures or understanding the artifact model |

## Architecture

```
┌─────────────┐
│   Agent     │
│  (creates)  │
└──────┬──────┘
       │ CreateRequest
       ▼
┌─────────────────┐
│ Validation      │ (request_validation.go)
│ - ID format     │
│ - Name rules    │
│ - Content check │
└──────┬──────────┘
       │ Valid Request
       ▼
┌─────────────────┐
│ Store Interface │ (store.go)
│ - Create()      │
│ - Get()         │
│ - List()        │
│ - Delete()      │
└──────┬──────────┘
       │
       ▼
┌─────────────────┐
│ Implementation  │
│ - InMemory      │ (inmemory.go)
│ - [Future: DB]  │
│ - [Future: S3]  │
└─────────────────┘
```

## Patterns

### Interface-Based Storage Abstraction

Store interface allows swapping storage backends (in-memory, database, cloud) without changing consumer code. All implementations must handle context cancellation.

See `artifact/store.go` for reference.

### Request Validation Separation

Validation logic is separated from storage logic. Validate requests before passing to Store implementations to ensure consistent error handling.

See `artifact/request_validation.go` for reference.

### Concurrent-Safe In-Memory Storage

InMemory implementation uses sync.Map for lock-free concurrent reads/writes. Good pattern for test fixtures and development.

See `artifact/inmemory.go` for reference.

### Structured Request Objects

Operations use typed request structs (CreateRequest, GetRequest, etc.) rather than raw parameters. Enables validation and future extensibility.

See `artifact/types.go` for reference.

### Context-First API Design

All Store methods accept context.Context as first parameter for cancellation, timeouts, and tracing propagation.

See `artifact/store.go` for reference.

## Rules

### Always (Do These)

| Rule | Rationale |
|------|------------|
| CRITICAL: Validate all artifact requests before passing to Store implementations | Store implementations assume valid input. Validation ensures consistent error messages and prevents invalid data from entering storage. |
| CRITICAL: Pass context.Context as first parameter to all Store methods | Enables cancellation, timeouts, and distributed tracing. Required for production-grade storage implementations. |
| HIGH: Use sync.Map or equivalent for concurrent in-memory storage | Multiple goroutines may access artifacts simultaneously. sync.Map provides lock-free reads and safe concurrent writes. |
| HIGH: Return ErrNotFound when artifact doesn't exist in Get/Delete operations | Consistent error handling across implementations. Consumers expect this specific error for missing artifacts. |
| MEDIUM: Include artifact ID in error messages for debugging | Makes troubleshooting easier when operations fail. See validation error messages for examples. |
| MEDIUM: Test both valid and invalid inputs comprehensively | Validation logic has many edge cases (empty IDs, special characters, nil content). Tests document expected behavior. |

### Never (Don't Do These)

| Rule | Severity | Rationale |
|------|----------|------------|
| NEVER store artifacts without validating the request first | CRITICAL | Invalid artifacts (empty IDs, malformed names) can corrupt storage and cause downstream failures. Always validate before Create(). |
| NEVER ignore context cancellation in Store implementations | CRITICAL | Ignoring context can cause resource leaks and hung operations. Check ctx.Err() for long-running operations. |
| NEVER use regular maps without locking for concurrent artifact storage | CRITICAL | Concurrent map access causes panics. Use sync.Map, sync.RWMutex, or channels for concurrent access. |
| NEVER return nil for List() when no artifacts exist | HIGH | Return empty slice []Artifact{} instead of nil to avoid nil pointer dereferences in consumers. |
| NEVER modify Artifact objects after storing them | HIGH | Artifacts should be immutable after creation. Modifications can cause race conditions and unexpected behavior. |
| NEVER assume artifact IDs are globally unique across stores | MEDIUM | IDs are unique within a Store instance but not across different storage backends. Include store context when needed. |

### Ask First

- **Adding new storage backend implementations (database, S3, etc.)** - Must ensure implementation correctly handles context, concurrency, and error semantics. Review Store interface contract carefully.
- **Changing artifact ID format or validation rules** - ID format changes affect all consumers and may break existing stored artifacts. Requires migration strategy.
- **Adding new fields to Artifact struct** - Changes affect serialization, storage, and all consumers. Consider backward compatibility and migration path.
- **Implementing artifact versioning or history** - Significant architectural change. Current design assumes single version per ID. Needs design review.

## Commands

### test-artifact

Run all artifact package tests including validation and storage tests

```bash
go test ./artifact/...
```

### test-artifact-verbose

Run tests with verbose output to see individual test cases

```bash
go test -v ./artifact/...
```

### test-artifact-coverage

Run tests with coverage report to ensure validation paths are tested

```bash
go test -cover ./artifact/...
```

### benchmark-inmemory

Benchmark in-memory storage performance for concurrent operations

```bash
go test -bench=. -benchmem ./artifact/
```

## Testing

Table-driven tests for validation logic with comprehensive edge cases. Unit tests for Store implementations verify interface contract (create, get, list, delete, error handling). Tests use in-memory implementation as reference.

```bash
go test ./artifact/...
go test -race ./artifact/...
```

Test directory: `artifact/`

## Common Tasks

### Implement a new storage backend

1. 1. Review artifact/store.go to understand the Store interface contract
2. 2. Study artifact/inmemory.go as reference implementation
3. 3. Implement all Store methods: Create, Get, List, Delete
4. 4. Handle context cancellation in long-running operations (check ctx.Err())
5. 5. Return ErrNotFound for missing artifacts in Get/Delete
6. 6. Return empty slice (not nil) from List when no artifacts exist
7. 7. Ensure thread-safety for concurrent operations
8. 8. Write tests covering all Store interface methods and error cases
9. 9. Test with -race flag to verify concurrent access safety

### Add new artifact validation rule

1. 1. Add validation logic to artifact/request_validation.go
2. 2. Return descriptive error with artifact ID for debugging
3. 3. Add test cases to artifact/request_validation_test.go covering valid and invalid inputs
4. 4. Update documentation if validation affects API contract
5. 5. Consider backward compatibility if changing existing rules

### Store and retrieve an artifact

1. 1. Create artifact.CreateRequest with ID, Name, Content, and Metadata
2. 2. Validate request using validation functions (if not done by Store)
3. 3. Call store.Create(ctx, request) with appropriate context
4. 4. Handle errors (validation failures, storage errors)
5. 5. To retrieve: call store.Get(ctx, artifact.GetRequest{ID: id})
6. 6. Handle ErrNotFound if artifact doesn't exist
7. 7. Use returned Artifact object (do not modify after retrieval)

### Debug artifact storage issues

1. 1. Check validation errors first - most issues are invalid requests
2. 2. Verify artifact ID format matches validation rules
3. 3. Check context isn't cancelled before storage operation
4. 4. For concurrent issues: run tests with -race flag
5. 5. Review error messages - they include artifact IDs for tracing
6. 6. Use List() to verify artifact was actually stored
7. 7. Check Store implementation logs if available

