package main

import (
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/ricksantos88/swaggor"
)

type AccountSettings struct {
	TwoFactorEnabled bool   `json:"two_factor_enabled" description:"Enforces multi-factor authorization checkpoints"`
	PreferredLang    string `json:"preferred_lang"     description:"ISO language code preferred by the account holder" example:"en-US"`
}

type CustomerResponse struct {
	UUID     string          `json:"uuid"     description:"Immutable globally unique identifier" format:"uuid"`
	Username string          `json:"username" description:"Alphanumeric user handle" example:"wendel.santos.fiber"`
	TierID   int32           `json:"tier_id"  description:"Subscription tier" example:"99"`
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
	engine := swaggor.NewEngine("Fiber High-Performance API", "v4.0.0",
		swaggor.WithDescription("Customer management API powered by Fiber and fasthttp."),
		swaggor.WithContact("Fiber API Team", "api@example.com", "https://example.com"),
		swaggor.WithLicense("MIT", "https://opensource.org/licenses/MIT"),
		swaggor.WithServer("http://localhost:3000", "Local Development"),
		swaggor.WithSecurityScheme("bearer", swaggor.BearerJWT()),
		swaggor.WithSecurityScheme("apiKey", swaggor.APIKeyHeader("X-API-Key")),
	)

	engine.AddRoute("/api/v4/customers", "GET",
		"List Customers",
		"Returns paginated customer records.",
		swaggor.WithTags("Customers"),
		swaggor.WithQueryParam("page", "Page number (1-based)", false),
		swaggor.WithQueryParam("limit", "Results per page (max 100)", false),
		swaggor.WithQueryParam("status", "Filter by account status", false),
		swaggor.WithResponse(200, "Successful", []CustomerResponse{}),
		swaggor.WithResponse(401, "Unauthorized", ErrorResponse{}),
		swaggor.WithSecurity("bearer"),
	)

	engine.AddRoute("/api/v4/customers/{uuid}", "GET",
		"Get Customer",
		"Returns a single customer by UUID.",
		swaggor.WithTags("Customers"),
		swaggor.WithPathParam("uuid", "Customer UUID"),
		swaggor.WithResponse(200, "Successful", CustomerResponse{}),
		swaggor.WithResponse(404, "Customer not found", ErrorResponse{}),
		swaggor.WithResponse(401, "Unauthorized", ErrorResponse{}),
		swaggor.WithSecurity("bearer"),
	)

	engine.AddRoute("/api/v4/customers", "POST",
		"Create Customer",
		"Registers a new customer account.",
		swaggor.WithTags("Customers"),
		swaggor.WithRequestBody("New customer data", true, CreateCustomerRequest{}),
		swaggor.WithResponse(201, "Customer created", CustomerResponse{}),
		swaggor.WithResponse(400, "Validation error", ErrorResponse{}),
		swaggor.WithResponse(401, "Unauthorized", ErrorResponse{}),
		swaggor.WithSecurity("bearer"),
	)

	engine.AddRoute("/api/v4/customers/{uuid}", "PUT",
		"Replace Customer",
		"Full replacement of a customer record.",
		swaggor.WithTags("Customers"),
		swaggor.WithPathParam("uuid", "Customer UUID"),
		swaggor.WithRequestBody("Replacement data", true, CreateCustomerRequest{}),
		swaggor.WithResponse(200, "Customer updated", CustomerResponse{}),
		swaggor.WithResponse(404, "Customer not found", ErrorResponse{}),
		swaggor.WithSecurity("bearer"),
	)

	engine.AddRoute("/api/v4/customers/{uuid}", "PATCH",
		"Update Customer",
		"Partial update of a customer record.",
		swaggor.WithTags("Customers"),
		swaggor.WithPathParam("uuid", "Customer UUID"),
		swaggor.WithRequestBody("Partial data", false, CreateCustomerRequest{}),
		swaggor.WithResponse(200, "Customer updated", CustomerResponse{}),
		swaggor.WithResponse(404, "Customer not found", ErrorResponse{}),
		swaggor.WithSecurity("bearer"),
	)

	engine.AddRoute("/api/v4/customers/{uuid}", "DELETE",
		"Delete Customer",
		"Permanently removes a customer account.",
		swaggor.WithTags("Customers"),
		swaggor.WithPathParam("uuid", "Customer UUID"),
		swaggor.WithResponse(204, "Deleted successfully", nil),
		swaggor.WithResponse(404, "Customer not found", ErrorResponse{}),
		swaggor.WithSecurity("bearer"),
	)

	app := fiber.New(fiber.Config{AppName: "Fiber Swagger Integration"})

	app.Get("/api/v4/customers", func(c *fiber.Ctx) error {
		return c.JSON([]CustomerResponse{
			{
				UUID:     "7a9b2c1d-8f3e-4d5b-9a1c-3b2d1e4f5a6b",
				Username: "wendel.santos.fiber",
				TierID:   99,
				Status:   "active",
				Settings: AccountSettings{TwoFactorEnabled: true, PreferredLang: "en-US"},
			},
		})
	})

	app.Get("/api/v4/customers/:uuid", func(c *fiber.Ctx) error {
		return c.JSON(CustomerResponse{
			UUID:     c.Params("uuid"),
			Username: "wendel.santos.fiber",
			TierID:   99,
			Status:   "active",
			Settings: AccountSettings{TwoFactorEnabled: true, PreferredLang: "en-US"},
		})
	})

	app.Post("/api/v4/customers", func(c *fiber.Ctx) error {
		c.Status(http.StatusCreated)
		return c.JSON(CustomerResponse{
			UUID:     "new-uuid-here",
			Username: "new.user",
			TierID:   1,
			Status:   "active",
		})
	})

	app.Delete("/api/v4/customers/:uuid", func(c *fiber.Ctx) error {
		return c.SendStatus(http.StatusNoContent)
	})

	// fiber usa fasthttp internamente, então net/http handlers precisam do adaptor
	app.All("/swaggor/*", adaptor.HTTPHandler(engine.Handler()))

	log.Println("API:        http://localhost:3000/api/v4/customers")
	log.Println("Swagger UI: http://localhost:3000/swaggor/")

	if err := app.Listen(":3000"); err != nil {
		log.Fatalf("fiber: %v", err)
	}
}
