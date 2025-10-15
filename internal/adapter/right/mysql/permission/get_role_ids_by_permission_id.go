package mysqlpermissionadapter

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// GetRoleIDsByPermissionID retorna IDs de roles vinculados a uma permissão
func (pa *PermissionAdapter) GetRoleIDsByPermissionID(ctx context.Context, tx *sql.Tx, permissionID int64) ([]int64, error) {
	ctx, spanEnd, logger, err := startPermissionOperation(ctx)
	if err != nil {
		return nil, err
	}
	defer spanEnd()

	logger = logger.With("permission_id", permissionID)

	query := `SELECT DISTINCT role_id FROM role_permissions WHERE permission_id = ?`

	rows, queryErr := tx.QueryContext(ctx, query, permissionID)
	if queryErr != nil {
		utils.SetSpanError(ctx, queryErr)
		logger.Error("mysql.permission.get_role_ids_by_permission_id.query_error", "error", queryErr)
		return nil, fmt.Errorf("query role ids by permission id: %w", queryErr)
	}
	defer rows.Close()

	var roleIDs []int64
	for rows.Next() {
		var roleID int64
		if scanErr := rows.Scan(&roleID); scanErr != nil {
			utils.SetSpanError(ctx, scanErr)
			logger.Error("mysql.permission.get_role_ids_by_permission_id.scan_error", "error", scanErr)
			return nil, fmt.Errorf("scan role id by permission id: %w", scanErr)
		}
		roleIDs = append(roleIDs, roleID)
	}

	if err = rows.Err(); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.permission.get_role_ids_by_permission_id.rows_error", "error", err)
		return nil, fmt.Errorf("rows iteration role ids by permission id: %w", err)
	}

	return roleIDs, nil
}
