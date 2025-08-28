package mysqlpermissionadapter

import (
	"context"
	"database/sql"
)

// DeleteRolePermission remove um role_permission pelo ID
func (pa *PermissionAdapter) DeleteRolePermission(ctx context.Context, tx *sql.Tx, rolePermissionID int64) error {
	query := `
		DELETE FROM role_permissions 
		WHERE id = ?
	`

	_, err := pa.Delete(ctx, tx, query, rolePermissionID)
	return err
}
