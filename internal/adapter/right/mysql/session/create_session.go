package sessionmysqladapter

import (
	"context"
	"database/sql"
	"fmt"

	sessionconverters "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/session/converters"
	sessionmodel "github.com/projeto-toq/toq_server/internal/core/model/session_model"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// CreateSession inserts a new session row into the sessions table and mutates the domain model with the generated ID.
//
// Behavior:
//   - Converts domain → entity to honor NULL semantics (sql.Null*)
//   - Executes INSERT with InstrumentedAdapter to emit tracing/metrics/logging
//   - Writes all fields, including optional metadata (user_agent, ip, device_id)
//   - Sets session ID on the domain model after successful insert
//
// Parameters:
//   - ctx: Carries tracing/logging context
//   - tx: Optional transaction (nil for standalone)
//   - session: Domain session model to persist (must have user_id, refresh_hash, expires_at)
//
// Returns:
//   - error: Infrastructure errors only (constraint violations, connectivity, tx errors)
//
// Error Semantics:
//   - Returns fmt.Errorf with context; does NOT wrap into domain/HTTP errors
//   - Leaves span marked on infra failures
//
// Usage:
//   - Called after generating refresh hash during login/refresh flows
//   - Caller owns transaction lifecycle (begin/commit/rollback)
func (sa *SessionAdapter) CreateSession(ctx context.Context, tx *sql.Tx, session sessionmodel.SessionInterface) error {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `INSERT INTO sessions 
			(user_id, refresh_hash, token_jti, expires_at, absolute_expires_at, created_at, rotated_at, user_agent, ip, device_id, rotation_counter, last_refresh_at, revoked)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	entity := sessionconverters.SessionDomainToEntity(session)

	result, execErr := sa.ExecContext(ctx, tx, "insert", query,
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
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.session.create_session.exec_error", "err", execErr)
		return fmt.Errorf("create session: %w", execErr)
	}

	id, lastIDErr := result.LastInsertId()
	if lastIDErr != nil {
		utils.SetSpanError(ctx, lastIDErr)
		logger.Error("mysql.session.create_session.last_insert_error", "err", lastIDErr)
		return fmt.Errorf("create session last insert id: %w", lastIDErr)
	}

	session.SetID(id)
	return nil
}
