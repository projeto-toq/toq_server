package mysqlpermissionadapter

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
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

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.permission.deactivate_all_user_roles.prepare_error", "error", err)
		return fmt.Errorf("prepare deactivate all user roles statement: %w", err)
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, userID)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.permission.deactivate_all_user_roles.exec_error", "error", err)
		return fmt.Errorf("execute deactivate all user roles: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.permission.deactivate_all_user_roles.rows_affected_error", "error", err)
		return fmt.Errorf("rows affected deactivate all user roles: %w", err)
	}

	logger.Debug("mysql.permission.deactivate_all_user_roles.success", "rows_affected", rowsAffected)
	return nil
}
