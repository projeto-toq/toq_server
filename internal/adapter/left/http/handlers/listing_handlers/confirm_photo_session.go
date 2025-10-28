package listinghandlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/middlewares"
	listingservices "github.com/projeto-toq/toq_server/internal/core/service/listing_service"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

// ConfirmPhotoSession confirma uma reserva previamente criada, gerando agendamento definitivo.
//
//	@Summary   Confirm a photo session reservation
//	@Tags      Listing Photo Sessions
//	@Accept    json
//	@Produce   json
//	@Param     request body      dto.ConfirmPhotoSessionRequest true "Confirmation payload" Extensions(x-example={"listingId":1001,"slotId":2002,"reservationToken":"c36b754f-6c37-4c15-8f25-9d77ddf9bb3e"})
//	@Success   200     {object} dto.ConfirmPhotoSessionResponse
//	@Failure   400     {object} dto.ErrorResponse "Invalid payload"
//	@Failure   401     {object} dto.ErrorResponse "Unauthorized"
//	@Failure   403     {object} dto.ErrorResponse "Forbidden"
//	@Failure   404     {object} dto.ErrorResponse "Reservation not found"
//	@Failure   409     {object} dto.ErrorResponse "Slot unavailable"
//	@Failure   500     {object} dto.ErrorResponse "Internal error"
//	@Router    /listings/photo-session/confirm [post]
//	@Security  BearerAuth
func (lh *ListingHandler) ConfirmPhotoSession(c *gin.Context) {
	baseCtx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)
	ctx, spanEnd, err := coreutils.GenerateTracer(baseCtx)
	if err != nil {
		httperrors.SendHTTPError(c, http.StatusInternalServerError, "TRACER_ERROR", "Failed to generate tracer")
		return
	}
	defer spanEnd()

	if _, ok := middlewares.GetUserInfoFromContext(c); !ok {
		httperrors.SendHTTPError(c, http.StatusInternalServerError, "INTERNAL_CONTEXT_MISSING", "User context not found")
		return
	}

	var request dto.ConfirmPhotoSessionRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httperrors.SendHTTPError(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid confirmation payload")
		return
	}

	if request.ListingID <= 0 || request.SlotID == 0 || request.ReservationToken == "" {
		httperrors.SendHTTPError(c, http.StatusBadRequest, "INVALID_REQUEST", "listingId, slotId and reservationToken are required")
		return
	}

	input := listingservices.ConfirmPhotoSessionInput{
		ListingID:        request.ListingID,
		SlotID:           request.SlotID,
		ReservationToken: request.ReservationToken,
	}

	output, err := lh.listingService.ConfirmPhotoSession(ctx, input)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.ConfirmPhotoSessionResponse{
		PhotoSessionID: output.PhotoSessionID,
		SlotID:         output.SlotID,
		ScheduledStart: output.ScheduledStart.UTC().Format(time.RFC3339),
		ScheduledEnd:   output.ScheduledEnd.UTC().Format(time.RFC3339),
	})
}
