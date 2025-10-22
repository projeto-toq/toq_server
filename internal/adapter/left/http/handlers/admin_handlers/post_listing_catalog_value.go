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

// CreateListingCatalogValue handles POST /admin/listing/catalog.
//
//	@Summary	Create a listing catalog value
//	@Tags		Admin Listings
//	@Accept		json
//	@Produce	json
//	@Param		request	body	dto.ListingCatalogCreateRequest	true	"Creation payload"
//	@Success	201	{object}	dto.ListingCatalogValueResponse
//	@Failure	400	{object}	map[string]any
//	@Failure	401	{object}	map[string]any
//	@Failure	403	{object}	map[string]any
//	@Failure	409	{object}	map[string]any
//	@Failure	500	{object}	map[string]any
//	@Router		/admin/listing/catalog [post]
func (h *AdminHandler) CreateListingCatalogValue(c *gin.Context) {
	ctx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)
	var req dto.ListingCatalogCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httperrors.SendHTTPErrorObj(c, httperrors.ConvertBindError(err))
		return
	}

	input := listingservices.CreateCatalogValueInput{
		Category:    req.Category,
		Slug:        req.Slug,
		Label:       req.Label,
		Description: req.Description,
		IsActive:    req.IsActive,
	}

	value, err := h.listingService.CreateCatalogValue(ctx, input)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	c.JSON(http.StatusCreated, httpconv.ToListingCatalogValueResponse(value))
}
