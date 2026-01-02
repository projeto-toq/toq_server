package visithandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/converters"
	dto "github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

// CreateVisit handles POST /visits.
//
// @Summary     Request a visit
// @Description Validates lead time (visits.min_hours_ahead), horizon (visits.max_days_ahead), and listing availability before creating a visit.
// @Tags        Visits
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       Authorization header string true "Bearer token"
// @Param       X-Device-Id   header string false "Device ID"
// @Param       request body dto.CreateVisitRequest true "Visit payload"
// @Success     201 {object} dto.VisitResponse
// @Failure     400 {object} dto.ErrorResponse
// @Failure     401 {object} dto.ErrorResponse
// @Failure     403 {object} dto.ErrorResponse
// @Failure     404 {object} dto.ErrorResponse
// @Failure     409 {object} dto.ErrorResponse
// @Failure     500 {object} dto.ErrorResponse
// @Router      /visits [post]
func (h *VisitHandler) CreateVisit(c *gin.Context) {
	baseCtx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)

	var req dto.CreateVisitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httperrors.SendHTTPErrorObj(c, httperrors.ConvertBindError(err))
		return
	}

	input, err := converters.CreateVisitDTOToInput(req)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	ctx := coreutils.ContextWithLogger(baseCtx)
	visit, svcErr := h.visitService.CreateVisit(ctx, input)
	if svcErr != nil {
		httperrors.SendHTTPErrorObj(c, svcErr)
		return
	}

	c.JSON(http.StatusCreated, converters.VisitDomainToResponse(visit))
}
