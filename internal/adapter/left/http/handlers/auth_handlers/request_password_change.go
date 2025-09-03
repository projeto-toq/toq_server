package authhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/giulio-alfieri/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/giulio-alfieri/toq_server/internal/adapter/left/http/http_errors"
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
//	@Failure		429		{object}	dto.ErrorResponse	"Too many requests"
//	@Failure		500		{object}	dto.ErrorResponse	"Internal server error"
//	@Router			/auth/password/request [post]
func (ah *AuthHandler) RequestPasswordChange(c *gin.Context) {
	ctx, spanEnd, err := utils.GenerateTracer(c.Request.Context())
	if err != nil {
		httperrors.SendHTTPError(c, http.StatusInternalServerError, "TRACER_ERROR", "Failed to generate tracer")
		return
	}
	defer spanEnd()

	// Parse request
	var request dto.RequestPasswordChangeRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httperrors.SendHTTPError(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request format")
		return
	}

	// Call service (privacy-preserving: always return 200 on not found)
	if err = ah.userService.RequestPasswordChange(ctx, request.NationalID); err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	// Success response (generic message to avoid user enumeration)
	c.JSON(http.StatusOK, dto.RequestPasswordChangeResponse{Message: "If the account exists, a code has been sent"})
}
