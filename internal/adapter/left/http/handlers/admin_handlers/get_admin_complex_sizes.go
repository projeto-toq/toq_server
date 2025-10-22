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

// GetAdminComplexSizes handles GET /admin/complexes/sizes
//
//	@Summary	List complex sizes
//	@Tags		Admin Complexes
//	@Produce	json
//	@Param		complexId	query	int	false	"Complex identifier"
//	@Param		page		query	int	false	"Page number"
//	@Param		limit	query	int	false	"Page size"
//	@Success	200	{object}	dto.AdminListComplexSizesResponse
//	@Failure	400	{object}	map[string]any
//	@Failure	401	{object}	map[string]any
//	@Failure	403	{object}	map[string]any
//	@Failure	500	{object}	map[string]any
//	@Router		/admin/complexes/sizes [get]
func (h *AdminHandler) GetAdminComplexSizes(c *gin.Context) {
	ctx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)

	var req dto.AdminListComplexSizesRequest
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

	input := complexservices.ListComplexSizesInput{
		ComplexID: req.ComplexID,
		Page:      page,
		Limit:     limit,
	}

	sizes, err := h.complexService.ListComplexSizes(ctx, input)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	response := dto.AdminListComplexSizesResponse{
		Sizes: make([]dto.ComplexSizeResponse, 0, len(sizes)),
		Page:  page,
		Limit: limit,
	}

	for _, size := range sizes {
		response.Sizes = append(response.Sizes, httpconv.ToComplexSizeResponse(size))
	}

	c.JSON(http.StatusOK, response)
}
