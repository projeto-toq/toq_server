package mysqlpermissionadapter

import (
	"context"
	"database/sql"
	"fmt"

	permissionconverters "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql/permission/converters"
	permissionentities "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql/permission/entities"
	permissionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/permission_model"
)

// GetPermissionByID busca uma permissão pelo ID
func (pa *PermissionAdapter) GetPermissionByID(ctx context.Context, tx *sql.Tx, permissionID int64) (permissionmodel.PermissionInterface, error) {
	query := `
		SELECT id, name, slug, resource, action, description, conditions, is_active
		FROM permissions 
		WHERE id = ?
	`

	results, err := pa.Read(ctx, tx, query, permissionID)
	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, nil // Não encontrado
	}

	row := results[0]
	if len(row) != 8 {
		return nil, fmt.Errorf("unexpected number of columns: expected 8, got %d", len(row))
	}

	entity := &permissionentities.PermissionEntity{
		ID:          row[0].(int64),
		Name:        string(row[1].([]byte)),
		Slug:        string(row[2].([]byte)),
		Resource:    string(row[3].([]byte)),
		Action:      string(row[4].([]byte)),
		Description: string(row[5].([]byte)),
		IsActive:    row[7].(int64) == 1,
	}

	// Handle conditions (pode ser NULL)
	if row[6] != nil {
		conditionsStr := string(row[6].([]byte))
		entity.Conditions = &conditionsStr
	}

	return permissionconverters.PermissionEntityToDomain(entity), nil
}
