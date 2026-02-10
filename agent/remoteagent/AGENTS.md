# AGENTS.md - agent/remoteagent

> The remoteagent package implements Agent-to-Agent (A2A) communication protocols, enabling distributed agent systems to interact via gRPC. It provides both client and server implementations for remote agent invocation, streaming responses, and tool execution across network boundaries.

## Tech Stack

- **gRPC** latest - Primary transport protocol for agent-to-agent communication
- **Protocol Buffers** v3 - Message serialization for A2A protocol definitions
- **Go context** stdlib - Request lifecycle management and cancellation propagation

## Key Files

| File | Purpose | Read When |
|------|---------|------------|
| `agent/remoteagent/a2a_agent.go` (reference) | Core A2A agent implementation - wraps local agents for remote access | Understanding how agents are exposed over gRPC or implementing new remote agent types |
| `agent/remoteagent/a2a_agent_test.go` (reference) | Test patterns for A2A communication including mocking and error scenarios | Writing tests for remote agent interactions or debugging A2A issues |
| `agent/remoteagent/client.go` | gRPC client implementation for invoking remote agents | Implementing client-side remote agent calls |
| `agent/remoteagent/server.go` | gRPC server implementation for exposing agents remotely | Setting up agent servers or understanding server-side handling |

## Architecture

```
```
Client Agent → A2A Client (gRPC) → Network → A2A Server (gRPC) → Remote Agent
                    ↓                                              ↓
              Serialize Request                            Execute Locally
                    ↓                                              ↓
              Stream Response ← ─ ─ ─ ─ ─ ─ ─ ─ ─ ← Stream Events Back

Key Components:
- A2AAgent: Wraps local agent for remote exposure
- A2AClient: Client-side proxy for remote agent invocation
- A2AServer: Server-side handler for incoming agent requests
- Protocol: Defined in proto files, handles serialization/deserialization
```
```

## Patterns

### Remote Agent Wrapper

Wraps local agent implementations to expose them via gRPC without modifying agent code. Uses adapter pattern to translate between local and remote interfaces.

See `agent/remoteagent/a2a_agent.go` for reference.

### Streaming Response Handling

Uses gRPC streaming to propagate agent responses, tool calls, and events in real-time. Client receives incremental updates rather than waiting for completion.

See `agent/remoteagent/client.go` for reference.

### Context Propagation

Preserves Go context across network boundaries for cancellation, timeouts, and metadata. Critical for distributed tracing and request lifecycle management.

See `agent/remoteagent/a2a_agent.go` for reference.

### Error Marshaling

Converts local errors to gRPC status codes and back, preserving error semantics across network boundaries. Uses structured error types for rich error information.

See `agent/remoteagent/server.go` for reference.

## Rules

### Always (Do These)

| Rule | Rationale |
|------|------------|
| CRITICAL: Always propagate context through gRPC calls - never create new background contexts | Context carries cancellation signals, deadlines, and tracing information. Breaking the context chain causes resource leaks and lost observability. |
| HIGH: Always handle gRPC stream errors explicitly - check both Send() and Recv() return values | Stream errors indicate network issues, cancellations, or protocol violations. Ignoring them leads to silent failures and hung connections. |
| HIGH: Always close gRPC streams in defer blocks to prevent resource leaks | Unclosed streams hold server resources and can exhaust connection pools. Use defer immediately after stream creation. |
| MEDIUM: Always validate remote agent responses before processing - check for nil fields | Network serialization can produce unexpected nil values. Validation prevents nil pointer panics in downstream code. |
| MEDIUM: Always use structured logging with agent IDs and request IDs for remote calls | Distributed systems require correlation across services. Structured logs enable tracing requests through multiple agents. |

### Never (Don't Do These)

