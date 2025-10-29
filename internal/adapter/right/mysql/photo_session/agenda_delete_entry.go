package mysqlphotosessionadapter

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (a *PhotoSessionAdapter) DeleteEntryByID(ctx context.Context, tx *sql.Tx, entryID uint64) error {
	ctx, spanEnd, err := withTracer(ctx)
	if err != nil {
		return err
	}
	if spanEnd != nil {
		defer spanEnd()
	}

	exec := a.executor(tx)
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	result, err := exec.ExecContext(ctx, `DELETE FROM photographer_agenda_entries WHERE id = ?`, entryID)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.photo_session.delete_entry.exec_error", "entry_id", entryID, "err", err)
		return fmt.Errorf("delete agenda entry: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.photo_session.delete_entry.rows_error", "entry_id", entryID, "err", err)
		return fmt.Errorf("rows affected agenda entry: %w", err)
	}

	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
