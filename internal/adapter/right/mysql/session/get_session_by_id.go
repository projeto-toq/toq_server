package sessionmysqladapter

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

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

	row := sa.QueryRowContext(ctx, tx, "get_session_by_id", query, id)
	session, err = sa.mapSessionFromScanner(ctx, row, "get_session_by_id")
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.session.get_session_by_id.scan_error", "session_id", id, "error", err)
		return nil, fmt.Errorf("get session by id: %w", err)
	}

	return
}
