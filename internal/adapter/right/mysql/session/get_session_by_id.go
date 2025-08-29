package sessionmysqladapter

import (
	"context"
	"database/sql"
	"log/slog"

	sessionconverters "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql/session/converters"
	sessionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/session_model"
	
	
"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (sa *SessionAdapter) GetSessionByID(ctx context.Context, tx *sql.Tx, id int64) (session sessionmodel.SessionInterface, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	query := `SELECT id, user_id, refresh_hash, token_jti, expires_at, absolute_expires_at, created_at, rotated_at, user_agent, ip, device_id, rotation_counter, last_refresh_at, revoked FROM sessions WHERE id = ?`

	entities, err := sa.Read(ctx, tx, query, id)
	if err != nil {
		slog.Error("sessionmysqladapter/GetSessionByID: error executing Read", "error", err)
		return nil, utils.ErrInternalServer
	}

	if len(entities) == 0 {
		return nil, utils.ErrInternalServer
	}

	if len(entities) > 1 {
		slog.Error("sessionmysqladapter/GetSessionByID: multiple sessions found with the same ID", "ID", id)
		return nil, utils.ErrInternalServer
	}

	session, err = sessionconverters.SessionEntityToDomain(entities[0])
	if err != nil {
		return
	}

	return
}
