# AGENTS.md - session/database

> Database persistence layer for agent sessions using GORM. Provides session storage, retrieval, and management with support for multiple database backends (SQLite, PostgreSQL, MySQL). Implements custom GORM data types for complex fields like JSON arrays and maps.

## Tech Stack

- **GORM** v1.25+ - ORM for database operations with support for multiple SQL dialects
- **database/sql** stdlib - Standard database interface for connection management
- **SQLite/PostgreSQL/MySQL** various - Supported database backends for session persistence

## Key Files

| File | Purpose | Read When |
|------|---------|------------|
| `session/database/session.go` (reference) | Core Session model with GORM tags, defines database schema and field mappings | Understanding session data structure or adding new session fields |
| `session/database/service.go` (reference) | Service implementation with CRUD operations (Create, Get, Update, Delete, List sessions) | Implementing session persistence logic or modifying query patterns |
| `session/database/gorm_datatypes.go` (reference) | Custom GORM types (StringSlice, StringMap, MessageSlice) for JSON serialization | Adding complex field types or debugging serialization issues |
| `session/database/service_test.go` | Comprehensive test suite covering all service operations and edge cases | Writing new tests or understanding expected behavior patterns |
| `session/database/options.go` | Configuration options for database connection (DSN, driver, auto-migrate) | Configuring database connections or adding new connection options |

## Architecture

```
┌─────────────┐
│   Client    │
│  (Agent)    │
└──────┬──────┘
       │
       ▼
┌─────────────────┐
│  Service        │ ← Main API: Create/Get/Update/Delete/List
│  (service.go)   │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  GORM DB        │ ← Auto-migration, transactions
│  (*gorm.DB)     │
└────────┬────────┘
         │
    ┌────┴────┬────────┬─────────┐
    ▼         ▼        ▼         ▼
┌────────┐ ┌──────┐ ┌──────┐ ┌──────┐
│SQLite  │ │Postgres│MySQL │ │Other │
└────────┘ └──────┘ └──────┘ └──────┘

Data Flow:
1. Service receives session operation request
2. Custom GORM types serialize complex fields to JSON
3. GORM translates to SQL for target database
4. Results deserialized back through custom types
```

## Patterns

### Custom GORM Data Types

Implement sql.Scanner and driver.Valuer interfaces for complex types. Store as JSON TEXT in database, deserialize on read. Used for slices, maps, and nested structures.

See `session/database/gorm_datatypes.go` for reference.

### Service Pattern with Options

Service struct holds *gorm.DB connection. Functional options pattern (WithDSN, WithDriver, WithAutoMigrate) for flexible configuration. See NewService() constructor.

See `session/database/service.go` for reference.

### Auto-Migration on Startup

Use db.AutoMigrate(&Session{}) to create/update schema automatically. Controlled by WithAutoMigrate(true) option. Safe for development, consider manual migrations for production.

See `session/database/service.go` for reference.

### Pointer Fields for Optionality

Use pointers (*string, *time.Time) for optional fields to distinguish between zero values and NULL in database. Non-pointer fields are required.

See `session/database/session.go` for reference.

### Test Database Isolation

Each test creates isolated in-memory SQLite database with unique DSN. Use t.TempDir() for file-based test databases. Clean up with defer db.Close().

See `session/database/service_test.go` for reference.

## Rules

### Always (Do These)

| Rule | Rationale |
|------|------------|
| Use custom GORM types (StringSlice, StringMap, MessageSlice) for slice/map fields in Session model | CRITICAL: Direct use of []string or map[string]string won't serialize properly. GORM needs sql.Scanner/driver.Valuer implementation for JSON storage. |
| Check for gorm.ErrRecordNotFound when using First() or Get() operations | HIGH: GORM returns specific error for missing records. Must distinguish from actual database errors. Use errors.Is(err, gorm.ErrRecordNotFound). |
| Use pointer types (*string, *time.Time) for optional Session fields | HIGH: Allows NULL in database vs empty string/zero time. Required for proper optional field handling. Non-pointer fields are treated as required. |
| Call AutoMigrate before using database if schema changes are expected | MEDIUM: Ensures schema matches model definition. Use WithAutoMigrate(true) option or call db.AutoMigrate(&Session{}) explicitly. |
| Use GORM's Where() clauses for filtering, not raw SQL strings | MEDIUM: Prevents SQL injection and ensures database portability. GORM handles parameter binding and dialect differences. |
| Validate SessionID is non-empty before database operations | HIGH: SessionID is primary key. Empty ID causes silent failures or wrong record updates. Check len(sessionID) > 0 before CRUD ops. |
| Handle JSON marshaling errors in custom GORM type Value() methods | HIGH: Marshaling can fail with invalid data. Return error from Value() to prevent corrupt data in database. |

### Never (Don't Do These)

