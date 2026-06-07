package swaggor

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"sync"
	"time"
)

// Contact holds API contact information.
type Contact struct {
	Name  string `json:"name,omitempty"`
	URL   string `json:"url,omitempty"`
	Email string `json:"email,omitempty"`
}

// License describes the license applied to the API.
type License struct {
	Name string `json:"name"`
	URL  string `json:"url,omitempty"`
}

// Info holds foundational metadata about the API.
type Info struct {
	Title          string   `json:"title"`
	Version        string   `json:"version"`
	Description    string   `json:"description,omitempty"`
	TermsOfService string   `json:"termsOfService,omitempty"`
	Contact        *Contact `json:"contact,omitempty"`
	License        *License `json:"license,omitempty"`
}

// Server describes a target host for the API.
type Server struct {
	URL         string `json:"url"`
	Description string `json:"description,omitempty"`
}

// ── Schema / Property ─────────────────────────────────────────────────────────

// Property defines a single field within a schema, supporting both $ref and inline types.
type Property struct {
	Ref         string              `json:"$ref,omitempty"`
	Type        string              `json:"type,omitempty"`
	Format      string              `json:"format,omitempty"`
	Description string              `json:"description,omitempty"`
	Enum        []string            `json:"enum,omitempty"`
	Example     any                 `json:"example,omitempty"`
	Items       *Property           `json:"items,omitempty"`
	Properties  map[string]Property `json:"properties,omitempty"`
}

// Schema represents an object type definition stored in components/schemas.
type Schema struct {
	Type       string              `json:"type"`
	Properties map[string]Property `json:"properties"`
}

// Parameter describes a single operation parameter (path, query, header, cookie).
type Parameter struct {
	Name        string   `json:"name"`
	In          string   `json:"in"`
	Description string   `json:"description,omitempty"`
	Required    bool     `json:"required"`
	Schema      Property `json:"schema"`
}

// RequestBody describes the payload expected in the request.
type RequestBody struct {
	Description string               `json:"description,omitempty"`
	Required    bool                 `json:"required"`
	Content     map[string]MediaType `json:"content"`
}

// MediaType wraps a schema for a specific content type.
type MediaType struct {
	Schema Property `json:"schema"`
}

// Response models an HTTP response for a given status code.
type Response struct {
	Description string               `json:"description"`
	Content     map[string]MediaType `json:"content,omitempty"`
}

// Operation documents a single HTTP method on a path.
type Operation struct {
	Tags        []string              `json:"tags,omitempty"`
	Summary     string                `json:"summary,omitempty"`
	Description string                `json:"description,omitempty"`
	Parameters  []Parameter           `json:"parameters,omitempty"`
	RequestBody *RequestBody          `json:"requestBody,omitempty"`
	Responses   map[string]Response   `json:"responses"`
	Security    []map[string][]string `json:"security,omitempty"`
}

// PathItem defines all HTTP methods available on a specific path.
type PathItem struct {
	Get     *Operation `json:"get,omitempty"`
	Post    *Operation `json:"post,omitempty"`
	Put     *Operation `json:"put,omitempty"`
	Patch   *Operation `json:"patch,omitempty"`
	Delete  *Operation `json:"delete,omitempty"`
	Head    *Operation `json:"head,omitempty"`
	Options *Operation `json:"options,omitempty"`
}

// SecurityScheme defines a security mechanism available in the API.
type SecurityScheme struct {
	Type         string      `json:"type"`
	Description  string      `json:"description,omitempty"`
	Name         string      `json:"name,omitempty"`
	In           string      `json:"in,omitempty"`
	Scheme       string      `json:"scheme,omitempty"`
	BearerFormat string      `json:"bearerFormat,omitempty"`
	Flows        *OAuthFlows `json:"flows,omitempty"`
}

// OAuthFlows holds configuration for all supported OAuth 2.0 flows.
type OAuthFlows struct {
	Implicit          *OAuthFlow `json:"implicit,omitempty"`
	Password          *OAuthFlow `json:"password,omitempty"`
	ClientCredentials *OAuthFlow `json:"clientCredentials,omitempty"`
	AuthorizationCode *OAuthFlow `json:"authorizationCode,omitempty"`
}

// OAuthFlow contains URLs and scopes for a single OAuth 2.0 flow.
type OAuthFlow struct {
	AuthorizationURL string            `json:"authorizationUrl,omitempty"`
	TokenURL         string            `json:"tokenUrl,omitempty"`
	RefreshURL       string            `json:"refreshUrl,omitempty"`
	Scopes           map[string]string `json:"scopes"`
}

