package mysqlpermissionadapter

import (
	"context"
	"database/sql"
	"fmt"

	permissionconverters "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql/permission/converters"
	permissionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/permission_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

// UpdateRole atualiza um role existente
func (pa *PermissionAdapter) UpdateRole(ctx context.Context, tx *sql.Tx, role permissionmodel.RoleInterface) (err error) {
	ctx, spanEnd, logger, err := startPermissionOperation(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	entity := permissionconverters.RoleDomainToEntity(role)
	if entity == nil {
		logger.Warn("mysql.permission.update_role.empty_entity")
		return nil
	}

	logger = logger.With(
		"role_id", entity.ID,
		"role_slug", entity.Slug,
	)

	query := `
		UPDATE roles 
		SET name = ?, slug = ?, description = ?, is_system_role = ?, is_active = ?
		WHERE id = ?
	`

	rowsAffected, err := pa.Update(ctx, tx, query,
		entity.Name,
		entity.Slug,
		entity.Description,
		entity.IsSystemRole,
		entity.IsActive,
		entity.ID,
	)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.permission.update_role.exec_error", "error", err)
		return fmt.Errorf("update role: %w", err)
	}

	logger.Debug("mysql.permission.update_role.success", "rows_affected", rowsAffected)
	return nil
}
