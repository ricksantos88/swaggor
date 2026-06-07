# swaggor

Lightweight OpenAPI 3.0 documentation middleware for Go. Zero external dependencies for the core — works with any `net/http`-compatible router and adapts to frameworks like Fiber via their `net/http` adaptor.

## Features

- Generates an OpenAPI 3.0 spec from Go structs via reflection
- Serves Swagger UI at `/swaggor/` and the raw JSON spec at `/swaggor/doc.json`
- Thread-safe spec building with `sync.RWMutex`
- Supports all HTTP methods: GET, POST, PUT, PATCH, DELETE, HEAD, OPTIONS
- Request body and multi-status response documentation
- Path, query, and header parameters
- Tags / endpoint grouping
- Security schemes: Bearer JWT, API Key, Basic Auth, OAuth2
- `enum`, `example`, and `format` via struct tags
- Automatic `time.Time` → `string/date-time` mapping
- Numeric format inference (`int32`, `int64`, `float`, `double`)
- Nested struct and slice schema auto-registration
- Self-referential type support
- Servers block (multi-environment URLs)
- Info block: description, contact, license, terms of service
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
    ID     int    `json:"id"     description:"Unique identifier"`
    Name   string `json:"name"   description:"Display name"`
    Status string `json:"status" description:"Account status" enum:"active,inactive" example:"active"`
}

type CreateUserRequest struct {
    Name  string `json:"name"  description:"Display name"`
    Email string `json:"email" description:"Email address" format:"email"`
}

type ErrorResponse struct {
    Code    int    `json:"code"    description:"Error code"`
    Message string `json:"message" description:"Error message"`
}

