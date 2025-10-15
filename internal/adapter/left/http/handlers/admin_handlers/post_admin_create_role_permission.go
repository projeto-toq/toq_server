package adminhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	dto "github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	permissionservice "github.com/projeto-toq/toq_server/internal/core/service/permission_service"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

// PostAdminCreateRolePermission handles POST /admin/role-permissions
//
//	@Summary	Create a role-permission association
//	@Tags		Admin
//	@Accept		json
//	@Produce	json
//	@Param		payload	body	dto.AdminCreateRolePermissionRequest	true	"Role-permission payload"
//	@Success	201	{object}	dto.AdminRolePermissionResponse
//	@Failure	400	{object}	map[string]any
//	@Failure	401	{object}	map[string]any
//	@Failure	403	{object}	map[string]any
//	@Failure	404	{object}	map[string]any
//	@Failure	409	{object}	map[string]any
//	@Failure	500	{object}	map[string]any
//	@Router		/admin/role-permissions [post]
func (h *AdminHandler) PostAdminCreateRolePermission(c *gin.Context) {
	ctx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)

	var req dto.AdminCreateRolePermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httperrors.SendHTTPErrorObj(c, httperrors.ConvertBindError(err))
		return
	}

	input := permissionservice.CreateRolePermissionInput{
		RoleID:       req.RoleID,
		PermissionID: req.PermissionID,
		Granted:      req.Granted,
		Conditions:   req.Conditions,
	}

	rolePermission, err := h.permissionService.CreateRolePermission(ctx, input)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	resp := dto.AdminRolePermissionResponse{
		ID:      rolePermission.GetID(),
		Message: "role permission created",
	}

	c.JSON(http.StatusCreated, resp)
}
