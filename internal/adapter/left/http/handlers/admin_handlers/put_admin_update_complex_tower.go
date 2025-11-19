package adminhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	httpconv "github.com/projeto-toq/toq_server/internal/adapter/left/http/converters"
	dto "github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	propertycoverageservice "github.com/projeto-toq/toq_server/internal/core/service/property_coverage_service"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

// PutAdminUpdateComplexTower handles PUT /admin/complexes/towers
//
//	@Summary	Update a complex tower
//	@Tags		Admin Complexes
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

	input := propertycoverageservice.UpdateComplexTowerInput{
		ID: req.ID,
		CreateComplexTowerInput: propertycoverageservice.CreateComplexTowerInput{
			VerticalComplexID: req.ComplexID,
			Tower:             req.Tower,
			Floors:            req.Floors,
			TotalUnits:        req.TotalUnits,
			UnitsPerFloor:     req.UnitsPerFloor,
		},
	}

	tower, svcErr := h.propertyCoverageService.UpdateComplexTower(ctx, input)
	if svcErr != nil {
		httperrors.SendHTTPErrorObj(c, svcErr)
		return
	}

	c.JSON(http.StatusOK, httpconv.ToComplexTowerResponse(tower))
}
