package adminhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	dto "github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

// DeleteAdminPermission handles DELETE /admin/permissions
//
//	@Summary	Hard delete a permission
//	@Tags		Admin
//	@Accept		json
//	@Produce	json
//	@Param		payload	body	dto.AdminDeletePermissionRequest	true	"Permission identifier"
//	@Success	200	{object}	dto.AdminPermissionResponse
//	@Failure	400	{object}	map[string]any
//	@Failure	401	{object}	map[string]any
//	@Failure	403	{object}	map[string]any
//	@Failure	404	{object}	map[string]any
//	@Failure	500	{object}	map[string]any
//	@Router		/admin/permissions [delete]
func (h *AdminHandler) DeleteAdminPermission(c *gin.Context) {
	ctx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)

	var req dto.AdminDeletePermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httperrors.SendHTTPErrorObj(c, httperrors.ConvertBindError(err))
		return
	}

	if err := h.permissionService.DeletePermission(ctx, req.ID); err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	resp := dto.AdminPermissionResponse{
		ID:      req.ID,
		Message: "permission deleted",
	}

	c.JSON(http.StatusOK, resp)
}
