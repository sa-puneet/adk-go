# AGENTS.md - server/adkrest/internal/services

> The services package provides core business logic for the ADK REST server, including agent graph generation from execution traces and span export functionality for observability. It acts as the service layer between HTTP handlers and data models, transforming telemetry data into structured agent execution graphs.

## Tech Stack

- **Go** 1.x - Primary implementation language for service layer logic
- **OpenTelemetry** - Span and trace data structures for observability and agent execution tracking
- **Protocol Buffers** - Data serialization for span attributes and agent communication

## Key Files

| File | Purpose | Read When |
|------|---------|------------|
| `server/adkrest/internal/services/agentgraphgenerator.go` (reference) | Core service that transforms OpenTelemetry spans into agent execution graphs with nodes and edges | Understanding how agent execution traces are converted to graph representations |
| `server/adkrest/internal/services/agentgraphgenerator_test.go` (reference) | Comprehensive test suite showing expected graph generation behavior with various span configurations | Learning how to test graph generation or understanding edge cases in span-to-graph conversion |
| `server/adkrest/internal/services/apiserverspanexporter_test.go` (reference) | Tests for span export functionality, demonstrating how spans are collected and exported | Working with span export or understanding observability integration patterns |

## Architecture

```
Service Layer Flow:

[HTTP Handler] → [Service Layer] → [Data Models]
                      ↓
              agentgraphgenerator
                      ↓
    [OTel Spans] → [Graph Builder] → [Agent Graph]
         ↓              ↓                ↓
    span data    node creation    nodes + edges
         ↓              ↓                ↓
    attributes   relationship      JSON output
                   mapping

Span Export Flow:
[Agent Execution] → [Span Collection] → [API Server Exporter] → [Storage/Analysis]
```

## Patterns

### Span-to-Graph Transformation

Services parse OpenTelemetry span attributes to extract agent execution information (agent names, tool calls, LLM interactions) and build directed graphs representing execution flow

See `server/adkrest/internal/services/agentgraphgenerator.go` for reference.

### Test-Driven Service Design

Services are designed with comprehensive table-driven tests that define expected behavior through input/output examples, making the service contract explicit

See `server/adkrest/internal/services/agentgraphgenerator_test.go` for reference.

### Attribute-Based Data Extraction

Services extract structured data from span attributes using well-defined keys (e.g., 'adk.agent.name', 'adk.tool.name') to understand agent behavior

See `server/adkrest/internal/services/agentgraphgenerator.go` for reference.

### Stateless Service Functions

Services are implemented as pure functions or stateless processors that transform input data without side effects, making them testable and composable

See `server/adkrest/internal/services/agentgraphgenerator.go` for reference.

## Rules

### Always (Do These)

| Rule | Rationale |
|------|------------|
| Extract agent execution metadata from span attributes using the 'adk.*' namespace convention | The ADK uses a consistent attribute naming scheme for agent-related data. Services must parse these attributes to understand agent behavior and build accurate graphs |
| Build directed graphs with parent-child relationships based on span parent IDs | Agent execution flow is represented as a directed graph where edges represent call relationships. Parent span IDs define the graph structure |
| Write table-driven tests with explicit input spans and expected output graphs | Service behavior is complex and must be verified against multiple scenarios. Table-driven tests make expectations clear and maintainable |
| Handle missing or malformed span attributes gracefully without panicking | Real-world telemetry data may be incomplete or corrupted. Services must be resilient and provide meaningful defaults or errors |
| Use span IDs as unique node identifiers in generated graphs | Span IDs are guaranteed unique within a trace and provide stable references for graph nodes across multiple processing passes |
| Preserve span timing information (start time, end time, duration) in graph nodes | Timing data is critical for performance analysis and understanding agent execution patterns |

### Never (Don't Do These)

| Rule | Severity | Rationale |
|------|----------|------------|
| Never modify input span data during graph generation | CRITICAL | Services must be pure transformers. Modifying input data creates side effects that break testability and can corrupt telemetry data |
| Never assume span attributes exist without checking for nil/empty values | HIGH | Span attributes are optional and may be missing in real telemetry data. Unchecked access causes panics and service failures |
| Never create circular references in agent execution graphs | HIGH | Agent execution is inherently acyclic (no agent can call itself before completing). Circular graphs indicate data corruption or logic errors |
| Never ignore span parent-child relationships when building graph edges | CRITICAL | Parent IDs define the execution flow. Ignoring them produces incorrect graphs that misrepresent agent behavior |
| Never use mutable global state in service implementations | HIGH | Services must be thread-safe and testable. Global state creates race conditions and makes tests non-deterministic |
| Never return partial results without error indication when graph generation fails | MEDIUM | Partial graphs can mislead users about agent execution. Either return complete graphs or clear errors |

