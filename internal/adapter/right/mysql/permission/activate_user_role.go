package mysqlpermissionadapter

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// ActivateUserRole ativa um role específico do usuário
func (pa *PermissionAdapter) ActivateUserRole(ctx context.Context, tx *sql.Tx, userID, roleID int64) error {
	ctx, spanEnd, logger, err := startPermissionOperation(ctx)
	if err != nil {
		return err
	}
	defer spanEnd()

	logger = logger.With("user_id", userID, "role_id", roleID)

	query := `
		UPDATE user_roles
		SET is_active = 1
		WHERE user_id = ? AND role_id = ?
	`

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.permission.activate_user_role.prepare_error", "error", err)
		return fmt.Errorf("prepare activate user role statement: %w", err)
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, userID, roleID)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.permission.activate_user_role.exec_error", "error", err)
		return fmt.Errorf("execute activate user role: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.permission.activate_user_role.rows_affected_error", "error", err)
		return fmt.Errorf("rows affected activate user role: %w", err)
	}

	if rowsAffected == 0 {
		logger.Warn("mysql.permission.activate_user_role.no_rows")
		return nil
	}

	logger.Debug("mysql.permission.activate_user_role.success", "rows_affected", rowsAffected)
	return nil
}
