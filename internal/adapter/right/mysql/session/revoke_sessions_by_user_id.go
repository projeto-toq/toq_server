package sessionmysqladapter

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (sa *SessionAdapter) RevokeSessionsByUserID(ctx context.Context, tx *sql.Tx, userID int64) error {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `UPDATE sessions SET revoked = true WHERE user_id = ? AND revoked = false`

	_, err = sa.Update(ctx, tx, query, userID)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.session.revoke_sessions_by_user_id.update_error", "user_id", userID, "error", err)
		return fmt.Errorf("revoke sessions by user id: %w", err)
	}

	return nil
}
