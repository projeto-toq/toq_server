package mysqlpermissionadapter

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// HasUserPermission verifica se um usuário tem uma permissão específica de forma otimizada
func (p *PermissionAdapter) HasUserPermission(ctx context.Context, tx *sql.Tx, userID int64, resource, action string) (bool, error) {
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

	row := p.QueryRowContext(ctx, tx, "select", query, userID, resource, action)
	var exists int64
	if scanErr := row.Scan(&exists); scanErr != nil {
		if errors.Is(scanErr, sql.ErrNoRows) {
			logger.Debug("mysql.permission.has_user_permission.success", "has_permission", false)
			return false, nil
		}

		utils.SetSpanError(ctx, scanErr)
		logger.Error("mysql.permission.has_user_permission.scan_error", "error", scanErr)
		return false, fmt.Errorf("has user permission scan: %w", scanErr)
	}

	hasPermission := exists > 0
	logger.Debug("mysql.permission.has_user_permission.success", "has_permission", hasPermission)
	return hasPermission, nil
}
