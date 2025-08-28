package listinghandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/giulio-alfieri/toq_server/internal/adapter/left/http/dto"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

// GetAllListings handles getting all listings with pagination and filters
// Note: Service method not fully implemented yet, returning 501
//
//	@Summary		Get all listings with pagination and filters
//	@Description	Get all listings with optional filters for status, location, price range, etc.
//	@Tags			Listings
//	@Accept			json
//	@Produce		json
//	@Param			page		query		int		false	"Page number (default: 1)"			minimum(1)
//	@Param			limit		query		int		false	"Items per page (default: 10)"		minimum(1)	maximum(100)
//	@Param			status		query		string	false	"Filter by listing status"
//	@Param			userId		query		string	false	"Filter by user ID"
//	@Param			zipCode		query		string	false	"Filter by zip code"
//	@Param			minPrice	query		int		false	"Minimum price filter"
//	@Param			maxPrice	query		int		false	"Maximum price filter"
//	@Success		200			{object}	dto.GetAllListingsResponse
//	@Failure		400			{object}	dto.ErrorResponse	"Invalid request parameters"
//	@Failure		401			{object}	dto.ErrorResponse	"Unauthorized"
//	@Failure		403			{object}	dto.ErrorResponse	"Forbidden"
//	@Failure		500			{object}	dto.ErrorResponse	"Internal server error"
//	@Failure		501			{object}	dto.ErrorResponse	"Service method not implemented"
//	@Router			/listings [get]
//	@Security		BearerAuth
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
