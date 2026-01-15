package listinghandlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/converters"
	dto "github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	globalmodel "github.com/projeto-toq/toq_server/internal/core/model/global_model"
	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
	listingrepository "github.com/projeto-toq/toq_server/internal/core/port/right/repository/listing_repository"
	listingservices "github.com/projeto-toq/toq_server/internal/core/service/listing_service"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

// ListListings retrieves a paginated list of active listing versions with optional filters and sorting
//
//	@Summary      List active listing versions with filters and sorting
//	@Description  Retrieves a paginated list of active listing versions (versions linked via listing_identities.active_version_id).
//	              Visibility rules: Owners are auto-scoped to their own listings; Realtors can see listings from any owner but
//	              are forced to status PUBLISHED and active versions only. includeAllVersions=true is ignored for realtors.
//	              Supports filtering by status, code, title, location (zipCode, city, neighborhood), owner (userId),
//	              and price/size ranges. Results can be sorted by id (creation date proxy) or status.
//	              Default sorting: id DESC (newest first).
//	@Tags         Listings
//	@Produce      json
//	@Security     BearerAuth
//	@Param        Authorization       header  string  true   "Bearer token for authentication" Extensions(x-example=Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...)
//	@Param        page                query   int     false  "Page number (1-indexed)" minimum(1) default(1) Extensions(x-example=1)
//	@Param        limit               query   int     false  "Items per page" minimum(1) maximum(100) default(20) Extensions(x-example=20)
//	@Param        sortBy              query   string  false  "Field to sort by" Enums(id, status, zipCode, city, neighborhood, street, number, state, complex) default(id) Extensions(x-example=id)
//	@Param        sortOrder           query   string  false  "Sort direction" Enums(asc, desc) default(desc) Extensions(x-example=desc)
//	@Param        status              query   string  false  "Filter by listing status (enum name or numeric)" Extensions(x-example="PUBLISHED")
//	@Param        code                query   int     false  "Filter by exact listing code" Extensions(x-example=1024)
//	@Param        title               query   string  false  "Filter by listing title/description (supports '*' wildcard)" Extensions(x-example="*garden*")
//	@Param        userId              query   int     false  "Filter by owner user ID (owners auto-filtered to their own listings)" Extensions(x-example=55)
//	@Param        zipCode             query   string  false  "Filter by zip code (digits only; supports '*' wildcard)" Extensions(x-example="06543*")
//	@Param        city                query   string  false  "Filter by city (supports '*' wildcard)" Extensions(x-example="*Paulista*")
//	@Param        neighborhood        query   string  false  "Filter by neighborhood (supports '*' wildcard)" Extensions(x-example="*Centro*")
//	@Param        street              query   string  false  "Filter by street (supports '*' wildcard)" Extensions(x-example="*Paulista*")
//	@Param        number              query   string  false  "Filter by address number (supports '*' wildcard and S/N)" Extensions(x-example="12*")
//	@Param        complement          query   string  false  "Filter by complement (supports '*' wildcard)" Extensions(x-example="*Bloco B*")
//	@Param        complex             query   string  false  "Filter by complex/condominium name (supports '*' wildcard)" Extensions(x-example="*Residencial Atlântico*")
//	@Param        state               query   string  false  "Filter by state (UF); accepts wildcard but prefer exact two-letter code" Extensions(x-example="SP")
//	@Param        minSell             query   number  false  "Minimum sell price" Extensions(x-example=100000)
//	@Param        maxSell             query   number  false  "Maximum sell price" Extensions(x-example=900000)
//	@Param        minRent             query   number  false  "Minimum rent price" Extensions(x-example=1500)
//	@Param        maxRent             query   number  false  "Maximum rent price" Extensions(x-example=8000)
//	@Param        minLandSize         query   number  false  "Minimum land size in square meters" Extensions(x-example=120.5)
//	@Param        maxLandSize         query   number  false  "Maximum land size in square meters" Extensions(x-example=500.75)
//	@Param        minSuites           query   int     false  "Minimum suite count (from feature 'Suites')" Extensions(x-example=2)
//	@Param        maxSuites           query   int     false  "Maximum suite count (from feature 'Suites')" Extensions(x-example=4)
//	@Param        propertyTypes       query   []int   false  "Filter by property types (bitmask values)" Extensions(x-example=[1,2])
//	@Param        transactionTypes    query   []int   false  "Filter by transaction types (catalog numeric values)" Extensions(x-example=[1,2])
//	@Param        propertyUse         query   string  false  "Filter by property use" Enums(RESIDENTIAL, COMMERCIAL) Extensions(x-example="RESIDENTIAL")
//	@Param        acceptsExchange     query   bool    false  "Filter listings that accept exchange" Extensions(x-example=true)
//	@Param        acceptsFinancing    query   bool    false  "Filter listings that accept financing" Extensions(x-example=true)
//	@Param        onlySold            query   bool    false  "Return only sold listings" Extensions(x-example=false)
//	@Param        onlyNewListings     query   bool    false  "Return only listings created within configured recency window" Extensions(x-example=false)
//	@Param        onlyPriceChanged    query   bool    false  "Return only listings with price updates within configured recency window" Extensions(x-example=false)
//	@Param        includeAllVersions  query   bool    false  "Include all versions (active + draft). Default: false (active only)" Extensions(x-example=false)
//	@Success      200                 {object}  dto.ListListingsResponse         "Paginated list of listings with metadata"
//	@Failure      400                 {object}  dto.ErrorResponse                "Invalid request parameters (malformed sortBy, sortOrder, or filter values)"
//	@Failure      401                 {object}  dto.ErrorResponse                "Unauthorized (missing or invalid token)"
//	@Failure      403                 {object}  dto.ErrorResponse                "Forbidden (user lacks permission to access this resource)"
//	@Failure      422                 {object}  dto.ErrorResponse                "Validation failed (invalid enum values, range errors)" Extensions(x-example={"code":422,"message":"Validation failed","details":{"field":"sortBy","error":"Invalid sort field"}})
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

	minSuites, err := parseOptionalNonNegativeInt(req.MinSuites)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, coreutils.ValidationError("minSuites", err.Error()))
		return
	}
	maxSuites, err := parseOptionalNonNegativeInt(req.MaxSuites)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, coreutils.ValidationError("maxSuites", err.Error()))
		return
	}
	if minSuites != nil && maxSuites != nil && *minSuites > *maxSuites {
		httperrors.SendHTTPErrorObj(c, coreutils.ValidationError("minSuites", "minSuites cannot be greater than maxSuites"))
		return
	}

	propertyTypes, err := parsePropertyTypes(req.PropertyTypes)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, coreutils.ValidationError("propertyTypes", err.Error()))
		return
	}

	transactionTypes, err := parseTransactionTypes(req.TransactionTypes)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, coreutils.ValidationError("transactionTypes", err.Error()))
		return
	}

	propertyUse, err := parsePropertyUse(req.PropertyUse)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, coreutils.ValidationError("propertyUse", err.Error()))
		return
	}

	newerThanHours := resolveRecencyWindow(req.OnlyNewListings, lh.config.NewListingHoursThreshold)
	priceUpdatedWithin := resolveRecencyWindow(req.OnlyPriceChanged, lh.config.PriceChangedHoursThreshold)

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
		Street:             strings.TrimSpace(req.Street),
		City:               strings.TrimSpace(req.City),
		Neighborhood:       strings.TrimSpace(req.Neighborhood),
		Number:             strings.TrimSpace(req.Number),
		Complement:         strings.TrimSpace(req.Complement),
		Complex:            strings.TrimSpace(req.Complex),
		State:              strings.TrimSpace(req.State),
		UserID:             userIDPtr,
		MinSellPrice:       minSell,
		MaxSellPrice:       maxSell,
		MinRentPrice:       minRent,
		MaxRentPrice:       maxRent,
		MinLandSize:        minLand,
		MaxLandSize:        maxLand,
		MinSuites:          minSuites,
		MaxSuites:          maxSuites,
		PropertyTypes:      propertyTypes,
		TransactionTypes:   transactionTypes,
		PropertyUse:        propertyUse,
		AcceptsExchange:    req.AcceptsExchange,
		AcceptsFinancing:   req.AcceptsFinancing,
		OnlySold:           req.OnlySold,
		OnlyNewerThanHours: newerThanHours,
		PriceUpdatedWithin: priceUpdatedWithin,
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
		"id":           "id",
		"status":       "status",
		"zipcode":      "zipCode",
		"city":         "city",
		"neighborhood": "neighborhood",
		"street":       "street",
		"number":       "number",
		"state":        "state",
		"complex":      "complex",
	}

	if normalized, ok := allowed[trimmed]; ok {
		return normalized, nil
	}

	return "", fmt.Errorf("invalid sort field (allowed: id, status, zipCode, city, neighborhood, street, number, state, complex)")
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

	status, err := listingmodel.ParseListingStatus(trimmed)
	if err != nil {
		return nil, fmt.Errorf("invalid listing status: %w", err)
	}

	return &status, nil
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

