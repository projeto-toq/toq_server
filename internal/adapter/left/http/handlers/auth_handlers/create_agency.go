package authhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	httputils "github.com/projeto-toq/toq_server/internal/adapter/left/http/utils"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

// CreateAgency handles agency account creation (public endpoint)
//
//	@Summary		Create agency account
//	@Description	Create a new agency account with user information
//	@Tags			Authentication
//	@Accept			json
//	@Produce		json
//	@Param		X-Device-Id	header	string	false	"Device ID (UUIDv4). Optional but recommended when providing deviceToken for automatic sign-in"
//	@Param			request	body		dto.CreateAgencyRequest	true	"Agency creation data (include optional deviceToken for push notifications)"
//	@Success		201		{object}	dto.CreateAgencyResponse
//	@Failure		400		{object}	dto.ErrorResponse	"Invalid request format"
//	@Failure		422		{object}	dto.ErrorResponse	"Validation failed"
//	@Failure		409		{object}	dto.ErrorResponse	"User already exists"
//	@Failure		500		{object}	dto.ErrorResponse	"Internal server error"
//	@Router			/auth/agency [post]
func (ah *AuthHandler) CreateAgency(c *gin.Context) {
	// Observação: tracing de request já é provido por TelemetryMiddleware; evitamos spans duplicados aqui.
	ctx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)

	// Parse request
	var request dto.CreateAgencyRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httperrors.SendHTTPErrorObj(c, httputils.MapBindingError(err))
		return
	}

	// Validate and parse date fields with precise attribution
	bornAt, creciValidity, derr := httputils.ValidateUserDates(request.Agency, "agency")
	if derr != nil {
		httperrors.SendHTTPErrorObj(c, derr)
		return
	}

	// Create user model from DTO (using parsed dates)
	user, err := ah.createUserFromDTO(request.Agency, bornAt, creciValidity)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, coreutils.NewHTTPError(http.StatusUnprocessableEntity, "Validation failed"))
		return
	}

	// Extract request context for security logging and session metadata
	reqContext := coreutils.ExtractRequestContext(c)

	// Call service: cria a conta e autentica via SignIn padrão
	tokens, err := ah.userService.CreateAgency(ctx, user, request.Agency.Password, request.DeviceToken, reqContext.IPAddress, reqContext.UserAgent)
	if err != nil {
		if derr, ok := err.(coreutils.DomainError); ok {
			httperrors.SendHTTPErrorObj(c, derr)
			return
		}
		httperrors.SendHTTPErrorObj(c, coreutils.NewHTTPError(http.StatusConflict, "Failed to create agency"))
		return
	}

	// Success response
	c.JSON(http.StatusCreated, dto.CreateAgencyResponse{
		Tokens: dto.TokensResponse{
			AccessToken:  tokens.AccessToken,
			RefreshToken: tokens.RefreshToken,
		},
	})
}
