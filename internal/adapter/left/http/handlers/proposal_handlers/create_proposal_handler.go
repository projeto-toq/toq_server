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

// CreateProposal handles realtor submissions of new proposals.
//
// @Summary     Submit a proposal for a published listing
// @Description Allows authenticated realtors to send a free-text proposal and optional PDF attachment (≤1MB) for a listing identity they do not own.
// @Tags        Proposals
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       Authorization header string true "Bearer <token>"
// @Param       request body dto.CreateProposalRequest true "Proposal payload with optional PDF encoded in base64"
// @Success     201 {object} dto.ProposalResponse
// @Failure     400 {object} dto.ErrorResponse "Invalid payload or PDF too large"
// @Failure     401 {object} dto.ErrorResponse "Authentication required"
// @Failure     409 {object} dto.ErrorResponse "Listing already has pending/accepted proposal"
// @Failure     422 {object} dto.ErrorResponse "Business validation failure"
// @Failure     500 {object} dto.ErrorResponse "Infrastructure failure"
// @Router      /proposals [post]
func (h *ProposalHandler) CreateProposal(c *gin.Context) {
	baseCtx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)

	actor, err := converters.ProposalActorFromContext(c)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	var request dto.CreateProposalRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httperrors.SendHTTPErrorObj(c, httputils.MapBindingError(err))
		return
	}

	input, err := converters.CreateProposalDTOToInput(request, actor)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	ctx := coreutils.ContextWithLogger(baseCtx)
	result, svcErr := h.proposalService.CreateProposal(ctx, input)
	if svcErr != nil {
		httperrors.SendHTTPErrorObj(c, svcErr)
		return
	}

	c.JSON(http.StatusCreated, converters.ProposalDomainToResponse(result))
}
