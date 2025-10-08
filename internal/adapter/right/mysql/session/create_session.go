package sessionmysqladapter

import (
	"context"
	"database/sql"
	"fmt"

	sessionconverters "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql/session/converters"
	sessionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/session_model"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (sa *SessionAdapter) CreateSession(ctx context.Context, tx *sql.Tx, session sessionmodel.SessionInterface) (err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `INSERT INTO sessions 
			(user_id, refresh_hash, token_jti, expires_at, absolute_expires_at, created_at, rotated_at, user_agent, ip, device_id, rotation_counter, last_refresh_at, revoked)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	entity := sessionconverters.SessionDomainToEntity(ctx, session)

	id, err := sa.Create(ctx, tx, query,
		entity.UserID,
		entity.RefreshHash,
		entity.TokenJTI,
		entity.ExpiresAt,
		entity.AbsoluteExpiresAt,
		entity.CreatedAt,
		entity.RotatedAt,
		entity.UserAgent,
		entity.IP,
		entity.DeviceID,
		entity.RotationCounter,
		entity.LastRefreshAt,
		entity.Revoked)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.session.create_session.create_error", "error", err)
		return fmt.Errorf("create session: %w", err)
	}

	session.SetID(id)
	return
}
