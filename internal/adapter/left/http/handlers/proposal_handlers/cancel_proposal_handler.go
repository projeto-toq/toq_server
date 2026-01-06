package proposalhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/converters"
	dto "github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	httputils "github.com/projeto-toq/toq_server/internal/adapter/left/http/utils"
	proposalservice "github.com/projeto-toq/toq_server/internal/core/service/proposal_service"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

// CancelProposal lets the realtor withdraw a pending proposal.
//
// @Summary     Cancel a proposal before owner decision
// @Description Realtor-owned proposals can be cancelled anytime before acceptance, generating notifications to the owner.
// @Tags        Proposals
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       request body dto.CancelProposalRequest true "Proposal identifier"
// @Success     204 {string} string "No Content"
// @Failure     400,401,403,409,500 {object} dto.ErrorResponse
// @Router      /proposals/cancel [post]
func (h *ProposalHandler) CancelProposal(c *gin.Context) {
	baseCtx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)

	actor, err := converters.ProposalActorFromContext(c)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	var request dto.CancelProposalRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httperrors.SendHTTPErrorObj(c, httputils.MapBindingError(err))
		return
	}

	ctx := coreutils.ContextWithLogger(baseCtx)
	if svcErr := h.proposalService.CancelProposal(ctx, proposalservice.StatusChangeInput{
		ProposalID: request.ProposalID,
		Actor:      actor,
	}); svcErr != nil {
		httperrors.SendHTTPErrorObj(c, svcErr)
		return
	}

	c.Status(http.StatusNoContent)
}
