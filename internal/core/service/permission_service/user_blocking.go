package permissionservice

import (
	"context"
	"database/sql"
	"log/slog"
	"time"

	permissionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/permission_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

// BlockUserTemporarily blocks a user temporarily by changing their status to StatusTempBlocked
func (ps *permissionServiceImpl) BlockUserTemporarily(ctx context.Context, tx *sql.Tx, userID int64, reason string) error {
	blockedUntil := time.Now().UTC().Add(usermodel.TempBlockDuration)

	err := ps.permissionRepository.BlockUserTemporarily(ctx, tx, userID, blockedUntil, reason)
	if err != nil {
		slog.Error("Failed to block user temporarily", "userID", userID, "error", err)
		return utils.InternalError("Failed to block user")
	}

	// Clear user permissions cache after status change
	if errCache := ps.ClearUserPermissionsCache(ctx, userID); errCache != nil {
		slog.Warn("Failed to clear user permissions cache after blocking", "userID", userID, "error", errCache)
	}

	slog.Info("User blocked temporarily", "userID", userID, "blockedUntil", blockedUntil, "reason", reason)
	return nil
}

// UnblockUser unblocks a user by changing their status back to StatusActive
func (ps *permissionServiceImpl) UnblockUser(ctx context.Context, tx *sql.Tx, userID int64) error {
	err := ps.permissionRepository.UnblockUser(ctx, tx, userID)
	if err != nil {
		slog.Error("Failed to unblock user", "userID", userID, "error", err)
		return utils.InternalError("Failed to unblock user")
	}

	// Clear user permissions cache after status change
	if errCache := ps.ClearUserPermissionsCache(ctx, userID); errCache != nil {
		slog.Warn("Failed to clear user permissions cache after unblocking", "userID", userID, "error", errCache)
	}

	slog.Info("User unblocked successfully", "userID", userID)
	return nil
}

// IsUserTempBlocked checks if a user is temporarily blocked
func (ps *permissionServiceImpl) IsUserTempBlocked(ctx context.Context, userID int64) (bool, error) {
	userRole, err := ps.permissionRepository.GetActiveUserRoleByUserID(ctx, nil, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		slog.Error("Failed to get user role for temp block check", "userID", userID, "error", err)
		return false, utils.InternalError("Failed to check user status")
	}

	// Check if user is temp blocked and if the block hasn't expired
	if userRole.GetStatus() == permissionmodel.StatusTempBlocked {
		blockedUntil := userRole.GetBlockedUntil()
		if blockedUntil != nil && time.Now().UTC().Before(*blockedUntil) {
			return true, nil
		}
	}

	return false, nil
}

// GetExpiredTempBlockedUsers returns all users whose temporary block has expired
func (ps *permissionServiceImpl) GetExpiredTempBlockedUsers(ctx context.Context) ([]permissionmodel.UserRoleInterface, error) {
	tx, err := ps.globalService.StartTransaction(ctx)
	if err != nil {
		slog.Error("Failed to start transaction for getting expired temp blocked users", "error", err)
		return nil, utils.InternalError("Failed to start transaction")
	}
	defer func() {
		if err != nil {
			err = ps.globalService.RollbackTransaction(ctx, tx)
		} else {
			err = ps.globalService.CommitTransaction(ctx, tx)
		}
	}()

	users, err := ps.permissionRepository.GetExpiredTempBlockedUsers(ctx, tx)
	if err != nil {
		slog.Error("Failed to get expired temp blocked users", "error", err)
		return nil, utils.InternalError("Failed to get expired blocked users")
	}

	return users, nil
}
