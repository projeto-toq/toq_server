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

// GetListingBlockEntries handles GET /schedules/listing/block.
//
// @Summary	List blocking entries for a listing agenda
// @Description	Returns blocking agenda entries for a listing owned by the authenticated user.
// @Tags		Listing Schedules
// @Produce	json
// @Param	listingId	query	int64	true	"Listing identifier" example(3241)
// @Param	rangeFrom	query	string	false	"Start of time range (RFC3339)" example(2025-01-01T00:00:00Z)
// @Param	rangeTo	query	string	false	"End of time range (RFC3339)" example(2025-01-07T23:59:59Z)
// @Param	page	query	int	false	"Page number" example(1)
// @Param	limit	query	int	false	"Items per page" example(20)
// @Success	200	{object}	dto.ListingBlockEntriesResponse
// @Failure	400	{object}	dto.ErrorResponse
// @Failure	401	{object}	dto.ErrorResponse
// @Failure	403	{object}	dto.ErrorResponse
// @Failure	404	{object}	dto.ErrorResponse
// @Failure	500	{object}	dto.ErrorResponse
// @Router	/schedules/listing/block [get]
// @Security	BearerAuth
func (h *ScheduleHandler) GetListingBlockEntries(c *gin.Context) {
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

	var req dto.ListingBlockEntriesQuery
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

	filter := schedulemodel.BlockEntriesFilter{
		OwnerID:    userInfo.ID,
		ListingID:  req.ListingID,
		Range:      rangeFilter,
		Pagination: pagination,
	}

	ctx = coreutils.ContextWithLogger(ctx)
	result, serviceErr := h.scheduleService.ListBlockEntries(ctx, filter)
	if serviceErr != nil {
		httperrors.SendHTTPErrorObj(c, serviceErr)
		return
	}

	page, limit := schedulePaginationValues(pagination)
	response := converters.ScheduleBlockEntriesToDTO(result, page, limit)

	c.JSON(http.StatusOK, response)
}
