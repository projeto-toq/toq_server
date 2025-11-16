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

// CompleteUploadBatch handles batch upload confirmation and processing initiation
//
//	@Summary		Confirm upload completion and start processing
//	@Description	Validates S3 uploads via HEAD requests checking checksums, consolidates asset titles/sequences, updates batch status to RECEIVED, registers processing job and enqueues it to SQS. Returns job ID and estimated processing duration.
//	@Tags			Listings Media
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		dto.CompleteUploadBatchRequest		true	"Upload confirmation with object keys and checksums"
//	@Success		202		{object}	dto.CompleteUploadBatchResponse		"Processing job enqueued"
//	@Failure		400		{object}	dto.ErrorResponse	"Invalid request or checksum mismatch"
//	@Failure		401		{object}	dto.ErrorResponse	"Unauthorized"
//	@Failure		403		{object}	dto.ErrorResponse	"Forbidden"
//	@Failure		404		{object}	dto.ErrorResponse	"Batch not found or not in PENDING_UPLOAD"
//	@Failure		500		{object}	dto.ErrorResponse	"Internal server error"
//	@Router			/api/v2/listings/media/uploads/complete [post]
func (lh *ListingHandler) CompleteUploadBatch(c *gin.Context) {
	baseCtx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)
	ctx, spanEnd, err := coreutils.GenerateTracer(baseCtx)
	if err != nil {
		httperrors.SendHTTPError(c, http.StatusInternalServerError, "TRACER_ERROR", "Failed to generate tracer")
		return
	}
	defer spanEnd()

	var request dto.CompleteUploadBatchRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httperrors.SendHTTPErrorObj(c, httputils.MapBindingError(err))
		return
	}

	input := converters.DTOToCompleteUploadBatchInput(request)
	output, err := lh.mediaProcessingService.CompleteUploadBatch(ctx, input)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	response := converters.CompleteUploadBatchOutputToDTO(output)
	c.JSON(http.StatusAccepted, response)
}
