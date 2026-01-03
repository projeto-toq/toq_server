package sessionmysqladapter

import (
	"context"
	"database/sql"
	"fmt"

	sessionmodel "github.com/projeto-toq/toq_server/internal/core/model/session_model"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// GetActiveSessionsByUserID lists all active (non-revoked, non-expired) sessions for a user.
//
// Behavior:
//   - Applies revoked = false and expires_at > UTC_TIMESTAMP()
//   - Returns empty slice when no active sessions exist (no error)
//   - Uses shared mapper for consistent NULL handling and conversions
//
// Parameters:
//   - ctx: Tracing/logging context
//   - tx: Optional transaction
//   - userID: Owner user ID (FK users.id)
//
// Returns:
//   - sessions: Slice of domain sessions (can be empty)
//   - error: Infrastructure errors only; sql.ErrNoRows is not used for list operations
//
// Observability:
//   - Starts span, logs query/scan errors, marks span on infra failures
func (sa *SessionAdapter) GetActiveSessionsByUserID(ctx context.Context, tx *sql.Tx, userID int64) (sessions []sessionmodel.SessionInterface, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `SELECT id, user_id, refresh_hash, token_jti, expires_at, absolute_expires_at, created_at, rotated_at, user_agent, ip, device_id, rotation_counter, last_refresh_at, revoked 
			FROM sessions 
			WHERE user_id = ? AND revoked = false AND expires_at > UTC_TIMESTAMP()`

	rows, err := sa.QueryContext(ctx, tx, "get_active_sessions_by_user_id", query, userID)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.session.get_active_sessions_by_user_id.query_error", "user_id", userID, "error", err)
		return nil, fmt.Errorf("get active sessions by user id: %w", err)
	}
	defer rows.Close()

	sessions = make([]sessionmodel.SessionInterface, 0)
	for rows.Next() {
		session, mapErr := sa.mapSessionFromScanner(ctx, rows, "get_active_sessions_by_user_id")
		if mapErr != nil {
			utils.SetSpanError(ctx, mapErr)
			logger.Error("mysql.session.get_active_sessions_by_user_id.scan_error", "user_id", userID, "error", mapErr)
			return nil, fmt.Errorf("get active sessions by user id: %w", mapErr)
		}
		sessions = append(sessions, session)
	}

	if err = rows.Err(); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.session.get_active_sessions_by_user_id.rows_error", "user_id", userID, "error", err)
		return nil, fmt.Errorf("get active sessions by user id: %w", err)
	}

	return
}
