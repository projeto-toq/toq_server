package listinghandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	dto "github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

// RemoveFavoriteListing removes a listing from the authenticated user's favorites.
//
// @Summary     Unfavorite a listing
// @Description Removes the favorite relation for the given listing identity. Operation is idempotent.
// @Tags        Listings
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       Authorization header string true "Bearer token" Extensions(x-example=Bearer <token>)
// @Param       request body dto.FavoriteListingRequest true "Favorite request"
// @Success     204 "Favorite removed"
// @Failure     400 {object} dto.ErrorResponse "Invalid payload"
// @Failure     401 {object} dto.ErrorResponse "Unauthorized"
// @Failure     500 {object} dto.ErrorResponse "Internal error"
// @Router      /listings/favorites [delete]
func (lh *ListingHandler) RemoveFavoriteListing(c *gin.Context) {
	ctx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)

	var req dto.FavoriteListingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httperrors.SendHTTPErrorObj(c, httperrors.ConvertBindError(err))
		return
	}

	if req.ListingIdentityID <= 0 {
		httperrors.SendHTTPErrorObj(c, coreutils.ValidationError("listingIdentityId", "listingIdentityId must be greater than zero"))
		return
	}

	if err := lh.listingService.RemoveFavoriteListing(ctx, req.ListingIdentityID); err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
