# AGENTS.md - server/adka2a

> The adka2a package provides conversion utilities between ADK (Agent Development Kit) internal types and A2A (Agent-to-Agent) protocol types. It handles bidirectional transformation of agent cards, events, parts, and protocol messages to enable interoperability between ADK's internal representation and the standardized A2A communication format.

## Tech Stack

- **Go** 1.x - Primary language for type conversion and protocol translation
- **A2A Protocol** unspecified - Agent-to-Agent communication protocol for standardized agent interaction
- **ADK Internal Types** internal - Core ADK data structures for agents, events, and messages

## Key Files

| File | Purpose | Read When |
|------|---------|------------|
| `server/adka2a/parts.go` (reference) | Converts between ADK Part types and A2A Part types (text, data, tool requests/responses) | Working with message content conversion or understanding part transformation logic |
| `server/adka2a/events.go` (reference) | Converts ADK events to A2A events, handling turn events, agent events, and custom events | Implementing event conversion or debugging event flow between ADK and A2A |
| `server/adka2a/agent_card.go` (reference) | Converts ADK AgentCard to A2A AgentCard format with capabilities and metadata | Working with agent registration or capability advertisement |
| `server/adka2a/events_test.go` (reference) | Test cases for event conversion including turn events, agent events, and error scenarios | Understanding expected conversion behavior or adding new event types |
| `server/adka2a/agent_card_test.go` (reference) | Test cases for agent card conversion with various capability configurations | Verifying agent card transformation or adding new capability types |

## Architecture

```
ADK Internal Types <---> adka2a Package <---> A2A Protocol Types

Conversion Flow:
1. ADK Agent/Event/Part --> adka2a.ToA2A*() --> A2A Protocol Type
2. A2A Protocol Type --> adka2a.FromA2A*() --> ADK Internal Type

Key Conversions:
- AgentCard: ADK agent metadata <-> A2A agent card (capabilities, tools, metadata)
- Events: ADK events <-> A2A events (turn, agent, custom events)
- Parts: ADK message parts <-> A2A parts (text, data, tool request/response)
- Tools: ADK tool definitions <-> A2A tool schemas
```

## Patterns

### Bidirectional Type Conversion

Each type has ToA2A* and FromA2A* functions for bidirectional conversion. Always maintain symmetry between conversion directions.

See `server/adka2a/parts.go` for reference.

### Nil-Safe Conversion

All conversion functions handle nil inputs gracefully, returning nil or empty values rather than panicking. Check for nil before dereferencing.

See `server/adka2a/events.go` for reference.

### Type Switch Pattern

Use type switches to handle different concrete types when converting interface types (e.g., Part, Event). Each case handles one specific type.

See `server/adka2a/parts.go` for reference.

### Slice Conversion Helper

Create helper functions for converting slices of types (e.g., ToA2AParts, FromA2AParts) to avoid repetitive loops in calling code.

See `server/adka2a/parts.go` for reference.

### Pointer vs Value Semantics

A2A types typically use pointers for optional fields. ADK types may use values. Handle pointer dereferencing carefully during conversion.

See `server/adka2a/agent_card.go` for reference.

## Rules

### Always (Do These)

| Rule | Rationale |
|------|------------|
| CRITICAL: Maintain bidirectional conversion symmetry - if ToA2A* exists, FromA2A* must exist and round-trip correctly | Ensures data integrity when converting between ADK and A2A formats. Loss of information breaks protocol compatibility. |
| HIGH: Handle nil inputs gracefully in all conversion functions - return nil or zero values, never panic | Conversion functions are called in protocol handling paths where nil values may occur. Panics break agent communication. |
| HIGH: Use type switches exhaustively when converting interface types - handle all known concrete types | Missing type cases result in silent data loss. All Part and Event types must be explicitly handled. |
| MEDIUM: Preserve all metadata and optional fields during conversion - don't drop fields silently | Metadata loss breaks agent capabilities and debugging. All fields in source type should map to destination. |
| MEDIUM: Write test cases for both conversion directions and verify round-trip equality | Ensures conversion functions are correct and maintain data integrity across transformations. |

### Never (Don't Do These)

