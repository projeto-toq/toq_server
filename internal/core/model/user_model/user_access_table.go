package usermodel

// HTTP Access Control Tables
// Each map defines which HTTP endpoints each user role can access

// OwnerHTTPPrivileges defines endpoints accessible by Owner role
var OwnerHTTPPrivileges = map[string]bool{
	// Authentication (public endpoints - always allowed)
	"POST /api/v1/auth/signout": true,

	// Profile management
	"GET /api/v1/user/profile":            true,
	"PUT /api/v1/user/profile":            true,
	"DELETE /api/v1/user/account":         true,
	"GET /api/v1/user/onboarding":         true,
	"GET /api/v1/user/roles":              true,
	"GET /api/v1/user/home":               true,
	"PUT /api/v1/user/opt-status":         true,
	"POST /api/v1/user/photo/upload-url":  true,
	"GET /api/v1/user/profile/thumbnails": true,

	// Email/Phone change
	"POST /api/v1/user/email/request": true,
	"POST /api/v1/user/email/confirm": true,
	"POST /api/v1/user/email/resend":  true,
	"POST /api/v1/user/phone/request": true,
	"POST /api/v1/user/phone/confirm": true,
	"POST /api/v1/user/phone/resend":  true,

	// Role management (Owner/Realtor only)
	"POST /api/v1/user/role/alternative": true,
	"POST /api/v1/user/role/switch":      true,

	// Listing management (Owner side)
	"GET /api/v1/listings":                 true,
	"POST /api/v1/listings":                true,
	"GET /api/v1/listings/search":          true,
	"GET /api/v1/listings/options":         true,
	"GET /api/v1/listings/features/base":   true,
	"GET /api/v1/listings/:id":             true,
	"PUT /api/v1/listings/:id":             true,
	"DELETE /api/v1/listings/:id":          true,
	"POST /api/v1/listings/:id/end-update": true,
	"GET /api/v1/listings/:id/status":      true,
	"POST /api/v1/listings/:id/approve":    true,
	"POST /api/v1/listings/:id/reject":     true,
	"POST /api/v1/listings/:id/suspend":    true,
	"POST /api/v1/listings/:id/release":    true,
	"POST /api/v1/listings/:id/copy":       true,
	"GET /api/v1/listings/:id/visits":      true,
	"GET /api/v1/listings/:id/offers":      true,

	// Visit management (Owner side)
	"GET /api/v1/visits":              true,
	"POST /api/v1/visits/:id/approve": true,
	"POST /api/v1/visits/:id/reject":  true,

	// Offer management (Owner side)
	"GET /api/v1/offers":              true,
	"POST /api/v1/offers/:id/approve": true,
	"POST /api/v1/offers/:id/reject":  true,

	// Evaluations
	"POST /api/v1/realtors/:id/evaluate": true,
}

// RealtorHTTPPrivileges defines endpoints accessible by Realtor role
var RealtorHTTPPrivileges = map[string]bool{
	// Authentication
	"POST /api/v1/auth/signout": true,

	// Profile management
	"GET /api/v1/user/profile":            true,
	"PUT /api/v1/user/profile":            true,
	"DELETE /api/v1/user/account":         true,
	"GET /api/v1/user/onboarding":         true,
	"GET /api/v1/user/roles":              true,
	"GET /api/v1/user/home":               true,
	"PUT /api/v1/user/opt-status":         true,
	"POST /api/v1/user/photo/upload-url":  true,
	"GET /api/v1/user/profile/thumbnails": true,

	// Email/Phone change
	"POST /api/v1/user/email/request": true,
	"POST /api/v1/user/email/confirm": true,
	"POST /api/v1/user/email/resend":  true,
	"POST /api/v1/user/phone/request": true,
	"POST /api/v1/user/phone/confirm": true,
	"POST /api/v1/user/phone/resend":  true,

	// Role management (Owner/Realtor only)
	"POST /api/v1/user/role/alternative": true,
	"POST /api/v1/user/role/switch":      true,

	// Realtor specific operations
	"POST /api/v1/realtor/creci/verify":      true,
	"POST /api/v1/realtor/creci/upload-url":  true,
	"POST /api/v1/realtor/invitation/accept": true,
	"POST /api/v1/realtor/invitation/reject": true,
	"GET /api/v1/realtor/agency":             true,
	"DELETE /api/v1/realtor/agency":          true,

	// Listing operations (Realtor side)
	"GET /api/v1/listings":                    true,
	"GET /api/v1/listings/search":             true,
	"GET /api/v1/listings/options":            true,
	"GET /api/v1/listings/features/base":      true,
	"GET /api/v1/listings/:id":                true,
	"GET /api/v1/listings/favorites":          true,
	"POST /api/v1/listings/:id/share":         true,
	"POST /api/v1/listings/:id/favorite":      true,
	"DELETE /api/v1/listings/:id/favorite":    true,
	"POST /api/v1/listings/:id/visit/request": true,
	"GET /api/v1/listings/:id/visits":         true,
	"POST /api/v1/listings/:id/offers":        true,
	"GET /api/v1/listings/:id/offers":         true,

	// Visit management (Realtor side)
	"GET /api/v1/visits":              true,
	"DELETE /api/v1/visits/:id":       true,
	"POST /api/v1/visits/:id/confirm": true,

	// Offer management (Realtor side)
	"GET /api/v1/offers":           true,
	"PUT /api/v1/offers/:id":       true,
	"DELETE /api/v1/offers/:id":    true,
	"POST /api/v1/offers/:id/send": true,

	// Evaluations
	"POST /api/v1/owners/:id/evaluate": true,
}

