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

// GetListingAvailability handles GET /schedules/listing/availability.
//
// @Summary	List listing availability
// @Description	Returns the available slots for a listing within the provided time range.
// @Tags	Schedules
// @Produce	json
// @Param	listingId	query	int64	true	"Listing identifier"
// @Param	rangeFrom	query	string	false	"Start of time range (RFC3339)"
// @Param	rangeTo	query	string	false	"End of time range (RFC3339)"
// @Param	slotDurationMinute	query	int	false	"Desired slot duration in minutes"
// @Param	page	query	int	false	"Page number"
// @Param	limit	query	int	false	"Items per page"
// @Param	timezone	query	string	false	"Timezone identifier (IANA)" default(America/Sao_Paulo)
// @Success	200	{object}	dto.ScheduleAvailabilityResponse
// @Failure	400	{object}	dto.ErrorResponse
// @Failure	401	{object}	dto.ErrorResponse
// @Failure	403	{object}	dto.ErrorResponse
// @Failure	404	{object}	dto.ErrorResponse
// @Failure	500	{object}	dto.ErrorResponse
// @Router	/schedules/listing/availability [get]
// @Security	BearerAuth
func (h *ScheduleHandler) GetListingAvailability(c *gin.Context) {
	baseCtx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)
	ctx, spanEnd, err := coreutils.GenerateTracer(baseCtx)
	if err != nil {
		httperrors.SendHTTPError(c, http.StatusInternalServerError, "TRACER_ERROR", "Failed to generate tracer")
		return
	}
	defer spanEnd()

	if _, ok := middlewares.GetUserInfoFromContext(c); !ok {
		httperrors.SendHTTPError(c, http.StatusUnauthorized, "UNAUTHENTICATED", "User context not found")
		return
	}

	var req dto.ScheduleAvailabilityQuery
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

	filter := schedulemodel.AvailabilityFilter{
		ListingID:          req.ListingID,
		Range:              rangeFilter,
		SlotDurationMinute: req.SlotDurationMinute,
		Pagination:         pagination,
		Timezone:           req.Timezone,
	}

	ctx = coreutils.ContextWithLogger(ctx)
	result, serviceErr := h.scheduleService.GetAvailability(ctx, filter)
	if serviceErr != nil {
		httperrors.SendHTTPErrorObj(c, serviceErr)
		return
	}

	page, limit := schedulePaginationValues(pagination)
	response := converters.ScheduleAvailabilityToDTO(result, page, limit)

	c.JSON(http.StatusOK, response)
}
