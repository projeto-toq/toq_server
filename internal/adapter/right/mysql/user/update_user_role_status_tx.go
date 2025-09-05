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
		JOIN roles r ON r.id = ur.role_id
		SET ur.status = ?
		WHERE ur.user_id = ? AND ur.is_active = 1 AND r.slug = ?`

	var execer interface {
		ExecContext(context.Context, string, ...any) (sql.Result, error)
	}
	if tx != nil {
		execer = tx
	} else {
		execer = ua.db.GetDB()
	}

	res, err := execer.ExecContext(ctx, q, int(status), userID, role)
	if err != nil {
		slog.Error("mysqluseradapter/UpdateUserRoleStatus: error executing update", "userID", userID, "role", role, "status", status, "error", err)
		return err
	}
	if res != nil {
		if rows, rerr := res.RowsAffected(); rerr == nil && rows == 0 {
			// Nenhuma linha atualizada indica ausência de papel ativo com o slug informado
			slog.Warn("mysqluseradapter/UpdateUserRoleStatus: no rows affected", "userID", userID, "role", role, "status", status)
			return sql.ErrNoRows
		}
	}
	slog.Debug("Updated user role status (tx)", "userID", userID, "role", role, "status", status)
	return nil
}
