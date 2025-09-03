package userservices

import (
	"context"
	"database/sql"

	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

// GetProfile retrieves a user's profile by their ID.
// It generates a signed URL for the user's photo if it exists.
func (us *userService) GetProfile(ctx context.Context) (user usermodel.UserInterface, err error) {
	// Obter o ID do usuário do contexto (SSOT)
	userID, err := us.globalService.GetUserIDFromContext(ctx)
	if err != nil || userID == 0 {
		return nil, utils.ErrInternalServer
	}
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	// Iniciar uma transação somente leitura para otimizar o fluxo de leitura
	tx, err := us.globalService.StartReadOnlyTransaction(ctx)
	if err != nil {
		return
	}

	user, err = us.repo.GetUserByID(ctx, tx, userID)
	if err != nil {
		us.globalService.RollbackTransaction(ctx, tx)
		// Translate repository errors into DomainError for adapter serialization
		if err == sql.ErrNoRows {
			return nil, utils.NotFoundError("User")
		}
		return nil, utils.InternalError("Failed to get user by ID")
	}

	err = us.globalService.CommitTransaction(ctx, tx)
	if err != nil {
		us.globalService.RollbackTransaction(ctx, tx)
		return
	}

	return
}