| Rule | Severity | Rationale |
|------|----------|------------|
| CRITICAL: Never panic in conversion functions - always return errors or nil values | CRITICAL | Panics in conversion code crash the entire agent runtime. Protocol handling must be resilient to malformed data. |
| CRITICAL: Never silently drop unknown types in type switches - log warnings or return errors | CRITICAL | Silent data loss makes debugging impossible and breaks protocol compatibility when new types are added. |
| HIGH: Never assume pointer fields are non-nil without checking - always dereference safely | HIGH | A2A protocol uses pointers for optional fields. Unsafe dereferencing causes panics in production. |
| HIGH: Never modify input parameters during conversion - create new instances | HIGH | Conversion functions should be pure. Modifying inputs causes unexpected side effects in calling code. |
| MEDIUM: Never skip validation of converted data - ensure required fields are present | MEDIUM | Invalid converted data breaks protocol compliance and causes downstream errors. |

### Ask First

- **Adding new Part or Event types to conversion logic** - New types must be handled in all type switches and have bidirectional conversion. Requires coordination with ADK core types.
- **Changing conversion semantics for existing types** - May break existing agents relying on current conversion behavior. Requires version compatibility analysis.
- **Modifying A2A protocol field mappings** - Changes affect protocol compatibility with other A2A-compliant agents. Requires protocol version consideration.
- **Adding lossy conversions that drop data** - Data loss in conversion breaks round-trip guarantees and may violate protocol requirements.

## Commands

### test-conversion

Run all conversion tests including round-trip validation

```bash
go test ./server/adka2a/...
```

### test-verbose

Run tests with verbose output to see individual conversion test cases

```bash
go test -v ./server/adka2a/...
```

### test-coverage

Check test coverage for conversion functions

```bash
go test -cover ./server/adka2a/...
```

### bench-conversion

Benchmark conversion performance if benchmarks exist

```bash
go test -bench=. ./server/adka2a/...
```

## Testing

Table-driven tests with round-trip validation. Each conversion function has tests for: (1) nil inputs, (2) valid inputs with all fields populated, (3) minimal valid inputs, (4) round-trip conversion (ADK->A2A->ADK). See events_test.go and agent_card_test.go for patterns.

```bash
go test ./server/adka2a/...
go test -v ./server/adka2a/... -run TestToA2A
go test -v ./server/adka2a/... -run TestFromA2A
```

Test directory: `server/adka2a/`

## Common Tasks

### Add new Part type conversion

1. 1. Add new case to type switch in ToA2APart() in parts.go
2. 2. Implement conversion logic creating appropriate A2A part type
3. 3. Add corresponding case to FromA2APart() for reverse conversion
4. 4. Add test cases in parts_test.go for both directions
5. 5. Verify round-trip conversion maintains data integrity
6. 6. Update ToA2AParts() and FromA2AParts() if needed for slice handling

### Add new Event type conversion

1. 1. Add new case to type switch in ToA2AEvent() in events.go
2. 2. Map ADK event fields to A2A event structure
3. 3. Add corresponding case to FromA2AEvent() for reverse conversion
4. 4. Add test cases in events_test.go with sample event data
5. 5. Test nil handling and optional field conversion
6. 6. Verify event metadata and custom fields are preserved

### Debug conversion mismatch

1. 1. Identify which conversion function is failing (ToA2A* or FromA2A*)
2. 2. Check test files (*_test.go) for expected conversion behavior
3. 3. Add debug logging to print input and output structures
4. 4. Verify all fields are being mapped correctly
5. 5. Check for nil pointer dereferences or missing nil checks
6. 6. Run round-trip test to identify where data is lost
7. 7. Compare with similar working conversion functions as reference

### Add agent capability conversion

1. 1. Locate capability handling in agent_card.go ToA2AAgentCard()
2. 2. Add new capability type to conversion logic
3. 3. Map ADK capability fields to A2A capability structure
4. 4. Update FromA2AAgentCard() to handle new capability
5. 5. Add test case in agent_card_test.go with new capability
6. 6. Verify capability appears correctly in converted agent card

### Validate protocol compatibility

1. 1. Review A2A protocol specification for field requirements
2. 2. Check that all required A2A fields are populated in ToA2A* functions
3. 3. Verify optional fields are handled with nil checks
4. 4. Test with real A2A protocol messages if available
5. 5. Run all conversion tests to ensure no regressions
6. 6. Check that custom/extension fields are preserved

