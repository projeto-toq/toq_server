package mysqlscheduleadapter

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/projeto-toq/toq_server/internal/adapter/right/mysql/schedule/converters"
	"github.com/projeto-toq/toq_server/internal/adapter/right/mysql/schedule/entity"
	schedulemodel "github.com/projeto-toq/toq_server/internal/core/model/schedule_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

const ownerSummaryMaxPageSize = 50

func (a *ScheduleAdapter) ListOwnerSummary(ctx context.Context, tx *sql.Tx, filter schedulemodel.OwnerSummaryFilter) (schedulemodel.OwnerSummaryResult, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return schedulemodel.OwnerSummaryResult{}, err
	}
	defer spanEnd()
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	conditions := []string{"a.owner_id = ?"}
	args := []any{filter.OwnerID}

	if !filter.Range.From.IsZero() {
		conditions = append(conditions, "e.ends_at > ?")
		args = append(args, filter.Range.From)
	}
	if !filter.Range.To.IsZero() {
		conditions = append(conditions, "e.starts_at < ?")
		args = append(args, filter.Range.To)
	}

	if len(filter.ListingIDs) > 0 {
		placeholders := make([]string, len(filter.ListingIDs))
		for i, id := range filter.ListingIDs {
			placeholders[i] = "?"
			args = append(args, id)
		}
		conditions = append(conditions, fmt.Sprintf("a.listing_id IN (%s)", strings.Join(placeholders, ",")))
	}

	where := strings.Join(conditions, " AND ")

	countQuery := fmt.Sprintf("SELECT COUNT(DISTINCT a.listing_id) FROM listing_agenda_entries e INNER JOIN listing_agendas a ON a.id = e.agenda_id WHERE %s", where)

	var total int64
	countStmt, cleanup, prepareErr := a.PrepareContext(ctx, tx, "select", countQuery)
	if prepareErr != nil {
		utils.SetSpanError(ctx, prepareErr)
		logger.Error("mysql.schedule.owner_summary.prepare_count_error", "err", prepareErr)
		return schedulemodel.OwnerSummaryResult{}, fmt.Errorf("prepare owner summary count: %w", prepareErr)
	}
	defer cleanup()

	if scanErr := countStmt.QueryRowContext(ctx, args...).Scan(&total); scanErr != nil {
		utils.SetSpanError(ctx, scanErr)
		logger.Error("mysql.schedule.owner_summary.count_error", "owner_id", filter.OwnerID, "err", scanErr)
		return schedulemodel.OwnerSummaryResult{}, fmt.Errorf("count owner summary listings: %w", scanErr)
	}

	limit, offset := defaultPagination(filter.Pagination.Limit, filter.Pagination.Page, ownerSummaryMaxPageSize)

	listQuery := fmt.Sprintf(`
		SELECT DISTINCT a.listing_id
		FROM listing_agenda_entries e
		INNER JOIN listing_agendas a ON a.id = e.agenda_id
		WHERE %s
		ORDER BY a.listing_id
		LIMIT ? OFFSET ?
	`, where)

	listArgs := append(append([]any{}, args...), limit, offset)

	listingRows, listErr := a.QueryContext(ctx, tx, "select", listQuery, listArgs...)
	if listErr != nil {
		utils.SetSpanError(ctx, listErr)
		logger.Error("mysql.schedule.owner_summary.listings_query_error", "owner_id", filter.OwnerID, "err", listErr)
		return schedulemodel.OwnerSummaryResult{}, fmt.Errorf("query owner summary listings: %w", listErr)
	}
	defer listingRows.Close()

	listingIDs := make([]int64, 0)
	for listingRows.Next() {
		var id int64
		if scanErr := listingRows.Scan(&id); scanErr != nil {
			utils.SetSpanError(ctx, scanErr)
			logger.Error("mysql.schedule.owner_summary.listings_scan_error", "owner_id", filter.OwnerID, "err", scanErr)
			return schedulemodel.OwnerSummaryResult{}, fmt.Errorf("scan owner summary listing: %w", scanErr)
		}
		listingIDs = append(listingIDs, id)
	}

	if rowsErr := listingRows.Err(); rowsErr != nil {
		utils.SetSpanError(ctx, rowsErr)
		logger.Error("mysql.schedule.owner_summary.listings_rows_error", "owner_id", filter.OwnerID, "err", rowsErr)
		return schedulemodel.OwnerSummaryResult{}, fmt.Errorf("iterate owner summary listings: %w", rowsErr)
	}

	if len(listingIDs) == 0 {
		return schedulemodel.OwnerSummaryResult{Items: []schedulemodel.OwnerSummaryItem{}, Total: total}, nil
	}

	entriesResult, entriesErr := a.fetchSummaryEntries(ctx, tx, listingIDs, filter)
	if entriesErr != nil {
		return schedulemodel.OwnerSummaryResult{}, entriesErr
	}

	items := make([]schedulemodel.OwnerSummaryItem, 0, len(listingIDs))
	for _, id := range listingIDs {
		entries := entriesResult[id]
		items = append(items, schedulemodel.OwnerSummaryItem{ListingID: id, Entries: entries})
	}

	return schedulemodel.OwnerSummaryResult{Items: items, Total: total}, nil
}

