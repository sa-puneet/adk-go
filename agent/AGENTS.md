# AGENTS.md - agent

> The agent package provides the core Agent abstraction for building LLM-powered agents with tool execution, context management, and streaming capabilities. It implements a flexible agent architecture that handles conversation state, tool invocation, and model interaction through a clean interface pattern.

## Tech Stack

- **Go** 1.18+ - Core implementation language with generics support for type-safe agent operations
- **context.Context** stdlib - Request-scoped context propagation for cancellation, deadlines, and metadata
- **genkit** internal - Model abstraction layer for LLM interactions and tool definitions

## Key Files

| File | Purpose | Read When |
|------|---------|------------|
| `agent/agent.go` (reference) | Core Agent type definition with Run/Stream methods and configuration options | Understanding agent lifecycle, execution flow, or adding new agent capabilities |
| `agent/context.go` (reference) | Context management for agent state, conversation history, and metadata propagation | Working with agent state, session management, or context-aware operations |
| `agent/agent_test.go` (reference) | Test patterns for agent behavior, mocking, and validation strategies | Writing tests for agent functionality or understanding expected behavior |
| `agent/options.go` | Functional options pattern for agent configuration (likely exists based on patterns) | Configuring agents or adding new configuration options |

## Architecture

```
Agent Request Flow:

[Client] --> Agent.Run(ctx, input)
              |
              v
         [Context Setup]
         - Extract/create agent context
         - Load conversation history
         - Prepare tool registry
              |
              v
         [Model Interaction]
         - Format prompt with history
         - Call LLM with tools
         - Handle streaming if enabled
              |
              v
         [Tool Execution Loop]
         - Parse tool calls from response
         - Execute tools with context
         - Append results to history
         - Continue if more tools needed
              |
              v
         [Response Assembly]
         - Collect final text response
         - Update context state
         - Return result

Context Flow:
context.Context (stdlib) --> AgentContext (metadata) --> Tool Execution --> Model Calls
```

## Patterns

### Functional Options Pattern

Agent configuration uses functional options (WithX functions) for flexible, extensible initialization without breaking changes

See `agent/agent.go` for reference.

### Context-Based State Management

Agent state and conversation history stored in context.Context using typed keys, enabling request-scoped state without global variables

See `agent/context.go` for reference.

### Interface-Based Model Abstraction

Agent depends on model interfaces (not concrete types) allowing any LLM provider to be plugged in

See `agent/agent.go` for reference.

### Streaming Iterator Pattern

Stream method returns iterator for incremental response processing, supporting both streaming and non-streaming models uniformly

See `agent/agent.go` for reference.

### Tool Registry Pattern

Tools registered with agent and made available to model during execution, with automatic marshaling of tool calls and results

See `agent/agent.go` for reference.

## Rules

### Always (Do These)

| Rule | Rationale |
|------|------------|
| CRITICAL: Always propagate context.Context through all agent operations - never create new background contexts | Context carries cancellation signals, deadlines, agent state, and conversation history. Breaking the chain loses all request-scoped data and prevents proper cleanup |
| HIGH: Always use functional options (WithX functions) for agent configuration, never modify Agent struct fields directly | Maintains backward compatibility and allows validation during construction. Direct field access bypasses initialization logic |
| HIGH: Always append tool results and model responses to conversation history in context before next iteration | History continuity is essential for multi-turn conversations and tool execution loops. Missing history causes context loss and repeated tool calls |
| MEDIUM: Use typed context keys (not string keys) for storing agent-specific data in context | Prevents key collisions and provides type safety when retrieving values from context |
| HIGH: Check for context cancellation (ctx.Err()) before expensive operations like model calls or tool execution | Prevents wasted work and ensures responsive cancellation behavior for long-running agent operations |
| MEDIUM: Return structured errors with context about which phase failed (model call, tool execution, parsing) | Enables proper error handling and debugging in multi-step agent workflows |

### Never (Don't Do These)

