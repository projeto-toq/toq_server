package listinghandlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	dto "github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
	listingservices "github.com/projeto-toq/toq_server/internal/core/service/listing_service"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

// ListListings retrieves a paginated list of active listing versions with optional filters and sorting
//
//	@Summary      List active listing versions with filters and sorting
//	@Description  Retrieves a paginated list of active listing versions (versions linked via listing_identities.active_version_id).
//	              By default, only active versions are returned. Use includeAllVersions=true to retrieve all versions (active + draft).
//	              Supports filtering by status, code, title, location (zipCode, city, neighborhood), owner (userId),
//	              and price/size ranges. Results can be sorted by id (creation date proxy) or status.
//	              Default sorting: id DESC (newest first).
//	@Tags         Listings
//	@Produce      json
//	@Security     BearerAuth
//	@Param        Authorization       header  string  true   "Bearer token for authentication" example(Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...)
//	@Param        page                query   int     false  "Page number (1-indexed)" minimum(1) default(1) example(1)
//	@Param        limit               query   int     false  "Items per page" minimum(1) maximum(100) default(20) example(20)
//	@Param        sortBy              query   string  false  "Field to sort by" Enums(id, status) default(id) example(id)
//	@Param        sortOrder           query   string  false  "Sort direction" Enums(asc, desc) default(desc) example(desc)
//	@Param        status              query   string  false  "Filter by listing status (enum name or numeric)" example("PUBLISHED")
//	@Param        code                query   int     false  "Filter by exact listing code" example(1024)
//	@Param        title               query   string  false  "Filter by listing title/description (supports '*' wildcard)" example("*garden*")
//	@Param        userId              query   int     false  "Filter by owner user ID (owners auto-filtered to their own listings)" example(55)
//	@Param        zipCode             query   string  false  "Filter by zip code (digits only; supports '*' wildcard)" example("06543*")
//	@Param        city                query   string  false  "Filter by city (supports '*' wildcard)" example("*Paulista*")
//	@Param        neighborhood        query   string  false  "Filter by neighborhood (supports '*' wildcard)" example("*Centro*")
//	@Param        minSell             query   number  false  "Minimum sell price" example(100000)
//	@Param        maxSell             query   number  false  "Maximum sell price" example(900000)
//	@Param        minRent             query   number  false  "Minimum rent price" example(1500)
//	@Param        maxRent             query   number  false  "Maximum rent price" example(8000)
//	@Param        minLandSize         query   number  false  "Minimum land size in square meters" example(120.5)
//	@Param        maxLandSize         query   number  false  "Maximum land size in square meters" example(500.75)
//	@Param        includeAllVersions  query   bool    false  "Include all versions (active + draft). Default: false (active only)" example(false)
//	@Success      200                 {object}  dto.ListListingsResponse         "Paginated list of listings with metadata"
//	@Failure      400                 {object}  dto.ErrorResponse                "Invalid request parameters (malformed sortBy, sortOrder, or filter values)"
//	@Failure      401                 {object}  dto.ErrorResponse                "Unauthorized (missing or invalid token)"
//	@Failure      403                 {object}  dto.ErrorResponse                "Forbidden (user lacks permission to access this resource)"
//	@Failure      422                 {object}  dto.ErrorResponse                "Validation failed (invalid enum values, range errors)" example({"code":422,"message":"Validation failed","details":{"field":"sortBy","error":"Invalid sort field"}})
//	@Failure      500                 {object}  dto.ErrorResponse                "Internal server error"
//	@Router       /listings [get]
func (lh *ListingHandler) ListListings(c *gin.Context) {
	// Note: tracing already provided by TelemetryMiddleware
	ctx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)

	var req dto.ListListingsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		httperrors.SendHTTPErrorObj(c, httperrors.ConvertBindError(err))
		return
	}

	// Parse and validate sortBy (default: id)
	sortBy, err := parseSortBy(req.SortBy)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, coreutils.ValidationError("sortBy", err.Error()))
		return
	}

	// Parse and validate sortOrder (default: desc)
	sortOrder, err := parseSortOrder(req.SortOrder)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, coreutils.ValidationError("sortOrder", err.Error()))
		return
	}

	// Parse optional status filter
	statusPtr, err := parseListingStatus(strings.TrimSpace(req.Status))
	if err != nil {
		httperrors.SendHTTPErrorObj(c, coreutils.ValidationError("status", err.Error()))
		return
	}

	// Parse optional numeric filters with range validation
	codePtr, err := parseOptionalUint32(req.Code)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, coreutils.ValidationError("code", err.Error()))
		return
	}

	userIDPtr, err := parseOptionalInt64(req.UserID)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, coreutils.ValidationError("userId", err.Error()))
		return
	}

	// Parse and validate price range filters
	minSell, err := parseOptionalFloat64(req.MinSellPrice)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, coreutils.ValidationError("minSell", err.Error()))
		return
	}
	maxSell, err := parseOptionalFloat64(req.MaxSellPrice)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, coreutils.ValidationError("maxSell", err.Error()))
		return
	}
	if minSell != nil && maxSell != nil && *minSell > *maxSell {
		httperrors.SendHTTPErrorObj(c, coreutils.ValidationError("minSell", "minSell cannot be greater than maxSell"))
		return
	}

	minRent, err := parseOptionalFloat64(req.MinRentPrice)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, coreutils.ValidationError("minRent", err.Error()))
		return
	}
	maxRent, err := parseOptionalFloat64(req.MaxRentPrice)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, coreutils.ValidationError("maxRent", err.Error()))
		return
	}
	if minRent != nil && maxRent != nil && *minRent > *maxRent {
		httperrors.SendHTTPErrorObj(c, coreutils.ValidationError("minRent", "minRent cannot be greater than maxRent"))
		return
	}

	// Parse and validate land size range filters
	minLand, err := parseOptionalFloat64(req.MinLandSize)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, coreutils.ValidationError("minLandSize", err.Error()))
		return
	}
	maxLand, err := parseOptionalFloat64(req.MaxLandSize)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, coreutils.ValidationError("maxLandSize", err.Error()))
		return
	}
	if minLand != nil && maxLand != nil && *minLand > *maxLand {
		httperrors.SendHTTPErrorObj(c, coreutils.ValidationError("minLandSize", "minLandSize cannot be greater than maxLandSize"))
		return
	}

	// Extract authenticated user info for permission filtering
	userInfo, infoErr := coreutils.GetUserInfoFromGinContext(c)
	if infoErr != nil {
		httperrors.SendHTTPErrorObj(c, coreutils.AuthenticationError("User info not found"))
		return
	}

	// Build service input with all filters and sorting
	input := listingservices.ListListingsInput{
		Page:               req.Page,
		Limit:              req.Limit,
		SortBy:             sortBy,
		SortOrder:          sortOrder,
		Status:             statusPtr,
		Code:               codePtr,
		Title:              strings.TrimSpace(req.Title),
		ZipCode:            strings.TrimSpace(req.ZipCode),
		City:               strings.TrimSpace(req.City),
		Neighborhood:       strings.TrimSpace(req.Neighborhood),
		UserID:             userIDPtr,
		MinSellPrice:       minSell,
		MaxSellPrice:       maxSell,
		MinRentPrice:       minRent,
		MaxRentPrice:       maxRent,
		MinLandSize:        minLand,
		MaxLandSize:        maxLand,
		IncludeAllVersions: req.IncludeAllVersions,
		RequesterUserID:    userInfo.ID,
		RequesterRoleSlug:  userInfo.RoleSlug,
	}

	// Call service layer for business logic execution
	result, listErr := lh.listingService.ListListings(ctx, input)
	if listErr != nil {
		httperrors.SendHTTPErrorObj(c, listErr)
		return
	}

	// Convert domain models to response DTOs
	data := make([]dto.ListingResponse, 0, len(result.Items))
	for _, item := range result.Items {
		data = append(data, toListingResponse(item))
	}

	// Build response with pagination metadata
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

