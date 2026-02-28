# AGENTS.md - tool/mcptoolset

> The mcptoolset package provides a collection and management system for MCP (Model Context Protocol) tools. It implements a ToolSet that aggregates multiple tools, handles tool registration, lookup, and execution with proper error handling and context propagation.

## Tech Stack

- **Go** 1.x - Core implementation language for tool management and execution
- **MCP Protocol** - Model Context Protocol for tool definitions and execution

## Key Files

| File | Purpose | Read When |
|------|---------|------------|
| `tool/mcptoolset/set.go` (reference) | Core ToolSet implementation - manages collection of tools, registration, and lookup | Understanding how tools are organized and accessed |
| `tool/mcptoolset/tool.go` (reference) | Tool interface and execution logic - defines how individual tools work | Implementing new tools or understanding tool execution flow |
| `tool/mcptoolset/set_test.go` (reference) | Comprehensive test suite showing usage patterns and edge cases | Learning how to use ToolSet API or writing new tests |

## Architecture

```
Tool Registration & Execution Flow:

1. Tool Definition
   └─> Implement Tool interface (Name, Description, InputSchema, Execute)

2. ToolSet Creation
   └─> New() creates empty set
   └─> Add() registers tools
   └─> Merge() combines multiple sets

3. Tool Discovery
   └─> List() returns all available tools
   └─> Get(name) retrieves specific tool

4. Tool Execution
   └─> Execute(ctx, name, input) runs tool
   └─> Context propagation for cancellation
   └─> Error handling with wrapped errors

5. Tool Metadata
   └─> Name uniquely identifies tool
   └─> Description for documentation
   └─> InputSchema defines expected parameters
```

## Patterns

### Tool Interface Pattern

Tools implement a standard interface with Name, Description, InputSchema, and Execute methods. This allows uniform handling of diverse tool types.

See `tool/mcptoolset/tool.go` for reference.

### Immutable Collection Pattern

ToolSet uses internal map for storage but returns copies/slices to prevent external modification. Add() and Merge() create new state rather than mutating.

See `tool/mcptoolset/set.go` for reference.

### Context Propagation Pattern

All Execute methods accept context.Context as first parameter for cancellation, timeouts, and request-scoped values.

See `tool/mcptoolset/tool.go` for reference.

### Error Wrapping Pattern

Errors are wrapped with context using fmt.Errorf with %w verb to maintain error chains for debugging.

See `tool/mcptoolset/set.go` for reference.

### Name-Based Lookup Pattern

Tools are indexed by unique string names. Get() returns tool by name, Execute() takes name as parameter.

See `tool/mcptoolset/set.go` for reference.

## Rules

### Always (Do These)

| Rule | Rationale |
|------|------------|
| CRITICAL: Tools MUST have unique names within a ToolSet - duplicate names will cause the second tool to be ignored | ToolSet uses map[string]Tool internally. Adding duplicate names silently overwrites previous tool, leading to unexpected behavior. |
| HIGH: Always pass context.Context as first parameter to Execute() methods for cancellation support | Enables proper timeout handling, request cancellation, and resource cleanup. Critical for long-running tool operations. |
| HIGH: Wrap errors with context using fmt.Errorf with %w verb to maintain error chains | Preserves original error information while adding context about where/why failure occurred. Essential for debugging. |
| MEDIUM: Return descriptive errors when tools are not found - include the tool name in error message | Makes debugging easier by clearly identifying which tool was requested but missing. |
| MEDIUM: Implement all four Tool interface methods (Name, Description, InputSchema, Execute) when creating new tools | Required for tool to be properly registered and usable. Missing methods will cause compilation errors. |
| MEDIUM: Use List() to get all available tools rather than accessing internal map directly | Maintains encapsulation and prevents external code from depending on internal implementation details. |

### Never (Don't Do These)

| Rule | Severity | Rationale |
|------|----------|------------|
| CRITICAL: Never modify the ToolSet's internal map directly - always use Add() or Merge() | CRITICAL | ToolSet internals are private. Direct modification breaks encapsulation and can cause race conditions if concurrent access occurs. |
| CRITICAL: Never ignore context cancellation in Execute() implementations - check ctx.Err() for long operations | CRITICAL | Ignoring context cancellation can cause resource leaks and prevent graceful shutdown. Tools may continue running after client disconnects. |
| HIGH: Never return nil from Execute() without an error - always return valid result or error | HIGH | Callers expect either valid result or error. Returning (nil, nil) creates ambiguous state that's hard to handle. |
| HIGH: Never use empty strings for tool names - names must be non-empty and unique | HIGH | Empty names break lookup logic and make debugging impossible. Tool identification relies on unique names. |
| MEDIUM: Never assume tool exists without checking - always handle 'tool not found' errors from Get() or Execute() | MEDIUM | Tool availability may change at runtime. Unchecked access causes panics or incorrect behavior. |
| MEDIUM: Never mutate input parameters in Execute() - treat them as read-only | MEDIUM | Input mutation can cause unexpected side effects for callers and makes debugging difficult. Follow functional programming principles. |

