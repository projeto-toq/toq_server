package holidayhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/converters"
	dto "github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/middlewares"
	holidayservices "github.com/projeto-toq/toq_server/internal/core/service/holiday_service"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

// UpdateCalendarDate handles PUT /admin/holidays/dates.
//
// @Summary      Update holiday date
// @Description  Updates an existing holiday date within a calendar.
// @Tags         Admin Holidays
// @Accept       json
// @Produce      json
// @Param        request body dto.HolidayCalendarDateUpdateRequest true "Holiday date payload" Extensions(x-example={"id":10,"calendarId":42,"holidayDate":"2025-12-25T00:00:00Z","label":"Natal","recurrent":true})
// @Success      200 {object} dto.HolidayCalendarDateResponse
// @Failure      400 {object} dto.ErrorResponse
// @Failure      401 {object} dto.ErrorResponse
// @Failure      403 {object} dto.ErrorResponse
// @Failure      404 {object} dto.ErrorResponse
// @Failure      500 {object} dto.ErrorResponse
// @Router       /admin/holidays/dates [put]
// @Security     BearerAuth
func (h *HolidayHandler) UpdateCalendarDate(c *gin.Context) {
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

	var req dto.HolidayCalendarDateUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httperrors.SendHTTPError(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request payload")
		return
	}

	holidayDate, err := parseHolidayDate("holidayDate", req.HolidayDate)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	input := holidayservices.UpdateCalendarDateInput{
		ID:          req.ID,
		CalendarID:  req.CalendarID,
		HolidayDate: holidayDate,
		Label:       req.Label,
		Recurrent:   req.Recurrent,
	}

	ctx = coreutils.ContextWithLogger(ctx)
	date, serviceErr := h.holidayService.UpdateCalendarDate(ctx, input)
	if serviceErr != nil {
		httperrors.SendHTTPErrorObj(c, serviceErr)
		return
	}

	response := converters.HolidayCalendarDateToDTO(date)
	c.JSON(http.StatusOK, response)
}
