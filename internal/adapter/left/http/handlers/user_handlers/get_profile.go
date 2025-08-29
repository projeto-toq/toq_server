package userhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

// GetProfile handles getting user profile
//
//	@Summary		Get user profile
//	@Description	Get the current authenticated user's profile information
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	map[string]interface{}	"Profile data with user information"
//	@Failure		401	{object}	dto.ErrorResponse	"Unauthorized"
//	@Failure		403	{object}	dto.ErrorResponse	"Forbidden"
//	@Failure		500	{object}	dto.ErrorResponse	"Internal server error"
//	@Router			/user/profile [get]
//	@Security		BearerAuth
func (uh *UserHandler) GetProfile(c *gin.Context) {
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

	// Call service
	user, err := uh.userService.GetProfile(ctx, userInfos.ID)
	if err != nil {
		utils.SendHTTPError(c, http.StatusInternalServerError, "GET_PROFILE_FAILED", "Failed to get profile")
		return
	}

	// Success response
	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"id":           user.GetID(),
			"email":        user.GetEmail(),
			"phone_number": user.GetPhoneNumber(),
			"full_name":    user.GetFullName(),
			"nick_name":    user.GetNickName(),
			"national_id":  user.GetNationalID(),
			"active_role": gin.H{
				"id":            user.GetActiveRole().GetID(),
				"role":          user.GetActiveRole().GetRole(),
				"active":        user.GetActiveRole().IsActive(),
				"status":        user.GetActiveRole().GetStatus(),
				"status_reason": user.GetActiveRole().GetStatusReason(),
			},
			"born_at":  user.GetBornAt().Format("2006-01-02"),
			"zip_code": user.GetZipCode(),
			"street":   user.GetStreet(),
			"city":     user.GetCity(),
			"state":    user.GetState(),
		},
	})
}
