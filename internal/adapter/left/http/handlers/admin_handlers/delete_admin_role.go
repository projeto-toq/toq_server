package adminhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	dto "github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

// DeleteAdminRole handles DELETE /admin/roles
//
//	@Summary      Delete (deactivate) a role
//	@Tags         Admin Roles
//	@Accept       json
//	@Produce      json
//	@Param        request  body  dto.AdminDeleteRoleRequest  true  "Role deletion payload"
//	@Success      204  "Role deactivated"
//	@Failure      400  {object}  map[string]any
//	@Failure      401  {object}  map[string]any
//	@Failure      403  {object}  map[string]any
//	@Failure      404  {object}  map[string]any
//	@Failure      409  {object}  map[string]any
//	@Failure      500  {object}  map[string]any
//	@Router       /admin/roles [delete]
func (h *AdminHandler) DeleteAdminRole(c *gin.Context) {
	ctx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)
	var req dto.AdminDeleteRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httperrors.SendHTTPErrorObj(c, httperrors.ConvertBindError(err))
		return
	}
	if req.ID <= 0 {
		httperrors.SendHTTPErrorObj(c, coreutils.ValidationError("id", "Invalid role id"))
		return
	}

	if err := h.permissionService.DeleteRole(ctx, req.ID); err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
