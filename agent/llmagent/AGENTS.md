# AGENTS.md - agent/llmagent

> The llmagent package implements LLM-powered agents with state management, tool execution, and multi-turn conversation capabilities. It provides both stateless and stateful agent implementations that integrate with LLM providers, handle tool calls, and manage conversation history with configurable memory strategies.

## Tech Stack

- **Go** 1.x - Core implementation language for agent logic and state management
- **genai** - LLM model abstraction layer for generating responses and handling tool calls
- **agent/tool** - Tool registration and execution framework for agent capabilities

## Key Files

| File | Purpose | Read When |
|------|---------|------------|
| `agent/llmagent/llmagent.go` (reference) | Core stateless LLM agent implementation with tool execution and response generation | Understanding basic agent structure, tool integration, or response handling |
| `agent/llmagent/state_agent.go` (reference) | Stateful agent wrapper that manages conversation history and memory strategies | Implementing multi-turn conversations or state persistence |
| `agent/llmagent/options.go` (reference) | Configuration options for agents including system prompts, tools, and memory settings | Configuring agent behavior or adding new configuration options |
| `agent/llmagent/memory.go` (reference) | Memory strategy implementations (sliding window, token-based) for conversation history | Understanding or modifying how conversation history is managed |
| `agent/llmagent/llmagent_test.go` | Test patterns for stateless agent behavior, tool execution, and error handling | Writing tests or understanding expected agent behavior |
| `agent/llmagent/state_agent_test.go` | Test patterns for stateful agents, memory strategies, and conversation flow | Testing stateful behavior or memory management |

## Architecture

```
Request Flow:
1. Client → StateAgent.Run(ctx, input) [state_agent.go]
2. StateAgent loads history from memory strategy [memory.go]
3. StateAgent → LLMAgent.Run(ctx, history+input) [llmagent.go]
4. LLMAgent → Model.Generate(ctx, messages) [genai integration]
5. Model returns response with potential tool calls
6. LLMAgent executes tools via ToolRegistry [tool integration]
7. If tools executed → loop back to step 4 with tool results
8. LLMAgent returns final response
9. StateAgent updates memory with new messages
10. Response → Client

Key Components:
- LLMAgent: Stateless executor (tool calls, generation)
- StateAgent: Stateful wrapper (history, memory)
- Memory: Strategy pattern for history management
- ToolRegistry: Tool discovery and execution
```

## Patterns

### Functional Options Pattern

Agent configuration uses variadic option functions (WithSystemPrompt, WithTools, WithMemory) for flexible initialization

See `agent/llmagent/options.go` for reference.

### Strategy Pattern for Memory

Memory interface with multiple implementations (SlidingWindowMemory, TokenBasedMemory) allows pluggable history management

See `agent/llmagent/memory.go` for reference.

### Stateless Core with Stateful Wrapper

LLMAgent is stateless for single-turn execution; StateAgent wraps it to add conversation state management

See `agent/llmagent/state_agent.go` for reference.

### Tool Call Loop Pattern

Agent iteratively calls LLM → execute tools → call LLM with results until no more tool calls needed

See `agent/llmagent/llmagent.go` for reference.

### Context Propagation

All operations accept context.Context for cancellation, timeouts, and request-scoped values

See `agent/llmagent/llmagent.go` for reference.

## Rules

### Always (Do These)

| Rule | Rationale |
|------|------------|
| CRITICAL: Pass context.Context as first parameter to all agent Run methods and propagate it through tool executions | Enables proper cancellation, timeout handling, and request tracing across the entire agent execution chain |
| HIGH: Use StateAgent for multi-turn conversations; use LLMAgent only for single-turn stateless operations | StateAgent manages conversation history and memory; LLMAgent requires caller to manage all history |
| HIGH: Configure memory strategy explicitly when creating StateAgent (WithMemory option) | Default memory behavior may not suit all use cases; explicit configuration prevents unbounded memory growth |
| HIGH: Register tools before agent creation using WithTools() option, not after | Tools must be available during agent initialization for proper LLM function declaration |
| MEDIUM: Include system prompts via WithSystemPrompt() to guide agent behavior | System prompts establish agent personality, constraints, and operational guidelines |
| MEDIUM: Handle tool execution errors gracefully - they should be passed back to LLM as tool results | LLM can reason about tool failures and potentially retry or use alternative approaches |
| MEDIUM: Check response.FinishReason to detect truncation, safety blocks, or other non-standard completions | Finish reasons indicate whether response is complete or was interrupted by limits/filters |

### Never (Don't Do These)

