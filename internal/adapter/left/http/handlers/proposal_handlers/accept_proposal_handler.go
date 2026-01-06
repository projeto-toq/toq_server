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

// AcceptProposal confirms the owner's approval and updates listing flags.
// @Summary     Accept a proposal
// @Description Owners can accept a single pending proposal per listing, triggering notifications and flag updates.
// @Tags        Proposals
// @Security    BearerAuth
// @Accept      json
// @Produce     json
// @Param       request body dto.AcceptProposalRequest true "Proposal identifier"
// @Success     200 {object} dto.ProposalResponse
// @Failure     400,401,403,409,500 {object} dto.ErrorResponse
// @Router      /proposals/accept [post]
func (h *ProposalHandler) AcceptProposal(c *gin.Context) {
	baseCtx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)

	actor, err := converters.ProposalActorFromContext(c)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	var request dto.AcceptProposalRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httperrors.SendHTTPErrorObj(c, httputils.MapBindingError(err))
		return
	}

	ctx := coreutils.ContextWithLogger(baseCtx)
	result, svcErr := h.proposalService.AcceptProposal(ctx, proposalservice.StatusChangeInput{
		ProposalID: request.ProposalID,
		Actor:      actor,
	})
	if svcErr != nil {
		httperrors.SendHTTPErrorObj(c, svcErr)
		return
	}

	c.JSON(http.StatusOK, converters.ProposalDomainToResponse(result))
}
