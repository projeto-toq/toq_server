package middlewares

import coreutils "github.com/projeto-toq/toq_server/internal/core/utils"

// Common utility functions for middlewares

// isAuthRequiredEndpoint checks if an endpoint requires authentication
func isAuthRequiredEndpoint(path string) bool {
	// Delegate to core utils for a single source of truth
	return !coreutils.IsPublicEndpoint(path)
}
