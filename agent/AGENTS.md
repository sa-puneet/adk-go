# AGENTS.md - agent

> The agent package provides the core Agent abstraction for building LLM-powered agents with tool execution, context management, and streaming capabilities. It implements a flexible agent architecture that handles LLM interactions, tool calling, state management, and error handling through a clean interface pattern.

## Tech Stack

- **Go** 1.21+ - Core implementation language with generics support for type-safe agent operations
- **context.Context** stdlib - Request-scoped context propagation, cancellation, and metadata storage
- **genkit** internal - LLM model abstraction, tool definitions, and streaming interfaces

## Key Files

| File | Purpose | Read When |
|------|---------|------------|
| `agent/agent.go` (reference) | Core Agent interface and implementation with Run/Stream methods, tool execution, and state management | Understanding agent lifecycle, tool calling patterns, or implementing new agent types |
| `agent/context.go` (reference) | Context key management for agent metadata, state storage, and request-scoped data propagation | Working with agent context, state management, or custom middleware |
| `agent/agent_test.go` (reference) | Comprehensive test suite showing agent usage patterns, tool execution, and error handling | Learning how to test agents, mock tools, or validate agent behavior |
| `agent/options.go` | Agent configuration options including model selection, tools, system prompts, and callbacks | Configuring agents or adding new configuration options |
| `agent/state.go` | State management for agent execution including history, messages, and tool results | Understanding agent state lifecycle or implementing state persistence |

## Architecture

```
User Request → Agent.Run/Stream → Context Setup → State Init → LLM Generate → Tool Detection → Tool Execution → Result Merge → LLM Continue → Response Stream → Final State

Key Components:
- Agent: Main orchestrator (agent.go)
- Context: Request-scoped metadata (context.go)
- State: Execution history and tool results (state.go)
- Options: Configuration and callbacks (options.go)
- Tools: Executable functions registered with agent

Data Flow:
1. Input message + context → Agent
2. Agent creates/updates State with history
3. LLM generates response (may include tool calls)
4. Agent executes tools if requested
5. Tool results fed back to LLM
6. Loop until final text response
7. State returned with full history
```

## Patterns

### Functional Options Pattern

Agent configuration uses variadic Option functions for flexible, extensible setup without breaking changes

See `agent/options.go` for reference.

### Context-Based State Management

Agent metadata and state stored in context.Context using typed keys, enabling middleware and request-scoped data

See `agent/context.go` for reference.

### Streaming Iterator Pattern

Agent.Stream returns iter.Seq2[*genkit.Candidate, error] for lazy evaluation and backpressure handling

See `agent/agent.go` for reference.

### Tool Execution Loop

Agent automatically detects tool calls in LLM responses, executes them, and feeds results back until text response received

See `agent/agent.go` for reference.

### Immutable State Updates

State objects are copied and updated rather than mutated, enabling safe concurrent access and history tracking

See `agent/state.go` for reference.

### Interface-Based Abstraction

Agent interface allows multiple implementations (basic, advanced) while maintaining consistent API

See `agent/agent.go` for reference.

## Rules

### Always (Do These)

| Rule | Rationale |
|------|------------|
| CRITICAL: Always propagate context.Context through all agent operations - never create new background contexts | Context carries cancellation signals, deadlines, tracing data, and agent state. Breaking the chain loses all request-scoped data and prevents proper cleanup |
| CRITICAL: Check for tool calls in LLM responses and execute them before returning final response | Agents must complete tool execution loops. Returning partial responses breaks the agent contract and leaves operations incomplete |
| HIGH: Use functional options (WithModel, WithTools, etc.) for all agent configuration | Maintains backward compatibility, enables optional features, and provides clear configuration intent |
| HIGH: Store agent state in context using NewContextWithState/StateFromContext, never as global variables | Enables concurrent agent execution, proper request isolation, and middleware composition |
| HIGH: Return errors immediately from agent operations - do not swallow or log-and-continue | Agents must fail fast on errors to prevent cascading failures and invalid state propagation |
| MEDIUM: Use iter.Seq2[*genkit.Candidate, error] for streaming responses, not channels | Go 1.23+ iterators provide better composition, automatic cleanup, and standard library integration |
| MEDIUM: Include tool results in state history for LLM context in subsequent calls | LLMs need full conversation history including tool interactions to maintain coherent responses |
| MEDIUM: Validate tool definitions before registering them with agents | Invalid tool schemas cause runtime failures during LLM tool calling that are hard to debug |

### Never (Don't Do These)

