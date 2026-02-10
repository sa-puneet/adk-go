# AGENTS.md - artifact

> The artifact package provides storage and retrieval mechanisms for conversation artifacts (files, data, resources) generated during agent interactions. It defines interfaces for artifact stores with an in-memory implementation and comprehensive request validation for artifact operations.

## Tech Stack

- **Go** 1.x - Core implementation language with interface-based design for artifact storage abstraction
- **sync.RWMutex** stdlib - Thread-safe concurrent access to in-memory artifact storage

## Key Files

| File | Purpose | Read When |
|------|---------|------------|
| `artifact/inmemory.go` (reference) | In-memory implementation of Store interface with thread-safe map storage | Understanding the reference implementation or implementing custom stores |
| `artifact/request_validation_test.go` (reference) | Comprehensive validation test suite covering all artifact operation edge cases | Understanding validation rules, error handling patterns, or adding new validations |
| `artifact/store.go` | Core Store interface definition and artifact data structures | Implementing new store backends or understanding the contract |
| `artifact/validation.go` | Request validation logic for Create, Get, List, Update, Delete operations | Understanding validation requirements or debugging validation errors |

## Architecture

```
┌─────────────┐
│   Client    │
│  (Agent)    │
└──────┬──────┘
       │ Create/Get/List/Update/Delete
       ▼
┌─────────────────────┐
│  Validation Layer   │ ◄── validation.go
│  - validateCreate() │
│  - validateGet()    │
│  - validateList()   │
└──────┬──────────────┘
       │ Valid Request
       ▼
┌─────────────────────┐
│   Store Interface   │ ◄── store.go
└──────┬──────────────┘
       │
       ▼
┌─────────────────────┐
│ InMemory/Custom     │ ◄── inmemory.go
│ Implementation      │
│ (RWMutex protected) │
└─────────────────────┘
```

## Patterns

### Interface-Based Storage Abstraction

Store interface allows pluggable backends (in-memory, database, cloud storage). All implementations must handle thread-safety and validation.

See `artifact/store.go, artifact/inmemory.go` for reference.

### Request Validation Before Storage

All operations validate requests (non-empty IDs, valid filters) before touching storage. Validation returns specific error types for different failure modes.

See `artifact/validation.go` for reference.

### Thread-Safe Map with RWMutex

InMemory store uses sync.RWMutex for concurrent read/write access. Multiple readers allowed, exclusive writer lock for modifications.

See `artifact/inmemory.go` for reference.

### Immutable Artifact Metadata

Artifacts have immutable fields (ID, CreatedAt) and mutable fields (Content, UpdatedAt). Updates preserve creation metadata.

See `artifact/inmemory.go (Update method)` for reference.

### Filter-Based Listing

List operations support optional filtering by ConversationID and/or TurnID. Empty filters return all artifacts.

See `artifact/inmemory.go (List method)` for reference.

## Rules

### Always (Do These)

| Rule | Rationale |
|------|------------|
| CRITICAL: Validate all request fields before storage operations. Never skip validation for performance. | Invalid data in storage causes cascading failures. Validation catches errors early with clear messages. |
| CRITICAL: Preserve immutable fields (ID, CreatedAt) during Update operations. Only modify Content and UpdatedAt. | Changing IDs breaks references; changing CreatedAt corrupts audit trails. Updates must maintain data integrity. |
| HIGH: Use RWMutex correctly - RLock for reads, Lock for writes. Always defer Unlock immediately after lock. | Incorrect locking causes deadlocks or race conditions. Deferred unlock ensures cleanup even on panic. |
| HIGH: Return ErrNotFound when artifact doesn't exist in Get/Update/Delete. Never return nil error with nil artifact. | Consistent error handling allows callers to distinguish missing vs. other errors. Nil artifact with nil error is ambiguous. |
| MEDIUM: Clone artifacts when returning from storage to prevent external mutation of internal state. | Returning pointers to internal storage allows callers to modify stored data, breaking encapsulation. |
| MEDIUM: Test all validation edge cases with table-driven tests. Include empty strings, nil values, invalid combinations. | Validation is security-critical. Comprehensive tests prevent bypasses and ensure consistent error messages. |

### Never (Don't Do These)

