package authhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/giulio-alfieri/toq_server/internal/adapter/left/http/dto"
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
//	@Param			request	body		dto.CreateRealtorRequest	true	"Realtor creation data"
//	@Success		201		{object}	dto.CreateRealtorResponse
//	@Failure		400		{object}	dto.ErrorResponse	"Invalid request format"
//	@Failure		409		{object}	dto.ErrorResponse	"User already exists"
//	@Failure		500		{object}	dto.ErrorResponse	"Internal server error"
//	@Router			/auth/realtor [post]
func (ah *AuthHandler) CreateRealtor(c *gin.Context) {
	ctx, spanEnd, err := utils.GenerateTracer(c.Request.Context())
	if err != nil {
		utils.SendHTTPError(c, http.StatusInternalServerError, "TRACER_ERROR", "Failed to generate tracer")
		return
	}
	defer spanEnd()

	// Parse request
	var request dto.CreateRealtorRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.SendHTTPError(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request format")
		return
	}

	// Create user model from DTO
	user, err := ah.createUserFromDTO(request.Realtor, permissionmodel.RoleSlugRealtor)
	if err != nil {
		utils.SendHTTPError(c, http.StatusBadRequest, "INVALID_USER_DATA", "Invalid user data")
		return
	}

	// Call service
	tokens, err := ah.userService.CreateRealtor(ctx, user)
	if err != nil {
		utils.SendHTTPError(c, http.StatusConflict, "USER_CREATION_FAILED", "Failed to create realtor")
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
