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
//	@Description	Allows partial updates for draft listings. Omitted fields remain unchanged; present fields (including null/empty) overwrite stored values. **IMPORTANT TAX RULES**: (1) IPTU (property tax) requires exactly ONE field: either `annualTax` OR `monthlyTax`, never both simultaneously. (2) Laudêmio (ground rent) is optional but if provided, use either `annualGroundRent` OR `monthlyGroundRent`, never both.
//	@Tags		Listings
//	@Accept		json
//	@Produce	json
//	@Param	request	body	dto.UpdateListingRequestSwagger	true	"Payload for update (listingIdentityId and listingVersionId are required)" Extensions(x-example={"listingIdentityId":1024,"listingVersionId":5001,"owner":"myself","features":[{"featureId":101,"quantity":2},{"featureId":205,"quantity":1}],"landSize":423.5,"corner":true,"nonBuildable":12.75,"buildable":410.75,"delivered":"furnished","whoLives":"tenant","description":"Apartamento amplo com vista panoramica","transaction":"sale","sellNet":1200000,"rentNet":8500,"condominium":1200.5,"monthlyTax":283.40,"monthlyGroundRent":150,"exchange":true,"exchangePercentual":50,"exchangePlaces":[{"neighborhood":"Vila Mariana","city":"Sao Paulo","state":"SP"},{"neighborhood":"Centro","city":"Campinas","state":"SP"}],"installment":"short_term","financing":true,"financingBlockers":["pending_probate","other"],"guarantees":[{"priority":1,"guarantee":"security_deposit"},{"priority":2,"guarantee":"surety_bond"}],"visit":"client","tenantName":"Joao da Silva","tenantEmail":"joao.silva@example.com","tenantPhone":"+5511912345678","title":"Apartamento 3 dormitorios com piscina","accompanying":"assistant","completionForecast":"2026-06","landBlock":"A","landLot":"15","landFront":12.5,"landSide":30.0,"landBack":12.5,"landTerrainType":"plano","hasKmz":true,"kmzFile":"https://storage.exemplo.com/terrenos/lote15.kmz","buildingFloors":8,"unitTower":"Torre B","unitFloor":5,"unitNumber":"502","warehouseManufacturingArea":850.5,"warehouseSector":"industrial","warehouseHasPrimaryCabin":true,"warehouseCabinKva":150.0,"warehouseGroundFloor":4.2,"warehouseFloorResistance":2500.0,"warehouseZoning":"ZI-2","warehouseHasOfficeArea":true,"warehouseOfficeArea":120.0,"warehouseAdditionalFloors":[{"floorName":"Mezanino","floorOrder":1,"floorHeight":3.5},{"floorName":"Segundo Piso","floorOrder":2,"floorHeight":3.2}],"storeHasMezzanine":true,"storeMezzanineArea":45.0})
//	@Success	200	{object}	dto.UpdateListingResponse
//	@Failure	400	{object}	dto.ErrorResponse	"Invalid payload"
//	@Failure	401	{object}	dto.ErrorResponse	"Unauthorized"
//	@Failure	403	{object}	dto.ErrorResponse	"Forbidden"
//	@Failure	404	{object}	dto.ErrorResponse	"Not found"
//	@Failure	409	{object}	dto.ErrorResponse	"Conflict"
//	@Failure	422	{object}	dto.ErrorResponse	"Tax field conflict: both annual and monthly values provided for IPTU or Laudêmio"
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

	// Validate listingIdentityId
	if !request.ListingIdentityID.IsPresent() || request.ListingIdentityID.IsNull() {
		httperrors.SendHTTPError(c, http.StatusBadRequest, "MISSING_IDENTITY_ID", "listingIdentityId must be provided in the request body")
		return
	}
	identityID, ok := request.ListingIdentityID.Value()
	if !ok || identityID <= 0 {
		httperrors.SendHTTPError(c, http.StatusBadRequest, "INVALID_IDENTITY_ID", "listingIdentityId is invalid")
		return
	}

	// Validate listingVersionId
	if !request.ListingVersionID.IsPresent() || request.ListingVersionID.IsNull() {
		httperrors.SendHTTPError(c, http.StatusBadRequest, "MISSING_VERSION_ID", "listingVersionId must be provided in the request body")
		return
	}

	versionID, ok := request.ListingVersionID.Value()
	if !ok || versionID <= 0 {
		httperrors.SendHTTPError(c, http.StatusBadRequest, "INVALID_VERSION_ID", "listingVersionId is invalid")
		return
	}

	input, err := converters.UpdateListingRequestToInput(request)
	if err != nil {
		httperrors.SendHTTPError(c, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}
	input.ListingIdentityID = identityID
	input.VersionID = versionID

	if err := lh.listingService.UpdateListing(baseCtx, input); err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.UpdateListingResponse{
		Success: true,
		Message: "Listing updated",
	})
}
