package adminhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	dto "github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	permissionservice "github.com/projeto-toq/toq_server/internal/core/service/permission_service"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

// PutAdminUpdateRolePermission handles PUT /admin/role-permissions
//
//	@Summary	Update a role-permission association
//	@Tags		Admin
//	@Accept		json
//	@Produce	json
//	@Param		payload	body	dto.AdminUpdateRolePermissionRequest	true	"Role-permission payload"
//	@Success	200	{object}	dto.AdminRolePermissionResponse
//	@Failure	400	{object}	map[string]any
//	@Failure	401	{object}	map[string]any
//	@Failure	403	{object}	map[string]any
//	@Failure	404	{object}	map[string]any
//	@Failure	500	{object}	map[string]any
//	@Router		/admin/role-permissions [put]
func (h *AdminHandler) PutAdminUpdateRolePermission(c *gin.Context) {
	ctx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)

	var req dto.AdminUpdateRolePermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httperrors.SendHTTPErrorObj(c, httperrors.ConvertBindError(err))
		return
	}

	input := permissionservice.UpdateRolePermissionInput{
		ID:      req.ID,
		Granted: req.Granted,
	}

	rolePermission, err := h.permissionService.UpdateRolePermission(ctx, input)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	resp := dto.AdminRolePermissionResponse{
		ID:      rolePermission.GetID(),
		Message: "role permission updated",
	}

	c.JSON(http.StatusOK, resp)
}
