package sessionmysqladapter

import (
	"context"
	"database/sql"
	"fmt"

	sessionconverters "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql/session/converters"
	sessionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/session_model"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

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

	entities, err := sa.Read(ctx, tx, query, hash)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.session.get_active_session_by_refresh_hash.read_error", "hash", hash, "error", err)
		return nil, fmt.Errorf("get active session by refresh hash: %w", err)
	}

	if len(entities) == 0 {
		logger.Debug("mysql.session.get_active_session_by_refresh_hash.not_found", "hash", hash)
		return nil, sql.ErrNoRows
	}

	session, err = sessionconverters.SessionEntityToDomain(entities[0])
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.session.get_active_session_by_refresh_hash.convert_error", "hash", hash, "error", err)
		return
	}

	return
}
