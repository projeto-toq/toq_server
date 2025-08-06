package userservices

import (
	"context"
	"database/sql"

	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (us *userService) CreateRoot(ctx context.Context, root usermodel.UserInterface) (err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	tx, err := us.globalService.StartTransaction(ctx)
	if err != nil {
		return
	}

	err = us.createRoot(ctx, tx, root)
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

func (us *userService) createRoot(ctx context.Context, tx *sql.Tx, root usermodel.UserInterface) (err error) {

	root.SetPassword(us.encryptPassword(root.GetPassword()))
	err = us.repo.CreateUser(ctx, tx, root)
	if err != nil {
		return
	}

	err = us.AddFirstUserRole(ctx, tx, root, usermodel.RoleRoot)
	if err != nil {
		return
	}

	err = us.globalService.CreateAudit(ctx, tx, globalmodel.TableUsers, "Criado novo usuário com papel de Administrador", root.GetID())
	if err != nil {
		return
	}

	return
}
