package adminhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	httpconv "github.com/projeto-toq/toq_server/internal/adapter/left/http/converters"
	dto "github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	complexservices "github.com/projeto-toq/toq_server/internal/core/service/complex_service"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

// PutAdminUpdateComplexTower handles PUT /admin/complexes/towers
//
//	@Summary	Update a complex tower
//	@Tags		Admin
//	@Accept		json
//	@Produce	json
//	@Param		request	body	dto.AdminUpdateComplexTowerRequest	true	"Complex tower payload"
//	@Success	200	{object}	dto.ComplexTowerResponse
//	@Failure	400	{object}	map[string]any
//	@Failure	401	{object}	map[string]any
//	@Failure	403	{object}	map[string]any
//	@Failure	404	{object}	map[string]any
//	@Failure	500	{object}	map[string]any
//	@Router		/admin/complexes/towers [put]
func (h *AdminHandler) PutAdminUpdateComplexTower(c *gin.Context) {
	ctx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)

	var req dto.AdminUpdateComplexTowerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httperrors.SendHTTPErrorObj(c, httperrors.ConvertBindError(err))
		return
	}

	input := complexservices.UpdateComplexTowerInput{
		ID:            req.ID,
		ComplexID:     req.ComplexID,
		Tower:         req.Tower,
		Floors:        req.Floors,
		TotalUnits:    req.TotalUnits,
		UnitsPerFloor: req.UnitsPerFloor,
	}

	tower, err := h.complexService.UpdateComplexTower(ctx, input)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	c.JSON(http.StatusOK, httpconv.ToComplexTowerResponse(tower))
}
