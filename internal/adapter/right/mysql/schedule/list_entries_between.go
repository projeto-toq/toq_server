package mysqlscheduleadapter

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/projeto-toq/toq_server/internal/adapter/right/mysql/schedule/converters"
	"github.com/projeto-toq/toq_server/internal/adapter/right/mysql/schedule/entity"
	schedulemodel "github.com/projeto-toq/toq_server/internal/core/model/schedule_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (a *ScheduleAdapter) ListEntriesBetween(ctx context.Context, tx *sql.Tx, agendaID uint64, from, to time.Time) ([]schedulemodel.AgendaEntryInterface, error) {
	ctx, spanEnd, err := withTracer(ctx)
	if err != nil {
		return nil, err
	}
	if spanEnd != nil {
		defer spanEnd()
	}

	exec := a.executor(tx)
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `SELECT id, agenda_id, entry_type, starts_at, ends_at, blocking, reason, visit_id, photo_booking_id FROM listing_agenda_entries WHERE agenda_id = ? AND ends_at > ? AND starts_at < ? ORDER BY starts_at`

	rows, err := exec.QueryContext(ctx, query, agendaID, from, to)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.schedule.list_entries_between.query_error", "agenda_id", agendaID, "err", err)
		return nil, fmt.Errorf("query agenda entries between: %w", err)
	}
	defer rows.Close()

	entries := make([]schedulemodel.AgendaEntryInterface, 0)
	for rows.Next() {
		var entryEntity entity.EntryEntity
		if err = rows.Scan(&entryEntity.ID, &entryEntity.AgendaID, &entryEntity.EntryType, &entryEntity.StartsAt, &entryEntity.EndsAt, &entryEntity.Blocking, &entryEntity.Reason, &entryEntity.VisitID, &entryEntity.PhotoBookingID); err != nil {
			utils.SetSpanError(ctx, err)
			logger.Error("mysql.schedule.list_entries_between.scan_error", "agenda_id", agendaID, "err", err)
			return nil, fmt.Errorf("scan agenda entry between: %w", err)
		}
		entries = append(entries, converters.ToEntryModel(entryEntity))
	}

	if err = rows.Err(); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.schedule.list_entries_between.rows_error", "agenda_id", agendaID, "err", err)
		return nil, fmt.Errorf("iterate agenda entries between: %w", err)
	}

	return entries, nil
}
