package authhandlers

import (
	"net/http"

	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/giulio-alfieri/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/giulio-alfieri/toq_server/internal/adapter/left/http/http_errors"
	httputils "github.com/giulio-alfieri/toq_server/internal/adapter/left/http/utils"
	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

// CreateRealtor handles realtor account creation (public endpoint)
//
//	@Summary		Create realtor account
//	@Description	Create a new realtor account with user and CRECI information
//	@Tags			Authentication
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.CreateRealtorRequest	true	"Realtor creation data (include optional deviceToken for push notifications)"
//	@Success		201		{object}	dto.CreateRealtorResponse
//	@Failure		400		{object}	dto.ErrorResponse	"Validation error (invalid input data)"
//	@Failure		422		{object}	dto.ErrorResponse	"Validation failed (see details)"
//	@Failure		409		{object}	dto.ErrorResponse	"User already exists"
//	@Failure		500		{object}	dto.ErrorResponse	"Internal server error"
//	@Router			/auth/realtor [post]
func (ah *AuthHandler) CreateRealtor(c *gin.Context) {
	// Observação: tracing de request já é provido por TelemetryMiddleware; evitamos spans duplicados aqui.
	ctx := c.Request.Context()

	// Parse request
	var request dto.CreateRealtorRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httperrors.SendHTTPErrorObj(c, httputils.MapBindingError(err))
		return
	}

	// Validate and parse date fields with precise attribution
	bornAt, creciValidity, derr := httputils.ValidateUserDates(request.Realtor, "realtor")
	if derr != nil {
		httperrors.SendHTTPErrorObj(c, derr)
		return
	}

	// Create user model from DTO (using parsed dates)
	user, err := ah.createUserFromDTO(request.Realtor, bornAt, creciValidity)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, utils.NewHTTPError(http.StatusUnprocessableEntity, "Validation failed"))
		return
	}

	// Extract request context for security logging and session metadata
	reqContext := utils.ExtractRequestContext(c)

	// Debug: rastrear valores de contexto antes de chamar o serviço
	headerDeviceID := c.GetHeader("X-Device-Id")
	ctxDeviceID, _ := ctx.Value(globalmodel.DeviceIDKey).(string)
	slog.Debug("auth.create_realtor.debug",
		"device_token", request.DeviceToken,
		"ip", reqContext.IPAddress,
		"user_agent", reqContext.UserAgent,
		"header_device_id", headerDeviceID,
		"ctx_device_id", ctxDeviceID,
	)

	// Call service: cria a conta e autentica via SignIn padrão
	tokens, err := ah.userService.CreateRealtor(ctx, user, request.Realtor.Password, request.DeviceToken, reqContext.IPAddress, reqContext.UserAgent)
	if err != nil {
		if derr, ok := err.(utils.DomainError); ok {
			httperrors.SendHTTPErrorObj(c, derr)
			return
		}
		httperrors.SendHTTPErrorObj(c, utils.NewHTTPError(http.StatusConflict, "Failed to create realtor"))
		return
	}

	// Success response
	c.JSON(http.StatusCreated, dto.CreateRealtorResponse{
		Tokens: dto.TokensResponse{
			AccessToken:  tokens.AccessToken,
			RefreshToken: tokens.RefreshToken,
		},
	})
}
