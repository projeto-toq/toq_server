package userhandlers

import (
	userhandlerport "github.com/projeto-toq/toq_server/internal/core/port/left/http/userhandler"
	complexservice "github.com/projeto-toq/toq_server/internal/core/service/complex_service"
	globalservice "github.com/projeto-toq/toq_server/internal/core/service/global_service"
	permissionservice "github.com/projeto-toq/toq_server/internal/core/service/permission_service"
	userservice "github.com/projeto-toq/toq_server/internal/core/service/user_service"
)

// UserHandler implementa os handlers HTTP para operações de usuário
type UserHandler struct {
	userService       userservice.UserServiceInterface
	globalService     globalservice.GlobalServiceInterface
	complexService    complexservice.ComplexServiceInterface
	permissionService permissionservice.PermissionServiceInterface
}

// NewUserHandlerAdapter cria uma nova instância de UserHandler
func NewUserHandlerAdapter(
	userService userservice.UserServiceInterface,
	globalService globalservice.GlobalServiceInterface,
	complexService complexservice.ComplexServiceInterface,
	permissionService permissionservice.PermissionServiceInterface,
) userhandlerport.UserHandlerPort {
	return &UserHandler{
		userService:       userService,
		globalService:     globalService,
		complexService:    complexService,
		permissionService: permissionService,
	}
}
