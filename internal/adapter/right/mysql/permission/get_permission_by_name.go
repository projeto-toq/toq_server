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

// GetPermissionByName busca uma permissão pelo nome
func (pa *PermissionAdapter) GetPermissionByName(ctx context.Context, tx *sql.Tx, name string) (permission permissionmodel.PermissionInterface, err error) {
	ctx, spanEnd, logger, err := startPermissionOperation(ctx)
	if err != nil {
		return nil, err
	}
	defer spanEnd()

	logger = logger.With("permission_name", name)

	query := `
		SELECT id, name, CONCAT(resource, ':', action) AS slug, resource, action, description, conditions, is_active
		FROM permissions 
		WHERE name = ?
	`

	var (
		id          int64
		nameOut     string
		slug        string
		resource    string
		action      string
		description string
		conditions  sql.NullString
		isActiveInt int64
	)

	err = tx.QueryRowContext(ctx, query, name).Scan(
		&id, &nameOut, &slug, &resource, &action, &description, &conditions, &isActiveInt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Debug("mysql.permission.get_permission_by_name.not_found")
			return nil, nil
		}
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.permission.get_permission_by_name.scan_error", "error", err)
		return nil, fmt.Errorf("get permission by name scan: %w", err)
	}

	entity := &permissionentities.PermissionEntity{
		ID:          id,
		Name:        nameOut,
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
		logger.Error("mysql.permission.get_permission_by_name.convert_error", "error", convertErr)
		return nil, fmt.Errorf("convert permission entity to domain: %w", convertErr)
	}

	logger.Debug("mysql.permission.get_permission_by_name.success")
	return permission, nil
}
