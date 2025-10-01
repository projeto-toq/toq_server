package dto

// CSPPolicyResponse defines the payload returned by the CSP policy endpoints.
// swagger:model CSPPolicyResponse
type CSPPolicyResponse struct {
	// Version is the optimistic-lock value of the stored policy.
	Version int64 `json:"version"`
	// Directives contains the CSP directives in key/value format.
	Directives map[string]string `json:"directives"`
}

// UpdateCSPPolicyRequest represents the incoming payload used to update the CSP policy.
// swagger:model UpdateCSPPolicyRequest
type UpdateCSPPolicyRequest struct {
	// Version is the optimistic-lock control. Use 0 to create a new policy.
	Version int64 `json:"version"`
	// Directives must include at least the default-src directive.
	Directives map[string]string `json:"directives" binding:"required"`
}
