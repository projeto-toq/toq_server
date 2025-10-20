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

// PostAdminCreateComplexTower handles POST /admin/complexes/towers
//
//	@Summary	Create a complex tower
//	@Tags		Admin
//	@Accept		json
//	@Produce	json
//	@Param		request	body	dto.AdminCreateComplexTowerRequest	true	"Complex tower payload"
//	@Success	201	{object}	dto.ComplexTowerResponse
//	@Failure	400	{object}	map[string]any
//	@Failure	401	{object}	map[string]any
//	@Failure	403	{object}	map[string]any
//	@Failure	404	{object}	map[string]any
//	@Failure	500	{object}	map[string]any
//	@Router		/admin/complexes/towers [post]
func (h *AdminHandler) PostAdminCreateComplexTower(c *gin.Context) {
	ctx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)

	var req dto.AdminCreateComplexTowerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httperrors.SendHTTPErrorObj(c, httperrors.ConvertBindError(err))
		return
	}

	input := complexservices.CreateComplexTowerInput{
		ComplexID:     req.ComplexID,
		Tower:         req.Tower,
		Floors:        req.Floors,
		TotalUnits:    req.TotalUnits,
		UnitsPerFloor: req.UnitsPerFloor,
	}

	tower, err := h.complexService.CreateComplexTower(ctx, input)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	c.JSON(http.StatusCreated, httpconv.ToComplexTowerResponse(tower))
}
