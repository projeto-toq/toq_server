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

// UpdateCalendar handles PUT /admin/holidays/calendars.
//
// @Summary		Update holiday calendar
// @Description	Updates an existing holiday calendar metadata.
// @Tags		Admin Holidays
// @Accept		json
// @Produce	json
// @Param	request	body	dto.HolidayCalendarUpdateRequest	true	"Calendar update payload" Extensions(x-example={"id":42,"name":"Calendario Campinas","scope":"CITY","state":"SP","city":"Campinas","isActive":true,"timezone":"America/Sao_Paulo"})
// @Success	200	{object}	dto.HolidayCalendarResponse
// @Failure	400	{object}	dto.ErrorResponse
// @Failure	401	{object}	dto.ErrorResponse
// @Failure	403	{object}	dto.ErrorResponse
// @Failure	404	{object}	dto.ErrorResponse
// @Failure	500	{object}	dto.ErrorResponse
// @Router		/admin/holidays/calendars [put]
// @Security	BearerAuth
func (h *HolidayHandler) UpdateCalendar(c *gin.Context) {
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

	var req dto.HolidayCalendarUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httperrors.SendHTTPError(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request payload")
		return
	}

	scope, err := parseHolidayScope(req.Scope)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	input := holidayservices.UpdateCalendarInput{
		ID:       req.ID,
		Name:     req.Name,
		Scope:    scope,
		State:    req.State,
		City:     req.City,
		IsActive: req.IsActive,
		Timezone: req.Timezone,
	}

	ctx = coreutils.ContextWithLogger(ctx)
	calendar, serviceErr := h.holidayService.UpdateCalendar(ctx, input)
	if serviceErr != nil {
		httperrors.SendHTTPErrorObj(c, serviceErr)
		return
	}

	response := converters.HolidayCalendarToDTO(calendar)
	c.JSON(http.StatusOK, response)
}
