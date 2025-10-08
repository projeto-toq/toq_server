package mysqlpermissionadapter

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

// HasUserPermission verifica se um usuário tem uma permissão específica de forma otimizada
func (pa *PermissionAdapter) HasUserPermission(ctx context.Context, tx *sql.Tx, userID int64, resource, action string) (bool, error) {
	ctx, spanEnd, logger, err := startPermissionOperation(ctx)
	if err != nil {
		return false, err
	}
	defer spanEnd()

	logger = logger.With("user_id", userID, "resource", resource, "action", action)

	query := `
		SELECT 1
		FROM permissions p
		INNER JOIN role_permissions rp ON p.id = rp.permission_id AND rp.granted = 1
		INNER JOIN roles r ON rp.role_id = r.id AND r.is_active = 1
		INNER JOIN user_roles ur ON r.id = ur.role_id AND ur.is_active = 1
		WHERE ur.user_id = ? 
		  AND p.resource = ?
		  AND p.action = ?
		  AND p.is_active = 1
		  AND (ur.expires_at IS NULL OR ur.expires_at > NOW())
		LIMIT 1
	`

	results, err := pa.Read(ctx, tx, query, userID, resource, action)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.permission.has_user_permission.read_error", "error", err)
		return false, fmt.Errorf("has user permission read: %w", err)
	}

	hasPermission := len(results) > 0
	logger.Debug("mysql.permission.has_user_permission.success", "has_permission", hasPermission)
	return hasPermission, nil
}
