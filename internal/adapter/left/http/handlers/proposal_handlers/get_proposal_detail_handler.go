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

// GetProposalDetail returns proposal metadata plus documents (base64 encoded) to the owner or the realtor.
// @Summary     Retrieve proposal detail with attachments
// @Description Provides proposal metadata (including createdAt/receivedAt/respondedAt), realtor profile summary with acceptedProposals/photoUrl and the binary (base64) of any PDF attachments for authorized actors.
// @Tags        Proposals
// @Security    BearerAuth
// @Accept      json
// @Produce     json
// @Param       request body dto.GetProposalDetailRequest true "Proposal identifier"
// @Success     200 {object} dto.ProposalDetailResponse "Includes documents[].base64Payload, realtor.photoUrl and proposal timeline"
// @Failure     400,401,403,404,500 {object} dto.ErrorResponse
// @Router      /proposals/detail [post]
func (h *ProposalHandler) GetProposalDetail(c *gin.Context) {
	baseCtx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)

	actor, err := converters.ProposalActorFromContext(c)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	var request dto.GetProposalDetailRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httperrors.SendHTTPErrorObj(c, httputils.MapBindingError(err))
		return
	}

	ctx := coreutils.ContextWithLogger(baseCtx)
	detail, svcErr := h.proposalService.GetProposalDetail(ctx, proposalservice.DetailInput{
		ProposalID: request.ProposalID,
		Actor:      actor,
	})
	if svcErr != nil {
		httperrors.SendHTTPErrorObj(c, svcErr)
		return
	}

	c.JSON(http.StatusOK, converters.ProposalDetailToResponse(detail))
}
