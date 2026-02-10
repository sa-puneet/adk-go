# AGENTS.md - tool/mcptoolset

> The mcptoolset package provides a registry and management system for MCP (Model Context Protocol) tools. It handles tool registration, validation, execution, and schema generation for tools that can be invoked by AI agents, with support for both synchronous and asynchronous execution patterns.

## Tech Stack

- **Go** 1.x - Core implementation language for tool registry and execution
- **MCP Protocol** - Model Context Protocol for tool definitions and schemas
- **JSON Schema** - Tool input validation and schema generation

## Key Files

| File | Purpose | Read When |
|------|---------|------------|
| `tool/mcptoolset/set.go` (reference) | Core ToolSet implementation - registry for managing multiple tools with Add/Get/List/Execute operations | Understanding tool registration, lookup, and execution flow |
| `tool/mcptoolset/tool.go` (reference) | Tool wrapper and execution logic - handles individual tool invocation, schema generation, and error handling | Understanding how individual tools are wrapped and executed |
| `tool/mcptoolset/set_test.go` (reference) | Comprehensive test suite showing usage patterns, error cases, and expected behaviors | Learning correct usage patterns and edge cases |

## Architecture

```
Tool Registration & Execution Flow:

1. REGISTRATION:
   ToolSet.Add(tool) → validates → stores in registry
   ↓
   Tool wrapped with metadata (name, description, schema)

2. DISCOVERY:
   ToolSet.List() → returns all registered tool definitions
   ToolSet.Get(name) → retrieves specific tool

3. EXECUTION:
   ToolSet.Execute(ctx, name, args) → 
   ↓
   Lookup tool by name
   ↓
   Validate arguments against schema
   ↓
   Invoke tool.Call(ctx, args)
   ↓
   Return result or error

4. SCHEMA GENERATION:
   Tool → InputSchema() → JSON Schema for validation
```

## Patterns

### Tool Registry Pattern

Central registry (ToolSet) manages multiple tools with name-based lookup. Tools are registered once and can be executed multiple times. Thread-safe with mutex protection.

See `tool/mcptoolset/set.go - see ToolSet struct and Add/Get methods` for reference.

### Tool Wrapper Pattern

Raw tool implementations are wrapped in a Tool struct that adds metadata, schema generation, and standardized execution interface. Separates tool logic from protocol concerns.

See `tool/mcptoolset/tool.go - see Tool struct and NewTool function` for reference.

### Context-Based Execution

All tool executions accept context.Context as first parameter for cancellation, timeouts, and request-scoped values. Critical for managing long-running operations.

See `tool/mcptoolset/tool.go - Call method signature` for reference.

### Schema-Driven Validation

Tools define input schemas that are used for validation before execution. Schemas are generated from tool definitions and exposed via InputSchema() method.

See `tool/mcptoolset/tool.go - InputSchema method` for reference.

### Error Wrapping

Errors are wrapped with context using fmt.Errorf with %w verb to maintain error chains. Provides clear error messages with operation context.

See `tool/mcptoolset/set.go - Execute method error handling` for reference.

## Rules

### Always (Do These)

| Rule | Rationale |
|------|------------|
| CRITICAL: Register tools with unique names - duplicate names will cause registration to fail | ToolSet enforces unique tool names to prevent ambiguity during execution. Attempting to add a tool with an existing name returns an error. |
| CRITICAL: Pass context.Context as first parameter to all Execute calls | Context enables cancellation, timeouts, and proper resource cleanup. Missing context prevents graceful shutdown of long-running operations. |
| HIGH: Check errors from ToolSet.Add() during registration | Add can fail due to duplicate names or invalid tool definitions. Ignoring these errors leads to missing tools at runtime. |
| HIGH: Use ToolSet.Get() to check tool existence before execution | Execute will fail if tool doesn't exist. Pre-checking with Get allows better error handling and user feedback. |
| MEDIUM: Implement proper JSON marshaling for tool arguments | Tool arguments are passed as JSON. Ensure your argument types can be properly marshaled/unmarshaled. |
| MEDIUM: Return structured data from tool implementations | Tool results should be JSON-serializable. Complex types should implement proper marshaling. |

