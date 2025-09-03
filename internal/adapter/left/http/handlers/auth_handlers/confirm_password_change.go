package authhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/giulio-alfieri/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/giulio-alfieri/toq_server/internal/adapter/left/http/http_errors"
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
//	@Failure		409		{object}	dto.ErrorResponse	"Password change not pending"
//	@Failure		422		{object}	dto.ErrorResponse	"Invalid verification code or password validation failed"
//	@Failure		410		{object}	dto.ErrorResponse	"Verification code expired"
//	@Failure		500		{object}	dto.ErrorResponse	"Internal server error"
//	@Router			/auth/password/confirm [post]
func (ah *AuthHandler) ConfirmPasswordChange(c *gin.Context) {
	ctx, spanEnd, err := utils.GenerateTracer(c.Request.Context())
	if err != nil {
		httperrors.SendHTTPError(c, http.StatusInternalServerError, "TRACER_ERROR", "Failed to generate tracer")
		return
	}
	defer spanEnd()

	// Parse request
	var request dto.ConfirmPasswordChangeRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httperrors.SendHTTPError(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request format")
		return
	}

	// Call service with structured error propagation
	if err = ah.userService.ConfirmPasswordChange(ctx, request.NationalID, request.NewPassword, request.Code); err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	// Success response
	c.JSON(http.StatusOK, dto.ConfirmPasswordChangeResponse{
		Message: "Password changed successfully",
	})
}
