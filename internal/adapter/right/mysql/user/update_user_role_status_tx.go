package mysqluseradapter

import (
	"context"
	"database/sql"
	"fmt"

	permissionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/permission_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

// UpdateUserRoleStatus updates the active user role status for a specific role using a transaction.
// It affects only the active role row matching the provided role slug.
func (ua *UserAdapter) UpdateUserRoleStatus(ctx context.Context, tx *sql.Tx, userID int64, role permissionmodel.RoleSlug, status permissionmodel.UserRoleStatus) error {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

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
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.user.update_user_role_status_tx.update_error", "user_id", userID, "role", role, "status", status, "error", err)
		return fmt.Errorf("update user role status: %w", err)
	}
	if res != nil {
		if rows, rerr := res.RowsAffected(); rerr == nil {
			if rows == 0 {
				// Nenhuma linha atualizada indica ausência de papel ativo com o slug informado
				errNoRows := sql.ErrNoRows
				utils.SetSpanError(ctx, errNoRows)
				logger.Error("mysql.user.update_user_role_status_tx.no_rows", "user_id", userID, "role", role, "status", status, "error", errNoRows)
				return errNoRows
			}
		} else {
			logger.Warn("mysql.user.update_user_role_status_tx.rows_affected_warning", "user_id", userID, "role", role, "status", status, "error", rerr)
		}
	}
	logger.Debug("mysql.user.update_user_role_status_tx.success", "user_id", userID, "role", role, "status", status)
	return nil
}
