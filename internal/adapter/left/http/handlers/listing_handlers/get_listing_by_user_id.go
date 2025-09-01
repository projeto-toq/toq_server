package listinghandlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/giulio-alfieri/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/giulio-alfieri/toq_server/internal/adapter/left/http/http_errors"
	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

// GetListingByUserId handles getting all listings for a specific user
func (lh *ListingHandler) GetListingByUserId(c *gin.Context) {
	ctx, spanEnd, err := utils.GenerateTracer(c.Request.Context())
	if err != nil {
		httperrors.SendHTTPError(c, http.StatusInternalServerError, "TRACER_ERROR", "Failed to generate tracer")
		return
	}
	defer spanEnd()

	// Get user info from context (set by auth middleware)
	_, exists := c.Get(string(globalmodel.TokenKey))
	if !exists {
		httperrors.SendHTTPError(c, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	// Get user ID from URL parameter
	userIDStr := c.Param("userId")
	if userIDStr == "" {
		httperrors.SendHTTPError(c, http.StatusBadRequest, "MISSING_USER_ID", "User ID is required")
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		httperrors.SendHTTPError(c, http.StatusBadRequest, "INVALID_USER_ID", "Invalid user ID")
		return
	}

	// Call service to get listings by user
	listings, err := lh.listingService.GetAllListingsByUser(ctx, userID)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	// Convert to response DTOs
	listingResponses := make([]dto.ListingResponse, 0, len(listings))
	for _, listing := range listings {
		listingResponses = append(listingResponses, dto.ListingResponse{
			ID:           listing.ID(),
			Title:        "", // Title not available in basic listing model
			Description:  listing.Description(),
			Price:        0, // Price not available in basic listing model
			Status:       string(listing.Status()),
			PropertyType: int(listing.ListingType()),
			ZipCode:      listing.ZipCode(),
			Number:       listing.Number(),
			UserID:       listing.UserID(),
			ComplexID:    "", // ComplexID not easily accessible
			CreatedAt:    "", // CreatedAt not available in basic model
			UpdatedAt:    "", // UpdatedAt not available in basic model
		})
	}

	// Success response
	c.JSON(http.StatusOK, gin.H{
		"data": listingResponses,
	})
}
