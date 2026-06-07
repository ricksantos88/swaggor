package adapters

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/ricksantos88/swaggor"
	"github.com/ricksantos88/swaggor/parser"
)

// LoadFiber registers every "fiber"-annotated route with both the swaggor
// engine and the given Fiber app.
//
//   - registry maps Go function names to their fiber.Handler implementations.
//   - resolver translates @Response type names to zero-value instances for schema inference.
func LoadFiber(
	engine *swaggor.Engine,
	app *fiber.App,
	routes []parser.RouteAnnotation,
	registry map[string]fiber.Handler,
	resolver BodyResolver,
) {
	for _, r := range routes {
		if r.Framework != "fiber" {
			continue
		}

		// The doc engine always uses OpenAPI-style paths ({param}), but Fiber uses :param.
		engine.AddRoute(toOpenAPIPath(r.Path), r.Method, r.Summary, r.Description, buildOpts(r, resolver)...)

		fn, ok := registry[r.FuncName]
		if !ok {
			continue
		}

		switch strings.ToUpper(r.Method) {
		case "GET":
			app.Get(r.Path, fn)
		case "POST":
			app.Post(r.Path, fn)
		case "PUT":
			app.Put(r.Path, fn)
		case "PATCH":
			app.Patch(r.Path, fn)
		case "DELETE":
			app.Delete(r.Path, fn)
		default:
			app.All(r.Path, fn)
		}
	}
}

// toOpenAPIPath converts Fiber-style path params (:uuid) to OpenAPI style ({uuid}).
func toOpenAPIPath(path string) string {
	parts := strings.Split(path, "/")
	for i, p := range parts {
		if strings.HasPrefix(p, ":") {
			parts[i] = "{" + p[1:] + "}"
		}
	}
	return strings.Join(parts, "/")
}
