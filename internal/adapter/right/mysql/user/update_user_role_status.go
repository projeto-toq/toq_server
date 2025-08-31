package mysqluseradapter

import (
	"context"
	"log/slog"
)

func (ua *UserAdapter) UpdateUserRoleStatusByUserID(ctx context.Context, userID int64, status int) (err error) {
	query := `UPDATE user_roles SET status = ?, updated_at = NOW() WHERE user_id = ? AND is_active = 1`

	_, err = ua.db.GetDB().ExecContext(ctx, query, status, userID)
	if err != nil {
		slog.Error("mysqluseradapter/UpdateUserRoleStatusByUserID: error executing update", "userID", userID, "status", status, "error", err)
		return
	}

	slog.Debug("Updated user role status", "userID", userID, "status", status)
	return
}
