package userservices

import (
	"context"
	"database/sql"
	"fmt"

	permissionmodel "github.com/projeto-toq/toq_server/internal/core/model/permission_model"
	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (us *userService) SwitchUserRole(ctx context.Context, roleSlug permissionmodel.RoleSlug) (tokens usermodel.Tokens, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return tokens, utils.InternalError("Failed to generate tracer")
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	// Obter o ID do usuário do contexto (SSOT)
	userID, err := us.globalService.GetUserIDFromContext(ctx)
	if err != nil || userID == 0 {
		return tokens, utils.AuthenticationError("")
	}

	// Start transaction
	tx, txErr := us.globalService.StartTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("user.switch_role.tx_start_error", "error", txErr)
		return tokens, utils.InternalError("Failed to start transaction")
	}
	defer func() {
		if err != nil {
			if rbErr := us.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				logger.Error("user.switch_role.tx_rollback_error", "error", rbErr)
			}
		}
	}()

	tokens, err = us.switchUserRole(ctx, tx, userID, roleSlug)
	if err != nil {
		return
	}

	if commitErr := us.globalService.CommitTransaction(ctx, tx); commitErr != nil {
		utils.SetSpanError(ctx, commitErr)
		logger.Error("user.switch_role.tx_commit_error", "error", commitErr)
		return tokens, utils.InternalError("Failed to commit transaction")
	}

	return
}

func (us *userService) switchUserRole(ctx context.Context, tx *sql.Tx, userID int64, roleSlug permissionmodel.RoleSlug) (tokens usermodel.Tokens, err error) {
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	if roleSlug != permissionmodel.RoleSlugOwner && roleSlug != permissionmodel.RoleSlugRealtor {
		return tokens, utils.AuthorizationError("Only owners or realtors can switch roles")
	}

	// Recupera as informações dos roles permitidos (alvo e complementar)
	role, roleErr := us.permissionService.GetRoleBySlugWithTx(ctx, tx, roleSlug)
	if roleErr != nil {
		utils.SetSpanError(ctx, roleErr)
		logger.Error("user.switch_role.get_role_error", "error", roleErr, "role_slug", roleSlug)
		err = roleErr
		return
	}

	otherSlug := permissionmodel.RoleSlugOwner
	if roleSlug == permissionmodel.RoleSlugOwner {
		otherSlug = permissionmodel.RoleSlugRealtor
	}

	otherRole, otherRoleErr := us.permissionService.GetRoleBySlugWithTx(ctx, tx, otherSlug)
	if otherRoleErr != nil {
		utils.SetSpanError(ctx, otherRoleErr)
		logger.Error("user.switch_role.get_role_error", "error", otherRoleErr, "role_slug", otherSlug)
		err = otherRoleErr
		return
	}

	allowedRoleIDs := map[permissionmodel.RoleSlug]int64{
		roleSlug:  role.GetID(),
		otherSlug: otherRole.GetID(),
	}

	activeRole, activeErr := us.permissionService.GetActiveUserRoleWithTx(ctx, tx, userID)
	if activeErr != nil {
		utils.SetSpanError(ctx, activeErr)
		logger.Error("user.switch_role.get_active_role_error", "error", activeErr, "user_id", userID)
		err = activeErr
		return
	}

	if activeRole == nil {
		return tokens, utils.ErrUserActiveRoleMissing
	}

	currentRoleID := activeRole.GetRoleID()
	var currentSlug permissionmodel.RoleSlug
	for slug, id := range allowedRoleIDs {
		if id == currentRoleID {
			currentSlug = slug
			break
		}
	}

	if currentSlug == "" {
		return tokens, utils.AuthorizationError("Only owners or realtors can switch roles")
	}

	if currentRoleID == role.GetID() {
		return tokens, utils.BadRequest("Requested role is already active")
	}

	// Verifica se o usuário possui a role alvo cadastrada
	userRoles, rolesErr := us.permissionService.GetUserRolesWithTx(ctx, tx, userID)
	if rolesErr != nil {
		utils.SetSpanError(ctx, rolesErr)
		logger.Error("user.switch_role.get_user_roles_error", "error", rolesErr, "user_id", userID)
		err = rolesErr
		return
	}

	var hasTarget bool
	for _, userRole := range userRoles {
		if userRole.GetRoleID() == role.GetID() {
			hasTarget = true
			break
		}
	}

	if len(userRoles) < 2 || !hasTarget {
		return tokens, utils.BadRequest(fmt.Sprintf("User must have role '%s' assigned to switch", roleSlug))
	}

	if err = us.permissionService.SwitchActiveRoleWithTx(ctx, tx, userID, role.GetID()); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("user.switch_role.switch_active_role_error", "error", err, "user_id", userID, "role_id", role.GetID())
		return
	}

	// Buscar usuário atualizado
	user, err := us.repo.GetUserByID(ctx, tx, userID)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("user.switch_role.get_user_error", "error", err, "user_id", userID)
		return
	}

	// Gerar novos tokens
	tokens, err = us.CreateTokens(ctx, tx, user, false)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("user.switch_role.create_tokens_error", "error", err, "user_id", userID)
		return
	}

	logger.Info("user.switch_role.success", "user_id", userID, "from_role", currentSlug.String(), "to_role", roleSlug.String())

	return
}
