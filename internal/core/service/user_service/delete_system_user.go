package userservices

import (
	"context"
	"time"

	derrors "github.com/projeto-toq/toq_server/internal/core/derrors"
	globalmodel "github.com/projeto-toq/toq_server/internal/core/model/global_model"
	permissionmodel "github.com/projeto-toq/toq_server/internal/core/model/permission_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// DeleteSystemUser performs logical deletion of a system user and deactivates all their roles.
// User data and role history are preserved for audit purposes.
func (us *userService) DeleteSystemUser(ctx context.Context, input DeleteSystemUserInput) error {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return utils.InternalError("Failed to generate tracer")
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	if input.UserID <= 0 {
		return utils.ValidationError("userId", "User id must be positive")
	}

	if currentID, cerr := us.globalService.GetUserIDFromContext(ctx); cerr == nil && currentID == input.UserID {
		return derrors.ErrCannotDeleteLoggedUser
	}

	tx, txErr := us.globalService.StartTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("admin.users.delete.tx_start_failed", "error", txErr)
		return utils.InternalError("")
	}

	var opErr error
	defer func() {
		if opErr != nil {
			if rbErr := us.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				logger.Error("admin.users.delete.tx_rollback_failed", "error", rbErr)
			}
		}
	}()

	existing, userErr := us.repo.GetUserByID(ctx, tx, input.UserID)
	if userErr != nil {
		utils.SetSpanError(ctx, userErr)
		logger.Error("admin.users.delete.get_user_failed", "user_id", input.UserID, "error", userErr)
		if errorsIsNoRows(userErr) {
			opErr = utils.NotFoundError("user")
		} else {
			opErr = utils.InternalError("")
		}
		return opErr
	}

	if existing.IsDeleted() {
		opErr = derrors.ErrUserAlreadyDeleted
		return opErr
	}

	activeRole := existing.GetActiveRole()
	if activeRole == nil || activeRole.GetRole() == nil || !activeRole.GetRole().GetIsSystemRole() {
		opErr = derrors.ErrSystemUserRoleMismatch
		return opErr
	}

	roleSlug := permissionmodel.RoleSlug(activeRole.GetRole().GetSlug())
	if roleSlug == "" {
		opErr = derrors.ErrSystemUserRoleMismatch
		return opErr
	}

	if us.sessionRepo != nil {
		if revokeErr := us.sessionRepo.RevokeSessionsByUserID(ctx, tx, existing.GetID()); revokeErr != nil {
			logger.Warn("admin.users.delete.revoke_sessions_failed", "user_id", existing.GetID(), "error", revokeErr)
		}
		if deleteErr := us.sessionRepo.DeleteSessionsByUserID(ctx, tx, existing.GetID()); deleteErr != nil {
			utils.SetSpanError(ctx, deleteErr)
			logger.Error("admin.users.delete.delete_sessions_failed", "user_id", existing.GetID(), "error", deleteErr)
			opErr = utils.InternalError("")
			return opErr
		}
	}

	if removeTokensErr := us.deviceTokenRepo.RemoveAllByUserID(existing.GetID()); removeTokensErr != nil {
		logger.Warn("admin.users.delete.remove_tokens_failed", "user_id", existing.GetID(), "error", removeTokensErr)
	}

	// Mark user as deleted (preserving all data for audit)
	existing.SetDeleted(true)
	existing.SetLastActivityAt(time.Now().UTC())

	if updateErr := us.repo.UpdateUserByID(ctx, tx, existing); updateErr != nil {
		utils.SetSpanError(ctx, updateErr)
		logger.Error("admin.users.delete.update_user_failed", "user_id", existing.GetID(), "error", updateErr)
		opErr = utils.InternalError("")
		return opErr
	}

	// Deactivate all user roles (soft delete: is_active = 0 for all roles)
	if deactivateErr := us.repo.DeactivateAllUserRoles(ctx, tx, existing.GetID()); deactivateErr != nil {
		utils.SetSpanError(ctx, deactivateErr)
		logger.Error("admin.users.delete.deactivate_roles_failed", "user_id", existing.GetID(), "error", deactivateErr)
		opErr = utils.InternalError("")
		return opErr
	}

	if auditErr := us.globalService.CreateAudit(ctx, tx, globalmodel.TableUsers, "System user deleted (data and role history preserved for audit)", existing.GetID()); auditErr != nil {
		utils.SetSpanError(ctx, auditErr)
		logger.Error("admin.users.delete.audit_failed", "user_id", existing.GetID(), "error", auditErr)
		opErr = utils.InternalError("")
		return opErr
	}

	if commitErr := us.globalService.CommitTransaction(ctx, tx); commitErr != nil {
		utils.SetSpanError(ctx, commitErr)
		logger.Error("admin.users.delete.tx_commit_failed", "user_id", existing.GetID(), "error", commitErr)
		return utils.InternalError("")
	}

	logger.Info("admin.users.delete.success", "user_id", existing.GetID())
	return nil
}
