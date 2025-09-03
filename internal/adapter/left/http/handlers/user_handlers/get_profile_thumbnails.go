package userhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/giulio-alfieri/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/giulio-alfieri/toq_server/internal/adapter/left/http/http_errors"
	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
)

// GetProfileThumbnails returns signed URLs for all profile photo sizes
//
//	@Summary      Get profile photo thumbnails
//	@Description  Returns signed URLs for original, small, medium, and large profile photos
//	@Tags         User
//	@Produce      json
//	@Success      200  {object}  dto.GetProfileThumbnailsResponse
//	@Failure      401  {object}  dto.ErrorResponse  "Unauthorized"
//	@Failure      403  {object}  dto.ErrorResponse  "Forbidden"
//	@Failure      500  {object}  dto.ErrorResponse  "Internal server error"
//	@Router       /user/profile/thumbnails [get]
//	@Security     BearerAuth
func (uh *UserHandler) GetProfileThumbnails(c *gin.Context) {
	ctx := c.Request.Context()

	// Get user information from context (set by middleware)
	userInfos, exists := c.Get(string(globalmodel.TokenKey))
	if !exists {
		httperrors.SendHTTPError(c, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	userInfo := userInfos.(usermodel.UserInfos)

	// Call service to get profile thumbnails
	thumbnails, err := uh.userService.GetProfileThumbnails(ctx, userInfo.ID)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	// Prepare response
	response := dto.GetProfileThumbnailsResponse{
		OriginalURL: thumbnails.OriginalURL,
		SmallURL:    thumbnails.SmallURL,
		MediumURL:   thumbnails.MediumURL,
		LargeURL:    thumbnails.LargeURL,
	}

	c.JSON(http.StatusOK, response)
}
