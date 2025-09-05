package userservices

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	permissionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/permission_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (us *userService) SwitchUserRole(ctx context.Context, roleSlug permissionmodel.RoleSlug) (tokens usermodel.Tokens, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return tokens, utils.InternalError("Failed to generate tracer")
	}
	defer spanEnd()

	// Obter o ID do usuário do contexto (SSOT)
	userID, err := us.globalService.GetUserIDFromContext(ctx)
	if err != nil || userID == 0 {
		return tokens, utils.AuthenticationError("")
	}

	// Start transaction
	tx, txErr := us.globalService.StartTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		slog.Error("user.switch_role.tx_start_error", "error", txErr)
		return tokens, utils.InternalError("Failed to start transaction")
	}
	defer func() {
		if err != nil {
			if rbErr := us.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				slog.Error("user.switch_role.tx_rollback_error", "error", rbErr)
			}
		}
	}()

	tokens, err = us.switchUserRole(ctx, tx, userID, roleSlug)
	if err != nil {
		return
	}

	if commitErr := us.globalService.CommitTransaction(ctx, tx); commitErr != nil {
		utils.SetSpanError(ctx, commitErr)
		slog.Error("user.switch_role.tx_commit_error", "error", commitErr)
		return tokens, utils.InternalError("Failed to commit transaction")
	}

	return
}

func (us *userService) switchUserRole(ctx context.Context, tx *sql.Tx, userID int64, roleSlug permissionmodel.RoleSlug) (tokens usermodel.Tokens, err error) {
	// Verificar se o usuário tem múltiplos roles usando permission service
	userRoles, err := us.permissionService.GetUserRolesWithTx(ctx, tx, userID)
	if err != nil {
		utils.SetSpanError(ctx, err)
		slog.Error("user.switch_role.get_user_roles_error", "error", err, "user_id", userID)
		return
	}

	if len(userRoles) == 1 {
		// Only one role available; switching is not applicable
		return tokens, utils.ConflictError("User has only one role; cannot switch")
	}

	// Verificar se o role solicitado existe para o usuário
	roleExists := false
	for _, userRole := range userRoles {
		// Usar a role do banco diretamente para comparar com o slug fornecido
		role := userRole.GetRole()
		if role != nil {
			roleSlugFromDB := permissionmodel.RoleSlug(role.GetSlug())
			if roleSlugFromDB == roleSlug {
				roleExists = true
				break
			}
		}
	}

	if !roleExists {
		return tokens, utils.NotFoundError(fmt.Sprintf("role '%s' for user", roleSlug))
	}

	// Usar permission service diretamente para buscar role
	role, roleErr := us.permissionService.GetRoleBySlugWithTx(ctx, tx, roleSlug)
	if roleErr != nil {
		utils.SetSpanError(ctx, roleErr)
		slog.Error("user.switch_role.get_role_error", "error", roleErr, "role_slug", roleSlug)
		err = roleErr
		return
	}

	err = us.permissionService.SwitchActiveRoleWithTx(ctx, tx, userID, role.GetID())
	if err != nil {
		utils.SetSpanError(ctx, err)
		slog.Error("user.switch_role.switch_active_role_error", "error", err, "user_id", userID, "role_id", role.GetID())
		return
	}

	// Buscar usuário atualizado
	user, err := us.repo.GetUserByID(ctx, tx, userID)
	if err != nil {
		utils.SetSpanError(ctx, err)
		slog.Error("user.switch_role.get_user_error", "error", err, "user_id", userID)
		return
	}

	// Gerar novos tokens
	tokens, err = us.CreateTokens(ctx, tx, user, false)
	if err != nil {
		utils.SetSpanError(ctx, err)
		slog.Error("user.switch_role.create_tokens_error", "error", err, "user_id", userID)
		return
	}

	return
}
