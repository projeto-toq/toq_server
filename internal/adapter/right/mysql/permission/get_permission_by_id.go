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
func (pa *PermissionAdapter) GetPermissionByID(ctx context.Context, tx *sql.Tx, permissionID int64) (permission permissionmodel.PermissionInterface, err error) {
	ctx, spanEnd, logger, err := startPermissionOperation(ctx)
	if err != nil {
		return nil, err
	}
	defer spanEnd()

	logger = logger.With("permission_id", permissionID)

	query := `
		SELECT id, name, CONCAT(resource, ':', action) AS slug, resource, action, description, conditions, is_active
		FROM permissions 
		WHERE id = ?
	`

	var (
		id          int64
		name        string
		slug        string
		resource    string
		action      string
		description string
		conditions  sql.NullString
		isActiveInt int64
	)

	err = tx.QueryRowContext(ctx, query, permissionID).Scan(
		&id, &name, &slug, &resource, &action, &description, &conditions, &isActiveInt,
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
		ID:          id,
		Name:        name,
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

	permission, convertErr := permissionconverters.PermissionEntityToDomain(entity)
	if convertErr != nil {
		utils.SetSpanError(ctx, convertErr)
		logger.Error("mysql.permission.get_permission_by_id.convert_error", "error", convertErr)
		return nil, fmt.Errorf("convert permission entity to domain: %w", convertErr)
	}

	logger.Debug("mysql.permission.get_permission_by_id.success")
	return permission, nil
}
