package listinghandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/middlewares"
	listingservices "github.com/projeto-toq/toq_server/internal/core/service/listing_service"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

// CreateDraftVersion handles creating a new draft version from an active listing
//
//	@Summary		Create draft version
//	@Description	Create a new draft version from the current active listing version
//	@Tags			Listings
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.CreateDraftVersionRequest	true	"Draft version creation data"
//	@Success		201		{object}	dto.CreateDraftVersionResponse
//	@Failure		400		{object}	dto.ErrorResponse	"Invalid request format"
//	@Failure		401		{object}	dto.ErrorResponse	"Unauthorized"
//	@Failure		403		{object}	dto.ErrorResponse	"Forbidden"
//	@Failure		404		{object}	dto.ErrorResponse	"Active version not found"
//	@Failure		409		{object}	dto.ErrorResponse	"Draft already exists or listing is published"
//	@Failure		410		{object}	dto.ErrorResponse	"Listing is permanently closed"
//	@Failure		423		{object}	dto.ErrorResponse	"Listing is locked in workflow"
//	@Failure		500		{object}	dto.ErrorResponse	"Internal server error"
//	@Router			/listings/versions/draft [post]
//	@Security		BearerAuth
func (lh *ListingHandler) CreateDraftVersion(c *gin.Context) {
	baseCtx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)
	ctx, spanEnd, err := coreutils.GenerateTracer(baseCtx)
	if err != nil {
		httperrors.SendHTTPError(c, http.StatusInternalServerError, "TRACER_ERROR", "Failed to generate tracer")
		return
	}
	defer spanEnd()

	// Get user info from context (set by auth middleware)
	if _, ok := middlewares.GetUserInfoFromContext(c); !ok {
		httperrors.SendHTTPError(c, http.StatusInternalServerError, "INTERNAL_CONTEXT_MISSING", "User context not found")
		return
	}

	// Parse request body using DTO
	var request dto.CreateDraftVersionRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httperrors.SendHTTPError(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request format")
		return
	}

	input := listingservices.CreateDraftVersionInput{
		ListingIdentityID: request.ListingIdentityID,
	}

	output, err := lh.listingService.CreateDraftVersion(ctx, input)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	// Success response
	c.JSON(http.StatusCreated, dto.CreateDraftVersionResponse{
		VersionID: output.VersionID,
		Version:   output.Version,
		Status:    output.Status,
	})
}
