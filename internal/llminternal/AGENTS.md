# AGENTS.md - internal/llminternal

> The llminternal package provides core LLM agent functionality including content processing, instruction handling, agent transfer mechanisms, and model abstraction. It serves as the internal implementation layer for agent-based workflows, handling message transformation, tool execution, and multi-agent coordination.

## Tech Stack

- **Go** 1.21+ - Primary language for LLM agent implementation
- **genai SDK** - Google Generative AI integration for model interactions
- **context package** stdlib - Request lifecycle and cancellation management

## Key Files

| File | Purpose | Read When |
|------|---------|------------|
| `internal/llminternal/agent_transfer.go` (reference) | Implements agent transfer protocol for multi-agent coordination | Understanding how agents delegate work to other agents |
| `internal/llminternal/instruction_processor.go` (reference) | Processes and transforms instructions for LLM consumption | Working with prompt engineering or instruction formatting |
| `internal/llminternal/contents_processor.go` (reference) | Handles content transformation between different message formats | Implementing message parsing or content type conversions |
| `internal/llminternal/contents_processor_test.go` (reference) | Comprehensive test suite for content processing logic | Understanding expected behavior for message transformations |
| `internal/llminternal/model_client.go` | Abstraction layer for LLM model interactions | Implementing or modifying model communication patterns |

## Architecture

```
Agent Request → Instruction Processor → Contents Processor → Model Client → LLM API
                                                                    ↓
                                                          Agent Transfer Handler
                                                                    ↓
                                                            Delegate to Sub-Agent

Key Flow:
1. Instructions are processed and formatted (instruction_processor.go)
2. Content is transformed to LLM-compatible format (contents_processor.go)
3. Model client handles API communication (model_client.go)
4. Agent transfer mechanism enables multi-agent workflows (agent_transfer.go)
5. Responses flow back through processors for normalization
```

## Patterns

### Content Transformation Pipeline

Multi-stage processing of messages through type-specific handlers. Each content type (text, tool call, tool response) has dedicated transformation logic.

See `internal/llminternal/contents_processor.go` for reference.

### Context-First Design

All operations accept context.Context as first parameter for cancellation, timeouts, and request-scoped values.

See `internal/llminternal/model_client.go` for reference.

### Agent Transfer Protocol

Structured mechanism for agents to delegate work to specialized sub-agents with explicit transfer metadata.

See `internal/llminternal/agent_transfer.go` for reference.

### Error Wrapping with Context

Errors are wrapped with fmt.Errorf to preserve stack context and provide actionable error messages.

See `internal/llminternal/instruction_processor.go` for reference.

### Table-Driven Testing

Test cases defined as structs with input/expected output pairs for comprehensive coverage.

See `internal/llminternal/contents_processor_test.go` for reference.

## Rules

### Always (Do These)

| Rule | Rationale |
|------|------------|
| Pass context.Context as the first parameter to all functions that perform I/O or may need cancellation | Enables proper request lifecycle management, timeout handling, and graceful cancellation across agent workflows |
| Validate content types before processing and return descriptive errors for unsupported types | Content processors must handle multiple message formats; explicit validation prevents runtime panics and provides clear debugging information |
| Wrap errors with context using fmt.Errorf with %w verb to preserve error chains | Maintains error causality for debugging while adding operation-specific context |
| Use table-driven tests with named test cases for all content transformation logic | Content processing has many edge cases; structured tests ensure comprehensive coverage and make test intent clear |
| Include agent transfer metadata when delegating to sub-agents | Transfer metadata enables proper tracking, debugging, and coordination in multi-agent systems |
| Normalize message formats when converting between internal and external representations | Different LLM providers use different message schemas; normalization ensures consistent behavior |

### Never (Don't Do These)