// parseOptionalNonNegativeInt parses optional integer (>= 0) from string
func parseOptionalNonNegativeInt(raw string) (*int, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return nil, nil
	}
	value, err := strconv.Atoi(trimmed)
	if err != nil {
		return nil, fmt.Errorf("invalid integer value")
	}
	if value < 0 {
		return nil, fmt.Errorf("value cannot be negative")
	}
	return &value, nil
}

func parsePropertyTypes(raw []uint16) ([]globalmodel.PropertyType, error) {
	if len(raw) == 0 {
		return nil, nil
	}
	result := make([]globalmodel.PropertyType, 0, len(raw))
	for _, v := range raw {
		if v == 0 {
			return nil, fmt.Errorf("property type must be greater than zero")
		}
		result = append(result, globalmodel.PropertyType(v))
	}
	return result, nil
}

func parseTransactionTypes(raw []uint8) ([]listingmodel.TransactionType, error) {
	if len(raw) == 0 {
		return nil, nil
	}
	result := make([]listingmodel.TransactionType, 0, len(raw))
	for _, v := range raw {
		if v == 0 {
			return nil, fmt.Errorf("transaction type must be greater than zero")
		}
		result = append(result, listingmodel.TransactionType(v))
	}
	return result, nil
}

func parsePropertyUse(raw string) (listingrepository.PropertyUseFilter, error) {
	trimmed := strings.TrimSpace(strings.ToUpper(raw))
	switch trimmed {
	case "":
		return listingrepository.PropertyUseUndefined, nil
	case "RESIDENTIAL":
		return listingrepository.PropertyUseResidential, nil
	case "COMMERCIAL":
		return listingrepository.PropertyUseCommercial, nil
	default:
		return listingrepository.PropertyUseUndefined, fmt.Errorf("invalid propertyUse (allowed: RESIDENTIAL, COMMERCIAL)")
	}
}

func resolveRecencyWindow(enabled bool, hours int) *int {
	if !enabled || hours <= 0 {
		return nil
	}
	value := hours
	return &value
}

// toListingResponse converts service item to DTO response
func toListingResponse(item listingservices.ListListingsItem) dto.ListingResponse {
	listing := item.Listing

	price := listing.SellNet()
	if price == 0 {
		price = listing.RentNet()
	}

	complexValue := ""
	if listing.HasComplex() {
		complexValue = strings.TrimSpace(listing.Complex())
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
		PropertyType:      converters.BuildListingPropertyTypeDTO(listing.ListingType()),
		ZipCode:           listing.ZipCode(),
		Street:            strings.TrimSpace(listing.Street()),
		Number:            listing.Number(),
		Complement:        strings.TrimSpace(listing.Complement()),
		Neighborhood:      strings.TrimSpace(listing.Neighborhood()),
		City:              strings.TrimSpace(listing.City()),
		State:             strings.TrimSpace(listing.State()),
		Complex:           complexValue,
		UserID:            listing.UserID(),
		FavoritesCount:    item.FavoritesCount,
		IsFavorite:        item.IsFavorite,
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
