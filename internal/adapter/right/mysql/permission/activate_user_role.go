package mysqlpermissionadapter

import (
	"context"
	"database/sql"
	"log/slog"

"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

// ActivateUserRole ativa um role específico do usuário
func (pa *PermissionAdapter) ActivateUserRole(ctx context.Context, tx *sql.Tx, userID, roleID int64) error {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return err
	}
	defer spanEnd()

	query := `
		UPDATE user_roles
		SET is_active = 1
		WHERE user_id = ? AND role_id = ?
	`

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		slog.Error("mysqlpermissionadapter/ActivateUserRole: error preparing statement", "error", err)
		return err
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, userID, roleID)
	if err != nil {
		slog.Error("mysqlpermissionadapter/ActivateUserRole: error executing update", "error", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		slog.Error("mysqlpermissionadapter/ActivateUserRole: error getting rows affected", "error", err)
		return err
	}

	if rowsAffected == 0 {
		slog.Warn("No user role found to activate", "userID", userID, "roleID", roleID)
	}

	slog.Info("Activated user role", "userID", userID, "roleID", roleID, "rowsAffected", rowsAffected)
	return nil
}
