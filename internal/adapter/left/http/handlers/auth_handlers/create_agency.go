package authhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/giulio-alfieri/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/giulio-alfieri/toq_server/internal/adapter/left/http/http_errors"
	permissionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/permission_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

// CreateAgency handles agency account creation (public endpoint)
//
//	@Summary		Create agency account
//	@Description	Create a new agency account with user information
//	@Tags			Authentication
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.CreateAgencyRequest	true	"Agency creation data"
//	@Success		201		{object}	dto.CreateAgencyResponse
//	@Failure		400		{object}	dto.ErrorResponse	"Invalid request format"
//	@Failure		422		{object}	dto.ErrorResponse	"Validation failed"
//	@Failure		409		{object}	dto.ErrorResponse	"User already exists"
//	@Failure		500		{object}	dto.ErrorResponse	"Internal server error"
//	@Router			/auth/agency [post]
func (ah *AuthHandler) CreateAgency(c *gin.Context) {
	ctx, spanEnd, err := utils.GenerateTracer(c.Request.Context())
	if err != nil {
		httperrors.SendHTTPErrorObj(c, utils.NewHTTPError(http.StatusInternalServerError, "Failed to generate tracer"))
		return
	}
	defer spanEnd()

	// Parse request
	var request dto.CreateAgencyRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httperrors.SendHTTPErrorObj(c, utils.NewHTTPError(http.StatusBadRequest, "Invalid request format"))
		return
	}

	// Create user model from DTO
	user, err := ah.createUserFromDTO(request.Agency, permissionmodel.RoleSlugAgency)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, utils.NewHTTPError(http.StatusBadRequest, "Invalid user data"))
		return
	}

	// Call service
	tokens, err := ah.userService.CreateAgency(ctx, user)
	if err != nil {
		if derr, ok := err.(utils.DomainError); ok {
			httperrors.SendHTTPErrorObj(c, derr)
			return
		}
		httperrors.SendHTTPErrorObj(c, utils.NewHTTPError(http.StatusConflict, "Failed to create agency"))
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
