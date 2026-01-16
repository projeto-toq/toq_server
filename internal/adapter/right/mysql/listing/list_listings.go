package mysqllistingadapter

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	listingconverters "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/listing/converters"
	globalmodel "github.com/projeto-toq/toq_server/internal/core/model/global_model"
	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
	listingrepository "github.com/projeto-toq/toq_server/internal/core/port/right/repository/listing_repository"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

const suitesCountSubquery = `(
	SELECT COALESCE(SUM(f.qty), 0)
	FROM features f
	INNER JOIN base_features bf ON bf.id = f.feature_id
	WHERE f.listing_version_id = lv.id
	  AND bf.feature = 'Suites'
)`

// ListListings retrieves a filtered, sorted, and paginated list of listing versions
//
// This function executes a complex SQL query with multiple optional filters and dynamic sorting.
// By default, only active versions are returned (versions linked via listing_identities.active_version_id).
//
// Query Structure:
//   - Base: SELECT from listing_versions JOIN listing_identities
//   - WHERE: deleted=0 + optional filters (status, code, title, location, prices, sizes)
//   - Active filter: lv.id = li.active_version_id (unless includeAllVersions=true)
//   - ORDER BY: Dynamic based on filter.SortBy and filter.SortOrder
//   - LIMIT/OFFSET: Pagination
//
// Sorting Options:
//   - id: Order by listing version ID (proxy for creation date - higher ID = newer)
//   - status: Order by status enum value
//
// Parameters:
//   - ctx: Context for tracing, cancellation, and logging
//   - tx: Database transaction (can be nil for standalone queries, must not be nil for consistency)
//   - filter: ListListingsFilter with pagination, sorting, and filter criteria
//
// Returns:
//   - result: ListListingsResult with listing records and total count
//   - error: Database errors (query execution, scan errors, count errors)
//
// Business Rules:
//   - Only non-deleted versions (lv.deleted = 0) and identities (li.deleted = 0)
//   - Active versions only by default (lv.id = li.active_version_id)
//   - Wildcard search uses SQL LIKE with '%' pattern
//   - Price/size filters use >= and <= operators
func (la *ListingAdapter) ListListings(ctx context.Context, tx *sql.Tx, filter listingrepository.ListListingsFilter) (listingrepository.ListListingsResult, error) {
	// Initialize tracing for observability (metrics + distributed tracing)
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return listingrepository.ListListingsResult{}, err
	}
	defer spanEnd()

	// Attach logger to context for request_id/trace_id propagation
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	// Normalize pagination defaults (defensive programming)
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.Limit <= 0 {
		filter.Limit = 20
	}

	// Build WHERE conditions array (all conditions AND-ed)
	conditions := []string{
		"lv.deleted = 0", // Only non-deleted versions
		"li.deleted = 0", // Only non-deleted identities
	}
	args := make([]any, 0)

	// Add active version filter by default (unless includeAllVersions=true)
	if !filter.IncludeAllVersions {
		conditions = append(conditions, "lv.id = li.active_version_id")
	}

	// Optional filter: status
	if filter.Status != nil {
		conditions = append(conditions, "lv.status = ?")
		args = append(args, int(*filter.Status))
	}

	// Optional filter: exact code match
	if filter.Code != nil {
		conditions = append(conditions, "lv.code = ?")
		args = append(args, int64(*filter.Code))
	}

	// Optional filter: wildcard title/description search (LIKE pattern)
	if filter.Title != "" {
		conditions = append(conditions, "(COALESCE(lv.title, '') LIKE ? OR COALESCE(lv.description, '') LIKE ?)")
		args = append(args, filter.Title, filter.Title)
	}

	// Optional filter: wildcard full address search (concatenated fields without spaces)
	if filter.Address != "" {
		conditions = append(conditions, "LOWER(REPLACE(CONCAT_WS(' ', COALESCE(lv.zip_code,''), COALESCE(lv.city,''), COALESCE(lv.neighborhood,''), COALESCE(lv.street,''), COALESCE(lv.number,''), COALESCE(lv.complement,''), COALESCE(lv.state,'')), ' ', '')) LIKE LOWER(REPLACE(?, ' ', ''))")
		args = append(args, filter.Address)
	}

	// Optional filter: owner user ID
	if filter.UserID != nil {
		conditions = append(conditions, "lv.user_id = ?")
		args = append(args, *filter.UserID)
	}

	if len(filter.PropertyTypes) > 0 {
		placeholders := makePlaceholders(len(filter.PropertyTypes))
		conditions = append(conditions, fmt.Sprintf("lv.type IN (%s)", placeholders))
		for _, pt := range filter.PropertyTypes {
			args = append(args, pt)
		}
	}

	effectiveTransactionTypes := normalizeTransactionTypes(filter.TransactionTypes)
	if len(effectiveTransactionTypes) > 0 {
		placeholders := makePlaceholders(len(effectiveTransactionTypes))
		conditions = append(conditions, fmt.Sprintf("lv.transaction IN (%s)", placeholders))
		for _, tt := range effectiveTransactionTypes {
			args = append(args, tt)
		}
	}

	switch filter.PropertyUse {
	case listingrepository.PropertyUseResidential:
		residentialTypes := []globalmodel.PropertyType{
			globalmodel.Apartment,
			globalmodel.House,
			globalmodel.OffPlanHouse,
			globalmodel.ResidencialLand,
		}
		placeholders := makePlaceholders(len(residentialTypes))
		conditions = append(conditions, fmt.Sprintf("lv.type IN (%s)", placeholders))
		for _, pt := range residentialTypes {
			args = append(args, pt)
		}
	case listingrepository.PropertyUseCommercial:
		commercialTypes := []globalmodel.PropertyType{
			globalmodel.CommercialStore,
			globalmodel.CommercialFloor,
			globalmodel.Suite,
			globalmodel.CommercialLand,
			globalmodel.Warehouse,
			globalmodel.Building,
		}
		placeholders := makePlaceholders(len(commercialTypes))
		conditions = append(conditions, fmt.Sprintf("lv.type IN (%s)", placeholders))
		for _, pt := range commercialTypes {
			args = append(args, pt)
		}
	}

	if filter.AcceptsExchange != nil {
		conditions = append(conditions, "lv.exchange = ?")
		args = append(args, boolToTinyInt(*filter.AcceptsExchange))
	}

	if filter.AcceptsFinancing != nil {
		conditions = append(conditions, "lv.financing = ?")
		args = append(args, boolToTinyInt(*filter.AcceptsFinancing))
	}

	if filter.OnlySold {
		conditions = append(conditions, "lv.status = ?")
		args = append(args, int(listingmodel.StatusClosed))
	}

	if filter.OnlyNewerThanHours != nil {
		conditions = append(conditions, "TIMESTAMPDIFF(HOUR, lv.created_at, NOW()) <= ?")
		args = append(args, *filter.OnlyNewerThanHours)
	}

	if filter.PriceUpdatedWithin != nil {
		conditions = append(conditions,
			"lv.price_updated_at IS NOT NULL",
			"lv.price_updated_at > lv.created_at",
			"TIMESTAMPDIFF(HOUR, lv.price_updated_at, NOW()) <= ?",
		)
		args = append(args, *filter.PriceUpdatedWithin)
	}

	// Optional filter: sell price range
	if filter.MinSellPrice != nil {
		conditions = append(conditions, "COALESCE(lv.sell_net, 0) >= ?")
		args = append(args, *filter.MinSellPrice)
	}
	if filter.MaxSellPrice != nil {
		conditions = append(conditions, "COALESCE(lv.sell_net, 0) <= ?")
		args = append(args, *filter.MaxSellPrice)
	}

	// Optional filter: rent price range
	if filter.MinRentPrice != nil {
		conditions = append(conditions, "COALESCE(lv.rent_net, 0) >= ?")
		args = append(args, *filter.MinRentPrice)
	}
	if filter.MaxRentPrice != nil {
		conditions = append(conditions, "COALESCE(lv.rent_net, 0) <= ?")
		args = append(args, *filter.MaxRentPrice)
	}

	// Optional filter: land size range
	if filter.MinLandSize != nil {
		conditions = append(conditions, "COALESCE(lv.land_size, 0) >= ?")
		args = append(args, *filter.MinLandSize)
	}
	if filter.MaxLandSize != nil {
		conditions = append(conditions, "COALESCE(lv.land_size, 0) <= ?")
		args = append(args, *filter.MaxLandSize)
	}

	// Optional filter: suite count derived from features
	if filter.MinSuites != nil {
		conditions = append(conditions, fmt.Sprintf("%s >= ?", suitesCountSubquery))
		args = append(args, *filter.MinSuites)
	}
	if filter.MaxSuites != nil {
		conditions = append(conditions, fmt.Sprintf("%s <= ?", suitesCountSubquery))
		args = append(args, *filter.MaxSuites)
	}

	// Construct WHERE clause (all conditions AND-ed)
	whereClause := "WHERE " + strings.Join(conditions, " AND ")

	// Construct ORDER BY clause based on sortBy and sortOrder
	orderByClause := buildOrderByClause(filter.SortBy, filter.SortOrder)

	// Base SELECT with explicit column list (never use SELECT *)
	baseSelect := fmt.Sprintf(`SELECT
%s
FROM listing_versions lv
INNER JOIN listing_identities li ON li.id = lv.listing_identity_id`, listingSelectColumns)

	// Construct full query with WHERE, ORDER BY, LIMIT, OFFSET
	listQuery := baseSelect + " " + whereClause + " " + orderByClause + " LIMIT ? OFFSET ?"
	listArgs := append([]any{}, args...)
	offset := (filter.Page - 1) * filter.Limit
	listArgs = append(listArgs, filter.Limit, offset)

	// Execute query via InstrumentedAdapter (auto-generates metrics + tracing)
	rows, queryErr := la.QueryContext(ctx, tx, "select", listQuery, listArgs...)
	if queryErr != nil {
		// Mark span as error for distributed tracing analysis
		utils.SetSpanError(ctx, queryErr)
		logger.Error("mysql.listing.list.query_error", "error", queryErr)
		return listingrepository.ListListingsResult{}, fmt.Errorf("list listings query: %w", queryErr)
	}
	defer rows.Close()

	// Initialize result container
	result := listingrepository.ListListingsResult{}

	// Scan all result rows into entities
	for rows.Next() {
		entity, scanErr := scanListingEntity(rows)
		if scanErr != nil {
			utils.SetSpanError(ctx, scanErr)
			logger.Error("mysql.listing.list.scan_error", "error", scanErr)
			return listingrepository.ListListingsResult{}, fmt.Errorf("scan listing row: %w", scanErr)
		}

		// Convert entity to domain model (separation of concerns)
		listing := listingconverters.ListingEntityToDomain(entity)
		if listing != nil {
			result.Records = append(result.Records, listingrepository.ListingRecord{
				Listing: listing,
			})
		}
	}

	// Check for iteration errors (connection issues, etc.)
	if err = rows.Err(); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.listing.list.rows_error", "error", err)
		return listingrepository.ListListingsResult{}, fmt.Errorf("iterate listing rows: %w", err)
	}

	// Execute count query to get total matching records (for pagination metadata)
	countQuery := "SELECT COUNT(*) FROM listing_versions lv INNER JOIN listing_identities li ON li.id = lv.listing_identity_id " + whereClause
	var total int64
	if countErr := la.QueryRowContext(ctx, tx, "select", countQuery, args...).Scan(&total); countErr != nil {
		utils.SetSpanError(ctx, countErr)
		logger.Error("mysql.listing.list.count_error", "error", countErr)
		return listingrepository.ListListingsResult{}, fmt.Errorf("count listings: %w", countErr)
	}
	result.Total = total

	return result, nil
}

