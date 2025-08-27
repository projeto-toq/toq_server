package userhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

// SignOut handles user sign out
func (uh *UserHandler) SignOut(c *gin.Context) {
	ctx, spanEnd, err := utils.GenerateTracer(c.Request.Context())
	if err != nil {
		utils.SendHTTPError(c, http.StatusInternalServerError, "TRACER_ERROR", "Failed to generate tracer")
		return
	}
	defer spanEnd()

	// Get user info from context (set by auth middleware)
	infos, exists := c.Get(string(globalmodel.TokenKey))
	if !exists {
		utils.SendHTTPError(c, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	userInfos := infos.(usermodel.UserInfos)

	// Parse request body
	var request struct {
		DeviceToken  string `json:"device_token"`
		RefreshToken string `json:"refresh_token"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utils.SendHTTPError(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request format")
		return
	}

	// Call service
	if err := uh.userService.SignOut(ctx, userInfos.ID, request.DeviceToken, request.RefreshToken); err != nil {
		utils.SendHTTPError(c, http.StatusInternalServerError, "SIGNOUT_FAILED", "Failed to sign out")
		return
	}

	// Success response
	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully signed out",
	})
}
