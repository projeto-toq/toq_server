package mysqlpermissionadapter

import (
	"context"
	"database/sql"
)

// HasUserPermission verifica se um usuário tem uma permissão específica de forma otimizada
func (pa *PermissionAdapter) HasUserPermission(ctx context.Context, tx *sql.Tx, userID int64, resource, action string) (bool, error) {
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
		return false, err
	}

	return len(results) > 0, nil
}
