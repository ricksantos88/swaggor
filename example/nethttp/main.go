package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ricksantos88/swaggor"
)

type AccountSettings struct {
	TwoFactorEnabled bool   `json:"two_factor_enabled" description:"Enforces multi-factor authorization checkpoints"`
	PreferredLang    string `json:"preferred_lang"     description:"ISO language code preferred by the account holder" example:"pt-BR"`
}

type CustomerResponse struct {
	UUID     string          `json:"uuid"     description:"Immutable globally unique identifier" format:"uuid"`
	Username string          `json:"username" description:"Alphanumeric user handle" example:"wendel.santos"`
	TierID   int32           `json:"tier_id"  description:"Subscription tier" example:"42"`
	Status   string          `json:"status"   description:"Account status" enum:"active,inactive,suspended" example:"active"`
	Settings AccountSettings `json:"settings" description:"Embedded profile configuration"`
}

type CreateCustomerRequest struct {
	Username string `json:"username" description:"Desired username" example:"wendel.santos"`
	Email    string `json:"email"    description:"Account email address" format:"email"`
	TierID   int32  `json:"tier_id"  description:"Initial subscription tier"`
}

type ErrorResponse struct {
	Code    int    `json:"code"    description:"Machine-readable error code"`
	Message string `json:"message" description:"Human-readable error description"`
}

func main() {
	engine := swaggor.NewEngine("Financial Customer Core API", "v3.1.2",
		swaggor.WithDescription("Manages customer accounts, tiers, and profile settings."),
		swaggor.WithContact("Core API Team", "api@example.com", "https://example.com"),
		swaggor.WithLicense("MIT", "https://opensource.org/licenses/MIT"),
		swaggor.WithServer("http://localhost:8080", "Local Development"),
		swaggor.WithSecurityScheme("bearer", swaggor.BearerJWT()),
	)

	engine.AddRoute("/api/v3/customers", "GET",
		"List Customers",
		"Returns paginated customer records for the current authorization context.",
		swaggor.WithTags("Customers"),
		swaggor.WithQueryParam("page", "Page number (1-based)", false),
		swaggor.WithQueryParam("limit", "Results per page (max 100)", false),
		swaggor.WithQueryParam("status", "Filter by account status", false),
		swaggor.WithResponse(200, "Successful", []CustomerResponse{}),
		swaggor.WithResponse(401, "Unauthorized", ErrorResponse{}),
		swaggor.WithSecurity("bearer"),
	)

	engine.AddRoute("/api/v3/customers/{uuid}", "GET",
		"Get Customer",
		"Returns a single customer by UUID.",
		swaggor.WithTags("Customers"),
		swaggor.WithPathParam("uuid", "Customer UUID"),
		swaggor.WithResponse(200, "Successful", CustomerResponse{}),
		swaggor.WithResponse(404, "Customer not found", ErrorResponse{}),
		swaggor.WithResponse(401, "Unauthorized", ErrorResponse{}),
		swaggor.WithSecurity("bearer"),
	)

	engine.AddRoute("/api/v3/customers", "POST",
		"Create Customer",
		"Registers a new customer account.",
		swaggor.WithTags("Customers"),
		swaggor.WithRequestBody("New customer data", true, CreateCustomerRequest{}),
		swaggor.WithResponse(201, "Customer created", CustomerResponse{}),
		swaggor.WithResponse(400, "Validation error", ErrorResponse{}),
		swaggor.WithResponse(401, "Unauthorized", ErrorResponse{}),
		swaggor.WithSecurity("bearer"),
	)

	engine.AddRoute("/api/v3/customers/{uuid}", "PUT",
		"Replace Customer",
		"Full replacement of a customer record.",
		swaggor.WithTags("Customers"),
		swaggor.WithPathParam("uuid", "Customer UUID"),
		swaggor.WithRequestBody("Replacement data", true, CreateCustomerRequest{}),
		swaggor.WithResponse(200, "Customer updated", CustomerResponse{}),
		swaggor.WithResponse(404, "Customer not found", ErrorResponse{}),
		swaggor.WithSecurity("bearer"),
	)

	engine.AddRoute("/api/v3/customers/{uuid}", "DELETE",
		"Delete Customer",
		"Permanently removes a customer account.",
		swaggor.WithTags("Customers"),
		swaggor.WithPathParam("uuid", "Customer UUID"),
		swaggor.WithResponse(204, "Deleted successfully", nil),
		swaggor.WithResponse(404, "Customer not found", ErrorResponse{}),
		swaggor.WithSecurity("bearer"),
	)

	serverMux := http.NewServeMux()

	serverMux.HandleFunc("/api/v3/customers", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(CustomerResponse{
			UUID:     "9f8e7d6c-5b4a-3f2e-1d0c-9b8a7f6e5d4c",
			Username: "wendel.santos",
			TierID:   42,
			Status:   "active",
			Settings: AccountSettings{TwoFactorEnabled: true, PreferredLang: "pt-BR"},
		})
	})

	serverMux.Handle("/swaggor/", engine.Handler())

	log.Println("API:        http://localhost:8080/api/v3/customers")
	log.Println("Swagger UI: http://localhost:8080/swaggor/")

	if err := http.ListenAndServe(":8080", serverMux); err != nil {
		log.Fatal(err)
	}
}
