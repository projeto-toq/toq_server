package authhandlers

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	httputils "github.com/projeto-toq/toq_server/internal/adapter/left/http/utils"
	globalmodel "github.com/projeto-toq/toq_server/internal/core/model/global_model"
	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

// CreateOwner handles owner account creation (public endpoint)
//
//	@Summary		Create owner account
//	@Description	Create a new owner account with user information. Address fields come from CEP lookup, except number, neighborhood and complement which honor the request payload.
//	@Tags			Authentication
//	@Accept			json
//	@Produce		json
//	@Param		X-Device-Id	header	string	false	"Device ID (UUIDv4). Optional but recommended when providing deviceToken for automatic sign-in"
//	@Param			request	body		dto.CreateOwnerRequest	true	"Owner creation data (include optional deviceToken for push notifications)"
//	@Success		201		{object}	dto.CreateOwnerResponse
//	@Failure		400		{object}	dto.ErrorResponse	"Validation error (invalid input data)"
//	@Failure		409		{object}	dto.ErrorResponse	"User already exists"
//	@Failure		500		{object}	dto.ErrorResponse	"Internal server error"
//	@Failure		422		{object}	dto.ErrorResponse	"Validation failed (see details)"
//	@Router			/auth/owner [post]
func (ah *AuthHandler) CreateOwner(c *gin.Context) {
	// Observação: tracing de request já é provido por TelemetryMiddleware; evitamos spans duplicados aqui.
	ctx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)
	logger := coreutils.LoggerFromContext(ctx)

	// Parse request
	var request dto.CreateOwnerRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httperrors.SendHTTPErrorObj(c, httputils.MapBindingError(err))
		return
	}

	// Validate and parse date fields with precise attribution
	bornAt, creciValidity, derr := httputils.ValidateUserDates(request.Owner, "owner")
	if derr != nil {
		httperrors.SendHTTPErrorObj(c, derr)
		return
	}

	// Create user model from DTO (using parsed dates)
	user, err := ah.createUserFromDTO(request.Owner, bornAt, creciValidity)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, coreutils.NewHTTPError(http.StatusUnprocessableEntity, "Validation failed"))
		return
	}

	// Extract request context for security logging and session metadata
	reqContext := coreutils.ExtractRequestContext(c)

	rawDeviceID := c.GetHeader("X-Device-Id")
	trimmedDeviceID := strings.TrimSpace(rawDeviceID)
	if trimmedDeviceID == "" {
		httperrors.SendHTTPErrorObj(c, coreutils.NewHTTPErrorWithSource(http.StatusBadRequest, "X-Device-Id header is required"))
		return
	}
	if _, err := uuid.Parse(trimmedDeviceID); err != nil {
		httperrors.SendHTTPErrorObj(c, coreutils.NewHTTPErrorWithSource(http.StatusBadRequest, "X-Device-Id must be a valid UUID"))
		return
	}

	trimmedDeviceToken := strings.TrimSpace(request.DeviceToken)
	if trimmedDeviceToken == "" {
		httperrors.SendHTTPErrorObj(c, coreutils.NewHTTPErrorWithSource(http.StatusBadRequest, "deviceToken is required"))
		return
	}

	ctx = context.WithValue(ctx, globalmodel.DeviceIDKey, trimmedDeviceID)
	c.Request = c.Request.WithContext(ctx)

	// Debug: rastrear valores de contexto antes de chamar o serviço
	ctxDeviceID, _ := ctx.Value(globalmodel.DeviceIDKey).(string)
	logger.Debug("auth.create_owner.debug",
		"device_token", trimmedDeviceToken,
		"ip", reqContext.IPAddress,
		"user_agent", reqContext.UserAgent,
		"header_device_id", trimmedDeviceID,
		"ctx_device_id", ctxDeviceID,
	)

	// Call service: cria a conta e autentica via SignIn padrão
	tokens, err := ah.userService.CreateOwner(ctx, user, request.Owner.Password, trimmedDeviceToken, reqContext.IPAddress, reqContext.UserAgent)
	if err != nil {
		if derr, ok := err.(coreutils.DomainError); ok {
			httperrors.SendHTTPErrorObj(c, derr)
			return
		}
		// Fallback: conflito genérico
		httperrors.SendHTTPErrorObj(c, coreutils.NewHTTPError(http.StatusConflict, "Failed to create owner"))
		return
	}

	// Success response
	c.JSON(http.StatusCreated, dto.CreateOwnerResponse{
		Tokens: dto.TokensResponse{
			AccessToken:  tokens.AccessToken,
			RefreshToken: tokens.RefreshToken,
		},
	})
}

// createUserFromDTO converts DTO to User model
func (ah *AuthHandler) createUserFromDTO(dtoUser dto.UserCreateRequest, bornAt time.Time, creciValidity *time.Time) (usermodel.UserInterface, error) {
	user := usermodel.NewUser()

	// Set user data
	user.SetFullName(dtoUser.FullName)
	user.SetNickName(dtoUser.NickName)
	user.SetNationalID(dtoUser.NationalID)
	user.SetCreciNumber(dtoUser.CreciNumber)
	user.SetCreciState(dtoUser.CreciState)
	if creciValidity != nil && !creciValidity.IsZero() {
		user.SetCreciValidity(*creciValidity)
	}
	user.SetBornAt(bornAt)
	user.SetPhoneNumber(dtoUser.PhoneNumber)
	user.SetEmail(dtoUser.Email)
	user.SetZipCode(dtoUser.ZipCode)
	user.SetStreet(dtoUser.Street)
	user.SetNumber(dtoUser.Number)
	user.SetComplement(dtoUser.Complement)
	user.SetNeighborhood(dtoUser.Neighborhood)
	user.SetCity(dtoUser.City)
	user.SetState(dtoUser.State)
	user.SetPassword(dtoUser.Password)

	// Note: Active role will be set by the service layer, not the handler

	user.SetOptStatus(false) // Default opt-in status
	user.SetDeleted(false)

	return user, nil
}