// AgencyHTTPPrivileges defines endpoints accessible by Agency role
var AgencyHTTPPrivileges = map[string]bool{
	// Authentication
	"POST /api/v1/auth/signout": true,

	// Profile management
	"GET /api/v1/user/profile":            true,
	"PUT /api/v1/user/profile":            true,
	"DELETE /api/v1/user/account":         true,
	"GET /api/v1/user/onboarding":         true,
	"GET /api/v1/user/roles":              true,
	"GET /api/v1/user/home":               true,
	"PUT /api/v1/user/opt-status":         true,
	"POST /api/v1/user/photo/upload-url":  true,
	"GET /api/v1/user/profile/thumbnails": true,

	// Email/Phone change
	"POST /api/v1/user/email/request": true,
	"POST /api/v1/user/email/confirm": true,
	"POST /api/v1/user/email/resend":  true,
	"POST /api/v1/user/phone/request": true,
	"POST /api/v1/user/phone/confirm": true,
	"POST /api/v1/user/phone/resend":  true,

	// Agency specific operations
	"POST /api/v1/agency/documents/upload-url": true,
	"POST /api/v1/agency/invite-realtor":       true,
	"GET /api/v1/agency/realtors":              true,
	"GET /api/v1/agency/realtors/:id":          true,
	"DELETE /api/v1/agency/realtors/:id":       true,

	// Listing operations (basic read access)
	"GET /api/v1/listings":            true,
	"GET /api/v1/listings/search":     true,
	"GET /api/v1/listings/:id":        true,
	"GET /api/v1/listings/:id/visits": true,
	"GET /api/v1/listings/:id/offers": true,

	// Limited visit and offer access
	"GET /api/v1/visits": true,
	"GET /api/v1/offers": true,
}

// RootHTTPPrivileges defines endpoints accessible by Root role (full access)
var RootHTTPPrivileges = map[string]bool{
	// Root (Admin) has access to everything - this will be handled differently
	// in the IsHTTPEndpointAllowed function by returning true for all endpoints
}

// GetHTTPPrivilegesForRole returns the privileges map for a given role
func GetHTTPPrivilegesForRole(role UserRole) map[string]bool {
	switch role {
	case RoleRoot:
		return RootHTTPPrivileges // Root has full access
	case RoleOwner:
		return OwnerHTTPPrivileges
	case RoleRealtor:
		return RealtorHTTPPrivileges
	case RoleAgency:
		return AgencyHTTPPrivileges
	default:
		return make(map[string]bool) // No access for unknown roles
	}
}

// IsHTTPEndpointAllowed checks if a user role has access to a specific HTTP endpoint
func IsHTTPEndpointAllowed(role UserRole, method, path string) bool {
	// Root role has access to everything
	if role == RoleRoot {
		return true
	}

	endpoint := method + " " + path
	privileges := GetHTTPPrivilegesForRole(role)

	// Check exact match first
	if allowed, exists := privileges[endpoint]; exists {
		return allowed
	}

	// Check for parameterized routes (e.g., /api/v1/listings/:id)
	for pattern := range privileges {
		if matchesParameterizedRoute(endpoint, pattern) {
			return privileges[pattern]
		}
	}

	return false
}

// matchesParameterizedRoute checks if an endpoint matches a parameterized route pattern
func matchesParameterizedRoute(endpoint, pattern string) bool {
	endpointParts := splitMethodPath(endpoint)
	patternParts := splitMethodPath(pattern)

	if len(endpointParts) != 2 || len(patternParts) != 2 {
		return false
	}

	// Method must match exactly
	if endpointParts[0] != patternParts[0] {
		return false
	}

	// Check if paths match with parameters
	return matchesPathWithParams(endpointParts[1], patternParts[1])
}

// splitMethodPath splits "METHOD /path" into ["METHOD", "/path"]
func splitMethodPath(methodPath string) []string {
	parts := make([]string, 2)
	spaceIndex := -1
	for i, char := range methodPath {
		if char == ' ' {
			spaceIndex = i
			break
		}
	}
	if spaceIndex == -1 {
		return []string{methodPath} // Invalid format
	}
	parts[0] = methodPath[:spaceIndex]
	parts[1] = methodPath[spaceIndex+1:]
	return parts
}

// matchesPathWithParams checks if a path matches a pattern with :param placeholders
func matchesPathWithParams(path, pattern string) bool {
	pathParts := splitPath(path)
	patternParts := splitPath(pattern)

	if len(pathParts) != len(patternParts) {
		return false
	}

	for i := 0; i < len(pathParts); i++ {
		// Skip parameter placeholders (start with :)
		if len(patternParts[i]) > 0 && patternParts[i][0] == ':' {
			continue
		}
		// Exact match required for non-parameter parts
		if pathParts[i] != patternParts[i] {
			return false
		}
	}

	return true
}

// splitPath splits "/api/v1/users/:id" into ["api", "v1", "users", ":id"]
func splitPath(path string) []string {
	if len(path) == 0 {
		return []string{}
	}

	// Remove leading slash
	if path[0] == '/' {
		path = path[1:]
	}

	if len(path) == 0 {
		return []string{}
	}

	parts := []string{}
	current := ""

	for _, char := range path {
		if char == '/' {
			if len(current) > 0 {
				parts = append(parts, current)
				current = ""
			}
		} else {
			current += string(char)
		}
	}

	if len(current) > 0 {
		parts = append(parts, current)
	}

	return parts
}
