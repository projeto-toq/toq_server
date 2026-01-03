package sessionmysqladapter

import (
	"context"
	"database/sql"

	sessionrepository "github.com/projeto-toq/toq_server/internal/core/port/right/repository/session_repository"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// Ensure interface method is implemented
var _ sessionrepository.SessionRepoPortInterface = (*SessionAdapter)(nil)

// DeleteExpiredSessions removes expired sessions (sliding or absolute) regardless of revoked flag.
//
// Behavior:
//   - Expires when expires_at < UTC_TIMESTAMP() OR absolute_expires_at < UTC_TIMESTAMP()
//   - Uses LIMIT to cap cleanup batch size (for cron/maintenance jobs)
//   - Returns number of rows deleted; logs RowsAffected errors as warnings
//
// Parameters:
//   - ctx: Tracing/logging context
//   - tx: Optional transaction (nil for standalone cron job)
//   - limit: Maximum rows to delete in this call
//
// Returns:
//   - int64: Rows deleted (best-effort if RowsAffected errors occur)
//   - error: Infrastructure errors only (query/connection/tx)
func (sa *SessionAdapter) DeleteExpiredSessions(ctx context.Context, tx *sql.Tx, limit int) (int64, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return 0, err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `DELETE FROM sessions 
		WHERE ((expires_at IS NOT NULL AND expires_at < UTC_TIMESTAMP()) 
			OR (absolute_expires_at IS NOT NULL AND absolute_expires_at < UTC_TIMESTAMP())) 
		LIMIT ?`

	res, err := sa.ExecContext(ctx, tx, "delete_expired_sessions", query, limit)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.session.delete_expired_sessions.exec_error", "limit", limit, "error", err)
		return 0, err
	}

	n, err := res.RowsAffected()
	if err != nil {
		logger.Warn("mysql.session.delete_expired_sessions.rows_affected_error", "error", err)
	}
	return n, nil
}
