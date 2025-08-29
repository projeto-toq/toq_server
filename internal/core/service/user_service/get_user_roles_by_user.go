package userservices

import (
	"context"
	"database/sql"

	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (us *userService) GetUserRolesByUser(ctx context.Context, userID int64) (roles []usermodel.UserRoleInterface, err error) {

	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	// Start a database transaction
	tx, err := us.globalService.StartTransaction(ctx)
	if err != nil {
		return
	}

	roles, err = us.getUserRolesByUser(ctx, tx, userID)
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

func (us *userService) getUserRolesByUser(ctx context.Context, tx *sql.Tx, userID int64) (roles []usermodel.UserRoleInterface, err error) {

	// Read the realtors user with given status from the database
	roles, err = us.repo.GetUserRolesByUserID(ctx, tx, userID)
	if err != nil {
		return
	}

	return
}