| Rule | Severity | Rationale |
|------|----------|------------|
| CRITICAL: Never modify conversation history directly in StateAgent - always use the Memory interface | CRITICAL | Direct modification bypasses memory strategy logic (windowing, token limits) and causes state inconsistencies |
| CRITICAL: Never ignore context cancellation in long-running agent operations or tool executions | CRITICAL | Ignoring cancellation leads to resource leaks and unresponsive agents that continue after client disconnect |
| HIGH: Never assume tool calls will succeed - always handle tool execution errors | HIGH | Tool failures are common (network issues, invalid params); unhandled errors crash agent execution |
| HIGH: Never create multiple StateAgent instances for the same conversation without shared memory | HIGH | Each StateAgent maintains independent history; multiple instances lose conversation context |
| HIGH: Never pass nil Model to NewLLMAgent - it will panic on first Run call | HIGH | Model is required for generation; nil check happens at runtime causing panic |
| MEDIUM: Never rely on specific tool call order - LLM may call tools in any sequence or parallel | MEDIUM | Tool execution order is determined by LLM reasoning, not guaranteed by agent implementation |
| MEDIUM: Never store sensitive data in system prompts without considering memory persistence | MEDIUM | System prompts are stored in conversation history and may be logged or persisted |

### Ask First

- **Adding new Memory strategy implementations** - Memory strategies affect performance and behavior across all agents; coordinate with existing implementations in memory.go
- **Changing tool execution loop logic or retry behavior** - Core agent behavior affects all users; changes may break existing tool implementations or conversation flows
- **Modifying how conversation history is structured or stored** - History format changes may break existing memory implementations and persistence layers
- **Adding agent-level caching or memoization** - Caching affects statefulness guarantees and may interact unexpectedly with memory strategies
- **Changing error handling for LLM generation failures** - Error handling patterns affect all agent users; changes should be consistent with broader error strategy

## Commands

### test-llmagent

Run all tests for llmagent package including stateless and stateful agent tests

```bash
go test ./agent/llmagent/...
```

### test-verbose

Run tests with verbose output to see individual test execution

```bash
go test -v ./agent/llmagent/...
```

### test-coverage

Run tests with coverage reporting

```bash
go test -cover ./agent/llmagent/...
```

### test-specific

Run specific test by name (e.g., TestLLMAgent_Run, TestStateAgent_Memory)

```bash
go test ./agent/llmagent -run TestName
```

## Testing

Table-driven tests with mock LLM models and tools. Tests cover stateless agent behavior, stateful conversation flows, memory strategies, tool execution, and error handling. Mock implementations verify correct LLM calls and tool invocations.

```bash
go test ./agent/llmagent/...
go test -v ./agent/llmagent/... -run TestLLMAgent
go test -v ./agent/llmagent/... -run TestStateAgent
```

Test directory: `agent/llmagent/`

## Common Tasks

### Create a basic stateless agent

1. Import agent/llmagent and your model implementation
2. Create model instance: model := yourprovider.NewModel(...)
3. Create agent: agent := llmagent.NewLLMAgent(model, llmagent.WithSystemPrompt("..."))
4. Run agent: response, err := agent.Run(ctx, []genai.Message{...})
5. Handle response and errors appropriately
6. See agent/llmagent/llmagent_test.go for examples

### Create a stateful agent with conversation history

1. Create base LLM agent as above
2. Choose memory strategy: mem := llmagent.NewSlidingWindowMemory(10) or llmagent.NewTokenBasedMemory(4096)
3. Create state agent: stateAgent := llmagent.NewStateAgent(agent, llmagent.WithMemory(mem))
4. Run conversations: response, err := stateAgent.Run(ctx, userMessage)
5. History is automatically managed by memory strategy
6. See agent/llmagent/state_agent_test.go for examples

### Add tools to an agent

1. Define tool functions matching tool.Tool interface
2. Create tool registry: registry := tool.NewRegistry()
3. Register tools: registry.Register("tool_name", toolFunc)
4. Create agent with tools: agent := llmagent.NewLLMAgent(model, llmagent.WithTools(registry))
5. Agent will automatically execute tool calls from LLM
6. Tool results are passed back to LLM for reasoning
7. See agent/llmagent/llmagent_test.go for tool integration examples

### Implement custom memory strategy

1. Create type implementing llmagent.Memory interface
2. Implement Add(msg genai.Message) to store messages
3. Implement Get() []genai.Message to retrieve history
4. Implement Clear() to reset history
5. Apply your strategy logic in Add() (e.g., pruning, summarization)
6. Use with StateAgent: llmagent.WithMemory(yourMemory)
7. See agent/llmagent/memory.go for reference implementations

### Handle agent errors and edge cases

1. Always check error return from agent.Run()
2. Check response.FinishReason for non-standard completions
3. Handle context cancellation: if ctx.Err() != nil { ... }
4. For tool errors, check if they're passed to LLM as tool results
5. Consider retry logic for transient LLM failures
6. Log errors with sufficient context for debugging
7. See agent/llmagent/llmagent_test.go for error handling patterns

### Debug agent behavior

1. Enable verbose logging in your model implementation
2. Inspect conversation history: stateAgent.Memory.Get()
3. Check tool execution by logging in tool functions
4. Verify system prompt is correctly set
5. Test with mock model to isolate agent logic from LLM variability
6. Use debugger to step through tool execution loop
7. Check response.FinishReason and response.Content for unexpected values

