# AGENTS.md - agent/llmagent

> The llmagent package implements LLM-powered agents with state management, tool execution, and multi-turn conversation capabilities. It provides both stateless and stateful agent implementations that integrate with LLM models, handle tool calls, and manage conversation history with support for streaming responses and error recovery.

## Tech Stack

- **Go** 1.x - Core implementation language for agent logic and state management
- **genai** - LLM model integration and content generation
- **context** stdlib - Request lifecycle management and cancellation

## Key Files

| File | Purpose | Read When |
|------|---------|------------|
| `agent/llmagent/llmagent.go` (reference) | Core stateless LLM agent implementation with tool execution and streaming support | Understanding basic agent execution flow and tool integration |
| `agent/llmagent/state_agent.go` (reference) | Stateful agent wrapper that maintains conversation history and manages multi-turn interactions | Implementing agents that need to maintain context across multiple turns |
| `agent/llmagent/llmagent_test.go` (reference) | Test patterns for agent execution, tool calling, and error handling | Writing tests or understanding expected agent behavior |
| `agent/llmagent/state_agent_test.go` (reference) | Test patterns for stateful agents, history management, and state persistence | Testing stateful agent implementations or history management |
| `agent/llmagent/options.go` | Configuration options for agent behavior (max iterations, tool choice, etc.) | Configuring agent behavior or understanding available options |
| `agent/llmagent/types.go` (reference) | Core type definitions for agents, responses, and state management | Understanding agent interfaces and data structures |

## Architecture

```
Agent Execution Flow:

[User Input] → [Agent.Run()]
       ↓
[Build Prompt + History] → [LLM Model]
       ↓
[Response Analysis]
       ├─→ [Text Response] → [Return to User]
       └─→ [Tool Calls Detected]
              ↓
       [Execute Tools in Parallel]
              ↓
       [Append Tool Results to History]
              ↓
       [Loop Back to LLM] (max iterations check)
              ↓
       [Final Response]

Stateful Agent adds:
[StateAgent] wraps [Agent]
     ↓
[Maintains History]
     ↓
[Persists State Between Calls]
```

## Patterns

### Functional Options Pattern

Agent configuration uses functional options (WithMaxIterations, WithToolChoice, etc.) for flexible initialization

See `agent/llmagent/options.go` for reference.

### Streaming Response Pattern

Agents support streaming responses via channels, allowing real-time token delivery to clients

See `agent/llmagent/llmagent.go` for reference.

### Tool Execution Loop

Agents iterate between LLM calls and tool execution until a final answer is reached or max iterations hit

See `agent/llmagent/llmagent.go` for reference.

### State Wrapper Pattern

StateAgent wraps base Agent to add history management without modifying core logic

See `agent/llmagent/state_agent.go` for reference.

### Parallel Tool Execution

Multiple tool calls from a single LLM response are executed concurrently using goroutines

See `agent/llmagent/llmagent.go` for reference.

### Context Propagation

Context is threaded through all operations for cancellation and timeout support

See `agent/llmagent/llmagent.go` for reference.

## Rules

### Always (Do These)

| Rule | Rationale |
|------|------------|
| CRITICAL: Always check MaxIterations to prevent infinite loops in tool execution cycles | Agents can get stuck in tool-calling loops if LLM keeps requesting tools without providing final answers. MaxIterations prevents runaway execution and resource exhaustion. |
| HIGH: Always append tool results to history before next LLM call | LLM needs to see tool execution results to generate informed responses. Missing tool results breaks the agent's reasoning chain. |
| HIGH: Always propagate context.Context through agent operations | Enables proper cancellation, timeout handling, and request tracing. Critical for production reliability. |
| MEDIUM: Use StateAgent when conversation history is needed across multiple turns | StateAgent automatically manages history persistence. Using base Agent for multi-turn conversations loses context. |
| HIGH: Execute multiple tool calls in parallel when possible | LLMs often request multiple independent tools. Parallel execution significantly reduces latency. |
| MEDIUM: Include system instructions in agent configuration for consistent behavior | System instructions guide agent behavior and tool usage patterns. Missing instructions leads to unpredictable responses. |
| HIGH: Handle streaming errors by checking channel closure and error returns | Streaming responses can fail mid-stream. Proper error handling prevents silent failures and data corruption. |

### Never (Don't Do These)

