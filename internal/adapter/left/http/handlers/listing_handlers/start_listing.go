package listinghandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/giulio-alfieri/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/giulio-alfieri/toq_server/internal/adapter/left/http/http_errors"
	"github.com/giulio-alfieri/toq_server/internal/adapter/left/http/middlewares"
	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

// StartListing handles creating a new listing
//
//	@Summary		Start a new listing
//	@Description	Create a new listing with basic information (zip code, number, property type)
//	@Tags			Listings
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.StartListingRequest	true	"Listing creation data"
//	@Success		201		{object}	dto.StartListingResponse
//	@Failure		400		{object}	dto.ErrorResponse	"Invalid request format"
//	@Failure		401		{object}	dto.ErrorResponse	"Unauthorized"
//	@Failure		403		{object}	dto.ErrorResponse	"Forbidden"
//	@Failure		500		{object}	dto.ErrorResponse	"Internal server error"
//	@Router			/listings [post]
//	@Security		BearerAuth
func (lh *ListingHandler) StartListing(c *gin.Context) {
	ctx, spanEnd, err := utils.GenerateTracer(c.Request.Context())
	if err != nil {
		httperrors.SendHTTPError(c, http.StatusInternalServerError, "TRACER_ERROR", "Failed to generate tracer")
		return
	}
	defer spanEnd()

	// Get user info from context (set by auth middleware)
	if _, ok := middlewares.GetUserInfoFromContext(c); !ok {
		// Se chegar aqui, é erro de pipeline (middleware deveria ter setado)
		httperrors.SendHTTPError(c, http.StatusInternalServerError, "INTERNAL_CONTEXT_MISSING", "User context not found")
		return
	}

	// Parse request body using DTO
	var request dto.StartListingRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httperrors.SendHTTPError(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request format")
		return
	}

	// Call service to start listing
	listing, err := lh.listingService.StartListing(
		ctx,
		request.ZipCode,
		request.Number,
		globalmodel.PropertyType(request.PropertyType),
	)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	// Success response
	c.JSON(http.StatusCreated, dto.StartListingResponse{
		ID: listing.ID(),
	})
}
