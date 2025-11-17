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

// RetryMediaBatch re-enqueues a terminal batch for reprocessing
//
//	@Summary		Retry media batch processing
//	@Description	Accepts batches in FAILED or READY status and creates a new processing job reusing existing raw S3 objects. Updates batch status to PROCESSING and enqueues the job to SQS with retry flag. Always returns 202 with new job ID.
//	@Tags			Listings Media
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		dto.RetryMediaBatchRequest		true	"Batch identification and retry reason"
//	@Success		202		{object}	dto.RetryMediaBatchResponse		"Retry job enqueued"
//	@Failure		400		{object}	dto.ErrorResponse	"Invalid request or batch not in terminal status"
//	@Failure		401		{object}	dto.ErrorResponse	"Unauthorized"
//	@Failure		403		{object}	dto.ErrorResponse	"Forbidden"
//	@Failure		404		{object}	dto.ErrorResponse	"Batch not found"
//	@Failure		500		{object}	dto.ErrorResponse	"Internal server error"
//	@Router			/listings/media/uploads/retry [post]
func (lh *ListingHandler) RetryMediaBatch(c *gin.Context) {
	baseCtx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)
	ctx, spanEnd, err := coreutils.GenerateTracer(baseCtx)
	if err != nil {
		httperrors.SendHTTPError(c, http.StatusInternalServerError, "TRACER_ERROR", "Failed to generate tracer")
		return
	}
	defer spanEnd()

	var request dto.RetryMediaBatchRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httperrors.SendHTTPErrorObj(c, httputils.MapBindingError(err))
		return
	}

	input := converters.DTOToRetryMediaBatchInput(request)
	output, err := lh.mediaProcessingService.RetryMediaBatch(ctx, input)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	response := converters.RetryMediaBatchOutputToDTO(output)
	c.JSON(http.StatusAccepted, response)
}
