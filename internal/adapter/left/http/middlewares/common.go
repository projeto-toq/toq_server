package middlewares

// Common utility functions for middlewares

// isAuthRequiredEndpoint checks if an endpoint requires authentication
func isAuthRequiredEndpoint(path string) bool {
	publicEndpoints := []string{
		"/api/v1/auth/owner",
		"/api/v1/auth/realtor",
		"/api/v1/auth/agency",
		"/api/v1/auth/signin",
		"/api/v1/auth/refresh",
		"/api/v1/auth/password/request",
		"/api/v1/auth/password/confirm",
		"/healthz",
		"/readyz",
	}

	for _, endpoint := range publicEndpoints {
		if path == endpoint {
			return false
		}
	}
	return true
}

// isAccessControlRequired checks if an endpoint requires access control
func isAccessControlRequired(method, path string) bool {
	publicEndpoints := map[string]bool{
		"POST /api/v1/auth/signin":              false,
		"POST /api/v1/auth/signup":              false,
		"POST /api/v1/auth/forgot-password":     false,
		"POST /api/v1/auth/reset-password":      false,
		"POST /api/v1/auth/verify-email":        false,
		"POST /api/v1/auth/resend-verification": false,
		"GET /api/v1/health":                    false,
		"GET /api/v1/status":                    false,
	}

	endpoint := method + " " + path
	required, exists := publicEndpoints[endpoint]
	if exists {
		return required // false for public endpoints
	}
	return true // true for all other endpoints
}
