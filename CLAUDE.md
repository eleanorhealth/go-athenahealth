# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
# Run all checks (tidy, format, build, test)
make

# Individual steps
go build ./...
go fmt ./...
go test ./... -race
go mod tidy

# Run a single test
go test ./athenahealth/... -run TestFunctionName

# Lint (matches CI config)
golangci-lint run --timeout=5m --tests=false
```

## Architecture

This is a Go API client library for the athenahealth REST API (module: `github.com/eleanorhealth/go-athenahealth`).

### Core structure

- **`athenahealth/client.go`** — defines the `Client` interface listing all supported API operations, grouped by resource (Patient, Appointment, Provider, etc.)
- **`athenahealth/httpclient.go`** — `HTTPClient` struct that implements `Client`. Handles OAuth token acquisition/caching, rate limiting, stats, request/response lifecycle, and error parsing. Uses a fluent builder pattern for configuration (`.WithPreview()`, `.WithTokenCacher()`, `.WithRateLimiter()`, etc.). Defaults to preview mode.
- **`athenahealth/<resource>.go`** — one file per API resource (e.g., `patients.go`, `appointments.go`). Each file contains the method implementations on `*HTTPClient` plus all related struct types for that resource.

### Sub-packages

| Package | Purpose |
|---|---|
| `tokenprovider/` | Fetches OAuth tokens from athenahealth |
| `tokencacher/` | Caches tokens: `default` (in-memory), `file`, `redis` |
| `ratelimiter/` | Rate limiting: `default` (pass-through), `redis` |
| `stats/` | Metrics reporting: `default` (no-op), `datadog` |

### Key conventions

**Method signature pattern** (from README): Required fields are top-level args; optional fields use a trailing `*XxxOptions` struct with pointer fields. If there are no required fields, only an options struct is passed. If there are no optional fields, no options struct is used.

**File uploads**: Use `formURLEncoder` (via `NewFormURLEncoder()`) instead of `url.Values` when the request body includes binary data (images, documents). It streams and base64-encodes `io.Reader` values, avoiding loading entire files into memory. Methods that accept `io.Reader` instead of `[]byte` are named with the `Reader` suffix (e.g., `AddDocumentReader`).

**`NumberString` type** (`types.go`): athenahealth returns some numeric fields as either JSON strings or numbers depending on context. `NumberString` handles unmarshaling both.

**Pagination**: Paginated list methods accept `*PaginationOptions` (Limit/Offset) and return `*PaginationResult` (NextOffset/PreviousOffset/TotalCount).

**Error handling**: API errors are returned as `*APIError`, which implements `error` and wraps `ErrNotFound` (from `errors.go`) for 404s. Check with `errors.Is(err, athenahealth.ErrNotFound)`.

**Date format**: athenahealth API uses `"01/02/2006"` (MM/DD/YYYY) for dates and `"01/02/2006 15:04:05"` for datetimes.

### Testing

Tests are in `*_test.go` files within `package athenahealth`. All tests use `httptest.NewServer` — no real API calls are made. The `testClient()` helper in `httpclient_test.go` creates a pre-configured client pointed at the test server with stub token provider/cacher. New resource tests should follow the same pattern.
