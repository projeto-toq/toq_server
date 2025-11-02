package mysqlpermissionadapter

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// DeactivateAllUserRoles desativa todos os roles de um usuário
func (pa *PermissionAdapter) DeactivateAllUserRoles(ctx context.Context, tx *sql.Tx, userID int64) error {
	ctx, spanEnd, logger, err := startPermissionOperation(ctx)
	if err != nil {
		return err
	}
	defer spanEnd()

	logger = logger.With("user_id", userID)

	query := `
		UPDATE user_roles
		SET is_active = 0
		WHERE user_id = ?
	`

	result, execErr := pa.ExecContext(ctx, tx, "update", query, userID)
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.permission.deactivate_all_user_roles.exec_error", "error", execErr)
		return fmt.Errorf("execute deactivate all user roles: %w", execErr)
	}

	rowsAffected, rowsErr := result.RowsAffected()
	if rowsErr != nil {
		utils.SetSpanError(ctx, rowsErr)
		logger.Error("mysql.permission.deactivate_all_user_roles.rows_affected_error", "error", rowsErr)
		return fmt.Errorf("rows affected deactivate all user roles: %w", rowsErr)
	}

	logger.Debug("mysql.permission.deactivate_all_user_roles.success", "rows_affected", rowsAffected)
	return nil
}
