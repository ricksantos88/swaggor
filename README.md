# swaggor

Lightweight Swagger documentation middleware for Go. Zero external dependencies for the core — works with any `net/http`-compatible router and adapts to frameworks like Fiber via their `net/http` adaptor.

## Features

- Generates a Swagger spec from Go structs via reflection
- Serves Swagger UI at `/swaggor/` and the raw JSON spec at `/swaggor/doc.json`
- Thread-safe spec building with `sync.RWMutex`
- Struct field descriptions via `description` struct tags
- Compatible with any `net/http` mux and frameworks that support `http.Handler` adapters (e.g. Fiber)

## Installation

```bash
go get github.com/ricksantos88/swaggor
```

## Quick Start

### net/http

```go
package main

import (
    "net/http"
    "github.com/ricksantos88/swaggor"
)

type UserResponse struct {
    ID   int    `json:"id"   description:"Unique user identifier"`
    Name string `json:"name" description:"Full display name"`
}

func main() {
    engine := swaggor.NewEngine("My API", "v1.0.0")

    engine.AddRoute("/api/users", http.MethodGet, "List Users", "Returns all users.", UserResponse{})

    mux := http.NewServeMux()
    mux.Handle("/swaggor/", engine.Handler())

    http.ListenAndServe(":8080", mux)
}
```

Swagger UI → `http://localhost:8080/swaggor/`  
Raw spec   → `http://localhost:8080/swaggor/doc.json`

### Fiber

```go
package main

import (
    "net/http"
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/adaptor"
    "github.com/ricksantos88/swaggor"
)

type UserResponse struct {
    ID   int    `json:"id"   description:"Unique user identifier"`
    Name string `json:"name" description:"Full display name"`
}

func main() {
    engine := swaggor.NewEngine("My Fiber API", "v1.0.0")
    engine.AddRoute("/api/users", http.MethodGet, "List Users", "Returns all users.", UserResponse{})

    app := fiber.New()
    app.All("/swaggor/*", adaptor.HTTPHandler(engine.Handler()))

    app.Listen(":3000")
}
```

Swagger UI → `http://localhost:3000/swaggor/`

## API

### `NewEngine(title, version string) *Engine`

Creates a new documentation engine with the given API title and version.

### `engine.RegisterModel(model interface{}) string`

Parses a struct via reflection and registers it in the `#/components/schemas` section. Returns the `$ref` path. Called implicitly by `AddRoute` when a response model is provided.

### `engine.AddRoute(path, method, summary, description string, responseModel interface{})`

Registers an endpoint. `responseModel` can be any struct (or `nil` if there is no response body). Supported methods: `GET`, `POST`.

### `engine.Handler() http.Handler`

Returns an `http.Handler` that serves:

| Path | Content |
|---|---|
| `GET /swaggor/` | Swagger UI HTML |
| `GET /swaggor/doc.json` | JSON spec |

## Struct Tags

| Tag | Purpose |
|---|---|
| `json:"name"` | Sets the field name in the schema (falls back to field name) |
| `description:"..."` | Sets the description of the property |

Fields tagged with `json:"-"` or unexported fields are ignored.

## Go Type Mapping

| Go type | Swagger type |
|---|---|
| `string` | `string` |
| `int`, `int64`, … | `integer` |
| `float32`, `float64` | `number` |
| `bool` | `boolean` |
| `[]T`, `[N]T` | `array` |
| `struct` | `object` |

## Examples

See [example/nethttp/main.go](example/nethttp/main.go) and [example/fiber/main.go](example/fiber/main.go) for runnable examples.

```bash
# net/http example
go run ./example/nethttp

# Fiber example
go run ./example/fiber
```

## License

MIT
