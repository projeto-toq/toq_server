package mysqlphotosessionadapter

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (a *PhotoSessionAdapter) UpdateEntrySourceID(ctx context.Context, tx *sql.Tx, entryID uint64, sourceID uint64) error {
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

	result, err := exec.ExecContext(ctx, `UPDATE photographer_agenda_entries SET source_id = ? WHERE id = ?`, sourceID, entryID)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.photo_session.update_entry_source.exec_error", "entry_id", entryID, "err", err)
		return fmt.Errorf("update agenda entry source id: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.photo_session.update_entry_source.rows_error", "entry_id", entryID, "err", err)
		return fmt.Errorf("rows affected agenda entry source id: %w", err)
	}

	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
