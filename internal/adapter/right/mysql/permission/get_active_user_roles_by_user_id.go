package mysqlpermissionadapter

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	permissionconverters "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql/permission/converters"
	permissionentities "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql/permission/entities"
	permissionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/permission_model"
)

// GetActiveUserRolesByUserID busca apenas os user_roles ativos de um usuário
func (pa *PermissionAdapter) GetActiveUserRolesByUserID(ctx context.Context, tx *sql.Tx, userID int64) ([]permissionmodel.UserRoleInterface, error) {
	query := `
		SELECT id, user_id, role_id, is_active, expires_at
		FROM user_roles 
		WHERE user_id = ? 
		  AND is_active = 1
		  AND (expires_at IS NULL OR expires_at > NOW())
		ORDER BY id
	`

	results, err := pa.Read(ctx, tx, query, userID)
	if err != nil {
		return nil, err
	}

	userRoles := make([]permissionmodel.UserRoleInterface, 0, len(results))
	for _, row := range results {
		if len(row) != 5 {
			return nil, fmt.Errorf("unexpected number of columns: expected 5, got %d", len(row))
		}

		entity := &permissionentities.UserRoleEntity{
			ID:       row[0].(int64),
			UserID:   row[1].(int64),
			RoleID:   row[2].(int64),
			IsActive: row[3].(int64) == 1,
		}

		// Handle expires_at (pode ser NULL)
		if row[4] != nil {
			if expiresAt, ok := row[4].(time.Time); ok {
				entity.ExpiresAt = &expiresAt
			}
		}

		userRole := permissionconverters.UserRoleEntityToDomain(entity)
		if userRole != nil {
			userRoles = append(userRoles, userRole)
		}
	}

	return userRoles, nil
}
