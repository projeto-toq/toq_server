package mysqlpermissionadapter

import (
	"context"
	"database/sql"
	"log/slog"
	"time"

	permissionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/permission_model"
)

// BlockUserTemporarily blocks a user temporarily by updating their user_role status and blocked_until
func (pa *PermissionAdapter) BlockUserTemporarily(ctx context.Context, tx *sql.Tx, userID int64, blockedUntil time.Time, reason string) error {
	// O schema de user_roles não contém status_reason nem updated_at; manter apenas campos válidos.
	query := `
		UPDATE user_roles 
		SET status = ?, blocked_until = ?
		WHERE user_id = ? AND is_active = 1
	`

	_, err := tx.ExecContext(ctx, query, permissionmodel.StatusTempBlocked, blockedUntil, userID)
	if err != nil {
		slog.Error("Failed to block user temporarily", "userID", userID, "error", err)
		return err
	}

	slog.Info("User blocked temporarily in database", "userID", userID, "blockedUntil", blockedUntil, "reason", reason)
	return nil
}

// UnblockUser unblocks a user by setting their status back to active and clearing blocked_until
func (pa *PermissionAdapter) UnblockUser(ctx context.Context, tx *sql.Tx, userID int64) error {
	// Remover colunas inexistentes no schema (status_reason, updated_at).
	query := `
		UPDATE user_roles 
		SET status = ?, blocked_until = NULL
		WHERE user_id = ? AND status = ? AND is_active = 1
	`

	_, err := tx.ExecContext(ctx, query,
		permissionmodel.StatusActive,
		userID,
		permissionmodel.StatusTempBlocked,
	)
	if err != nil {
		slog.Error("Failed to unblock user", "userID", userID, "error", err)
		return err
	}

	slog.Info("User unblocked in database", "userID", userID)
	return nil
}

// GetExpiredTempBlockedUsers returns users whose temporary block has expired
func (pa *PermissionAdapter) GetExpiredTempBlockedUsers(ctx context.Context, tx *sql.Tx) ([]permissionmodel.UserRoleInterface, error) {
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
		slog.Error("Failed to query expired temp blocked users", "error", err)
		return nil, err
	}
	defer rows.Close()

	var userRoles []permissionmodel.UserRoleInterface

	for rows.Next() {
		userRole := permissionmodel.NewUserRole()
		var id, userID, roleID int64
		var isActive bool
		var status permissionmodel.UserRoleStatus
		var expiresAt, blockedUntil *time.Time

		err := rows.Scan(&id, &userID, &roleID, &isActive, &status, &expiresAt, &blockedUntil)
		if err != nil {
			slog.Error("Failed to scan expired temp blocked user row", "error", err)
			continue
		}

		userRole.SetID(id)
		userRole.SetUserID(userID)
		userRole.SetRoleID(roleID)
		userRole.SetIsActive(isActive)
		userRole.SetStatus(status)
		userRole.SetExpiresAt(expiresAt)
		userRole.SetBlockedUntil(blockedUntil)

		userRoles = append(userRoles, userRole)
	}

	if err := rows.Err(); err != nil {
		slog.Error("Error iterating expired temp blocked users rows", "error", err)
		return nil, err
	}

	slog.Debug("Found expired temp blocked users", "count", len(userRoles))
	return userRoles, nil
}
