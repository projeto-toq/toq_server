package sessionmysqladapter

import (
	"context"
	"database/sql"
	"log/slog"

	sessionconverters "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql/session/converters"
	sessionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/session_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (sa *SessionAdapter) GetActiveSessionsByUserID(ctx context.Context, tx *sql.Tx, userID int64) (sessions []sessionmodel.SessionInterface, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	query := `SELECT id, user_id, refresh_hash, token_jti, expires_at, absolute_expires_at, created_at, rotated_at, user_agent, ip, device_id, rotation_counter, last_refresh_at, revoked 
			FROM sessions 
			WHERE user_id = ? AND revoked = false AND expires_at > UTC_TIMESTAMP()`

	entities, err := sa.Read(ctx, tx, query, userID)
	if err != nil {
		slog.Error("sessionmysqladapter/GetActiveSessionsByUserID: error executing Read", "error", err)
		return nil, status.Error(codes.Internal, "Failed to get active sessions by user ID")
	}

	sessions = make([]sessionmodel.SessionInterface, 0, len(entities))
	for _, entity := range entities {
		session, err := sessionconverters.SessionEntityToDomain(entity)
		if err != nil {
			return nil, err
		}
		sessions = append(sessions, session)
	}

	return
}
