package listinghandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/giulio-alfieri/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/giulio-alfieri/toq_server/internal/adapter/left/http/http_errors"
	coreutils "github.com/giulio-alfieri/toq_server/internal/core/utils"
)

// GetOptions handles getting available options for listings
//
//	@Summary		Get listing options
//	@Description	Get available property types and statuses for listings based on location
//	@Tags			Listings
//	@Accept			json
//	@Produce		json
//	@Param			zipCode	query		string					true	"Zip code for location-based options"
//	@Param			number	query		string					true	"Property number"
//	@Success		200		{object}	dto.GetOptionsResponse
//	@Failure		400		{object}	dto.ErrorResponse	"Missing required parameters"
//	@Failure		401		{object}	dto.ErrorResponse	"Unauthorized"
//	@Failure		403		{object}	dto.ErrorResponse	"Forbidden"
//	@Failure		500		{object}	dto.ErrorResponse	"Internal server error"
//	@Router			/listings/options [get]
//	@Security		BearerAuth
func (lh *ListingHandler) GetOptions(c *gin.Context) {
	baseCtx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)
	ctx, spanEnd, err := coreutils.GenerateTracer(baseCtx)
	if err != nil {
		httperrors.SendHTTPError(c, http.StatusInternalServerError, "TRACER_ERROR", "Failed to generate tracer")
		return
	}
	defer spanEnd()

	// Get zipCode and number from query parameters
	zipCode := c.Query("zipCode")
	number := c.Query("number")

	if zipCode == "" || number == "" {
		httperrors.SendHTTPError(c, http.StatusBadRequest, "INVALID_REQUEST", "zipCode and number are required")
		return
	}

	// Call service to get options
	types, err := lh.listingService.GetOptions(ctx, zipCode, number)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	// Convert types to property type options
	propertyTypes := make([]dto.PropertyTypeOption, 0, len(types))
	for _, t := range types {
		propertyTypes = append(propertyTypes, dto.PropertyTypeOption{
			ID:   int(t),
			Name: getPropertyTypeName(int(t)), // Helper function to get name from ID
		})
	}

	// Static status options (these could be moved to a service later)
	statuses := []dto.StatusOption{
		{ID: "DRAFT", Name: "Rascunho"},
		{ID: "ACTIVE", Name: "Ativo"},
		{ID: "INACTIVE", Name: "Inativo"},
		{ID: "SOLD", Name: "Vendido"},
		{ID: "RENTED", Name: "Alugado"},
	}

	// Success response
	c.JSON(http.StatusOK, dto.GetOptionsResponse{
		PropertyTypes: propertyTypes,
		Statuses:      statuses,
	})
}

// Helper function to map property type ID to name
// This could be moved to a more appropriate place or fetched from a service
func getPropertyTypeName(id int) string {
	switch id {
	case 1:
		return "Casa"
	case 2:
		return "Apartamento"
	case 3:
		return "Terreno"
	case 4:
		return "Comercial"
	case 5:
		return "Rural"
	default:
		return "Outros"
	}
}
