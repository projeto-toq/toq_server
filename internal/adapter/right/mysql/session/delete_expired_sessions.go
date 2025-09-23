package sessionmysqladapter

import (
	"context"
	"database/sql"
	"log/slog"

	sessionrepository "github.com/giulio-alfieri/toq_server/internal/core/port/right/repository/session_repository"
)

// Ensure interface method is implemented
var _ sessionrepository.SessionRepoPortInterface = (*SessionAdapter)(nil)

// DeleteExpiredSessions deletes sessions that are revoked and past expiry or absolutely expired; returns affected rows
func (sa *SessionAdapter) DeleteExpiredSessions(ctx context.Context, tx *sql.Tx, limit int) (int64, error) {
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
		slog.Error("sessionmysqladapter/DeleteExpiredSessions: delete failed", "error", err)
		return 0, err
	}
	n, _ := res.RowsAffected()
	return n, nil
}
