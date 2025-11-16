package listinghandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/handlers/listing_handlers/converters"
	httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	httputils "github.com/projeto-toq/toq_server/internal/adapter/left/http/utils"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

// GetBatchStatus retrieves detailed batch status for UI polling
//
//	@Summary		Get media batch processing status
//	@Description	Returns batch status (PENDING_UPLOAD/RECEIVED/PROCESSING/READY/FAILED) along with asset details (titles, sequences, object keys). Used by frontend to poll progress and determine when downloads are available.
//	@Tags			Listings Media
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		dto.GetBatchStatusRequest	true	"Batch identification"
//	@Success		200		{object}	dto.GetBatchStatusResponse	"Batch status retrieved"
//	@Failure		400		{object}	dto.ErrorResponse	"Invalid request"
//	@Failure		401		{object}	dto.ErrorResponse	"Unauthorized"
//	@Failure		403		{object}	dto.ErrorResponse	"Forbidden"
//	@Failure		404		{object}	dto.ErrorResponse	"Batch not found"
//	@Failure		500		{object}	dto.ErrorResponse	"Internal server error"
//	@Router			/api/v2/listings/media/status [post]
func (lh *ListingHandler) GetBatchStatus(c *gin.Context) {
	baseCtx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)
	ctx, spanEnd, err := coreutils.GenerateTracer(baseCtx)
	if err != nil {
		httperrors.SendHTTPError(c, http.StatusInternalServerError, "TRACER_ERROR", "Failed to generate tracer")
		return
	}
	defer spanEnd()

	var request dto.GetBatchStatusRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httperrors.SendHTTPErrorObj(c, httputils.MapBindingError(err))
		return
	}

	input := converters.DTOToGetBatchStatusInput(request)
	output, err := lh.mediaProcessingService.GetBatchStatus(ctx, input)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	response := converters.GetBatchStatusOutputToDTO(output)
	c.JSON(http.StatusOK, response)
}
