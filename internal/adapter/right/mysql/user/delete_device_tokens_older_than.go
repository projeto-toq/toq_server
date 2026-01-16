package mysqluseradapter

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// DeleteDeviceTokensOlderThan removes device tokens whose updated_at is older than the provided cutoff.
//
// Parameters:
//   - ctx: Context for tracing/logging.
//   - tx: Optional transaction (nil for standalone maintenance).
//   - cutoff: Tokens with updated_at earlier than this timestamp are removed.
//   - limit: Maximum rows to delete in this batch (defensive bound).
//
// Returns the number of rows deleted and infrastructure errors, if any. Zero rows is success.
func (ua *UserAdapter) DeleteDeviceTokensOlderThan(ctx context.Context, tx *sql.Tx, cutoff time.Time, limit int) (int64, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return 0, err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	if limit <= 0 {
		limit = 500
	}

	query := `DELETE FROM device_tokens WHERE updated_at < ? LIMIT ?`

	res, execErr := ua.ExecContext(ctx, tx, "delete", query, cutoff, limit)
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.user.delete_device_tokens_older.exec_error", "cutoff", cutoff, "limit", limit, "error", execErr)
		return 0, fmt.Errorf("delete stale device tokens: %w", execErr)
	}

	rows, raErr := res.RowsAffected()
	if raErr != nil {
		logger.Warn("mysql.user.delete_device_tokens_older.rows_affected_warning", "error", raErr)
		return 0, nil
	}

	if rows > 0 {
		logger.Debug("mysql.user.delete_device_tokens_older.success", "deleted", rows, "cutoff", cutoff, "limit", limit)
	}
	return rows, nil
}
