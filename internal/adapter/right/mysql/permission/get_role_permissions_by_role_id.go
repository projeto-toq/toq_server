package mysqlpermissionadapter

import (
	"context"
	"database/sql"
	"fmt"

	permissionconverters "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql/permission/converters"
	permissionentities "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql/permission/entities"
	permissionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/permission_model"
)

// GetRolePermissionsByRoleID busca todas as associações role_permission de um role
func (pa *PermissionAdapter) GetRolePermissionsByRoleID(ctx context.Context, tx *sql.Tx, roleID int64) ([]permissionmodel.RolePermissionInterface, error) {
	query := `
		SELECT id, role_id, permission_id, granted, conditions
		FROM role_permissions 
		WHERE role_id = ?
		ORDER BY id
	`

	results, err := pa.Read(ctx, tx, query, roleID)
	if err != nil {
		return nil, err
	}

	rolePermissions := make([]permissionmodel.RolePermissionInterface, 0, len(results))
	for _, row := range results {
		if len(row) != 5 {
			return nil, fmt.Errorf("unexpected number of columns: expected 5, got %d", len(row))
		}

		entity := &permissionentities.RolePermissionEntity{
			ID:           row[0].(int64),
			RoleID:       row[1].(int64),
			PermissionID: row[2].(int64),
			Granted:      row[3].(int64) == 1,
		}

		// conditions (pode ser NULL)
		if row[4] != nil {
			conditionsStr := string(row[4].([]byte))
			entity.Conditions = &conditionsStr
		}

		rolePermission := permissionconverters.RolePermissionEntityToDomain(entity)
		if rolePermission != nil {
			rolePermissions = append(rolePermissions, rolePermission)
		}
	}

	return rolePermissions, nil
}
