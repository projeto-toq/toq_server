package sessionmysqladapter

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// DeleteSessionsByUserID permanently removes all sessions for a given user (active or revoked).
//
// Behavior:
//   - Hard deletes rows; cannot be undone
//   - Accepts optional transaction; always uses provided tx when not nil
//   - RowsAffected is logged but not bubbled up
//
// Parameters:
//   - ctx: Tracing/logging context
//   - tx: Optional transaction
//   - userID: Owner user ID
//
// Returns:
//   - error: Infrastructure errors only; missing rows treated as no-op
func (sa *SessionAdapter) DeleteSessionsByUserID(ctx context.Context, tx *sql.Tx, userID int64) error {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `DELETE FROM sessions WHERE user_id = ?`

	// Use executor when tx is provided
	if tx != nil {
		result, execErr := sa.ExecContext(ctx, tx, "delete", query, userID)
		if execErr != nil {
			utils.SetSpanError(ctx, execErr)
			logger.Error("mysql.session.delete_sessions_by_user_id.exec_error", "user_id", userID, "err", execErr)
			return fmt.Errorf("delete sessions by user id: %w", execErr)
		}
		if _, rowsErr := result.RowsAffected(); rowsErr != nil {
			logger.Warn("mysql.session.delete_sessions_by_user_id.rows_affected_error", "user_id", userID, "err", rowsErr)
		}
		return nil
	}

	// Fallback if no transaction (should rarely happen in our flows)
	result, execErr := sa.ExecContext(ctx, nil, "delete", query, userID)
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.session.delete_sessions_by_user_id.exec_error", "user_id", userID, "err", execErr)
		return fmt.Errorf("delete sessions by user id: %w", execErr)
	}
	if _, rowsErr := result.RowsAffected(); rowsErr != nil {
		logger.Warn("mysql.session.delete_sessions_by_user_id.rows_affected_error", "user_id", userID, "err", rowsErr)
	}
	return nil
}
