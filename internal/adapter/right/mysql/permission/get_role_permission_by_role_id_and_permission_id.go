package mysqlpermissionadapter

import (
	"context"
	"database/sql"
	"fmt"

	permissionconverters "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql/permission/converters"
	permissionentities "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql/permission/entities"
	permissionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/permission_model"
)

// GetRolePermissionByRoleIDAndPermissionID busca um role_permission específico pela combinação role_id + permission_id
func (pa *PermissionAdapter) GetRolePermissionByRoleIDAndPermissionID(ctx context.Context, tx *sql.Tx, roleID, permissionID int64) (permissionmodel.RolePermissionInterface, error) {
	query := `
		SELECT id, role_id, permission_id, granted, conditions
		FROM role_permissions 
		WHERE role_id = ? AND permission_id = ?
		LIMIT 1
	`

	results, err := pa.Read(ctx, tx, query, roleID, permissionID)
	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, nil // Não encontrado
	}

	row := results[0]
	if len(row) != 5 {
		return nil, fmt.Errorf("unexpected number of columns: expected 5, got %d", len(row))
	}

	entity := &permissionentities.RolePermissionEntity{
		ID:           row[0].(int64),
		RoleID:       row[1].(int64),
		PermissionID: row[2].(int64),
		Granted:      row[3].(int64) == 1,
	}

	// Handle conditions (pode ser NULL)
	if row[4] != nil {
		conditionsStr := string(row[4].([]byte))
		entity.Conditions = &conditionsStr
	}

	return permissionconverters.RolePermissionEntityToDomain(entity), nil
}
