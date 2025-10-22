package adminhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	httpconv "github.com/projeto-toq/toq_server/internal/adapter/left/http/converters"
	dto "github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	complexmodel "github.com/projeto-toq/toq_server/internal/core/model/complex_model"
	globalmodel "github.com/projeto-toq/toq_server/internal/core/model/global_model"
	complexservices "github.com/projeto-toq/toq_server/internal/core/service/complex_service"
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
//	@Param		state	query	string	false	"Complex state"
//	@Param		sector	query	int	false	"Sector identifier"
//	@Param		propertyType	query	int	false	"Property type identifier"
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

	var sector *complexmodel.Sector
	if req.Sector != nil {
		converted := complexmodel.Sector(*req.Sector)
		sector = &converted
	}

	var propertyType *globalmodel.PropertyType
	if req.PropertyType != nil {
		converted := globalmodel.PropertyType(*req.PropertyType)
		propertyType = &converted
	}

	input := complexservices.ListComplexesInput{
		Name:         req.Name,
		ZipCode:      req.ZipCode,
		City:         req.City,
		State:        req.State,
		Sector:       sector,
		PropertyType: propertyType,
		Page:         page,
		Limit:        limit,
	}

	complexes, err := h.complexService.ListComplexes(ctx, input)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	response := dto.AdminListComplexesResponse{
		Complexes: httpconv.ToComplexResponses(complexes),
		Page:      page,
		Limit:     limit,
	}

	c.JSON(http.StatusOK, response)
}
