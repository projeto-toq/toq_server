package mysqlvisitadapter

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/projeto-toq/toq_server/internal/adapter/right/mysql/visit/converters"
	"github.com/projeto-toq/toq_server/internal/adapter/right/mysql/visit/entities"
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
	query := fmt.Sprintf(`
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
			active.zip_code,
			active.street,
			active.number,
			active.complement,
			active.neighborhood,
			active.city,
			active.state,
			active.title,
			active.description
		FROM listing_visits lv
		JOIN listing_identities li ON li.id = lv.listing_identity_id
		LEFT JOIN listing_versions active ON active.id = li.active_version_id AND active.deleted = 0
		WHERE %s
		ORDER BY %s ASC
		LIMIT ? OFFSET ?
	`, scheduledStartExpr, scheduledEndExpr, where, scheduledStartExpr)

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
		var visitEntity entities.VisitEntity
		var listingID sql.NullInt64
		var listingVersion sql.NullInt64
		var listingZip sql.NullString
		var listingStreet sql.NullString
		var listingNumber sql.NullString
		var listingComplement sql.NullString
		var listingNeighborhood sql.NullString
		var listingCity sql.NullString
		var listingState sql.NullString
		var listingTitle sql.NullString
		var listingDescription sql.NullString

		if err = rows.Scan(
			&visitEntity.ID,
			&visitEntity.ListingIdentityID,
			&visitEntity.ListingVersion,
			&visitEntity.RequesterUserID,
			&visitEntity.OwnerUserID,
			&visitEntity.ScheduledStart,
			&visitEntity.ScheduledEnd,
			&visitEntity.Status,
			&visitEntity.Source,
			&visitEntity.Notes,
			&visitEntity.RejectionReason,
			&visitEntity.FirstOwnerActionAt,
			&visitEntity.RequestedAt,
			&listingID,
			&listingVersion,
			&listingZip,
			&listingStreet,
			&listingNumber,
			&listingComplement,
			&listingNeighborhood,
			&listingCity,
			&listingState,
			&listingTitle,
			&listingDescription,
		); err != nil {
			utils.SetSpanError(ctx, err)
			logger.Error("mysql.visit.list.scan_error", "err", err)
			return listingmodel.VisitListResult{}, fmt.Errorf("scan visit: %w", err)
		}

		if !listingID.Valid || !listingVersion.Valid {
			errMissing := fmt.Errorf("listing version not found for listing_identity_id=%d", visitEntity.ListingIdentityID)
			utils.SetSpanError(ctx, errMissing)
			logger.Error("mysql.visit.list.missing_listing_version", "listing_identity_id", visitEntity.ListingIdentityID)
			return listingmodel.VisitListResult{}, errMissing
		}

		requiredSnapshotFields := []struct {
			valid bool
			name  string
		}{
			{listingZip.Valid, "zip_code"},
			{listingStreet.Valid, "street"},
			{listingNeighborhood.Valid, "neighborhood"},
			{listingCity.Valid, "city"},
			{listingState.Valid, "state"},
		}

		for _, field := range requiredSnapshotFields {
			if field.valid {
				continue
			}
			errMissing := fmt.Errorf("listing snapshot missing %s for listing_identity_id=%d", field.name, visitEntity.ListingIdentityID)
			utils.SetSpanError(ctx, errMissing)
			logger.Error("mysql.visit.list.missing_listing_field", "listing_identity_id", visitEntity.ListingIdentityID, "field", field.name)
			return listingmodel.VisitListResult{}, errMissing
		}

		visitModel := converters.ToVisitModel(visitEntity)
		listingModel := listingmodel.NewListing()
		listingModel.SetID(listingID.Int64)
		listingModel.SetListingIdentityID(visitEntity.ListingIdentityID)
		listingModel.SetVersion(uint8(listingVersion.Int64))
		listingModel.SetZipCode(listingZip.String)
		listingModel.SetStreet(listingStreet.String)
		if listingNumber.Valid {
			listingModel.SetNumber(listingNumber.String)
		}
		if listingComplement.Valid && strings.TrimSpace(listingComplement.String) != "" {
			listingModel.SetComplement(strings.TrimSpace(listingComplement.String))
		}
		listingModel.SetNeighborhood(listingNeighborhood.String)
		listingModel.SetCity(listingCity.String)
		listingModel.SetState(listingState.String)
		if listingTitle.Valid {
			listingModel.SetTitle(listingTitle.String)
		} else {
			listingModel.UnsetTitle()
		}
		if listingDescription.Valid {
			listingModel.SetDescription(listingDescription.String)
		} else {
			listingModel.UnsetDescription()
		}

		visits = append(visits, listingmodel.VisitWithListing{Visit: visitModel, Listing: listingModel})
	}

	if err = rows.Err(); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.visit.list.rows_error", "err", err)
		return listingmodel.VisitListResult{}, fmt.Errorf("iterate visits: %w", err)
	}

	return listingmodel.VisitListResult{Visits: visits, Total: total}, nil
}
