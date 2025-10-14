package listinghandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	httpconv "github.com/projeto-toq/toq_server/internal/adapter/left/http/converters"
	dto "github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

// ListCatalogValues handles GET /listings/catalog
//
//	@Summary	Listar valores ativos do catálogo de listings
//	@Tags		Listings
//	@Accept		json
//	@Produce	json
//	@Param		request	body	dto.ListingCatalogRequest	true	"Payload de consulta"
//	@Success	200	{object}	dto.ListingCatalogValuesResponse
//	@Failure	400	{object}	map[string]any
//	@Failure	401	{object}	map[string]any
//	@Failure	403	{object}	map[string]any
//	@Failure	500	{object}	map[string]any
//	@Router		/listings/catalog [get]
//	@Security	BearerAuth
func (lh *ListingHandler) ListCatalogValues(c *gin.Context) {
	ctx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)

	var req dto.ListingCatalogRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httperrors.SendHTTPErrorObj(c, httperrors.ConvertBindError(err))
		return
	}

	values, err := lh.listingService.ListCatalogValues(ctx, req.Category, false)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	c.JSON(http.StatusOK, httpconv.ToListingCatalogValuesResponse(values))
}
