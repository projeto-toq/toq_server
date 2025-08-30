package userhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
) // SignOut handles user sign out
// @Summary		Sign out user
// @Description	Sign out the current user and invalidate their session tokens
// @Tags			User
// @Accept			json
// @Produce		json
// @Param			request	body		object					true	"Sign out data"
// @Param			request.device_token	body		string	false	"Device token to invalidate"
// @Param			request.refresh_token	body		string	false	"Refresh token to invalidate"
// @Success		200		{object}	map[string]string	"Sign out confirmation message"
// @Failure		400		{object}	dto.ErrorResponse	"Invalid request format"
// @Failure		401		{object}	dto.ErrorResponse	"Unauthorized"
// @Failure		403		{object}	dto.ErrorResponse	"Forbidden"
// @Failure		500		{object}	dto.ErrorResponse	"Internal server error"
// @Router			/user/signout [post]
// @Security		BearerAuth
func (uh *UserHandler) SignOut(c *gin.Context) {
	ctx, spanEnd, err := utils.GenerateTracer(c.Request.Context())
	if err != nil {
		utils.SendHTTPError(c, http.StatusInternalServerError, "TRACER_ERROR", "Failed to generate tracer")
		return
	}
	defer spanEnd()

	// Get user info from context using Context Utils (set by auth middleware)
	userInfos, err := utils.GetUserInfoFromGinContext(c)
	if err != nil {
		utils.SendHTTPError(c, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

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
