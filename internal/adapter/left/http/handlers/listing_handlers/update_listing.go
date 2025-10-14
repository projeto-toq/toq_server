package listinghandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/converters"
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/middlewares"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

// UpdateListing updates an existing listing.
//
//	@Summary	Update a listing
//	@Description	Allows partial updates for draft listings. Omitted fields remain unchanged; present fields (including null/empty) overwrite stored values.
//	@Tags		Listings
//	@Accept		json
//	@Produce	json
//	@Param		request	body	dto.UpdateListingRequest	true	"Payload for update (ID must be provided in the body)" Extensions(x-example={"id":98765,"owner":2,"features":[{"featureId":101,"quantity":2},{"featureId":205,"quantity":1}],"landSize":423.5,"corner":true,"nonBuildable":12.75,"buildable":410.75,"delivered":1,"whoLives":3,"description":"Apartamento amplo com vista panoramica","transaction":2,"sellNet":1200000,"rentNet":8500,"condominium":1200.5,"annualTax":3400.75,"annualGroundRent":1800,"exchange":true,"exchangePercentual":50,"exchangePlaces":[{"neighborhood":"Vila Mariana","city":"Sao Paulo","state":"SP"},{"neighborhood":"Centro","city":"Campinas","state":"SP"}],"installment":2,"financing":true,"financingBlockers":[4,7],"guarantees":[{"priority":1,"guarantee":33},{"priority":2,"guarantee":34}],"visit":3,"tenantName":"Joao da Silva","tenantEmail":"joao.silva@example.com","tenantPhone":"+55 11 91234-5678","accompanying":2})
//	@Success	200	{object}	dto.UpdateListingResponse
//	@Failure	400	{object}	dto.ErrorResponse	"Invalid payload"
//	@Failure	401	{object}	dto.ErrorResponse	"Unauthorized"
//	@Failure	403	{object}	dto.ErrorResponse	"Forbidden"
//	@Failure	404	{object}	dto.ErrorResponse	"Not found"
//	@Failure	409	{object}	dto.ErrorResponse	"Conflict"
//	@Failure	500	{object}	dto.ErrorResponse	"Internal error"
//	@Router		/listings [put]
//	@Security	BearerAuth
func (lh *ListingHandler) UpdateListing(c *gin.Context) {
	baseCtx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)

	if _, ok := middlewares.GetUserInfoFromContext(c); !ok {
		httperrors.SendHTTPError(c, http.StatusInternalServerError, "INTERNAL_CONTEXT_MISSING", "User context not found")
		return
	}

	var request dto.UpdateListingRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httperrors.SendHTTPError(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request format")
		return
	}

	if !request.ID.IsPresent() || request.ID.IsNull() {
		httperrors.SendHTTPError(c, http.StatusBadRequest, "MISSING_ID", "Listing ID must be provided in the request body")
		return
	}

	listingID, ok := request.ID.Value()
	if !ok || listingID <= 0 {
		httperrors.SendHTTPError(c, http.StatusBadRequest, "INVALID_ID", "Listing ID is invalid")
		return
	}

	input, err := converters.UpdateListingRequestToInput(request)
	if err != nil {
		httperrors.SendHTTPError(c, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}
	input.ID = listingID

	if err := lh.listingService.UpdateListing(baseCtx, input); err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.UpdateListingResponse{
		Success: true,
		Message: "Listing updated",
	})
}
