package authhandlers

import (
	globalservice "github.com/projeto-toq/toq_server/internal/core/service/global_service"
	userservice "github.com/projeto-toq/toq_server/internal/core/service/user_service"
	"github.com/projeto-toq/toq_server/internal/core/utils/hmacauth"
)

// AuthHandler implementa os handlers HTTP para operações de autenticação
type AuthHandler struct {
	userService   userservice.UserServiceInterface
	globalService globalservice.GlobalServiceInterface
	hmacValidator *hmacauth.Validator
}

// NewAuthHandlerAdapter cria uma nova instância de AuthHandler
func NewAuthHandlerAdapter(
	userService userservice.UserServiceInterface,
	globalService globalservice.GlobalServiceInterface,
	hmacValidator *hmacauth.Validator,
) *AuthHandler {
	return &AuthHandler{
		userService:   userService,
		globalService: globalService,
		hmacValidator: hmacValidator,
	}
}
