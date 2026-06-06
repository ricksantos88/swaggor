package main

import (
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/ricksantos88/swaggor"
)

// AccountSettings maps administrative options for customer profiles.
type AccountSettings struct {
	TwoFactorEnabled bool   `json:"two_factor_enabled" description:"Enforces multi-factor authorization checkpoints"`
	PreferredLang    string `json:"preferred_lang" description:"ISO standard language code preferred by account customer"`
}

// CustomerResponse represents the production enterprise representation layout.
type CustomerResponse struct {
	UUID     string          `json:"uuid" description:"Immutable unique global transaction trace domain key identity"`
	Username string          `json:"username" description:"Alphanumeric user identification handle mapped globally"`
	TierID   int             `json:"tier_id" description:"Subscription classification ranking system layer identification code"`
	Settings AccountSettings `json:"settings" description:"Embedded user profile system configurations"`
}

func main() {
	// 1. Initialize our Swaggor Engine
	engine := swaggor.NewEngine("Fiber High-Performance API", "v4.0.0")

	// 2. Register the endpoint metadata dynamically
	engine.AddRoute(
		"/api/v4/customers",
		http.MethodGet,
		"Fetch Customer Information",
		"Retrieves sanitized domain customer entities using the Fiber runtime.",
		CustomerResponse{},
	)

	// 3. Initialize the Fiber application instance
	app := fiber.New(fiber.Config{
		AppName: "Fiber Swagger Integration",
	})

	// 4. Implement the business logic route using native Fiber context (*fiber.Ctx)
	app.Get("/api/v4/customers", func(c *fiber.Ctx) error {
		response := CustomerResponse{
			UUID:     "7a9b2c1d-8f3e-4d5b-9a1c-3b2d1e4f5a6b",
			Username: "wendel.santos.fiber",
			TierID:   99,
			Settings: AccountSettings{
				TwoFactorEnabled: true,
				PreferredLang:    "en-US",
			},
		}

		// Fiber has built-in JSON marshaling
		return c.JSON(response)
	})

	// 5. The Magic: Adapt the net/http Handler into a Fiber Handler
	// We use app.All and a wildcard route "/swaggor/*" so the internal router
	// of our library can handle both "/" (HTML) and "/doc.json" (JSON).
	app.All("/swaggor/*", adaptor.HTTPHandler(engine.Handler()))

	log.Println("[INFO] Fiber Application active at: http://localhost:3000/api/v4/customers")
	log.Println("[INFO] Swagger UI Documentation at: http://localhost:3000/swaggor/")

	// 6. Start the fasthttp server
	if err := app.Listen(":3000"); err != nil {
		log.Fatalf("[CRITICAL] Fiber infrastructure failure: %v", err)
	}
}