| Rule | Severity | Rationale |
|------|----------|------------|
| CRITICAL: Never allow empty or whitespace-only artifact IDs in any operation. | CRITICAL | Empty IDs break indexing, cause lookup failures, and create ambiguous references. All operations require valid IDs. |
| CRITICAL: Never modify the original artifact pointer passed to Update. Always create new instances. | CRITICAL | Modifying input parameters causes side effects in caller code and violates function purity expectations. |
| HIGH: Never hold locks while performing I/O or expensive operations. Keep critical sections minimal. | HIGH | Long-held locks block all other operations, causing performance degradation and potential deadlocks. |
| HIGH: Never return internal map references or allow direct map access from outside the package. | HIGH | Exposing internal maps breaks encapsulation and allows external code to bypass thread-safety mechanisms. |
| MEDIUM: Never assume artifacts exist without checking. Always handle ErrNotFound gracefully. | MEDIUM | Missing artifacts are normal in distributed systems. Code must handle not-found as expected case, not exceptional. |

### Ask First

- **Adding new fields to Artifact struct** - Changes affect serialization, storage backends, and all consumers. Requires migration strategy for existing data.
- **Implementing persistent storage backend (database, cloud)** - Must maintain interface contract, handle connection pooling, transactions, and migration from in-memory. Requires architecture review.
- **Changing validation rules (e.g., allowing empty IDs)** - Validation rules are security boundaries. Relaxing them may break assumptions in dependent code or create vulnerabilities.
- **Adding caching layer or performance optimizations** - Must not break thread-safety guarantees or introduce cache invalidation bugs. Requires careful concurrency analysis.
- **Modifying error types or error messages** - Callers may depend on specific error types for control flow. Changes can break error handling in dependent packages.

## Commands

### test

Run all artifact package tests including validation and storage tests

```bash
go test ./artifact/...
```

### test-verbose

Run tests with verbose output showing individual test cases

```bash
go test -v ./artifact/...
```

### test-race

Run tests with race detector to catch concurrency issues in InMemory store

```bash
go test -race ./artifact/...
```

### bench

Run benchmarks to measure storage operation performance and memory allocation

```bash
go test -bench=. -benchmem ./artifact/...
```

### coverage

Generate and view test coverage report for artifact package

```bash
go test -coverprofile=coverage.out ./artifact/... && go tool cover -html=coverage.out
```

## Testing

Table-driven tests for validation logic with comprehensive edge cases. Unit tests for InMemory store covering CRUD operations, concurrency, and error conditions. Tests verify both happy paths and error scenarios.

```bash
go test ./artifact/...
go test -race ./artifact/...
go test -v -run TestValidate ./artifact/...
```

Test directory: `artifact/`

## Common Tasks

### Implement custom storage backend

1. 1. Create new file (e.g., postgres.go) in artifact/ package
2. 2. Define struct implementing Store interface from artifact/store.go
3. 3. Implement all 5 methods: Create, Get, List, Update, Delete
4. 4. Use validation functions from artifact/validation.go before storage operations
5. 5. Handle ErrNotFound consistently with inmemory.go pattern
6. 6. Add thread-safety if needed (connection pooling, transactions)
7. 7. Write tests following patterns in request_validation_test.go

### Add new validation rule

1. 1. Identify which operation needs validation (Create, Get, List, Update, Delete)
2. 2. Add validation logic to appropriate validate* function in artifact/validation.go
3. 3. Return descriptive error message explaining what's invalid
4. 4. Add test cases to artifact/request_validation_test.go covering valid and invalid inputs
5. 5. Verify existing tests still pass: go test ./artifact/...
6. 6. Document the new validation requirement in function comments

### Debug artifact not found error

1. 1. Check if artifact ID is correct and not empty/whitespace
2. 2. Verify artifact was created successfully (check Create return value)
3. 3. For List operations, check filter values (ConversationID, TurnID)
4. 4. Add logging in Get/List to see what IDs are being queried
5. 5. Check if using correct Store instance (not mixing multiple stores)
6. 6. See artifact/inmemory.go Get() and List() for reference implementation

### Add new artifact field

1. 1. Add field to Artifact struct in artifact/store.go
2. 2. Update Create method in artifact/inmemory.go to initialize field
3. 3. Update Update method to handle field modification (if mutable)
4. 4. Add validation for new field in artifact/validation.go if needed
5. 5. Update tests in artifact/request_validation_test.go
6. 6. Consider backward compatibility: can old code handle new field?
7. 7. Document field purpose and mutability in struct comments

### Optimize artifact listing performance

1. 1. Profile current List performance: go test -bench=BenchmarkList -cpuprofile=cpu.out
2. 2. Identify bottleneck (lock contention, filtering, copying)
3. 3. For filtering: consider indexing by ConversationID/TurnID (see inmemory.go List)
4. 4. For copying: consider lazy copying or copy-on-write patterns
5. 5. Maintain thread-safety: use RWMutex correctly
6. 6. Add benchmark tests to verify improvement
7. 7. Run race detector: go test -race to ensure no new races

