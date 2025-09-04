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
	ctx, end, _ := utils.GenerateTracer(ctx)
	defer end()

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
	ctx, end, _ := utils.GenerateTracer(ctx)
	defer end()

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
	ctx, end, _ := utils.GenerateTracer(ctx)
	defer end()

	// Start a transaction for read operations when caller doesn't manage one
	tx, err := ps.globalService.StartTransaction(ctx)
	if err != nil {
		slog.Error("Failed to start transaction for temp block check", "userID", userID, "error", err)
		return false, utils.InternalError("Failed to start transaction")
	}
	defer func() {
		if err != nil {
			if rbErr := ps.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				slog.Error("Failed to rollback tx for temp block check", "userID", userID, "error", rbErr)
			}
		}
	}()

	blocked, ierr := ps.IsUserTempBlockedWithTx(ctx, tx, userID)
	if ierr != nil {
		return false, ierr
	}
	if err = ps.globalService.CommitTransaction(ctx, tx); err != nil {
		slog.Error("Failed to commit tx for temp block check", "userID", userID, "error", err)
		return false, utils.InternalError("Failed to commit transaction")
	}
	return blocked, nil
}

// IsUserTempBlockedWithTx checks if a user is temporarily blocked using the provided transaction
func (ps *permissionServiceImpl) IsUserTempBlockedWithTx(ctx context.Context, tx *sql.Tx, userID int64) (bool, error) {
	ctx, end, _ := utils.GenerateTracer(ctx)
	defer end()

	userRole, err := ps.permissionRepository.GetActiveUserRoleByUserID(ctx, tx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		slog.Error("Failed to get user role for temp block check", "userID", userID, "error", err)
		return false, utils.InternalError("Failed to check user status")
	}

	if userRole == nil {
		return false, nil
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
	ctx, end, _ := utils.GenerateTracer(ctx)
	defer end()

	tx, err := ps.globalService.StartTransaction(ctx)
	if err != nil {
		slog.Error("Failed to start transaction for getting expired temp blocked users", "error", err)
		return nil, utils.InternalError("Failed to start transaction")
	}
	defer func() {
		if err != nil {
			if rbErr := ps.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				slog.Error("Failed to rollback tx after error when getting expired temp blocked users", "error", rbErr)
			}
		} else {
			if cmErr := ps.globalService.CommitTransaction(ctx, tx); cmErr != nil {
				slog.Error("Failed to commit tx when getting expired temp blocked users", "error", cmErr)
			}
		}
	}()

	users, err := ps.permissionRepository.GetExpiredTempBlockedUsers(ctx, tx)
	if err != nil {
		slog.Error("Failed to get expired temp blocked users", "error", err)
		return nil, utils.InternalError("Failed to get expired blocked users")
	}

	return users, nil
}
