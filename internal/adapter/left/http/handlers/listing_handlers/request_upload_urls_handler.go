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

// RequestUploadURLs handles signed URL issuance for listing media uploads
//
//	@Summary		Request upload URLs for a listing media batch
//	@Description	Validates permissions (photographer booking or owner for project media), enforces listing status PENDING_PHOTO_PROCESSING, creates a media batch and returns signed PUT URLs with required headers (Content-Type, checksum). Rejects duplicates (sequence/title) and unsupported media types.
//	@Tags			Listings Media
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		dto.RequestUploadURLsRequest	true	"Manifest with client-side file metadata"
//	@Success		201		{object}	dto.RequestUploadURLsResponse	"Batch created with signed URLs"
//	@Failure		400		{object}	dto.ErrorResponse	"Invalid request format"
//	@Failure		401		{object}	dto.ErrorResponse	"Unauthorized"
//	@Failure		403		{object}	dto.ErrorResponse	"Forbidden - User not allowed to upload for this listing"
//	@Failure		404		{object}	dto.ErrorResponse	"Listing not found or not in PENDING_PHOTO_PROCESSING"
//	@Failure		409		{object}	dto.ErrorResponse	"There is another open batch"
//	@Failure		500		{object}	dto.ErrorResponse	"Internal server error"
//	@Router			/listings/media/uploads [post]
func (lh *ListingHandler) RequestUploadURLs(c *gin.Context) {
	baseCtx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)
	ctx, spanEnd, err := coreutils.GenerateTracer(baseCtx)
	if err != nil {
		httperrors.SendHTTPError(c, http.StatusInternalServerError, "TRACER_ERROR", "Failed to generate tracer")
		return
	}
	defer spanEnd()

	var request dto.RequestUploadURLsRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httperrors.SendHTTPErrorObj(c, httputils.MapBindingError(err))
		return
	}

	input := converters.DTOToRequestUploadURLsInput(request)
	// Inject RequestedBy from context
	if userID, ok := coreutils.GetUserIDFromContext(ctx); ok {
		input.RequestedBy = userID
	}

	output, err := lh.mediaProcessingService.RequestUploadURLs(ctx, input)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	response := converters.RequestUploadURLsOutputToDTO(output)
	c.JSON(http.StatusCreated, response)
}
