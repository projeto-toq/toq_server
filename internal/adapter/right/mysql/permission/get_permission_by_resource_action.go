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

// GetPermissionByResourceAndAction busca permissão pela combinação resource/action
func (pa *PermissionAdapter) GetPermissionByResourceAndAction(ctx context.Context, tx *sql.Tx, resource, action string) (permissionmodel.PermissionInterface, error) {
	ctx, spanEnd, logger, err := startPermissionOperation(ctx)
	if err != nil {
		return nil, err
	}
	defer spanEnd()

	logger = logger.With("resource", resource, "action", action)

	query := `
        SELECT id, name, CONCAT(resource, ':', action) AS slug, resource, action, description, conditions, is_active
        FROM permissions
        WHERE resource = ? AND action = ?
    `

	var (
		id          int64
		name        string
		slug        string
		resourceOut string
		actionOut   string
		description string
		conditions  sql.NullString
		isActiveInt int64
	)

	err = tx.QueryRowContext(ctx, query, resource, action).Scan(&id, &name, &slug, &resourceOut, &actionOut, &description, &conditions, &isActiveInt)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Debug("mysql.permission.get_permission_by_resource_action.not_found")
			return nil, nil
		}
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.permission.get_permission_by_resource_action.scan_error", "error", err)
		return nil, fmt.Errorf("get permission by resource/action scan: %w", err)
	}

	entity := &permissionentities.PermissionEntity{
		ID:          id,
		Name:        name,
		Slug:        slug,
		Resource:    resourceOut,
		Action:      actionOut,
		Description: description,
		IsActive:    isActiveInt == 1,
	}
	if conditions.Valid {
		v := conditions.String
		entity.Conditions = &v
	}

	permission, convertErr := permissionconverters.PermissionEntityToDomain(entity)
	if convertErr != nil {
		utils.SetSpanError(ctx, convertErr)
		logger.Error("mysql.permission.get_permission_by_resource_action.convert_error", "error", convertErr)
		return nil, fmt.Errorf("convert permission entity to domain: %w", convertErr)
	}

	return permission, nil
}
