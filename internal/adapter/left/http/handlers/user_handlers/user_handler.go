package userhandlers

import (
	userhandlerport "github.com/giulio-alfieri/toq_server/internal/core/port/left/http/userhandler"
	complexservice "github.com/giulio-alfieri/toq_server/internal/core/service/complex_service"
	globalservice "github.com/giulio-alfieri/toq_server/internal/core/service/global_service"
	userservice "github.com/giulio-alfieri/toq_server/internal/core/service/user_service"
)

// UserHandler implementa os handlers HTTP para operações de usuário
type UserHandler struct {
	userService    userservice.UserServiceInterface
	globalService  globalservice.GlobalServiceInterface
	complexService complexservice.ComplexServiceInterface
}

// NewUserHandlerAdapter cria uma nova instância de UserHandler
func NewUserHandlerAdapter(
	userService userservice.UserServiceInterface,
	globalService globalservice.GlobalServiceInterface,
	complexService complexservice.ComplexServiceInterface,
) userhandlerport.UserHandlerPort {
	return &UserHandler{
		userService:    userService,
		globalService:  globalService,
		complexService: complexService,
	}
}
