package mysqluseradapter

import (
	"context"
	"fmt"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (ua *UserAdapter) UpdateUserRoleStatusByUserID(ctx context.Context, userID int64, status int) (err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `UPDATE user_roles SET status = ? WHERE user_id = ? AND is_active = 1`

	result, execErr := ua.ExecContext(ctx, nil, "update", query, status, userID)
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.user.update_user_role_status.exec_error", "user_id", userID, "status", status, "error", execErr)
		return fmt.Errorf("update user role status by user: %w", execErr)
	}

	if _, rowsErr := result.RowsAffected(); rowsErr != nil {
		utils.SetSpanError(ctx, rowsErr)
		logger.Error("mysql.user.update_user_role_status.rows_affected_error", "user_id", userID, "status", status, "error", rowsErr)
		return fmt.Errorf("update user role status rows affected: %w", rowsErr)
	}

	logger.Debug("mysql.user.update_user_role_status.success", "user_id", userID, "status", status)
	return nil
}
