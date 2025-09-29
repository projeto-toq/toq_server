package mysqlpermissionadapter

import (
	"context"
	"database/sql"
	"fmt"
)

// GetActiveUserIDsByRoleID retorna IDs de usuários ativos associados a um role específico
func (pa *PermissionAdapter) GetActiveUserIDsByRoleID(ctx context.Context, tx *sql.Tx, roleID int64) ([]int64, error) {
	query := `
        SELECT DISTINCT
            ur.user_id
        FROM user_roles ur
        INNER JOIN roles r ON r.id = ur.role_id AND r.is_active = 1
        WHERE ur.role_id = ?
          AND ur.is_active = 1
          AND (ur.expires_at IS NULL OR ur.expires_at > NOW())
    `

	results, err := pa.Read(ctx, tx, query, roleID)
	if err != nil {
		return nil, err
	}

	userIDs := make([]int64, 0, len(results))
	for _, row := range results {
		if len(row) != 1 {
			return nil, fmt.Errorf("unexpected number of columns: expected 1, got %d", len(row))
		}

		value, ok := row[0].(int64)
		if !ok {
			return nil, fmt.Errorf("unexpected column type for user_id: %T", row[0])
		}
		userIDs = append(userIDs, value)
	}

	return userIDs, nil
}
