package adminhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	httpconv "github.com/projeto-toq/toq_server/internal/adapter/left/http/converters"
	dto "github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	globalmodel "github.com/projeto-toq/toq_server/internal/core/model/global_model"
	propertycoveragemodel "github.com/projeto-toq/toq_server/internal/core/model/property_coverage_model"
	propertycoverageservice "github.com/projeto-toq/toq_server/internal/core/service/property_coverage_service"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

// GetAdminComplexes handles GET /admin/complexes
//
//	@Summary	List complexes
//	@Tags		Admin Complexes
//	@Produce	json
//	@Param		name	query	string	false	"Complex name filter"
//	@Param		zipCode	query	string	false	"Complex zip code"
//	@Param		city	query	string	false	"Complex city"
//	@Param		number	query	string	false	"Complex number"
//	@Param		state	query	string	false	"Complex state"
//	@Param		sector	query	int	false	"Sector identifier"
//	@Param		propertyType	query	int	false	"Property type identifier"
//	@Param		coverageType	query	string	false	"Coverage type (VERTICAL, HORIZONTAL, STANDALONE)"
//	@Param		page	query	int	false	"Page number"
//	@Param		limit	query	int	false	"Page size"
//	@Success	200	{object}	dto.AdminListComplexesResponse
//	@Failure	400	{object}	map[string]any
//	@Failure	401	{object}	map[string]any
//	@Failure	403	{object}	map[string]any
//	@Failure	500	{object}	map[string]any
//	@Router		/admin/complexes [get]
func (h *AdminHandler) GetAdminComplexes(c *gin.Context) {
	ctx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)

	var req dto.AdminListComplexesRequest
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

	var sector *propertycoveragemodel.Sector
	if req.Sector != nil {
		converted := propertycoveragemodel.Sector(*req.Sector)
		sector = &converted
	}

	var propertyType *globalmodel.PropertyType
	if req.PropertyType != nil {
		converted := globalmodel.PropertyType(*req.PropertyType)
		propertyType = &converted
	}

	kind, err := parseOptionalCoverageKind(req.CoverageType)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	input := propertycoverageservice.ListComplexesInput{
		Name:         req.Name,
		ZipCode:      req.ZipCode,
		Number:       req.Number,
		City:         req.City,
		State:        req.State,
		Sector:       sector,
		PropertyType: propertyType,
		Kind:         kind,
		Page:         page,
		Limit:        limit,
	}

	complexes, svcErr := h.propertyCoverageService.ListComplexes(ctx, input)
	if svcErr != nil {
		httperrors.SendHTTPErrorObj(c, svcErr)
		return
	}

	response := dto.AdminListComplexesResponse{
		Complexes: httpconv.ToComplexResponses(complexes),
		Page:      page,
		Limit:     limit,
	}

	c.JSON(http.StatusOK, response)
}
