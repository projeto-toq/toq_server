package mysqlpermissionadapter

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// GetActiveUserIDsByRoleID retorna IDs de usuários ativos associados a um role específico
func (p *PermissionAdapter) GetActiveUserIDsByRoleID(ctx context.Context, tx *sql.Tx, roleID int64) ([]int64, error) {
	// Initialize tracing for observability
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, err
	}
	defer spanEnd()

	// Attach logger to context for request_id/trace_id propagation
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `
        SELECT DISTINCT
            ur.user_id
        FROM user_roles ur
        INNER JOIN roles r ON r.id = ur.role_id AND r.is_active = 1
        WHERE ur.role_id = ?
          AND ur.is_active = 1
          AND (ur.expires_at IS NULL OR ur.expires_at > NOW())
    `

	rows, readErr := p.QueryContext(ctx, tx, "select", query, roleID)
	if readErr != nil {
		utils.SetSpanError(ctx, readErr)
		logger.Error("mysql.user.get_active_user_ids_by_role_id.read_error", "error", readErr)
		return nil, fmt.Errorf("get active user ids by role id read: %w", readErr)
	}
	defer rows.Close()

	userIDs := make([]int64, 0)
	for rows.Next() {
		var userID int64
		if scanErr := rows.Scan(&userID); scanErr != nil {
			utils.SetSpanError(ctx, scanErr)
			logger.Error("mysql.user.get_active_user_ids_by_role_id.scan_error", "error", scanErr)
			return nil, fmt.Errorf("scan active user id by role id: %w", scanErr)
		}
		userIDs = append(userIDs, userID)
	}

	if rowsErr := rows.Err(); rowsErr != nil {
		utils.SetSpanError(ctx, rowsErr)
		logger.Error("mysql.user.get_active_user_ids_by_role_id.rows_error", "error", rowsErr)
		return nil, fmt.Errorf("iterate active user ids by role id: %w", rowsErr)
	}

	logger.Debug("mysql.user.get_active_user_ids_by_role_id.success", "count", len(userIDs))
	return userIDs, nil
}
