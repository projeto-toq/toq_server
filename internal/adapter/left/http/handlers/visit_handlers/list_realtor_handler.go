package visithandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/converters"
	dto "github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

// ListVisitsRealtor handles GET /visits/realtor.
//
// @Summary     List visits for requesters
// @Description Lists visits for the authenticated requester/realtor with status/type/time filters (RFC3339) and pagination (max 50 per page).
// @Tags        Visits
// @Produce     json
// @Security    BearerAuth
// @Param       listingIdentityId query int false "Listing identity filter" example(123)
// @Param       status            query []string false "Statuses" collectionFormat(multi)
// @Param       type              query []string false "Visit types" collectionFormat(multi)
// @Param       from              query string false "Start date/time (RFC3339)" example("2025-01-01T00:00:00Z")
// @Param       to                query string false "End date/time (RFC3339)" example("2025-01-31T23:59:59Z")
// @Param       page              query int false "Page" default(1)
// @Param       limit             query int false "Page size" default(20)
// @Success     200 {object} dto.VisitListResponse
// @Failure     400 {object} dto.ErrorResponse
// @Failure     401 {object} dto.ErrorResponse
// @Failure     403 {object} dto.ErrorResponse
// @Failure     404 {object} dto.ErrorResponse
// @Failure     500 {object} dto.ErrorResponse
// @Router      /visits/realtor [get]
func (h *VisitHandler) ListVisitsRealtor(c *gin.Context) {
	baseCtx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)

	var query dto.VisitListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		httperrors.SendHTTPErrorObj(c, httperrors.ConvertBindError(err))
		return
	}

	userInfo, infoErr := coreutils.GetUserInfoFromGinContext(c)
	if infoErr != nil {
		httperrors.SendHTTPErrorObj(c, infoErr)
		return
	}

	filter, err := buildVisitFilter(query, userInfo.ID, false)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	ctx := coreutils.ContextWithLogger(baseCtx)
	result, svcErr := h.visitService.ListVisits(ctx, filter)
	if svcErr != nil {
		httperrors.SendHTTPErrorObj(c, svcErr)
		return
	}

	c.JSON(http.StatusOK, converters.VisitListDetailToResponse(result))
}
