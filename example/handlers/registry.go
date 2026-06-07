package handlers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// NetHTTPRegistry maps function names to their net/http implementations.
// The adapters.LoadNetHTTP loader uses this to wire real functions from
// the names extracted by the parser — without reflection on function values.
var NetHTTPRegistry = map[string]http.HandlerFunc{
	"ListCustomers":   ListCustomers,
	"GetCustomer":     GetCustomer,
	"CreateCustomer":  CreateCustomer,
	"ReplaceCustomer": ReplaceCustomer,
	"DeleteCustomer":  DeleteCustomer,
}

// FiberRegistry maps function names to their Fiber implementations.
var FiberRegistry = map[string]fiber.Handler{
	"ListCustomersFiber":   ListCustomersFiber,
	"GetCustomerFiber":     GetCustomerFiber,
	"CreateCustomerFiber":  CreateCustomerFiber,
	"ReplaceCustomerFiber": ReplaceCustomerFiber,
	"DeleteCustomerFiber":  DeleteCustomerFiber,
}
