package listinghandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/giulio-alfieri/toq_server/internal/adapter/left/http/dto"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

// GetAllListings handles getting all listings with pagination and filters
// Note: Service method not fully implemented yet, returning 501
func (lh *ListingHandler) GetAllListings(c *gin.Context) {
	_, spanEnd, err := utils.GenerateTracer(c.Request.Context())
	if err != nil {
		utils.SendHTTPError(c, http.StatusInternalServerError, "TRACER_ERROR", "Failed to generate tracer")
		return
	}
	defer spanEnd()

	// Parse query parameters
	var request dto.GetAllListingsRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		utils.SendHTTPError(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid query parameters")
		return
	}

	// Service method GetAllListings not implemented yet
	utils.SendHTTPError(c, http.StatusNotImplemented, "NOT_IMPLEMENTED", "GetAllListings service method not implemented yet")
}
