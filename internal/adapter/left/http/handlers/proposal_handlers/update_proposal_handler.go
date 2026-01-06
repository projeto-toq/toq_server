package proposalhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/converters"
	dto "github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	httputils "github.com/projeto-toq/toq_server/internal/adapter/left/http/utils"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

// UpdateProposal allows the author to edit a pending proposal.
//
// @Summary     Edit a pending proposal
// @Description Updates text and/or PDF while the proposal is still pending, ensuring idempotency and actor ownership.
// @Tags        Proposals
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       Authorization header string true "Bearer <token>"
// @Param       request body dto.UpdateProposalRequest true "Editable fields"
// @Success     200 {object} dto.ProposalResponse
// @Failure     400,401,403,409,422,500 {object} dto.ErrorResponse
// @Router      /proposals [put]
func (h *ProposalHandler) UpdateProposal(c *gin.Context) {
	baseCtx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)

	actor, err := converters.ProposalActorFromContext(c)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	var request dto.UpdateProposalRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httperrors.SendHTTPErrorObj(c, httputils.MapBindingError(err))
		return
	}

	input, err := converters.UpdateProposalDTOToInput(request, actor)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	ctx := coreutils.ContextWithLogger(baseCtx)
	result, svcErr := h.proposalService.UpdateProposal(ctx, input)
	if svcErr != nil {
		httperrors.SendHTTPErrorObj(c, svcErr)
		return
	}

	c.JSON(http.StatusOK, converters.ProposalDomainToResponse(result))
}
