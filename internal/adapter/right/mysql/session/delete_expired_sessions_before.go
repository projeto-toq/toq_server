package sessionmysqladapter

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	sessionrepository "github.com/projeto-toq/toq_server/internal/core/port/right/repository/session_repository"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// Ensure interface compliance
var _ sessionrepository.SessionRepoPortInterface = (*SessionAdapter)(nil)

// DeleteExpiredSessionsBefore removes sessions whose effective expiry is older than cutoff.
// Uses the generated column effective_expires_at (COALESCE of absolute_expires_at and expires_at) for efficient pruning.
func (sa *SessionAdapter) DeleteExpiredSessionsBefore(ctx context.Context, tx *sql.Tx, cutoff time.Time, limit int) (int64, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return 0, err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	if limit <= 0 {
		limit = 500
	}

	query := `DELETE FROM sessions
        WHERE effective_expires_at IS NOT NULL AND effective_expires_at < ?
        LIMIT ?`

	res, execErr := sa.ExecContext(ctx, tx, "delete_expired_sessions_before", query, cutoff, limit)
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.session.delete_expired_sessions_before.exec_error", "cutoff", cutoff, "limit", limit, "error", execErr)
		return 0, fmt.Errorf("delete expired sessions by cutoff: %w", execErr)
	}

	rows, raErr := res.RowsAffected()
	if raErr != nil {
		logger.Warn("mysql.session.delete_expired_sessions_before.rows_affected_warning", "error", raErr)
		return 0, nil
	}

	if rows > 0 {
		logger.Debug("mysql.session.delete_expired_sessions_before.success", "deleted", rows, "cutoff", cutoff, "limit", limit)
	}
	return rows, nil
}
