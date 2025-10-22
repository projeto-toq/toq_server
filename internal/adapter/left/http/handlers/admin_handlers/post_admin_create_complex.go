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

// PostAdminCreateComplex handles POST /admin/complexes
//
//	@Summary	Create a complex
//	@Tags		Admin Complexes
//	@Accept		json
//	@Produce	json
//	@Param		request	body	dto.AdminCreateComplexRequest	true	"Complex payload"
//	@Success	201	{object}	dto.ComplexResponse
//	@Failure	400	{object}	map[string]any
//	@Failure	401	{object}	map[string]any
//	@Failure	403	{object}	map[string]any
//	@Failure	404	{object}	map[string]any
//	@Failure	500	{object}	map[string]any
//	@Router		/admin/complexes [post]
func (h *AdminHandler) PostAdminCreateComplex(c *gin.Context) {
	ctx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)

	var req dto.AdminCreateComplexRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httperrors.SendHTTPErrorObj(c, httperrors.ConvertBindError(err))
		return
	}

	input := complexservices.CreateComplexInput{
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

	complexEntity, err := h.complexService.CreateComplex(ctx, input)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	c.JSON(http.StatusCreated, httpconv.ToComplexResponse(complexEntity))
}