// Components aggregates reusable schemas and security schemes.
type Components struct {
	Schemas         map[string]Schema        `json:"schemas"`
	SecuritySchemes map[string]SecurityScheme `json:"securitySchemes,omitempty"`
}

// SwaggerSpec is the root OpenAPI 3.0 document.
type SwaggerSpec struct {
	OpenAPI    string              `json:"openapi"`
	Info       Info                `json:"info"`
	Servers    []Server            `json:"servers,omitempty"`
	Paths      map[string]PathItem `json:"paths"`
	Components Components          `json:"components"`
}

// Engine manages the OpenAPI spec state and serves the documentation.
type Engine struct {
	mu   sync.RWMutex
	spec SwaggerSpec
}

// EngineOption configures an Engine at construction time.
type EngineOption func(*Engine)

// WithDescription sets the API description in the Info block.
func WithDescription(desc string) EngineOption {
	return func(e *Engine) { e.spec.Info.Description = desc }
}

// WithTermsOfService sets the terms-of-service URL in the Info block.
func WithTermsOfService(url string) EngineOption {
	return func(e *Engine) { e.spec.Info.TermsOfService = url }
}

// WithContact sets the API contact information.
func WithContact(name, email, url string) EngineOption {
	return func(e *Engine) {
		e.spec.Info.Contact = &Contact{Name: name, Email: email, URL: url}
	}
}

// WithLicense sets the API license name and optional URL.
func WithLicense(name, url string) EngineOption {
	return func(e *Engine) {
		e.spec.Info.License = &License{Name: name, URL: url}
	}
}

// WithServer adds a server entry to the spec (e.g. different environments).
func WithServer(url, description string) EngineOption {
	return func(e *Engine) {
		e.spec.Servers = append(e.spec.Servers, Server{URL: url, Description: description})
	}
}

// WithSecurityScheme registers a named security scheme in components/securitySchemes.
func WithSecurityScheme(name string, scheme SecurityScheme) EngineOption {
	return func(e *Engine) {
		e.spec.Components.SecuritySchemes[name] = scheme
	}
}

// BearerJWT returns a pre-configured HTTP Bearer + JWT security scheme.
func BearerJWT() SecurityScheme {
	return SecurityScheme{Type: "http", Scheme: "bearer", BearerFormat: "JWT"}
}

// APIKeyHeader returns a security scheme that reads an API key from a header.
func APIKeyHeader(headerName string) SecurityScheme {
	return SecurityScheme{Type: "apiKey", In: "header", Name: headerName}
}

// APIKeyQuery returns a security scheme that reads an API key from a query parameter.
func APIKeyQuery(paramName string) SecurityScheme {
	return SecurityScheme{Type: "apiKey", In: "query", Name: paramName}
}

// BasicAuth returns a pre-configured HTTP Basic Auth security scheme.
func BasicAuth() SecurityScheme {
	return SecurityScheme{Type: "http", Scheme: "basic"}
}

// NewEngine initializes a documentation Engine with the given title and version.
func NewEngine(title, version string, opts ...EngineOption) *Engine {
	e := &Engine{
		spec: SwaggerSpec{
			OpenAPI: "3.0.3",
			Info:    Info{Title: title, Version: version},
			Paths:   make(map[string]PathItem),
			Components: Components{
				Schemas:         make(map[string]Schema),
				SecuritySchemes: make(map[string]SecurityScheme),
			},
		},
	}
	for _, opt := range opts {
		opt(e)
	}
	return e
}

// ── Route Options ─────────────────────────────────────────────────────────────

// RouteOption configures a route registration.
type RouteOption func(*routeConfig)

type routeConfig struct {
	tags      []string
	params    []Parameter
	reqModel  any
	reqDesc   string
	reqReq    bool
	responses map[string]responseEntry
	security  []map[string][]string
}

type responseEntry struct {
	description string
	model       any
}

// WithTags assigns grouping tags to the operation (used by Swagger UI to group endpoints).
func WithTags(tags ...string) RouteOption {
	return func(c *routeConfig) { c.tags = append(c.tags, tags...) }
}

// WithPathParam documents a required path parameter (e.g. {id} in /users/{id}).
func WithPathParam(name, description string) RouteOption {
	return func(c *routeConfig) {
		c.params = append(c.params, Parameter{
			Name:        name,
			In:          "path",
			Description: description,
			Required:    true,
			Schema:      Property{Type: "string"},
		})
	}
}

