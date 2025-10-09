package sessionmysqladapter

import (
	"context"
	"database/sql"

	sessionrepository "github.com/projeto-toq/toq_server/internal/core/port/right/repository/session_repository"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// Ensure interface method is implemented
var _ sessionrepository.SessionRepoPortInterface = (*SessionAdapter)(nil)

// DeleteExpiredSessions deletes sessions that are revoked and past expiry or absolutely expired; returns affected rows
func (sa *SessionAdapter) DeleteExpiredSessions(ctx context.Context, tx *sql.Tx, limit int) (int64, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return 0, err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	// Execute within provided transaction when available; fallback to direct DB exec otherwise
	execer := interface {
		ExecContext(context.Context, string, ...any) (sql.Result, error)
	}(sa.db.DB)
	if tx != nil {
		execer = tx
	}

	res, err := execer.ExecContext(ctx, `DELETE FROM sessions 
					WHERE ((expires_at IS NOT NULL AND expires_at < UTC_TIMESTAMP()) 
						OR (absolute_expires_at IS NOT NULL AND absolute_expires_at < UTC_TIMESTAMP())) 
					LIMIT ?`, limit)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.session.delete_expired_sessions.exec_error", "limit", limit, "error", err)
		return 0, err
	}
	n, _ := res.RowsAffected()
	return n, nil
}
