package userhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/giulio-alfieri/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/giulio-alfieri/toq_server/internal/adapter/left/http/http_errors"
	coreutils "github.com/giulio-alfieri/toq_server/internal/core/utils"
	"github.com/giulio-alfieri/toq_server/internal/core/utils/validators"
)

// ConfirmPhoneChange handles confirming a phone number change with verification code
//
//	@Summary		Confirm phone number change
//	@Description	Confirm a phone number change using the verification code sent to the new phone
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.ConfirmPhoneChangeRequest	true	"Phone change confirmation data"
//	@Success		200		{object}	dto.ConfirmPhoneChangeResponse
//	@Failure		400		{object}	dto.ErrorResponse	"Invalid request format or code"
//	@Failure		401		{object}	dto.ErrorResponse	"Unauthorized"
//	@Failure		500		{object}	dto.ErrorResponse	"Internal server error"
//	@Router			/user/phone/confirm [post]
//	@Security		BearerAuth
func (uh *UserHandler) ConfirmPhoneChange(c *gin.Context) {
	// Enrich context with request info and user
	ctx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)

	// Parse request body using DTO
	var request dto.ConfirmPhoneChangeRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httperrors.SendHTTPError(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request format")
		return
	}

	// Validate verification code format
	if err := validators.ValidateCode(request.Code); err != nil {
		httperrors.SendHTTPError(c, http.StatusBadRequest, "INVALID_CODE", "Invalid code format")
		return
	}

	// Call service to confirm phone change (no tokens returned)
	err := uh.userService.ConfirmPhoneChange(ctx, request.Code)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	// Return success response
	c.JSON(http.StatusOK, dto.ConfirmPhoneChangeResponse{Message: "Phone changed successfully"})
}
