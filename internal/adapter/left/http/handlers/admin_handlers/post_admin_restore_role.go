package adminhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	dto "github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

// RestoreAdminRole handles POST /admin/roles/restore
//
//	@Summary    Reactivate a role
//	@Tags       Admin Roles
//	@Accept     json
//	@Produce    json
//	@Param      request body dto.AdminRestoreRoleRequest true "Role reactivation payload"
//	@Success    200 {object} dto.AdminRoleResponse
//	@Failure    400 {object} map[string]any
//	@Failure    401 {object} map[string]any
//	@Failure    403 {object} map[string]any
//	@Failure    404 {object} map[string]any
//	@Failure    409 {object} map[string]any
//	@Failure    500 {object} map[string]any
//	@Router     /admin/roles/restore [post]
func (h *AdminHandler) RestoreAdminRole(c *gin.Context) {
	ctx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)
	var req dto.AdminRestoreRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httperrors.SendHTTPErrorObj(c, httperrors.ConvertBindError(err))
		return
	}

	role, err := h.permissionService.RestoreRole(ctx, req.ID)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	resp := dto.AdminRoleResponse{
		ID:      role.GetID(),
		Message: "Role reactivated",
	}

	c.JSON(http.StatusOK, resp)
}
