package userhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/giulio-alfieri/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/giulio-alfieri/toq_server/internal/adapter/left/http/http_errors"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
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
	// Generate tracer for observability
	ctx, spanEnd, err := utils.GenerateTracer(c.Request.Context())
	if err != nil {
		httperrors.SendHTTPError(c, http.StatusInternalServerError, "TRACER_ERROR", "Failed to generate tracer")
		return
	}
	defer spanEnd()

	// Get user information from context using Context Utils (set by auth middleware)
	userInfo, err := utils.GetUserInfoFromGinContext(c)
	if err != nil {
		httperrors.SendHTTPError(c, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

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

	// Enrich context with request information for service layer
	enrichedCtx := utils.EnrichContextWithRequestInfo(ctx, c)

	// Call service to confirm phone change
	tokens, err := uh.userService.ConfirmPhoneChange(enrichedCtx, userInfo.ID, request.Code)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	// Prepare successful response with new tokens
	response := dto.ConfirmPhoneChangeResponse{
		Tokens: dto.TokensResponse{
			AccessToken:  tokens.AccessToken,
			RefreshToken: tokens.RefreshToken,
		},
	}

	// Return success response
	c.JSON(http.StatusOK, response)
}
