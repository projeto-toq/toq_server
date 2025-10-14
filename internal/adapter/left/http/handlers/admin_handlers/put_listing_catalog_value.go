package adminhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	httpconv "github.com/projeto-toq/toq_server/internal/adapter/left/http/converters"
	dto "github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	listingservices "github.com/projeto-toq/toq_server/internal/core/service/listing_service"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

// UpdateListingCatalogValue handles PUT /admin/listing/catalog.
//
//	@Summary	Update a listing catalog value
//	@Tags		Admin
//	@Accept		json
//	@Produce	json
//	@Param		request	body	dto.ListingCatalogUpdateRequest	true	"Partial update payload"
//	@Success	200	{object}	dto.ListingCatalogValueResponse
//	@Failure	400	{object}	map[string]any
//	@Failure	401	{object}	map[string]any
//	@Failure	403	{object}	map[string]any
//	@Failure	404	{object}	map[string]any
//	@Failure	409	{object}	map[string]any
//	@Failure	500	{object}	map[string]any
//	@Router		/admin/listing/catalog [put]
func (h *AdminHandler) UpdateListingCatalogValue(c *gin.Context) {
	ctx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)
	var req dto.ListingCatalogUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httperrors.SendHTTPErrorObj(c, httperrors.ConvertBindError(err))
		return
	}

	input := listingservices.UpdateCatalogValueInput{
		Category:    req.Category,
		ID:          req.ID,
		Slug:        req.Slug,
		Label:       req.Label,
		Description: req.Description,
		IsActive:    req.IsActive,
	}

	value, err := h.listingService.UpdateCatalogValue(ctx, input)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	c.JSON(http.StatusOK, httpconv.ToListingCatalogValueResponse(value))
}
