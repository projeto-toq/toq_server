package userservices

import (
	"context"
	"database/sql"

	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (us *userService) UpdateProfile(ctx context.Context, user usermodel.UserInterface) (err error) {
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

	err = us.updateProfile(ctx, tx, user)
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

func (us *userService) updateProfile(
	ctx context.Context,
	tx *sql.Tx,
	user usermodel.UserInterface,
) (err error) {
	//recover the user before update it
	current, err := us.repo.GetUserByID(ctx, tx, user.GetID())
	if err != nil {
		return
	}

	//update the current with the new data
	current.SetNickName(user.GetNickName())
	current.SetBornAt(user.GetBornAt())
	current.SetZipCode(user.GetZipCode())
	current.SetStreet(user.GetStreet())
	current.SetNumber(user.GetNumber())
	current.SetComplement(user.GetComplement())
	current.SetNeighborhood(user.GetNeighborhood())
	current.SetCity(user.GetCity())
	current.SetState(user.GetState())

	err = us.repo.UpdateUserByID(ctx, tx, current)
	if err != nil {
		return
	}

	err = us.globalService.CreateAudit(ctx, tx, globalmodel.TableUsers, "Usuário atualizou o perfil")
	if err != nil {
		return
	}

	return
}
