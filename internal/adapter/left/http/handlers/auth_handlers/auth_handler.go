package authhandlers

import (
	userservice "github.com/giulio-alfieri/toq_server/internal/core/service/user_service"
)

// AuthHandler implementa os handlers HTTP para operações de autenticação
type AuthHandler struct {
	userService userservice.UserServiceInterface
}

// NewAuthHandlerAdapter cria uma nova instância de AuthHandler
func NewAuthHandlerAdapter(
	userService userservice.UserServiceInterface,
) *AuthHandler {
	return &AuthHandler{
		userService: userService,
	}
}
