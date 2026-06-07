package handlers

// AccountSettings holds user-level preferences.
type AccountSettings struct {
	TwoFactorEnabled bool   `json:"two_factor_enabled" description:"Enforces multi-factor authorization checkpoints"`
	PreferredLang    string `json:"preferred_lang"     description:"ISO language code preferred by the account holder" example:"pt-BR"`
}

// CustomerResponse is the outbound representation of a customer.
type CustomerResponse struct {
	UUID     string          `json:"uuid"     description:"Immutable globally unique identifier" format:"uuid"`
	Username string          `json:"username" description:"Alphanumeric user handle"            example:"wendel.santos"`
	TierID   int32           `json:"tier_id"  description:"Subscription tier"                   example:"42"`
	Status   string          `json:"status"   description:"Account status" enum:"active,inactive,suspended" example:"active"`
	Settings AccountSettings `json:"settings" description:"Embedded profile configuration"`
}

// CreateCustomerRequest is the inbound payload for customer creation / replacement.
type CreateCustomerRequest struct {
	Username string `json:"username" description:"Desired username"      example:"wendel.santos"`
	Email    string `json:"email"    description:"Account email address" format:"email"`
	TierID   int32  `json:"tier_id"  description:"Initial subscription tier"`
}

// ErrorResponse is the standard error envelope.
type ErrorResponse struct {
	Code    int    `json:"code"    description:"Machine-readable error code"`
	Message string `json:"message" description:"Human-readable error description"`
}
