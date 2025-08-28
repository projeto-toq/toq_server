package mysqlpermissionadapter

import (
	"context"
	"database/sql"
)

// DeleteUserRole remove um user_role pelo ID
func (pa *PermissionAdapter) DeleteUserRole(ctx context.Context, tx *sql.Tx, userRoleID int64) error {
	query := `
		DELETE FROM user_roles 
		WHERE id = ?
	`

	_, err := pa.Delete(ctx, tx, query, userRoleID)
	return err
}
