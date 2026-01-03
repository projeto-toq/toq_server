package sessionmysqladapter

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// MarkSessionRotated sets rotated_at to UTC_TIMESTAMP() for a session ID.
//
// Behavior:
//   - Used during refresh flows to timestamp the latest rotation
//   - Does not change rotation_counter or revoke status
//
// Parameters:
//   - ctx: Tracing/logging context
//   - tx: Optional transaction
//   - id: Session primary key
//
// Returns:
//   - error: Infrastructure errors; rows affected errors surfaced; missing rows treated as no-op
func (sa *SessionAdapter) MarkSessionRotated(ctx context.Context, tx *sql.Tx, id int64) error {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `UPDATE sessions SET rotated_at = UTC_TIMESTAMP() WHERE id = ?`

	result, execErr := sa.ExecContext(ctx, tx, "update", query, id)
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.session.mark_session_rotated.exec_error", "session_id", id, "err", execErr)
		return fmt.Errorf("mark session rotated: %w", execErr)
	}

	if _, rowsErr := result.RowsAffected(); rowsErr != nil {
		utils.SetSpanError(ctx, rowsErr)
		logger.Error("mysql.session.mark_session_rotated.rows_error", "session_id", id, "err", rowsErr)
		return fmt.Errorf("mark session rotated rows affected: %w", rowsErr)
	}

	return nil
}
