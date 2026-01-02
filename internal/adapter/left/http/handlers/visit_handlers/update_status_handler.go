package visithandlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/converters"
	dto "github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

// UpdateVisitStatus handles POST /visits/status for status transitions.
//
// @Summary     Update visit status
// @Description Approve, reject, cancel, complete, or mark a visit as no-show.
// @Tags        Visits
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       request body dto.UpdateVisitStatusRequest true "Status payload"
// @Success     200 {object} dto.VisitResponse
// @Failure     400 {object} dto.ErrorResponse
// @Failure     401 {object} dto.ErrorResponse
// @Failure     403 {object} dto.ErrorResponse
// @Failure     404 {object} dto.ErrorResponse
// @Failure     409 {object} dto.ErrorResponse
// @Failure     500 {object} dto.ErrorResponse
// @Router      /visits/status [post]
func (h *VisitHandler) UpdateVisitStatus(c *gin.Context) {
	baseCtx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)

	var req dto.UpdateVisitStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httperrors.SendHTTPErrorObj(c, httperrors.ConvertBindError(err))
		return
	}

	action := strings.ToUpper(strings.TrimSpace(req.Action))
	ctx := coreutils.ContextWithLogger(baseCtx)

	var visit listingmodel.VisitInterface
	var svcErr error

	switch action {
	case "APPROVE":
		visit, svcErr = h.visitService.ApproveVisit(ctx, req.VisitID, strings.TrimSpace(req.Notes))
	case "REJECT":
		reason := strings.TrimSpace(req.RejectionReason)
		if reason == "" {
			httperrors.SendHTTPErrorObj(c, coreutils.ValidationError("rejectionReason", "is required for REJECT"))
			return
		}
		visit, svcErr = h.visitService.RejectVisit(ctx, req.VisitID, reason)
	case "CANCEL":
		reason := strings.TrimSpace(req.RejectionReason)
		if reason == "" {
			httperrors.SendHTTPErrorObj(c, coreutils.ValidationError("rejectionReason", "is required for CANCEL"))
			return
		}
		visit, svcErr = h.visitService.CancelVisit(ctx, req.VisitID, reason)
	case "COMPLETE":
		visit, svcErr = h.visitService.CompleteVisit(ctx, req.VisitID, strings.TrimSpace(req.Notes))
	case "NO_SHOW":
		visit, svcErr = h.visitService.MarkNoShow(ctx, req.VisitID, strings.TrimSpace(req.Notes))
	default:
		httperrors.SendHTTPErrorObj(c, coreutils.ValidationError("action", "must be one of APPROVE, REJECT, CANCEL, COMPLETE, NO_SHOW"))
		return
	}

	if svcErr != nil {
		httperrors.SendHTTPErrorObj(c, svcErr)
		return
	}

	c.JSON(http.StatusOK, converters.VisitDomainToResponse(visit))
}