func (a *ScheduleAdapter) fetchSummaryEntries(ctx context.Context, tx *sql.Tx, listingIDs []int64, filter schedulemodel.OwnerSummaryFilter) (map[int64][]schedulemodel.SummaryEntry, error) {
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	placeholders := make([]string, len(listingIDs))
	args := make([]any, 0, len(listingIDs)+2)
	for i, id := range listingIDs {
		placeholders[i] = "?"
		args = append(args, id)
	}

	conditions := []string{fmt.Sprintf("a.listing_id IN (%s)", strings.Join(placeholders, ","))}
	if !filter.Range.From.IsZero() {
		conditions = append(conditions, "e.ends_at > ?")
		args = append(args, filter.Range.From)
	}
	if !filter.Range.To.IsZero() {
		conditions = append(conditions, "e.starts_at < ?")
		args = append(args, filter.Range.To)
	}

	where := strings.Join(conditions, " AND ")

	query := fmt.Sprintf(`
		SELECT e.id, e.agenda_id, e.entry_type, e.starts_at, e.ends_at, e.blocking, e.reason, e.visit_id, e.photo_booking_id, a.listing_id
		FROM listing_agenda_entries e
		INNER JOIN listing_agendas a ON a.id = e.agenda_id
		WHERE %s
		ORDER BY a.listing_id, e.starts_at
	`, where)

	rows, queryErr := a.QueryContext(ctx, tx, "select", query, args...)
	if queryErr != nil {
		utils.SetSpanError(ctx, queryErr)
		logger.Error("mysql.schedule.owner_summary.entries_query_error", "err", queryErr)
		return nil, fmt.Errorf("query owner summary entries: %w", queryErr)
	}
	defer rows.Close()

	result := make(map[int64][]schedulemodel.SummaryEntry, len(listingIDs))

	for rows.Next() {
		var entryEntity entity.EntryEntity
		var listingID int64
		if scanErr := rows.Scan(
			&entryEntity.ID,
			&entryEntity.AgendaID,
			&entryEntity.EntryType,
			&entryEntity.StartsAt,
			&entryEntity.EndsAt,
			&entryEntity.Blocking,
			&entryEntity.Reason,
			&entryEntity.VisitID,
			&entryEntity.PhotoBookingID,
			&listingID,
		); scanErr != nil {
			utils.SetSpanError(ctx, scanErr)
			logger.Error("mysql.schedule.owner_summary.entries_scan_error", "err", scanErr)
			return nil, fmt.Errorf("scan owner summary entry: %w", scanErr)
		}

		entry := converters.ToEntryModel(entryEntity)
		summary := schedulemodel.SummaryEntry{
			EntryType: entry.EntryType(),
			StartsAt:  entry.StartsAt(),
			EndsAt:    entry.EndsAt(),
			Blocking:  entry.Blocking(),
		}
		result[listingID] = append(result[listingID], summary)
	}

	if rowsErr := rows.Err(); rowsErr != nil {
		utils.SetSpanError(ctx, rowsErr)
		logger.Error("mysql.schedule.owner_summary.entries_rows_error", "err", rowsErr)
		return nil, fmt.Errorf("iterate owner summary entries: %w", rowsErr)
	}

	return result, nil
}
