package listingservices

import (
	"context"
	"strings"

	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
	permissionmodel "github.com/projeto-toq/toq_server/internal/core/model/permission_model"
	listingrepository "github.com/projeto-toq/toq_server/internal/core/port/right/repository/listing_repository"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// ListListingsInput captures filters, pagination, and sorting for listing search.
//
// This input struct aggregates all parameters needed to retrieve a filtered, sorted, and paginated
// list of listing versions. Includes role-based filtering (owners auto-scoped to their own listings).
type ListListingsInput struct {
	// Pagination
	Page  int // 1-indexed page number (default: 1)
	Limit int // Items per page (default: 20, max: 100)

	// Sorting
	SortBy    string // Field to sort by: id, status (default: id)
	SortOrder string // Sort direction: asc, desc (default: desc)

	// Filters
	Status             *listingmodel.ListingStatus // Optional listing status filter
	Code               *uint32                     // Optional exact code filter
	Title              string                      // Optional wildcard title/description search
	ZipCode            string                      // Optional wildcard zip code filter
	City               string                      // Optional wildcard city filter
	Neighborhood       string                      // Optional wildcard neighborhood filter
	UserID             *int64                      // Optional owner user ID filter
	MinSellPrice       *float64                    // Optional minimum sell price
	MaxSellPrice       *float64                    // Optional maximum sell price
	MinRentPrice       *float64                    // Optional minimum rent price
	MaxRentPrice       *float64                    // Optional maximum rent price
	MinLandSize        *float64                    // Optional minimum land size (sq meters)
	MaxLandSize        *float64                    // Optional maximum land size (sq meters)
	MinSuites          *int                        // Optional minimum suite count (derived from features)
	MaxSuites          *int                        // Optional maximum suite count (derived from features)
	IncludeAllVersions bool                        // true: all versions; false: active only (default)

	// Security context
	RequesterUserID   int64                    // Authenticated user ID
	RequesterRoleSlug permissionmodel.RoleSlug // Authenticated user role
}

// ListListingsOutput encapsulates listings and paging metadata.
//
// Contains the filtered listing collection and pagination metadata for UI rendering.
type ListListingsOutput struct {
	Items []ListListingsItem // Listing items for current page
	Total int64              // Total count matching filters (all pages)
	Page  int                // Current page number
	Limit int                // Items per page
}

// ListListingsItem wraps a listing entity with metadata for response assembly.
//
// Currently contains only the listing entity, but designed for future extension
// (e.g., favorite status, offer counts, etc.).
type ListListingsItem struct {
	Listing listingmodel.ListingInterface // Listing domain entity
}

// ListListings returns listings filtered, sorted, and paginated for admin panel or owner consumption.
//
// This method orchestrates the complete listing retrieval flow:
//  1. Applies role-based security (owners auto-filtered to own listings)
//  2. Validates and normalizes pagination and sorting parameters
//  3. Starts read-only transaction for consistency
//  4. Constructs repository filter with all parameters
//  5. Retrieves listings from repository with sorting applied
//  6. Rolls back transaction (read-only, no commit needed)
//
// Business Rules:
//   - Owners (RoleSlugOwner) are automatically scoped to their own listings (userID filter enforced)
//   - Default pagination: page=1, limit=20
//   - Default sorting: id DESC (newest listings first)
//   - Active versions only by default (unless includeAllVersions=true)
//   - Wildcard search supports '*' character for partial matches
//
// Parameters:
//   - ctx: Context for tracing, cancellation, and logging. Must contain request metadata.
//   - input: ListListingsInput with all filters, pagination, and sorting parameters
//
// Returns:
//   - output: ListListingsOutput with paginated listing collection and metadata
//   - err: Infrastructure error (500) for database failures; never returns domain errors
//
// Side Effects:
//   - None (read-only operation)
func (ls *listingService) ListListings(ctx context.Context, input ListListingsInput) (ListListingsOutput, error) {
	// Initialize tracing for distributed observability
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return ListListingsOutput{}, utils.InternalError("Failed to generate tracer")
	}
	defer spanEnd()

	// Ensure logger propagation with request_id and trace_id
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	// Security: Owners can only see their own listings (auto-enforce userID filter)
	if input.RequesterRoleSlug == permissionmodel.RoleSlugOwner {
		ownerID := input.RequesterUserID
		input.UserID = &ownerID
		logger.Debug("listing.list.owner_scope_enforced", "user_id", ownerID)
	}

	// Normalize pagination defaults
	if input.Page <= 0 {
		input.Page = 1
	}
	if input.Limit <= 0 {
		input.Limit = 20
	}

	// Normalize sorting defaults
	if input.SortBy == "" {
		input.SortBy = "id"
	}
	if input.SortOrder == "" {
		input.SortOrder = "desc"
	}

	// Start read-only transaction for consistent snapshot view
	tx, txErr := ls.gsi.StartReadOnlyTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("listing.list.tx_start_failed", "error", txErr)
		return ListListingsOutput{}, utils.InternalError("")
	}
	defer func() {
		// Always rollback read-only transactions (no commit needed)
		_ = ls.gsi.RollbackTransaction(ctx, tx)
	}()

	// Sanitize zipCode filter (remove non-numeric/wildcard characters)
	zipFilter := sanitizeZipFilter(input.ZipCode)

	// Build repository filter with all parameters
	repoFilter := listingrepository.ListListingsFilter{
		Page:               input.Page,
		Limit:              input.Limit,
		SortBy:             input.SortBy,
		SortOrder:          input.SortOrder,
		Status:             input.Status,
		Code:               input.Code,
		Title:              utils.NormalizeSearchPattern(input.Title),
		ZipCode:            utils.NormalizeSearchPattern(zipFilter),
		City:               utils.NormalizeSearchPattern(input.City),
		Neighborhood:       utils.NormalizeSearchPattern(input.Neighborhood),
		UserID:             input.UserID,
		MinSellPrice:       input.MinSellPrice,
		MaxSellPrice:       input.MaxSellPrice,
		MinRentPrice:       input.MinRentPrice,
		MaxRentPrice:       input.MaxRentPrice,
		MinLandSize:        input.MinLandSize,
		MaxLandSize:        input.MaxLandSize,
		MinSuites:          input.MinSuites,
		MaxSuites:          input.MaxSuites,
		IncludeAllVersions: input.IncludeAllVersions,
	}

	// Call repository to execute query with filters and sorting
	result, listErr := ls.listingRepository.ListListings(ctx, tx, repoFilter)
	if listErr != nil {
		utils.SetSpanError(ctx, listErr)
		logger.Error("listing.list.repo_error", "error", listErr)
		return ListListingsOutput{}, utils.InternalError("")
	}

	// Convert repository records to service output items
	items := make([]ListListingsItem, 0, len(result.Records))
	for _, record := range result.Records {
		items = append(items, ListListingsItem{
			Listing: record.Listing,
		})
	}

	return ListListingsOutput{
		Items: items,
		Total: result.Total,
		Page:  repoFilter.Page,
		Limit: repoFilter.Limit,
	}, nil
}

// sanitizeZipFilter removes non-numeric characters from zip code filter (except wildcards)
//
// Allows only digits, '*', and '%' for SQL LIKE pattern matching.
// Returns empty string if input is blank or contains only invalid characters.
func sanitizeZipFilter(raw string) string {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return ""
	}
	var builder strings.Builder
	for _, r := range trimmed {
		switch {
		case r >= '0' && r <= '9':
			builder.WriteRune(r)
		case r == '*', r == '%':
			builder.WriteRune(r)
		}
	}
	return builder.String()
}
