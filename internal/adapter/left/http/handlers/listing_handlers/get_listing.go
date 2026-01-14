package listinghandlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/converters"
	dto "github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/middlewares"
	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"
	listingservices "github.com/projeto-toq/toq_server/internal/core/service/listing_service"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

const ownerDetailPhotoVariant = "medium"

// GetListing retrieves comprehensive details of a listing by its identity ID
//
// This handler returns the ACTIVE version of a listing (referenced by listing_identities.active_version_id).
// If a draft version exists, its metadata is included in the response via draftVersionId field.
//
// Visibility rules:
//   - Owners: can access only their own listings (listing_identities.user_id == authenticated user_id)
//   - Realtors: can access listings from any owner only when the listing status is PUBLISHED
//
// Returned data includes:
//   - All listing version fields (address, property type, transaction details, prices, etc.)
//   - Enriched catalog values (owner, delivered, whoLives, transaction, etc.) with slug and label
//   - Features with descriptions and quantities
//   - Exchange places (if exchange is enabled)
//   - Financing blockers (if financing is disabled)
//   - Guarantees with priority (for rent transactions)
//   - Photo session booking ID (if active booking exists)
//   - Version metadata (activeVersionId, draftVersionId, version number, status)
//
// Business Rules:
//
//   - Returns HTTP 403 Forbidden if requester is not the owner (when role is owner) or if realtor requests a non-published listing
//
//   - Returns HTTP 404 Not Found if listing identity does not exist
//
//   - Returns HTTP 400 Bad Request if listingIdentityId is invalid or missing
//
//     @Summary		Get listing details by identity ID
//     @Description	Retrieves comprehensive details of a listing including active version, draft metadata (if exists),
//     enriched catalog values, features, guarantees, exchange places, and photo session status.
//     Only the listing owner can access details (ownership validated via listing_identities.user_id).
//     Returns the ACTIVE version by default (listing_identities.active_version_id).
//     @Tags			Listings
//     @Accept			json
//     @Produce		json
//     @Security		BearerAuth
//     @Param			Authorization	header	string						true	"Bearer token for authentication"	Extensions(x-example=Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...)
//     @Param			request			body	dto.GetListingDetailRequest	true	"Listing identity identifier"
//     @Success		200				{object}	dto.ListingDetailResponse	"Listing details successfully retrieved"
//     @Failure		400				{object}	dto.ErrorResponse			"Invalid request format (missing listingIdentityId or invalid value)"
//     @Failure		401				{object}	dto.ErrorResponse			"Unauthorized (missing or invalid token)"
//     @Failure		403				{object}	dto.ErrorResponse			"Forbidden (requester is not the listing owner)"
//     @Failure		404				{object}	dto.ErrorResponse			"Listing identity not found"
//     @Failure		500				{object}	dto.ErrorResponse			"Internal server error (database failure, transaction error)"
//     @Router			/listings/detail [post]
func (lh *ListingHandler) GetListing(c *gin.Context) {
	// Note: request tracing already provided by TelemetryMiddleware
	ctx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)

	// Validate authenticated user context (set by AuthMiddleware)
	if _, ok := middlewares.GetUserInfoFromContext(c); !ok {
		httperrors.SendHTTPError(c, http.StatusInternalServerError, "INTERNAL_CONTEXT_MISSING", "User context not found")
		return
	}

	// Parse and validate request body
	var req dto.GetListingDetailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httperrors.SendHTTPError(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request format")
		return
	}

	// Validate listingIdentityId (must be positive integer)
	if req.ListingIdentityID <= 0 {
		httperrors.SendHTTPErrorObj(c, coreutils.ValidationError("listingIdentityId", "listingIdentityId must be greater than zero"))
		return
	}

	// Call service layer to retrieve listing details
	// Service validates ownership and fetches active version with enriched catalog data
	detail, serviceErr := lh.listingService.GetListingDetail(ctx, req.ListingIdentityID)
	if serviceErr != nil {
		// SendHTTPErrorObj converts domain errors to appropriate HTTP responses:
		// - 403 if not owner
		// - 404 if listing not found
		// - 500 for infrastructure errors
		httperrors.SendHTTPErrorObj(c, serviceErr)
		return
	}

	// Enrich owner metadata with photo URL before serialization
	lh.attachOwnerPhoto(ctx, &detail)

	// Convert service output to DTO response
	response := converters.ListingDetailToDTO(detail)
	c.JSON(http.StatusOK, response)
}

// attachOwnerPhoto enriches ListingDetailOutput with owner photo download URL using userService impersonation.
func (lh *ListingHandler) attachOwnerPhoto(ctx context.Context, detail *listingservices.ListingDetailOutput) {
	if detail == nil || detail.OwnerDetail == nil || lh.userService == nil {
		return
	}
	ownerID := detail.OwnerDetail.ID
	if ownerID <= 0 {
		return
	}

	impersonatedCtx := coreutils.SetUserInContext(ctx, usermodel.UserInfos{ID: ownerID})
	photoURL, err := lh.userService.GetPhotoDownloadURL(impersonatedCtx, ownerDetailPhotoVariant)
	if err != nil {
		logger := coreutils.LoggerFromContext(ctx)
		logger.Debug("listing.detail.owner_photo_error", "owner_id", ownerID, "err", err)
		return
	}

	detail.OwnerDetail.PhotoURL = photoURL
}
