package authhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/giulio-alfieri/toq_server/internal/adapter/left/http/dto"
"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

// ConfirmPasswordChange handles password change confirmation (public endpoint)
//
//	@Summary		Confirm password change
//	@Description	Confirm password change using verification code
//	@Tags			Authentication
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.ConfirmPasswordChangeRequest	true	"Password change confirmation data"
//	@Success		200		{object}	dto.ConfirmPasswordChangeResponse
//	@Failure		400		{object}	dto.ErrorResponse	"Invalid request format"
//	@Failure		401		{object}	dto.ErrorResponse	"Invalid verification code"
//	@Failure		404		{object}	dto.ErrorResponse	"User not found"
//	@Failure		422		{object}	dto.ErrorResponse	"Password validation failed"
//	@Failure		500		{object}	dto.ErrorResponse	"Internal server error"
//	@Router			/auth/password/confirm [post]
func (ah *AuthHandler) ConfirmPasswordChange(c *gin.Context) {
	ctx, spanEnd, err := utils.GenerateTracer(c.Request.Context())
	if err != nil {
		utils.SendHTTPError(c, http.StatusInternalServerError, "TRACER_ERROR", "Failed to generate tracer")
		return
	}
	defer spanEnd()

	// Parse request
	var request dto.ConfirmPasswordChangeRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.SendHTTPError(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request format")
		return
	}

	// Call service
	err = ah.userService.ConfirmPasswordChange(ctx, request.NationalID, request.Code, request.NewPassword)
	if err != nil {
		utils.SendHTTPError(c, http.StatusInternalServerError, "PASSWORD_CHANGE_CONFIRMATION_FAILED", "Failed to confirm password change")
		return
	}

	// Success response
	c.JSON(http.StatusOK, dto.ConfirmPasswordChangeResponse{
		Message: "Password changed successfully",
	})
}
