package userservices

import (
	"context"
	"database/sql"

	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (us *userService) UpdatePhoto(ctx context.Context, userID int64, photo []byte) (err error) {
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

	err = us.updatePhoto(ctx, tx, userID, photo)
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

func (us *userService) updatePhoto(ctx context.Context, tx *sql.Tx, userID int64, photo []byte) (err error) {
	user, err := us.repo.GetUserByID(ctx, tx, userID)
	if err != nil {
		return
	}

	user.SetPhoto(photo)

	err = us.repo.UpdateUserPhotoByID(ctx, tx, user)
	if err != nil {
		return
	}

	err = us.globalService.CreateAudit(ctx, tx, globalmodel.TableUsers, "Usuário atualizou a foto")
	if err != nil {
		return
	}
	return

}
