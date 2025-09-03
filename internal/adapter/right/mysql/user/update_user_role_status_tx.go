package mysqluseradapter

import (
	"context"
	"database/sql"
	"log/slog"

	permissionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/permission_model"
)

// UpdateUserRoleStatus updates the active user role status for a specific role using a transaction.
// It affects only the active role row matching the provided role slug.
func (ua *UserAdapter) UpdateUserRoleStatus(ctx context.Context, tx *sql.Tx, userID int64, role permissionmodel.RoleSlug, status permissionmodel.UserRoleStatus) error {
	// Comentário: para evitar atualização incorreta, limitamos por is_active=1 e role slug.
	const q = `
        UPDATE user_roles ur
        JOIN base_roles br ON br.id = ur.role_id
        SET ur.status = ?, ur.updated_at = NOW()
        WHERE ur.user_id = ? AND ur.is_active = 1 AND br.slug = ?`

	var execer interface {
		ExecContext(context.Context, string, ...any) (sql.Result, error)
	}
	if tx != nil {
		execer = tx
	} else {
		execer = ua.db.GetDB()
	}

	if _, err := execer.ExecContext(ctx, q, int(status), userID, role); err != nil {
		slog.Error("mysqluseradapter/UpdateUserRoleStatus: error executing update", "userID", userID, "role", role, "status", status, "error", err)
		return err
	}
	slog.Debug("Updated user role status (tx)", "userID", userID, "role", role, "status", status)
	return nil
}
