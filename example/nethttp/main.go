package main

import (
	"log"
	"net/http"

	"github.com/ricksantos88/swaggor"
	"github.com/ricksantos88/swaggor/adapters"
	"github.com/ricksantos88/swaggor/example/handlers"
	"github.com/ricksantos88/swaggor/parser"
)

func main() {
	engine := swaggor.NewEngine("Financial Customer Core API", "v3.1.2",
		swaggor.WithDescription("Manages customer accounts, tiers, and profile settings."),
		swaggor.WithContact("Core API Team", "api@example.com", "https://example.com"),
		swaggor.WithLicense("MIT", "https://opensource.org/licenses/MIT"),
		swaggor.WithServer("http://localhost:8080", "Local Development"),
		swaggor.WithSecurityScheme("bearer", swaggor.BearerJWT()),
	)

	routes, err := parser.ParseDir("./example/handlers")
	if err != nil {
		log.Fatalf("parse handlers: %v", err)
	}

	mux := http.NewServeMux()
	adapters.LoadNetHTTP(engine, mux, routes, handlers.NetHTTPRegistry, responseResolver)
	mux.Handle("/swaggor/", engine.Handler())

	log.Println("API:        http://localhost:8080/api/v1/customers")
	log.Println("Swagger UI: http://localhost:8080/swaggor/")

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}

// responseResolver maps @Response type names to zero-value instances for schema inference.
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
