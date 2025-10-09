package userhandlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/middlewares"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

func (uh *UserHandler) GoHome(c *gin.Context) {
	ctx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)

	// Get user information from context (set by middleware)
	userInfo, ok := middlewares.GetUserInfoFromContext(c)
	if !ok {
		// Se chegar aqui, é erro de pipeline (middleware deveria ter setado)
		httperrors.SendHTTPError(c, http.StatusInternalServerError, "INTERNAL_CONTEXT_MISSING", "User context not found")
		return
	}

	// Call service to get home data
	user, err := uh.userService.Home(ctx, userInfo.ID)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	// Prepare response with welcome message
	currentRole := coreutils.GetUserRoleSlugFromUserRole(user.GetActiveRole())
	// TODO: Código temporário comentado - RoleStatus removido dos tokens
	// isActive := utils.IsProfileActiveFromStatus(userInfo.RoleStatus)

	response := dto.GoHomeResponse{
		// TODO: Mensagem simplificada até refatoração completa do handler
		Message: fmt.Sprintf("Welcome %s. Seu Role é %s", user.GetNickName(), currentRole),
		// Message: fmt.Sprintf("Welcome %s. Seu Role é %s, e seu perfil está %v", user.GetNickName(), currentRole, isActive),
	}

	c.JSON(http.StatusOK, response)
}