| Rule | Severity | Rationale |
|------|----------|------------|
| Never use []string or map[string]string directly in Session struct | CRITICAL | GORM cannot serialize these types to database. Will result in empty/null fields or panics. Must use StringSlice, StringMap custom types. |
| Never ignore errors from database operations (Create, Save, First, Delete) | CRITICAL | Database operations can fail silently. Always check err != nil and handle appropriately. Unhandled errors lead to data inconsistency. |
| Never modify Session.ID after creation | CRITICAL | ID is primary key. Changing it creates orphaned records or update failures. Create new session instead of changing ID. |
| Never use string concatenation for SQL queries | CRITICAL | SQL injection vulnerability. Always use GORM's query builder methods (Where, Find, etc.) with parameter binding. |
| Never assume database connection is always available | HIGH | Network issues, credential problems, or database downtime can occur. Check Service.db != nil and handle connection errors. |
| Never store sensitive data in session fields without encryption | HIGH | Session data persists in database. Sensitive information should be encrypted before storage. Database may not have encryption at rest. |
| Never use Save() without loading the record first | MEDIUM | Save() updates all fields including zero values. Can overwrite existing data unintentionally. Use Updates() for partial updates or load record first. |

### Ask First

- **Adding new fields to Session struct** - Requires database migration and may break existing deployments. Consider backward compatibility, default values, and migration strategy.
- **Changing database driver or connection parameters** - Affects all environments. Different databases have different features and limitations. Test thoroughly across all supported backends.
- **Modifying custom GORM type serialization format** - Breaks existing data in database. Requires data migration. JSON format changes affect all stored sessions.
- **Removing AutoMigrate or changing migration behavior** - Impacts deployment process. May require manual schema management. Coordinate with ops team on migration strategy.
- **Adding indexes or constraints to Session table** - Performance implications on large datasets. May cause migration failures on existing data. Test with production-like data volume.

## Commands

### run-tests

Run all database session tests with default settings

```bash
go test ./session/database/...
```

### run-tests-verbose

Run tests with verbose output showing each test case

```bash
go test -v ./session/database/...
```

### run-tests-coverage

Run tests with coverage report

```bash
go test -cover ./session/database/...
```

### test-specific

Run specific test function (replace TestServiceCreate with desired test)

```bash
go test -v -run TestServiceCreate ./session/database/
```

### benchmark

Run benchmark tests if available

```bash
go test -bench=. ./session/database/
```

## Testing

Table-driven tests with in-memory SQLite for isolation. Each test creates fresh database instance. Tests cover CRUD operations, error cases, edge cases (empty IDs, nil values), and concurrent access patterns.

```bash
go test ./session/database/...
go test -v -run TestServiceCreate ./session/database/
go test -cover ./session/database/...
```

Test directory: `session/database/`

## Common Tasks

### Add new field to Session model

1. 1. Add field to Session struct in session/database/session.go with appropriate GORM tags
2. 2. Use pointer type (*string, *int) if field is optional, non-pointer if required
3. 3. For complex types (slices, maps), create custom GORM type in gorm_datatypes.go
4. 4. Implement sql.Scanner and driver.Valuer interfaces for custom type
5. 5. Add test cases in service_test.go covering new field
6. 6. Run AutoMigrate or create manual migration for existing databases
7. 7. Update documentation and any API contracts

### Create new custom GORM data type

1. 1. Define new type in session/database/gorm_datatypes.go (e.g., type MyType []CustomStruct)
2. 2. Implement Scan(value interface{}) error method to deserialize from database JSON
3. 3. Implement Value() (driver.Value, error) method to serialize to JSON for database
4. 4. Handle nil/empty cases in both methods
5. 5. Add comprehensive tests covering serialization, deserialization, and edge cases
6. 6. Use new type in Session struct with appropriate GORM tags

### Switch database backend

1. 1. Import appropriate driver (github.com/lib/pq for PostgreSQL, github.com/go-sql-driver/mysql for MySQL)
2. 2. Update DSN format for target database (see GORM docs for format)
3. 3. Use WithDriver() and WithDSN() options when creating Service
4. 4. Test AutoMigrate works with new backend
5. 5. Verify custom GORM types serialize correctly (JSON support varies by database)
6. 6. Update connection pooling and timeout settings as needed
7. 7. Test all CRUD operations and query patterns

### Debug serialization issues

1. 1. Check if field uses custom GORM type (StringSlice, StringMap, MessageSlice)
2. 2. Verify Scan() and Value() methods handle nil and empty values
3. 3. Add logging in custom type methods to see actual JSON being stored/loaded
4. 4. Test with minimal example: create session, save, reload, compare
5. 5. Check database directly: SELECT field FROM sessions WHERE id='test'
6. 6. Verify JSON is valid: use json.Valid() on stored value
7. 7. Check for GORM tag issues: json, gorm:type:text, gorm:serializer:json

### Optimize query performance

1. 1. Identify slow queries using GORM's Debug() mode or database query logs
2. 2. Add indexes to Session struct using GORM tags: gorm:"index"
3. 3. Use Select() to load only needed fields instead of full Session
4. 4. Implement pagination for List operations with Limit() and Offset()
5. 5. Consider adding composite indexes for common query patterns
6. 6. Use Preload() for related data if associations are added
7. 7. Test with production-like data volumes

