package userservices

import (
	"context"
	"log/slog"

	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (us *userService) GetAgencyOfRealtor(ctx context.Context, realtorID int64) (agency usermodel.UserInterface, err error) {

	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, utils.InternalError("Failed to generate tracer")
	}
	defer spanEnd()

	// Start a database transaction
	tx, err := us.globalService.StartTransaction(ctx)
	if err != nil {
		slog.Error("user.get_agency_of_realtor.tx_start_error", "err", err)
		return nil, utils.InternalError("Failed to start transaction")
	}
	defer func() {
		if err != nil {
			if rbErr := us.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				slog.Error("user.get_agency_of_realtor.tx_rollback_error", "err", rbErr)
			}
		}
	}()

	agency, err = us.repo.GetAgencyOfRealtor(ctx, tx, realtorID)
	if err != nil {
		return nil, err
	}

	// Commit the transaction
	err = us.globalService.CommitTransaction(ctx, tx)
	if err != nil {
		slog.Error("user.get_agency_of_realtor.tx_commit_error", "err", err)
		return nil, utils.InternalError("Failed to commit transaction")
	}
	return
}
