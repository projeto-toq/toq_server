package proposalhandlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/converters"
	dto "github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	httputils "github.com/projeto-toq/toq_server/internal/adapter/left/http/utils"
	proposalservice "github.com/projeto-toq/toq_server/internal/core/service/proposal_service"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

// RejectProposal stores the owner's reason and informs the realtor.
// @Summary     Reject a proposal with reason
// @Description Owners must provide a free text reason; realtor receives a push notification with the status update.
// @Tags        Proposals
// @Security    BearerAuth
// @Accept      json
// @Produce     json
// @Param       request body dto.RejectProposalRequest true "Proposal identifier and reason"
// @Success     200 {object} dto.ProposalResponse
// @Failure     400,401,403,409,422,500 {object} dto.ErrorResponse
// @Router      /proposals/reject [post]
func (h *ProposalHandler) RejectProposal(c *gin.Context) {
	baseCtx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)

	actor, err := converters.ProposalActorFromContext(c)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	var request dto.RejectProposalRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httperrors.SendHTTPErrorObj(c, httputils.MapBindingError(err))
		return
	}

	ctx := coreutils.ContextWithLogger(baseCtx)
	result, svcErr := h.proposalService.RejectProposal(ctx, proposalservice.StatusChangeInput{
		ProposalID: request.ProposalID,
		Actor:      actor,
		Reason:     strings.TrimSpace(request.Reason),
	})
	if svcErr != nil {
		httperrors.SendHTTPErrorObj(c, svcErr)
		return
	}

	c.JSON(http.StatusOK, converters.ProposalDomainToResponse(result))
}
