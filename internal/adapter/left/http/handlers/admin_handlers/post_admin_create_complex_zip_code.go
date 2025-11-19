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

// PostAdminCreateComplexZipCode handles POST /admin/complexes/zip-codes
//
//	@Summary	Create a complex zip code
//	@Tags		Admin Complexes
//	@Accept		json
//	@Produce	json
//	@Param		request	body	dto.AdminCreateComplexZipCodeRequest	true	"Complex zip code payload"
//	@Success	201	{object}	dto.ComplexZipCodeResponse
//	@Failure	400	{object}	map[string]any
//	@Failure	401	{object}	map[string]any
//	@Failure	403	{object}	map[string]any
//	@Failure	404	{object}	map[string]any
//	@Failure	500	{object}	map[string]any
//	@Router		/admin/complexes/zip-codes [post]
func (h *AdminHandler) PostAdminCreateComplexZipCode(c *gin.Context) {
	ctx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)

	var req dto.AdminCreateComplexZipCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httperrors.SendHTTPErrorObj(c, httperrors.ConvertBindError(err))
		return
	}

	input := propertycoverageservice.CreateComplexZipCodeInput{
		HorizontalComplexID: req.ComplexID,
		ZipCode:             req.ZipCode,
	}

	zipCode, err := h.propertyCoverageService.CreateComplexZipCode(ctx, input)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	c.JSON(http.StatusCreated, httpconv.ToComplexZipCodeResponse(zipCode))
}
