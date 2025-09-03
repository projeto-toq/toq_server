package userservices

import (
	"context"
	"database/sql"
	"fmt"

	permissionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/permission_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (us *userService) SwitchUserRole(ctx context.Context, roleSlug permissionmodel.RoleSlug) (tokens usermodel.Tokens, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	// Obter o ID do usuário do contexto (SSOT)
	userID, err := us.globalService.GetUserIDFromContext(ctx)
	if err != nil || userID == 0 {
		return tokens, utils.ErrInternalServer
	}

	// Start transaction
	tx, err := us.globalService.StartTransaction(ctx)
	if err != nil {
		return
	}

	tokens, err = us.switchUserRole(ctx, tx, userID, roleSlug)
	if err != nil {
		us.globalService.RollbackTransaction(ctx, tx)
		return
	}

	err = us.globalService.CommitTransaction(ctx, tx)
	if err != nil {
		us.globalService.RollbackTransaction(ctx, tx)
		return
	}

	return
}

func (us *userService) switchUserRole(ctx context.Context, tx *sql.Tx, userID int64, roleSlug permissionmodel.RoleSlug) (tokens usermodel.Tokens, err error) {
	// Verificar se o usuário tem múltiplos roles usando permission service
	userRoles, err := us.permissionService.GetUserRolesWithTx(ctx, tx, userID)
	if err != nil {
		return
	}

	if len(userRoles) == 1 {
		err = utils.ErrInternalServer
		return
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
		err = utils.NewHTTPError(400, fmt.Sprintf("Role '%s' not found for user", roleSlug))
		return
	}

	// Usar permission service diretamente para buscar role
	role, roleErr := us.permissionService.GetRoleBySlugWithTx(ctx, tx, roleSlug)
	if roleErr != nil {
		err = roleErr
		return
	}

	err = us.permissionService.SwitchActiveRoleWithTx(ctx, tx, userID, role.GetID())
	if err != nil {
		return
	}

	// Buscar usuário atualizado
	user, err := us.repo.GetUserByID(ctx, tx, userID)
	if err != nil {
		return
	}

	// Gerar novos tokens
	tokens, err = us.CreateTokens(ctx, tx, user, false)
	if err != nil {
		return
	}

	return
}
