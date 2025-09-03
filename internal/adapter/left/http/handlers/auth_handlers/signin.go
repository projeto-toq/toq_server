package authhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/giulio-alfieri/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/giulio-alfieri/toq_server/internal/adapter/left/http/http_errors"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
	"github.com/google/uuid"
)

// SignIn handles user authentication (public endpoint)
//
//	@Summary		User sign in
//	@Description	Authenticate user with national ID and password. When sending a deviceToken, provide X-Device-Id (UUIDv4) for per-device associations.
//	@Tags			Authentication
//	@Accept			json
//	@Produce		json
//	@Param			X-Device-Id	header	string	false	"Device ID (UUIDv4). Required when request contains deviceToken."
//	@Param			request	body		dto.SignInRequest	true	"Sign in credentials"
//	@Success		200		{object}	dto.SignInResponse	"Successful authentication"
//	@Failure		400		{object}	dto.ErrorResponse	"Invalid request format or missing/invalid device ID when deviceToken provided"
//	@Failure		401		{object}	dto.ErrorResponse	"Invalid credentials"
//	@Failure		403		{object}	dto.ErrorResponse	"No active user roles"
//	@Failure		423		{object}	dto.ErrorResponse	"Account temporarily locked due to security measures"
//	@Failure		429		{object}	dto.ErrorResponse	"Too many attempts"
//	@Failure		500		{object}	dto.ErrorResponse	"Internal server error"
//	@Router			/auth/signin [post]
func (ah *AuthHandler) SignIn(c *gin.Context) {
	// Observação: tracing de request já é provido por TelemetryMiddleware; evitamos spans duplicados aqui.
	ctx := c.Request.Context()

	// Extract request context for security logging
	reqContext := utils.ExtractRequestContext(c)

	// Parse request
	var request dto.SignInRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httperrors.SendHTTPError(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request format")
		return
	}

	// Enforce device header contract when deviceToken is provided (development environment, no backward-compat)
	if request.DeviceToken != "" {
		headerDeviceID := c.GetHeader("X-Device-Id")
		if headerDeviceID == "" {
			httperrors.SendHTTPError(c, http.StatusBadRequest, "MISSING_DEVICE_ID", "X-Device-Id header is required when deviceToken is provided")
			return
		}
		if _, err := uuid.Parse(headerDeviceID); err != nil {
			httperrors.SendHTTPError(c, http.StatusBadRequest, "INVALID_DEVICE_ID", "X-Device-Id must be a valid UUID")
			return
		}
	}

	// Call service with enhanced context
	tokens, err := ah.userService.SignInWithContext(ctx, request.NationalID, request.Password, request.DeviceToken, reqContext.IPAddress, reqContext.UserAgent)
	if err != nil {
		if derr, ok := err.(utils.DomainError); ok {
			switch derr.Code() {
			case http.StatusUnauthorized:
				httperrors.SendHTTPError(c, http.StatusUnauthorized, "INVALID_CREDENTIALS", "Invalid credentials")
			case http.StatusLocked:
				httperrors.SendHTTPError(c, http.StatusLocked, "ACCOUNT_BLOCKED", "Account temporarily blocked due to security measures")
			case http.StatusForbidden:
				httperrors.SendHTTPError(c, http.StatusForbidden, "NO_ACTIVE_ROLES", "No active user roles")
			case http.StatusTooManyRequests:
				httperrors.SendHTTPError(c, http.StatusTooManyRequests, "TOO_MANY_ATTEMPTS", "Too many failed attempts")
			case http.StatusInternalServerError:
				httperrors.SendHTTPError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
			default:
				httperrors.SendHTTPError(c, http.StatusUnauthorized, "AUTHENTICATION_FAILED", "Authentication failed")
			}
		} else {
			httperrors.SendHTTPError(c, http.StatusUnauthorized, "INVALID_CREDENTIALS", "Invalid credentials")
		}
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
