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

// RestoreListingCatalogValue handles POST /admin/listing/catalog/restore.
//
//	@Summary    Reactivate a listing catalog value
//	@Tags       Admin
//	@Accept     json
//	@Produce    json
//	@Param      request body dto.ListingCatalogRestoreRequest true "Reactivation payload"
//	@Success    200 {object} dto.ListingCatalogValueResponse
//	@Failure    400 {object} map[string]any
//	@Failure    401 {object} map[string]any
//	@Failure    403 {object} map[string]any
//	@Failure    404 {object} map[string]any
//	@Failure    409 {object} map[string]any
//	@Failure    500 {object} map[string]any
//	@Router     /admin/listing/catalog/restore [post]
func (h *AdminHandler) RestoreListingCatalogValue(c *gin.Context) {
	ctx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)
	var req dto.ListingCatalogRestoreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httperrors.SendHTTPErrorObj(c, httperrors.ConvertBindError(err))
		return
	}

	input := listingservices.RestoreCatalogValueInput{
		Category: req.Category,
		ID:       req.ID,
	}

	value, err := h.listingService.RestoreCatalogValue(ctx, input)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	c.JSON(http.StatusOK, httpconv.ToListingCatalogValueResponse(value))
}
