package mysqlscheduleadapter

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/projeto-toq/toq_server/internal/adapter/right/mysql/schedule/converters"
	schedulemodel "github.com/projeto-toq/toq_server/internal/core/model/schedule_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (a *ScheduleAdapter) InsertEntry(ctx context.Context, tx *sql.Tx, entry schedulemodel.AgendaEntryInterface) (uint64, error) {
	ctx, spanEnd, err := withTracer(ctx)
	if err != nil {
		return 0, err
	}
	if spanEnd != nil {
		defer spanEnd()
	}

	exec := a.executor(tx)
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	entity := converters.ToEntryEntity(entry)

	query := `INSERT INTO listing_agenda_entries (agenda_id, entry_type, starts_at, ends_at, blocking, reason, visit_id, photo_booking_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	result, err := exec.ExecContext(ctx, query, entity.AgendaID, entity.EntryType, entity.StartsAt, entity.EndsAt, entity.Blocking, entity.Reason, entity.VisitID, entity.PhotoBookingID)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.schedule.insert_entry.exec_error", "agenda_id", entity.AgendaID, "err", err)
		return 0, fmt.Errorf("insert agenda entry: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.schedule.insert_entry.last_id_error", "agenda_id", entity.AgendaID, "err", err)
		return 0, fmt.Errorf("agenda entry last insert id: %w", err)
	}

	entry.SetID(uint64(id))
	return uint64(id), nil
}
