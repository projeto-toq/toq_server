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

// ListMedia retrieves a paginated list of media assets for a listing
//
//	@Summary		List media assets
//	@Description	Retrieves a paginated list of media assets for a listing, with optional filtering by asset type.
//	@Tags			Listings Media
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			listingIdentityId	query		int		true	"Listing Identity ID"
//	@Param			assetType			query		string	false	"Asset Type"
//	@Param			page				query		int		false	"Page number"
//	@Param			limit				query		int		false	"Items per page"
//	@Param			sort				query		string	false	"Sort field"
//	@Param			order				query		string	false	"Sort order"
//	@Success		200		{object}	dto.ListMediaResponse	"Media assets list"
//	@Failure		400		{object}	dto.ErrorResponse	"Invalid request"
//	@Failure		401		{object}	dto.ErrorResponse	"Unauthorized"
//	@Failure		403		{object}	dto.ErrorResponse	"Forbidden"
//	@Failure		404		{object}	dto.ErrorResponse	"Listing not found"
//	@Failure		500		{object}	dto.ErrorResponse	"Internal server error"
//	@Router			/listings/media [get]
func (h *MediaProcessingHandler) ListMedia(c *gin.Context) {
	baseCtx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)
	ctx, spanEnd, err := coreutils.GenerateTracer(baseCtx)
	if err != nil {
		httperrors.SendHTTPError(c, http.StatusInternalServerError, "TRACER_ERROR", "Failed to generate tracer")
		return
	}
	defer spanEnd()

	var request dto.ListMediaRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		httperrors.SendHTTPErrorObj(c, httputils.MapBindingError(err))
		return
	}

	input := converters.DTOToListMediaInput(request)
	output, err := h.service.ListMedia(ctx, input)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	response := converters.ListMediaOutputToDTO(output)
	c.JSON(http.StatusOK, response)
}
