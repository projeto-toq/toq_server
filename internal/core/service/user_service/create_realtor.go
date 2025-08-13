package userservices

import (
	"context"
	"database/sql"

	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (us *userService) CreateRealtor(ctx context.Context, realtor usermodel.UserInterface) (tokens usermodel.Tokens, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	tx, err := us.globalService.StartTransaction(ctx)
	if err != nil {
		return
	}

	tokens, err = us.createRealtor(ctx, tx, realtor)
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

func (us *userService) createRealtor(ctx context.Context, tx *sql.Tx, realtor usermodel.UserInterface) (tokens usermodel.Tokens, err error) {

	err = us.ValidateUserData(ctx, tx, realtor, usermodel.RoleRealtor)
	if err != nil {
		return
	}

	err = us.repo.CreateUser(ctx, tx, realtor)
	if err != nil {
		return
	}

	err = us.AddFirstUserRole(ctx, tx, realtor, usermodel.RoleRealtor)
	if err != nil {
		return
	}

	err = us.CreateUserValidations(ctx, tx, realtor)
	if err != nil {
		return
	}

	err = us.googleCloudService.CreateUserBucket(ctx, realtor.GetID())
	if err != nil {
		return
	}

	tokens, err = us.CreateTokens(ctx, tx, realtor, false)
	if err != nil {
		return
	}

	err = us.globalService.CreateAudit(ctx, tx, globalmodel.TableUsers, "Criado novo usuário com papel de Corretor", realtor.GetID())
	if err != nil {
		return
	}

	return
}
