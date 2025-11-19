package complexhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	httpconv "github.com/projeto-toq/toq_server/internal/adapter/left/http/converters"
	dto "github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	propertycoverageservice "github.com/projeto-toq/toq_server/internal/core/service/property_coverage_service"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

// GetComplexByAddress handles GET /complex
//
//	@Summary	Get complex details by address
//	@Description Retrieves full complex details (including towers, sizes, and zip codes) based on ZipCode and Number.
//	@Tags		Complex
//	@Accept		json
//	@Produce	json
//	@Param		zipCode	query	string	true	"Complex zip code"
//	@Param		number	query	string	false	"Complex street number"
//	@Success	200	{object}	dto.ComplexResponse
//	@Failure	400	{object}	dto.ErrorResponse
//	@Failure	404	{object}	dto.ErrorResponse
//	@Failure	500	{object}	dto.ErrorResponse
//	@Router		/complex [get]
func (h *ComplexHandler) GetComplexByAddress(c *gin.Context) {
	ctx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)

	var query dto.GetComplexByAddressQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		httperrors.SendHTTPErrorObj(c, httperrors.ConvertBindError(err))
		return
	}

	input := propertycoverageservice.GetComplexByAddressInput{
		ZipCode: query.ZipCode,
		Number:  query.Number,
	}

	complex, err := h.propertyCoverageService.GetComplexByAddress(ctx, input)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	c.JSON(http.StatusOK, httpconv.ToComplexResponse(complex))
}
