package mysqlpermissionadapter

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

// DeactivateAllUserRoles desativa todos os roles de um usuário
func (pa *PermissionAdapter) DeactivateAllUserRoles(ctx context.Context, tx *sql.Tx, userID int64) error {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return err
	}
	defer spanEnd()

	query := `
		UPDATE user_roles
		SET is_active = 0
		WHERE user_id = ?
	`

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		slog.Error("mysqlpermissionadapter/DeactivateAllUserRoles: error preparing statement", "error", err)
		return err
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, userID)
	if err != nil {
		slog.Error("mysqlpermissionadapter/DeactivateAllUserRoles: error executing update", "error", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		slog.Error("mysqlpermissionadapter/DeactivateAllUserRoles: error getting rows affected", "error", err)
		return err
	}

	slog.Debug("Deactivated user roles", "userID", userID, "rowsAffected", rowsAffected)
	return nil
}
