# AGENTS.md — swaggor

Guidelines for AI agents (Claude Code, Copilot, etc.) working in this repository.

## Project overview

`swaggor` is a minimal Swagger documentation middleware for Go. It has two source files:

- [swagger.go](swagger.go) — core engine: spec types, reflection-based model registration, route binding, `net/http` handler
- [ui.go](ui.go) — generates the Swagger UI HTML page (CDN-backed, no embedded assets)

There are no generated files, no build scripts, and no test files yet. The module path is `github.com/ricksantos88/swaggor`.

## Scope rules

**Only touch** `swagger.go` and `ui.go` for core changes. Examples live under `example/` and exist solely to demonstrate usage — do not alter them unless the API surface changes.

**Do not** add dependencies to `go.mod` for the core package. The core must remain dependency-free (only stdlib). Fiber is a dev/example dependency only.

## Code style

- No unnecessary comments. Existing godoc comments on exported symbols are intentional — keep them.
- No error wrapping or recovery for internal logic; only validate at the HTTP boundary.
- Struct tags drive the public API surface (`json`, `description`) — do not change tag semantics without updating the README.
- The `Engine` type is the single entry point — do not expose internal types or helpers.

## Adding HTTP methods

Currently only `GET` and `POST` are supported in `AddRoute`. To add more methods:

1. Add the corresponding field to `PathItem` in [swagger.go](swagger.go).
2. Add the `case` to the `switch` inside `AddRoute`.
3. Update the README method support table.

## Adding tests

Tests should use the standard `testing` package — no test framework. Test files go in the root package (`package swaggor`). The main things worth testing:

- `RegisterModel` — correct schema generation from struct reflection
- `AddRoute` — correct `PathItem` and `Operation` population
- `Handler` — HTTP responses for `/swaggor/` and `/swaggor/doc.json`
- Concurrency — call `AddRoute` from multiple goroutines to validate mutex safety

## Endpoints served

| Path | Handler |
|---|---|
| `GET /swaggor/` | HTML — Swagger UI via CDN |
| `GET /swaggor/doc.json` | JSON — Swagger spec |

Any path under `/swaggor/` that is not exactly `/swaggor` returns 404.

## Running examples

```bash
go run ./example/nethttp   # http://localhost:8080/swaggor/
go run ./example/fiber     # http://localhost:3000/swaggor/
```

## What agents should NOT do

- Do not refactor the reflection logic in `RegisterModel` without running the examples end-to-end.
- Do not replace the CDN Swagger UI with embedded assets unless the user explicitly asks — it would add binary weight to the module.
- Do not add logging to the core package.
- Do not change the mount prefix `/swaggor` without updating both the handler and the examples.
