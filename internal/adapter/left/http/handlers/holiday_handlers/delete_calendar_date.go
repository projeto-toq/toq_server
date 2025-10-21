package holidayhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	dto "github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/middlewares"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

// DeleteCalendarDate handles DELETE /admin/holidays/dates.
//
// @Summary		Delete holiday date
// @Description	Removes a holiday date from a calendar.
// @Tags		Admin Holidays
// @Accept		json
// @Produce	json
// @Param		request	body	dto.HolidayCalendarDateDeleteRequest	true	"Holiday date delete request" Extensions(x-example={"id":1001})
// @Success	204	"Holiday date deleted"
// @Failure	400	{object}	dto.ErrorResponse
// @Failure	401	{object}	dto.ErrorResponse
// @Failure	403	{object}	dto.ErrorResponse
// @Failure	404	{object}	dto.ErrorResponse
// @Failure	500	{object}	dto.ErrorResponse
// @Router		/admin/holidays/dates [delete]
// @Security	BearerAuth
func (h *HolidayHandler) DeleteCalendarDate(c *gin.Context) {
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

	var req dto.HolidayCalendarDateDeleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httperrors.SendHTTPError(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request payload")
		return
	}

	ctx = coreutils.ContextWithLogger(ctx)
	if err := h.holidayService.DeleteCalendarDate(ctx, req.ID); err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
