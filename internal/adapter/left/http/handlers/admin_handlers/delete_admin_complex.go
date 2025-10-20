package adminhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	dto "github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

// DeleteAdminComplex handles DELETE /admin/complexes
//
//	@Summary	Delete a complex
//	@Tags		Admin
//	@Accept		json
//	@Produce	json
//	@Param		request	body	dto.AdminDeleteComplexRequest	true	"Complex deletion payload"
//	@Success	204	"Complex deleted"
//	@Failure	400	{object}	map[string]any
//	@Failure	401	{object}	map[string]any
//	@Failure	403	{object}	map[string]any
//	@Failure	404	{object}	map[string]any
//	@Failure	500	{object}	map[string]any
//	@Router		/admin/complexes [delete]
func (h *AdminHandler) DeleteAdminComplex(c *gin.Context) {
	ctx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)

	var req dto.AdminDeleteComplexRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httperrors.SendHTTPErrorObj(c, httperrors.ConvertBindError(err))
		return
	}

	if err := h.complexService.DeleteComplex(ctx, req.ID); err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
