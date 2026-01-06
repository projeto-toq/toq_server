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

// ListRealtorProposals returns paginated history filtered by realtor context.
// @Summary     List realtor proposals
// @Description Returns paginated proposals created by the authenticated realtor, supporting filters by status and listing identity.
// @Tags        Proposals
// @Security    BearerAuth
// @Produce     json
// @Param       status query []string false "Status filter" collectionFormat(multi)
// @Param       listingIdentityId query int false "Listing identity"
// @Param       page query int false "Page number" default(1)
// @Param       pageSize query int false "Page size" default(20)
// @Success     200 {object} dto.ListProposalsResponse
// @Failure     400,401,500 {object} dto.ErrorResponse
// @Router      /proposals/realtor [get]
func (h *ProposalHandler) ListRealtorProposals(c *gin.Context) {
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
	result, svcErr := h.proposalService.ListRealtorProposals(ctx, filter)
	if svcErr != nil {
		httperrors.SendHTTPErrorObj(c, svcErr)
		return
	}

	c.JSON(http.StatusOK, converters.ProposalListToResponse(result))
}
