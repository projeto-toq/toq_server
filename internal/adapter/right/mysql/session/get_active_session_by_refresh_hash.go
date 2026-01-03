package sessionmysqladapter

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	sessionmodel "github.com/projeto-toq/toq_server/internal/core/model/session_model"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// GetActiveSessionByRefreshHash returns the active (non-revoked, non-expired) session for a refresh hash.
//
// Behavior:
//   - Filters on revoked = false and expires_at > UTC_TIMESTAMP() (sliding expiration)
//   - Does NOT check absolute_expires_at here; absolute check is enforced in services
//   - Returns sql.ErrNoRows when no active session matches (not found or expired)
//
// Parameters:
//   - ctx: Tracing/logging context
//   - tx: Optional transaction
//   - hash: SHA-256 refresh token hash (unique)
//
// Returns:
//   - session: Domain model populated from DB
//   - error: sql.ErrNoRows when absent/expired; wrapped infra errors otherwise
//
// Observability:
//   - Starts span via utils.GenerateTracer and propagates contextual logger
//   - Logs debug on not found; logs error and marks span on infra failures
func (sa *SessionAdapter) GetActiveSessionByRefreshHash(ctx context.Context, tx *sql.Tx, hash string) (session sessionmodel.SessionInterface, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `SELECT id, user_id, refresh_hash, token_jti, expires_at, absolute_expires_at, created_at, rotated_at, user_agent, ip, device_id, rotation_counter, last_refresh_at, revoked 
			FROM sessions 
			WHERE refresh_hash = ? AND revoked = false AND expires_at > UTC_TIMESTAMP()`

	row := sa.QueryRowContext(ctx, tx, "get_active_session_by_refresh_hash", query, hash)
	session, err = sa.mapSessionFromScanner(ctx, row, "get_active_session_by_refresh_hash")
	if errors.Is(err, sql.ErrNoRows) {
		logger.Debug("mysql.session.get_active_session_by_refresh_hash.not_found", "hash", hash)
		return nil, sql.ErrNoRows
	}
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.session.get_active_session_by_refresh_hash.scan_error", "hash", hash, "error", err)
		return nil, fmt.Errorf("get active session by refresh hash: %w", err)
	}

	return
}
