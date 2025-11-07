package mysqlpermissionadapter

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// DeleteRolePermission remove um role_permission pelo ID
func (p *PermissionAdapter) DeleteRolePermission(ctx context.Context, tx *sql.Tx, rolePermissionID int64) (err error) {
	ctx, spanEnd, logger, err := startPermissionOperation(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	logger = logger.With("role_permission_id", rolePermissionID)

	query := `
		DELETE FROM role_permissions 
		WHERE id = ?
	`

	result, execErr := p.ExecContext(ctx, tx, "delete", query, rolePermissionID)
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.permission.delete_role_permission.exec_error", "error", execErr)
		return fmt.Errorf("delete role permission: %w", execErr)
	}

	if _, rowsErr := result.RowsAffected(); rowsErr != nil {
		utils.SetSpanError(ctx, rowsErr)
		logger.Error("mysql.permission.delete_role_permission.rows_affected_error", "error", rowsErr)
		return fmt.Errorf("role permission delete rows affected: %w", rowsErr)
	}

	logger.Debug("mysql.permission.delete_role_permission.success")
	return nil
}
