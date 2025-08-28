package mysqlpermissionadapter

import (
	"context"
	"database/sql"
	"fmt"

	permissionconverters "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql/permission/converters"
	permissionentities "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql/permission/entities"
	permissionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/permission_model"
)

// GetAllRoles busca todos os roles
func (pa *PermissionAdapter) GetAllRoles(ctx context.Context, tx *sql.Tx) ([]permissionmodel.RoleInterface, error) {
	query := `
		SELECT id, name, slug, description, is_system_role, is_active
		FROM roles 
		ORDER BY name
	`

	results, err := pa.Read(ctx, tx, query)
	if err != nil {
		return nil, err
	}

	roles := make([]permissionmodel.RoleInterface, 0, len(results))
	for _, row := range results {
		if len(row) != 6 {
			return nil, fmt.Errorf("unexpected number of columns: expected 6, got %d", len(row))
		}

		entity := &permissionentities.RoleEntity{
			ID:           row[0].(int64),
			Name:         string(row[1].([]byte)),
			Slug:         string(row[2].([]byte)),
			Description:  string(row[3].([]byte)),
			IsSystemRole: row[4].(int64) == 1,
			IsActive:     row[5].(int64) == 1,
		}

		role := permissionconverters.RoleEntityToDomain(entity)
		if role != nil {
			roles = append(roles, role)
		}
	}

	return roles, nil
}
