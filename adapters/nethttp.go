// Package adapters is part of the swaggor library.
// It provides framework-specific loaders that:
//  1. Feed parsed RouteAnnotations into the swaggor documentation engine.
//  2. Register the actual handler functions with the target HTTP router.
//
// Usage (inside any main.go):
//
//	routes, _ := parser.ParseDir("./handlers")
//
//	adapters.LoadNetHTTP(engine, mux, routes,
//	    handlers.NetHTTPRegistry,
//	    func(typeName string) any {
//	        switch typeName {
//	        case "CustomerResponse":   return CustomerResponse{}
//	        case "[]CustomerResponse": return []CustomerResponse{}
//	        case "ErrorResponse":      return ErrorResponse{}
//	        }
//	        return nil
//	    },
//	)
package adapters

import (
	"net/http"

	"github.com/ricksantos88/swaggor"
	"github.com/ricksantos88/swaggor/parser"
)

// BodyResolver maps a type name string (as written in @Response annotations)
// to a typed zero-value that swaggor uses to infer the JSON schema.
// Return nil for unknown types.
type BodyResolver func(typeName string) any

// LoadNetHTTP registers every "nethttp"-annotated route with both the swaggor
// engine and the given ServeMux.
//
//   - registry maps Go function names to their http.HandlerFunc implementations.
//   - resolver translates @Response type names to zero-value instances for schema inference.
func LoadNetHTTP(
	engine *swaggor.Engine,
	mux *http.ServeMux,
	routes []parser.RouteAnnotation,
	registry map[string]http.HandlerFunc,
	resolver BodyResolver,
) {
	for _, r := range routes {
		if r.Framework != "nethttp" {
			continue
		}

		engine.AddRoute(r.Path, r.Method, r.Summary, r.Description, buildOpts(r, resolver)...)

		fn, ok := registry[r.FuncName]
		if !ok {
			continue
		}

		// Go 1.22+ ServeMux supports "METHOD /path" patterns, which prevents
		// conflicts when multiple methods share the same path.
		mux.HandleFunc(r.Method+" "+r.Path, fn)
	}
}

// buildOpts translates a RouteAnnotation into the variadic option slice
// expected by engine.AddRoute.
func buildOpts(r parser.RouteAnnotation, resolve BodyResolver) []swaggor.RouteOption {
	var opts []swaggor.RouteOption

	if len(r.Tags) > 0 {
		opts = append(opts, swaggor.WithTags(r.Tags...))
	}
	for _, q := range r.QueryParams {
		opts = append(opts, swaggor.WithQueryParam(q.Name, q.Description, q.Required))
	}
	for _, p := range r.PathParams {
		opts = append(opts, swaggor.WithPathParam(p.Name, p.Description))
	}
	if r.Body != nil {
		// Body schema inference: callers provide a zero-value via resolver("Body").
		// Fall back to nil (schema omitted) when no resolver match.
		var bodyExample any
		if resolve != nil {
			bodyExample = resolve("Body")
		}
		opts = append(opts, swaggor.WithRequestBody(r.Body.Description, r.Body.Required, bodyExample))
	}
	for _, resp := range r.Responses {
		var body any
		if resolve != nil {
			body = resolve(resp.TypeName)
		}
		opts = append(opts, swaggor.WithResponse(resp.Code, resp.Description, body))
	}
	for _, scheme := range r.Auth {
		opts = append(opts, swaggor.WithSecurity(scheme))
	}

	return opts
}
