package userservices

import (
	"context"
	"database/sql"

	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (us *userService) GetRealtorsByAgency(ctx context.Context, agencyID int64) (realtors []usermodel.UserInterface, err error) {

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

	realtors, err = us.getRealtorsByAgency(ctx, tx, agencyID)
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

func (us *userService) getRealtorsByAgency(ctx context.Context, tx *sql.Tx, agencyID int64) (realtors []usermodel.UserInterface, err error) {

	// Read the realtors user with given status from the database
	realtors, err = us.repo.GetRealtorsByAgency(ctx, tx, agencyID)
	if err != nil {
		return
	}

	return
}