// buildOrderByClause constructs ORDER BY SQL clause based on sortBy field and sortOrder direction
//
// Validates sortBy against allowed fields and constructs safe SQL without injection risks.
// Falls back to "id DESC" if invalid parameters are provided.
//
// Allowed sortBy fields:
//   - id: Order by listing version ID (proxy for creation date)
//   - status: Order by status enum value
//
// Parameters:
//   - sortBy: Field name (id, status)
//   - sortOrder: Direction (asc, desc)
//
// Returns:
//   - SQL ORDER BY clause (e.g., "ORDER BY lv.id DESC")
func buildOrderByClause(sortBy, sortOrder string) string {
	// Map sortBy input to actual column names (validate against whitelist)
	columnMap := map[string]string{
		"id":           "lv.id",
		"status":       "lv.status",
		"zipcode":      "lv.zip_code",
		"city":         "lv.city",
		"neighborhood": "lv.neighborhood",
		"street":       "lv.street",
		"number":       "lv.number",
		"state":        "lv.state",
		"complex":      "lv.complex",
	}

	column, ok := columnMap[strings.ToLower(sortBy)]
	if !ok {
		// Invalid sortBy - fall back to default
		column = "lv.id"
	}

	// Validate sortOrder (only asc/desc allowed)
	direction := strings.ToUpper(sortOrder)
	if direction != "ASC" && direction != "DESC" {
		// Invalid sortOrder - fall back to default
		direction = "DESC"
	}

	return fmt.Sprintf("ORDER BY %s %s", column, direction)
}

