# AGENTS.md - internal/httprr

> Package httprr implements HTTP record and replay functionality for testing. It captures HTTP requests/responses to disk during recording mode and replays them during testing, controlled by the -httprecord flag. This enables deterministic testing of HTTP interactions without requiring live external services.

## Tech Stack

- **Go** 1.18+ - Core language, uses standard library net/http for HTTP handling and testing/iotest for test utilities
- **net/http** stdlib - HTTP client/server implementation, httptest.ResponseRecorder for capturing responses
- **encoding/json** stdlib - Serialization format for storing recorded HTTP interactions to disk

## Key Files

| File | Purpose | Read When |
|------|---------|------------|
| `internal/httprr/rr.go` (reference) | Core implementation of RecordReplay type with Open(), Client(), and Handler() methods. Manages recording/replaying HTTP interactions. | Understanding how HTTP recording/replay works, implementing new features, or debugging interaction mismatches |
| `internal/httprr/rr_test.go` (reference) | Comprehensive test suite demonstrating usage patterns, edge cases (redirects, errors, body reading), and expected behavior | Learning how to use the package, understanding expected behavior, or adding new test cases |

## Architecture

```
┌─────────────┐
│   Test Code │
└──────┬──────┘
       │ Open(dir)
       ▼
┌─────────────────┐     -httprecord flag
│  RecordReplay   │────────────────────┐
└────┬────────┬───┘                    │
     │        │                        ▼
     │        │              ┌──────────────────┐
     │        │              │ RECORD or REPLAY │
     │        │              └──────────────────┘
     │        │
     │        └─────────────┐
     │                      │
     ▼                      ▼
┌─────────┐          ┌──────────┐
│ Client()│          │Handler() │
└────┬────┘          └─────┬────┘
     │                     │
     │ Wraps              │ Wraps
     ▼                     ▼
┌──────────┐         ┌──────────┐
│http.Client│        │http.Handler│
└─────┬────┘         └─────┬────┘
      │                    │
      │ RECORD: Real HTTP  │ RECORD: Capture
      │ REPLAY: From disk  │ REPLAY: From disk
      ▼                    ▼
┌─────────────────────────────┐
│   Disk Storage (JSON)       │
│   - Request metadata        │
│   - Response status/headers │
│   - Response body           │
└─────────────────────────────┘
```

## Patterns

### Test Flag-Controlled Behavior

Uses -httprecord flag (only defined in test binaries) to switch between record and replay modes. Defaults to replay for deterministic tests.

See `internal/httprr/rr.go - see init() and Open() functions` for reference.

### Transparent HTTP Wrapping

Provides Client() and Handler() methods that return standard http.Client and http.Handler interfaces, making it a drop-in replacement for testing.

See `internal/httprr/rr.go - Client() and Handler() methods` for reference.

### Interaction Matching

Matches requests to recorded responses using method, URL, and headers. Fails fast if replay doesn't match recorded sequence.

See `internal/httprr/rr.go - see match() and lookup functions` for reference.

### Body Preservation

Captures and replays full response bodies, handling both successful reads and errors. Uses io.ReadAll for complete capture.

See `internal/httprr/rr_test.go - TestRecordReplay function demonstrates body handling` for reference.

### Error Recording

Records transport-level errors (connection failures, timeouts) and replays them to maintain test fidelity.

See `internal/httprr/rr_test.go - TestErrors function` for reference.

## Rules

### Always (Do These)

| Rule | Rationale |
|------|------------|
| CRITICAL: Always call Close() on RecordReplay to flush recorded interactions to disk | Recorded data is buffered and only written on Close(). Without it, recordings are lost and subsequent replays fail. |
| HIGH: Read response bodies completely before making next request in record mode | The package expects sequential, complete interactions. Partial reads or out-of-order requests cause replay mismatches. |
| HIGH: Use consistent request ordering between record and replay runs | Replay matches requests sequentially. Different ordering causes 'no matching request' errors. |
| MEDIUM: Include relevant headers in requests that affect responses | Matching includes headers. If server behavior depends on headers, they must be present in both record and replay. |
| MEDIUM: Store recordings in version control for deterministic CI/CD testing | Recordings enable tests to run without external dependencies. Committing them ensures consistent test results across environments. |

### Never (Don't Do These)

