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
