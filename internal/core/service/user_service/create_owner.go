package userservices

import (
	"context"
	"database/sql"

	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (us *userService) CreateOwner(ctx context.Context, owner usermodel.UserInterface) (tokens usermodel.Tokens, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	tx, err := us.globalService.StartTransaction(ctx)
	if err != nil {
		return
	}

	tokens, err = us.createOwner(ctx, tx, owner)
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

func (us *userService) createOwner(ctx context.Context, tx *sql.Tx, owner usermodel.UserInterface) (tokens usermodel.Tokens, err error) {

	legacyRole, err := us.getLegacyRoleBySlug("owner")
	if err != nil {
		return
	}

	err = us.ValidateUserData(ctx, tx, owner, legacyRole)
	if err != nil {
		return
	}

	err = us.repo.CreateUser(ctx, tx, owner)
	if err != nil {
		return
	}

	err = us.assignRoleToUser(ctx, tx, owner.GetID(), "owner")
	if err != nil {
		return
	}

	err = us.CreateUserValidations(ctx, tx, owner)
	if err != nil {
		return
	}

	err = us.cloudStorageService.CreateUserFolder(ctx, owner.GetID())
	if err != nil {
		return
	}

	tokens, err = us.CreateTokens(ctx, tx, owner, false)
	if err != nil {
		return
	}

	err = us.globalService.CreateAudit(ctx, tx, globalmodel.TableUsers, "Criado novo usuário com papel de Proprietário", owner.GetID())
	if err != nil {
		return
	}

	return
}
