package holidayhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/converters"
	dto "github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/middlewares"
	holidaymodel "github.com/projeto-toq/toq_server/internal/core/model/holiday_model"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

// ListCalendarDates handles GET /admin/holidays/dates.
//
// @Summary		List holiday dates
// @Description	Returns holiday dates for a calendar with optional range filters.
// @Tags		Admin Holidays
// @Accept		json
// @Produce	json
// @Param		 calendarId	query	int	true	"Calendar identifier" example(42)
// @Param		 from	query	string	false	"Start date (RFC3339)" example("2025-12-01T00:00:00Z")
// @Param		 to	query	string	false	"End date (RFC3339)" example("2026-01-10T23:59:59Z")
// @Param		 page	query	int	false	"Page number" example(1)
// @Param		 limit	query	int	false	"Page size" example(50)
// @Success	200	{object}	dto.HolidayCalendarDatesListResponse
// @Failure	400	{object}	dto.ErrorResponse
// @Failure	401	{object}	dto.ErrorResponse
// @Failure	403	{object}	dto.ErrorResponse
// @Failure	500	{object}	dto.ErrorResponse
// @Router		/admin/holidays/dates [get]
// @Security	BearerAuth
func (h *HolidayHandler) ListCalendarDates(c *gin.Context) {
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

	var req dto.HolidayCalendarDatesListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		httperrors.SendHTTPError(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid query parameters")
		return
	}

	from, err := parseOptionalHolidayDate("from", req.From)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	to, err := parseOptionalHolidayDate("to", req.To)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	page, limit := sanitizeHolidayPagination(req.Page, req.Limit)

	filter := holidaymodel.CalendarDatesFilter{
		CalendarID: req.CalendarID,
		From:       from,
		To:         to,
		Timezone:   req.Timezone,
		Page:       page,
		Limit:      limit,
	}

	ctx = coreutils.ContextWithLogger(ctx)
	result, serviceErr := h.holidayService.ListCalendarDates(ctx, filter)
	if serviceErr != nil {
		httperrors.SendHTTPErrorObj(c, serviceErr)
		return
	}

	response := converters.HolidayCalendarDatesToListDTO(result.Dates, page, limit, result.Total)
	c.JSON(http.StatusOK, response)
}
