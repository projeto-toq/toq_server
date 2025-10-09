package sessionmysqladapter

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// DeleteSessionsByUserID permanently removes all sessions for a given user
func (sa *SessionAdapter) DeleteSessionsByUserID(ctx context.Context, tx *sql.Tx, userID int64) error {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `DELETE FROM sessions WHERE user_id = ?`

	// Use helper when tx is provided
	if tx != nil {
		if _, err := sa.Delete(ctx, tx, query, userID); err != nil {
			utils.SetSpanError(ctx, err)
			logger.Error("mysql.session.delete_sessions_by_user_id.delete_error", "user_id", userID, "error", err)
			return fmt.Errorf("delete sessions by user id: %w", err)
		}
		return nil
	}

	// Fallback if no transaction (should rarely happen in our flows)
	res, err := sa.db.DB.ExecContext(ctx, query, userID)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.session.delete_sessions_by_user_id.exec_error", "user_id", userID, "error", err)
		return fmt.Errorf("delete sessions by user id: %w", err)
	}
	if _, err := res.RowsAffected(); err != nil {
		logger.Warn("mysql.session.delete_sessions_by_user_id.rows_affected_error", "user_id", userID, "error", err)
	}
	return nil
}
