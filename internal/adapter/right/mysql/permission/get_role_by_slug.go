package mysqlpermissionadapter

import (
	"context"
	"database/sql"
	"fmt"

	permissionconverters "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/permission/converters"
	permissionentities "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/permission/entities"
	permissionmodel "github.com/projeto-toq/toq_server/internal/core/model/permission_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// GetRoleBySlug busca um role pelo slug
func (pa *PermissionAdapter) GetRoleBySlug(ctx context.Context, tx *sql.Tx, slug string) (role permissionmodel.RoleInterface, err error) {
	ctx, spanEnd, logger, err := startPermissionOperation(ctx)
	if err != nil {
		return nil, err
	}
	defer spanEnd()

	logger = logger.With("slug", slug)

	query := `
		SELECT id, name, slug, description, is_system_role, is_active
		FROM roles 
		WHERE slug = ?
	`

	var (
		id          int64
		name        string
		slugOut     string
		description string
		isSystemInt int64
		isActiveInt int64
	)

	err = tx.QueryRowContext(ctx, query, slug).Scan(
		&id, &name, &slugOut, &description, &isSystemInt, &isActiveInt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Debug("mysql.permission.get_role_by_slug.not_found")
			return nil, nil
		}
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.permission.get_role_by_slug.scan_error", "error", err)
		return nil, fmt.Errorf("get role by slug scan: %w", err)
	}

	entity := &permissionentities.RoleEntity{
		ID:           id,
		Name:         name,
		Slug:         slugOut,
		Description:  description,
		IsSystemRole: isSystemInt == 1,
		IsActive:     isActiveInt == 1,
	}

	role = permissionconverters.RoleEntityToDomain(entity)
	logger.Debug("mysql.permission.get_role_by_slug.success")
	return role, nil
}
