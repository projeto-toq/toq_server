package adminhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	dto "github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

// DeleteAdminComplexSize handles DELETE /admin/complexes/sizes
//
//	@Summary	Delete a complex size
//	@Tags		Admin Complexes
//	@Accept		json
//	@Produce	json
//	@Param		request	body	dto.AdminDeleteComplexSizeRequest	true	"Complex size deletion payload"
//	@Success	204	"Complex size deleted"
//	@Failure	400	{object}	map[string]any
//	@Failure	401	{object}	map[string]any
//	@Failure	403	{object}	map[string]any
//	@Failure	404	{object}	map[string]any
//	@Failure	500	{object}	map[string]any
//	@Router		/admin/complexes/sizes [delete]
func (h *AdminHandler) DeleteAdminComplexSize(c *gin.Context) {
	ctx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)

	var req dto.AdminDeleteComplexSizeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httperrors.SendHTTPErrorObj(c, httperrors.ConvertBindError(err))
		return
	}

	if err := h.propertyCoverageService.DeleteComplexSize(ctx, req.ID); err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
