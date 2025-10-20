package adminhandlers

import (
	complexservice "github.com/projeto-toq/toq_server/internal/core/service/complex_service"
	listingservices "github.com/projeto-toq/toq_server/internal/core/service/listing_service"
	permissionservice "github.com/projeto-toq/toq_server/internal/core/service/permission_service"
	userservices "github.com/projeto-toq/toq_server/internal/core/service/user_service"
)

type AdminHandler struct {
	userService       userservices.UserServiceInterface
	listingService    listingservices.ListingServiceInterface
	permissionService permissionservice.PermissionServiceInterface
	complexService    complexservice.ComplexServiceInterface
}

func NewAdminHandlerAdapter(
	userService userservices.UserServiceInterface,
	listingService listingservices.ListingServiceInterface,
	permissionService permissionservice.PermissionServiceInterface,
	complexService complexservice.ComplexServiceInterface,
) *AdminHandler {
	return &AdminHandler{
		userService:       userService,
		listingService:    listingService,
		permissionService: permissionService,
		complexService:    complexService,
	}
}
