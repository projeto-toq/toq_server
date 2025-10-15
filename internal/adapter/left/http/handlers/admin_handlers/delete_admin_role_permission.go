package adminhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	dto "github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

// DeleteAdminRolePermission handles DELETE /admin/role-permissions
//
//	@Summary	Delete a role-permission association
//	@Tags		Admin
//	@Accept		json
//	@Produce	json
//	@Param		payload	body	dto.AdminDeleteRolePermissionRequest	true	"Role-permission payload"
//	@Success	200	{object}	dto.AdminRolePermissionResponse
//	@Failure	400	{object}	map[string]any
//	@Failure	401	{object}	map[string]any
//	@Failure	403	{object}	map[string]any
//	@Failure	404	{object}	map[string]any
//	@Failure	500	{object}	map[string]any
//	@Router		/admin/role-permissions [delete]
func (h *AdminHandler) DeleteAdminRolePermission(c *gin.Context) {
	ctx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)

	var req dto.AdminDeleteRolePermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httperrors.SendHTTPErrorObj(c, httperrors.ConvertBindError(err))
		return
	}

	if err := h.permissionService.DeleteRolePermission(ctx, req.ID); err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	resp := dto.AdminRolePermissionResponse{
		ID:      req.ID,
		Message: "role permission deleted",
	}

	c.JSON(http.StatusOK, resp)
}
