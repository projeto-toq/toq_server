package mysqlscheduleadapter

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	scheduleconverters "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/schedule/converters"
	scheduleentity "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/schedule/entities"
	schedulemodel "github.com/projeto-toq/toq_server/internal/core/model/schedule_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

const agendaDetailMaxPageSize = 100

// ListAgendaEntries returns paginated agenda entries for a listing within a date window.
//
// Parameters:
//   - ctx: request-scoped context for tracing/logging.
//   - tx: optional transaction for consistent reads.
//   - filter: owner/listing ids, time range, and pagination config.
//
// Returns: AgendaEntriesPage with entries slice (empty when none) and total count; infra errors bubbled.
// Observability: tracer span, logger propagation, span error marking on infra failures.
func (a *ScheduleAdapter) ListAgendaEntries(ctx context.Context, tx *sql.Tx, filter schedulemodel.AgendaDetailFilter) (schedulemodel.AgendaEntriesPage, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return schedulemodel.AgendaEntriesPage{}, err
	}
	defer spanEnd()
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	conditions := []string{"a.owner_id = ?", "a.listing_identity_id = ?"}
	args := []any{filter.OwnerID, filter.ListingIdentityID}

	if !filter.Range.From.IsZero() {
		conditions = append(conditions, "e.ends_at > ?")
		args = append(args, filter.Range.From)
	}
	if !filter.Range.To.IsZero() {
		conditions = append(conditions, "e.starts_at < ?")
		args = append(args, filter.Range.To)
	}

	where := strings.Join(conditions, " AND ")

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM listing_agenda_entries e INNER JOIN listing_agendas a ON a.id = e.agenda_id WHERE %s", where)

	var total int64
	if scanErr := a.QueryRowContext(ctx, tx, "select", countQuery, args...).Scan(&total); scanErr != nil {
		utils.SetSpanError(ctx, scanErr)
		logger.Error("mysql.schedule.agenda_detail.count_error", "listing_identity_id", filter.ListingIdentityID, "err", scanErr)
		return schedulemodel.AgendaEntriesPage{}, fmt.Errorf("count agenda entries: %w", scanErr)
	}

	limit, offset := defaultPagination(filter.Pagination.Limit, filter.Pagination.Page, agendaDetailMaxPageSize)

	query := fmt.Sprintf(`
		SELECT e.id, e.agenda_id, e.entry_type, e.starts_at, e.ends_at, e.blocking, e.reason, e.visit_id, e.photo_booking_id
		FROM listing_agenda_entries e
		INNER JOIN listing_agendas a ON a.id = e.agenda_id
		WHERE %s
		ORDER BY e.starts_at
		LIMIT ? OFFSET ?
	`, where)

	argsWithPagination := append(append([]any{}, args...), limit, offset)

	rows, queryErr := a.QueryContext(ctx, tx, "select", query, argsWithPagination...)
	if queryErr != nil {
		utils.SetSpanError(ctx, queryErr)
		logger.Error("mysql.schedule.agenda_detail.query_error", "listing_identity_id", filter.ListingIdentityID, "err", queryErr)
		return schedulemodel.AgendaEntriesPage{}, fmt.Errorf("query agenda entries: %w", queryErr)
	}
	defer rows.Close()

	entries := make([]schedulemodel.AgendaEntryInterface, 0)
	for rows.Next() {
		var entryEntity scheduleentity.EntryEntity
		if scanErr := rows.Scan(&entryEntity.ID, &entryEntity.AgendaID, &entryEntity.EntryType, &entryEntity.StartsAt, &entryEntity.EndsAt, &entryEntity.Blocking, &entryEntity.Reason, &entryEntity.VisitID, &entryEntity.PhotoBookingID); scanErr != nil {
			utils.SetSpanError(ctx, scanErr)
			logger.Error("mysql.schedule.agenda_detail.scan_error", "listing_identity_id", filter.ListingIdentityID, "err", scanErr)
			return schedulemodel.AgendaEntriesPage{}, fmt.Errorf("scan agenda entry: %w", scanErr)
		}
		entries = append(entries, scheduleconverters.EntryEntityToDomain(entryEntity))
	}

	if rowsErr := rows.Err(); rowsErr != nil {
		utils.SetSpanError(ctx, rowsErr)
		logger.Error("mysql.schedule.agenda_detail.rows_error", "listing_identity_id", filter.ListingIdentityID, "err", rowsErr)
		return schedulemodel.AgendaEntriesPage{}, fmt.Errorf("iterate agenda entries: %w", rowsErr)
	}

	return schedulemodel.AgendaEntriesPage{Entries: entries, Total: total}, nil
}
