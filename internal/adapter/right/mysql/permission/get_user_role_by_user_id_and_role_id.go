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

// GetUserRoleByUserIDAndRoleID busca um user_role específico pela combinação user_id + role_id
func (pa *PermissionAdapter) GetUserRoleByUserIDAndRoleID(ctx context.Context, tx *sql.Tx, userID, roleID int64) (permissionmodel.UserRoleInterface, error) {
	query := `
		SELECT id, user_id, role_id, is_active, status, expires_at
		FROM user_roles 
		WHERE user_id = ? AND role_id = ?
		LIMIT 1
	`

	results, err := pa.Read(ctx, tx, query, userID, roleID)
	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, nil // Não encontrado
	}

	row := results[0]
	if len(row) != 6 {
		return nil, fmt.Errorf("unexpected number of columns: expected 6, got %d", len(row))
	}

	entity := &permissionentities.UserRoleEntity{
		ID:       row[0].(int64),
		UserID:   row[1].(int64),
		RoleID:   row[2].(int64),
		IsActive: row[3].(int64) == 1,
		Status:   row[4].(int64),
	}

	// Handle expires_at (pode ser NULL)
	if row[5] != nil {
		if expiresAt, ok := row[5].(time.Time); ok {
			entity.ExpiresAt = &expiresAt
		}
	}

	return permissionconverters.UserRoleEntityToDomain(entity), nil
}
