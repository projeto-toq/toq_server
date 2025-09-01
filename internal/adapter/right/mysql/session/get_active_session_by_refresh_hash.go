package sessionmysqladapter

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

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

	query := `SELECT id, user_id, refresh_hash, token_jti, expires_at, absolute_expires_at, created_at, rotated_at, user_agent, ip, device_id, rotation_counter, last_refresh_at, revoked 
			FROM sessions 
			WHERE refresh_hash = ? AND revoked = false AND expires_at > UTC_TIMESTAMP()`

	entities, err := sa.Read(ctx, tx, query, hash)
	if err != nil {
		slog.Error("sessionmysqladapter/GetActiveSessionByRefreshHash: error executing Read", "error", err)
		return nil, fmt.Errorf("get active session by refresh hash: %w", err)
	}

	if len(entities) == 0 {
		slog.Debug("sessionmysqladapter/GetActiveSessionByRefreshHash: no active session found", "hash", hash)
		return nil, sql.ErrNoRows
	}

	session, err = sessionconverters.SessionEntityToDomain(entities[0])
	if err != nil {
		return
	}

	return
}
