package sessionmysqladapter

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (sa *SessionAdapter) RevokeSession(ctx context.Context, tx *sql.Tx, id int64) error {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `UPDATE sessions SET revoked = true WHERE id = ?`

	result, execErr := sa.ExecContext(ctx, tx, "update", query, id)
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.session.revoke_session.exec_error", "session_id", id, "err", execErr)
		return fmt.Errorf("revoke session: %w", execErr)
	}

	if _, rowsErr := result.RowsAffected(); rowsErr != nil {
		utils.SetSpanError(ctx, rowsErr)
		logger.Error("mysql.session.revoke_session.rows_error", "session_id", id, "err", rowsErr)
		return fmt.Errorf("revoke session rows affected: %w", rowsErr)
	}

	return nil
}
