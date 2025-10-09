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

	_, err = ua.db.GetDB().ExecContext(ctx, query, userID)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.user.reset_wrong_signin_attempts.delete_error", "user_id", userID, "error", err)
		return fmt.Errorf("reset wrong signin attempts: %w", err)
	}

	logger.Debug("mysql.user.reset_wrong_signin_attempts.success", "user_id", userID)
	return nil
}
