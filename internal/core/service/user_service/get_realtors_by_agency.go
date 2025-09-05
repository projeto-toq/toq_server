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
		utils.SetSpanError(ctx, err)
		slog.Error("user.get_realtors_by_agency.tx_start_error", "error", err)
		return nil, utils.InternalError("Failed to start transaction")
	}
	defer func() {
		if err != nil {
			if rbErr := us.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				slog.Error("user.get_realtors_by_agency.tx_rollback_error", "error", rbErr)
			}
		}
	}()

	realtors, err = us.getRealtorsByAgency(ctx, tx, agencyID)
	if err != nil {
		utils.SetSpanError(ctx, err)
		slog.Error("user.get_realtors_by_agency.read_realtors_error", "error", err, "agency_id", agencyID)
		return nil, utils.MapRepositoryError(err, "Realtors not found for agency")
	}

	// Commit the transaction
	err = us.globalService.CommitTransaction(ctx, tx)
	if err != nil {
		utils.SetSpanError(ctx, err)
		slog.Error("user.get_realtors_by_agency.tx_commit_error", "error", err)
		return nil, utils.InternalError("Failed to commit transaction")
	}
	return
}

func (us *userService) getRealtorsByAgency(ctx context.Context, tx *sql.Tx, agencyID int64) (realtors []usermodel.UserInterface, err error) {

	// Read the realtors user with given status from the database
	realtors, err = us.repo.GetRealtorsByAgency(ctx, tx, agencyID)
	if err != nil {
		utils.SetSpanError(ctx, err)
		slog.Error("user.get_realtors_by_agency.read_realtors_error", "error", err, "agency_id", agencyID)
		return
	}

	return
}
