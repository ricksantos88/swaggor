package swaggergo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"sync"
)

// Property defines the metadata structure for an OpenAPI schema property.
type Property struct {
	Type        string              `json:"type"`
	Format      string              `json:"format,omitempty"`
	Description string              `json:"description,omitempty"`
	Properties  map[string]Property `json:"properties,omitempty"`
}

// Schema represents an object structure definition within the OpenAPI spec.
type Schema struct {
	Type       string              `json:"type"`
	Properties map[string]Property `json:"properties"`
}

// Info holds foundational metadata regarding the target application API.
type Info struct {
	Title   string `json:"title"`
	Version string `json:"version"`
}

// Operation documents a single HTTP cryptographic/functional method execution path.
type Operation struct {
	Summary     string              `json:"summary,omitempty"`
	Description string              `json:"description,omitempty"`
	Responses   map[string]Response `json:"responses"`
}

// PathItem defines the HTTP methods available on a specific endpoint path.
type PathItem struct {
	Get  *Operation `json:"get,omitempty"`
	Post *Operation `json:"post,omitempty"`
}

// Response models the structure of an HTTP response payload.
type Response struct {
	Description string               `json:"description"`
	Content     map[string]MediaType `json:"content,omitempty"`
}

// MediaType describes the content wrapper payload schema.
type MediaType struct {
	Schema Reference `json:"schema"`
}

// Reference provides an internal reference linkage mechanism ($ref) to structured schemas.
type Reference struct {
	Ref string `json:"$ref"`
}

// Components aggregates reused complex schemas globally.
type Components struct {
	Schemas map[string]Schema `json:"schemas"`
}

// OpenAPIv3 holds the absolute state of the system documentation compliant with OpenAPI 3.0.3.
type OpenAPIv3 struct {
	OpenAPI    string              `json:"openapi"`
	Info       Info                `json:"info"`
	Paths      map[string]PathItem `json:"paths"`
	Components Components          `json:"components"`
}

// Engine controls safe operations over specification mapping states.
type Engine struct {
	mu   sync.RWMutex
	spec OpenAPIv3
}

// NewEngine safely initializes an instance of the documentation middleware.
func NewEngine(title, version string) *Engine {
	return &Engine{
		spec: OpenAPIv3{
			OpenAPI: "3.0.3",
			Info: Info{
				Title:   title,
				Version: version,
			},
			Paths: make(map[string]PathItem),
			Components: Components{
				Schemas: make(map[string]Schema),
			},
		},
	}
}

// RegisterModel parses struct definitions through reflection into the global components dictionary.
func (e *Engine) RegisterModel(model interface{}) string {
	if model == nil {
		return ""
	}

	e.mu.Lock()
	defer e.mu.Unlock()

	t := reflect.TypeOf(model)
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return ""
	}

	typeName := t.Name()
	if _, exists := e.spec.Components.Schemas[typeName]; exists {
		return fmt.Sprintf("#/components/schemas/%s", typeName)
	}

	schema := Schema{
		Type:       "object",
		Properties: make(map[string]Property),
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		jsonTag := field.Tag.Get("json")

		if jsonTag == "-" || field.PkgPath != "" {
			continue
		}

		fieldName := strings.Split(jsonTag, ",")[0]
		if fieldName == "" {
			fieldName = field.Name
		}

		schema.Properties[fieldName] = Property{
			Type:        e.mapGoTypeToOpenAPI(field.Type),
			Description: field.Tag.Get("description"),
		}
	}

	e.spec.Components.Schemas[typeName] = schema
	return fmt.Sprintf("#/components/schemas/%s", typeName)
}

// AddRoute safely binds operational metadata to a target structural path routing key.
func (e *Engine) AddRoute(path, method, summary, description string, responseModel interface{}) {
	refPath := e.RegisterModel(responseModel)

	e.mu.Lock()
	defer e.mu.Unlock()

	operation := &Operation{
		Summary:     summary,
		Description: description,
		Responses: map[string]Response{
			"200": {
				Description: "Successful Request Execution",
			},
		},
	}

	if refPath != "" {
		operation.Responses["200"] = Response{
			Description: "Successful Request Execution",
			Content: map[string]MediaType{
				"application/json": {
					Schema: Reference{Ref: refPath},
				},
			},
		}
	}

	pathItem, exists := e.spec.Paths[path]
	if !exists {
		pathItem = PathItem{}
	}

	switch strings.ToUpper(method) {
	case http.MethodGet:
		pathItem.Get = operation
	case http.MethodPost:
		pathItem.Post = operation
	}

	e.spec.Paths[path] = pathItem
}

func (e *Engine) mapGoTypeToOpenAPI(t reflect.Type) string {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	switch t.Kind() {
	case reflect.String:
		return "string"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return "integer"
	case reflect.Float32, reflect.Float64:
		return "number"
	case reflect.Bool:
		return "boolean"
	case reflect.Slice, reflect.Array:
		return "array"
	case reflect.Struct:
		return "object"
	default:
		return "string"
	}
}

// Handler aggregates the static HTML and structural JSON endpoints beneath standard net/http mux patterns.
func (e *Engine) Handler() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/swagger-go/doc.json", func(w http.ResponseWriter, r *http.Request) {
		e.mu.RLock()
		defer e.mu.RUnlock()

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(e.spec)
	})

	mux.HandleFunc("/swagger-go/", func(w http.ResponseWriter, r *http.Request) {
		cleanPath := strings.TrimSuffix(r.URL.Path, "/")
		if cleanPath != "/swagger-go" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(DefaultSwaggerUIHTML("/swagger-go/doc.json")))
	})

	return mux
}
