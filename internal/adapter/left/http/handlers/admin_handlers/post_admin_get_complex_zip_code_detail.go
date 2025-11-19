package adminhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	httpconv "github.com/projeto-toq/toq_server/internal/adapter/left/http/converters"
	dto "github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

// PostAdminGetComplexZipCodeDetail handles POST /admin/complexes/zip-codes/detail
//
//	@Summary	Get complex zip code detail
//	@Tags		Admin Complexes
//	@Accept		json
//	@Produce	json
//	@Param		request	body	dto.AdminGetComplexZipCodeDetailRequest	true	"Complex zip code detail payload"
//	@Success	200	{object}	dto.ComplexZipCodeResponse
//	@Failure	400	{object}	map[string]any
//	@Failure	401	{object}	map[string]any
//	@Failure	403	{object}	map[string]any
//	@Failure	404	{object}	map[string]any
//	@Failure	500	{object}	map[string]any
//	@Router		/admin/complexes/zip-codes/detail [post]
func (h *AdminHandler) PostAdminGetComplexZipCodeDetail(c *gin.Context) {
	ctx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)

	var req dto.AdminGetComplexZipCodeDetailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httperrors.SendHTTPErrorObj(c, httperrors.ConvertBindError(err))
		return
	}

	zipCode, err := h.propertyCoverageService.GetComplexZipCodeDetail(ctx, req.ID)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	c.JSON(http.StatusOK, httpconv.ToComplexZipCodeResponse(zipCode))
}
