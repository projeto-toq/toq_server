package sessionmysqladapter

import (
	"context"
	"database/sql"
	"fmt"

	sessionconverters "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/session/converters"
	sessionmodel "github.com/projeto-toq/toq_server/internal/core/model/session_model"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (sa *SessionAdapter) GetSessionByID(ctx context.Context, tx *sql.Tx, id int64) (session sessionmodel.SessionInterface, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `SELECT id, user_id, refresh_hash, token_jti, expires_at, absolute_expires_at, created_at, rotated_at, user_agent, ip, device_id, rotation_counter, last_refresh_at, revoked FROM sessions WHERE id = ?`

	entities, err := sa.Read(ctx, tx, query, id)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.session.get_session_by_id.read_error", "session_id", id, "error", err)
		return nil, fmt.Errorf("get session by id: %w", err)
	}

	if len(entities) == 0 {
		return nil, sql.ErrNoRows
	}

	if len(entities) > 1 {
		err := fmt.Errorf("multiple sessions found for id %d", id)
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.session.get_session_by_id.multiple_results", "session_id", id, "error", err)
		return nil, fmt.Errorf("multiple sessions found for id %d", id)
	}

	session, err = sessionconverters.SessionEntityToDomain(entities[0])
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.session.get_session_by_id.convert_error", "session_id", id, "error", err)
		return
	}

	return
}
