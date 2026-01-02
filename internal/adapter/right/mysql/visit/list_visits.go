package mysqlvisitadapter

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/projeto-toq/toq_server/internal/adapter/right/mysql/visit/converters"
	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

const visitsMaxPageSize = 50

// ListVisits retrieves a paginated and filtered list of visits from listing_visits table.
// Returns a result set with total count and matching visit records.
//
// This function supports dynamic filtering by multiple criteria and executes two queries:
// 1. COUNT query for total matching records (pagination metadata)
// 2. SELECT query with LIMIT/OFFSET for current page
//
// Parameters:
//   - ctx: Context for tracing and logging
//   - tx: Database transaction (can be nil for read-only queries)
//   - filter: VisitListFilter with optional filters and pagination parameters
//
// Returns:
//   - result: VisitListResult containing visits array and total count
//   - error: Database errors or scan failures
//
// Filter Criteria (all optional, combined with AND):
//   - ListingID: Filter by specific listing (*int64)
//   - OwnerID: Filter by property owner (*int64)
//   - RealtorID: Filter by realtor (*int64)
//   - Statuses: Filter by status array ([]VisitStatus) - uses IN clause
//   - From: Filter visits ending after this time (*time.Time) - inclusive
//   - To: Filter visits starting before this time (*time.Time) - inclusive
//
// Pagination:
//   - Limit: Items per page (default/max: 50) - prevents memory/performance issues
//   - Page: Page number (1-indexed, default: 1)
//   - Offset calculated as: (page - 1) * limit
//
// Sorting:
//   - Default: ORDER BY scheduled_start ASC (chronological order)
//   - Fixed sort to maintain predictable pagination behavior
//
// Performance Notes:
//   - Recommended indexes: idx_visits_listing_id, idx_visits_scheduled_start
//   - COUNT query executed before SELECT (ensure consistency)
//   - Empty filters return all visits (use with caution in production)
func (a *VisitAdapter) ListVisits(ctx context.Context, tx *sql.Tx, filter listingmodel.VisitListFilter) (listingmodel.VisitListResult, error) {
	// Initialize tracing for observability
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return listingmodel.VisitListResult{}, err
	}
	defer spanEnd()

	// Ensure logger propagation with request_id and trace_id
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	// Build dynamic WHERE clause based on provided filters
	conditions := make([]string, 0)
	args := make([]any, 0)

	scheduledStartExpr := "CAST(CONCAT(scheduled_date, ' ', scheduled_time_start) AS DATETIME)"
	scheduledEndExpr := "CAST(CONCAT(scheduled_date, ' ', scheduled_time_end) AS DATETIME)"

	// Filter by listing identity (exact match)
	if filter.ListingIdentityID != nil {
		conditions = append(conditions, "listing_identity_id = ?")
		args = append(args, *filter.ListingIdentityID)
	}

	// Filter by owner user (exact match)
	if filter.OwnerUserID != nil {
		conditions = append(conditions, "listing_identity_id IN (SELECT id FROM listing_identities WHERE user_id = ?)")
		args = append(args, *filter.OwnerUserID)
	}

	// Filter by requester user (exact match)
	if filter.RequesterUserID != nil {
		conditions = append(conditions, "user_id = ?")
		args = append(args, *filter.RequesterUserID)
	}

	// Filter by status array (IN clause for multiple statuses)
	if len(filter.Statuses) > 0 {
		placeholders := make([]string, len(filter.Statuses))
		for i, status := range filter.Statuses {
			placeholders[i] = "?"
			args = append(args, string(status))
		}
		conditions = append(conditions, fmt.Sprintf("status IN (%s)", strings.Join(placeholders, ",")))
	}

	// Filter by time range (visits ending after 'from')
	if filter.From != nil {
		conditions = append(conditions, fmt.Sprintf("%s >= ?", scheduledEndExpr))
		args = append(args, *filter.From)
	}

	// Filter by time range (visits starting before 'to')
	if filter.To != nil {
		conditions = append(conditions, fmt.Sprintf("%s <= ?", scheduledStartExpr))
		args = append(args, *filter.To)
	}

	// Default WHERE clause (always true if no filters provided)
	where := "1=1"
	if len(conditions) > 0 {
		where = strings.Join(conditions, " AND ")
	}

	// Execute COUNT query first for total records (pagination metadata)
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM listing_visits WHERE %s", where)
	var total int64
	countRow := a.QueryRowContext(ctx, tx, "list_visits_count", countQuery, args...)
	if err = countRow.Scan(&total); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.visit.list.count_error", "err", err)
		return listingmodel.VisitListResult{}, fmt.Errorf("count visits: %w", err)
	}

	// Normalize pagination parameters (max 50 per page to prevent DoS)
	limit, offset := defaultPagination(filter.Limit, filter.Page, visitsMaxPageSize)

	// Execute main SELECT query with pagination
	// Note: ORDER BY reconstructed scheduled_start ensures chronological listing
	query := fmt.Sprintf(`
		SELECT 
			id,
			listing_identity_id,
			listing_version,
			user_id,
			(SELECT user_id FROM listing_identities li WHERE li.id = listing_visits.listing_identity_id LIMIT 1) AS owner_user_id,
			%s AS scheduled_start,
			%s AS scheduled_end,
			status,
			source,
			notes,
			rejection_reason,
			first_owner_action_at,
			requested_at
		FROM listing_visits
		WHERE %s
		ORDER BY %s ASC
		LIMIT ? OFFSET ?
	`, scheduledStartExpr, scheduledEndExpr, where, scheduledStartExpr)

	// Append limit and offset to args (after filter args)
	params := append(make([]any, 0, len(args)+2), args...)
	params = append(params, limit, offset)

	// Execute query using instrumented adapter
	rows, err := a.QueryContext(ctx, tx, "list_visits", query, params...)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.visit.list.query_error", "err", err)
		return listingmodel.VisitListResult{}, fmt.Errorf("query visits: %w", err)
	}
	defer rows.Close()

	// Iterate through result set and convert each row to domain model
	visits := make([]listingmodel.VisitInterface, 0)
	for rows.Next() {
		visitEntity, scanErr := scanVisitEntity(rows)
		if scanErr != nil {
			utils.SetSpanError(ctx, scanErr)
			logger.Error("mysql.visit.list.scan_error", "err", scanErr)
			return listingmodel.VisitListResult{}, fmt.Errorf("scan visit: %w", scanErr)
		}
		// Convert entity to domain model and append to results
		visits = append(visits, converters.ToVisitModel(visitEntity))
	}

	// Check for iteration errors (connection issues, context cancellation)
	if err = rows.Err(); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.visit.list.rows_error", "err", err)
		return listingmodel.VisitListResult{}, fmt.Errorf("iterate visits: %w", err)
	}

	return listingmodel.VisitListResult{Visits: visits, Total: total}, nil
}