| Rule | Severity | Rationale |
|------|----------|------------|
| CRITICAL: Never store mutable state in Agent struct fields that changes during execution | CRITICAL | Agent instances may be reused across multiple concurrent requests. Mutable state causes race conditions and data leaks between requests |
| CRITICAL: Never ignore errors from tool execution - always propagate or handle explicitly | CRITICAL | Silent tool failures lead to incorrect agent behavior and confusing responses. Model needs to know when tools fail |
| HIGH: Never modify the input context's values directly - always create new context with additional values | HIGH | Context is immutable by design. Use context.WithValue to create derived contexts |
| HIGH: Never assume conversation history exists in context - always check and initialize if missing | HIGH | First request in a session won't have history. Nil pointer panics occur if not checked |
| MEDIUM: Never block indefinitely in tool execution without respecting context cancellation | MEDIUM | Prevents hung agents and resource leaks when clients disconnect or timeout |
| HIGH: Never return partial results without indicating error state when agent execution fails mid-stream | HIGH | Clients need to know if response is complete or truncated due to error |

### Ask First

- **Adding new fields to Agent struct** - May break existing code or require migration. Consider functional options instead for backward compatibility
- **Changing context key types or values** - Breaks all code that reads from context. Requires coordinated update across codebase
- **Modifying tool execution order or loop termination logic** - Core agent behavior change that affects all users. May cause infinite loops or premature termination
- **Adding automatic retries or error recovery in agent execution** - Can mask real issues and cause unexpected behavior. Should be opt-in via configuration
- **Changing how conversation history is stored or serialized** - Affects session persistence and may break existing stored sessions
- **Modify configuration file: agent/run_config.go** - Configuration changes can affect all environments

## Commands

### test

Run all agent package tests

```bash
go test ./agent/...
```

### test-verbose

Run tests with verbose output to see individual test execution

```bash
go test -v ./agent/...
```

### test-coverage

Generate and view test coverage report

```bash
go test -coverprofile=coverage.out ./agent/... && go tool cover -html=coverage.out
```

### bench

Run benchmarks to measure agent performance

```bash
go test -bench=. -benchmem ./agent/...
```

### race

Run tests with race detector to find concurrency issues

```bash
go test -race ./agent/...
```

## Testing

Table-driven tests with mock models and tools. Tests focus on: (1) context propagation correctness, (2) tool execution loops, (3) error handling paths, (4) streaming behavior, (5) conversation history management. See agent/agent_test.go for patterns.

```bash
go test ./agent/...
go test -race ./agent/...
go test -v -run TestAgentRun ./agent/...
```

Test directory: `agent/`

## Common Tasks

### Create a new agent with tools

1. Import agent package and model interface
2. Define tools using tool definition format (see tool package)
3. Create agent: agent.New(model, agent.WithTools(tools...))
4. Call agent.Run(ctx, input) or agent.Stream(ctx, input)
5. Handle response or iterate over stream

### Add conversation history to agent context

1. Retrieve existing history from storage/session
2. Create context with history: ctx = agent.WithHistory(ctx, history)
3. Pass context to agent.Run or agent.Stream
4. Extract updated history from returned context for persistence

### Implement a custom tool for agent

1. Define tool function with signature: func(ctx context.Context, input ToolInput) (ToolOutput, error)
2. Create tool definition with name, description, and schema
3. Register tool with agent using WithTools option
4. Tool will be automatically called by agent when model requests it
5. See agent/agent_test.go for mock tool examples

### Handle streaming agent responses

1. Call iter := agent.Stream(ctx, input)
2. Loop: for chunk := range iter { process(chunk) }
3. Check iter.Err() after loop completes to detect errors
4. Accumulate chunks to build complete response if needed
5. Handle partial results gracefully if error occurs mid-stream

### Debug agent execution flow

1. Enable verbose logging if available in model implementation
2. Add logging in tool implementations to trace execution
3. Inspect context values to verify state propagation
4. Check conversation history in context after each turn
5. Use -v flag with tests to see detailed execution: go test -v ./agent/...
6. Run with race detector if suspecting concurrency issues: go test -race ./agent/...

### Add agent configuration option

1. Define option function: func WithX(value T) AgentOption
2. Option should return func(*Agent) that modifies agent config
3. Add field to Agent struct if needed (prefer avoiding mutable state)
4. Document option behavior and when to use it
5. Add test case in agent_test.go validating option effect
6. See existing WithTools, WithX patterns in agent/agent.go

