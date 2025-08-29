package userhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/giulio-alfieri/toq_server/internal/adapter/left/http/dto"
	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (uh *UserHandler) GetOnboardingStatus(c *gin.Context) {
	ctx := c.Request.Context()

	// Get user information from context (set by middleware)
	userInfos, exists := c.Get(string(globalmodel.TokenKey))
	if !exists {
		utils.SendHTTPError(c, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	userInfo := userInfos.(usermodel.UserInfos)

	// Call service to get onboarding status
	status, reason, err := uh.userService.GetOnboardingStatus(ctx, userInfo.ID)
	if err != nil {
		utils.SendHTTPError(c, http.StatusInternalServerError, "GET_ONBOARDING_STATUS_FAILED", "Failed to get onboarding status")
		return
	}

	// Prepare response
	response := dto.GetOnboardingStatusResponse{
		Status: status,
		Reason: reason,
	}

	c.JSON(http.StatusOK, response)
}
