package mysqlpermissionadapter

import (
	"context"
	"database/sql"
	"log/slog"

	permissionconverters "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql/permission/converters"
	permissionentities "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql/permission/entities"
	permissionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/permission_model"
)

// GetRoleByID busca um role pelo ID
func (pa *PermissionAdapter) GetRoleByID(ctx context.Context, tx *sql.Tx, roleID int64) (permissionmodel.RoleInterface, error) {
	query := `
		SELECT id, name, slug, description, is_system_role, is_active
		FROM roles 
		WHERE id = ?
	`

	var (
		id          int64
		name        string
		slug        string
		description string
		isSystemInt int64
		isActiveInt int64
	)

	err := tx.QueryRowContext(ctx, query, roleID).Scan(
		&id, &name, &slug, &description, &isSystemInt, &isActiveInt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		slog.Error("mysqlpermissionadapter/GetRoleByID: error scanning row", "error", err)
		return nil, err
	}

	entity := &permissionentities.RoleEntity{
		ID:           id,
		Name:         name,
		Slug:         slug,
		Description:  description,
		IsSystemRole: isSystemInt == 1,
		IsActive:     isActiveInt == 1,
	}

	return permissionconverters.RoleEntityToDomain(entity), nil
}
