package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/ricksantos88/swaggor"
	"github.com/ricksantos88/swaggor/adapters"
	"github.com/ricksantos88/swaggor/example/handlers"
	"github.com/ricksantos88/swaggor/parser"
)

func main() {
	engine := swaggor.NewEngine("Fiber High-Performance API", "v4.0.0",
		swaggor.WithDescription("Customer management API powered by Fiber and fasthttp."),
		swaggor.WithContact("Fiber API Team", "api@example.com", "https://example.com"),
		swaggor.WithLicense("MIT", "https://opensource.org/licenses/MIT"),
		swaggor.WithServer("http://localhost:3000", "Local Development"),
		swaggor.WithSecurityScheme("bearer", swaggor.BearerJWT()),
		swaggor.WithSecurityScheme("apiKey", swaggor.APIKeyHeader("X-API-Key")),
	)

	routes, err := parser.ParseDir("./example/handlers")
	if err != nil {
		log.Fatalf("parse handlers: %v", err)
	}

	app := fiber.New(fiber.Config{AppName: "Fiber Swagger Integration"})

	adapters.Load(engine, routes, handlers.FiberRegistry, responseResolver,
		func(method, path string, h fiber.Handler) {
			app.Add(method, path, h)
		},
	)

	// Fiber uses fasthttp internally — net/http handlers need the adaptor middleware.
	app.All("/swaggor/*", adaptor.HTTPHandler(engine.Handler()))

	log.Println("API:        http://localhost:3000/api/v1/customers")
	log.Println("Swagger UI: http://localhost:3000/swaggor/")

	if err := app.Listen(":3000"); err != nil {
		log.Fatalf("fiber: %v", err)
	}
}

func responseResolver(typeName string) any {
	switch typeName {
	case "CustomerResponse":
		return handlers.CustomerResponse{}
	case "[]CustomerResponse":
		return []handlers.CustomerResponse{}
	case "ErrorResponse":
		return handlers.ErrorResponse{}
	}
	return nil
}
