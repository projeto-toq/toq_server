package mysqlphotosessionadapter

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// DeleteTimeOff removes a time-off entry by ID.
func (a *PhotoSessionAdapter) DeleteTimeOff(ctx context.Context, tx *sql.Tx, timeOffID uint64) error {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `
		DELETE FROM photographer_time_off
		WHERE id = ?
	`

	res, execErr := tx.ExecContext(ctx, query, timeOffID)
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.photo_session.delete_time_off.exec_error", "err", execErr)
		return fmt.Errorf("delete photographer time off: %w", execErr)
	}

	rows, rowsErr := res.RowsAffected()
	if rowsErr != nil {
		utils.SetSpanError(ctx, rowsErr)
		logger.Error("mysql.photo_session.delete_time_off.rows_error", "err", rowsErr)
		return fmt.Errorf("delete photographer time off rows affected: %w", rowsErr)
	}

	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}
