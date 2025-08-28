package authhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/giulio-alfieri/toq_server/internal/adapter/left/http/dto"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

// RequestPasswordChange handles password change request (public endpoint)
//
//	@Summary		Request password change
//	@Description	Initiate password change process by sending verification code
//	@Tags			Authentication
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.RequestPasswordChangeRequest	true	"Password change request data"
//	@Success		200		{object}	dto.RequestPasswordChangeResponse
//	@Failure		400		{object}	dto.ErrorResponse	"Invalid request format"
//	@Failure		404		{object}	dto.ErrorResponse	"User not found"
//	@Failure		429		{object}	dto.ErrorResponse	"Too many requests"
//	@Failure		500		{object}	dto.ErrorResponse	"Internal server error"
//	@Router			/auth/password/request [post]
func (ah *AuthHandler) RequestPasswordChange(c *gin.Context) {
	ctx, spanEnd, err := utils.GenerateTracer(c.Request.Context())
	if err != nil {
		utils.SendHTTPError(c, http.StatusInternalServerError, "TRACER_ERROR", "Failed to generate tracer")
		return
	}
	defer spanEnd()

	// Parse request
	var request dto.RequestPasswordChangeRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.SendHTTPError(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request format")
		return
	}

	// Call service
	err = ah.userService.RequestPasswordChange(ctx, request.NationalID)
	if err != nil {
		utils.SendHTTPError(c, http.StatusInternalServerError, "PASSWORD_CHANGE_REQUEST_FAILED", "Failed to request password change")
		return
	}

	// Success response
	c.JSON(http.StatusOK, dto.RequestPasswordChangeResponse{
		Message: "Password change code sent successfully",
	})
}