| Rule | Severity | Rationale |
|------|----------|------------|
| CRITICAL: Never modify agent history directly - use provided methods | CRITICAL | Direct history modification bypasses state management logic and can corrupt conversation state, leading to incorrect LLM responses. |
| CRITICAL: Never ignore MaxIterations - always enforce iteration limits | CRITICAL | Removing iteration limits allows infinite loops that can exhaust API quotas, cause timeouts, and waste resources. |
| HIGH: Never execute tools without proper error handling and recovery | HIGH | Tool execution failures must be captured and reported back to LLM. Silent failures break agent reasoning and produce incorrect results. |
| HIGH: Never block indefinitely on streaming channels without context cancellation | HIGH | Streaming operations must respect context cancellation to prevent goroutine leaks and hung requests. |
| MEDIUM: Never assume tool calls will complete successfully - always handle failures | MEDIUM | Tools can fail for many reasons (network, permissions, invalid input). Agents must gracefully handle failures and report them to LLM. |
| HIGH: Never reuse Agent instances across concurrent requests without proper synchronization | HIGH | Agent state is not thread-safe. Concurrent use causes race conditions and corrupted state. |
| MEDIUM: Never skip validation of tool responses before passing to LLM | MEDIUM | Invalid tool responses can confuse LLM or cause parsing errors. Validation ensures clean data flow. |

### Ask First

- **Changing MaxIterations default value** - Affects all agent behavior and resource usage. May impact production systems and API costs.
- **Modifying tool execution parallelization strategy** - Changes performance characteristics and may affect tool execution order assumptions.
- **Altering history management in StateAgent** - Impacts conversation continuity and state persistence. Breaking changes affect all stateful agents.
- **Adding new agent lifecycle hooks or callbacks** - Affects agent execution flow and may introduce breaking changes to existing implementations.
- **Changing streaming response format or channel behavior** - Breaking change for all streaming consumers. Requires coordinated updates across codebase.

## Commands

### test

Run all llmagent tests

```bash
go test ./agent/llmagent/...
```

### test-verbose

Run tests with verbose output

```bash
go test -v ./agent/llmagent/...
```

### test-coverage

Run tests with coverage report

```bash
go test -cover ./agent/llmagent/...
```

### bench

Run benchmarks for agent performance

```bash
go test -bench=. ./agent/llmagent/...
```

## Testing

Table-driven tests with mock LLM models and tools. Tests cover stateless agents, stateful agents, streaming, error handling, and tool execution patterns. Mock implementations allow testing without real LLM API calls.

```bash
go test ./agent/llmagent/...
go test -race ./agent/llmagent/...
```

Test directory: `agent/llmagent/`

## Common Tasks

### Create a basic stateless agent

1. Import agent/llmagent package
2. Create LLM model instance (see genai integration)
3. Define tools using tool package
4. Initialize agent: agent.New(model, tools, options...)
5. Call agent.Run(ctx, input) for single execution
6. See agent/llmagent/llmagent_test.go for examples

### Create a stateful agent with history

1. Create base agent as above
2. Wrap with StateAgent: stateAgent := agent.NewStateAgent(baseAgent)
3. Call stateAgent.Run(ctx, input) - history is automatically managed
4. Access history via stateAgent.History() if needed
5. See agent/llmagent/state_agent_test.go for examples

### Implement streaming responses

1. Create agent as normal
2. Call agent.RunStream(ctx, input) to get response channel
3. Iterate over channel: for resp := range responseChan
4. Check resp.Error for errors during streaming
5. Handle resp.Content for incremental text
6. See agent/llmagent/llmagent.go RunStream implementation

### Configure agent behavior

1. Use functional options during agent creation
2. WithMaxIterations(n) - limit tool execution loops
3. WithToolChoice(choice) - control tool selection
4. WithSystemInstruction(text) - set agent behavior guidelines
5. See agent/llmagent/options.go for all available options

### Handle tool execution errors

1. Tools should return errors for failures
2. Agent automatically captures tool errors
3. Errors are formatted and sent back to LLM
4. LLM can retry or provide alternative responses
5. Check agent response for final error status
6. See agent/llmagent/llmagent_test.go error handling tests

### Debug agent execution flow

1. Enable verbose logging if available
2. Inspect agent.History() to see conversation flow
3. Check iteration count against MaxIterations
4. Examine tool call requests and responses in history
5. Use context with timeout to prevent hangs
6. See agent/llmagent/state_agent_test.go for history inspection patterns

