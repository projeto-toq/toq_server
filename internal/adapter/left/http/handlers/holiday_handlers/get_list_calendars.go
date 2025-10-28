package holidayhandlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/converters"
	dto "github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/middlewares"
	holidaymodel "github.com/projeto-toq/toq_server/internal/core/model/holiday_model"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

// ListCalendars handles GET /admin/holidays/calendars.
//
// @Summary		List holiday calendars
// @Description	Returns holiday calendars with filtering and pagination options.
// @Tags		Admin Holidays
// @Accept		json
// @Produce	json
// @Param		 scope	query	string	false	"Calendar scope (NATIONAL|STATE|CITY)" example("STATE")
// @Param		 state	query	string	false	"State abbreviation" example("SP")
// @Param		 city	query	string	false	"City name" example("Campinas")
// @Param		 search	query	string	false	"Free text search" example("Christmas")
// @Param		 onlyActive	query	bool	false	"Filter active calendars" example(true)
// @Param		 page	query	int	false	"Page number" example(1)
// @Param		 limit	query	int	false	"Page size" example(20)
// @Success	200	{object}	dto.HolidayCalendarsListResponse
// @Failure	400	{object}	dto.ErrorResponse
// @Failure	401	{object}	dto.ErrorResponse
// @Failure	403	{object}	dto.ErrorResponse
// @Failure	500	{object}	dto.ErrorResponse
// @Router		/admin/holidays/calendars [get]
// @Security	BearerAuth
func (h *HolidayHandler) ListCalendars(c *gin.Context) {
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

	var req dto.HolidayCalendarsListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		httperrors.SendHTTPError(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid query parameters")
		return
	}

	scopeValue, err := parseHolidayScope(req.Scope)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	var scopePtr *holidaymodel.CalendarScope
	if scopeValue != "" {
		scopePtr = new(holidaymodel.CalendarScope)
		*scopePtr = scopeValue
	}

	state := strings.TrimSpace(req.State)
	city := strings.TrimSpace(req.City)

	var statePtr *string
	if state != "" {
		statePtr = &state
	}

	var cityPtr *string
	if city != "" {
		cityPtr = &city
	}

	page, limit := sanitizeHolidayPagination(req.Page, req.Limit)

	filter := holidaymodel.CalendarListFilter{
		Scope:      scopePtr,
		State:      statePtr,
		City:       cityPtr,
		SearchTerm: strings.TrimSpace(req.SearchTerm),
		OnlyActive: req.OnlyActive,
		Page:       page,
		Limit:      limit,
	}

	ctx = coreutils.ContextWithLogger(ctx)
	result, serviceErr := h.holidayService.ListCalendars(ctx, filter)
	if serviceErr != nil {
		httperrors.SendHTTPErrorObj(c, serviceErr)
		return
	}

	response := converters.HolidayCalendarsToListDTO(result.Calendars, page, limit, result.Total)
	c.JSON(http.StatusOK, response)
}
