package adminhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	httpconv "github.com/projeto-toq/toq_server/internal/adapter/left/http/converters"
	dto "github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

// PostAdminGetComplexTowerDetail handles POST /admin/complexes/towers/detail
//
//	@Summary	Get complex tower detail
//	@Tags		Admin Complexes
//	@Accept		json
//	@Produce	json
//	@Param		request	body	dto.AdminGetComplexTowerDetailRequest	true	"Complex tower detail payload"
//	@Success	200	{object}	dto.ComplexTowerResponse
//	@Failure	400	{object}	map[string]any
//	@Failure	401	{object}	map[string]any
//	@Failure	403	{object}	map[string]any
//	@Failure	404	{object}	map[string]any
//	@Failure	500	{object}	map[string]any
//	@Router		/admin/complexes/towers/detail [post]
func (h *AdminHandler) PostAdminGetComplexTowerDetail(c *gin.Context) {
	ctx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)

	var req dto.AdminGetComplexTowerDetailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httperrors.SendHTTPErrorObj(c, httperrors.ConvertBindError(err))
		return
	}

	tower, err := h.propertyCoverageService.GetComplexTowerDetail(ctx, req.ID)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	c.JSON(http.StatusOK, httpconv.ToComplexTowerResponse(tower))
}
