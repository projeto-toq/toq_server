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
		slog.Error("user.get_profile.tx_start_error", "err", txErr)
		return nil, utils.InternalError("Failed to start transaction")
	}
	defer func() {
		if err != nil {
			if rbErr := us.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				slog.Error("user.get_profile.tx_rollback_error", "err", rbErr)
			}
		}
	}()

	user, err = us.repo.GetUserByID(ctx, tx, userID)
	if err != nil {
		// Translate repository errors into DomainError for adapter serialization
		if err == sql.ErrNoRows {
			return nil, utils.NotFoundError("User")
		}
		return nil, utils.InternalError("Failed to get user by ID")
	}

	if cmErr := us.globalService.CommitTransaction(ctx, tx); cmErr != nil {
		slog.Error("user.get_profile.tx_commit_error", "err", cmErr)
		return nil, utils.InternalError("Failed to commit transaction")
	}

	return
}
