package mysqlpermissionadapter

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

// DeleteRolePermission remove um role_permission pelo ID
func (pa *PermissionAdapter) DeleteRolePermission(ctx context.Context, tx *sql.Tx, rolePermissionID int64) (err error) {
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

	_, err = pa.Delete(ctx, tx, query, rolePermissionID)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.permission.delete_role_permission.exec_error", "error", err)
		return fmt.Errorf("delete role permission: %w", err)
	}

	logger.Debug("mysql.permission.delete_role_permission.success")
	return nil
}
