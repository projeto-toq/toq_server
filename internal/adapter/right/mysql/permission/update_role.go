package mysqlpermissionadapter

import (
	"context"
	"database/sql"
	"fmt"

	permissionconverters "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/permission/converters"
	permissionmodel "github.com/projeto-toq/toq_server/internal/core/model/permission_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// UpdateRole atualiza um role existente
func (p *PermissionAdapter) UpdateRole(ctx context.Context, tx *sql.Tx, role permissionmodel.RoleInterface) (err error) {
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

	result, execErr := p.ExecContext(ctx, tx, "update", query,
		entity.Name,
		entity.Slug,
		entity.Description,
		entity.IsSystemRole,
		entity.IsActive,
		entity.ID,
	)
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.permission.update_role.exec_error", "error", execErr)
		return fmt.Errorf("update role: %w", execErr)
	}

	rowsAffected, rowsErr := result.RowsAffected()
	if rowsErr != nil {
		utils.SetSpanError(ctx, rowsErr)
		logger.Error("mysql.permission.update_role.rows_affected_error", "error", rowsErr)
		return fmt.Errorf("role update rows affected: %w", rowsErr)
	}

	logger.Debug("mysql.permission.update_role.success", "rows_affected", rowsAffected)
	return nil
}
