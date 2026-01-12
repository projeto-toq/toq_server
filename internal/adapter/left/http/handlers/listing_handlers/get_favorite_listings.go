package listinghandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	dto "github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

// GetFavoriteListings returns the authenticated user's favorite listings with pagination.
//
// @Summary     List favorite listings
// @Description Retrieves the user's favorites ordered by most recent favorite action.
// @Tags        Listings
// @Produce     json
// @Security    BearerAuth
// @Param       Authorization header string true "Bearer token" Extensions(x-example=Bearer <token>)
// @Param       page  query int false "Page number (default 1)" minimum(1) default(1)
// @Param       limit query int false "Items per page (default 20, max 100)" minimum(1) maximum(100) default(20)
// @Success     200 {object} dto.ListListingsResponse "Favorite listings"
// @Failure     400 {object} dto.ErrorResponse "Invalid query params"
// @Failure     401 {object} dto.ErrorResponse "Unauthorized"
// @Failure     500 {object} dto.ErrorResponse "Internal error"
// @Router      /listings/favorites [get]
func (lh *ListingHandler) GetFavoriteListings(c *gin.Context) {
	ctx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)

	var req dto.ListFavoriteListingsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		httperrors.SendHTTPErrorObj(c, httperrors.ConvertBindError(err))
		return
	}

	inputPage := req.Page
	inputLimit := req.Limit

	result, err := lh.listingService.ListFavoriteListings(ctx, inputPage, inputLimit)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	data := make([]dto.ListingResponse, 0, len(result.Items))
	for _, item := range result.Items {
		data = append(data, toListingResponse(item))
	}

	resp := dto.ListListingsResponse{
		Data: data,
		Pagination: dto.PaginationResponse{
			Page:       result.Page,
			Limit:      result.Limit,
			Total:      result.Total,
			TotalPages: computeTotalPages(result.Total, result.Limit),
		},
	}

	c.JSON(http.StatusOK, resp)
}
