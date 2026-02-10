# AGENTS.md - runner

> The runner package provides the core execution engine for AI agents, managing the agent lifecycle, tool execution, and LLM interaction loops. It orchestrates the flow between user input, LLM responses, tool calls, and final outputs while handling context, errors, and state management.

## Tech Stack

- **Go** 1.x - Core implementation language for agent execution engine
- **context** stdlib - Manages cancellation, timeouts, and request-scoped values throughout agent execution
- **testing** stdlib - Unit testing framework with table-driven tests for runner behavior

## Key Files

| File | Purpose | Read When |
|------|---------|------------|
| `runner/runner.go` (reference) | Main runner implementation with Run() method that executes agent loops, handles tool calls, and manages LLM interactions | Understanding agent execution flow, implementing new runner features, or debugging agent behavior |
| `runner/runner_test.go` (reference) | Comprehensive test suite demonstrating runner behavior with mocked LLMs and tools, showing expected execution patterns | Understanding expected runner behavior, writing new tests, or validating changes to execution logic |

## Architecture

```
User Input → Runner.Run() → LLM Request → LLM Response Analysis:
  ├─ Text Response? → Return to User
  ├─ Tool Call? → Execute Tool → Add Result to History → Loop to LLM
  └─ Error? → Handle & Return

Key Components:
- Runner: Orchestrates execution loop
- Agent: Provides tools, system prompt, model config
- LLM: Generates responses and tool calls
- Tools: Execute actions and return results
- Context: Manages cancellation and timeouts
- History: Maintains conversation state
```

## Patterns

### Agent Execution Loop

Iterative loop that sends messages to LLM, processes responses (text or tool calls), executes tools, and continues until final response. Uses history accumulation pattern.

See `runner/runner.go - Run() method` for reference.

### Tool Call Resolution

When LLM returns tool calls, runner looks up tools by name from agent, executes them with provided arguments, and appends results to message history before next LLM call.

See `runner/runner.go - tool execution logic in Run()` for reference.

### Context Propagation

Context passed through entire execution chain (Run → LLM → Tools) enabling cancellation, timeouts, and request-scoped values throughout agent lifecycle.

See `runner/runner.go - ctx parameter threading` for reference.

### Mock-Based Testing

Tests use mock LLMs and tools to verify runner behavior without external dependencies. Table-driven tests cover various scenarios (tool calls, errors, text responses).

See `runner/runner_test.go - TestRun with mock implementations` for reference.

### Error Wrapping

Errors wrapped with context using fmt.Errorf with %w verb to maintain error chain while adding contextual information about where failure occurred.

See `runner/runner.go - error returns in Run()` for reference.

## Rules

### Always (Do These)

| Rule | Rationale |
|------|------------|
| CRITICAL: Pass context.Context as first parameter to Run() and propagate it to all LLM and tool calls | Enables proper cancellation, timeout handling, and request tracing. Breaking this prevents graceful shutdown and can cause resource leaks. |
| HIGH: Accumulate all messages (user, assistant, tool results) in history slice before each LLM call | LLM needs full conversation context to generate appropriate responses. Missing history causes context loss and incorrect agent behavior. |
| HIGH: Look up tools from agent.Tools() map by name when processing tool calls | Tool resolution must use agent's registered tools. Direct tool access bypasses agent configuration and breaks tool isolation. |
| MEDIUM: Return immediately when LLM response contains text content (non-tool response) | Text responses indicate agent has completed its work. Continuing the loop would cause unnecessary LLM calls and incorrect behavior. |
| HIGH: Wrap errors with context using fmt.Errorf with %w to maintain error chain | Preserves original error for debugging while adding context about where in execution flow error occurred. Critical for troubleshooting. |
| MEDIUM: Execute all tool calls from a single LLM response before making next LLM request | LLM may request multiple tools in one response. All must execute and results added to history before next LLM call for proper context. |

### Never (Don't Do These)

