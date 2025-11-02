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

	result, execErr := pa.ExecContext(ctx, tx, "update", query, userID, roleID)
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.permission.activate_user_role.exec_error", "error", execErr)
		return fmt.Errorf("execute activate user role: %w", execErr)
	}

	rowsAffected, rowsErr := result.RowsAffected()
	if rowsErr != nil {
		utils.SetSpanError(ctx, rowsErr)
		logger.Error("mysql.permission.activate_user_role.rows_affected_error", "error", rowsErr)
		return fmt.Errorf("rows affected activate user role: %w", rowsErr)
	}

	if rowsAffected == 0 {
		logger.Warn("mysql.permission.activate_user_role.no_rows")
		return nil
	}

	logger.Debug("mysql.permission.activate_user_role.success", "rows_affected", rowsAffected)
	return nil
}
