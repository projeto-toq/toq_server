package mysqlpermissionadapter

import (
	"context"
	"database/sql"
	"fmt"

	permissionconverters "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/permission/converters"
	permissionmodel "github.com/projeto-toq/toq_server/internal/core/model/permission_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// CreatePermission cria uma nova permissão no banco de dados
func (pa *PermissionAdapter) CreatePermission(ctx context.Context, tx *sql.Tx, permission permissionmodel.PermissionInterface) (err error) {
	ctx, spanEnd, logger, err := startPermissionOperation(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	entity, convertErr := permissionconverters.PermissionDomainToEntity(permission)
	if convertErr != nil {
		utils.SetSpanError(ctx, convertErr)
		logger.Error("mysql.permission.create_permission.convert_error", "error", convertErr)
		return fmt.Errorf("convert permission domain to entity: %w", convertErr)
	}
	if entity == nil {
		logger.Warn("mysql.permission.create_permission.empty_entity")
		return nil
	}

	logger = logger.With(
		"resource", entity.Resource,
		"action", entity.Action,
	)

	query := `
		INSERT INTO permissions (name, resource, action, description, conditions, is_active)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	id, err := pa.Create(ctx, tx, query,
		entity.Name,
		entity.Resource,
		entity.Action,
		entity.Description,
		entity.Conditions,
		entity.IsActive,
	)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.permission.create_permission.exec_error", "error", err)
		return fmt.Errorf("create permission: %w", err)
	}

	permission.SetID(id)
	logger.Debug("mysql.permission.create_permission.success", "permission_id", id)
	return nil
}
