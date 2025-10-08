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

// GetUserPermissionsByResource busca permissões de um usuário filtradas por resource
func (pa *PermissionAdapter) GetUserPermissionsByResource(ctx context.Context, tx *sql.Tx, userID int64, resource string) (permissions []permissionmodel.PermissionInterface, err error) {
	ctx, spanEnd, logger, err := startPermissionOperation(ctx)
	if err != nil {
		return nil, err
	}
	defer spanEnd()

	logger = logger.With("user_id", userID, "resource", resource)

	query := `
		SELECT DISTINCT 
			p.id, p.name, CONCAT(p.resource, ':', p.action) AS slug, p.resource, p.action, p.description, p.conditions, p.is_active
		FROM permissions p
		INNER JOIN role_permissions rp ON p.id = rp.permission_id AND rp.granted = 1
		INNER JOIN roles r ON rp.role_id = r.id AND r.is_active = 1
		INNER JOIN user_roles ur ON r.id = ur.role_id AND ur.is_active = 1
		WHERE ur.user_id = ? 
		  AND p.resource = ?
		  AND p.is_active = 1
		  AND (ur.expires_at IS NULL OR ur.expires_at > NOW())
		ORDER BY p.action
	`

	results, err := pa.Read(ctx, tx, query, userID, resource)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.permission.get_user_permissions_by_resource.read_error", "error", err)
		return nil, fmt.Errorf("get user permissions by resource read: %w", err)
	}

	permissions = make([]permissionmodel.PermissionInterface, 0, len(results))
	for index, row := range results {
		if len(row) != 8 {
			errColumns := fmt.Errorf("unexpected number of columns: expected 8, got %d", len(row))
			utils.SetSpanError(ctx, errColumns)
			logger.Error("mysql.permission.get_user_permissions_by_resource.columns_mismatch", "row_index", index, "error", errColumns)
			return nil, errColumns
		}

		entity := &permissionentities.PermissionEntity{
			ID:       row[0].(int64),
			Name:     string(row[1].([]byte)),
			Slug:     string(row[2].([]byte)),
			Resource: string(row[3].([]byte)),
			Action:   string(row[4].([]byte)),
			IsActive: row[7].(int64) == 1,
		}

		// description (pode ser NULL)
		if row[5] != nil {
			entity.Description = string(row[5].([]byte))
		}
		// conditions (pode ser NULL)
		if row[6] != nil {
			conditionsStr := string(row[6].([]byte))
			entity.Conditions = &conditionsStr
		}

		permission, convertErr := permissionconverters.PermissionEntityToDomain(entity)
		if convertErr != nil {
			utils.SetSpanError(ctx, convertErr)
			logger.Error("mysql.permission.get_user_permissions_by_resource.convert_error", "row_index", index, "error", convertErr)
			return nil, fmt.Errorf("convert permission entity to domain: %w", convertErr)
		}
		if permission != nil {
			permissions = append(permissions, permission)
		}
	}

	logger.Debug("mysql.permission.get_user_permissions_by_resource.success", "count", len(permissions))
	return permissions, nil
}
