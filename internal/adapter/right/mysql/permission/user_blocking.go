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

	// O schema de user_roles não contém status_reason nem updated_at; manter apenas campos válidos.
	query := `
		UPDATE user_roles 
		SET status = ?, blocked_until = ?
		WHERE user_id = ? AND is_active = 1
	`

	if _, err = tx.ExecContext(ctx, query, permissionmodel.StatusTempBlocked, blockedUntil, userID); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.permission.block_user_temporarily.exec_error", "error", err)
		return fmt.Errorf("block user temporarily exec: %w", err)
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

	// Remover colunas inexistentes no schema (status_reason, updated_at).
	query := `
		UPDATE user_roles 
		SET status = ?, blocked_until = NULL
		WHERE user_id = ? AND status = ? AND is_active = 1
	`

	if _, err = tx.ExecContext(ctx, query,
		permissionmodel.StatusActive,
		userID,
		permissionmodel.StatusTempBlocked,
	); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.permission.unblock_user.exec_error", "error", err)
		return fmt.Errorf("unblock user exec: %w", err)
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

	rows, err := tx.QueryContext(ctx, query, permissionmodel.StatusTempBlocked)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.permission.get_expired_temp_blocked_users.query_error", "error", err)
		return nil, fmt.Errorf("query expired temp blocked users: %w", err)
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
