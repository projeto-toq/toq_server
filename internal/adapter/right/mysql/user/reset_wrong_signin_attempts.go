package mysqluseradapter

import (
	"context"
	"log/slog"
)

func (ua *UserAdapter) ResetUserWrongSigninAttempts(ctx context.Context, userID int64) (err error) {
	query := `DELETE FROM temp_wrong_signin WHERE user_id = ?`

	_, err = ua.db.GetDB().ExecContext(ctx, query, userID)
	if err != nil {
		slog.Error("mysqluseradapter/ResetUserWrongSigninAttempts: error executing delete", "userID", userID, "error", err)
		return
	}

	slog.Debug("Reset wrong signin attempts for user", "userID", userID)
	return
}
