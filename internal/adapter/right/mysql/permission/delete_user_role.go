package mysqlpermissionadapter

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

// DeleteUserRole remove um user_role pelo ID
func (pa *PermissionAdapter) DeleteUserRole(ctx context.Context, tx *sql.Tx, userRoleID int64) (err error) {
	ctx, spanEnd, logger, err := startPermissionOperation(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	logger = logger.With("user_role_id", userRoleID)

	query := `
		DELETE FROM user_roles 
		WHERE id = ?
	`

	_, err = pa.Delete(ctx, tx, query, userRoleID)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.permission.delete_user_role.exec_error", "error", err)
		return fmt.Errorf("delete user role: %w", err)
	}

	logger.Debug("mysql.permission.delete_user_role.success")
	return nil
}
