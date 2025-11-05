package listinghandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/converters"
	dto "github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/middlewares"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

// GetListing trata a rota GET /listings/detail.
//
//	@Summary		Get listing detail
//	@Description	Returns all fields of a listing given its identifier, including active photo session ID if exists.
//	@Tags			Listings
//	@Accept			json
//	@Produce		json
//	@Param			request	body	dto.GetListingDetailRequest	true	"Listing identifier"
//	@Success		200	{object}	dto.ListingDetailResponse
//	@Failure		400	{object}	dto.ErrorResponse
//	@Failure		401	{object}	dto.ErrorResponse
//	@Failure		403	{object}	dto.ErrorResponse
//	@Failure		404	{object}	dto.ErrorResponse
//	@Failure		500	{object}	dto.ErrorResponse
//	@Router			/listings/detail [get]
//	@Security		BearerAuth
func (lh *ListingHandler) GetListing(c *gin.Context) {
	baseCtx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)
	ctx, spanEnd, err := coreutils.GenerateTracer(baseCtx)
	if err != nil {
		httperrors.SendHTTPError(c, http.StatusInternalServerError, "TRACER_ERROR", "Failed to generate tracer")
		return
	}
	defer spanEnd()

	if _, ok := middlewares.GetUserInfoFromContext(c); !ok {
		httperrors.SendHTTPError(c, http.StatusInternalServerError, "INTERNAL_CONTEXT_MISSING", "User context not found")
		return
	}

	var req dto.GetListingDetailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httperrors.SendHTTPError(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request format")
		return
	}

	if req.ListingID <= 0 {
		httperrors.SendHTTPErrorObj(c, coreutils.ValidationError("listingId", "listingId must be greater than zero"))
		return
	}

	detail, serviceErr := lh.listingService.GetListingDetail(ctx, req.ListingID)
	if serviceErr != nil {
		httperrors.SendHTTPErrorObj(c, serviceErr)
		return
	}

	response := converters.ListingDetailToDTO(detail)
	c.JSON(http.StatusOK, response)
}
