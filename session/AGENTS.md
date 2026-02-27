# AGENTS.md - session

> The session package provides session management for AI agents, implementing in-memory storage for conversation history, state persistence, and message tracking. It handles session lifecycle, memory management with configurable limits, and provides thread-safe operations for concurrent access.

## Tech Stack

- **Go** 1.x - Core implementation language for session management
- **sync.RWMutex** stdlib - Thread-safe concurrent access to session data
- **context** stdlib - Cancellation and timeout handling for session operations

## Key Files

| File | Purpose | Read When |
|------|---------|------------|
| `session/inmemory.go` (reference) | Core in-memory session store implementation with thread-safe operations | Understanding session storage, memory management, or implementing new session backends |
| `session/inmemory_test.go` (reference) | Comprehensive test suite demonstrating session operations, edge cases, and concurrency patterns | Writing tests for session functionality or understanding expected behavior |
| `session/session.go` (reference) | Session interface definitions and core types | Understanding session contracts or implementing custom session stores |
| `session/manager.go` | Session lifecycle management and coordination | Implementing session creation, retrieval, or cleanup logic |

## Architecture

```
Session Flow:

[Agent Request] → [Session Manager]
                       |
                       v
              [Session Store Interface]
                       |
                       v
              [InMemory Implementation]
                       |
    +------------------+------------------+
    |                  |                  |
    v                  v                  v
[Session Data]  [Message History]  [State Storage]
    |                  |                  |
    v                  v                  v
[RWMutex Lock]  [Memory Limits]   [Cleanup/GC]

Data Flow:
1. Agent creates/retrieves session via Manager
2. Manager delegates to Store (InMemory)
3. Store locks for thread-safety
4. Operations: Add messages, update state, retrieve history
5. Memory limits enforced (max messages, max size)
6. Cleanup on session close or timeout
```

## Patterns

### Thread-Safe Session Storage

Uses sync.RWMutex for concurrent read/write access. Multiple readers allowed, exclusive writer lock. Always lock before accessing session data maps.

See `session/inmemory.go` for reference.

### Memory-Bounded Storage

Implements configurable limits (max messages, max total size) to prevent unbounded memory growth. Automatically evicts oldest messages when limits exceeded.

See `session/inmemory.go` for reference.

### Context-Aware Operations

All session operations accept context.Context for cancellation and timeout handling. Check context before long operations.

See `session/inmemory.go` for reference.

### Interface-Based Design

Session storage defined by interfaces, allowing multiple implementations (in-memory, Redis, database). InMemory is reference implementation.

See `session/session.go` for reference.

### Immutable Session IDs

Session IDs are strings, treated as immutable keys. Never modify session ID after creation. Use as map keys throughout.

See `session/inmemory.go` for reference.

### Graceful Degradation

When memory limits reached, oldest messages evicted rather than failing. Logs warnings but continues operation.

See `session/inmemory.go` for reference.

## Rules

### Always (Do These)

| Rule | Rationale |
|------|------------|
| CRITICAL: Always acquire appropriate lock (RLock for reads, Lock for writes) before accessing session data structures | Session stores are accessed concurrently by multiple goroutines. Race conditions will corrupt session state and cause data loss or panics. |
| CRITICAL: Always defer unlock immediately after acquiring lock to prevent deadlocks | Early returns or panics without unlock will deadlock all other goroutines waiting for the lock. |
| HIGH: Check context.Err() before performing expensive operations or long-running tasks | Respects cancellation and timeouts, prevents wasted work on cancelled requests. |
| HIGH: Enforce memory limits when adding messages or state to prevent unbounded growth | In-memory sessions can exhaust server memory if not bounded. Must evict old data when limits reached. |
| MEDIUM: Return ErrSessionNotFound when session ID doesn't exist, not generic errors | Callers need to distinguish between missing sessions and other errors for proper handling (create vs retry). |
| MEDIUM: Copy session data when returning to prevent external mutation of internal state | Returning pointers to internal data allows external code to bypass locks and corrupt state. |
| MEDIUM: Initialize session with empty collections (messages, state) never nil | Prevents nil pointer panics when accessing session fields. Empty collections are safer than nil. |

### Never (Don't Do These)

