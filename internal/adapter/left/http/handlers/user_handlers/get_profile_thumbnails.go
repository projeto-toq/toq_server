package userhandlers

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/giulio-alfieri/toq_server/internal/adapter/left/http/dto"
	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (uh *UserHandler) GetProfileThumbnails(c *gin.Context) {
	ctx := c.Request.Context()

	// Get user information from context (set by middleware)
	userInfos, exists := c.Get(string(globalmodel.TokenKey))
	if !exists {
		utils.SendHTTPError(c, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	userInfo := userInfos.(usermodel.UserInfos)

	// Call service to get profile thumbnails
	thumbnails, err := uh.userService.GetProfileThumbnails(ctx, userInfo.ID)
	if err != nil {
		slog.Error("failed to get profile thumbnails", "error", err, "userID", userInfo.ID)
		utils.SendHTTPError(c, http.StatusInternalServerError, "GET_PROFILE_THUMBNAILS_FAILED", "Failed to get profile thumbnails")
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