| Rule | Severity | Rationale |
|------|----------|------------|
| Never ignore context cancellation signals in long-running operations | CRITICAL | Ignoring context cancellation can lead to resource leaks, hung requests, and inability to terminate agent workflows |
| Never mutate input parameters; always create new objects for transformations | HIGH | Content processors are used in concurrent agent workflows; mutation causes race conditions and unpredictable behavior |
| Never panic on invalid input; return errors instead | HIGH | Agent systems must be resilient to malformed input from LLMs or external sources; panics crash the entire process |
| Never log sensitive data (API keys, user content) in error messages or debug output | CRITICAL | LLM interactions may contain PII or confidential information; logging creates security vulnerabilities |
| Never assume tool call responses are well-formed JSON without validation | HIGH | LLMs can generate malformed tool calls; processing without validation causes runtime errors |
| Never create circular agent transfer chains without depth limits | CRITICAL | Unbounded agent delegation can cause infinite loops and resource exhaustion |

### Ask First

- **Adding new content type handlers to the processor pipeline** - New content types affect the entire message transformation flow; requires coordination with model client and API layer
- **Modifying agent transfer protocol or metadata structure** - Changes impact multi-agent coordination and may break existing agent implementations
- **Changing error handling patterns or error types** - Error handling is used throughout the codebase; changes affect error recovery and debugging capabilities
- **Adding dependencies on external LLM provider SDKs** - Provider dependencies affect abstraction layer design and may introduce version conflicts
- **Modifying instruction processing logic that affects prompt formatting** - Prompt changes can significantly impact LLM behavior and agent performance; requires testing across models

## Commands

### test-package

Run all tests in the llminternal package

```bash
go test ./internal/llminternal/...
```

### test-verbose

Run tests with verbose output to see individual test cases

```bash
go test -v ./internal/llminternal/...
```

### test-coverage

Run tests with coverage reporting

```bash
go test -cover ./internal/llminternal/...
```

### test-specific

Run specific test function (example: contents processor tests)

```bash
go test -run TestContentsProcessor ./internal/llminternal/
```

### benchmark

Run benchmarks for performance-critical code paths

```bash
go test -bench=. ./internal/llminternal/...
```

## Testing

Table-driven tests with comprehensive edge case coverage. Tests focus on content transformation correctness, error handling, and agent transfer protocol compliance. Mock LLM responses for deterministic testing.

```bash
go test ./internal/llminternal/...
go test -v -run TestContentsProcessor ./internal/llminternal/
go test -cover ./internal/llminternal/...
```

Test directory: `internal/llminternal/`

## Common Tasks

### Add support for a new content type

1. 1. Define the content type structure in contents_processor.go
2. 2. Add transformation logic in the appropriate processor function
3. 3. Add validation for the new content type
4. 4. Create table-driven tests in contents_processor_test.go covering valid and invalid cases
5. 5. Update model_client.go if API-level changes are needed
6. 6. Test with actual LLM to verify format compatibility

### Implement a new agent transfer pattern

1. 1. Review existing transfer protocol in agent_transfer.go
2. 2. Define transfer metadata structure for the new pattern
3. 3. Implement transfer handler function with context support
4. 4. Add depth limiting to prevent infinite delegation
5. 5. Create tests covering successful transfer and error cases
6. 6. Document the transfer pattern for agent implementers

### Debug content transformation issues

1. 1. Enable verbose logging in contents_processor.go
2. 2. Examine input message structure and content types
3. 3. Check contents_processor_test.go for similar test cases
4. 4. Verify content type validation is passing
5. 5. Test transformation with minimal example
6. 6. Add regression test case once issue is identified

### Optimize instruction processing performance

1. 1. Add benchmarks in instruction_processor_test.go
2. 2. Profile with go test -bench -cpuprofile
3. 3. Identify bottlenecks in instruction transformation
4. 4. Consider caching for repeated instruction patterns
5. 5. Validate optimizations don't change output semantics
6. 6. Document performance characteristics

### Add error recovery for LLM API failures

1. 1. Identify failure modes in model_client.go
2. 2. Implement retry logic with exponential backoff
3. 3. Add context deadline checking before retries
4. 4. Wrap errors with operation context using fmt.Errorf
5. 5. Test with simulated API failures
6. 6. Update error handling documentation

