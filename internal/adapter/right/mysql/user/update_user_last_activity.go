package mysqluseradapter

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (ua *UserAdapter) UpdateUserLastActivity(ctx context.Context, tx *sql.Tx, id int64) (err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `UPDATE users SET last_activity_at = ? WHERE id = ?;`

	result, execErr := ua.ExecContext(ctx, tx, "update", query,
		time.Now().UTC(),
		id,
	)
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.user.update_user_last_activity.exec_error", "error", execErr)
		return fmt.Errorf("update user last_activity: %w", execErr)
	}

	if _, rowsErr := result.RowsAffected(); rowsErr != nil {
		utils.SetSpanError(ctx, rowsErr)
		logger.Error("mysql.user.update_user_last_activity.rows_affected_error", "error", rowsErr)
		return fmt.Errorf("user last activity update rows affected: %w", rowsErr)
	}

	return
}
