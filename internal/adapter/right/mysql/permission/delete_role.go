package mysqlpermissionadapter

import (
	"context"
	"database/sql"
)

// DeleteRole remove um role pelo ID
func (pa *PermissionAdapter) DeleteRole(ctx context.Context, tx *sql.Tx, roleID int64) error {
	query := `
		DELETE FROM roles 
		WHERE id = ?
	`

	_, err := pa.Delete(ctx, tx, query, roleID)
	return err
}
