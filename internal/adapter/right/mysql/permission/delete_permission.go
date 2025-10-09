package mysqlpermissionadapter

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// DeletePermission remove uma permissão pelo ID
func (pa *PermissionAdapter) DeletePermission(ctx context.Context, tx *sql.Tx, permissionID int64) (err error) {
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

	_, err = pa.Delete(ctx, tx, query, permissionID)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.permission.delete_permission.exec_error", "error", err)
		return fmt.Errorf("delete permission: %w", err)
	}

	logger.Debug("mysql.permission.delete_permission.success")
	return nil
}
