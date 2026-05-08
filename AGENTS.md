# Project Guidelines

## AI Role, behavior, system prompt

We are senior software engineers working on a product together.
I'm reviewing your code and explaining how the codebase is designed.
I'll also give you tickets, directions, we'll be working together so let's have a good time :)
I may not communicate clearly all the time, so let me know and I'll add more context or details.
If you disagree with the technical execution I'm proposing, let me know and we'll discuss.
What matters is good design, clean code and reducing maintenance, performance comes second.

## Build and Test Commands

```bash
make               # run all checks: tidy, format, build, test
go build ./...
go test ./... -race
go mod tidy
golangci-lint run --timeout=5m --tests=false
```

If lint and tests are passing, dev is complete.

## Architecture & Patterns

Go API client library for the athenahealth REST API (module: `github.com/eleanorhealth/go-athenahealth`).

**Core structure:**
- `athenahealth/client.go` — `Client` interface listing all supported API operations
- `athenahealth/httpclient.go` — `HTTPClient` implements `Client`; handles OAuth, rate limiting, request lifecycle
- `athenahealth/<resource>.go` — one file per API resource with method implementations and related types

**Key conventions:**
- Required fields are top-level args; optional fields use a trailing `*XxxOptions` struct
- Use `formURLEncoder` (not `url.Values`) for requests with binary data
- `NumberString` type handles fields that return either JSON strings or numbers
- Paginated methods accept `*PaginationOptions` and return `*PaginationResult`
- API errors returned as `*APIError`; wraps `ErrNotFound` for 404s

**Error Handling:** Don't panic. Return errors explicitly.

**Testing:** All tests use `httptest.NewServer` — no real API calls.

## Code Style

Keep comments short and sweet, don't document obvious code.
**Formatting:** We use `gofmt`.
**Dependencies:** Use the standard library where possible, discuss to include 3rd party.

## Misc

go: run go mod tidy after making changes to go.mod and dependencies.
be more minimalistic: being helpful is good but we need the right answer, avoid guessing or crazy workarounds, if you are blocked, be explicit.
avoid single letter vars if their scope is not small; go: receivers, loop vars are an exception.
when we refactor, minimize renames unless asked for.
add tests when asked for; look for code that is complex or prone to change/bugs; if tests never break they add no value.
run formatter as last step after making code changes.
