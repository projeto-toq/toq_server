package userservices

import (
	"context"
	"database/sql"
	"errors"

	"github.com/projeto-toq/toq_server/internal/core/events"
	auditmodel "github.com/projeto-toq/toq_server/internal/core/model/audit_model"
	permissionmodel "github.com/projeto-toq/toq_server/internal/core/model/permission_model"
	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"
	auditservice "github.com/projeto-toq/toq_server/internal/core/service/audit_service"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// DeleteAccount deletes the current authenticated user's account.
// It revokes all sessions, removes device tokens, and marks the account as deleted.
// User data and role history are preserved for audit purposes.
// Idempotent: if already deleted, returns success and expired tokens.
func (us *userService) DeleteAccount(ctx context.Context) (tokens usermodel.Tokens, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		ctx = utils.ContextWithLogger(ctx)
		utils.LoggerFromContext(ctx).Error("user.delete_account.tracer_error", "error", err)
		return
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	// Resolve userID from context
	userID, err := us.globalService.GetUserIDFromContext(ctx)
	if err != nil || userID == 0 {
		return tokens, utils.AuthenticationError("")
	}

	tx, txErr := us.globalService.StartTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("user.delete_account.tx_start_error", "error", txErr)
		return tokens, utils.InternalError("Failed to start transaction")
	}
	defer func() {
		if err != nil {
			if rbErr := us.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				logger.Error("user.delete_account.tx_rollback_error", "error", rbErr)
			}
		}
	}()

	// Load user with active role (repository returns complete aggregate)
	user, err := us.repo.GetUserByID(ctx, tx, userID)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("user.delete_account.get_user_error", "error", err, "user_id", userID)
		return
	}

	// Validate domain invariant: repository already populated active role
	activeRole := user.GetActiveRole()
	if activeRole == nil || activeRole.GetRole() == nil {
		// Inconsistência interna: por domínio, usuário deve ter role ativa
		err = utils.InternalError("Active role missing unexpectedly")
		utils.SetSpanError(ctx, err)
		logger.Error("user.delete_account.active_role_missing", "user_id", userID)
		return tokens, err
	}

	var publishSessionsEvent bool
	tokens, publishSessionsEvent, err = us.deleteAccount(ctx, tx, user)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("user.delete_account.delete_account_error", "error", err, "user_id", userID)
		return
	}

	if err = us.globalService.CommitTransaction(ctx, tx); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("user.delete_account.tx_commit_error", "error", err)
		return tokens, utils.InternalError("Failed to commit transaction")
	}

	if publishSessionsEvent {
		us.globalService.GetEventBus().Publish(events.SessionEvent{Type: events.SessionsRevoked, UserID: userID})
	}

	return
}

func (us *userService) deleteAccount(ctx context.Context, tx *sql.Tx, user usermodel.UserInterface) (tokens usermodel.Tokens, publishSessionsEvent bool, err error) {
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	//delete the account dependencies
	activeRole := user.GetActiveRole()

	if activeRole != nil && activeRole.GetRole() != nil {
		roleSlug := permissionmodel.RoleSlug(activeRole.GetRole().GetSlug())
		switch roleSlug {
		case permissionmodel.RoleSlugOwner:
			err = us.CleanOwnerPending(ctx, user)
			if err != nil {
				// Do not mark span here; called service will handle its own infra errors
				return
			}
		case permissionmodel.RoleSlugRealtor:
			err = us.CleanRealtorPending(ctx, user)
			if err != nil {
				// Do not mark span here; called service will handle its own infra errors
				return
			}
		case permissionmodel.RoleSlugAgency:
			err = us.CleanAgencyPending(ctx, user)
			if err != nil {
				// Do not mark span here; called service will handle its own infra errors
				return
			}
		}
	}

	// Revoke all active sessions for the user (security first)
	if us.sessionRepo != nil {
		if err2 := us.sessionRepo.RevokeSessionsByUserID(ctx, tx, user.GetID()); err2 != nil {
			logger.Warn("user.delete_account.revoke_sessions_warning", "error", err2, "user_id", user.GetID())
		}
		if err2 := us.sessionRepo.DeleteSessionsByUserID(ctx, tx, user.GetID()); err2 != nil {
			utils.SetSpanError(ctx, err2)
			logger.Error("user.delete_account.delete_sessions_error", "error", err2, "user_id", user.GetID())
			return tokens, publishSessionsEvent, utils.InternalError("Failed to delete user sessions")
		}
		publishSessionsEvent = true
	}

	// Remove all device tokens
	if err2 := us.repo.RemoveAllDeviceTokensByUserID(ctx, tx, user.GetID()); err2 != nil {
		logger.Warn("user.delete_account.remove_device_tokens_warning", "error", err2, "user_id", user.GetID())
	}

	// Mark user as deleted (preserving all data for audit)
	user.SetDeleted(true)

	err = us.repo.UpdateUserByID(ctx, tx, user)
	if err != nil {
		// Handle sql.ErrNoRows as success: happens when MySQL UPDATE finds no changes
		// (user was loaded in same transaction, so user exists, just no fields changed)
		if !errors.Is(err, sql.ErrNoRows) {
			// Real infrastructure error
			utils.SetSpanError(ctx, err)
			logger.Error("user.delete_account.update_user_error", "error", err, "user_id", user.GetID())
			return tokens, publishSessionsEvent, utils.InternalError("Failed to mark user as deleted")
		}
		// No changes needed = success, continue
	}

	auditRecord := auditservice.BuildRecordFromContext(
		ctx,
		user.GetID(),
		auditmodel.AuditTarget{Type: auditmodel.TargetUser, ID: user.GetID()},
		auditmodel.OperationDelete,
		map[string]any{"requested_by": "self"},
	)
	if err = us.auditService.RecordChange(ctx, tx, auditRecord); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("user.delete_account.audit_users_error", "error", err, "user_id", user.GetID())
		return
	}

	// Deactivate all user roles (soft delete: is_active = 0 for all roles)
	if err2 := us.repo.DeactivateAllUserRoles(ctx, tx, user.GetID()); err2 != nil {
		utils.SetSpanError(ctx, err2)
		logger.Error("user.delete_account.deactivate_roles_error", "error", err2, "user_id", user.GetID())
		return tokens, publishSessionsEvent, utils.InternalError("Failed to deactivate user roles")
	}

	auditRoles := auditservice.BuildRecordFromContext(
		ctx,
		user.GetID(),
		auditmodel.AuditTarget{Type: auditmodel.TargetUserRole, ID: activeRole.GetID()},
		auditmodel.OperationStatusChange,
		map[string]any{
			"action":      "deactivate_all",
			"role_slug":   string(permissionmodel.RoleSlug(activeRole.GetRole().GetSlug())),
			"user_id":     user.GetID(),
			"active_role": activeRole.GetID(),
		},
	)
	if err = us.auditService.RecordChange(ctx, tx, auditRoles); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("user.delete_account.audit_user_roles_error", "error", err, "user_id", user.GetID())
		return
	}

	// Generate expired tokens to ensure client logout on all devices
	tokens, err = us.CreateTokens(ctx, tx, user, true)
	if err != nil {
		return
	}

	return
}