### Ask First

- **Adding new methods to the Tool interface** - Changes Tool interface contract, requiring updates to all existing tool implementations. Coordinate with team to ensure backward compatibility.
- **Changing tool name format or naming conventions** - Tool names may be persisted or referenced externally. Changes could break existing integrations or stored configurations.
- **Modifying ToolSet's internal storage mechanism (e.g., from map to different structure)** - Could impact performance characteristics and concurrent access patterns. Needs performance testing and team review.
- **Adding caching or memoization to tool execution** - Changes execution semantics and may cause stale results. Requires careful consideration of cache invalidation strategy.
- **Implementing tool versioning or deprecation mechanisms** - Affects API stability and backward compatibility. Needs architectural discussion about migration paths.

## Commands

### test

Run all tests for the mcptoolset package

```bash
go test ./tool/mcptoolset/...
```

### test-verbose

Run tests with verbose output showing individual test results

```bash
go test -v ./tool/mcptoolset/...
```

### test-coverage

Run tests and show code coverage percentage

```bash
go test -cover ./tool/mcptoolset/...
```

### test-coverage-html

Generate and open HTML coverage report

```bash
go test -coverprofile=coverage.out ./tool/mcptoolset/... && go tool cover -html=coverage.out
```

### bench

Run benchmark tests if any exist

```bash
go test -bench=. ./tool/mcptoolset/...
```

## Testing

Table-driven tests with comprehensive coverage of success cases, error cases, and edge conditions. Tests verify tool registration, lookup, execution, and error handling. See set_test.go for examples of testing ToolSet operations including Add, Get, List, Execute, and Merge.

```bash
go test ./tool/mcptoolset/...
go test -v ./tool/mcptoolset/...
go test -cover ./tool/mcptoolset/...
```

Test directory: `tool/mcptoolset`

## Common Tasks

### Create a new tool

1. 1. Define struct that implements Tool interface from tool/mcptoolset/tool.go
2. 2. Implement Name() method returning unique string identifier
3. 3. Implement Description() method with human-readable description
4. 4. Implement InputSchema() method defining expected input structure
5. 5. Implement Execute(ctx, input) method with actual tool logic
6. 6. Handle context cancellation in Execute if operation is long-running
7. 7. Return proper errors wrapped with context using fmt.Errorf
8. 8. Add tests following patterns in set_test.go

### Register tools in a ToolSet

1. 1. Create new ToolSet: toolset := mcptoolset.New()
2. 2. Add individual tools: toolset.Add(myTool)
3. 3. Or merge multiple sets: toolset.Merge(otherToolSet)
4. 4. Verify registration: tools := toolset.List()
5. 5. Check for specific tool: tool := toolset.Get('tool-name')

### Execute a tool from ToolSet

1. 1. Create context: ctx := context.Background() or context.WithTimeout(...)
2. 2. Prepare input data matching tool's InputSchema
3. 3. Execute: result, err := toolset.Execute(ctx, 'tool-name', input)
4. 4. Check error: if err != nil { handle error }
5. 5. Process result based on tool's output format
6. 6. Handle 'tool not found' errors separately if needed

### Debug tool execution issues

1. 1. Verify tool is registered: tool := toolset.Get('tool-name')
2. 2. Check tool name matches exactly (case-sensitive)
3. 3. Verify input matches InputSchema requirements
4. 4. Add logging in Execute method to trace execution
5. 5. Check error wrapping chain using errors.Is() or errors.As()
6. 6. Test with context.WithTimeout to identify hanging operations
7. 7. Review set_test.go for similar test cases

### Combine multiple tool sources

1. 1. Create separate ToolSets for different tool categories
2. 2. Use Merge() to combine: combined := set1.Merge(set2)
3. 3. Be aware: duplicate names will use last-merged tool
4. 4. Consider prefixing tool names to avoid collisions
5. 5. Use List() to verify all expected tools are present

