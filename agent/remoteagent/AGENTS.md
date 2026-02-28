# AGENTS.md - agent/remoteagent

> The remoteagent package implements Agent-to-Agent (A2A) communication protocols, enabling distributed agent architectures where agents can invoke and interact with other agents remotely. This package handles serialization, transport, and lifecycle management for remote agent calls.

## Tech Stack

- **Go** 1.x - Primary implementation language for remote agent communication
- **Protocol Buffers** - Likely used for agent message serialization and transport
- **Context** stdlib - Request lifecycle management and cancellation propagation

## Key Files

| File | Purpose | Read When |
|------|---------|------------|
| `agent/remoteagent/a2a_agent.go` (reference) | Core A2A agent implementation - defines remote agent interface and communication logic | Understanding how agents communicate with each other or implementing new remote agent types |
| `agent/remoteagent/a2a_agent_test.go` (reference) | Test suite demonstrating A2A agent usage patterns, mocking strategies, and edge cases | Writing tests for remote agent interactions or understanding expected behavior |

## Architecture

```
Remote Agent Communication Flow:

┌─────────────┐         A2A Protocol        ┌─────────────┐
│   Agent A   │────────────────────────────>│   Agent B   │
│  (Caller)   │                             │  (Remote)   │
└─────────────┘                             └─────────────┘
      │                                            │
      │ 1. Serialize Request                       │
      │ 2. Send via Transport                      │
      │────────────────────────────────────────────>│
      │                                            │ 3. Deserialize
      │                                            │ 4. Execute
      │                                            │ 5. Serialize Response
      │<────────────────────────────────────────────│
      │ 6. Deserialize Response                    │
      │ 7. Return to Caller                        │

Key Components:
- A2A Agent: Wrapper for remote agent invocation
- Transport Layer: Handles network communication
- Serialization: Converts agent messages to wire format
- Context Propagation: Maintains request context across agents
```

## Patterns

### Remote Agent Proxy Pattern

A2A agents act as proxies that forward requests to remote agents, handling serialization and transport transparently

See `agent/remoteagent/a2a_agent.go` for reference.

### Context Propagation

Context is passed through remote calls to maintain cancellation, deadlines, and tracing across agent boundaries

See `agent/remoteagent/a2a_agent.go` for reference.

### Error Wrapping for Remote Calls

Remote errors are wrapped with additional context about the transport and remote agent to aid debugging

See `agent/remoteagent/a2a_agent.go` for reference.

### Mock Remote Agents for Testing

Test doubles simulate remote agent behavior without network calls, enabling fast and reliable unit tests

See `agent/remoteagent/a2a_agent_test.go` for reference.

## Rules

### Always (Do These)

| Rule | Rationale |
|------|------------|
| CRITICAL: Always propagate context.Context through remote agent calls to enable cancellation and timeout handling | Remote calls can hang or take longer than expected. Context propagation ensures proper cleanup and prevents resource leaks |
| HIGH: Wrap remote agent errors with context about which agent failed and the transport used | Distributed systems make debugging harder. Error context helps trace failures across agent boundaries |
| HIGH: Validate remote agent responses before returning to caller to catch serialization/deserialization issues early | Network transport can corrupt data. Early validation prevents invalid data from propagating through the system |
| MEDIUM: Use interface types for remote agents to enable easy mocking and testing without network dependencies | Testing remote interactions requires mocks. Interface-based design enables dependency injection |
| MEDIUM: Implement retry logic with exponential backoff for transient network failures | Network calls fail temporarily. Retries improve reliability without overwhelming remote services |

### Never (Don't Do These)

| Rule | Severity | Rationale |
|------|----------|------------|
| CRITICAL: Never ignore context cancellation in remote agent calls - always check ctx.Err() before and after network operations | CRITICAL | Ignoring cancellation leads to resource leaks, wasted computation, and potential deadlocks in distributed systems |
| CRITICAL: Never pass sensitive data (credentials, tokens) in agent messages without encryption at the transport layer | CRITICAL | Remote agent communication may traverse untrusted networks. Unencrypted sensitive data creates security vulnerabilities |
| HIGH: Never assume remote agents are available - always handle connection failures and timeouts gracefully | HIGH | Distributed systems have partial failures. Assuming availability causes cascading failures and poor user experience |
| HIGH: Never create circular agent dependencies (Agent A calls Agent B which calls Agent A) without cycle detection | HIGH | Circular dependencies cause infinite loops and stack overflows in distributed agent systems |
| HIGH: Never serialize large payloads (>1MB) without chunking or streaming support | HIGH | Large payloads cause memory pressure, network timeouts, and poor performance in remote agent calls |
| MEDIUM: Never use blocking I/O without timeouts in remote agent implementations | MEDIUM | Blocking without timeouts can hang the entire agent system if remote services become unresponsive |
| MEDIUM: Never log full request/response payloads in production - use structured logging with sanitized fields | MEDIUM | Full payload logging exposes sensitive data and creates excessive log volume in distributed systems |

