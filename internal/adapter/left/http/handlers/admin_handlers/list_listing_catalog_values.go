package adminhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	httpconv "github.com/projeto-toq/toq_server/internal/adapter/left/http/converters"
	dto "github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

// ListListingCatalogValues handles GET /admin/listing/catalog.
//
// This endpoint retrieves all values for a specific listing catalog category.
// It supports filtering by active status and is used primarily by the admin panel
// to manage catalog entries.
//
//	@Summary		List listing catalog values
//	@Description	Available categories: property_owner, property_delivered, who_lives, transaction_type, installment_plan, financing_blocker, visit_type, accompanying_type, guarantee_type, land_terrain_type, warehouse_sector.
//	@Tags			Admin Listings
//	@Produce		json
//	@Param			category		query		string	true	"Catalog category (property_owner, property_delivered, who_lives, transaction_type, installment_plan, financing_blocker, visit_type, accompanying_type, guarantee_type, land_terrain_type, warehouse_sector)"	Extensions(x-example=land_terrain_type)
//	@Param			includeInactive	query		bool	false	"Include inactive values"	default(false)
//	@Success		200	{object}	dto.ListingCatalogValuesResponse	"List of catalog values"
//	@Failure		400	{object}	map[string]any						"Invalid category or request parameters"
//	@Failure		401	{object}	map[string]any						"Unauthorized (missing or invalid token)"
//	@Failure		403	{object}	map[string]any						"Forbidden (insufficient permissions)"
//	@Failure		500	{object}	map[string]any						"Internal server error"
//	@Router			/admin/listing/catalog [get]
func (h *AdminHandler) ListListingCatalogValues(c *gin.Context) {
	// Note: request tracing already provided by TelemetryMiddleware
	ctx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)

	// Parse and validate query parameters
	var req dto.AdminListingCatalogRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		httperrors.SendHTTPErrorObj(c, httperrors.ConvertBindError(err))
		return
	}

	// Call service layer for catalog values retrieval
	// Service validates category and applies filtering
	values, err := h.listingService.ListCatalogValues(ctx, req.Category, req.IncludeInactive)
	if err != nil {
		// Service returns domain errors (validation, infra failures)
		// SendHTTPErrorObj converts to appropriate HTTP status
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	// Convert domain models to DTO response
	c.JSON(http.StatusOK, httpconv.ToListingCatalogValuesResponse(values))
}
