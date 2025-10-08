package mysqlpermissionadapter

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
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

	_, err = pa.Delete(ctx, tx, query, roleID)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.permission.delete_role.exec_error", "error", err)
		return fmt.Errorf("delete role: %w", err)
	}

	logger.Debug("mysql.permission.delete_role.success")
	return nil
}
