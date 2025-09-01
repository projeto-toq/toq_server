package listinghandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/giulio-alfieri/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/giulio-alfieri/toq_server/internal/adapter/left/http/http_errors"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

// GetBaseFeatures handles getting available base features for listings
//
//	@Summary		Get base features
//	@Description	Get all available base features that can be associated with listings
//	@Tags			Listings
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	dto.GetBaseFeaturesResponse
//	@Failure		401	{object}	dto.ErrorResponse	"Unauthorized"
//	@Failure		403	{object}	dto.ErrorResponse	"Forbidden"
//	@Failure		500	{object}	dto.ErrorResponse	"Internal server error"
//	@Router			/listings/features/base [get]
//	@Security		BearerAuth
func (lh *ListingHandler) GetBaseFeatures(c *gin.Context) {
	ctx, spanEnd, err := utils.GenerateTracer(c.Request.Context())
	if err != nil {
		httperrors.SendHTTPError(c, http.StatusInternalServerError, "TRACER_ERROR", "Failed to generate tracer")
		return
	}
	defer spanEnd()

	// Call service to get base features
	features, err := lh.listingService.GetBaseFeatures(ctx)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, err)
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