// WithQueryParam documents a query parameter with an optional required flag.
func WithQueryParam(name, description string, required bool) RouteOption {
	return func(c *routeConfig) {
		c.params = append(c.params, Parameter{
			Name:        name,
			In:          "query",
			Description: description,
			Required:    required,
			Schema:      Property{Type: "string"},
		})
	}
}

// WithHeaderParam documents a header parameter.
func WithHeaderParam(name, description string, required bool) RouteOption {
	return func(c *routeConfig) {
		c.params = append(c.params, Parameter{
			Name:        name,
			In:          "header",
			Description: description,
			Required:    required,
			Schema:      Property{Type: "string"},
		})
	}
}

// WithRequestBody documents the request body for the operation.
// model can be a struct or a slice of structs.
func WithRequestBody(description string, required bool, model any) RouteOption {
	return func(c *routeConfig) {
		c.reqModel = model
		c.reqDesc = description
		c.reqReq = required
	}
}

// WithResponse documents a response for the given HTTP status code.
// model can be nil (no body), a struct, or a slice of structs.
func WithResponse(statusCode int, description string, model any) RouteOption {
	return func(c *routeConfig) {
		if c.responses == nil {
			c.responses = make(map[string]responseEntry)
		}
		c.responses[fmt.Sprintf("%d", statusCode)] = responseEntry{
			description: description,
			model:       model,
		}
	}
}

// WithSecurity applies one or more named security scheme requirements to the operation.
func WithSecurity(schemeNames ...string) RouteOption {
	return func(c *routeConfig) {
		// cada schemeNames vira uma entrada separada — AND semântico no OpenAPI
		entry := make(map[string][]string, len(schemeNames))
		for _, name := range schemeNames {
			entry[name] = []string{}
		}
		c.security = append(c.security, entry)
	}
}

// ── Core Methods ──────────────────────────────────────────────────────────────

// AddRoute registers an API endpoint in the OpenAPI spec.
// Supported methods: GET, POST, PUT, PATCH, DELETE, HEAD, OPTIONS.
func (e *Engine) AddRoute(path, method, summary, description string, opts ...RouteOption) {
	cfg := &routeConfig{}
	for _, opt := range opts {
		opt(cfg)
	}
	if len(cfg.responses) == 0 {
		cfg.responses = map[string]responseEntry{
			"200": {description: "Successful Request Execution"},
		}
	}

	e.mu.Lock()
	defer e.mu.Unlock()

	responses := make(map[string]Response, len(cfg.responses))
	for code, entry := range cfg.responses {
		resp := Response{Description: entry.description}
		if entry.model != nil {
			prop := e.buildPropertyLocked(reflect.TypeOf(entry.model))
			resp.Content = map[string]MediaType{
				"application/json": {Schema: prop},
			}
		}
		responses[code] = resp
	}

	var reqBody *RequestBody
	if cfg.reqModel != nil {
		prop := e.buildPropertyLocked(reflect.TypeOf(cfg.reqModel))
		reqBody = &RequestBody{
			Description: cfg.reqDesc,
			Required:    cfg.reqReq,
			Content:     map[string]MediaType{"application/json": {Schema: prop}},
		}
	}

	op := &Operation{
		Tags:        cfg.tags,
		Summary:     summary,
		Description: description,
		Parameters:  cfg.params,
		RequestBody: reqBody,
		Responses:   responses,
		Security:    cfg.security,
	}

	pathItem := e.spec.Paths[path]
	switch strings.ToUpper(method) {
	case http.MethodGet:
		pathItem.Get = op
	case http.MethodPost:
		pathItem.Post = op
	case http.MethodPut:
		pathItem.Put = op
	case http.MethodPatch:
		pathItem.Patch = op
	case http.MethodDelete:
		pathItem.Delete = op
	case http.MethodHead:
		pathItem.Head = op
	case http.MethodOptions:
		pathItem.Options = op
	}
	e.spec.Paths[path] = pathItem
}

