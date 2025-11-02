package mysqluseradapter

import (
	"context"
	"fmt"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (ua *UserAdapter) ResetUserWrongSigninAttempts(ctx context.Context, userID int64) (err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `DELETE FROM temp_wrong_signin WHERE user_id = ?`

	result, execErr := ua.ExecContext(ctx, nil, "delete", query, userID)
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.user.reset_wrong_signin_attempts.exec_error", "user_id", userID, "error", execErr)
		return fmt.Errorf("reset wrong signin attempts: %w", execErr)
	}

	if _, rowsErr := result.RowsAffected(); rowsErr != nil {
		utils.SetSpanError(ctx, rowsErr)
		logger.Error("mysql.user.reset_wrong_signin_attempts.rows_affected_error", "user_id", userID, "error", rowsErr)
		return fmt.Errorf("reset wrong signin attempts rows affected: %w", rowsErr)
	}

	logger.Debug("mysql.user.reset_wrong_signin_attempts.success", "user_id", userID)
	return nil
}
