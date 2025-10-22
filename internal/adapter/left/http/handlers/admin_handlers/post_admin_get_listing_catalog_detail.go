package adminhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	httpconv "github.com/projeto-toq/toq_server/internal/adapter/left/http/converters"
	dto "github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

// PostAdminGetListingCatalogDetail handles POST /admin/listing/catalog/detail
//
//	@Summary	Get listing catalog value detail
//	@Tags		Admin Listings
//	@Accept		json
//	@Produce	json
//	@Param		request	body	dto.AdminGetListingCatalogDetailRequest	true	"Listing catalog detail payload"
//	@Success	200	{object}	dto.ListingCatalogValueResponse
//	@Failure	400	{object}	map[string]any
//	@Failure	401	{object}	map[string]any
//	@Failure	403	{object}	map[string]any
//	@Failure	404	{object}	map[string]any
//	@Failure	500	{object}	map[string]any
//	@Router		/admin/listing/catalog/detail [post]
func (h *AdminHandler) PostAdminGetListingCatalogDetail(c *gin.Context) {
	ctx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)

	var req dto.AdminGetListingCatalogDetailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httperrors.SendHTTPErrorObj(c, httperrors.ConvertBindError(err))
		return
	}

	if req.ID > 255 {
		httperrors.SendHTTPErrorObj(c, coreutils.ValidationError("id", "id must be between 1 and 255"))
		return
	}

	value, err := h.listingService.GetCatalogValueDetail(ctx, req.Category, uint8(req.ID))
	if err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	c.JSON(http.StatusOK, httpconv.ToListingCatalogValueResponse(value))
}
