package mysqlpermissionadapter

import (
	"context"
	"database/sql"
)

// DeletePermission remove uma permissão pelo ID
func (pa *PermissionAdapter) DeletePermission(ctx context.Context, tx *sql.Tx, permissionID int64) error {
	query := `
		DELETE FROM permissions 
		WHERE id = ?
	`

	_, err := pa.Delete(ctx, tx, query, permissionID)
	return err
}
