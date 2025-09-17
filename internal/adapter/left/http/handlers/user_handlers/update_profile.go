package userhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/giulio-alfieri/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/giulio-alfieri/toq_server/internal/adapter/left/http/http_errors"
	"github.com/giulio-alfieri/toq_server/internal/adapter/left/http/middlewares"
	userservices "github.com/giulio-alfieri/toq_server/internal/core/service/user_service"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

// UpdateProfile handles updating user profile
//
//	@Summary		Update user profile (non-sensitive fields)
//	@Description	Update profile fields like nickname, born_at, and address. Email/phone/password have dedicated flows.
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			request	body	dto.UpdateProfileRequest	true	"Update profile payload"
//	@Success		200		{object}	dto.UpdateProfileResponse	"Profile updated"
//	@Failure		400		{object}	dto.ErrorResponse			"Invalid request format"
//	@Failure		401		{object}	dto.ErrorResponse			"Unauthorized"
//	@Failure		422		{object}	dto.ErrorResponse			"Validation error"
//	@Failure		500		{object}	dto.ErrorResponse			"Internal server error"
//	@Router			/user/profile [put]
//	@Security		BearerAuth
func (uh *UserHandler) UpdateProfile(c *gin.Context) {
	ctx, spanEnd, err := utils.GenerateTracer(c.Request.Context())
	if err != nil {
		httperrors.SendHTTPError(c, http.StatusInternalServerError, "TRACER_ERROR", "Failed to generate tracer")
		return
	}
	defer spanEnd()

	// Get user info from Gin context (set by AuthMiddleware)
	userInfos, ok := middlewares.GetUserInfoFromContext(c)
	if !ok {
		// Se chegar aqui, é erro de pipeline (middleware deveria ter setado)
		httperrors.SendHTTPError(c, http.StatusInternalServerError, "INTERNAL_CONTEXT_MISSING", "User context not found")
		return
	}

	// Parse request body using DTO
	var request dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httperrors.SendHTTPError(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request format")
		return
	}

	// Monta input tipado para o serviço com campos opcionais
	in := userservices.UpdateProfileInput{UserID: userInfos.ID}
	if v := request.User.NickName; v != "" {
		in.NickName = &v
	}
	if v := request.User.BornAt; v != "" {
		in.BornAt = &v
	}
	if v := request.User.ZipCode; v != "" {
		in.ZipCode = &v
	}
	if v := request.User.Street; v != "" {
		in.Street = &v
	}
	if v := request.User.Number; v != "" {
		in.Number = &v
	}
	if v := request.User.Complement; v != "" {
		in.Complement = &v
	}
	if v := request.User.Neighborhood; v != "" {
		in.Neighborhood = &v
	}
	if v := request.User.City; v != "" {
		in.City = &v
	}
	if v := request.User.State; v != "" {
		in.State = &v
	}

	// Chama serviço para atualizar perfil
	if err := uh.userService.UpdateProfile(ctx, in); err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	// Success response using DTO
	response := dto.UpdateProfileResponse{
		Message: "Profile updated successfully",
	}
	c.JSON(http.StatusOK, response)
}
