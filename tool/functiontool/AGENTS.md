# AGENTS.md - tool/functiontool

> The functiontool package provides a framework for converting Go functions into AI-callable tools with automatic JSON schema generation and parameter validation. It handles type reflection, schema generation, and execution of functions with proper error handling and context management.

## Tech Stack

- **Go** 1.21+ - Core implementation language with reflection for dynamic function handling
- **encoding/json** stdlib - JSON schema generation and parameter marshaling/unmarshaling
- **reflect** stdlib - Runtime type inspection and function invocation

## Key Files

| File | Purpose | Read When |
|------|---------|------------|
| `tool/functiontool/function.go` (reference) | Core FunctionTool implementation with New() constructor and Run() execution | Understanding how Go functions are wrapped as AI tools |
| `tool/functiontool/schema.go` (reference) | JSON schema generation from Go types using reflection | Adding support for new parameter types or modifying schema generation |
| `tool/functiontool/function_test.go` (reference) | Comprehensive test suite showing usage patterns for various function signatures | Learning how to create function tools or troubleshooting type handling |
| `tool/functiontool/validation.go` | Parameter validation logic before function execution | Understanding input validation or adding custom validation rules |
| `tool/functiontool/types.go` | Type definitions and interfaces for function tool system | Understanding the type system or implementing custom types |

## Architecture

```
Function Registration Flow:
1. User provides Go function → New(name, desc, fn)
2. Reflection extracts signature → validateFunction()
3. Generate JSON schema → generateSchema()
4. Store metadata → FunctionTool struct

Execution Flow:
1. AI calls tool with JSON params → Run(ctx, input)
2. Unmarshal to map[string]any → json.Unmarshal
3. Convert to Go types → convertParameters()
4. Validate parameters → validateParams()
5. Call function via reflection → fn.Call()
6. Return result or error → marshal response
```

## Patterns

### Function Signature Validation

Functions must follow specific patterns: (ctx, params) → (result, error) or (params) → (result, error). First param can be context.Context, last return must be error.

See `tool/functiontool/function.go:validateFunction()` for reference.

### Struct Tag Schema Generation

Use struct tags `json:"name"` for field names and `description:"text"` for schema documentation. Required fields use `required:"true"`.

See `tool/functiontool/schema.go:generateStructSchema()` for reference.

### Type Reflection and Conversion

Automatic conversion between JSON types and Go types using reflect package. Handles primitives, structs, slices, maps, and pointers.

See `tool/functiontool/types.go:convertValue()` for reference.

### Error Wrapping

All errors are wrapped with context using fmt.Errorf with %w verb to maintain error chains for debugging.

See `tool/functiontool/function.go:Run()` for reference.

### Optional Context Parameter

Functions can optionally accept context.Context as first parameter for cancellation and timeout support.

See `tool/functiontool/function_test.go:TestFunctionWithContext` for reference.

## Rules

### Always (Do These)

| Rule | Rationale |
|------|------------|
| CRITICAL: Function signatures must return error as last return value | The reflection-based execution system expects error as final return to distinguish success/failure. Without it, panics or silent failures occur. |
| HIGH: Use struct tags for parameter documentation: `json:"name" description:"purpose" required:"true"` | Schema generation relies on struct tags to create accurate JSON schemas for AI models. Missing tags result in poor parameter descriptions. |
| HIGH: Validate function signature with New() before registration | New() performs compile-time-like checks at runtime. Invalid signatures cause runtime panics during execution. |
| MEDIUM: Export struct fields that should be accessible as tool parameters | Reflection can only access exported fields. Unexported fields are silently ignored in schema generation. |
| MEDIUM: Handle nil pointers in parameter structs explicitly | JSON unmarshaling may produce nil pointers for optional fields. Dereference checks prevent panics. |
| HIGH: Pass context.Context as first parameter for long-running functions | Enables cancellation and timeout control. AI agents may need to abort operations. |

### Never (Don't Do These)

