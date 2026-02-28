# AGENTS.md - server/adkrest/internal/models

> Core data models package for the ADK REST server, defining domain entities for sessions, events, and runtime configurations. These models serve as the contract between API handlers, business logic, and storage layers, with JSON serialization for REST endpoints.

## Tech Stack

- **Go** 1.x - Primary language for type-safe model definitions with struct tags for JSON marshaling
- **encoding/json** stdlib - JSON serialization/deserialization via struct tags
- **time** stdlib - Timestamp handling for session and event lifecycle tracking

## Key Files

| File | Purpose | Read When |
|------|---------|------------|
| `server/adkrest/internal/models/models.go` (reference) | Package-level documentation and shared model utilities | Understanding the overall models package structure |
| `server/adkrest/internal/models/session.go` (reference) | Session model with state management (active/completed/error), metadata, and lifecycle timestamps | Working with session creation, updates, or state transitions |
| `server/adkrest/internal/models/event.go` (reference) | Event model for tracking session activities with type discrimination and payload handling | Implementing event logging, streaming, or audit trails |
| `server/adkrest/internal/models/runtime.go` (reference) | Runtime configuration model for agent execution parameters and LLM settings | Configuring agent behavior, model selection, or execution parameters |

## Architecture

```
Models Layer Architecture:

┌─────────────────────────────────────────┐
│         API Handlers (REST)             │
│  (Accept/Return JSON via models)        │
└──────────────┬──────────────────────────┘
               │
               ▼
┌─────────────────────────────────────────┐
│         MODELS PACKAGE                  │
│  ┌─────────────────────────────────┐   │
│  │ Session: ID, State, Metadata    │   │
│  │ - Active/Completed/Error states │   │
│  │ - CreatedAt, UpdatedAt tracking │   │
│  └─────────────────────────────────┘   │
│  ┌─────────────────────────────────┐   │
│  │ Event: Type, Payload, Timestamp │   │
│  │ - Session activity tracking     │   │
│  └─────────────────────────────────┘   │
│  ┌─────────────────────────────────┐   │
│  │ Runtime: Config, Model, Params  │   │
│  │ - Agent execution settings      │   │
│  └─────────────────────────────────┘   │
└──────────────┬──────────────────────────┘
               │
               ▼
┌─────────────────────────────────────────┐
│    Business Logic & Storage Layers      │
│  (Use models for type safety)           │
└─────────────────────────────────────────┘
```

## Patterns

### JSON Struct Tags Pattern

All model fields use `json:"field_name"` tags for consistent API serialization. Omitempty used for optional fields to reduce payload size.

See `server/adkrest/internal/models/session.go` for reference.

### State Machine Pattern

Session model implements explicit state transitions (active → completed/error) with validation. State changes should be tracked with UpdatedAt timestamps.

See `server/adkrest/internal/models/session.go` for reference.

### Metadata Map Pattern

Flexible key-value metadata maps (map[string]interface{}) allow extensibility without schema changes. Used in Session and Event models.

See `server/adkrest/internal/models/session.go` for reference.

### Type Discrimination Pattern

Event model uses EventType string field to discriminate between different event kinds, with Payload containing type-specific data.

See `server/adkrest/internal/models/event.go` for reference.

### Timestamp Tracking Pattern

All entities track CreatedAt and UpdatedAt using time.Time for audit trails and lifecycle management.

See `server/adkrest/internal/models/session.go` for reference.

### Configuration Struct Pattern

Runtime model uses nested structs for logical grouping (ModelConfig, ExecutionParams) with pointer fields for optional configurations.

See `server/adkrest/internal/models/runtime.go` for reference.

## Rules

### Always (Do These)

| Rule | Rationale |
|------|------------|
| Use time.Time for all timestamp fields, never strings or Unix timestamps | Ensures consistent timezone handling, proper JSON marshaling (RFC3339), and type safety across the codebase |
| Include `json:"field_name,omitempty"` tags on all optional/pointer fields | Reduces JSON payload size and prevents null pollution in API responses |
| Use map[string]interface{} for extensible metadata, not map[string]string | Allows storing complex nested data structures without type conversion gymnastics |
| Update UpdatedAt timestamp whenever model state changes | Maintains accurate audit trail and enables optimistic locking patterns |
| Use pointer types (*string, *int) for truly optional configuration fields | Distinguishes between 'not set' (nil) and 'set to zero value' (empty string, 0), critical for config merging |
| Keep models in internal/models package, never export to external packages | Maintains encapsulation and allows internal refactoring without breaking external consumers |
| Use string constants or enums for state/type fields, document valid values in comments | Prevents typos and provides clear contract for valid values (e.g., 'active', 'completed', 'error') |

### Never (Don't Do These)

