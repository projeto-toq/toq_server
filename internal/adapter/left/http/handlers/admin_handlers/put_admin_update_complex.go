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

	kind, err := parseCoverageKind(req.CoverageType)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	input := propertycoverageservice.UpdateComplexInput{
		ID: req.ID,
		CreateComplexInput: propertycoverageservice.CreateComplexInput{
			Kind:             kind,
			Name:             req.Name,
			ZipCode:          req.ZipCode,
			Street:           req.Street,
			Number:           req.Number,
			Neighborhood:     req.Neighborhood,
			City:             req.City,
			State:            req.State,
			ReceptionPhone:   req.PhoneNumber,
			Sector:           propertycoveragemodel.Sector(*req.Sector),
			MainRegistration: req.MainRegistration,
			PropertyType:     globalmodel.PropertyType(*req.PropertyType),
		},
	}

	complexEntity, svcErr := h.propertyCoverageService.UpdateComplex(ctx, input)
	if svcErr != nil {
		httperrors.SendHTTPErrorObj(c, svcErr)
		return
	}

	c.JSON(http.StatusOK, httpconv.ToComplexResponse(complexEntity))
}
