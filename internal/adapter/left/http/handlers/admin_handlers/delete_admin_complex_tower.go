package adminhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	dto "github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

// DeleteAdminComplexTower handles DELETE /admin/complexes/towers
//
//	@Summary	Delete a complex tower
//	@Tags		Admin
//	@Accept		json
//	@Produce	json
//	@Param		request	body	dto.AdminDeleteComplexTowerRequest	true	"Complex tower deletion payload"
//	@Success	204	"Complex tower deleted"
//	@Failure	400	{object}	map[string]any
//	@Failure	401	{object}	map[string]any
//	@Failure	403	{object}	map[string]any
//	@Failure	404	{object}	map[string]any
//	@Failure	500	{object}	map[string]any
//	@Router		/admin/complexes/towers [delete]
func (h *AdminHandler) DeleteAdminComplexTower(c *gin.Context) {
	ctx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)

	var req dto.AdminDeleteComplexTowerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httperrors.SendHTTPErrorObj(c, httperrors.ConvertBindError(err))
		return
	}

	if err := h.complexService.DeleteComplexTower(ctx, req.ID); err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
