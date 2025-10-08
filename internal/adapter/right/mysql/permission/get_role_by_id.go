package mysqlpermissionadapter

import (
	"context"
	"database/sql"
	"fmt"

	permissionconverters "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql/permission/converters"
	permissionentities "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql/permission/entities"
	permissionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/permission_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

// GetRoleByID busca um role pelo ID
func (pa *PermissionAdapter) GetRoleByID(ctx context.Context, tx *sql.Tx, roleID int64) (role permissionmodel.RoleInterface, err error) {
	ctx, spanEnd, logger, err := startPermissionOperation(ctx)
	if err != nil {
		return nil, err
	}
	defer spanEnd()

	logger = logger.With("role_id", roleID)

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

	err = tx.QueryRowContext(ctx, query, roleID).Scan(
		&id, &name, &slug, &description, &isSystemInt, &isActiveInt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Debug("mysql.permission.get_role_by_id.not_found")
			return nil, nil
		}
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.permission.get_role_by_id.scan_error", "error", err)
		return nil, fmt.Errorf("get role by id scan: %w", err)
	}

	entity := &permissionentities.RoleEntity{
		ID:           id,
		Name:         name,
		Slug:         slug,
		Description:  description,
		IsSystemRole: isSystemInt == 1,
		IsActive:     isActiveInt == 1,
	}

	role = permissionconverters.RoleEntityToDomain(entity)
	logger.Debug("mysql.permission.get_role_by_id.success")
	return role, nil
}
