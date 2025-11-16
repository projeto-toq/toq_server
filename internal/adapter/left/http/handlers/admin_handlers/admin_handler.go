package adminhandlers

import (
	"github.com/gin-gonic/gin"
	complexservice "github.com/projeto-toq/toq_server/internal/core/service/complex_service"
	listingservices "github.com/projeto-toq/toq_server/internal/core/service/listing_service"
	permissionservice "github.com/projeto-toq/toq_server/internal/core/service/permission_service"
	userservices "github.com/projeto-toq/toq_server/internal/core/service/user_service"
)

// AdminHandler handles administrative operations for the TOQ Server API.
// It provides endpoints for managing users, permissions, roles, complexes, and system metadata.
// All handlers require authentication and admin-level permissions via AuthMiddleware and PermissionMiddleware.
type AdminHandler struct {
	userService       userservices.UserServiceInterface
	listingService    listingservices.ListingServiceInterface
	permissionService permissionservice.PermissionServiceInterface
	complexService    complexservice.ComplexServiceInterface
	router            *gin.Engine // Gin engine reference for route introspection
}

// NewAdminHandlerAdapter creates a new AdminHandler with injected service dependencies
// and a reference to the Gin router for route introspection.
//
// The router reference is stored to allow the GetAdminRoutes handler to introspect
// registered routes at request time. Routes are registered after handler creation,
// so we cannot cache them in the constructor.
//
// Parameters:
//   - userService: Service for user management operations
//   - listingService: Service for listing management operations
//   - permissionService: Service for permission and role management
//   - complexService: Service for complex (building) management
//   - router: Gin engine instance for route introspection
//
// Returns:
//   - *AdminHandler: Configured admin handler ready for route registration
func NewAdminHandlerAdapter(
	userService userservices.UserServiceInterface,
	listingService listingservices.ListingServiceInterface,
	permissionService permissionservice.PermissionServiceInterface,
	complexService complexservice.ComplexServiceInterface,
	router *gin.Engine,
) *AdminHandler {
	return &AdminHandler{
		userService:       userService,
		listingService:    listingService,
		permissionService: permissionService,
		complexService:    complexService,
		router:            router, // Store router reference for route introspection
	}
}