func makePlaceholders(count int) string {
	if count <= 0 {
		return ""
	}
	return strings.TrimRight(strings.Repeat("?,", count), ",")
}

func boolToTinyInt(value bool) int {
	if value {
		return 1
	}
	return 0
}

// normalizeTransactionTypes expands sale/rent filters to include listings marked as both.
// Rules:
// - sale => sale + both
// - rent => rent + both
// - sale + rent => sale + rent + both
// - both (alone) => both only
func normalizeTransactionTypes(requested []listingmodel.TransactionType) []listingmodel.TransactionType {
	if len(requested) == 0 {
		return nil
	}

	hasSale := false
	hasRent := false
	hasBoth := false

	for _, t := range requested {
		switch t {
		case listingmodel.TransactionType(1):
			hasSale = true
		case listingmodel.TransactionType(2):
			hasRent = true
		case listingmodel.TransactionType(3):
			hasBoth = true
		}
	}

	result := make([]listingmodel.TransactionType, 0, 3)
	switch {
	case hasSale && hasRent:
		result = append(result, listingmodel.TransactionType(1), listingmodel.TransactionType(2), listingmodel.TransactionType(3))
	case hasSale && !hasRent:
		result = append(result, listingmodel.TransactionType(1), listingmodel.TransactionType(3))
	case hasRent && !hasSale:
		result = append(result, listingmodel.TransactionType(2), listingmodel.TransactionType(3))
	case !hasSale && !hasRent && hasBoth:
		result = append(result, listingmodel.TransactionType(3))
	}

	return result
}