| Rule | Severity | Rationale |
|------|----------|------------|
| NEVER access session data structures without holding the appropriate lock | CRITICAL | Causes race conditions, data corruption, and crashes. Go race detector will catch this in tests. |
| NEVER modify session ID after creation or use mutable types as session IDs | CRITICAL | Session IDs are map keys. Modifying them breaks lookups and causes sessions to become unreachable. |
| NEVER return pointers to internal session data structures without copying | CRITICAL | Allows external code to modify internal state without locks, causing race conditions and corruption. |
| NEVER ignore context cancellation in long-running or blocking operations | HIGH | Causes goroutine leaks and wasted resources. Operations continue after client disconnects. |
| NEVER allow unbounded growth of session data (messages, state, metadata) | HIGH | Will exhaust server memory and cause OOM crashes. Must implement eviction or limits. |
| NEVER assume session exists without checking - always handle ErrSessionNotFound | HIGH | Sessions can expire, be deleted, or never exist. Missing checks cause nil pointer panics. |
| NEVER use global state or package-level variables for session storage | MEDIUM | Makes testing difficult, prevents multiple instances, and creates hidden dependencies. |
| NEVER log sensitive session data (user messages, state) at INFO level | MEDIUM | Session data may contain PII or secrets. Only log at DEBUG with explicit opt-in. |

### Ask First

- **Changing session storage backend from in-memory to persistent (Redis, database)** - Affects deployment architecture, requires new dependencies, changes failure modes and performance characteristics. Coordinate with ops team.
- **Modifying memory limit defaults or eviction policies** - Impacts all agents using sessions. May cause unexpected message loss or memory issues. Needs capacity planning and testing.
- **Adding new session metadata fields or changing session structure** - May break existing agents or require migration. Affects serialization if persistence added later.
- **Implementing session sharing or multi-agent access to same session** - Complex concurrency implications. May require distributed locking or conflict resolution strategies.
- **Adding session expiration or TTL mechanisms** - Affects session lifecycle and cleanup. May surprise users if sessions disappear unexpectedly.

## Commands

### test-session

Run all session package tests including concurrency tests

```bash
go test ./session/...
```

### test-session-race

Run tests with race detector to catch concurrency bugs (CRITICAL before committing)

```bash
go test -race ./session/...
```

### test-session-verbose

Run tests with verbose output to see individual test results

```bash
go test -v ./session/...
```

### bench-session

Run benchmarks to measure session operation performance and memory usage

```bash
go test -bench=. -benchmem ./session/...
```

### test-session-coverage

Generate and view test coverage report

```bash
go test -coverprofile=coverage.out ./session/... && go tool cover -html=coverage.out
```

## Testing

Table-driven tests with comprehensive edge cases. Heavy focus on concurrency testing with race detector. Tests cover: session lifecycle, concurrent access, memory limits, context cancellation, error conditions, and boundary cases.

```bash
go test ./session/...
go test -race ./session/...
go test -v -run TestInMemory ./session/...
```

Test directory: `session/`

## Common Tasks

### Add new session operation

1. 1. Define method signature in session interface (session/session.go)
2. 2. Implement in InMemoryStore with proper locking pattern: Lock/RLock → defer Unlock → check context → perform operation
3. 3. Add memory limit checks if operation adds data
4. 4. Return appropriate errors (ErrSessionNotFound, etc.)
5. 5. Write table-driven tests in session/inmemory_test.go
6. 6. Run with race detector: go test -race ./session/...
7. 7. Add benchmark if performance-critical

### Implement new session backend

1. 1. Study session/session.go interface definitions
2. 2. Review session/inmemory.go as reference implementation
3. 3. Implement all interface methods with proper error handling
4. 4. Ensure thread-safety appropriate for backend (locks, transactions, etc.)
5. 5. Handle context cancellation in all operations
6. 6. Copy test suite from inmemory_test.go and adapt
7. 7. Add backend-specific tests (connection failures, etc.)
8. 8. Document configuration and deployment requirements

### Debug session concurrency issue

1. 1. Run tests with race detector: go test -race ./session/...
2. 2. Check all session data access is lock-protected
3. 3. Verify defer unlock immediately follows lock acquisition
4. 4. Look for returned pointers to internal data (should copy instead)
5. 5. Check for lock held across external calls (can deadlock)
6. 6. Add targeted concurrency test reproducing the issue
7. 7. Use go test -race -count=100 for intermittent issues

### Optimize session memory usage

1. 1. Run benchmarks: go test -bench=. -benchmem ./session/...
2. 2. Profile with pprof: go test -memprofile=mem.out -bench=.
3. 3. Check message size calculations in checkMemoryLimits()
4. 4. Review eviction policy - oldest messages removed first
5. 5. Consider adjusting default limits based on usage patterns
6. 6. Test with realistic message sizes and counts
7. 7. Verify no memory leaks with long-running sessions

### Add session persistence

1. 1. Design serialization format (JSON, protobuf, etc.)
2. 2. Implement Save/Load methods in new backend
3. 3. Handle partial failures and corruption gracefully
4. 4. Add migration path from in-memory to persistent
5. 5. Test recovery from crashes and restarts
6. 6. Document backup and restore procedures
7. 7. Consider encryption for sensitive session data

