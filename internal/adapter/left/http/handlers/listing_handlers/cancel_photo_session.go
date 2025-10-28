package listinghandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/middlewares"
	listingservices "github.com/projeto-toq/toq_server/internal/core/service/listing_service"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

// CancelPhotoSession cancela uma sessão de fotos previamente agendada.
//
//	@Summary   Cancel a booked photo session
//	@Tags      Listing Photo Sessions
//	@Accept    json
//	@Produce   json
//	@Param     request body      dto.CancelPhotoSessionRequest true "Cancel payload" Extensions(x-example={"photoSessionId":3003})
//	@Success   200     {object} dto.APIResponse
//	@Failure   400     {object} dto.ErrorResponse "Invalid payload"
//	@Failure   401     {object} dto.ErrorResponse "Unauthorized"
//	@Failure   403     {object} dto.ErrorResponse "Forbidden"
//	@Failure   404     {object} dto.ErrorResponse "Photo session not found"
//	@Failure   500     {object} dto.ErrorResponse "Internal error"
//	@Router    /listings/photo-session/cancel [post]
//	@Security  BearerAuth
func (lh *ListingHandler) CancelPhotoSession(c *gin.Context) {
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

	var request dto.CancelPhotoSessionRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httperrors.SendHTTPError(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid cancel payload")
		return
	}

	if request.PhotoSessionID == 0 {
		httperrors.SendHTTPError(c, http.StatusBadRequest, "INVALID_REQUEST", "photoSessionId is required")
		return
	}

	input := listingservices.CancelPhotoSessionInput{
		PhotoSessionID: request.PhotoSessionID,
	}

	if err := lh.listingService.CancelPhotoSession(ctx, input); err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(nil, "Photo session cancelled"))
}