| Rule | Severity | Rationale |
|------|----------|------------|
| NEVER modify State objects in place - always create new copies with updates | CRITICAL | State mutation causes race conditions in concurrent agent execution and breaks history tracking |
| NEVER ignore context cancellation - always check ctx.Err() in loops and long operations | CRITICAL | Ignoring cancellation causes resource leaks, hangs, and prevents graceful shutdown |
| NEVER execute tools without validating their input schemas | HIGH | Invalid tool inputs cause panics or undefined behavior in tool implementations |
| NEVER return nil State from agent operations - return empty State instead | HIGH | Nil states cause nil pointer dereferences in downstream code expecting valid state objects |
| NEVER store mutable state in Agent struct fields - use context or State instead | HIGH | Agent instances may be reused across requests, causing state leakage between unrelated operations |
| NEVER block indefinitely in tool execution - always respect context deadlines | MEDIUM | Blocking tools prevent agent timeout handling and cause resource exhaustion |
| NEVER assume tool calls will succeed - always handle tool execution errors gracefully | MEDIUM | Tool failures are common (network issues, invalid inputs) and must be recoverable |

### Ask First

- **Adding new methods to the Agent interface** - Interface changes break all existing implementations. Consider adding new optional interfaces or using functional options instead
- **Changing State struct fields or adding required fields** - State is serialized and persisted. Schema changes require migration strategy and backward compatibility
- **Modifying tool execution order or loop termination conditions** - Changes affect all agents and may break existing tool chains or cause infinite loops
- **Adding new context keys or changing context key types** - Context keys must be unique across codebase. Collisions cause silent data corruption
- **Changing streaming behavior or iterator semantics** - Streaming contracts are relied upon by consumers. Breaking changes cause deadlocks or data loss
- **Modify configuration file: agent/run_config.go** - Configuration changes can affect all environments

## Commands

### test-agent

Run all agent package tests including unit and integration tests

```bash
go test ./agent/...
```

### test-agent-verbose

Run agent tests with verbose output showing individual test cases

```bash
go test -v ./agent/...
```

### test-agent-coverage

Generate and view test coverage report for agent package

```bash
go test -coverprofile=coverage.out ./agent/... && go tool cover -html=coverage.out
```

### bench-agent

Run agent performance benchmarks with memory allocation stats

```bash
go test -bench=. -benchmem ./agent/...
```

## Testing

Table-driven tests with mock LLMs and tools. Tests cover happy paths, error cases, tool execution loops, context propagation, and streaming behavior. See agent/agent_test.go for comprehensive examples.

```bash
go test ./agent/...
go test -race ./agent/...
go test -v -run TestAgentRun ./agent/...
```

Test directory: `agent/`

## Common Tasks

### Create a new agent with tools

1. 1. Define tools using genkit.DefineTool with input/output schemas
2. 2. Create agent: agent := agent.New(agent.WithModel(model), agent.WithTools(tools...))
3. 3. Call agent.Run(ctx, state, input) or agent.Stream(ctx, state, input)
4. 4. Handle returned State and errors appropriately
5. See agent/agent_test.go for complete examples

### Implement custom tool for agent

1. 1. Define input/output structs with json tags
2. 2. Use genkit.DefineTool(name, description, func(ctx, input) (output, error))
3. 3. Ensure tool respects context cancellation (check ctx.Err())
4. 4. Return structured errors for tool failures
5. 5. Register tool with agent using WithTools option
6. See agent/agent_test.go - mock tool implementations

### Add agent state persistence

1. 1. Extract State from agent.Run/Stream return value
2. 2. Serialize State to JSON/database (State is JSON-serializable)
3. 3. On next request, deserialize State and pass to agent.Run/Stream
4. 4. State contains full conversation history and tool results
5. See agent/state.go for State structure

### Implement streaming agent responses

1. 1. Call agent.Stream(ctx, state, input) to get iterator
2. 2. Iterate: for candidate, err := range agent.Stream(...) { handle candidate }
3. 3. Check err on each iteration for streaming errors
4. 4. Accumulate candidates to build final response
5. 5. Extract final State from last candidate
6. See agent/agent.go - Stream method implementation

### Add middleware to agent execution

1. 1. Create context with custom values before calling agent
2. 2. Use agent/context.go functions to store/retrieve agent metadata
3. 3. Implement callbacks using WithStartCallback/WithEndCallback options
4. 4. Access context in tools to read middleware-injected data
5. See agent/context.go for context key patterns

### Debug tool execution issues

1. 1. Enable verbose logging in agent options
2. 2. Check State.History for tool call requests and responses
3. 3. Verify tool input schemas match LLM output format
4. 4. Add logging in tool implementation to trace execution
5. 5. Check for context cancellation causing premature termination
6. See agent/agent_test.go for tool debugging patterns

