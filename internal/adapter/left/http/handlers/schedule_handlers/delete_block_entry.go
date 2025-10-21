package schedulehandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	dto "github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/middlewares"
	scheduleservices "github.com/projeto-toq/toq_server/internal/core/service/schedule_service"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

// DeleteBlockEntry handles DELETE /schedules/listing/block.
//
// @Summary		Delete a blocking entry
// @Description	Removes a blocking or temporary block entry from a listing agenda.
// @Tags		Schedules
// @Accept		json
// @Produce	json
// @Param		request	body	dto.ScheduleDeleteEntryRequest	true	"Block entry deletion payload" Extensions(x-example={"entryId":5021,"listingId":3241})
// @Success	204	"Entry deleted successfully"
// @Failure	400	{object}	dto.ErrorResponse
// @Failure	401	{object}	dto.ErrorResponse
// @Failure	403	{object}	dto.ErrorResponse
// @Failure	404	{object}	dto.ErrorResponse
// @Failure	500	{object}	dto.ErrorResponse
// @Router		/schedules/listing/block [delete]
// @Security	BearerAuth
func (h *ScheduleHandler) DeleteBlockEntry(c *gin.Context) {
	baseCtx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)
	ctx, spanEnd, err := coreutils.GenerateTracer(baseCtx)
	if err != nil {
		httperrors.SendHTTPError(c, http.StatusInternalServerError, "TRACER_ERROR", "Failed to generate tracer")
		return
	}
	defer spanEnd()

	userInfo, ok := middlewares.GetUserInfoFromContext(c)
	if !ok {
		httperrors.SendHTTPError(c, http.StatusUnauthorized, "UNAUTHENTICATED", "User context not found")
		return
	}

	var req dto.ScheduleDeleteEntryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httperrors.SendHTTPError(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request payload")
		return
	}

	input := scheduleservices.DeleteEntryInput{
		EntryID:   req.EntryID,
		ListingID: req.ListingID,
		OwnerID:   userInfo.ID,
	}

	ctx = coreutils.ContextWithLogger(ctx)
	if err := h.scheduleService.DeleteBlockEntry(ctx, input); err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
