# AGENTS.md - tool/loadartifactstool

> The loadartifactstool package implements a tool for loading and managing artifacts (files, data) within the ADK-Go agent framework. It provides functionality to retrieve artifacts by ID, handle artifact metadata, and integrate with the broader tool system for agent-based workflows.

## Tech Stack

- **Go** 1.x - Primary implementation language for the load artifacts tool
- **ADK-Go Tool System** internal - Integration with agent tool framework for artifact management

## Key Files

| File | Purpose | Read When |
|------|---------|------------|
| `tool/loadartifactstool/load_artifacts_tool.go` (reference) | Core implementation of the LoadArtifactsTool - defines tool structure, execution logic, and artifact loading mechanisms | Understanding how artifacts are loaded, tool registration, or modifying artifact retrieval logic |
| `tool/loadartifactstool/load_artifacts_tool_test.go` (reference) | Comprehensive test suite covering artifact loading scenarios, error cases, and tool behavior validation | Writing tests, understanding expected behavior, or debugging artifact loading issues |

## Architecture

```
┌─────────────────┐
│  Agent Request  │
│  (artifact_id)  │
└────────┬────────┘
         │
         ▼
┌─────────────────────┐
│ LoadArtifactsTool   │
│  - Validate input   │
│  - Resolve ID       │
└────────┬────────────┘
         │
         ▼
┌─────────────────────┐
│ Artifact Store/FS   │
│  - Fetch artifact   │
│  - Load metadata    │
└────────┬────────────┘
         │
         ▼
┌─────────────────────┐
│  Return Artifact    │
│  (content + meta)   │
└─────────────────────┘
```

## Patterns

### Tool Interface Implementation

Implements the standard ADK-Go tool interface with Name(), Description(), and Execute() methods for integration with agent workflows

See `tool/loadartifactstool/load_artifacts_tool.go` for reference.

### Artifact ID Resolution

Pattern for resolving artifact identifiers to actual artifact objects, handling both direct IDs and reference-based lookups

See `tool/loadartifactstool/load_artifacts_tool.go` for reference.

### Table-Driven Testing

Uses table-driven test pattern with multiple test cases covering success paths, error conditions, and edge cases

See `tool/loadartifactstool/load_artifacts_tool_test.go` for reference.

### Context Propagation

Properly propagates context.Context through tool execution for cancellation, timeouts, and request-scoped values

See `tool/loadartifactstool/load_artifacts_tool.go` for reference.

## Rules

### Always (Do These)

| Rule | Rationale |
|------|------------|
| CRITICAL: Validate artifact IDs before attempting to load - check for empty strings, invalid formats, or malicious paths | Prevents path traversal attacks, nil pointer dereferences, and ensures security when accessing file systems or artifact stores |
| HIGH: Return structured error messages that include the artifact ID and failure reason for debugging | Enables agents and users to understand why artifact loading failed and take corrective action |
| HIGH: Implement the complete tool interface (Name, Description, Execute) for proper tool system integration | Required for the tool to be discoverable and usable by agents in the ADK-Go framework |
| MEDIUM: Include artifact metadata (size, type, creation time) in responses when available | Provides agents with context about artifacts for better decision-making and workflow optimization |
| MEDIUM: Use context.Context for all I/O operations to support cancellation and timeouts | Prevents resource leaks and allows graceful shutdown of long-running artifact load operations |

### Never (Don't Do These)

| Rule | Severity | Rationale |
|------|----------|------------|
| NEVER load artifacts without validating the artifact ID format and permissions | CRITICAL | Could lead to arbitrary file access, path traversal vulnerabilities, or unauthorized data exposure |
| NEVER ignore context cancellation during artifact loading operations | HIGH | Leads to resource leaks, hanging operations, and inability to gracefully shutdown or timeout requests |
| NEVER return raw error messages from underlying storage systems without sanitization | HIGH | May expose internal paths, credentials, or system architecture details to agents or end users |
| NEVER load entire large artifacts into memory without size checks or streaming support | MEDIUM | Can cause out-of-memory errors and denial of service when handling large files or datasets |
| NEVER modify the tool's Name() or Description() after registration with the tool system | MEDIUM | Breaks tool discovery, caching, and agent workflow consistency |

### Ask First

- **Adding caching mechanisms for frequently accessed artifacts** - Requires decisions about cache invalidation strategy, memory limits, and consistency guarantees that affect system-wide behavior
- **Changing artifact ID format or resolution logic** - May break existing agents, workflows, and stored artifact references across the system
- **Adding support for remote artifact stores (S3, GCS, etc.)** - Introduces new dependencies, authentication requirements, and failure modes that need architectural review
- **Implementing artifact versioning or history tracking** - Requires coordination with artifact storage layer and may impact performance and storage requirements

## Commands

### test

Run all tests for the load artifacts tool

```bash
go test ./tool/loadartifactstool/...
```

### test-verbose

Run tests with verbose output to see individual test cases

```bash
go test -v ./tool/loadartifactstool/...
```

### test-coverage

Run tests with coverage reporting

```bash
go test -cover ./tool/loadartifactstool/...
```

### bench

Run benchmarks for artifact loading performance

```bash
go test -bench=. ./tool/loadartifactstool/...
```

## Testing

Table-driven tests covering success cases, error conditions (invalid IDs, missing artifacts, permission errors), and edge cases (empty artifacts, large files). Tests validate tool interface implementation, error message quality, and context handling.

```bash
go test ./tool/loadartifactstool/...
go test -race ./tool/loadartifactstool/...
go test -cover ./tool/loadartifactstool/...
```

Test directory: `tool/loadartifactstool/`

## Common Tasks

### Add support for a new artifact type

1. 1. Review tool/loadartifactstool/load_artifacts_tool.go to understand current artifact handling
2. 2. Add type detection logic in the artifact loading function
3. 3. Implement type-specific parsing or validation if needed
4. 4. Add test cases in load_artifacts_tool_test.go for the new type
5. 5. Update tool Description() to document the new artifact type support
6. 6. Test with go test ./tool/loadartifactstool/...

### Debug artifact loading failures

1. 1. Check test cases in load_artifacts_tool_test.go for similar failure scenarios
2. 2. Verify artifact ID format and validation logic in load_artifacts_tool.go
3. 3. Add logging to trace artifact resolution path
4. 4. Check context cancellation and timeout settings
5. 5. Verify artifact store/filesystem permissions and accessibility
6. 6. Add a regression test case for the specific failure

### Optimize artifact loading performance

1. 1. Add benchmarks in load_artifacts_tool_test.go for current performance baseline
2. 2. Profile with go test -cpuprofile to identify bottlenecks
3. 3. Consider implementing streaming for large artifacts instead of full load
4. 4. Add caching for frequently accessed artifacts (ask first - see rules_ask_first)
5. 5. Implement parallel loading for batch artifact requests if applicable
6. 6. Re-run benchmarks to validate improvements

### Integrate with a new artifact storage backend

1. 1. Review existing artifact store interface in load_artifacts_tool.go
2. 2. Create adapter/wrapper for new storage backend (S3, GCS, database, etc.)
3. 3. Implement artifact ID resolution for new backend's identifier format
4. 4. Add authentication and configuration handling
5. 5. Create mock implementation for testing without real backend access
6. 6. Add comprehensive tests covering new backend in load_artifacts_tool_test.go
7. 7. Document new backend requirements in tool Description()

