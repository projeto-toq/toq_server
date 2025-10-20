package complexhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	httpconv "github.com/projeto-toq/toq_server/internal/adapter/left/http/converters"
	dto "github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	complexservices "github.com/projeto-toq/toq_server/internal/core/service/complex_service"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

// ListSizesByAddress handles GET /complex/sizes
//
//	@Summary	List complex sizes by address
//	@Tags		Complex
//	@Accept		json
//	@Produce	json
//	@Param		zipCode	query	string	true	"Complex zip code"
//	@Param		number	query	string	false	"Complex street number"
//	@Success	200	{object}	dto.ListSizesByAddressResponse
//	@Failure	400	{object}	map[string]any
//	@Failure	404	{object}	map[string]any
//	@Failure	500	{object}	map[string]any
//	@Router		/complex/sizes [get]
func (h *ComplexHandler) ListSizesByAddress(c *gin.Context) {
	ctx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)

	var query dto.ListSizesByAddressQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		httperrors.SendHTTPErrorObj(c, httperrors.ConvertBindError(err))
		return
	}

	input := complexservices.ListSizesByAddressInput{
		ZipCode: query.ZipCode,
		Number:  query.Number,
	}

	sizes, err := h.complexService.ListSizesByAddress(ctx, input)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	response := dto.ListSizesByAddressResponse{
		Sizes: make([]dto.ComplexSizeResponse, 0, len(sizes)),
	}

	for _, size := range sizes {
		response.Sizes = append(response.Sizes, httpconv.ToComplexSizeResponse(size))
	}

	c.JSON(http.StatusOK, response)
}
