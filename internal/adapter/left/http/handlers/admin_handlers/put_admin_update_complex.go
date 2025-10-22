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

// PutAdminUpdateComplex handles PUT /admin/complexes
//
//	@Summary	Update a complex
//	@Tags		Admin Complexes
//	@Accept		json
//	@Produce	json
//	@Param		request	body	dto.AdminUpdateComplexRequest	true	"Complex payload"
//	@Success	200	{object}	dto.ComplexResponse
//	@Failure	400	{object}	map[string]any
//	@Failure	401	{object}	map[string]any
//	@Failure	403	{object}	map[string]any
//	@Failure	404	{object}	map[string]any
//	@Failure	500	{object}	map[string]any
//	@Router		/admin/complexes [put]
func (h *AdminHandler) PutAdminUpdateComplex(c *gin.Context) {
	ctx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)

	var req dto.AdminUpdateComplexRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httperrors.SendHTTPErrorObj(c, httperrors.ConvertBindError(err))
		return
	}

	input := complexservices.UpdateComplexInput{
		ID:               req.ID,
		Name:             req.Name,
		ZipCode:          req.ZipCode,
		Street:           req.Street,
		Number:           req.Number,
		Neighborhood:     req.Neighborhood,
		City:             req.City,
		State:            req.State,
		PhoneNumber:      req.PhoneNumber,
		Sector:           complexmodel.Sector(req.Sector),
		MainRegistration: req.MainRegistration,
		PropertyType:     globalmodel.PropertyType(req.PropertyType),
	}

	complexEntity, err := h.complexService.UpdateComplex(ctx, input)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	c.JSON(http.StatusOK, httpconv.ToComplexResponse(complexEntity))
}
