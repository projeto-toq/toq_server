package mysqluseradapter

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// ActivateUserRole ativa um role específico do usuário
func (ua *UserAdapter) ActivateUserRole(ctx context.Context, tx *sql.Tx, userID, roleID int64) error {
	// Initialize tracing for observability
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return err
	}
	defer spanEnd()

	// Attach logger to context for request_id/trace_id propagation
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `
		UPDATE user_roles
		SET is_active = 1
		WHERE user_id = ? AND role_id = ?
	`

	result, execErr := ua.ExecContext(ctx, tx, "update", query, userID, roleID)
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.user.activate_user_role.exec_error", "error", execErr)
		return fmt.Errorf("execute activate user role: %w", execErr)
	}

	rowsAffected, rowsErr := result.RowsAffected()
	if rowsErr != nil {
		utils.SetSpanError(ctx, rowsErr)
		logger.Error("mysql.user.activate_user_role.rows_affected_error", "error", rowsErr)
		return fmt.Errorf("rows affected activate user role: %w", rowsErr)
	}

	if rowsAffected == 0 {
		logger.Warn("mysql.user.activate_user_role.no_rows")
		return nil
	}

	logger.Debug("mysql.user.activate_user_role.success", "rows_affected", rowsAffected)
	return nil
}
