package mysqluseradapter

import (
	"context"
	"log/slog"
)

func (ua *UserAdapter) UpdateUserRoleStatusByUserID(ctx context.Context, userID int64, status int) (err error) {
	// Coluna updated_at não existe no schema de user_roles; remover do UPDATE.
	query := `UPDATE user_roles SET status = ? WHERE user_id = ? AND is_active = 1`

	_, err = ua.db.GetDB().ExecContext(ctx, query, status, userID)
	if err != nil {
		slog.Error("mysqluseradapter/UpdateUserRoleStatusByUserID: error executing update", "userID", userID, "status", status, "error", err)
		return
	}

	slog.Debug("Updated user role status", "userID", userID, "status", status)
	return
}
