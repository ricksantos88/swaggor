package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ricksantos88/swaggergo" // Imports the local module
)

// AccountSettings maps administrative options for customer profiles.
type AccountSettings struct {
	TwoFactorEnabled bool   `json:"two_factor_enabled" description:"Enforces multi-factor authorization checkpoints"`
	PreferredLang    string `json:"preferred_lang" description:"ISO standard language code preferred by account customer"`
}

// CustomerResponse represents the production enterprise representation layout for core client consumers.
type CustomerResponse struct {
	UUID     string          `json:"uuid" description:"Immutable unique global transaction trace domain key identity"`
	Username string          `json:"username" description:"Alphanumeric user identification handle mapped globally"`
	TierID   int             `json:"tier_id" description:"Subscription classification ranking system layer identification code"`
	Settings AccountSettings `json:"settings" description:"Embedded user profile system configurations"`
}

func main() {
	// 1. Initialize Engine metadata definitions
	engine := swaggergo.NewEngine("Financial Customer Core API", "v3.1.2")

	// 2. Bind endpoints and expected runtime payloads cleanly
	engine.AddRoute(
		"/api/v3/customers",
		http.MethodGet,
		"Fetch Customer Information",
		"Retrieves sanitized domain customer entities associated with current authorization contexts.",
		CustomerResponse{},
	)

	// 3. Build default network routing engine multiplexer
	serverMux := http.NewServeMux()

	// 4. Implement enterprise system live business endpoint
	serverMux.HandleFunc("/api/v3/customers", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")

		response := CustomerResponse{
			UUID:     "9f8e7d6c-5b4a-3f2e-1d0c-9b8a7f6e5d4c",
			Username: "wendel.santos",
			TierID:   42,
			Settings: AccountSettings{
				TwoFactorEnabled: true,
				PreferredLang:    "pt-BR",
			},
		}
		_ = json.NewEncoder(w).Encode(response)
	})

	// 5. Mount documentation interface directly under the specified default mount path
	serverMux.Handle("/swagger-go/", engine.Handler())

	log.Println("[INFO] Operating application lifecycle metrics engine on: http://localhost:8080/api/v3/customers")
	log.Println("[INFO] Interactive OpenAPI graphical visualization route: http://localhost:8080/swagger-go/")

	if err := http.ListenAndServe(":8080", serverMux); err != nil {
		log.Fatalf("[CRITICAL] Infrastructure binding failure event intercepted: %v", err)
	}
}
