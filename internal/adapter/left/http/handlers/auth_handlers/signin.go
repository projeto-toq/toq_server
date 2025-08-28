package authhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/giulio-alfieri/toq_server/internal/adapter/left/http/dto"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

// SignIn handles user authentication (public endpoint)
//
//	@Summary		User sign in
//	@Description	Authenticate user with national ID and password
//	@Tags			Authentication
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.SignInRequest	true	"Sign in credentials"
//	@Success		200		{object}	dto.SignInResponse
//	@Failure		400		{object}	dto.ErrorResponse	"Invalid request format"
//	@Failure		401		{object}	dto.ErrorResponse	"Invalid credentials"
//	@Failure		429		{object}	dto.ErrorResponse	"Too many attempts"
//	@Failure		500		{object}	dto.ErrorResponse	"Internal server error"
//	@Router			/auth/signin [post]
func (ah *AuthHandler) SignIn(c *gin.Context) {
	ctx, spanEnd, err := utils.GenerateTracer(c.Request.Context())
	if err != nil {
		utils.SendHTTPError(c, http.StatusInternalServerError, "TRACER_ERROR", "Failed to generate tracer")
		return
	}
	defer spanEnd()

	// Parse request
	var request dto.SignInRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.SendHTTPError(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request format")
		return
	}

	// Call service
	tokens, err := ah.userService.SignIn(ctx, request.NationalID, request.Password, request.DeviceToken)
	if err != nil {
		utils.SendHTTPError(c, http.StatusUnauthorized, "INVALID_CREDENTIALS", "Invalid credentials")
		return
	}

	// Success response
	c.JSON(http.StatusOK, dto.SignInResponse{
		Tokens: dto.TokensResponse{
			AccessToken:  tokens.AccessToken,
			RefreshToken: tokens.RefreshToken,
		},
	})
}
