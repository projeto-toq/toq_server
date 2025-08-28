package listinghandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/giulio-alfieri/toq_server/internal/adapter/left/http/dto"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

// GetBaseFeatures handles getting available base features for listings
func (lh *ListingHandler) GetBaseFeatures(c *gin.Context) {
	ctx, spanEnd, err := utils.GenerateTracer(c.Request.Context())
	if err != nil {
		utils.SendHTTPError(c, http.StatusInternalServerError, "TRACER_ERROR", "Failed to generate tracer")
		return
	}
	defer spanEnd()

	// Call service to get base features
	features, err := lh.listingService.GetBaseFeatures(ctx)
	if err != nil {
		utils.SendHTTPError(c, http.StatusInternalServerError, "GET_BASE_FEATURES_FAILED", "Failed to get base features")
		return
	}

	// Convert to response DTOs
	baseFeatures := make([]dto.BaseFeature, 0, len(features))
	for _, feature := range features {
		baseFeatures = append(baseFeatures, dto.BaseFeature{
			ID:          int(feature.ID()),
			Name:        feature.Feature(),
			Description: feature.Description(),
			Category:    "", // Category not available in the model
		})
	}

	// Success response
	c.JSON(http.StatusOK, dto.GetBaseFeaturesResponse{
		Features: baseFeatures,
	})
}