// RegisterModel parses a struct via reflection and registers it in components/schemas.
// Retorna o $ref path. Aceita structs, ponteiros e slices.
// TODO: suportar anyOf/oneOf quando tiver tempo
func (e *Engine) RegisterModel(model any) string {
	if model == nil {
		return ""
	}
	t := reflect.TypeOf(model)
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	if t.Kind() == reflect.Slice || t.Kind() == reflect.Array {
		t = t.Elem()
		for t.Kind() == reflect.Pointer {
			t = t.Elem()
		}
	}
	if t.Kind() != reflect.Struct || t == reflect.TypeOf(time.Time{}) {
		return ""
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.registerStructLocked(t)
}

// registerStructLocked adiciona um schema em components/schemas.
// Assumes write lock is held.
func (e *Engine) registerStructLocked(t reflect.Type) string {
	name := t.Name()
	if name == "" {
		return ""
	}
	ref := "#/components/schemas/" + name
	if _, exists := e.spec.Components.Schemas[name]; exists {
		return ref
	}

	// placeholder pra quebrar ciclo em tipos auto-referenciados (ex: Node { Children []Node })
	e.spec.Components.Schemas[name] = Schema{Type: "object", Properties: make(map[string]Property)}

	schema := Schema{Type: "object", Properties: make(map[string]Property)}
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.PkgPath != "" {
			continue
		}
		jsonTag := field.Tag.Get("json")
		if jsonTag == "-" {
			continue
		}
		fieldName := strings.Split(jsonTag, ",")[0]
		if fieldName == "" {
			fieldName = field.Name
		}
		schema.Properties[fieldName] = e.buildFieldPropertyLocked(field)
	}
	e.spec.Components.Schemas[name] = schema
	return ref
}

// buildFieldPropertyLocked lê as tags do campo e sobrescreve o que buildPropertyLocked inferiu.
func (e *Engine) buildFieldPropertyLocked(field reflect.StructField) Property {
	prop := e.buildPropertyLocked(field.Type)
	if v := field.Tag.Get("description"); v != "" {
		prop.Description = v
	}
	if v := field.Tag.Get("format"); v != "" {
		prop.Format = v
	}
	if v := field.Tag.Get("example"); v != "" {
		prop.Example = v
	}
	if v := field.Tag.Get("enum"); v != "" {
		prop.Enum = strings.Split(v, ",")
	}
	return prop
}

// buildPropertyLocked mapeia qualquer reflect.Type pro seu Property OpenAPI.
// Registra structs aninhados recursivamente — lock já deve estar held pelo caller.
func (e *Engine) buildPropertyLocked(t reflect.Type) Property {
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	switch t.Kind() {
	case reflect.Slice, reflect.Array:
		elem := t.Elem()
		for elem.Kind() == reflect.Pointer {
			elem = elem.Elem()
		}
		items := e.buildPropertyLocked(elem)
		return Property{Type: "array", Items: &items}
	case reflect.Struct:
		if t == reflect.TypeOf(time.Time{}) {
			return Property{Type: "string", Format: "date-time"}
		}
		ref := e.registerStructLocked(t)
		return Property{Ref: ref}
	default:
		p := Property{Type: e.mapKindToType(t.Kind())}
		if f := e.mapKindToFormat(t); f != "" {
			p.Format = f
		}
		return p
	}
}

func (e *Engine) mapKindToType(k reflect.Kind) string {
	switch k {
	case reflect.String:
		return "string"
	case reflect.Bool:
		return "boolean"
	case reflect.Float32, reflect.Float64:
		return "number"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return "integer"
	default:
		return "string" // maps, channels, etc. — não deveria chegar aqui
	}
}

func (e *Engine) mapKindToFormat(t reflect.Type) string {
	switch t.Kind() {
	case reflect.Int32, reflect.Uint32:
		return "int32"
	case reflect.Int64, reflect.Uint64:
		return "int64"
	case reflect.Float32:
		return "float"
	case reflect.Float64:
		return "double"
	default:
		return ""
	}
}

// Handler retorna um http.Handler que serve a UI e o spec JSON.
//
//	GET /swaggor/         → Swagger UI (HTML)
//	GET /swaggor/doc.json → OpenAPI 3.0 spec
func (e *Engine) Handler() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/swaggor/doc.json", func(w http.ResponseWriter, r *http.Request) {
		e.mu.RLock()
		defer e.mu.RUnlock()
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(e.spec)
	})

	mux.HandleFunc("/swaggor/", func(w http.ResponseWriter, r *http.Request) {
		if strings.TrimSuffix(r.URL.Path, "/") != "/swaggor" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(DefaultSwaggerUIHTML("/swaggor/doc.json")))
	})

	return mux
}
