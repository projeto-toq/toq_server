package mysqlphotosessionadapter

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (a *PhotoSessionAdapter) UpdateEntrySourceID(ctx context.Context, tx *sql.Tx, entryID uint64, sourceID uint64) error {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	result, execErr := a.ExecContext(ctx, tx, "update", `UPDATE photographer_agenda_entries SET source_id = ? WHERE id = ?`, sourceID, entryID)
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.photo_session.update_entry_source.exec_error", "entry_id", entryID, "err", execErr)
		return fmt.Errorf("update agenda entry source id: %w", execErr)
	}

	affected, rowsErr := result.RowsAffected()
	if rowsErr != nil {
		utils.SetSpanError(ctx, rowsErr)
		logger.Error("mysql.photo_session.update_entry_source.rows_error", "entry_id", entryID, "err", rowsErr)
		return fmt.Errorf("rows affected agenda entry source id: %w", rowsErr)
	}

	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
