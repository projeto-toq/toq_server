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

// CreateCalendar handles POST /admin/holidays/calendars.
//
// @Summary		Create holiday calendar
// @Description	Creates a new holiday calendar available for scheduling operations.
// @Tags		Admin Holidays
// @Accept		json
// @Produce	json
// @Param		request	body	dto.HolidayCalendarCreateRequest	true	"Calendar payload" Extensions(x-example={"name":"Feriados Sao Paulo","scope":"STATE","state":"SP","cityIbge":"","isActive":true})
// @Success	201	{object}	dto.HolidayCalendarResponse
// @Failure	400	{object}	dto.ErrorResponse
// @Failure	401	{object}	dto.ErrorResponse
// @Failure	403	{object}	dto.ErrorResponse
// @Failure	500	{object}	dto.ErrorResponse
// @Router		/admin/holidays/calendars [post]
// @Security	BearerAuth
func (h *HolidayHandler) CreateCalendar(c *gin.Context) {
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

	var req dto.HolidayCalendarCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httperrors.SendHTTPError(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request payload")
		return
	}

	scope, err := parseHolidayScope(req.Scope)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	input := holidayservices.CreateCalendarInput{
		Name:     req.Name,
		Scope:    scope,
		State:    req.State,
		CityIBGE: req.CityIBGE,
		IsActive: req.IsActive,
	}
	ctx = coreutils.ContextWithLogger(ctx)
	calendar, serviceErr := h.holidayService.CreateCalendar(ctx, input)
	if serviceErr != nil {
		httperrors.SendHTTPErrorObj(c, serviceErr)
		return
	}

	response := converters.HolidayCalendarToDTO(calendar)
	c.JSON(http.StatusCreated, response)
}
