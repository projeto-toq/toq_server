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

// PostAdminGetComplexDetail handles POST /admin/complexes/detail
//
//	@Summary	Get complex detail
//	@Tags		Admin Complexes
//	@Accept		json
//	@Produce	json
//	@Param		request	body	dto.AdminGetComplexDetailRequest	true	"Complex detail payload"
//	@Success	200	{object}	dto.ComplexResponse
//	@Failure	400	{object}	map[string]any
//	@Failure	401	{object}	map[string]any
//	@Failure	403	{object}	map[string]any
//	@Failure	404	{object}	map[string]any
//	@Failure	500	{object}	map[string]any
//	@Router		/admin/complexes/detail [post]
func (h *AdminHandler) PostAdminGetComplexDetail(c *gin.Context) {
	ctx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)

	var req dto.AdminGetComplexDetailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httperrors.SendHTTPErrorObj(c, httperrors.ConvertBindError(err))
		return
	}

	kind, err := parseCoverageKind(req.CoverageType)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	input := propertycoverageservice.GetComplexDetailInput{ID: req.ID, Kind: kind}

	complexEntity, svcErr := h.propertyCoverageService.GetComplexDetail(ctx, input)
	if svcErr != nil {
		httperrors.SendHTTPErrorObj(c, svcErr)
		return
	}

	c.JSON(http.StatusOK, httpconv.ToComplexResponse(complexEntity))
}
