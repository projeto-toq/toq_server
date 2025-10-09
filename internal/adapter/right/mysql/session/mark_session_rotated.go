package sessionmysqladapter

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (sa *SessionAdapter) MarkSessionRotated(ctx context.Context, tx *sql.Tx, id int64) error {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `UPDATE sessions SET rotated_at = UTC_TIMESTAMP() WHERE id = ?`

	_, err = sa.Update(ctx, tx, query, id)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.session.mark_session_rotated.update_error", "session_id", id, "error", err)
		return fmt.Errorf("mark session rotated: %w", err)
	}

	return nil
}
