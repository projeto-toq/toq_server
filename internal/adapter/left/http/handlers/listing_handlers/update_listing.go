package listinghandlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/middlewares"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

// UpdateListing handles updating an existing listing
func (lh *ListingHandler) UpdateListing(c *gin.Context) {
	baseCtx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)
	_, spanEnd, err := coreutils.GenerateTracer(baseCtx)
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

	_, err = strconv.ParseInt(listingIDStr, 10, 64)
	if err != nil {
		httperrors.SendHTTPError(c, http.StatusBadRequest, "INVALID_ID", "Invalid listing ID")
		return
	}

	// Parse request body using DTO
	var request dto.UpdateListingRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httperrors.SendHTTPError(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request format")
		return
	}

	// Note: The service UpdateListing expects a ListingInterface, not individual fields
	// This would require first getting the listing, then updating it, then calling the service
	// For now, we'll return a not implemented response since the service signature doesn't match our needs
	httperrors.SendHTTPError(c, http.StatusNotImplemented, "UPDATE_NOT_IMPLEMENTED", "Update listing service needs refactoring for HTTP usage")
}
