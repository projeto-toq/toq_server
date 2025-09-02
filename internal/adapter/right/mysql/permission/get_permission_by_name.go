package mysqlpermissionadapter

import (
	"context"
	"database/sql"
	"log/slog"

	permissionconverters "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql/permission/converters"
	permissionentities "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql/permission/entities"
	permissionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/permission_model"
)

// GetPermissionByName busca uma permissão pelo nome
func (pa *PermissionAdapter) GetPermissionByName(ctx context.Context, tx *sql.Tx, name string) (permissionmodel.PermissionInterface, error) {
	query := `
		SELECT id, name, CONCAT(resource, ':', action) AS slug, resource, action, description, conditions, is_active
		FROM permissions 
		WHERE name = ?
	`

	var (
		id          int64
		nameOut     string
		slug        string
		resource    string
		action      string
		description string
		conditions  sql.NullString
		isActiveInt int64
	)

	err := tx.QueryRowContext(ctx, query, name).Scan(
		&id, &nameOut, &slug, &resource, &action, &description, &conditions, &isActiveInt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		slog.Error("mysqlpermissionadapter/GetPermissionByName: error scanning row", "error", err)
		return nil, err
	}

	entity := &permissionentities.PermissionEntity{
		ID:          id,
		Name:        nameOut,
		Slug:        slug,
		Resource:    resource,
		Action:      action,
		Description: description,
		IsActive:    isActiveInt == 1,
	}
	if conditions.Valid {
		v := conditions.String
		entity.Conditions = &v
	}

	return permissionconverters.PermissionEntityToDomain(entity), nil
}
