package listinghandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	httpconv "github.com/projeto-toq/toq_server/internal/adapter/left/http/converters"
	dto "github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

var _ dto.ListingComplexItemResponse

// ListComplexes handles GET /listings/complexes
//
// @Summary     List all complexes (vertical and horizontal)
// @Description Returns every managed complex with only id and name fields. Standalone coverage is excluded.
// @Tags        Listings
// @Produce     json
// @Security    BearerAuth
// @Success     200 {array} dto.ListingComplexItemResponse "List of complexes"
// @Failure     500 {object} map[string]any
// @Router      /listings/complexes [get]
func (h *ListingHandler) ListComplexes(c *gin.Context) {
	// Tracing handled by middleware; enrich context for logging and request metadata.
	ctx := utils.EnrichContextWithRequestInfo(c.Request.Context(), c)

	complexes, err := h.propertyCoverageService.ListPublicComplexes(ctx)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	c.JSON(http.StatusOK, httpconv.ToListingComplexItems(complexes))
}
