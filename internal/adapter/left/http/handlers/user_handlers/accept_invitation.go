package userhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/giulio-alfieri/toq_server/internal/adapter/left/http/dto"
	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (uh *UserHandler) AcceptInvitation(c *gin.Context) {
	ctx := c.Request.Context()

	// Get user information from context (set by middleware)
	userInfos, exists := c.Get(string(globalmodel.TokenKey))
	if !exists {
		utils.SendHTTPError(c, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	userInfo := userInfos.(usermodel.UserInfos)

	// Call service to accept invitation
	if err := uh.userService.AcceptInvitation(ctx, userInfo.ID); err != nil {
		utils.SendHTTPError(c, http.StatusInternalServerError, "ACCEPT_INVITATION_FAILED", "Failed to accept invitation")
		return
	}

	// Prepare response
	response := dto.AcceptInvitationResponse{
		Message: "Invitation accepted successfully",
	}

	c.JSON(http.StatusOK, response)
}
