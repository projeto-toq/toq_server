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

// GetCalendarDateDetail handles POST /admin/holidays/dates/detail.
//
// @Summary      Get holiday date detail
// @Description  Retrieves details about a specific holiday date entry.
// @Tags         Admin Holidays
// @Accept       json
// @Produce      json
// @Param        request body dto.HolidayCalendarDateDetailRequest true "Holiday date identifier" Extensions(x-example={"id":10,"timezone":"America/Sao_Paulo"})
// @Success      200 {object} dto.HolidayCalendarDateResponse
// @Failure      400 {object} dto.ErrorResponse
// @Failure      401 {object} dto.ErrorResponse
// @Failure      403 {object} dto.ErrorResponse
// @Failure      404 {object} dto.ErrorResponse
// @Failure      500 {object} dto.ErrorResponse
// @Router       /admin/holidays/dates/detail [post]
// @Security     BearerAuth
func (h *HolidayHandler) GetCalendarDateDetail(c *gin.Context) {
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

	var req dto.HolidayCalendarDateDetailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httperrors.SendHTTPError(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request payload")
		return
	}

	ctx = coreutils.ContextWithLogger(ctx)
	date, serviceErr := h.holidayService.GetCalendarDateByID(ctx, req.ID, req.Timezone)
	if serviceErr != nil {
		httperrors.SendHTTPErrorObj(c, serviceErr)
		return
	}

	response := converters.HolidayCalendarDateToDTO(date)
	c.JSON(http.StatusOK, response)
}
