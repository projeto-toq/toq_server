package mysqlphotosessionadapter

import (
	"context"
	"database/sql"
	"fmt"

	photosessionmodel "github.com/projeto-toq/toq_server/internal/core/model/photo_session_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// UpdateTimeOff updates an existing photographer time-off entry.
func (a *PhotoSessionAdapter) UpdateTimeOff(ctx context.Context, tx *sql.Tx, timeOff photosessionmodel.PhotographerTimeOffInterface) error {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `
        UPDATE photographer_time_off
        SET start_date = ?, end_date = ?, reason = ?
        WHERE id = ?
    `

	res, execErr := tx.ExecContext(ctx, query, timeOff.StartDate(), timeOff.EndDate(), timeOff.Reason(), timeOff.ID())
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.photo_session.update_time_off.exec_error", "time_off_id", timeOff.ID(), "err", execErr)
		return fmt.Errorf("update photographer time off: %w", execErr)
	}

	rows, rowsErr := res.RowsAffected()
	if rowsErr != nil {
		utils.SetSpanError(ctx, rowsErr)
		logger.Error("mysql.photo_session.update_time_off.rows_error", "time_off_id", timeOff.ID(), "err", rowsErr)
		return fmt.Errorf("rows affected on update photographer time off: %w", rowsErr)
	}

	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}
