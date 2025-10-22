package schedulehandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/converters"
	dto "github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/middlewares"
	schedulemodel "github.com/projeto-toq/toq_server/internal/core/model/schedule_model"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

// GetOwnerSummary handles GET /schedules/owner/summary.
//
// @Summary	List owner agenda summary
// @Description	Returns a consolidated view of agenda entries for all listings owned by the authenticated user.
// @Tags	Schedules
// @Produce	json
// @Param	listingIds	query	[]int64	false	"Listing identifiers" collectionFormat(multi)
// @Param	rangeFrom	query	string	false	"Start of time range (RFC3339)"
// @Param	rangeTo	query	string	false	"End of time range (RFC3339)"
// @Param	page	query	int	false	"Page number"
// @Param	limit	query	int	false	"Items per page"
// @Success	200	{object}	dto.OwnerAgendaSummaryResponse
// @Failure	400	{object}	dto.ErrorResponse
// @Failure	401	{object}	dto.ErrorResponse
// @Failure	403	{object}	dto.ErrorResponse
// @Failure	500	{object}	dto.ErrorResponse
// @Router	/schedules/owner/summary [get]
// @Security	BearerAuth
func (h *ScheduleHandler) GetOwnerSummary(c *gin.Context) {
	baseCtx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)
	ctx, spanEnd, err := coreutils.GenerateTracer(baseCtx)
	if err != nil {
		httperrors.SendHTTPError(c, http.StatusInternalServerError, "TRACER_ERROR", "Failed to generate tracer")
		return
	}
	defer spanEnd()

	userInfo, ok := middlewares.GetUserInfoFromContext(c)
	if !ok {
		httperrors.SendHTTPError(c, http.StatusUnauthorized, "UNAUTHENTICATED", "User context not found")
		return
	}

	var req dto.OwnerAgendaSummaryQuery
	if err := c.ShouldBindQuery(&req); err != nil {
		httperrors.SendHTTPError(c, http.StatusBadRequest, "INVALID_QUERY", "Invalid query parameters")
		return
	}

	rangeFilter, err := parseScheduleRange(dto.ScheduleRangeRequest{From: req.RangeFrom, To: req.RangeTo})
	if err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	pagination := sanitizeSchedulePagination(dto.SchedulePaginationRequest{Page: req.Page, Limit: req.Limit})

	filter := schedulemodel.OwnerSummaryFilter{
		OwnerID:    userInfo.ID,
		ListingIDs: req.ListingIDs,
		Range:      rangeFilter,
		Pagination: pagination,
	}

	ctx = coreutils.ContextWithLogger(ctx)
	result, serviceErr := h.scheduleService.ListOwnerSummary(ctx, filter)
	if serviceErr != nil {
		httperrors.SendHTTPErrorObj(c, serviceErr)
		return
	}

	page, limit := schedulePaginationValues(pagination)
	response := converters.ScheduleOwnerSummaryToDTO(result, page, limit)

	c.JSON(http.StatusOK, response)
}
