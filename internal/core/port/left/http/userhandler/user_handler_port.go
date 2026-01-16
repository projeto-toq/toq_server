package userhandler

import "github.com/gin-gonic/gin"

// UserHandlerPort define a interface para handlers de usuário
type UserHandlerPort interface {
	// Authentication handlers
	SignOut(c *gin.Context)

	// Profile handlers
	GetProfile(c *gin.Context)
	UpdateProfile(c *gin.Context)
	DeleteAccount(c *gin.Context)
	UpdateOptStatus(c *gin.Context)
	PostPhotoUploadURL(c *gin.Context)
	PostPhotoDownloadURL(c *gin.Context)
	GetUserStatus(c *gin.Context)

	// Email/Phone change handlers
	RequestEmailChange(c *gin.Context)
	ConfirmEmailChange(c *gin.Context)
	ResendEmailChangeCode(c *gin.Context)
	RequestPhoneChange(c *gin.Context)
	ConfirmPhoneChange(c *gin.Context)
	ResendPhoneChangeCode(c *gin.Context)

	// Role management handlers
	AddAlternativeUserRole(c *gin.Context)
	SwitchUserRole(c *gin.Context)

	// Realtor handlers
	VerifyCreciDocuments(c *gin.Context)
	PostCreciUploadURL(c *gin.Context)
}
