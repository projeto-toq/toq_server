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

// GetAdminComplexTowers handles GET /admin/complexes/towers
//
//	@Summary	List complex towers
//	@Tags		Admin Complexes
//	@Produce	json
//	@Param		complexId	query	int	false	"Complex identifier"
//	@Param		tower	query	string	false	"Tower name filter"
//	@Param		page	query	int	false	"Page number"
//	@Param		limit	query	int	false	"Page size"
//	@Success	200	{object}	dto.AdminListComplexTowersResponse
//	@Failure	400	{object}	map[string]any
//	@Failure	401	{object}	map[string]any
//	@Failure	403	{object}	map[string]any
//	@Failure	500	{object}	map[string]any
//	@Router		/admin/complexes/towers [get]
func (h *AdminHandler) GetAdminComplexTowers(c *gin.Context) {
	ctx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)

	var req dto.AdminListComplexTowersRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		httperrors.SendHTTPErrorObj(c, httperrors.ConvertBindError(err))
		return
	}

	page := req.Page
	if page <= 0 {
		page = 1
	}

	limit := req.Limit
	if limit <= 0 {
		limit = 20
	}

	input := complexservices.ListComplexTowersInput{
		ComplexID: req.ComplexID,
		Tower:     req.Tower,
		Page:      page,
		Limit:     limit,
	}

	towers, err := h.complexService.ListComplexTowers(ctx, input)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	response := dto.AdminListComplexTowersResponse{
		Towers: make([]dto.ComplexTowerResponse, 0, len(towers)),
		Page:   page,
		Limit:  limit,
	}

	for _, tower := range towers {
		response.Towers = append(response.Towers, httpconv.ToComplexTowerResponse(tower))
	}

	c.JSON(http.StatusOK, response)
}
