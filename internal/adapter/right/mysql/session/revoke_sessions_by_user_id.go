package sessionmysqladapter

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (sa *SessionAdapter) RevokeSessionsByUserID(ctx context.Context, tx *sql.Tx, userID int64) error {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `UPDATE sessions SET revoked = true WHERE user_id = ? AND revoked = false`

	result, execErr := sa.ExecContext(ctx, tx, "update", query, userID)
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.session.revoke_sessions_by_user_id.exec_error", "user_id", userID, "err", execErr)
		return fmt.Errorf("revoke sessions by user id: %w", execErr)
	}

	if _, rowsErr := result.RowsAffected(); rowsErr != nil {
		logger.Warn("mysql.session.revoke_sessions_by_user_id.rows_affected_error", "user_id", userID, "err", rowsErr)
	}

	return nil
}
