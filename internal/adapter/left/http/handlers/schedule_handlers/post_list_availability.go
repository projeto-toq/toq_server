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

// PostListAvailability handles POST /schedules/listing/availability.
//
// @Summary		List listing availability
// @Description	Returns the available slots for a listing within the provided time range.
// @Tags		Schedules
// @Accept		json
// @Produce	json
// @Param		request	body	dto.ScheduleAvailabilityRequest	true	"Availability filter"
// @Success	200	{object}	dto.ScheduleAvailabilityResponse
// @Failure	400	{object}	dto.ErrorResponse
// @Failure	401	{object}	dto.ErrorResponse
// @Failure	403	{object}	dto.ErrorResponse
// @Failure	404	{object}	dto.ErrorResponse
// @Failure	500	{object}	dto.ErrorResponse
// @Router		/schedules/listing/availability [post]
// @Security	BearerAuth
func (h *ScheduleHandler) PostListAvailability(c *gin.Context) {
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

	var req dto.ScheduleAvailabilityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httperrors.SendHTTPError(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request payload")
		return
	}

	rangeFilter, err := parseScheduleRange(req.Range)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	pagination := sanitizeSchedulePagination(req.Pagination)

	filter := schedulemodel.AvailabilityFilter{
		ListingID:          req.ListingID,
		Range:              rangeFilter,
		SlotDurationMinute: req.SlotDurationMinute,
		Pagination:         pagination,
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
