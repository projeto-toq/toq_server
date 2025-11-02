package mysqlscheduleadapter

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/projeto-toq/toq_server/internal/adapter/right/mysql/schedule/converters"
	schedulemodel "github.com/projeto-toq/toq_server/internal/core/model/schedule_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (a *ScheduleAdapter) UpdateEntry(ctx context.Context, tx *sql.Tx, entry schedulemodel.AgendaEntryInterface) error {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return err
	}
	defer spanEnd()
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	entity := converters.ToEntryEntity(entry)

	query := `UPDATE listing_agenda_entries SET entry_type = ?, starts_at = ?, ends_at = ?, blocking = ?, reason = ?, visit_id = ?, photo_booking_id = ? WHERE id = ?`
	result, execErr := a.ExecContext(ctx, tx, "update", query, entity.EntryType, entity.StartsAt, entity.EndsAt, entity.Blocking, entity.Reason, entity.VisitID, entity.PhotoBookingID, entity.ID)
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.schedule.update_entry.exec_error", "entry_id", entity.ID, "err", execErr)
		return fmt.Errorf("update agenda entry: %w", execErr)
	}

	affected, rowsErr := result.RowsAffected()
	if rowsErr != nil {
		utils.SetSpanError(ctx, rowsErr)
		logger.Error("mysql.schedule.update_entry.rows_error", "entry_id", entity.ID, "err", rowsErr)
		return fmt.Errorf("agenda entry rows affected: %w", rowsErr)
	}

	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
