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

// PostFinishListingAgenda handles POST /schedules/listing/finish.
//
// @Summary	Confirm listing agenda creation
// @Description	Marks the listing agenda as finished and moves listing to pending photo scheduling.
// @Tags		Listing Schedules
// @Accept	json
// @Produce	json
// @Param	request	body	dto.ScheduleFinishAgendaRequest	true	"Finish agenda payload" Extensions(x-example={"listingId":3241})
// @Success	200	{object}	dto.APIResponse
// @Failure	400	{object}	dto.ErrorResponse
// @Failure	401	{object}	dto.ErrorResponse
// @Failure	403	{object}	dto.ErrorResponse
// @Failure	404	{object}	dto.ErrorResponse
// @Failure	409	{object}	dto.ErrorResponse
// @Failure	500	{object}	dto.ErrorResponse
// @Router	/schedules/listing/finish [post]
// @Security	BearerAuth
func (h *ScheduleHandler) PostFinishListingAgenda(c *gin.Context) {
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

	var req dto.ScheduleFinishAgendaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httperrors.SendHTTPError(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request payload")
		return
	}

	input := scheduleservices.FinishListingAgendaInput{
		ListingID: req.ListingID,
		OwnerID:   userInfo.ID,
		ActorID:   userInfo.ID,
	}

	ctx = coreutils.ContextWithLogger(ctx)
	if serviceErr := h.scheduleService.FinishListingAgenda(ctx, input); serviceErr != nil {
		httperrors.SendHTTPErrorObj(c, serviceErr)
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(nil, "Listing agenda finished"))
}