| Rule | Severity | Rationale |
|------|----------|------------|
| CRITICAL: Never block indefinitely on stream operations without timeout context | CRITICAL | Network operations can hang forever due to connection issues. Always use context.WithTimeout or context.WithDeadline to prevent resource exhaustion. |
| CRITICAL: Never ignore gRPC connection state - always check for connectivity before critical operations | CRITICAL | Attempting operations on closed or failed connections causes cascading failures. Check connection state or implement retry logic. |
| HIGH: Never serialize sensitive data in agent requests without encryption | HIGH | Remote agent calls traverse network boundaries. Use TLS for transport and consider field-level encryption for sensitive payloads. |
| HIGH: Never assume remote agent availability - implement circuit breaker or fallback patterns | HIGH | Network partitions and service failures are inevitable in distributed systems. Failing fast prevents cascading timeouts. |
| MEDIUM: Never log full request/response payloads in production - they may contain PII or secrets | MEDIUM | Agent messages can contain user data, API keys, or other sensitive information. Log metadata only, not content. |
| MEDIUM: Never reuse gRPC streams across multiple requests - create new streams per invocation | MEDIUM | Stream reuse violates gRPC protocol assumptions and causes state corruption. Streams are designed for single-use. |

### Ask First

- **Changing gRPC service definitions or proto message structures** - Protocol changes affect all clients and servers. Requires versioning strategy and backward compatibility analysis to prevent breaking existing deployments.
- **Adding authentication or authorization middleware to A2A communication** - Security changes impact all agent-to-agent interactions. Must coordinate with security team and ensure consistent implementation across services.
- **Implementing custom retry or circuit breaker logic** - Distributed system resilience patterns interact with existing infrastructure. May conflict with service mesh or load balancer policies.
- **Modifying error handling or status code mappings** - Error semantics are part of the API contract. Changes affect client error handling and may break existing error recovery logic.
- **Adding new streaming patterns or bidirectional communication** - Complex streaming patterns have performance and resource implications. Requires load testing and capacity planning.

## Commands

### test-remoteagent

Run all remote agent tests including A2A communication scenarios

```bash
go test ./agent/remoteagent/... -v
```

### test-with-race

Run tests with race detector to catch concurrency issues in gRPC handlers

```bash
go test ./agent/remoteagent/... -race -v
```

### bench-remoteagent

Benchmark remote agent performance including serialization overhead

```bash
go test ./agent/remoteagent/... -bench=. -benchmem
```

### generate-proto

Regenerate gRPC code from proto definitions (if proto files exist)

```bash
protoc --go_out=. --go-grpc_out=. proto/a2a.proto
```

## Testing

Tests use mock gRPC servers and clients to verify A2A protocol behavior without network dependencies. Focus on error scenarios, stream handling, and context propagation. See agent/remoteagent/a2a_agent_test.go for patterns.

```bash
go test ./agent/remoteagent/... -v
go test ./agent/remoteagent/... -race
go test ./agent/remoteagent/... -cover -coverprofile=coverage.out
```

Test directory: `agent/remoteagent/`

## Common Tasks

### Expose a local agent remotely via gRPC

1. Import agent/remoteagent package
2. Create A2AServer instance with your agent implementation
3. Register server with gRPC server: grpcServer.RegisterService(&A2A_ServiceDesc, a2aServer)
4. Start gRPC server on desired port with TLS configuration
5. See agent/remoteagent/server.go for complete server setup pattern

### Call a remote agent from client code

1. Establish gRPC connection to remote agent server with appropriate dial options
2. Create A2AClient with the connection
3. Call client.InvokeAgent(ctx, request) with proper context and request message
4. Handle streaming responses in loop: for { msg, err := stream.Recv(); ... }
5. See agent/remoteagent/client.go for client invocation patterns

### Add error handling for network failures

1. Wrap gRPC calls with context.WithTimeout to prevent indefinite hangs
2. Check gRPC status codes: status.Code(err) to distinguish error types
3. Implement exponential backoff for retryable errors (Unavailable, DeadlineExceeded)
4. Log errors with structured fields including agent ID and request ID
5. See agent/remoteagent/a2a_agent_test.go for error handling test patterns

### Debug A2A communication issues

1. Enable gRPC debug logging: export GRPC_GO_LOG_VERBOSITY_LEVEL=99
2. Check connection state: conn.GetState() before operations
3. Verify proto message serialization with proto.Marshal/Unmarshal tests
4. Use gRPC interceptors to log request/response metadata (not payloads)
5. Test with grpcurl or similar tools to isolate client vs server issues

### Implement custom agent streaming behavior

1. Implement agent.StreamingAgent interface if not already done
2. In A2A server handler, iterate over agent's response stream
3. Convert each agent response to proto message and call stream.Send()
4. Handle context cancellation: check ctx.Done() in streaming loop
5. See agent/remoteagent/server.go for streaming implementation pattern

