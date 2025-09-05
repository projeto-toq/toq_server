package userservices

import (
	"context"
	"log/slog"

	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (us *userService) GetUsers(ctx context.Context) (users []usermodel.UserInterface, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, utils.InternalError("Failed to generate tracer")
	}
	defer spanEnd()

	// Start transaction
	tx, txErr := us.globalService.StartTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		slog.Error("user.get_users.tx_start_error", "error", txErr)
		return nil, utils.InternalError("Failed to start transaction")
	}
	defer func() {
		if err != nil {
			if rbErr := us.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				slog.Error("user.get_users.tx_rollback_error", "error", rbErr)
			}
		}
	}()

	users, err = us.repo.GetUsers(ctx, tx)
	if err != nil {
		utils.SetSpanError(ctx, err)
		slog.Error("user.get_users.read_users_error", "error", err)
		return nil, utils.MapRepositoryError(err, "Users not found")
	}

	if cmErr := us.globalService.CommitTransaction(ctx, tx); cmErr != nil {
		utils.SetSpanError(ctx, cmErr)
		slog.Error("user.get_users.tx_commit_error", "error", cmErr)
		return nil, utils.InternalError("Failed to commit transaction")
	}

	return
}
