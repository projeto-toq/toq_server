package mysqlpermissionadapter

import (
	"context"
	"database/sql"
	"fmt"

	permissionconverters "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql/permission/converters"
	permissionentities "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql/permission/entities"
	permissionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/permission_model"
)

// GetGrantedPermissionsByRoleID busca todas as permissões concedidas a um role
func (pa *PermissionAdapter) GetGrantedPermissionsByRoleID(ctx context.Context, tx *sql.Tx, roleID int64) ([]permissionmodel.PermissionInterface, error) {
	query := `
		SELECT p.id, p.name, CONCAT(p.resource, ':', p.action) AS slug, p.resource, p.action, p.description, p.conditions, p.is_active
		FROM permissions p
		INNER JOIN role_permissions rp ON p.id = rp.permission_id
		WHERE rp.role_id = ? 
		  AND rp.granted = 1
		  AND p.is_active = 1
		ORDER BY p.resource, p.action
	`

	results, err := pa.Read(ctx, tx, query, roleID)
	if err != nil {
		return nil, err
	}

	permissions := make([]permissionmodel.PermissionInterface, 0, len(results))
	for _, row := range results {
		if len(row) != 8 {
			return nil, fmt.Errorf("unexpected number of columns: expected 8, got %d", len(row))
		}

		entity := &permissionentities.PermissionEntity{
			ID:       row[0].(int64),
			Name:     string(row[1].([]byte)),
			Slug:     string(row[2].([]byte)),
			Resource: string(row[3].([]byte)),
			Action:   string(row[4].([]byte)),
			IsActive: row[7].(int64) == 1,
		}

		// description (pode ser NULL)
		if row[5] != nil {
			entity.Description = string(row[5].([]byte))
		}
		// conditions (pode ser NULL)
		if row[6] != nil {
			conditionsStr := string(row[6].([]byte))
			entity.Conditions = &conditionsStr
		}

		permission := permissionconverters.PermissionEntityToDomain(entity)
		if permission != nil {
			permissions = append(permissions, permission)
		}
	}

	return permissions, nil
}
