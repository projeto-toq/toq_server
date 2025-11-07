package mysqlpermissionadapter

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// DeletePermission remove uma permissão pelo ID
func (p *PermissionAdapter) DeletePermission(ctx context.Context, tx *sql.Tx, permissionID int64) (err error) {
	ctx, spanEnd, logger, err := startPermissionOperation(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	logger = logger.With("permission_id", permissionID)

	query := `
		DELETE FROM permissions 
		WHERE id = ?
	`

	result, execErr := p.ExecContext(ctx, tx, "delete", query, permissionID)
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.permission.delete_permission.exec_error", "error", execErr)
		return fmt.Errorf("delete permission: %w", execErr)
	}

	if _, rowsErr := result.RowsAffected(); rowsErr != nil {
		utils.SetSpanError(ctx, rowsErr)
		logger.Error("mysql.permission.delete_permission.rows_affected_error", "error", rowsErr)
		return fmt.Errorf("permission delete rows affected: %w", rowsErr)
	}

	logger.Debug("mysql.permission.delete_permission.success")
	return nil
}
