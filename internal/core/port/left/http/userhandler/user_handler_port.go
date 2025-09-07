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
	GetUserRoles(c *gin.Context)
	GoHome(c *gin.Context)
	UpdateOptStatus(c *gin.Context)
	GetPhotoUploadURL(c *gin.Context)
	GetProfileThumbnails(c *gin.Context)
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

	// Agency handlers
	InviteRealtor(c *gin.Context)
	GetRealtorsByAgency(c *gin.Context)
	GetRealtorByID(c *gin.Context)
	DeleteRealtorByID(c *gin.Context)

	// Realtor handlers
	VerifyCreciDocuments(c *gin.Context)
	GetCreciUploadURL(c *gin.Context)
	AcceptInvitation(c *gin.Context)
	RejectInvitation(c *gin.Context)
	GetAgencyOfRealtor(c *gin.Context)
	DeleteAgencyOfRealtor(c *gin.Context)
}
