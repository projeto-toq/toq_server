package listinghandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/converters"
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/middlewares"
	listingservices "github.com/projeto-toq/toq_server/internal/core/service/listing_service"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

// ListListingVersions retrieves all versions for a specific listing identity
//
// @Summary     List listing versions by identity ID
// @Description Retrieves all versions attached to a listing identity, indicating which one is currently active.
//
//	Returns version metadata including ID, version number, status, and title.
//	Soft-deleted versions can be included by setting includeDeleted to true.
//
// @Tags        Listings
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       Authorization header string                           true  "Bearer token for authentication" Extensions(x-example=Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...)
// @Param       request       body   dto.ListListingVersionsRequest  true  "Listing identity ID and filter options"
// @Success     200           {object} dto.ListListingVersionsResponse "List of versions with active flag"
// @Failure     400           {object} dto.ErrorResponse               "Invalid request body or validation error"
// @Failure     401           {object} dto.ErrorResponse               "Unauthorized (missing or invalid token)"
// @Failure     403           {object} dto.ErrorResponse               "Forbidden (insufficient permissions)"
// @Failure     404           {object} dto.ErrorResponse               "Listing identity not found"
// @Failure     500           {object} dto.ErrorResponse               "Internal server error"
// @Router      /listings/versions [post]
func (lh *ListingHandler) ListListingVersions(c *gin.Context) {
	// Note: tracing already provided by TelemetryMiddleware; avoid duplicate spans
	baseCtx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)
	ctx, spanEnd, err := coreutils.GenerateTracer(baseCtx)
	if err != nil {
		httperrors.SendHTTPError(c, http.StatusInternalServerError, "TRACER_ERROR", "Failed to generate tracer")
		return
	}
	defer spanEnd()

	// Validate user context from JWT token (required for authorization)
	if _, ok := middlewares.GetUserInfoFromContext(c); !ok {
		httperrors.SendHTTPError(c, http.StatusInternalServerError, "INTERNAL_CONTEXT_MISSING", "User context not found")
		return
	}

	// Parse and validate JSON request body
	var request dto.ListListingVersionsRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httperrors.SendHTTPError(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request format")
		return
	}

	// Map DTO to service input (no transformation needed)
	input := listingservices.ListListingVersionsInput{
		ListingIdentityID: request.ListingIdentityID,
		IncludeDeleted:    request.IncludeDeleted,
	}

	// Call service layer for business logic execution
	output, serviceErr := lh.listingService.ListListingVersions(ctx, input)
	if serviceErr != nil {
		// Service layer error already logged; convert to HTTP response
		httperrors.SendHTTPErrorObj(c, serviceErr)
		return
	}

	// Convert service output to HTTP DTO
	response := converters.ListingVersionsToDTO(output)
	c.JSON(http.StatusOK, response)
}