| Rule | Severity | Rationale |
|------|----------|------------|
| CRITICAL: Never create functions without error return value | CRITICAL | System requires error as last return for proper error handling. Functions without error return will fail validation. |
| CRITICAL: Never modify function signature after calling New() | CRITICAL | Schema is generated once at creation. Runtime signature changes cause type mismatches and panics. |
| HIGH: Never use unexported struct fields for tool parameters | HIGH | Reflection cannot access unexported fields. They will be silently ignored, causing missing parameter errors. |
| HIGH: Never return multiple non-error values without wrapping in struct | HIGH | Tool system expects single result value (plus error). Multiple returns need struct wrapper for proper serialization. |
| MEDIUM: Never use channels, funcs, or unsafe types as parameters | MEDIUM | These types cannot be serialized to/from JSON. Schema generation will fail or produce invalid schemas. |
| HIGH: Never ignore context cancellation in long-running functions | HIGH | AI agents need ability to cancel operations. Ignoring context leads to resource leaks and hung operations. |

### Ask First

- **Adding support for new Go types in schema generation** - Type system is carefully designed for JSON compatibility. New types may break serialization or require special handling in multiple places (schema.go, types.go, validation.go)
- **Changing function signature validation rules** - Validation rules ensure runtime safety. Relaxing them could introduce panics or undefined behavior in production
- **Modifying error handling or wrapping patterns** - Error chains are used for debugging across the system. Changes affect observability and error reporting to AI agents
- **Adding custom struct tag support** - Tag parsing is centralized in schema generation. New tags need consistent handling and documentation

## Commands

### test

Run all function tool tests including type conversion and schema generation

```bash
go test ./tool/functiontool/...
```

### test-verbose

Run tests with verbose output to see individual test cases

```bash
go test -v ./tool/functiontool/...
```

### test-coverage

Run tests with coverage report

```bash
go test -cover ./tool/functiontool/...
```

### bench

Run benchmarks for reflection and schema generation performance

```bash
go test -bench=. ./tool/functiontool/...
```

## Testing

Table-driven tests with comprehensive type coverage. Tests validate function signatures, schema generation, parameter conversion, and execution for all supported Go types.

```bash
go test ./tool/functiontool/...
go test -race ./tool/functiontool/...
```

Test directory: `tool/functiontool/`

## Common Tasks

### Create a new function tool

1. 1. Define Go function with signature: func(ctx context.Context, params ParamStruct) (ResultStruct, error)
2. 2. Add struct tags to ParamStruct: `json:"field_name" description:"Field purpose" required:"true"`
3. 3. Call functiontool.New("tool_name", "Tool description", yourFunction)
4. 4. Register returned FunctionTool with agent's tool registry
5. See tool/functiontool/function_test.go:TestNew for examples

### Add parameter validation

1. 1. Implement validation logic inside your function (not in functiontool)
2. 2. Return descriptive errors for invalid parameters
3. 3. Use struct tags to mark required fields: `required:"true"`
4. 4. Consider adding custom validation in tool/functiontool/validation.go for reusable rules
5. See tool/functiontool/validation.go for validation patterns

### Support new parameter type

1. 1. Check if type is JSON-serializable (primitives, structs, slices, maps)
2. 2. Add type handling in tool/functiontool/schema.go:generateSchema()
3. 3. Add conversion logic in tool/functiontool/types.go:convertValue()
4. 4. Add test cases in function_test.go with new type
5. 5. Update documentation with supported type
6. WARNING: Ask before implementing - affects core type system

### Debug schema generation issues

1. 1. Check struct field tags: json, description, required
2. 2. Verify fields are exported (capitalized)
3. 3. Run test: go test -v ./tool/functiontool/ -run TestSchemaGeneration
4. 4. Compare generated schema with expected in test output
5. 5. Check tool/functiontool/schema.go:generateStructSchema() for type handling
6. See tool/functiontool/function_test.go for schema examples

### Handle context cancellation

1. 1. Add context.Context as first parameter: func(ctx context.Context, params T) (R, error)
2. 2. Check ctx.Done() in loops or before long operations
3. 3. Return ctx.Err() when cancelled
4. 4. Test with context.WithTimeout() or context.WithCancel()
5. See tool/functiontool/function_test.go:TestFunctionWithContext

