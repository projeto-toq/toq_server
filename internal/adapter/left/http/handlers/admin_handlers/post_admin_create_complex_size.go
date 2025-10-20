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

// PostAdminCreateComplexSize handles POST /admin/complexes/sizes
//
//	@Summary	Create a complex size
//	@Tags		Admin
//	@Accept		json
//	@Produce	json
//	@Param		request	body	dto.AdminCreateComplexSizeRequest	true	"Complex size payload"
//	@Success	201	{object}	dto.ComplexSizeResponse
//	@Failure	400	{object}	map[string]any
//	@Failure	401	{object}	map[string]any
//	@Failure	403	{object}	map[string]any
//	@Failure	404	{object}	map[string]any
//	@Failure	500	{object}	map[string]any
//	@Router		/admin/complexes/sizes [post]
func (h *AdminHandler) PostAdminCreateComplexSize(c *gin.Context) {
	ctx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)

	var req dto.AdminCreateComplexSizeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httperrors.SendHTTPErrorObj(c, httperrors.ConvertBindError(err))
		return
	}

	input := complexservices.CreateComplexSizeInput{
		ComplexID:   req.ComplexID,
		Size:        req.Size,
		Description: req.Description,
	}

	size, err := h.complexService.CreateComplexSize(ctx, input)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	c.JSON(http.StatusCreated, httpconv.ToComplexSizeResponse(size))
}