| Rule | Severity | Rationale |
|------|----------|------------|
| CRITICAL: Never ignore context cancellation - always check ctx.Err() or propagate context to blocking operations | CRITICAL | Ignoring cancellation causes goroutine leaks, resource exhaustion, and prevents graceful shutdown. Can crash production systems. |
| CRITICAL: Never modify the agent's tool map during execution | CRITICAL | Concurrent modification of shared tool map causes race conditions and undefined behavior. Tools must be immutable during run. |
| HIGH: Never skip adding tool results to history before next LLM call | HIGH | LLM needs tool results to continue reasoning. Skipping results causes infinite loops or incorrect responses as LLM lacks context. |
| HIGH: Never return partial results when tool execution fails | HIGH | Partial execution leaves agent in inconsistent state. Must return error to caller so they can handle failure appropriately. |
| MEDIUM: Never assume LLM response will contain either text OR tool calls - validate response structure | MEDIUM | LLM responses can be malformed or empty. Must validate structure before processing to avoid panics or incorrect behavior. |
| HIGH: Never create new runner instances per request - runner should be stateless and reusable | HIGH | Runner contains no request-specific state. Creating per-request wastes resources and suggests architectural misunderstanding. |

### Ask First

- **Adding state fields to Runner struct** - Runner is designed to be stateless and reusable. Adding state may indicate architectural issue or require thread-safety considerations.
- **Changing the order of history accumulation or tool execution** - Execution order is critical for correct LLM context. Changes can break agent reasoning and cause subtle bugs.
- **Adding retry logic or error recovery in the execution loop** - Error handling strategy affects all agents. Should be consistent across system and may need configuration options.
- **Modifying how context is propagated or adding context values** - Context usage patterns affect entire system. Changes should be coordinated with other packages and documented.
- **Adding hooks or callbacks in the execution loop** - Affects observability and extensibility patterns. Should align with system-wide instrumentation strategy.

## Commands

### test

Run all runner package tests

```bash
go test ./runner/...
```

### test-verbose

Run tests with verbose output showing each test case

```bash
go test -v ./runner/...
```

### test-coverage

Run tests with coverage report

```bash
go test -cover ./runner/...
```

### test-race

Run tests with race detector to catch concurrency issues

```bash
go test -race ./runner/...
```

### bench

Run benchmarks if present

```bash
go test -bench=. ./runner/...
```

## Testing

Table-driven tests with mock LLMs and tools to verify runner behavior across scenarios: successful tool calls, errors, text responses, multiple tool calls, and edge cases. Tests validate execution flow, history accumulation, error handling, and context propagation without external dependencies.

```bash
go test ./runner/...
go test -v ./runner/...
go test -race ./runner/...
```

Test directory: `runner/`

## Common Tasks

### Add new runner feature

1. 1. Review runner/runner.go to understand current execution flow
2. 2. Identify where in Run() loop new feature should integrate
3. 3. Ensure feature respects context cancellation and error handling patterns
4. 4. Add feature implementation maintaining stateless runner design
5. 5. Add test cases in runner/runner_test.go using mock pattern
6. 6. Verify feature works with tool calls, text responses, and error cases
7. 7. Run go test -race to check for concurrency issues

### Debug agent execution issue

1. 1. Add logging/tracing in runner/runner.go Run() method to track execution flow
2. 2. Verify context is properly propagated to LLM and tools
3. 3. Check history accumulation - ensure all messages added in correct order
4. 4. Validate tool lookup succeeds and tools execute with correct arguments
5. 5. Verify LLM responses are properly parsed (tool calls vs text)
6. 6. Check error wrapping maintains full error chain
7. 7. Review runner/runner_test.go for similar test scenarios

### Add test for new execution scenario

1. 1. Open runner/runner_test.go and review existing test structure
2. 2. Create mock LLM that returns responses matching your scenario
3. 3. Create mock tools if scenario involves tool execution
4. 4. Add test case to table-driven test or create new test function
5. 5. Verify history accumulation matches expected pattern
6. 6. Test both success and error paths for scenario
7. 7. Run go test -v to verify test passes

### Optimize runner performance

1. 1. Add benchmarks in runner/runner_test.go for target scenarios
2. 2. Run go test -bench=. -benchmem to establish baseline
3. 3. Profile with go test -cpuprofile=cpu.out -memprofile=mem.out
4. 4. Identify bottlenecks (likely in LLM calls or tool execution)
5. 5. Optimize without changing execution semantics or breaking tests
6. 6. Re-run benchmarks to verify improvement
7. 7. Ensure go test -race still passes after optimization

### Integrate new LLM provider

1. 1. Review runner/runner.go to see how LLM interface is used
2. 2. Implement LLM interface for new provider (see model package)
3. 3. Create mock version for testing in runner/runner_test.go
4. 4. Add test cases verifying runner works with new provider
5. 5. Test tool call format compatibility with runner expectations
6. 6. Verify context propagation and cancellation work correctly
7. 7. Test error handling for provider-specific errors

