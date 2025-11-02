package mysqlscheduleadapter

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/projeto-toq/toq_server/internal/adapter/right/mysql/schedule/converters"
	"github.com/projeto-toq/toq_server/internal/adapter/right/mysql/schedule/entity"
	schedulemodel "github.com/projeto-toq/toq_server/internal/core/model/schedule_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (a *ScheduleAdapter) GetEntryByID(ctx context.Context, tx *sql.Tx, entryID uint64) (schedulemodel.AgendaEntryInterface, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, err
	}
	defer spanEnd()
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `SELECT id, agenda_id, entry_type, starts_at, ends_at, blocking, reason, visit_id, photo_booking_id FROM listing_agenda_entries WHERE id = ?`
	row := a.QueryRowContext(ctx, tx, "select", query, entryID)

	var entryEntity entity.EntryEntity
	if err = row.Scan(&entryEntity.ID, &entryEntity.AgendaID, &entryEntity.EntryType, &entryEntity.StartsAt, &entryEntity.EndsAt, &entryEntity.Blocking, &entryEntity.Reason, &entryEntity.VisitID, &entryEntity.PhotoBookingID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.schedule.get_entry.scan_error", "entry_id", entryID, "err", err)
		return nil, fmt.Errorf("scan agenda entry: %w", err)
	}

	return converters.ToEntryModel(entryEntity), nil
}
