package sessionmysqladapter

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (sa *SessionAdapter) UpdateSessionRotation(ctx context.Context, tx *sql.Tx, id int64, rotationCounter int, lastRefreshAt time.Time) error {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `UPDATE sessions SET rotation_counter = ?, last_refresh_at = ? WHERE id = ?`

	result, execErr := sa.ExecContext(ctx, tx, "update", query, rotationCounter, lastRefreshAt, id)
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.session.update_session_rotation.exec_error", "session_id", id, "err", execErr)
		return fmt.Errorf("update session rotation: %w", execErr)
	}

	if _, rowsErr := result.RowsAffected(); rowsErr != nil {
		utils.SetSpanError(ctx, rowsErr)
		logger.Error("mysql.session.update_session_rotation.rows_error", "session_id", id, "err", rowsErr)
		return fmt.Errorf("update session rotation rows affected: %w", rowsErr)
	}

	return nil
}