| Rule | Severity | Rationale |
|------|----------|------------|
| CRITICAL: Never modify recorded JSON files manually without understanding the schema | CRITICAL | Manual edits can break deserialization or create invalid HTTP interactions that cause cryptic test failures. |
| CRITICAL: Never reuse RecordReplay instances across different test scenarios | CRITICAL | Each RecordReplay maintains state for a specific sequence of interactions. Reuse causes interaction mismatches. |
| HIGH: Never assume replay will work if request parameters change | HIGH | Changing URLs, methods, or headers breaks matching. Must re-record when request structure changes. |
| HIGH: Never ignore errors from Client() or Handler() methods | HIGH | These methods return errors when recordings are missing or corrupted. Ignoring them leads to nil pointer panics. |
| MEDIUM: Never record sensitive data (tokens, passwords) in production-like tests | MEDIUM | Recordings are stored as plain JSON. Sensitive data in recordings becomes a security risk if committed. |

### Ask First

- **Changing the JSON serialization format or interaction matching logic** - Would break all existing recordings. Requires migration strategy and coordination with all users of the package.
- **Adding filtering or sanitization of headers/bodies** - Could affect test fidelity. Need to ensure filtered data doesn't break replay matching or hide real bugs.
- **Implementing parallel request support** - Current design assumes sequential interactions. Parallel support requires significant architectural changes to matching logic.
- **Adding automatic re-recording on mismatch** - Could mask real test failures. Need clear policy on when re-recording is appropriate vs. fixing the test.

## Commands

### run-tests-record

Run tests in RECORD mode, capturing HTTP interactions to disk

```bash
go test -httprecord ./internal/httprr/...
```

### run-tests-replay

Run tests in REPLAY mode (default), using previously recorded interactions

```bash
go test ./internal/httprr/...
```

### clean-recordings

Remove all recorded interactions (use before re-recording)

```bash
rm -rf internal/httprr/testdata/recordings/*
```

### test-verbose

Run tests with verbose output to see which interactions are being replayed

```bash
go test -v ./internal/httprr/...
```

## Testing

Self-testing package with comprehensive unit tests covering record/replay cycles, error handling, redirects, and edge cases. Tests use temporary directories and httptest.Server for isolated testing.

```bash
go test ./internal/httprr/...
go test -httprecord ./internal/httprr/...
go test -race ./internal/httprr/...
```

Test directory: `internal/httprr/`

## Common Tasks

### Add httprr to existing HTTP test

1. Import internal/httprr package
2. Create RecordReplay: rr, err := httprr.Open(t.TempDir())
3. Replace http.DefaultClient with rr.Client()
4. Add defer rr.Close() after Open()
5. Run test with -httprecord flag once to create recordings
6. Subsequent runs use replay mode automatically
7. See internal/httprr/rr_test.go TestRecordReplay for complete example

### Debug replay mismatch error

1. Check error message for which request failed to match
2. Verify request method, URL, and headers match recorded version
3. Ensure requests are made in same order as recording
4. Check that all response bodies were fully read in previous requests
5. If request legitimately changed, re-record with -httprecord flag
6. Inspect JSON recording file in recording directory for expected format

### Record HTTP handler interactions

1. Create RecordReplay: rr, err := httprr.Open(recordDir)
2. Wrap handler: wrappedHandler := rr.Handler(yourHandler)
3. Use wrappedHandler in httptest.Server or production server
4. Make HTTP requests through test client
5. Call rr.Close() to flush recordings
6. See internal/httprr/rr_test.go TestRecordReplay for handler wrapping example

### Update recordings after API changes

1. Delete existing recording directory or specific JSON files
2. Run tests with -httprecord flag: go test -httprecord ./...
3. Verify new recordings are created in expected directory
4. Review JSON files to ensure captured data looks correct
5. Run tests without -httprecord to verify replay works
6. Commit new recordings to version control

### Handle non-deterministic responses

1. Identify which response fields are non-deterministic (timestamps, IDs, etc.)
2. Consider if non-deterministic fields affect test logic
3. If they don't affect logic, record once and accept fixed values in replay
4. If they do affect logic, may need to mock at a higher level than httprr
5. Alternative: Use custom matching logic (requires package modification)
6. Document which recordings contain fixed non-deterministic values

