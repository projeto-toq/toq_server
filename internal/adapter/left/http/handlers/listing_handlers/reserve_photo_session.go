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

// ReservePhotoSession cria uma reserva temporária de slot para sessão fotográfica.
//
//	@Summary   Reserve a photo session slot
//	@Tags      Listings
//	@Accept    json
//	@Produce   json
//	@Param     request body      dto.ReservePhotoSessionRequest true "Reservation request" Extensions(x-example={"listingId":1001,"slotId":2002})
//	@Success   200     {object} dto.ReservePhotoSessionResponse
//	@Failure   400     {object} dto.ErrorResponse "Invalid payload"
//	@Failure   401     {object} dto.ErrorResponse "Unauthorized"
//	@Failure   403     {object} dto.ErrorResponse "Forbidden"
//	@Failure   404     {object} dto.ErrorResponse "Slot or listing not found"
//	@Failure   409     {object} dto.ErrorResponse "Slot unavailable"
//	@Failure   500     {object} dto.ErrorResponse "Internal error"
//	@Router    /listings/photo-session/reserve [post]
//	@Security  BearerAuth
func (lh *ListingHandler) ReservePhotoSession(c *gin.Context) {
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

	var request dto.ReservePhotoSessionRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httperrors.SendHTTPError(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid reservation payload")
		return
	}

	if request.ListingID <= 0 || request.SlotID == 0 {
		httperrors.SendHTTPError(c, http.StatusBadRequest, "INVALID_REQUEST", "listingId and slotId are required")
		return
	}

	input := listingservices.ReservePhotoSessionInput{
		ListingID: request.ListingID,
		SlotID:    request.SlotID,
	}

	output, err := lh.listingService.ReservePhotoSession(ctx, input)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.ReservePhotoSessionResponse{
		SlotID:           output.SlotID,
		ReservationToken: output.ReservationToken,
		ExpiresAt:        output.ExpiresAt.UTC().Format(time.RFC3339),
	})
}