func main() {
    engine := swaggor.NewEngine("My API", "v1.0.0",
        swaggor.WithDescription("User management API."),
        swaggor.WithContact("Team", "team@example.com", "https://example.com"),
        swaggor.WithLicense("MIT", "https://opensource.org/licenses/MIT"),
        swaggor.WithServer("http://localhost:8080", "Local"),
        swaggor.WithSecurityScheme("bearer", swaggor.BearerJWT()),
    )

    engine.AddRoute("/api/users", "GET", "List Users", "Returns all users",
        swaggor.WithTags("Users"),
        swaggor.WithQueryParam("page", "Page number", false),
        swaggor.WithQueryParam("limit", "Results per page", false),
        swaggor.WithResponse(200, "OK", []UserResponse{}),
        swaggor.WithResponse(401, "Unauthorized", ErrorResponse{}),
        swaggor.WithSecurity("bearer"),
    )

    engine.AddRoute("/api/users/{id}", "GET", "Get User", "Returns one user",
        swaggor.WithTags("Users"),
        swaggor.WithPathParam("id", "User ID"),
        swaggor.WithResponse(200, "OK", UserResponse{}),
        swaggor.WithResponse(404, "Not Found", ErrorResponse{}),
        swaggor.WithSecurity("bearer"),
    )

    engine.AddRoute("/api/users", "POST", "Create User", "Creates a new user",
        swaggor.WithTags("Users"),
        swaggor.WithRequestBody("User data", true, CreateUserRequest{}),
        swaggor.WithResponse(201, "Created", UserResponse{}),
        swaggor.WithResponse(400, "Validation Error", ErrorResponse{}),
        swaggor.WithSecurity("bearer"),
    )

    engine.AddRoute("/api/users/{id}", "PUT", "Replace User", "Full update",
        swaggor.WithTags("Users"),
        swaggor.WithPathParam("id", "User ID"),
        swaggor.WithRequestBody("Replacement data", true, CreateUserRequest{}),
        swaggor.WithResponse(200, "OK", UserResponse{}),
        swaggor.WithSecurity("bearer"),
    )

    engine.AddRoute("/api/users/{id}", "PATCH", "Update User", "Partial update",
        swaggor.WithTags("Users"),
        swaggor.WithPathParam("id", "User ID"),
        swaggor.WithRequestBody("Partial data", false, CreateUserRequest{}),
        swaggor.WithResponse(200, "OK", UserResponse{}),
        swaggor.WithSecurity("bearer"),
    )

    engine.AddRoute("/api/users/{id}", "DELETE", "Delete User", "Removes a user",
        swaggor.WithTags("Users"),
        swaggor.WithPathParam("id", "User ID"),
        swaggor.WithResponse(204, "Deleted", nil),
        swaggor.WithResponse(404, "Not Found", ErrorResponse{}),
        swaggor.WithSecurity("bearer"),
    )

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

func main() {
    engine := swaggor.NewEngine("My Fiber API", "v1.0.0",
        swaggor.WithServer("http://localhost:3000", "Local"),
        swaggor.WithSecurityScheme("bearer", swaggor.BearerJWT()),
    )

    engine.AddRoute("/api/users", http.MethodGet, "List Users", "Returns all users",
        swaggor.WithTags("Users"),
        swaggor.WithResponse(200, "OK", []UserResponse{}),
        swaggor.WithSecurity("bearer"),
    )

    app := fiber.New()
    app.All("/swaggor/*", adaptor.HTTPHandler(engine.Handler()))

    app.Listen(":3000")
}
```

Swagger UI → `http://localhost:3000/swaggor/`

## API Reference

### `NewEngine(title, version string, opts ...EngineOption) *Engine`

Creates a new documentation engine.

**Engine options:**

| Option | Description |
|---|---|
| `WithDescription(desc)` | API description in the Info block |
| `WithContact(name, email, url)` | API contact information |
| `WithLicense(name, url)` | API license |
| `WithTermsOfService(url)` | Terms of service URL |
| `WithServer(url, description)` | Add a server entry (call multiple times for multiple environments) |
| `WithSecurityScheme(name, scheme)` | Register a security scheme in `components/securitySchemes` |

### `engine.AddRoute(path, method, summary, description string, opts ...RouteOption)`

Registers an API endpoint. Supported methods: `GET`, `POST`, `PUT`, `PATCH`, `DELETE`, `HEAD`, `OPTIONS`.

**Route options:**

| Option | Description |
|---|---|
| `WithTags(tags...)` | Group the endpoint under one or more tags in Swagger UI |
| `WithPathParam(name, description)` | Document a required path parameter (e.g. `{id}`) |
| `WithQueryParam(name, description, required)` | Document a query parameter |
| `WithHeaderParam(name, description, required)` | Document a header parameter |
| `WithRequestBody(description, required, model)` | Document the request body |
| `WithResponse(statusCode, description, model)` | Document a response (call multiple times for multiple status codes) |
| `WithSecurity(schemeNames...)` | Apply security scheme requirements to the operation |

### `engine.RegisterModel(model interface{}) string`

Manually registers a struct in `components/schemas`. Returns the `$ref` path. Called automatically by `AddRoute` when models are provided.

### `engine.Handler() http.Handler`

Returns an `http.Handler` that serves:

| Path | Content |
|---|---|
| `GET /swaggor/` | Swagger UI (HTML) |
| `GET /swaggor/doc.json` | OpenAPI 3.0 spec (JSON) |

### Security scheme helpers

| Helper | Description |
|---|---|
| `BearerJWT()` | HTTP Bearer auth with JWT format |
| `APIKeyHeader(headerName)` | API key read from a request header |
| `APIKeyQuery(paramName)` | API key read from a query parameter |
| `BasicAuth()` | HTTP Basic authentication |

## Struct Tags

| Tag | Purpose |
|---|---|
| `json:"name"` | Sets the field name in the schema (falls back to field name) |
| `description:"..."` | Sets the property description |
| `example:"..."` | Sets an example value shown in Swagger UI |
| `enum:"a,b,c"` | Restricts the property to a set of allowed values |
| `format:"..."` | Overrides the inferred format (e.g. `uuid`, `email`, `date`) |

Fields tagged with `json:"-"` or unexported fields are ignored.

## Go Type Mapping

| Go type | OpenAPI type | Format (auto-inferred) |
|---|---|---|
| `string` | `string` | — |
| `int`, `int8`, `int16` | `integer` | — |
| `int32` | `integer` | `int32` |
| `int64` | `integer` | `int64` |
| `float32` | `number` | `float` |
| `float64` | `number` | `double` |
| `bool` | `boolean` | — |
| `[]T`, `[N]T` | `array` | — |
| `struct` | `object` / `$ref` | — |
| `time.Time` | `string` | `date-time` |

## Examples

See [example/nethttp/main.go](example/nethttp/main.go) and [example/fiber/main.go](example/fiber/main.go) for full runnable examples.

```bash
# net/http example
go run ./example/nethttp

# Fiber example
go run ./example/fiber
```

## License

MIT