### Never (Don't Do These)

| Rule | Severity | Rationale |
|------|----------|------------|
| CRITICAL: Never register tools with nil or empty names | CRITICAL | Tool names are used as unique identifiers for lookup and execution. Empty names break the registry system. |
| CRITICAL: Never ignore context cancellation in tool implementations | CRITICAL | Ignoring context cancellation can cause resource leaks and prevent graceful shutdown. Always check ctx.Done() in long-running operations. |
| HIGH: Never modify ToolSet registry after concurrent access begins | HIGH | While ToolSet uses mutex for thread safety, adding tools during execution can cause race conditions. Register all tools during initialization. |
| HIGH: Never assume tool execution is synchronous without checking implementation | HIGH | Tools may perform I/O, network calls, or other async operations. Always handle potential delays and use context timeouts. |
| MEDIUM: Never return raw errors from tool implementations without context | MEDIUM | Wrap errors with operation context using fmt.Errorf to provide clear error messages for debugging. |

### Ask First

- **Adding new tool execution modes (async, streaming, etc.)** - Changes to execution model affect all tool implementations and may require protocol changes
- **Modifying tool registration validation logic** - Validation changes can break existing tool registrations and affect backward compatibility
- **Changing tool schema generation approach** - Schema changes affect client validation and may break existing integrations
- **Adding global tool interceptors or middleware** - Interceptors affect all tool executions and need careful design for performance and error handling

## Commands

### test

Run all tests for the mcptoolset package

```bash
go test ./tool/mcptoolset/...
```

### test-verbose

Run tests with verbose output to see individual test cases

```bash
go test -v ./tool/mcptoolset/...
```

### test-coverage

Run tests with coverage reporting

```bash
go test -cover ./tool/mcptoolset/...
```

### bench

Run benchmarks if any exist

```bash
go test -bench=. ./tool/mcptoolset/...
```

## Testing

Table-driven tests with comprehensive coverage of success paths, error cases, and edge conditions. Tests validate registration, lookup, execution, and error handling.

```bash
go test ./tool/mcptoolset/...
go test -v ./tool/mcptoolset/... -run TestToolSet_Add
go test -v ./tool/mcptoolset/... -run TestToolSet_Execute
```

Test directory: `tool/mcptoolset/`

## Common Tasks

### Register a new tool

1. Create a ToolSet: ts := mcptoolset.NewToolSet()
2. Define your tool implementation with Call method
3. Register: err := ts.Add(yourTool)
4. Check error for duplicate names or validation failures
5. See tool/mcptoolset/set_test.go for examples

### Execute a registered tool

1. Get context: ctx := context.Background() or context.WithTimeout(...)
2. Prepare arguments as map[string]any or JSON-serializable struct
3. Execute: result, err := ts.Execute(ctx, "toolName", args)
4. Handle errors (tool not found, execution failure, context cancellation)
5. See tool/mcptoolset/set_test.go TestToolSet_Execute for patterns

### List available tools

1. Call ts.List() to get all registered tool definitions
2. Each definition includes name, description, and input schema
3. Use for discovery, documentation, or UI generation
4. See tool/mcptoolset/set.go List method

### Implement a custom tool

1. Create struct implementing the tool interface
2. Implement Name() string method returning unique identifier
3. Implement Description() string method for documentation
4. Implement InputSchema() method returning JSON schema
5. Implement Call(ctx context.Context, args any) (any, error) for execution logic
6. Always respect context cancellation in Call implementation
7. See tool/mcptoolset/tool.go for Tool wrapper structure

### Handle tool execution errors

1. Check if error is tool not found: tool := ts.Get(name); if tool == nil {...}
2. Check if error is context cancellation: errors.Is(err, context.Canceled)
3. Wrap errors with context: fmt.Errorf("operation failed: %w", err)
4. Log errors with tool name and arguments for debugging
5. See tool/mcptoolset/set.go Execute method for error handling patterns

