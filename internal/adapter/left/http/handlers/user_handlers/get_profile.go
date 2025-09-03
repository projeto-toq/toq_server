package userhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/giulio-alfieri/toq_server/internal/adapter/left/http/converters"
	dto "github.com/giulio-alfieri/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/giulio-alfieri/toq_server/internal/adapter/left/http/http_errors"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

// Ensure Swag can resolve dto.ErrorResponse referenced in annotations without affecting runtime
type _ = dto.ErrorResponse

// GetProfile handles getting user profile
//
//	@Summary		Get user profile
//	@Description	Get the current authenticated user's profile information
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	dto.UserProfileResponse	"Profile data with user information"
//	@Failure		401	{object}	dto.ErrorResponse	"Unauthorized"
//	@Failure		403	{object}	dto.ErrorResponse	"Forbidden"
//	@Failure		404	{object}	dto.ErrorResponse	"User not found"
//	@Failure		500	{object}	dto.ErrorResponse	"Internal server error"
//	@Router			/user/profile [get]
//	@Security		BearerAuth
func (uh *UserHandler) GetProfile(c *gin.Context) {
	ctx, spanEnd, err := utils.GenerateTracer(c.Request.Context())
	if err != nil {
		httperrors.SendHTTPError(c, http.StatusInternalServerError, "TRACER_ERROR", "Failed to generate tracer")
		return
	}
	defer spanEnd()

	// Call service
	user, err := uh.userService.GetProfile(ctx)
	if err != nil {
		// Standardized DomainError passthrough
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	// Success response com DTO tipado e conversão segura
	resp := converters.ToUserProfileResponse(user)
	c.JSON(http.StatusOK, resp)
}
