package userhandlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/giulio-alfieri/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/giulio-alfieri/toq_server/internal/adapter/left/http/http_errors"
	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (uh *UserHandler) GoHome(c *gin.Context) {
	ctx := c.Request.Context()

	// Get user information from context (set by middleware)
	userInfos, exists := c.Get(string(globalmodel.TokenKey))
	if !exists {
		httperrors.SendHTTPError(c, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	userInfo := userInfos.(usermodel.UserInfos)

	// Call service to get home data
	user, err := uh.userService.Home(ctx, userInfo.ID)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	// Prepare response with welcome message
	currentRole := utils.GetUserRoleSlugFromUserRole(user.GetActiveRole())
	// TODO: Código temporário comentado - RoleStatus removido dos tokens
	// isActive := utils.IsProfileActiveFromStatus(userInfo.RoleStatus)

	response := dto.GoHomeResponse{
		// TODO: Mensagem simplificada até refatoração completa do handler
		Message: fmt.Sprintf("Welcome %s. Seu Role é %s", user.GetNickName(), currentRole),
		// Message: fmt.Sprintf("Welcome %s. Seu Role é %s, e seu perfil está %v", user.GetNickName(), currentRole, isActive),
	}

	c.JSON(http.StatusOK, response)
}
