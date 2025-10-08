package permissionservice

import (
	"context"
	"database/sql"
	"time"

	permissionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/permission_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

// BlockUserTemporarily blocks a user temporarily by changing their status to StatusTempBlocked
func (ps *permissionServiceImpl) BlockUserTemporarily(ctx context.Context, tx *sql.Tx, userID int64, reason string) error {
	ctx, end, _ := utils.GenerateTracer(ctx)
	defer end()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	blockedUntil := time.Now().UTC().Add(usermodel.TempBlockDuration)

	logger.Debug("permission.user.block.start", "user_id", userID, "reason", reason)

	err := ps.permissionRepository.BlockUserTemporarily(ctx, tx, userID, blockedUntil, reason)
	if err != nil {
		logger.Error("permission.user.block.db_failed", "user_id", userID, "error", err)
		utils.SetSpanError(ctx, err)
		return utils.InternalError("")
	}

	// Clear user permissions cache after status change
	if errCache := ps.ClearUserPermissionsCache(ctx, userID); errCache != nil {
		logger.Warn("permission.user.block.cache_clear_failed", "user_id", userID, "error", errCache)
	}

	logger.Info("permission.user.blocked", "user_id", userID, "blocked_until", blockedUntil, "reason", reason)
	return nil
}

// UnblockUser unblocks a user by changing their status back to StatusActive
func (ps *permissionServiceImpl) UnblockUser(ctx context.Context, tx *sql.Tx, userID int64) error {
	ctx, end, _ := utils.GenerateTracer(ctx)
	defer end()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	logger.Debug("permission.user.unblock.start", "user_id", userID)

	err := ps.permissionRepository.UnblockUser(ctx, tx, userID)
	if err != nil {
		logger.Error("permission.user.unblock.db_failed", "user_id", userID, "error", err)
		utils.SetSpanError(ctx, err)
		return utils.InternalError("")
	}

	// Clear user permissions cache after status change
	if errCache := ps.ClearUserPermissionsCache(ctx, userID); errCache != nil {
		logger.Warn("permission.user.unblock.cache_clear_failed", "user_id", userID, "error", errCache)
	}

	logger.Info("permission.user.unblocked", "user_id", userID)
	return nil
}

// IsUserTempBlocked checks if a user is temporarily blocked
func (ps *permissionServiceImpl) IsUserTempBlocked(ctx context.Context, userID int64) (bool, error) {
	ctx, end, _ := utils.GenerateTracer(ctx)
	defer end()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	logger.Debug("permission.user.temp_block.check.start", "user_id", userID)

	// Start a transaction for read operations when caller doesn't manage one
	tx, err := ps.globalService.StartTransaction(ctx)
	if err != nil {
		logger.Error("permission.user.temp_block.check.tx_start_failed", "user_id", userID, "error", err)
		utils.SetSpanError(ctx, err)
		return false, utils.InternalError("")
	}
	defer func() {
		if err != nil {
			if rbErr := ps.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				logger.Error("permission.user.temp_block.check.tx_rollback_failed", "user_id", userID, "error", rbErr)
				utils.SetSpanError(ctx, rbErr)
			}
		}
	}()

	blocked, ierr := ps.IsUserTempBlockedWithTx(ctx, tx, userID)
	if ierr != nil {
		return false, ierr
	}
	if err = ps.globalService.CommitTransaction(ctx, tx); err != nil {
		logger.Error("permission.user.temp_block.check.tx_commit_failed", "user_id", userID, "error", err)
		utils.SetSpanError(ctx, err)
		return false, utils.InternalError("")
	}
	return blocked, nil
}

// IsUserTempBlockedWithTx checks if a user is temporarily blocked using the provided transaction
func (ps *permissionServiceImpl) IsUserTempBlockedWithTx(ctx context.Context, tx *sql.Tx, userID int64) (bool, error) {
	ctx, end, _ := utils.GenerateTracer(ctx)
	defer end()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	userRole, err := ps.permissionRepository.GetActiveUserRoleByUserID(ctx, tx, userID)
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

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	logger.Debug("permission.user.temp_block.get_expired.start")

	tx, err := ps.globalService.StartTransaction(ctx)
	if err != nil {
		logger.Error("permission.user.temp_block.get_expired.tx_start_failed", "error", err)
		utils.SetSpanError(ctx, err)
		return nil, utils.InternalError("")
	}
	defer func() {
		if err != nil {
			if rbErr := ps.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				logger.Error("permission.user.temp_block.get_expired.tx_rollback_failed", "error", rbErr)
				utils.SetSpanError(ctx, rbErr)
			}
		} else {
			if cmErr := ps.globalService.CommitTransaction(ctx, tx); cmErr != nil {
				logger.Error("permission.user.temp_block.get_expired.tx_commit_failed", "error", cmErr)
				utils.SetSpanError(ctx, cmErr)
			}
		}
	}()

	users, err := ps.permissionRepository.GetExpiredTempBlockedUsers(ctx, tx)
	if err != nil {
		logger.Error("permission.user.temp_block.get_expired.db_failed", "error", err)
		utils.SetSpanError(ctx, err)
		return nil, utils.InternalError("")
	}

	return users, nil
}
