package mysqlscheduleadapter

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	scheduleconverters "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/schedule/converters"
	scheduleentity "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/schedule/entities"
	schedulemodel "github.com/projeto-toq/toq_server/internal/core/model/schedule_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// ListEntriesBetween returns agenda entries overlapping the half-open interval [from, to).
//
// Parameters:
//   - ctx: request-scoped context for tracing/logging.
//   - tx: optional transaction for consistent reads.
//   - agendaID: target agenda identifier.
//   - from/to: time window boundaries (inclusive start, exclusive end).
//
// Returns: slice of AgendaEntryInterface (empty when none) or infrastructure errors; sql.ErrNoRows is not used here.
// Observability: tracer span, logger propagation, span error marking on infra failures.
func (a *ScheduleAdapter) ListEntriesBetween(ctx context.Context, tx *sql.Tx, agendaID uint64, from, to time.Time) ([]schedulemodel.AgendaEntryInterface, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, err
	}
	defer spanEnd()
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `SELECT id, agenda_id, entry_type, starts_at, ends_at, blocking, reason, visit_id, photo_booking_id FROM listing_agenda_entries WHERE agenda_id = ? AND ends_at > ? AND starts_at < ? ORDER BY starts_at`

	rows, queryErr := a.QueryContext(ctx, tx, "select", query, agendaID, from, to)
	if queryErr != nil {
		utils.SetSpanError(ctx, queryErr)
		logger.Error("mysql.schedule.list_entries_between.query_error", "agenda_id", agendaID, "err", queryErr)
		return nil, fmt.Errorf("query agenda entries between: %w", queryErr)
	}
	defer rows.Close()

	entries := make([]schedulemodel.AgendaEntryInterface, 0)
	for rows.Next() {
		var entryEntity scheduleentity.EntryEntity
		if scanErr := rows.Scan(&entryEntity.ID, &entryEntity.AgendaID, &entryEntity.EntryType, &entryEntity.StartsAt, &entryEntity.EndsAt, &entryEntity.Blocking, &entryEntity.Reason, &entryEntity.VisitID, &entryEntity.PhotoBookingID); scanErr != nil {
			utils.SetSpanError(ctx, scanErr)
			logger.Error("mysql.schedule.list_entries_between.scan_error", "agenda_id", agendaID, "err", scanErr)
			return nil, fmt.Errorf("scan agenda entry between: %w", scanErr)
		}
		entries = append(entries, scheduleconverters.EntryEntityToDomain(entryEntity))
	}

	if rowsErr := rows.Err(); rowsErr != nil {
		utils.SetSpanError(ctx, rowsErr)
		logger.Error("mysql.schedule.list_entries_between.rows_error", "agenda_id", agendaID, "err", rowsErr)
		return nil, fmt.Errorf("iterate agenda entries between: %w", rowsErr)
	}

	return entries, nil
}
