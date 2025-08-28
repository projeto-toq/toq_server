package authhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/giulio-alfieri/toq_server/internal/adapter/left/http/dto"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

// RefreshToken handles token refresh (public endpoint)
//
//	@Summary		Refresh access token
//	@Description	Generate new access token using refresh token
//	@Tags			Authentication
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.RefreshTokenRequest	true	"Refresh token data"
//	@Success		200		{object}	dto.RefreshTokenResponse
//	@Failure		400		{object}	dto.ErrorResponse	"Invalid request format"
//	@Failure		401		{object}	dto.ErrorResponse	"Invalid refresh token"
//	@Failure		500		{object}	dto.ErrorResponse	"Internal server error"
//	@Router			/auth/refresh [post]
func (ah *AuthHandler) RefreshToken(c *gin.Context) {
	ctx, spanEnd, err := utils.GenerateTracer(c.Request.Context())
	if err != nil {
		utils.SendHTTPError(c, http.StatusInternalServerError, "TRACER_ERROR", "Failed to generate tracer")
		return
	}
	defer spanEnd()

	// Parse request
	var request dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.SendHTTPError(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request format")
		return
	}

	// Call service
	tokens, err := ah.userService.RefreshTokens(ctx, request.RefreshToken)
	if err != nil {
		utils.SendHTTPError(c, http.StatusUnauthorized, "INVALID_REFRESH_TOKEN", "Invalid refresh token")
		return
	}

	// Success response
	c.JSON(http.StatusOK, dto.RefreshTokenResponse{
		Tokens: dto.TokensResponse{
			AccessToken:  tokens.AccessToken,
			RefreshToken: tokens.RefreshToken,
		},
	})
}
