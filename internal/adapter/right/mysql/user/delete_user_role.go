package mysqluseradapter

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// DeleteUserRole remove um user_role pelo ID
func (ua *UserAdapter) DeleteUserRole(ctx context.Context, tx *sql.Tx, userRoleID int64) (err error) {
	// Initialize tracing for observability
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	// Attach logger to context for request_id/trace_id propagation
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `
		DELETE FROM user_roles 
		WHERE id = ?
	`

	result, execErr := ua.ExecContext(ctx, tx, "delete", query, userRoleID)
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.permission.delete_user_role.exec_error", "error", execErr)
		return fmt.Errorf("delete user role: %w", execErr)
	}

	if _, rowsErr := result.RowsAffected(); rowsErr != nil {
		utils.SetSpanError(ctx, rowsErr)
		logger.Error("mysql.permission.delete_user_role.rows_affected_error", "error", rowsErr)
		return fmt.Errorf("user role delete rows affected: %w", rowsErr)
	}

	logger.Debug("mysql.permission.delete_user_role.success")
	return nil
}
