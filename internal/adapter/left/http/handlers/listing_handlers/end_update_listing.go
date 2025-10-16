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

// EndUpdateListing finalizes a draft listing so it can move to the photo scheduling stage.
//
//	@Summary    Finalize listing update
//	@Description    Validates all required listing attributes and transitions status to photo scheduling.
//	@Tags       Listings
//	@Accept     json
//	@Produce    json
//	@Param      request body dto.EndUpdateListingRequest true "Listing identifier" Extensions(x-example={"listingId":98765})
//	@Success    200 {object} dto.EndUpdateListingResponse
//	@Failure    400 {object} dto.ErrorResponse "Validation error"
//	@Failure    401 {object} dto.ErrorResponse "Unauthorized"
//	@Failure    403 {object} dto.ErrorResponse "Forbidden"
//	@Failure    404 {object} dto.ErrorResponse "Listing not found"
//	@Failure    409 {object} dto.ErrorResponse "Conflict"
//	@Failure    500 {object} dto.ErrorResponse "Internal error"
//	@Router     /listings/end-update [post]
//	@Security   BearerAuth
func (lh *ListingHandler) EndUpdateListing(c *gin.Context) {
	baseCtx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)

	if _, ok := middlewares.GetUserInfoFromContext(c); !ok {
		httperrors.SendHTTPError(c, http.StatusInternalServerError, "INTERNAL_CONTEXT_MISSING", "User context not found")
		return
	}

	var request dto.EndUpdateListingRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httperrors.SendHTTPError(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request format")
		return
	}

	input := listingservices.EndUpdateListingInput{ListingID: request.ListingID}
	if err := lh.listingService.EndUpdateListing(baseCtx, input); err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.EndUpdateListingResponse{
		Success: true,
		Message: "Listing update finalized",
	})
}