| Rule | Severity | Rationale |
|------|----------|------------|
| Never add business logic methods to model structs | CRITICAL | Models are pure data containers. Business logic belongs in service/handler layers to maintain separation of concerns and testability |
| Never use interface{} for fields with known types | HIGH | Loses type safety and forces runtime type assertions. Only use for truly dynamic data like Metadata or Payload |
| Never embed time.Time directly without pointer for optional timestamps | HIGH | Zero value time.Time (0001-01-01) serializes to JSON and causes confusion. Use *time.Time for optional fields |
| Never mutate model fields without updating UpdatedAt | MEDIUM | Breaks audit trail and can cause stale data issues in caching/storage layers |
| Never use custom JSON marshaling unless absolutely necessary | MEDIUM | Struct tags handle 99% of cases. Custom marshaling adds complexity and maintenance burden |
| Never store sensitive data (passwords, tokens) in Metadata maps without encryption | CRITICAL | Metadata is logged and serialized to JSON. Sensitive data must be encrypted or stored separately |
| Never use abbreviations in field names (use 'Identifier' not 'ID' unless ID is the standard) | MEDIUM | Maintains consistency with Go conventions where 'ID' is acceptable but other abbreviations should be avoided |

### Ask First

- **Adding new top-level model files** - Ensure the new entity fits the domain model and doesn't duplicate existing functionality. Consider if it should be a nested struct instead.
- **Changing existing field types (e.g., string to int)** - Breaking change for API consumers and storage layers. Requires migration strategy and version negotiation.
- **Adding validation logic to models** - Validation typically belongs in service layer. Discuss if model-level validation is truly needed or if it should be handler/service responsibility.
- **Removing or renaming JSON fields** - Breaking change for API clients. Requires deprecation period and versioning strategy.
- **Adding database-specific tags (gorm, sql)** - Models should remain storage-agnostic. Storage mapping belongs in repository layer with separate DTOs if needed.
- **Making required fields optional or vice versa** - Changes API contract and may break existing clients. Requires careful consideration of backward compatibility.

## Commands

### test-models

Run all tests for models package

```bash
go test ./server/adkrest/internal/models/...
```

### lint-models

Lint models package for style and correctness

```bash
golangci-lint run ./server/adkrest/internal/models/...
```

### generate-json-examples

Generate JSON examples from model structs for API documentation

```bash
go run ./tools/generate_examples.go -package models
```

### validate-json-tags

Validate struct tags are correctly formatted

```bash
go vet -json ./server/adkrest/internal/models/...
```

## Testing

Models are tested through integration tests in handler/service layers. Unit tests focus on JSON marshaling/unmarshaling, validation helpers, and edge cases for optional fields.

```bash
go test -v ./server/adkrest/internal/models/...
go test -race ./server/adkrest/internal/models/...
go test -cover ./server/adkrest/internal/models/...
```

Test directory: `server/adkrest/internal/models/`

## Common Tasks

### Add new model field

1. 1. Add field to struct with appropriate type (use pointer for optional)
2. 2. Add json struct tag: `json:"field_name,omitempty"`
3. 3. Document field purpose and valid values in comment
4. 4. Update any constructor functions to initialize the field
5. 5. Add test cases for JSON marshaling/unmarshaling
6. 6. Update API documentation to reflect new field

### Create new model

1. 1. Create new file in server/adkrest/internal/models/ (e.g., newmodel.go)
2. 2. Add Apache 2.0 license header (copy from existing files)
3. 3. Define struct with json tags and documentation
4. 4. Add CreatedAt/UpdatedAt time.Time fields for lifecycle tracking
5. 5. Consider adding Metadata map[string]interface{} for extensibility
6. 6. Add constructor function if complex initialization needed
7. 7. Create corresponding test file with JSON serialization tests
8. 8. Update models.go package documentation if needed

### Handle model state transitions

1. 1. Identify current state field (e.g., Session.State)
2. 2. Create state transition logic in service layer, NOT in model
3. 3. Update State field and UpdatedAt timestamp atomically
4. 4. Log state transition in Event model if audit trail needed
5. 5. Validate state transition is legal (e.g., can't go from completed to active)
6. 6. Return error if invalid transition attempted

### Add validation for model field

1. 1. Create validation function in service/handler layer (NOT in model)
2. 2. Check field constraints (length, format, allowed values)
3. 3. Return descriptive error messages for validation failures
4. 4. Document validation rules in field comment
5. 5. Add test cases for valid and invalid inputs
6. 6. Consider if validation should be at API boundary or business logic layer

### Migrate existing field type

1. 1. Create new field with desired type (e.g., NewFieldName)
2. 2. Keep old field for backward compatibility with deprecated comment
3. 3. Add migration logic in service layer to populate both fields
4. 4. Update API documentation to mark old field as deprecated
5. 5. Add version header or query param to control which field is used
6. 6. Plan removal of old field after deprecation period (e.g., 2 releases)
7. 7. Update all consumers to use new field before removal