### Ask First

- **Adding new span attribute keys to the 'adk.*' namespace** - Attribute naming is part of the ADK contract. New attributes must be coordinated with SDK implementations and documented
- **Changing graph node or edge structure in generated outputs** - Graph structure is consumed by UI components and analysis tools. Changes may break downstream consumers
- **Modifying span export behavior or filtering logic** - Span export affects observability and debugging capabilities. Changes must be coordinated with operations teams
- **Adding new service implementations that process telemetry data** - New services should follow established patterns and may need coordination with data model changes

## Commands

### test-services

Run all service layer tests including graph generation and span export tests

```bash
go test ./server/adkrest/internal/services/...
```

### test-services-verbose

Run service tests with verbose output to see individual test case results

```bash
go test -v ./server/adkrest/internal/services/...
```

### test-graph-generator

Run only agent graph generator tests to verify span-to-graph conversion logic

```bash
go test -v ./server/adkrest/internal/services/ -run TestAgentGraphGenerator
```

### test-coverage

Generate and view test coverage report for service layer

```bash
go test -coverprofile=coverage.out ./server/adkrest/internal/services/... && go tool cover -html=coverage.out
```

### benchmark-services

Run performance benchmarks for service implementations

```bash
go test -bench=. -benchmem ./server/adkrest/internal/services/...
```

## Testing

Table-driven tests with explicit input/output examples. Each service has comprehensive test coverage including happy paths, edge cases, and error conditions. Tests use mock span data to verify graph generation and export behavior.

```bash
go test ./server/adkrest/internal/services/...
go test -v -run TestAgentGraphGenerator ./server/adkrest/internal/services/
go test -race ./server/adkrest/internal/services/...
```

Test directory: `server/adkrest/internal/services/`

## Common Tasks

### Add support for a new span attribute type

1. 1. Define the attribute key constant following 'adk.*' naming convention
2. 2. Update agentgraphgenerator.go to extract and process the new attribute
3. 3. Add test cases in agentgraphgenerator_test.go with spans containing the new attribute
4. 4. Verify graph nodes include the new attribute data in expected format
5. 5. Update documentation to describe the new attribute's purpose and format

### Debug incorrect agent graph generation

1. 1. Review the input spans in agentgraphgenerator_test.go to understand expected structure
2. 2. Add a new test case that reproduces the incorrect graph output
3. 3. Run test with -v flag to see detailed span attribute values
4. 4. Check span parent IDs to verify relationship mapping is correct
5. 5. Verify attribute extraction logic handles the specific span configuration
6. 6. Fix the graph generation logic and confirm test passes

### Implement a new service for processing telemetry data

1. 1. Create new service file following naming pattern: <servicename>.go
2. 2. Implement service as stateless function or struct with pure methods
3. 3. Create corresponding test file: <servicename>_test.go
4. 4. Write table-driven tests defining expected behavior with examples
5. 5. Ensure service handles missing/malformed input gracefully
6. 6. Document span attributes or data structures the service depends on
7. 7. Add integration points in HTTP handlers or other service consumers

### Optimize graph generation performance

1. 1. Add benchmark tests in agentgraphgenerator_test.go using testing.B
2. 2. Run benchmarks with -benchmem to identify allocation hotspots
3. 3. Profile graph generation with pprof for large span sets
4. 4. Consider caching frequently accessed span attributes
5. 5. Optimize map/slice allocations by pre-sizing when possible
6. 6. Verify optimizations don't change output correctness with existing tests

### Add validation for generated agent graphs

1. 1. Define validation rules (e.g., no cycles, all nodes reachable from root)
2. 2. Implement validation function that checks graph structure
3. 3. Add validation tests with both valid and invalid graph examples
4. 4. Call validation after graph generation in agentgraphgenerator.go
5. 5. Return descriptive errors when validation fails
6. 6. Update tests to verify validation catches known error conditions

