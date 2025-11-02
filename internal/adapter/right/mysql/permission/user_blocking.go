package mysqlpermissionadapter

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	permissionmodel "github.com/projeto-toq/toq_server/internal/core/model/permission_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// BlockUserTemporarily blocks a user temporarily by updating their user_role status and blocked_until
func (pa *PermissionAdapter) BlockUserTemporarily(ctx context.Context, tx *sql.Tx, userID int64, blockedUntil time.Time, reason string) error {
	ctx, spanEnd, logger, err := startPermissionOperation(ctx)
	if err != nil {
		return err
	}
	defer spanEnd()

	logger = logger.With("user_id", userID, "blocked_until", blockedUntil, "reason", reason)

	query := `
		UPDATE user_roles 
		SET status = ?, blocked_until = ?
		WHERE user_id = ? AND is_active = 1
	`

	result, execErr := pa.ExecContext(ctx, tx, "update", query,
		permissionmodel.StatusTempBlocked,
		blockedUntil,
		userID,
	)
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.permission.block_user_temporarily.exec_error", "error", execErr)
		return fmt.Errorf("block user temporarily exec: %w", execErr)
	}

	if _, rowsErr := result.RowsAffected(); rowsErr != nil {
		utils.SetSpanError(ctx, rowsErr)
		logger.Error("mysql.permission.block_user_temporarily.rows_affected_error", "error", rowsErr)
		return fmt.Errorf("block user temporarily rows affected: %w", rowsErr)
	}

	logger.Debug("mysql.permission.block_user_temporarily.success")
	return nil
}

// UnblockUser unblocks a user by setting their status back to active and clearing blocked_until
func (pa *PermissionAdapter) UnblockUser(ctx context.Context, tx *sql.Tx, userID int64) error {
	ctx, spanEnd, logger, err := startPermissionOperation(ctx)
	if err != nil {
		return err
	}
	defer spanEnd()

	logger = logger.With("user_id", userID)

	query := `
		UPDATE user_roles 
		SET status = ?, blocked_until = NULL
		WHERE user_id = ? AND status = ? AND is_active = 1
	`

	result, execErr := pa.ExecContext(ctx, tx, "update", query,
		permissionmodel.StatusActive,
		userID,
		permissionmodel.StatusTempBlocked,
	)
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.permission.unblock_user.exec_error", "error", execErr)
		return fmt.Errorf("unblock user exec: %w", execErr)
	}

	if _, rowsErr := result.RowsAffected(); rowsErr != nil {
		utils.SetSpanError(ctx, rowsErr)
		logger.Error("mysql.permission.unblock_user.rows_affected_error", "error", rowsErr)
		return fmt.Errorf("unblock user rows affected: %w", rowsErr)
	}

	logger.Debug("mysql.permission.unblock_user.success")
	return nil
}

// GetExpiredTempBlockedUsers returns users whose temporary block has expired
func (pa *PermissionAdapter) GetExpiredTempBlockedUsers(ctx context.Context, tx *sql.Tx) ([]permissionmodel.UserRoleInterface, error) {
	ctx, spanEnd, logger, err := startPermissionOperation(ctx)
	if err != nil {
		return nil, err
	}
	defer spanEnd()

	query := `
		SELECT ur.id, ur.user_id, ur.role_id, ur.is_active, ur.status, ur.expires_at, ur.blocked_until
		FROM user_roles ur
		WHERE ur.status = ? 
		  AND ur.blocked_until IS NOT NULL 
		  AND ur.blocked_until <= NOW()
		  AND ur.is_active = 1
	`

	rows, queryErr := pa.QueryContext(ctx, tx, "select", query, permissionmodel.StatusTempBlocked)
	if queryErr != nil {
		utils.SetSpanError(ctx, queryErr)
		logger.Error("mysql.permission.get_expired_temp_blocked_users.query_error", "error", queryErr)
		return nil, fmt.Errorf("query expired temp blocked users: %w", queryErr)
	}
	defer rows.Close()

	var userRoles []permissionmodel.UserRoleInterface

	index := 0
	for rows.Next() {
		index++
		userRole := permissionmodel.NewUserRole()
		var (
			id           int64
			userID       int64
			roleID       int64
			isActive     bool
			status       permissionmodel.UserRoleStatus
			expiresAt    sql.NullTime
			blockedUntil sql.NullTime
		)

		if err := rows.Scan(&id, &userID, &roleID, &isActive, &status, &expiresAt, &blockedUntil); err != nil {
			utils.SetSpanError(ctx, err)
			logger.Error("mysql.permission.get_expired_temp_blocked_users.scan_error", "row_index", index-1, "error", err)
			return nil, fmt.Errorf("scan expired temp blocked user row: %w", err)
		}

		userRole.SetID(id)
		userRole.SetUserID(userID)
		userRole.SetRoleID(roleID)
		userRole.SetIsActive(isActive)
		userRole.SetStatus(status)

		if expiresAt.Valid {
			t := expiresAt.Time
			userRole.SetExpiresAt(&t)
		}
		if blockedUntil.Valid {
			t := blockedUntil.Time
			userRole.SetBlockedUntil(&t)
		}

		userRoles = append(userRoles, userRole)
	}

	if err := rows.Err(); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.permission.get_expired_temp_blocked_users.rows_error", "error", err)
		return nil, fmt.Errorf("iterate expired temp blocked users rows: %w", err)
	}

	logger.Debug("mysql.permission.get_expired_temp_blocked_users.success", "count", len(userRoles))
	return userRoles, nil
}