// parseSortBy validates and normalizes sortBy query parameter
//
// Allowed values: id, status (case-insensitive)
// Default: id
func parseSortBy(raw string) (string, error) {
	trimmed := strings.TrimSpace(strings.ToLower(raw))
	if trimmed == "" {
		return "id", nil // Default sort by ID (creation order proxy)
	}

	allowed := map[string]string{
		"id":     "id",
		"status": "status",
	}

	if normalized, ok := allowed[trimmed]; ok {
		return normalized, nil
	}

	return "", fmt.Errorf("invalid sort field (allowed: id, status)")
}

// parseSortOrder validates and normalizes sortOrder query parameter
//
// Allowed values: asc, desc (case-insensitive)
// Default: desc
func parseSortOrder(raw string) (string, error) {
	trimmed := strings.TrimSpace(strings.ToLower(raw))
	if trimmed == "" {
		return "desc", nil // Default descending order
	}

	if trimmed == "asc" || trimmed == "desc" {
		return trimmed, nil
	}

	return "", fmt.Errorf("invalid sort order (allowed: asc, desc)")
}

// parseListingStatus converts string status to ListingStatus enum
//
// Accepts both enum names (e.g., "PUBLISHED") and numeric values.
// Returns nil if input is empty (no filter applied).
func parseListingStatus(raw string) (*listingmodel.ListingStatus, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return nil, nil
	}

	// Try numeric parsing first
	if numeric, numErr := strconv.Atoi(trimmed); numErr == nil {
		status := listingmodel.ListingStatus(numeric)
		if status >= listingmodel.StatusDraft && status <= listingmodel.StatusNeedsRevision {
			return &status, nil
		}
		return nil, fmt.Errorf("invalid listing status value")
	}

	// Try enum name parsing (normalize and map)
	upper := strings.ToUpper(trimmed)
	normalized := strings.ReplaceAll(strings.ReplaceAll(upper, " ", ""), "-", "")
	normalized = strings.ReplaceAll(normalized, "_", "")

	mapping := map[string]listingmodel.ListingStatus{
		"DRAFT":                    listingmodel.StatusDraft,
		"PENDINGAVAILABILITY":      listingmodel.StatusPendingAvailability,
		"PENDINGPHOTOSCHEDULING":   listingmodel.StatusPendingPhotoScheduling,
		"PENDINGPHOTOCONFIRMATION": listingmodel.StatusPendingPhotoConfirmation,
		"PHOTOSSCHEDULED":          listingmodel.StatusPhotosScheduled,
		"PENDINGPHOTOPROCESSING":   listingmodel.StatusPendingPhotoProcessing,
		"PENDINGOWNERAPPROVAL":     listingmodel.StatusPendingOwnerApproval,
		"REJECTEDBYOWNER":          listingmodel.StatusRejectedByOwner,
		"PENDINGADMINREVIEW":       listingmodel.StatusPendingAdminReview,
		"PUBLISHED":                listingmodel.StatusPublished,
		"UNDEROFFER":               listingmodel.StatusUnderOffer,
		"UNDERNEGOTIATION":         listingmodel.StatusUnderNegotiation,
		"CLOSED":                   listingmodel.StatusClosed,
		"SUSPENDED":                listingmodel.StatusSuspended,
		"EXPIRED":                  listingmodel.StatusExpired,
		"ARCHIVED":                 listingmodel.StatusArchived,
		"NEEDSREVISION":            listingmodel.StatusNeedsRevision,
	}
	if status, ok := mapping[normalized]; ok {
		return &status, nil
	}

	return nil, fmt.Errorf("invalid listing status")
}

