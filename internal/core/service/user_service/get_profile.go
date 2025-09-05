package userservices

import (
	"context"
	"database/sql"
	"log/slog"

	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

// GetProfile retrieves a user's profile by their ID.
// It generates a signed URL for the user's photo if it exists.
func (us *userService) GetProfile(ctx context.Context) (user usermodel.UserInterface, err error) {
	// Obter o ID do usuário do contexto (SSOT)
	userID, err := us.globalService.GetUserIDFromContext(ctx)
	if err != nil || userID == 0 {
		return nil, utils.InternalError("Failed to get environment")
	}
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, utils.InternalError("Failed to generate tracer")
	}
	defer spanEnd()

	// Iniciar uma transação somente leitura para otimizar o fluxo de leitura
	tx, txErr := us.globalService.StartReadOnlyTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		slog.Error("user.get_profile.tx_start_error", "error", txErr)
		return nil, utils.InternalError("Failed to start transaction")
	}
	defer func() {
		if err != nil {
			if rbErr := us.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				slog.Error("user.get_profile.tx_rollback_error", "error", rbErr)
			}
		}
	}()

	user, err = us.repo.GetUserByID(ctx, tx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, utils.NotFoundError("User")
		}
		utils.SetSpanError(ctx, err)
		slog.Error("user.get_profile.read_user_error", "error", err, "user_id", userID)
		return nil, utils.InternalError("Failed to get user by ID")
	}

	if cmErr := us.globalService.CommitTransaction(ctx, tx); cmErr != nil {
		utils.SetSpanError(ctx, cmErr)
		slog.Error("user.get_profile.tx_commit_error", "error", cmErr)
		return nil, utils.InternalError("Failed to commit transaction")
	}

	return
}
