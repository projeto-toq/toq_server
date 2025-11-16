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

// ListDownloadURLs generates signed download URLs for processed assets
//
//	@Summary		Get signed download URLs for processed media
//	@Description	Finds the most recent READY batch for a listing and generates time-limited signed S3 URLs for downloading processed assets (optimized images/videos). Includes preview URLs for thumbnails. Returns 204 if no READY batch exists.
//	@Tags			Listings Media
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		dto.ListDownloadURLsRequest		true	"Listing identification"
//	@Success		200		{object}	dto.ListDownloadURLsResponse	"Download URLs generated"
//	@Success		204		"No READY batch available"
//	@Failure		400		{object}	dto.ErrorResponse	"Invalid request"
//	@Failure		401		{object}	dto.ErrorResponse	"Unauthorized"
//	@Failure		403		{object}	dto.ErrorResponse	"Forbidden"
//	@Failure		404		{object}	dto.ErrorResponse	"Listing not found"
//	@Failure		500		{object}	dto.ErrorResponse	"Internal server error"
//	@Router			/api/v2/listings/media/downloads [post]
func (lh *ListingHandler) ListDownloadURLs(c *gin.Context) {
	baseCtx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)
	ctx, spanEnd, err := coreutils.GenerateTracer(baseCtx)
	if err != nil {
		httperrors.SendHTTPError(c, http.StatusInternalServerError, "TRACER_ERROR", "Failed to generate tracer")
		return
	}
	defer spanEnd()

	var request dto.ListDownloadURLsRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httperrors.SendHTTPErrorObj(c, httputils.MapBindingError(err))
		return
	}

	input := converters.DTOToListDownloadURLsInput(request)
	output, err := lh.mediaProcessingService.ListDownloadURLs(ctx, input)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	response := converters.ListDownloadURLsOutputToDTO(output)
	c.JSON(http.StatusOK, response)
}
