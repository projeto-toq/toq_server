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

// GetPermissionByAction retrieves a permission using the HTTP action identifier (METHOD:PATH)
func (pa *PermissionAdapter) GetPermissionByAction(ctx context.Context, tx *sql.Tx, action string) (permissionmodel.PermissionInterface, error) {
	ctx, spanEnd, logger, err := startPermissionOperation(ctx)
	if err != nil {
		return nil, err
	}
	defer spanEnd()

	logger = logger.With("action", action)

	query := `
        SELECT id, name, action, description, is_active
        FROM permissions
        WHERE action = ?
    `

	var (
		id          int64
		name        string
		actionOut   string
		description sql.NullString
		isActiveInt int64
	)

	row := pa.QueryRowContext(ctx, tx, "select", query, action)
	err = row.Scan(&id, &name, &actionOut, &description, &isActiveInt)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Debug("mysql.permission.get_permission_by_action.not_found")
			return nil, nil
		}
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.permission.get_permission_by_action.scan_error", "error", err)
		return nil, fmt.Errorf("get permission by action scan: %w", err)
	}

	entity := &permissionentities.PermissionEntity{
		ID:       id,
		Name:     name,
		Action:   actionOut,
		IsActive: isActiveInt == 1,
	}

	if description.Valid {
		desc := description.String
		entity.Description = &desc
	}

	permission, convertErr := permissionconverters.PermissionEntityToDomain(entity)
	if convertErr != nil {
		utils.SetSpanError(ctx, convertErr)
		logger.Error("mysql.permission.get_permission_by_action.convert_error", "error", convertErr)
		return nil, fmt.Errorf("convert permission entity to domain: %w", convertErr)
	}

	return permission, nil
}
