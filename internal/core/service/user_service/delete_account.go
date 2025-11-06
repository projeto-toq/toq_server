package userservices

import (
	"context"
	"database/sql"
	"errors"

	"github.com/projeto-toq/toq_server/internal/core/events"
	globalmodel "github.com/projeto-toq/toq_server/internal/core/model/global_model"
	permissionmodel "github.com/projeto-toq/toq_server/internal/core/model/permission_model"
	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// DeleteAccount deletes the current authenticated user's account.
// It revokes all sessions and removes all device tokens, masks PII, soft deletes the active role, and hard deletes the remaining roles.
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

	// Load user (somente usuários não deletados são retornados pelo repo)
	user, err := us.repo.GetUserByID(ctx, tx, userID)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("user.delete_account.get_user_error", "error", err, "user_id", userID)
		return
	}

	// Garantir role ativa carregada a partir do Permission Service
	activeRole, arErr := us.permissionService.GetActiveUserRoleWithTx(ctx, tx, userID)
	if arErr != nil {
		utils.SetSpanError(ctx, arErr)
		logger.Error("user.delete_account.get_active_role_error", "error", arErr, "user_id", userID)
		return tokens, utils.InternalError("")
	}
	if activeRole == nil || activeRole.GetRole() == nil {
		// Inconsistência interna: por domínio, usuário deve ter role ativa
		err = utils.InternalError("Active role missing unexpectedly")
		utils.SetSpanError(ctx, err)
		logger.Error("user.delete_account.active_role_missing", "user_id", userID)
		return tokens, err
	}
	user.SetActiveRole(activeRole)

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
	var activeRoleID int64
	var activeRoleSlug permissionmodel.RoleSlug
	if activeRole != nil && activeRole.GetRole() != nil {
		activeRoleID = activeRole.GetRole().GetID()
		activeRoleSlug = permissionmodel.RoleSlug(activeRole.GetRole().GetSlug())
	}

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
	if err2 := us.deviceTokenRepo.RemoveAllByUserID(user.GetID()); err2 != nil {
		logger.Warn("user.delete_account.remove_device_tokens_warning", "error", err2, "user_id", user.GetID())
	}

	us.setDeletedData(user)

	err = us.repo.UpdateUserByID(ctx, tx, user)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("user.delete_account.update_user_error", "error", err, "user_id", user.GetID())
		return
	}

	err = us.repo.UpdateUserPasswordByID(ctx, tx, user)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("user.delete_account.update_password_error", "error", err, "user_id", user.GetID())
		return
	}

	err = us.globalService.CreateAudit(ctx, tx, globalmodel.TableUsers, "Mascarado dados do usuário (conta apagada)")
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("user.delete_account.audit_users_error", "error", err, "user_id", user.GetID())
		return
	}

	// Delete user folder in cloud storage
	if us.cloudStorageService != nil {
		folderErr := us.DeleteUserFolder(ctx, user.GetID())
		if folderErr != nil {
			logger.Error("user.delete_account.delete_user_folder_error", "error", folderErr, "user_id", user.GetID())
			return tokens, publishSessionsEvent, utils.InternalError("Failed to delete user assets")
		}
	}

	// Remover todos os roles do usuário
	userRoles, err := us.permissionService.GetUserRolesWithTx(ctx, tx, user.GetID())
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("user.delete_account.get_user_roles_error", "error", err, "user_id", user.GetID())
		return
	}

	if err = us.markUserRolesAsDeleted(ctx, tx, user.GetID(), userRoles, activeRoleSlug); err != nil {
		return tokens, publishSessionsEvent, err
	}

	for _, userRole := range userRoles {
		role := userRole.GetRole()
		if role == nil {
			errMissingRole := utils.InternalError("Role details missing for user role")
			utils.SetSpanError(ctx, errMissingRole)
			logger.Error("user.delete_account.role_without_details", "user_id", user.GetID(), "user_role_id", userRole.GetID())
			return tokens, publishSessionsEvent, errMissingRole
		}

		if (activeRoleID != 0 && role.GetID() == activeRoleID) || (activeRoleSlug != "" && permissionmodel.RoleSlug(role.GetSlug()) == activeRoleSlug) {
			continue
		}

		err = us.permissionService.RemoveRoleFromUserWithTx(ctx, tx, user.GetID(), role.GetID())
		if err != nil {
			utils.SetSpanError(ctx, err)
			logger.Error("user.delete_account.remove_role_error", "error", err, "user_id", user.GetID(), "role_id", role.GetID())
			return
		}
	}

	err = us.globalService.CreateAudit(ctx, tx, globalmodel.TableUserRoles, "Apagados os papéis do usuário, pois a conta foi apagada")
	if err != nil {
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

func (us *userService) markUserRolesAsDeleted(ctx context.Context, tx *sql.Tx, userID int64, userRoles []permissionmodel.UserRoleInterface, activeRoleSlug permissionmodel.RoleSlug) error {
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	if activeRoleSlug == "" {
		for _, userRole := range userRoles {
			if !userRole.GetIsActive() {
				continue
			}
			role := userRole.GetRole()
			if role == nil {
				err := utils.InternalError("Active role details missing")
				utils.SetSpanError(ctx, err)
				logger.Error("user.delete_account.active_role_without_details", "user_id", userID, "user_role_id", userRole.GetID())
				return err
			}
			slug := permissionmodel.RoleSlug(role.GetSlug())
			if slug == "" {
				err := utils.InternalError("Active role slug missing")
				utils.SetSpanError(ctx, err)
				logger.Error("user.delete_account.active_role_without_slug", "user_id", userID, "user_role_id", userRole.GetID())
				return err
			}
			activeRoleSlug = slug
			break
		}
	}

	if activeRoleSlug == "" {
		err := utils.InternalError("Active role not found")
		utils.SetSpanError(ctx, err)
		logger.Error("user.delete_account.active_role_not_found", "user_id", userID)
		return err
	}

	if err := us.repo.UpdateUserRoleStatus(ctx, tx, userID, activeRoleSlug, permissionmodel.StatusDeleted); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.SetSpanError(ctx, err)
			logger.Error("user.delete_account.update_role_status_no_rows", "user_id", userID, "role_slug", activeRoleSlug)
			return utils.InternalError("Failed to update role status")
		}
		utils.SetSpanError(ctx, err)
		logger.Error("user.delete_account.update_role_status_error", "error", err, "user_id", userID, "role_slug", activeRoleSlug)
		return utils.InternalError("Failed to update role status")
	}

	logger.Info("user.delete_account.role_status_marked_deleted", "user_id", userID, "role_slug", activeRoleSlug)
	return nil
}
