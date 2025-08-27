package userhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/giulio-alfieri/toq_server/internal/adapter/left/http/dto"
	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (uh *UserHandler) RejectInvitation(c *gin.Context) {
	ctx := c.Request.Context()

	// Get user information from context (set by middleware)
	userInfos, exists := c.Get(string(globalmodel.TokenKey))
	if !exists {
		utils.SendHTTPError(c, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	userInfo := userInfos.(usermodel.UserInfos)

	// Call service to reject invitation
	if err := uh.userService.RejectInvitation(ctx, userInfo.ID); err != nil {
		utils.SendHTTPError(c, http.StatusInternalServerError, "REJECT_INVITATION_FAILED", "Failed to reject invitation")
		return
	}

	// Prepare response
	response := dto.RejectInvitationResponse{
		Message: "Invitation rejected successfully",
	}

	c.JSON(http.StatusOK, response)
}
