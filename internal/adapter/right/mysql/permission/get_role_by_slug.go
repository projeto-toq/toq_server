package mysqlpermissionadapter

import (
	"context"
	"database/sql"
	"fmt"

	permissionconverters "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql/permission/converters"
	permissionentities "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql/permission/entities"
	permissionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/permission_model"
)

// GetRoleBySlug busca um role pelo slug
func (pa *PermissionAdapter) GetRoleBySlug(ctx context.Context, tx *sql.Tx, slug string) (permissionmodel.RoleInterface, error) {
	query := `
		SELECT id, name, slug, description, is_system_role, is_active
		FROM roles 
		WHERE slug = ?
	`

	results, err := pa.Read(ctx, tx, query, slug)
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

	entity := &permissionentities.RoleEntity{
		ID:           row[0].(int64),
		Name:         string(row[1].([]byte)),
		Slug:         string(row[2].([]byte)),
		Description:  string(row[3].([]byte)),
		IsSystemRole: row[4].(int64) == 1,
		IsActive:     row[5].(int64) == 1,
	}

	return permissionconverters.RoleEntityToDomain(entity), nil
}
