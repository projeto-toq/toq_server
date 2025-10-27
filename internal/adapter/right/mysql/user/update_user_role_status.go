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

	_, err = ua.db.GetDB().ExecContext(ctx, query, status, userID)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.user.update_user_role_status.update_error", "user_id", userID, "status", status, "error", err)
		return fmt.Errorf("update user role status by user: %w", err)
	}

	logger.Debug("mysql.user.update_user_role_status.success", "user_id", userID, "status", status)
	return nil
}
