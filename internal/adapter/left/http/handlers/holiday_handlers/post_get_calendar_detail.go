package holidayhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/converters"
	dto "github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/middlewares"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

// GetCalendarDetail handles POST /admin/holidays/calendars/detail.
//
// @Summary		Get holiday calendar detail
// @Description	Returns the metadata of a specific holiday calendar.
// @Tags		Admin Holidays
// @Accept		json
// @Produce	json
// @Param	request	body	dto.HolidayCalendarDetailRequest	true	"Calendar detail request" Extensions(x-example={"id":42})
// @Success	200	{object}	dto.HolidayCalendarResponse
// @Failure	400	{object}	dto.ErrorResponse
// @Failure	401	{object}	dto.ErrorResponse
// @Failure	403	{object}	dto.ErrorResponse
// @Failure	404	{object}	dto.ErrorResponse
// @Failure	500	{object}	dto.ErrorResponse
// @Router		/admin/holidays/calendars/detail [post]
// @Security	BearerAuth
func (h *HolidayHandler) GetCalendarDetail(c *gin.Context) {
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

	var req dto.HolidayCalendarDetailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httperrors.SendHTTPError(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request payload")
		return
	}

	ctx = coreutils.ContextWithLogger(ctx)
	calendar, serviceErr := h.holidayService.GetCalendarByID(ctx, req.ID)
	if serviceErr != nil {
		httperrors.SendHTTPErrorObj(c, serviceErr)
		return
	}

	response := converters.HolidayCalendarToDTO(calendar)
	c.JSON(http.StatusOK, response)
}
