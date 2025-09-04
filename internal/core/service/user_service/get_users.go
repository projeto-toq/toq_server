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
		slog.Error("user.get_users.tx_start_error", "err", txErr)
		return nil, utils.InternalError("Failed to start transaction")
	}
	defer func() {
		if err != nil {
			if rbErr := us.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				slog.Error("user.get_users.tx_rollback_error", "err", rbErr)
			}
		}
	}()

	users, err = us.repo.GetUsers(ctx, tx)
	if err != nil {
		return
	}

	if cmErr := us.globalService.CommitTransaction(ctx, tx); cmErr != nil {
		slog.Error("user.get_users.tx_commit_error", "err", cmErr)
		return nil, utils.InternalError("Failed to commit transaction")
	}

	return
}
