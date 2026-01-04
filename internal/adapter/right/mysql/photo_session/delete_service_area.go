package mysqlphotosessionadapter

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// DeleteServiceArea removes a service area entry.
func (a *PhotoSessionAdapter) DeleteServiceArea(ctx context.Context, tx *sql.Tx, areaID uint64) error {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return err
	}
	defer spanEnd()
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `DELETE FROM photographer_service_areas WHERE id = ?`

	result, execErr := a.ExecContext(ctx, tx, "delete", query, areaID)
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.photo_session.service_area.delete.exec_error", "area_id", areaID, "err", execErr)
		return fmt.Errorf("delete service area: %w", execErr)
	}

	affected, rowsErr := result.RowsAffected()
	if rowsErr != nil {
		utils.SetSpanError(ctx, rowsErr)
		logger.Error("mysql.photo_session.service_area.delete.rows_affected_error", "area_id", areaID, "err", rowsErr)
		return fmt.Errorf("delete service area rows affected: %w", rowsErr)
	}

	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
