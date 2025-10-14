package adminhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	dto "github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

// DeleteListingCatalogValue handles DELETE /admin/listing/catalog.
//
//	@Summary	Deactivate a listing catalog value
//	@Tags		Admin
//	@Accept		json
//	@Produce	json
//	@Param		request	body	dto.ListingCatalogDeleteRequest	true	"Deactivation payload"
//	@Success	204	"Catalog value deactivated"
//	@Failure	400	{object}	map[string]any
//	@Failure	401	{object}	map[string]any
//	@Failure	403	{object}	map[string]any
//	@Failure	404	{object}	map[string]any
//	@Failure	500	{object}	map[string]any
//	@Router		/admin/listing/catalog [delete]
func (h *AdminHandler) DeleteListingCatalogValue(c *gin.Context) {
	ctx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)
	var req dto.ListingCatalogDeleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httperrors.SendHTTPErrorObj(c, httperrors.ConvertBindError(err))
		return
	}

	if err := h.listingService.DeleteCatalogValue(ctx, req.Category, req.ID); err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
