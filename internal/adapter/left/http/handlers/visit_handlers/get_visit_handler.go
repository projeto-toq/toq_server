package visithandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/converters"
	dto "github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

// swagger:route helper to keep dto import for specs
var _ dto.VisitResponse

// GetVisit handles POST /visits/detail.
//
// @Summary     Get visit detail
// @Description Retrieves a visit by its identifier provided in the request body.
// @Tags        Visits
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       Authorization header string true "Bearer token"
// @Param       request body dto.GetVisitDetailRequest true "Visit identifier"
// @Success     200 {object} dto.VisitResponse
// @Failure     400 {object} dto.ErrorResponse
// @Failure     401 {object} dto.ErrorResponse
// @Failure     403 {object} dto.ErrorResponse
// @Failure     404 {object} dto.ErrorResponse
// @Failure     500 {object} dto.ErrorResponse
// @Router      /visits/detail [post]
func (h *VisitHandler) GetVisit(c *gin.Context) {
	baseCtx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)

	var req dto.GetVisitDetailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httperrors.SendHTTPErrorObj(c, httperrors.ConvertBindError(err))
		return
	}

	ctx := coreutils.ContextWithLogger(baseCtx)
	detail, svcErr := h.visitService.GetVisit(ctx, req.VisitID)
	if svcErr != nil {
		httperrors.SendHTTPErrorObj(c, svcErr)
		return
	}

	c.JSON(http.StatusOK, converters.VisitDetailToResponse(detail))
}
