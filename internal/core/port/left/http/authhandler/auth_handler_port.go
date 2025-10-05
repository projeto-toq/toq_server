package authhandler

import "github.com/gin-gonic/gin"

// AuthHandlerPort define a interface para handlers de autenticação
type AuthHandlerPort interface {
	// Public authentication endpoints (no auth required)
	CreateOwner(c *gin.Context)
	CreateRealtor(c *gin.Context)
	CreateAgency(c *gin.Context)
	SignIn(c *gin.Context)
	RefreshToken(c *gin.Context)
	RequestPasswordChange(c *gin.Context)
	ConfirmPasswordChange(c *gin.Context)
	ValidateCPF(c *gin.Context)
	ValidateCNPJ(c *gin.Context)
	ValidateCEP(c *gin.Context)
}
