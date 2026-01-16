package listingservices

import (
	"context"

	globalmodel "github.com/projeto-toq/toq_server/internal/core/model/global_model"
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
	Address            string                      // Optional wildcard full address search (concatenated fields)
	UserID             *int64                      // Optional owner user ID filter
	MinSellPrice       *float64                    // Optional minimum sell price
	MaxSellPrice       *float64                    // Optional maximum sell price
	MinRentPrice       *float64                    // Optional minimum rent price
	MaxRentPrice       *float64                    // Optional maximum rent price
	MinLandSize        *float64                    // Optional minimum land size (sq meters)
	MaxLandSize        *float64                    // Optional maximum land size (sq meters)
	MinSuites          *int                        // Optional minimum suite count (derived from features)
	MaxSuites          *int                        // Optional maximum suite count (derived from features)
	PropertyTypes      []globalmodel.PropertyType  // Optional property type filters
	TransactionTypes   []listingmodel.TransactionType
	PropertyUse        listingrepository.PropertyUseFilter
	AcceptsExchange    *bool
	AcceptsFinancing   *bool
	OnlySold           bool
	OnlyNewerThanHours *int
	PriceUpdatedWithin *int
	IncludeAllVersions bool // true: all versions; false: active only (default)

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
	Listing        listingmodel.ListingInterface // Listing domain entity
	FavoritesCount int64                         // Total favorites for this listing identity
	IsFavorite     bool                          // Whether requester has favorited this listing
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
	switch input.RequesterRoleSlug {
	case permissionmodel.RoleSlugOwner:
		ownerID := input.RequesterUserID
		input.UserID = &ownerID
		logger.Debug("listing.list.owner_scope_enforced", "user_id", ownerID)
	case permissionmodel.RoleSlugRealtor:
		published := listingmodel.StatusPublished
		input.Status = &published
		input.UserID = nil
		input.IncludeAllVersions = false
		logger.Debug("listing.list.realtor_scope_published_only", "user_id", input.RequesterUserID)
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

	// Build repository filter with all parameters
	repoFilter := listingrepository.ListListingsFilter{
		Page:               input.Page,
		Limit:              input.Limit,
		SortBy:             input.SortBy,
		SortOrder:          input.SortOrder,
		Status:             input.Status,
		Code:               input.Code,
		Title:              utils.NormalizeSearchPattern(input.Title),
		Address:            utils.NormalizeSearchPattern(input.Address),
		UserID:             input.UserID,
		MinSellPrice:       input.MinSellPrice,
		MaxSellPrice:       input.MaxSellPrice,
		MinRentPrice:       input.MinRentPrice,
		MaxRentPrice:       input.MaxRentPrice,
		MinLandSize:        input.MinLandSize,
		MaxLandSize:        input.MaxLandSize,
		MinSuites:          input.MinSuites,
		MaxSuites:          input.MaxSuites,
		PropertyTypes:      input.PropertyTypes,
		TransactionTypes:   input.TransactionTypes,
		PropertyUse:        input.PropertyUse,
		AcceptsExchange:    input.AcceptsExchange,
		AcceptsFinancing:   input.AcceptsFinancing,
		OnlySold:           input.OnlySold,
		OnlyNewerThanHours: input.OnlyNewerThanHours,
		PriceUpdatedWithin: input.PriceUpdatedWithin,
		IncludeAllVersions: input.IncludeAllVersions,
	}

	// Call repository to execute query with filters and sorting
	result, listErr := ls.listingRepository.ListListings(ctx, tx, repoFilter)
	if listErr != nil {
		utils.SetSpanError(ctx, listErr)
		logger.Error("listing.list.repo_error", "error", listErr)
		return ListListingsOutput{}, utils.InternalError("")
	}

	identityIDs := make([]int64, 0, len(result.Records))
	seen := make(map[int64]struct{}, len(result.Records))
	for _, record := range result.Records {
		if id := record.Listing.IdentityID(); id > 0 {
			if _, ok := seen[id]; !ok {
				seen[id] = struct{}{}
				identityIDs = append(identityIDs, id)
			}
		}
	}

	favoriteCounts := make(map[int64]int64)
	if len(identityIDs) > 0 {
		countMap, countErr := ls.favoriteRepo.CountByListingIdentities(ctx, tx, identityIDs)
		if countErr != nil {
			utils.SetSpanError(ctx, countErr)
			logger.Error("listing.list.fav_count_error", "err", countErr)
			return ListListingsOutput{}, utils.InternalError("")
		}
		favoriteCounts = countMap
	}

	userFlags := make(map[int64]bool)
	if len(identityIDs) > 0 && input.RequesterUserID > 0 {
		flags, flagErr := ls.favoriteRepo.GetUserFlags(ctx, tx, identityIDs, input.RequesterUserID)
		if flagErr != nil {
			utils.SetSpanError(ctx, flagErr)
			logger.Error("listing.list.fav_flags_error", "err", flagErr)
			return ListListingsOutput{}, utils.InternalError("")
		}
		userFlags = flags
	}

	// Convert repository records to service output items
	items := make([]ListListingsItem, 0, len(result.Records))
	for _, record := range result.Records {
		identityID := record.Listing.IdentityID()
		items = append(items, ListListingsItem{
			Listing:        record.Listing,
			FavoritesCount: favoriteCounts[identityID],
			IsFavorite:     userFlags[identityID],
		})
	}

	return ListListingsOutput{
		Items: items,
		Total: result.Total,
		Page:  repoFilter.Page,
		Limit: repoFilter.Limit,
	}, nil
}
