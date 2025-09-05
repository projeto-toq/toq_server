package userservices

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/giulio-alfieri/toq_server/internal/core/events"
	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	permissionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/permission_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

// DeleteAccount deletes the current authenticated user's account.
// It revokes all sessions and removes all device tokens, masks PII, and removes roles.
// Idempotent: if already deleted, returns success and expired tokens.
func (us *userService) DeleteAccount(ctx context.Context) (tokens usermodel.Tokens, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	// Resolve userID from context
	userID, err := us.globalService.GetUserIDFromContext(ctx)
	if err != nil || userID == 0 {
		return tokens, utils.AuthenticationError("")
	}

	tx, txErr := us.globalService.StartTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		slog.Error("user.delete_account.tx_start_error", "error", txErr)
		return tokens, utils.InternalError("Failed to start transaction")
	}
	defer func() {
		if err != nil {
			if rbErr := us.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				slog.Error("user.delete_account.tx_rollback_error", "error", rbErr)
			}
		}
	}()

	// Load user first to support idempotency and side-effects in a single txn
	user, err := us.repo.GetUserByID(ctx, tx, userID)
	if err != nil {
		utils.SetSpanError(ctx, err)
		slog.Error("user.delete_account.get_user_error", "error", err, "user_id", userID)
		return
	}

	// Idempotency: if already deleted, just revoke any remaining sessions and return expired tokens
	if user.IsDeleted() {
		if us.sessionRepo != nil {
			if err2 := us.sessionRepo.RevokeSessionsByUserID(ctx, tx, userID); err2 != nil {
				slog.Warn("user.delete_account.revoke_sessions_warning", "error", err2, "user_id", userID)
			} else {
				us.globalService.GetEventBus().Publish(events.SessionEvent{Type: events.SessionsRevoked, UserID: userID})
			}
		}
		// Best-effort: remove device tokens
		if err2 := us.repo.RemoveAllDeviceTokens(ctx, tx, userID); err2 != nil {
			slog.Warn("user.delete_account.remove_device_tokens_warning", "error", err2, "user_id", userID)
		}
		// Return expired tokens
		tokens, err = us.CreateTokens(ctx, tx, user, true)
		if err != nil {
			return
		}
		if err = us.globalService.CommitTransaction(ctx, tx); err != nil {
			utils.SetSpanError(ctx, err)
			slog.Error("user.delete_account.tx_commit_error", "error", err)
			// best-effort: return success even if commit logging shows failure in edge case
		}
		return
	}

	tokens, err = us.deleteAccount(ctx, tx, userID)
	if err != nil {
		utils.SetSpanError(ctx, err)
		slog.Error("user.delete_account.delete_account_error", "error", err, "user_id", userID)
		return
	}

	if err = us.globalService.CommitTransaction(ctx, tx); err != nil {
		utils.SetSpanError(ctx, err)
		slog.Error("user.delete_account.tx_commit_error", "error", err)
		return tokens, utils.InternalError("Failed to commit transaction")
	}

	return
}

func (us *userService) deleteAccount(ctx context.Context, tx *sql.Tx, userId int64) (tokens usermodel.Tokens, err error) {
	user, err := us.repo.GetUserByID(ctx, tx, userId)
	if err != nil {
		utils.SetSpanError(ctx, err)
		slog.Error("user.delete_account.get_user_error", "error", err, "user_id", userId)
		return
	}
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
			slog.Warn("user.delete_account.revoke_sessions_warning", "error", err2, "user_id", user.GetID())
		} else {
			us.globalService.GetEventBus().Publish(events.SessionEvent{Type: events.SessionsRevoked, UserID: user.GetID()})
		}
	}

	// Remove all device tokens
	if err2 := us.repo.RemoveAllDeviceTokens(ctx, tx, user.GetID()); err2 != nil {
		slog.Warn("user.delete_account.remove_device_tokens_warning", "error", err2, "user_id", user.GetID())
	}

	us.setDeletedData(user)

	err = us.repo.UpdateUserByID(ctx, tx, user)
	if err != nil {
		utils.SetSpanError(ctx, err)
		slog.Error("user.delete_account.update_user_error", "error", err, "user_id", user.GetID())
		return
	}

	err = us.repo.UpdateUserPasswordByID(ctx, tx, user)
	if err != nil {
		utils.SetSpanError(ctx, err)
		slog.Error("user.delete_account.update_password_error", "error", err, "user_id", user.GetID())
		return
	}

	err = us.globalService.CreateAudit(ctx, tx, globalmodel.TableUsers, "Mascarado dados do usuário (conta apagada)")
	if err != nil {
		utils.SetSpanError(ctx, err)
		slog.Error("user.delete_account.audit_users_error", "error", err, "user_id", user.GetID())
		return
	}

	// Delete user folder in cloud storage
	if us.cloudStorageService != nil {
		folderErr := us.DeleteUserFolder(ctx, user.GetID())
		if folderErr != nil {
			// Best-effort: already marked on span inside DeleteUserFolder; just warn here
			slog.Warn("user.delete_account.delete_user_folder_warning", "error", folderErr, "user_id", user.GetID())
		}
	}

	// Remover todos os roles do usuário
	userRoles, err := us.permissionService.GetUserRolesWithTx(ctx, tx, user.GetID())
	if err != nil {
		utils.SetSpanError(ctx, err)
		slog.Error("user.delete_account.get_user_roles_error", "error", err, "user_id", user.GetID())
		return
	}

	for _, userRole := range userRoles {
		if userRole.GetRole() != nil {
			err = us.permissionService.RemoveRoleFromUserWithTx(ctx, tx, user.GetID(), userRole.GetRole().GetID())
			if err != nil {
				utils.SetSpanError(ctx, err)
				slog.Error("user.delete_account.remove_role_error", "error", err, "user_id", user.GetID(), "role_id", userRole.GetRole().GetID())
				return
			}
		}
	}

	err = us.globalService.CreateAudit(ctx, tx, globalmodel.TableUserRoles, "Apagados os papéis do usuário, pois a conta foi apagada")
	if err != nil {
		utils.SetSpanError(ctx, err)
		slog.Error("user.delete_account.audit_user_roles_error", "error", err, "user_id", user.GetID())
		return
	}

	// Generate expired tokens to ensure client logout on all devices
	tokens, err = us.CreateTokens(ctx, tx, user, true)
	if err != nil {
		return
	}

	return
}
