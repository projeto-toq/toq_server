package userhandlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/giulio-alfieri/toq_server/internal/adapter/left/http/dto"
	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (uh *UserHandler) GoHome(c *gin.Context) {
	ctx := c.Request.Context()

	// Get user information from context (set by middleware)
	userInfos, exists := c.Get(string(globalmodel.TokenKey))
	if !exists {
		utils.SendHTTPError(c, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	userInfo := userInfos.(usermodel.UserInfos)

	// Call service to get home data
	user, err := uh.userService.Home(ctx, userInfo.ID)
	if err != nil {
		utils.SendHTTPError(c, http.StatusInternalServerError, "GO_HOME_FAILED", "Failed to get home data")
		return
	}

	// Prepare response with welcome message
	response := dto.GoHomeResponse{
		Message: fmt.Sprintf("Welcome %s. Seu Role é %s, e seu perfil está %v", user.GetNickName(), userInfo.Role, userInfo.ProfileStatus),
	}

	c.JSON(http.StatusOK, response)
}
