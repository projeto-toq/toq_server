package userservices

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	globalmodel "github.com/projeto-toq/toq_server/internal/core/model/global_model"
	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// BlockUserTemporarily blocks a user by setting their status to StatusTempBlocked
func (us *userService) BlockUserTemporarily(ctx context.Context, userID int64) error {
	ctx, end, err := utils.GenerateTracer(ctx)
	if err != nil {
		ctx = utils.ContextWithLogger(ctx)
		utils.LoggerFromContext(ctx).Error("user.block_temp.tracer_error", "error", err)
		return utils.InternalError("Failed to generate tracer")
	}
	defer end()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	// Update all user roles to StatusTempBlocked
	err = us.repo.UpdateUserRoleStatusByUserID(ctx, userID, int(globalmodel.StatusTempBlocked))
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("user.block_temp.update_role_status_error", "user_id", userID, "error", err)
		return utils.InternalError(fmt.Sprintf("failed to block user %d temporarily", userID))
	}

	// Clear user permissions cache to force refresh on next signin attempt
	err = us.permissionService.ClearUserPermissionsCache(ctx, userID)
	if err != nil {
		logger.Warn("user.block_temp.clear_cache_failed", "user_id", userID, "error", err)
		// Don't return error here as the blocking was successful
	}

	logger.Info("user.block_temp.success", "user_id", userID)

	return nil
}

// UnblockUserTemporarily unblocks a user by setting their status back to StatusActive
func (us *userService) UnblockUserTemporarily(ctx context.Context, userID int64) error {
	ctx, end, err := utils.GenerateTracer(ctx)
	if err != nil {
		ctx = utils.ContextWithLogger(ctx)
		utils.LoggerFromContext(ctx).Error("user.unblock_temp.tracer_error", "error", err)
		return utils.InternalError("Failed to generate tracer")
	}
	defer end()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	// Update all user roles to StatusActive
	err = us.repo.UpdateUserRoleStatusByUserID(ctx, userID, int(globalmodel.StatusActive))
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("user.unblock_temp.update_role_status_error", "user_id", userID, "error", err)
		return utils.InternalError(fmt.Sprintf("failed to unblock user %d", userID))
	}

	// Reset wrong signin attempts counter
	err = us.repo.ResetUserWrongSigninAttempts(ctx, userID)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("user.unblock_temp.reset_wrong_signin_error", "user_id", userID, "error", err)
		return utils.InternalError(fmt.Sprintf("failed to reset signin attempts for user %d", userID))
	}

	// Clear user permissions cache to force refresh
	err = us.permissionService.ClearUserPermissionsCache(ctx, userID)
	if err != nil {
		logger.Warn("user.unblock_temp.clear_cache_failed", "user_id", userID, "error", err)
		// Don't return error here as the unblocking was successful
	}

	logger.Info("user.unblock_temp.success", "user_id", userID)

	return nil
}

// GetExpiredTempBlockedUsers returns all users whose temporary block has expired
func (us *userService) GetExpiredTempBlockedUsers(ctx context.Context) ([]usermodel.UserRoleInterface, error) {
	ctx, end, _ := utils.GenerateTracer(ctx)
	defer end()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	logger.Debug("permission.user.temp_block.get_expired.start")

	tx, err := us.globalService.StartTransaction(ctx)
	if err != nil {
		logger.Error("permission.user.temp_block.get_expired.tx_start_failed", "error", err)
		utils.SetSpanError(ctx, err)
		return nil, utils.InternalError("")
	}
	defer func() {
		if err != nil {
			if rbErr := us.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				logger.Error("permission.user.temp_block.get_expired.tx_rollback_failed", "error", rbErr)
				utils.SetSpanError(ctx, rbErr)
			}
		} else {
			if cmErr := us.globalService.CommitTransaction(ctx, tx); cmErr != nil {
				logger.Error("permission.user.temp_block.get_expired.tx_commit_failed", "error", cmErr)
				utils.SetSpanError(ctx, cmErr)
			}
		}
	}()

	users, err := us.repo.GetExpiredTempBlockedUsers(ctx, tx)
	if err != nil {
		logger.Error("permission.user.temp_block.get_expired.db_failed", "error", err)
		utils.SetSpanError(ctx, err)
		return nil, utils.InternalError("")
	}

	return users, nil
}

// IsUserTempBlocked checks if a user is temporarily blocked
func (us *userService) IsUserTempBlocked(ctx context.Context, userID int64) (bool, error) {
	ctx, end, _ := utils.GenerateTracer(ctx)
	defer end()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	logger.Debug("permission.user.temp_block.check.start", "user_id", userID)

	// Start a transaction for read operations when caller doesn't manage one
	tx, err := us.globalService.StartTransaction(ctx)
	if err != nil {
		logger.Error("permission.user.temp_block.check.tx_start_failed", "user_id", userID, "error", err)
		utils.SetSpanError(ctx, err)
		return false, utils.InternalError("")
	}
	defer func() {
		if err != nil {
			if rbErr := us.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				logger.Error("permission.user.temp_block.check.tx_rollback_failed", "user_id", userID, "error", rbErr)
				utils.SetSpanError(ctx, rbErr)
			}
		}
	}()

	blocked, ierr := us.IsUserTempBlockedWithTx(ctx, tx, userID)
	if ierr != nil {
		return false, ierr
	}
	if err = us.globalService.CommitTransaction(ctx, tx); err != nil {
		logger.Error("permission.user.temp_block.check.tx_commit_failed", "user_id", userID, "error", err)
		utils.SetSpanError(ctx, err)
		return false, utils.InternalError("")
	}
	return blocked, nil
}

// IsUserTempBlockedWithTx checks if a user is temporarily blocked using the provided transaction
func (us *userService) IsUserTempBlockedWithTx(ctx context.Context, tx *sql.Tx, userID int64) (bool, error) {
	ctx, end, _ := utils.GenerateTracer(ctx)
	defer end()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	userRole, err := us.repo.GetActiveUserRoleByUserID(ctx, tx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		logger.Error("permission.user.temp_block.check.db_failed", "user_id", userID, "error", err)
		utils.SetSpanError(ctx, err)
		return false, utils.InternalError("")
	}

	if userRole == nil {
		return false, nil
	}

	// Check if user is temp blocked and if the block hasn't expired
	if userRole.GetStatus() == globalmodel.StatusTempBlocked {
		blockedUntil := userRole.GetBlockedUntil()
		if blockedUntil != nil && time.Now().UTC().Before(*blockedUntil) {
			return true, nil
		}
	}

	return false, nil
}

// UnblockUser unblocks a user by changing their status back to StatusActive
func (us *userService) UnblockUser(ctx context.Context, tx *sql.Tx, userID int64) error {
	ctx, end, _ := utils.GenerateTracer(ctx)
	defer end()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	logger.Debug("permission.user.unblock.start", "user_id", userID)

	err := us.repo.UnblockUser(ctx, tx, userID)
	if err != nil {
		logger.Error("permission.user.unblock.db_failed", "user_id", userID, "error", err)
		utils.SetSpanError(ctx, err)
		return utils.InternalError("")
	}

	// Clear user permissions cache after status change
	if errCache := us.permissionService.ClearUserPermissionsCache(ctx, userID); errCache != nil {
		logger.Warn("permission.user.unblock.cache_clear_failed", "user_id", userID, "error", errCache)
	}

	logger.Info("permission.user.unblocked", "user_id", userID)
	return nil
}
