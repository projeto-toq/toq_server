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
//	@Summary	List listing catalog values
//	@Description	Available categories: property_owner, property_delivered, who_lives, transaction_type, installment_plan, financing_blocker, visit_type, accompanying_type, guarantee_type.
//	@Tags		Admin
//	@Produce	json
//	@Param		category	query	string	true	"Catalog category"
//	@Param		includeInactive	query	bool	false	"Include inactive values"
//	@Success	200	{object}	dto.ListingCatalogValuesResponse
//	@Failure	400	{object}	map[string]any
//	@Failure	401	{object}	map[string]any
//	@Failure	403	{object}	map[string]any
//	@Failure	500	{object}	map[string]any
//	@Router		/admin/listing/catalog [get]
func (h *AdminHandler) ListListingCatalogValues(c *gin.Context) {
	ctx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)
	var req dto.AdminListingCatalogRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		httperrors.SendHTTPErrorObj(c, httperrors.ConvertBindError(err))
		return
	}

	values, err := h.listingService.ListCatalogValues(ctx, req.Category, req.IncludeInactive)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	c.JSON(http.StatusOK, httpconv.ToListingCatalogValuesResponse(values))
}
