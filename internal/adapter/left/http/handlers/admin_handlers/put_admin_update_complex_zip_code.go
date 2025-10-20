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

// PutAdminUpdateComplexZipCode handles PUT /admin/complexes/zip-codes
//
//	@Summary	Update a complex zip code
//	@Tags		Admin
//	@Accept		json
//	@Produce	json
//	@Param		request	body	dto.AdminUpdateComplexZipCodeRequest	true	"Complex zip code payload"
//	@Success	200	{object}	dto.ComplexZipCodeResponse
//	@Failure	400	{object}	map[string]any
//	@Failure	401	{object}	map[string]any
//	@Failure	403	{object}	map[string]any
//	@Failure	404	{object}	map[string]any
//	@Failure	500	{object}	map[string]any
//	@Router		/admin/complexes/zip-codes [put]
func (h *AdminHandler) PutAdminUpdateComplexZipCode(c *gin.Context) {
	ctx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)

	var req dto.AdminUpdateComplexZipCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httperrors.SendHTTPErrorObj(c, httperrors.ConvertBindError(err))
		return
	}

	input := complexservices.UpdateComplexZipCodeInput{
		ID:        req.ID,
		ComplexID: req.ComplexID,
		ZipCode:   req.ZipCode,
	}

	zipCode, err := h.complexService.UpdateComplexZipCode(ctx, input)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	c.JSON(http.StatusOK, httpconv.ToComplexZipCodeResponse(zipCode))
}
