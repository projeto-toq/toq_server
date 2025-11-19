package adminhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	httpconv "github.com/projeto-toq/toq_server/internal/adapter/left/http/converters"
	dto "github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

// PostAdminGetComplexSizeDetail handles POST /admin/complexes/sizes/detail
//
//	@Summary	Get complex size detail
//	@Tags		Admin Complexes
//	@Accept		json
//	@Produce	json
//	@Param		request	body	dto.AdminGetComplexSizeDetailRequest	true	"Complex size detail payload"
//	@Success	200	{object}	dto.ComplexSizeResponse
//	@Failure	400	{object}	map[string]any
//	@Failure	401	{object}	map[string]any
//	@Failure	403	{object}	map[string]any
//	@Failure	404	{object}	map[string]any
//	@Failure	500	{object}	map[string]any
//	@Router		/admin/complexes/sizes/detail [post]
func (h *AdminHandler) PostAdminGetComplexSizeDetail(c *gin.Context) {
	ctx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)

	var req dto.AdminGetComplexSizeDetailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httperrors.SendHTTPErrorObj(c, httperrors.ConvertBindError(err))
		return
	}

	size, err := h.propertyCoverageService.GetComplexSizeDetail(ctx, req.ID)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	c.JSON(http.StatusOK, httpconv.ToComplexSizeResponse(size))
}
