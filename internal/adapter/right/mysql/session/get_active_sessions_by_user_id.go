package sessionmysqladapter

import (
	"context"
	"database/sql"
	"fmt"

	sessionconverters "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql/session/converters"
	sessionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/session_model"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

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

	entities, err := sa.Read(ctx, tx, query, userID)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.session.get_active_sessions_by_user_id.read_error", "user_id", userID, "error", err)
		return nil, fmt.Errorf("get active sessions by user id: %w", err)
	}

	sessions = make([]sessionmodel.SessionInterface, 0, len(entities))
	for _, entity := range entities {
		session, err := sessionconverters.SessionEntityToDomain(entity)
		if err != nil {
			utils.SetSpanError(ctx, err)
			logger.Error("mysql.session.get_active_sessions_by_user_id.convert_error", "user_id", userID, "error", err)
			return nil, err
		}
		sessions = append(sessions, session)
	}

	return
}
