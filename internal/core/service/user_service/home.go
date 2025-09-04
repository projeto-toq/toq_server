package userservices

import (
	"context"
	"log/slog"

	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (us *userService) Home(ctx context.Context, userID int64) (user usermodel.UserInterface, err error) {

	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, utils.InternalError("Failed to generate tracer")
	}
	defer spanEnd()

	tx, err := us.globalService.StartTransaction(ctx)
	if err != nil {
		slog.Error("user.home.tx_start_error", "err", err)
		return nil, utils.InternalError("Failed to start transaction")
	}
	defer func() {
		if err != nil {
			if rbErr := us.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				slog.Error("user.home.tx_rollback_error", "err", rbErr)
			}
		}
	}()

	user, err = us.repo.GetUserByID(ctx, tx, userID)
	if err != nil {
		slog.Error("Error getting user by ID", "error", err)
		return nil, utils.InternalError("Failed to get user")
	}

	// Commit the transaction
	err = us.globalService.CommitTransaction(ctx, tx)
	if err != nil {
		slog.Error("user.home.tx_commit_error", "err", err)
		return nil, utils.InternalError("Failed to commit transaction")
	}

	return
}
