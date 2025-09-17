package listinghandlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/giulio-alfieri/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/giulio-alfieri/toq_server/internal/adapter/left/http/http_errors"
	"github.com/giulio-alfieri/toq_server/internal/adapter/left/http/middlewares"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

// DeleteListing handles deleting an existing listing
func (lh *ListingHandler) DeleteListing(c *gin.Context) {
	ctx, spanEnd, err := utils.GenerateTracer(c.Request.Context())
	if err != nil {
		httperrors.SendHTTPError(c, http.StatusInternalServerError, "TRACER_ERROR", "Failed to generate tracer")
		return
	}
	defer spanEnd()

	// Get user info from context (set by auth middleware)
	if _, ok := middlewares.GetUserInfoFromContext(c); !ok {
		// Se chegar aqui, é erro de pipeline (middleware deveria ter setado)
		httperrors.SendHTTPError(c, http.StatusInternalServerError, "INTERNAL_CONTEXT_MISSING", "User context not found")
		return
	}

	// Get listing ID from URL parameter
	listingIDStr := c.Param("id")
	if listingIDStr == "" {
		httperrors.SendHTTPError(c, http.StatusBadRequest, "MISSING_ID", "Listing ID is required")
		return
	}

	listingID, err := strconv.ParseInt(listingIDStr, 10, 64)
	if err != nil {
		httperrors.SendHTTPError(c, http.StatusBadRequest, "INVALID_ID", "Invalid listing ID")
		return
	}

	// Call service to delete listing
	err = lh.listingService.DeleteListing(ctx, listingID)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	// Success response
	c.JSON(http.StatusOK, dto.DeleteListingResponse{
		Success: true,
		Message: "Listing deleted successfully",
	})
}
