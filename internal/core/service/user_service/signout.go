package userservices

import (
	"context"
	"database/sql"

	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (us *userService) SignOut(ctx context.Context, userID int64) (tokens usermodel.Tokens, err error) {

	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	tx, err := us.globalService.StartTransaction(ctx)
	if err != nil {
		return
	}

	tokens, err = us.signOut(ctx, tx, userID)
	if err != nil {
		us.globalService.RollbackTransaction(ctx, tx)
		return
	}

	// Commit the transaction
	err = us.globalService.CommitTransaction(ctx, tx)
	if err != nil {
		us.globalService.RollbackTransaction(ctx, tx)
		return
	}

	return
}

func (us *userService) signOut(ctx context.Context, tx *sql.Tx, userID int64) (tokens usermodel.Tokens, err error) {

	user, err := us.repo.GetUserByID(ctx, tx, userID)
	if err != nil {
		return
	}

	//generate the expired token
	tokens, err = us.CreateTokens(ctx, tx, user, true)
	if err != nil {
		return
	}
	return
}
