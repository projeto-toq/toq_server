package mysqlphotosessionadapter

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (a *PhotoSessionAdapter) DeleteEntryByID(ctx context.Context, tx *sql.Tx, entryID uint64) error {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return err
	}
	defer spanEnd()
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	result, execErr := a.ExecContext(ctx, tx, "delete", `DELETE FROM photographer_agenda_entries WHERE id = ?`, entryID)
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.photo_session.delete_entry.exec_error", "entry_id", entryID, "err", execErr)
		return fmt.Errorf("delete agenda entry: %w", execErr)
	}

	affected, rowsErr := result.RowsAffected()
	if rowsErr != nil {
		utils.SetSpanError(ctx, rowsErr)
		logger.Error("mysql.photo_session.delete_entry.rows_error", "entry_id", entryID, "err", rowsErr)
		return fmt.Errorf("rows affected agenda entry: %w", rowsErr)
	}

	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
