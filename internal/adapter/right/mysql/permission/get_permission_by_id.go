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

// GetPermissionByID busca uma permissão pelo ID
func (p *PermissionAdapter) GetPermissionByID(ctx context.Context, tx *sql.Tx, permissionID int64) (permission permissionmodel.PermissionInterface, err error) {
	ctx, spanEnd, logger, err := startPermissionOperation(ctx)
	if err != nil {
		return nil, err
	}
	defer spanEnd()

	logger = logger.With("permission_id", permissionID)

	query := `
		SELECT id, name, action, description, is_active
		FROM permissions 
		WHERE id = ?
	`

	var (
		id          int64
		name        string
		action      string
		description sql.NullString
		isActiveInt int64
	)

	row := p.QueryRowContext(ctx, tx, "select", query, permissionID)
	err = row.Scan(
		&id, &name, &action, &description, &isActiveInt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Debug("mysql.permission.get_permission_by_id.not_found")
			return nil, nil
		}
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.permission.get_permission_by_id.scan_error", "error", err)
		return nil, fmt.Errorf("get permission by id scan: %w", err)
	}

	entity := &permissionentities.PermissionEntity{
		ID:       id,
		Name:     name,
		Action:   action,
		IsActive: isActiveInt == 1,
	}
	if description.Valid {
		desc := description.String
		entity.Description = &desc
	}

	permission, convertErr := permissionconverters.PermissionEntityToDomain(entity)
	if convertErr != nil {
		utils.SetSpanError(ctx, convertErr)
		logger.Error("mysql.permission.get_permission_by_id.convert_error", "error", convertErr)
		return nil, fmt.Errorf("convert permission entity to domain: %w", convertErr)
	}

	logger.Debug("mysql.permission.get_permission_by_id.success")
	return permission, nil
}
