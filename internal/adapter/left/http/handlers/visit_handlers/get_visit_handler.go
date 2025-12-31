package visithandlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/converters"
	dto "github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

// swagger:route helper to keep dto import for specs
var _ dto.VisitResponse

// GetVisit handles GET /visits/{id}.
//
// @Summary     Get visit by ID
// @Description Retrieves a visit by its identifier.
// @Tags        Visits
// @Produce     json
// @Security    BearerAuth
// @Param       visitId path int true "Visit ID"
// @Success     200 {object} dto.VisitResponse
// @Failure     400 {object} dto.ErrorResponse
// @Failure     401 {object} dto.ErrorResponse
// @Failure     403 {object} dto.ErrorResponse
// @Failure     404 {object} dto.ErrorResponse
// @Failure     500 {object} dto.ErrorResponse
// @Router      /visits/{visitId} [get]
func (h *VisitHandler) GetVisit(c *gin.Context) {
	baseCtx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)

	idParam := c.Param("visitId")
	visitID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil || visitID <= 0 {
		httperrors.SendHTTPErrorObj(c, coreutils.ValidationError("visitId", "must be a positive integer"))
		return
	}

	ctx := coreutils.ContextWithLogger(baseCtx)
	visit, svcErr := h.visitService.GetVisit(ctx, visitID)
	if svcErr != nil {
		httperrors.SendHTTPErrorObj(c, svcErr)
		return
	}

	c.JSON(http.StatusOK, converters.VisitDomainToResponse(visit))
}
