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

// ListOwnerProposals lists proposals received by the owner.
// @Summary     List owner proposals
// @Description Owners can see all proposals received across listings with the same filters as realtors, scoped to identities they own.
// @Tags        Proposals
// @Security    BearerAuth
// @Produce     json
// @Param       status query []string false "Status filter"
// @Param       listingIdentityId query int false "Listing identity"
// @Param       page query int false "Page number"
// @Param       pageSize query int false "Page size"
// @Success     200 {object} dto.ListProposalsResponse
// @Failure     400,401,500 {object} dto.ErrorResponse
// @Router      /proposals/owner [get]
func (h *ProposalHandler) ListOwnerProposals(c *gin.Context) {
	baseCtx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)

	actor, err := converters.ProposalActorFromContext(c)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	var query dto.ListProposalsQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		httperrors.SendHTTPErrorObj(c, httputils.MapBindingError(err))
		return
	}

	filter, err := converters.ProposalListFilterFromQuery(query, actor)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	ctx := coreutils.ContextWithLogger(baseCtx)
	result, svcErr := h.proposalService.ListOwnerProposals(ctx, filter)
	if svcErr != nil {
		httperrors.SendHTTPErrorObj(c, svcErr)
		return
	}

	c.JSON(http.StatusOK, converters.ProposalListToResponse(result))
}
