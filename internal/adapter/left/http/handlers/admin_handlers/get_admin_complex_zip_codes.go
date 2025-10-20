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

// GetAdminComplexZipCodes handles GET /admin/complexes/zip-codes
//
//	@Summary	List complex zip codes
//	@Tags		Admin
//	@Produce	json
//	@Param		complexId	query	int	false	"Complex identifier"
//	@Param		zipCode		query	string	false	"Zip code filter"
//	@Param		page		query	int	false	"Page number"
//	@Param		limit		query	int	false	"Page size"
//	@Success	200	{object}	dto.AdminListComplexZipCodesResponse
//	@Failure	400	{object}	map[string]any
//	@Failure	401	{object}	map[string]any
//	@Failure	403	{object}	map[string]any
//	@Failure	500	{object}	map[string]any
//	@Router		/admin/complexes/zip-codes [get]
func (h *AdminHandler) GetAdminComplexZipCodes(c *gin.Context) {
	ctx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)

	var req dto.AdminListComplexZipCodesRequest
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

	input := complexservices.ListComplexZipCodesInput{
		ComplexID: req.ComplexID,
		ZipCode:   req.ZipCode,
		Page:      page,
		Limit:     limit,
	}

	zipCodes, err := h.complexService.ListComplexZipCodes(ctx, input)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	response := dto.AdminListComplexZipCodesResponse{
		ZipCodes: make([]dto.ComplexZipCodeResponse, 0, len(zipCodes)),
		Page:     page,
		Limit:    limit,
	}

	for _, zip := range zipCodes {
		response.ZipCodes = append(response.ZipCodes, httpconv.ToComplexZipCodeResponse(zip))
	}

	c.JSON(http.StatusOK, response)
}
