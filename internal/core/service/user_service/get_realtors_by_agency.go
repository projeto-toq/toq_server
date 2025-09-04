package userservices

import (
	"context"
	"database/sql"
	"log/slog"

	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (us *userService) GetRealtorsByAgency(ctx context.Context, agencyID int64) (realtors []usermodel.UserInterface, err error) {

	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, utils.InternalError("Failed to generate tracer")
	}
	defer spanEnd()

	// Start a database transaction
	tx, err := us.globalService.StartTransaction(ctx)
	if err != nil {
		slog.Error("user.get_realtors_by_agency.tx_start_error", "err", err)
		return nil, utils.InternalError("Failed to start transaction")
	}
	defer func() {
		if err != nil {
			if rbErr := us.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				slog.Error("user.get_realtors_by_agency.tx_rollback_error", "err", rbErr)
			}
		}
	}()

	realtors, err = us.getRealtorsByAgency(ctx, tx, agencyID)
	if err != nil {
		return
	}

	// Commit the transaction
	err = us.globalService.CommitTransaction(ctx, tx)
	if err != nil {
		slog.Error("user.get_realtors_by_agency.tx_commit_error", "err", err)
		return nil, utils.InternalError("Failed to commit transaction")
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