// parseOptionalUint32 parses optional uint32 from string
func parseOptionalUint32(raw string) (*uint32, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return nil, nil
	}
	value, err := strconv.ParseUint(trimmed, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("invalid numeric value")
	}
	v := uint32(value)
	return &v, nil
}

// parseOptionalInt64 parses optional int64 from string
func parseOptionalInt64(raw string) (*int64, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return nil, nil
	}
	value, err := strconv.ParseInt(trimmed, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid numeric value")
	}
	return &value, nil
}

// parseOptionalFloat64 parses optional float64 from string
func parseOptionalFloat64(raw string) (*float64, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return nil, nil
	}
	value, err := strconv.ParseFloat(trimmed, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid decimal value")
	}
	return &value, nil
}

// toListingResponse converts service item to DTO response
func toListingResponse(item listingservices.ListListingsItem) dto.ListingResponse {
	listing := item.Listing

	price := listing.SellNet()
	if price == 0 {
		price = listing.RentNet()
	}

	var draftVersionID *int64
	if draft, ok := listing.DraftVersion(); ok && draft != nil {
		if draftID := draft.ID(); draftID > 0 {
			draftVersionID = &draftID
		}
	}

	activeVersionID := listing.ActiveVersionID()
	if activeVersionID == 0 {
		activeVersionID = listing.ID()
	}

	return dto.ListingResponse{
		ID:                listing.ID(),
		ListingIdentityID: listing.IdentityID(),
		ListingUUID:       listing.UUID(),
		ActiveVersionID:   activeVersionID,
		DraftVersionID:    draftVersionID,
		Version:           listing.Version(),
		Title:             strings.TrimSpace(listing.Title()),
		Description:       strings.TrimSpace(listing.Description()),
		Price:             price,
		Status:            listing.Status().String(),
		PropertyType:      int(listing.ListingType()),
		ZipCode:           listing.ZipCode(),
		Number:            listing.Number(),
		UserID:            listing.UserID(),
	}
}

// computeTotalPages calculates total pages from total count and limit
func computeTotalPages(total int64, limit int) int {
	if limit <= 0 || total <= 0 {
		return 0
	}

	pages := int(total / int64(limit))
	if total%int64(limit) != 0 {
		pages++
	}

	if pages == 0 && total > 0 {
		return 1
	}

	return pages
}