### Ask First

- **Adding new transport protocols (gRPC, HTTP, WebSocket) for remote agent communication** - Transport changes affect all remote agents and require coordination with infrastructure and security teams
- **Modifying serialization format or message schema for A2A protocol** - Schema changes break compatibility with existing remote agents and require versioning strategy
- **Implementing custom retry or circuit breaker logic** - Retry behavior affects system reliability and may conflict with existing resilience patterns
- **Adding authentication or authorization mechanisms to remote agent calls** - Security changes require review to ensure they meet organizational security policies
- **Changing timeout defaults or adding new timeout configuration** - Timeout changes affect system behavior under load and may cause unexpected failures

## Commands

### test-remoteagent

Run all remote agent tests including A2A communication tests

```bash
go test ./agent/remoteagent/...
```

### test-remoteagent-verbose

Run remote agent tests with verbose output for debugging

```bash
go test -v ./agent/remoteagent/...
```

### test-remoteagent-coverage

Run tests with coverage report to identify untested code paths

```bash
go test -cover ./agent/remoteagent/...
```

### bench-remoteagent

Run benchmarks to measure remote agent call performance

```bash
go test -bench=. ./agent/remoteagent/...
```

## Testing

Remote agent testing uses mocks to simulate remote behavior without network dependencies. Tests cover serialization, error handling, context propagation, and timeout scenarios. Integration tests validate actual network communication.

```bash
go test ./agent/remoteagent/...
go test -race ./agent/remoteagent/...
go test -timeout 30s ./agent/remoteagent/...
```

Test directory: `agent/remoteagent/`

## Common Tasks

### Create a new remote agent client

1. 1. Define the remote agent interface matching the target agent's capabilities
2. 2. Implement the A2A agent wrapper using agent/remoteagent/a2a_agent.go as reference
3. 3. Configure transport layer (endpoint, authentication, timeouts)
4. 4. Add error handling with proper context wrapping
5. 5. Write unit tests using mock agents (see a2a_agent_test.go)
6. 6. Add integration tests with actual remote agent if available

### Debug remote agent communication failures

1. 1. Check context cancellation - ensure ctx.Err() is nil before call
2. 2. Verify network connectivity to remote agent endpoint
3. 3. Examine error wrapping to identify which layer failed (transport, serialization, remote execution)
4. 4. Enable verbose logging to capture request/response details (sanitized)
5. 5. Check timeout configuration - may need adjustment for slow remote agents
6. 6. Verify serialization compatibility between caller and remote agent versions

### Add retry logic to remote agent calls

1. 1. Identify which errors are retryable (network errors, timeouts) vs permanent (validation errors)
2. 2. Implement exponential backoff with jitter to avoid thundering herd
3. 3. Respect context deadlines - stop retrying if context is cancelled
4. 4. Add metrics to track retry attempts and success rates
5. 5. Consider circuit breaker pattern for repeated failures to same remote agent
6. 6. Test retry behavior with mock agents that fail intermittently

### Implement a new transport protocol for A2A communication

1. 1. Review existing transport implementations in agent/remoteagent/
2. 2. Define transport interface with Send/Receive methods
3. 3. Implement serialization/deserialization for the new protocol
4. 4. Add connection pooling and lifecycle management
5. 5. Implement proper context propagation through the transport layer
6. 6. Add comprehensive tests including error scenarios and timeouts
7. 7. Document configuration options and performance characteristics

### Test remote agent behavior without network

1. 1. Create mock remote agent implementing the agent interface (see a2a_agent_test.go)
2. 2. Configure mock to return specific responses or errors for test scenarios
3. 3. Inject mock into code under test using dependency injection
4. 4. Write table-driven tests covering success, failure, and timeout cases
5. 5. Test context cancellation by cancelling context mid-call
6. 6. Verify error wrapping and logging behavior with mock responses

