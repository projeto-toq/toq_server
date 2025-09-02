package authhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/giulio-alfieri/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/giulio-alfieri/toq_server/internal/adapter/left/http/http_errors"
	httputils "github.com/giulio-alfieri/toq_server/internal/adapter/left/http/utils"
	permissionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/permission_model"
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
//	@Failure		400		{object}	dto.ErrorResponse	"Invalid request format"
//	@Failure		422		{object}	dto.ErrorResponse	"Validation failed"
//	@Failure		409		{object}	dto.ErrorResponse	"User already exists"
//	@Failure		500		{object}	dto.ErrorResponse	"Internal server error"
//	@Router			/auth/realtor [post]
func (ah *AuthHandler) CreateRealtor(c *gin.Context) {
	ctx, spanEnd, err := utils.GenerateTracer(c.Request.Context())
	if err != nil {
		httperrors.SendHTTPErrorObj(c, utils.NewHTTPError(http.StatusInternalServerError, "Failed to generate tracer"))
		return
	}
	defer spanEnd()

	// Parse request
	var request dto.CreateRealtorRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httperrors.SendHTTPErrorObj(c, httputils.MapBindingError(err))
		return
	}

	// Create user model from DTO
	user, err := ah.createUserFromDTO(request.Realtor, permissionmodel.RoleSlugRealtor)
	if err != nil {
		// Parsing de datas (bornAt/creciValidity)
		httperrors.SendHTTPErrorObj(c, utils.NewHTTPError(http.StatusUnprocessableEntity, "Validation failed", map[string]any{
			"errors": []map[string]string{{"field": "realtor.bornAt", "message": "invalid date, expected YYYY-MM-DD"}},
		}))
		return
	}

	// Extract request context for security logging and session metadata
	reqContext := utils.ExtractRequestContext(c)

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
