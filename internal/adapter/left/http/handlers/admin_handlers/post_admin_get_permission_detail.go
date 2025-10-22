package adminhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	dto "github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

// PostAdminGetPermissionDetail handles POST /admin/permissions/detail
//
//	@Summary	Get permission detail
//	@Tags		Admin Permissions
//	@Accept		json
//	@Produce	json
//	@Param		request	body	dto.AdminGetPermissionDetailRequest	true	"Permission detail payload"
//	@Success	200	{object}	dto.AdminPermissionSummary
//	@Failure	400	{object}	map[string]any
//	@Failure	401	{object}	map[string]any
//	@Failure	403	{object}	map[string]any
//	@Failure	404	{object}	map[string]any
//	@Failure	500	{object}	map[string]any
//	@Router		/admin/permissions/detail [post]
func (h *AdminHandler) PostAdminGetPermissionDetail(c *gin.Context) {
	ctx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)

	var req dto.AdminGetPermissionDetailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httperrors.SendHTTPErrorObj(c, httperrors.ConvertBindError(err))
		return
	}

	permission, err := h.permissionService.GetPermissionByID(ctx, req.ID)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	c.JSON(http.StatusOK, toAdminPermissionSummary(permission))
}
