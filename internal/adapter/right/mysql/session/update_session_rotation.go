package sessionmysqladapter

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (sa *SessionAdapter) UpdateSessionRotation(ctx context.Context, tx *sql.Tx, id int64, rotationCounter int, lastRefreshAt time.Time) error {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `UPDATE sessions SET rotation_counter = ?, last_refresh_at = ? WHERE id = ?`

	_, err = sa.Update(ctx, tx, query, rotationCounter, lastRefreshAt, id)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.session.update_session_rotation.update_error", "session_id", id, "error", err)
		return fmt.Errorf("update session rotation: %w", err)
	}

	return nil
}
