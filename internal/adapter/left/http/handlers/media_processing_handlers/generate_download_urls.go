package mediaprocessinghandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/handlers/listing_handlers/converters"
	httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	httputils "github.com/projeto-toq/toq_server/internal/adapter/left/http/utils"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

// GenerateDownloadURLs generates signed download URLs for specific assets
//
//	@Summary		Generate signed download URLs
//	@Description	Generates time-limited signed S3 URLs for specific assets. Supports multiple assets/resolutions, including the final ZIP bundle (assetType=ZIP, resolution="zip" ignored). Valid resolutions: thumbnail (200px), small (400px), medium (800px), large (1200px), original, zip.
//	@Tags			Listings Media
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		dto.GenerateDownloadURLsRequest	true	"Download requests" Extensions(x-example={"listingIdentityId": 51, "requests": [{"assetType": "PHOTO_VERTICAL", "sequence": 1, "resolution": "thumbnail"}, {"assetType": "PHOTO_HORIZONTAL", "sequence": 2, "resolution": "large"}, {"assetType": "VIDEO_VERTICAL", "sequence": 1, "resolution": "original"}, {"assetType": "ZIP", "sequence": 0, "resolution": "zip"}]})
//	@Success		200		{object}	dto.GenerateDownloadURLsResponse	"Download URLs generated"
//	@Failure		400		{object}	dto.ErrorResponse	"Invalid request"
//	@Failure		401		{object}	dto.ErrorResponse	"Unauthorized"
//	@Failure		403		{object}	dto.ErrorResponse	"Forbidden"
//	@Failure		404		{object}	dto.ErrorResponse	"Listing not found"
//	@Failure		500		{object}	dto.ErrorResponse	"Internal server error"
//	@Router			/listings/media/download [post]
func (h *MediaProcessingHandler) GenerateDownloadURLs(c *gin.Context) {
	baseCtx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)
	ctx, spanEnd, err := coreutils.GenerateTracer(baseCtx)
	if err != nil {
		httperrors.SendHTTPError(c, http.StatusInternalServerError, "TRACER_ERROR", "Failed to generate tracer")
		return
	}
	defer spanEnd()

	var request dto.GenerateDownloadURLsRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httperrors.SendHTTPErrorObj(c, httputils.MapBindingError(err))
		return
	}

	input := converters.DTOToGenerateDownloadURLsInput(request)
	output, err := h.service.GenerateDownloadURLs(ctx, input)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	response := converters.GenerateDownloadURLsOutputToDTO(output)
	c.JSON(http.StatusOK, response)
}
