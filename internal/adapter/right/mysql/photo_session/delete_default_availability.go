package mysqlphotosessionadapter

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	photosessionmodel "github.com/projeto-toq/toq_server/internal/core/model/photo_session_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// DeleteDefaultAvailability removes default availability optionally filtered by weekday and period.
func (a *PhotoSessionAdapter) DeleteDefaultAvailability(ctx context.Context, tx *sql.Tx, photographerID uint64, weekday *time.Weekday, period *photosessionmodel.SlotPeriod) error {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `DELETE FROM photographer_default_availability WHERE photographer_user_id = ?`
	args := []any{photographerID}

	if weekday != nil {
		query += " AND weekday = ?"
		args = append(args, int(*weekday))
	}

	if period != nil {
		query += " AND period = ?"
		args = append(args, string(*period))
	}

	result, execErr := tx.ExecContext(ctx, query, args...)
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.photo_session.delete_default_availability.exec_error", "err", execErr)
		return fmt.Errorf("delete photographer default availability: %w", execErr)
	}

	affected, _ := result.RowsAffected()
	logger.Info("mysql.photo_session.delete_default_availability.success", "photographer_id", photographerID, "affected", affected)

	return nil
}
