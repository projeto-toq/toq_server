package userservices

import (
	"context"
	"database/sql"
	"fmt"

	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	
	
"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (us *userService) SwitchUserRole(ctx context.Context, userID int64, roleSlug string) (tokens usermodel.Tokens, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

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

func (us *userService) switchUserRole(ctx context.Context, tx *sql.Tx, userID int64, roleSlug string) (tokens usermodel.Tokens, err error) {
	// Verificar se o usuário tem múltiplos roles
	userRoles, err := us.repo.GetUserRolesByUserID(ctx, tx, userID)
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
		if userRole.GetRole() == usermodel.RoleRealtor && roleSlug == "realtor" {
			roleExists = true
			break
		}
		if userRole.GetRole() == usermodel.RoleOwner && roleSlug == "owner" {
			roleExists = true
			break
		}
		if userRole.GetRole() == usermodel.RoleAgency && roleSlug == "agency" {
			roleExists = true
			break
		}
	}

	if !roleExists {
		err = utils.NewHTTPError(400, fmt.Sprintf("Role '%s' not found for user", roleSlug))
		return
	}

	// Mapear roleSlug para roleID e usar permission service
	role, roleErr := us.getRoleBySlug(ctx, roleSlug)
	if roleErr != nil {
		err = roleErr
		return
	}

	err = us.permissionService.SwitchActiveRole(ctx, userID, role.GetID())
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
