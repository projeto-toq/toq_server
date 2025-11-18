package listinghandlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/middlewares"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

// GetListingByUserId handles getting all listings for a specific user
func (lh *ListingHandler) GetListingByUserId(c *gin.Context) {
	baseCtx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)
	ctx, spanEnd, err := coreutils.GenerateTracer(baseCtx)
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
		price := listing.SellNet()
		if price == 0 {
			price = listing.RentNet()
		}

		var draftVersionID *int64
		if draft, ok := listing.DraftVersion(); ok && draft != nil {
			if draftID := draft.ID(); draftID > 0 {
				draftVersionID = &draftID
			}
		}

		activeVersionID := listing.ActiveVersionID()
		if activeVersionID == 0 {
			activeVersionID = listing.ID()
		}

		complexValue := ""
		if listing.HasComplex() {
			complexValue = strings.TrimSpace(listing.Complex())
		}

		listingResponses = append(listingResponses, dto.ListingResponse{
			ID:                listing.ID(),
			ListingIdentityID: listing.IdentityID(),
			ListingUUID:       listing.UUID(),
			ActiveVersionID:   activeVersionID,
			DraftVersionID:    draftVersionID,
			Version:           listing.Version(),
			Title:             strings.TrimSpace(listing.Title()),
			Description:       listing.Description(),
			Price:             price,
			Status:            listing.Status().String(),
			PropertyType:      int(listing.ListingType()),
			ZipCode:           listing.ZipCode(),
			Number:            listing.Number(),
			Complex:           complexValue,
			UserID:            listing.UserID(),
			ComplexID:         "", // ComplexID not easily accessible
		})
	}

	// Success response
	c.JSON(http.StatusOK, gin.H{
		"data": listingResponses,
	})
}
