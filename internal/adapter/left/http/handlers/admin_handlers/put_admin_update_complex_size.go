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

// PutAdminUpdateComplexSize handles PUT /admin/complexes/sizes
//
//	@Summary	Update a complex size
//	@Tags		Admin Complexes
//	@Accept		json
//	@Produce	json
//	@Param		request	body	dto.AdminUpdateComplexSizeRequest	true	"Complex size payload"
//	@Success	200	{object}	dto.ComplexSizeResponse
//	@Failure	400	{object}	map[string]any
//	@Failure	401	{object}	map[string]any
//	@Failure	403	{object}	map[string]any
//	@Failure	404	{object}	map[string]any
//	@Failure	500	{object}	map[string]any
//	@Router		/admin/complexes/sizes [put]
func (h *AdminHandler) PutAdminUpdateComplexSize(c *gin.Context) {
	ctx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)

	var req dto.AdminUpdateComplexSizeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httperrors.SendHTTPErrorObj(c, httperrors.ConvertBindError(err))
		return
	}

	input := propertycoverageservice.UpdateComplexSizeInput{
		ID: req.ID,
		CreateComplexSizeInput: propertycoverageservice.CreateComplexSizeInput{
			VerticalComplexID: req.ComplexID,
			Size:              req.Size,
			Description:       req.Description,
		},
	}

	size, err := h.propertyCoverageService.UpdateComplexSize(ctx, input)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	c.JSON(http.StatusOK, httpconv.ToComplexSizeResponse(size))
}
