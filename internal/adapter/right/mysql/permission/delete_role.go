package mysqlpermissionadapter

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// DeleteRole remove um role pelo ID
func (pa *PermissionAdapter) DeleteRole(ctx context.Context, tx *sql.Tx, roleID int64) (err error) {
	ctx, spanEnd, logger, err := startPermissionOperation(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	logger = logger.With("role_id", roleID)

	query := `
		DELETE FROM roles 
		WHERE id = ?
	`

	result, execErr := pa.ExecContext(ctx, tx, "delete", query, roleID)
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.permission.delete_role.exec_error", "error", execErr)
		return fmt.Errorf("delete role: %w", execErr)
	}

	if _, rowsErr := result.RowsAffected(); rowsErr != nil {
		utils.SetSpanError(ctx, rowsErr)
		logger.Error("mysql.permission.delete_role.rows_affected_error", "error", rowsErr)
		return fmt.Errorf("role delete rows affected: %w", rowsErr)
	}

	logger.Debug("mysql.permission.delete_role.success")
	return nil
}
