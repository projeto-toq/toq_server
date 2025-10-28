package schedulehandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/converters"
	dto "github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/middlewares"
	scheduleservices "github.com/projeto-toq/toq_server/internal/core/service/schedule_service"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

// PostCreateBlockEntry handles POST /schedules/listing/block.
//
// @Summary		Create a blocking entry
// @Description	Creates a blocking or temporary blocking time range for a listing agenda.
// @Tags		Schedules
// @Accept		json
// @Produce	json
// @Param		request	body	dto.ScheduleBlockEntryRequest	true	"Block entry payload" Extensions(x-example={"listingId":3241,"entryType":"BLOCK","startsAt":"2025-06-15T09:00:00-03:00","endsAt":"2025-06-15T11:00:00-03:00","reason":"Janela de manutencao","timezone":"America/Sao_Paulo"})
// @Success	200	{object}	dto.ScheduleBlockEntryResponse
// @Failure	400	{object}	dto.ErrorResponse
// @Failure	401	{object}	dto.ErrorResponse
// @Failure	403	{object}	dto.ErrorResponse
// @Failure	404	{object}	dto.ErrorResponse
// @Failure	409	{object}	dto.ErrorResponse
// @Failure	500	{object}	dto.ErrorResponse
// @Router		/schedules/listing/block [post]
// @Security	BearerAuth
func (h *ScheduleHandler) PostCreateBlockEntry(c *gin.Context) {
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

	var req dto.ScheduleBlockEntryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httperrors.SendHTTPError(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request payload")
		return
	}

	typeValue, err := parseScheduleEntryType(req.EntryType)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	startsAt, err := parseScheduleTimestamp("startsAt", req.StartsAt)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	endsAt, err := parseScheduleTimestamp("endsAt", req.EndsAt)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	input := scheduleservices.CreateBlockEntryInput{
		ListingID: req.ListingID,
		OwnerID:   userInfo.ID,
		EntryType: typeValue,
		StartsAt:  startsAt,
		EndsAt:    endsAt,
		Reason:    req.Reason,
		ActorID:   userInfo.ID,
		Timezone:  req.Timezone,
	}

	ctx = coreutils.ContextWithLogger(ctx)
	entry, serviceErr := h.scheduleService.CreateBlockEntry(ctx, input)
	if serviceErr != nil {
		httperrors.SendHTTPErrorObj(c, serviceErr)
		return
	}

	response := dto.ScheduleBlockEntryResponse{Entry: converters.ScheduleEntryToDTO(entry)}
	c.JSON(http.StatusOK, response)
}
