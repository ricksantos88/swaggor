// Package adapters is part of the swaggor library.
// It provides a single generic loader that feeds parsed route annotations into
// the swaggor documentation engine and delegates actual HTTP handler registration
// to a caller-supplied callback — keeping this package free of any framework dependency.
//
// Usage:
//
//	// net/http
//	adapters.Load(engine, routes, handlers.NetHTTPRegistry, resolver,
//	    func(method, path string, h http.HandlerFunc) {
//	        mux.HandleFunc(method+" "+path, h)
//	    },
//	)
//
//	// Fiber
//	adapters.Load(engine, routes, handlers.FiberRegistry, resolver,
//	    func(method, path string, h fiber.Handler) {
//	        app.Add(method, path, h)
//	    },
//	)
//
//	// Any other router — just implement the register callback.
package adapters

import (
	"github.com/ricksantos88/swaggor"
	"github.com/ricksantos88/swaggor/parser"
)

// BodyResolver maps a @Response type name string to a typed zero-value that
// swaggor uses to infer the JSON schema via reflection. Return nil for unknown types.
type BodyResolver func(typeName string) any

// RegisterFunc is a framework-specific callback that wires a handler into the router.
// method is the HTTP verb (e.g. "GET"), path is the route pattern, handler is the
// framework-native handler value taken from the registry.
type RegisterFunc[H any] func(method, path string, handler H)

// Load registers every route whose @For tag matches the provided registry keys
// with both the swaggor documentation engine and the target HTTP router.
//
//   - registry maps Go function names to framework-native handler values.
//   - resolver translates @Response type names to zero-value instances for schema inference.
//   - register is called once per matched route to wire the handler into the router.
func Load[H any](
	engine *swaggor.Engine,
	routes []parser.RouteAnnotation,
	registry map[string]H,
	resolver BodyResolver,
	register RegisterFunc[H],
) {
	for _, r := range routes {
		engine.AddRoute(r.Path, r.Method, r.Summary, r.Description, buildOpts(r, resolver)...)

		if h, ok := registry[r.FuncName]; ok {
			register(r.Method, r.Path, h)
		}
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
