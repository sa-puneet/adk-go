# AGENTS.md - tool/agenttool

> The agenttool package provides the core tool abstraction for AI agents, enabling function calling and tool execution. It defines interfaces and implementations for converting Go functions into agent-callable tools with schema generation, parameter validation, and execution handling.

## Tech Stack

- **Go** 1.x - Core implementation language for tool abstraction and execution
- **genai** - Google AI SDK integration for tool schema and function declarations

## Key Files

| File | Purpose | Read When |
|------|---------|------------|
| `tool/agenttool/agent_tool.go` (reference) | Core tool interface and implementation - defines Tool interface, FromFunc factory, and execution logic | Understanding how to create tools from Go functions or implement custom tools |
| `tool/agenttool/agent_tool_test.go` (reference) | Comprehensive test suite showing tool creation patterns, parameter handling, and error cases | Learning how to properly test tools or understand expected behavior |

## Architecture

```
Tool Creation Flow:
1. Go Function → FromFunc(name, desc, fn) → Tool interface
2. Tool.Schema() → genai.Tool (JSON schema generation)
3. Agent receives function call → Tool.Execute(ctx, params)
4. Parameter validation → Type conversion → Function invocation
5. Result/Error → Structured response

Key Components:
- Tool interface: Name(), Description(), Schema(), Execute()
- funcTool: Concrete implementation wrapping Go functions
- Schema generation: Automatic parameter introspection
- Execution: Context-aware invocation with error handling
```

## Patterns

### Function-to-Tool Conversion

Use FromFunc to wrap any Go function as an agent tool. Function signature determines parameter schema automatically. Supports context.Context as first param, various return types (value, error, or both).

See `tool/agenttool/agent_tool.go - FromFunc implementation` for reference.

### Schema Generation

Tools automatically generate genai.Tool schemas with parameter types, descriptions, and requirements. Uses reflection to introspect function signatures and build JSON schemas.

See `tool/agenttool/agent_tool.go - Schema() method` for reference.

### Context-Aware Execution

Tools accept context.Context for cancellation and timeout handling. First parameter can be context.Context (optional), followed by actual tool parameters.

See `tool/agenttool/agent_tool_test.go - TestExecute with context` for reference.

### Type-Safe Parameter Handling

Parameters passed as map[string]any are validated and converted to function parameter types. Supports primitives, structs, and nested types with proper error reporting.

See `tool/agenttool/agent_tool_test.go - parameter validation tests` for reference.

## Rules

### Always (Do These)

| Rule | Rationale |
|------|------------|
| CRITICAL: Use FromFunc to create tools from Go functions - this is the primary factory method | FromFunc handles schema generation, parameter validation, and execution automatically. Direct Tool interface implementation should be rare. |
| HIGH: Include context.Context as first parameter in tool functions when cancellation/timeout is needed | Context enables proper cancellation propagation and timeout handling during tool execution. The framework automatically detects and passes context. |
| HIGH: Return (result, error) or just error from tool functions for proper error handling | Standard Go error handling pattern. Framework distinguishes between successful results and errors, enabling proper agent error recovery. |
| MEDIUM: Provide clear, descriptive names and descriptions when creating tools | Tool names and descriptions are exposed to LLMs for function selection. Clear descriptions improve agent decision-making. |
| MEDIUM: Use struct types for complex tool parameters instead of many individual parameters | Structs provide better organization, validation, and schema generation for complex parameter sets. |

### Never (Don't Do These)

| Rule | Severity | Rationale |
|------|----------|------------|
| CRITICAL: Never ignore errors returned from Tool.Execute() | CRITICAL | Execute errors indicate tool failures that must be handled. Ignoring them leads to silent failures and incorrect agent behavior. |
| HIGH: Never pass nil context to Execute() - use context.Background() at minimum | HIGH | Nil context causes panics in context-aware tools. Always provide a valid context, even if it's just Background(). |
| HIGH: Never create tools with functions that have unsupported parameter types without validation | HIGH | Unsupported types (channels, functions, etc.) cannot be serialized to JSON schemas and will cause runtime errors. |
| MEDIUM: Never modify tool schemas after creation | MEDIUM | Tool schemas are generated once and should be immutable. Modifying them breaks the contract with agents. |

### Ask First

- **Implementing custom Tool interface instead of using FromFunc** - FromFunc handles 99% of use cases. Custom implementations should only be needed for very specialized behavior like streaming or stateful tools.
- **Adding new parameter types to schema generation** - Schema generation must align with genai.Tool expectations. New types may require coordination with the broader SDK.
- **Changing tool execution semantics (e.g., async execution, retries)** - Execution model affects agent behavior and error handling. Changes should be coordinated with agent implementation patterns.

## Commands

### test

Run all agenttool tests including tool creation, execution, and error handling

```bash
go test ./tool/agenttool/...
```

### test-verbose

Run tests with verbose output to see detailed test execution

```bash
go test -v ./tool/agenttool/...
```

### test-coverage

Run tests with coverage reporting

```bash
go test -cover ./tool/agenttool/...
```

### benchmark

Run benchmarks if any exist for tool execution performance

```bash
go test -bench=. ./tool/agenttool/...
```

## Testing

Table-driven tests with comprehensive coverage of tool creation, parameter validation, execution paths, and error cases. Tests verify schema generation, type conversion, context handling, and error propagation.

```bash
go test ./tool/agenttool/...
go test -race ./tool/agenttool/...
```

Test directory: `tool/agenttool/`

## Common Tasks

### Create a simple tool from a Go function

1. Define a Go function with clear parameters and return types
2. Use agenttool.FromFunc(name, description, function) to create the tool
3. The function can return (result, error), just result, or just error
4. See tool/agenttool/agent_tool_test.go for examples of various function signatures

### Create a tool with context support

1. Add context.Context as the first parameter of your function
2. Implement context cancellation checks in long-running operations
3. Use agenttool.FromFunc to wrap the function - context is automatically handled
4. See tool/agenttool/agent_tool_test.go - TestExecute for context examples

### Create a tool with complex parameters

1. Define a struct type with fields for your parameters
2. Add json tags to struct fields for proper schema generation
3. Create function accepting the struct as parameter
4. Use FromFunc to create the tool - struct fields become schema properties
5. See tool/agenttool/agent_tool_test.go for struct parameter examples

### Execute a tool with parameters

1. Create a context (context.Background() or with timeout/cancellation)
2. Prepare parameters as map[string]any matching tool schema
3. Call tool.Execute(ctx, params)
4. Check returned error for execution failures
5. Handle the result based on your tool's return type
6. See tool/agenttool/agent_tool_test.go - TestExecute for execution patterns

### Get tool schema for agent registration

1. Create your tool using FromFunc
2. Call tool.Schema() to get genai.Tool schema
3. Pass schema to agent configuration or model
4. Schema includes function name, description, and parameter definitions
5. See tool/agenttool/agent_tool.go - Schema() method

### Debug tool parameter validation issues

1. Check that parameter names in map match function parameter names
2. Verify parameter types are compatible (string, int, float64, bool, structs)
3. Ensure required parameters are provided
4. Look at error messages from Execute() - they indicate specific validation failures
5. See tool/agenttool/agent_tool_test.go for examples of valid parameter patterns

