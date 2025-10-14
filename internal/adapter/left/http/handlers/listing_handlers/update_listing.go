package listinghandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/converters"
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/middlewares"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

// UpdateListing updates an existing listing.
//
//	@Summary	Update a listing
//	@Description	Allows partial updates for draft listings. Omitted fields remain unchanged; present fields (including null/empty) overwrite stored values.
//	@Tags		Listings
//	@Accept		json
//	@Produce	json
//	@Param		request	body	dto.UpdateListingRequest	true	"Payload for update (ID must be provided in the body)"
//	@Success	200	{object}	dto.UpdateListingResponse
//	@Failure	400	{object}	dto.ErrorResponse	"Invalid payload"
//	@Failure	401	{object}	dto.ErrorResponse	"Unauthorized"
//	@Failure	403	{object}	dto.ErrorResponse	"Forbidden"
//	@Failure	404	{object}	dto.ErrorResponse	"Not found"
//	@Failure	409	{object}	dto.ErrorResponse	"Conflict"
//	@Failure	500	{object}	dto.ErrorResponse	"Internal error"
//	@Router		/listings [put]
//	@Security	BearerAuth
func (lh *ListingHandler) UpdateListing(c *gin.Context) {
	baseCtx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)

	if _, ok := middlewares.GetUserInfoFromContext(c); !ok {
		httperrors.SendHTTPError(c, http.StatusInternalServerError, "INTERNAL_CONTEXT_MISSING", "User context not found")
		return
	}

	var request dto.UpdateListingRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httperrors.SendHTTPError(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request format")
		return
	}

	if !request.ID.IsPresent() || request.ID.IsNull() {
		httperrors.SendHTTPError(c, http.StatusBadRequest, "MISSING_ID", "Listing ID must be provided in the request body")
		return
	}

	listingID, ok := request.ID.Value()
	if !ok || listingID <= 0 {
		httperrors.SendHTTPError(c, http.StatusBadRequest, "INVALID_ID", "Listing ID is invalid")
		return
	}

	input, err := converters.UpdateListingRequestToInput(request)
	if err != nil {
		httperrors.SendHTTPError(c, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}
	input.ID = listingID

	if err := lh.listingService.UpdateListing(baseCtx, input); err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.UpdateListingResponse{
		Success: true,
		Message: "Listing updated",
	})
}
