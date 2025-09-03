package userhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/giulio-alfieri/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/giulio-alfieri/toq_server/internal/adapter/left/http/http_errors"
	permissionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/permission_model"
	coreutils "github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (uh *UserHandler) AddAlternativeUserRole(c *gin.Context) {
	ctx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)

	// Parse request body
	var request dto.AddAlternativeUserRoleRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httperrors.SendHTTPError(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request format")
		return
	}

	// Get user to determine current role - precisamos buscar o usuário para obter o role
	user, err := uh.userService.GetProfile(ctx)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	// Determine alternative role based on current active role
	var alternativeRole permissionmodel.RoleSlug
	currentRole := coreutils.GetUserRoleSlugFromUserRole(user.GetActiveRole())
	if currentRole == permissionmodel.RoleSlugOwner {
		alternativeRole = permissionmodel.RoleSlugRealtor
	} else {
		alternativeRole = permissionmodel.RoleSlugOwner
	}

	// Call service to add alternative role
	// Ainda passamos o userID explicitamente aqui pois o caso de uso adiciona role a um usuário específico (o próprio)
	// Poderíamos migrar para SSOT também, mas mantemos assinatura existente por ora.
	userInfo, _ := coreutils.GetUserInfoFromContext(ctx)
	if err := uh.userService.AddAlternativeRole(ctx, userInfo.ID, alternativeRole, request.CreciNumber, request.CreciState, request.CreciValidity); err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	// Prepare response
	response := dto.AddAlternativeUserRoleResponse{
		Message: "Alternative user role added successfully",
	}

	c.JSON(http.StatusOK, response)
}
