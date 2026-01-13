package mysqlvisitadapter

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

const (
	visitsMaxPageSize          = 50
	visitWithListingSelectBase = `
	SELECT 
		lv.id,
		lv.listing_identity_id,
		lv.listing_version,
		lv.user_id,
		li.user_id AS owner_user_id,
		%s AS scheduled_start,
		%s AS scheduled_end,
		lv.status,
		lv.source,
		lv.notes,
		lv.rejection_reason,
		lv.first_owner_action_at,
		lv.requested_at,
		active.id AS listing_id,
		active.version AS listing_version_number,
		active.type AS listing_type,
		active.zip_code,
		active.street,
		active.number,
		active.complement,
		active.neighborhood,
		active.city,
		active.state,
		active.title,
		active.description,
		owner.full_name AS owner_full_name,
		owner.created_at AS owner_created_at,
		orm.visit_avg_response_time_seconds,
		realtor.full_name AS realtor_full_name,
		realtor.created_at AS realtor_created_at,
		COALESCE(rv.total_visits, 0) AS realtor_total_visits
	FROM listing_visits lv
	JOIN listing_identities li ON li.id = lv.listing_identity_id
	JOIN users owner ON owner.id = li.user_id
	JOIN users realtor ON realtor.id = lv.user_id
	LEFT JOIN owner_response_metrics orm ON orm.user_id = li.user_id
	LEFT JOIN (
		SELECT user_id, COUNT(*) AS total_visits
		FROM listing_visits
		GROUP BY user_id
	) rv ON rv.user_id = lv.user_id
	LEFT JOIN listing_versions active ON active.id = li.active_version_id AND active.deleted = 0`
)

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

	scheduledStartExpr := "CAST(CONCAT(lv.scheduled_date, ' ', lv.scheduled_time_start) AS DATETIME)"
	scheduledEndExpr := "CAST(CONCAT(lv.scheduled_date, ' ', lv.scheduled_time_end) AS DATETIME)"

	// Filter by listing identity (exact match)
	if filter.ListingIdentityID != nil {
		conditions = append(conditions, "lv.listing_identity_id = ?")
		args = append(args, *filter.ListingIdentityID)
	}

	// Filter by owner user (exact match)
	if filter.OwnerUserID != nil {
		conditions = append(conditions, "li.user_id = ?")
		args = append(args, *filter.OwnerUserID)
	}

	// Filter by requester user (exact match)
	if filter.RequesterUserID != nil {
		conditions = append(conditions, "lv.user_id = ?")
		args = append(args, *filter.RequesterUserID)
	}

	// Filter by status array (IN clause for multiple statuses)
	if len(filter.Statuses) > 0 {
		placeholders := make([]string, len(filter.Statuses))
		for i, status := range filter.Statuses {
			placeholders[i] = "?"
			args = append(args, string(status))
		}
		conditions = append(conditions, fmt.Sprintf("lv.status IN (%s)", strings.Join(placeholders, ",")))
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
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM listing_visits lv JOIN listing_identities li ON li.id = lv.listing_identity_id WHERE %s", where)
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
	baseSelect := fmt.Sprintf(visitWithListingSelectBase, scheduledStartExpr, scheduledEndExpr)
	query := fmt.Sprintf(`
		%s
		WHERE %s
		ORDER BY %s ASC
		LIMIT ? OFFSET ?
	`, baseSelect, where, scheduledStartExpr)

	params := append(make([]any, 0, len(args)+2), args...)
	params = append(params, limit, offset)

	rows, err := a.QueryContext(ctx, tx, "list_visits", query, params...)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.visit.list.query_error", "err", err)
		return listingmodel.VisitListResult{}, fmt.Errorf("query visits: %w", err)
	}
	defer rows.Close()

	visits := make([]listingmodel.VisitWithListing, 0)
	for rows.Next() {
		entry, scanErr := scanVisitWithListingRow(rows)
		if scanErr != nil {
			utils.SetSpanError(ctx, scanErr)
			logger.Error("mysql.visit.list.scan_error", "err", scanErr)
			return listingmodel.VisitListResult{}, fmt.Errorf("scan visit: %w", scanErr)
		}
		visits = append(visits, entry)
	}

	if err = rows.Err(); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.visit.list.rows_error", "err", err)
		return listingmodel.VisitListResult{}, fmt.Errorf("iterate visits: %w", err)
	}

	return listingmodel.VisitListResult{Visits: visits, Total: total}, nil
}
