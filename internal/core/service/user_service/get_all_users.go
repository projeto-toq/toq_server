package userservices

import (
	"context"

	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// ListAllUsers retrieves all active users from the repository
//
// This is the service layer wrapper for the repository's ListAllUsers method.
// Handles transaction lifecycle (start, commit, rollback) and error mapping.
//
// Naming Convention: Follows Section 8.1.4 of guide - List* prefix for collection retrieval
//
// Parameters:
//   - ctx: Context for tracing, logging, cancellation
//
// Returns:
//   - users: Slice of UserInterface domain models (deleted=0 only)
//   - err: Mapped error via utils.MapRepositoryError or transaction errors
//
// Transaction Scope:
//   - Starts read transaction
//   - Calls repo.ListAllUsers(ctx, tx)
//   - Commits on success, rolls back on error
//
// Error Handling:
//   - sql.ErrNoRows → "Users not found" (404 equivalent)
//   - Other DB errors → Mapped via utils.MapRepositoryError
//   - Transaction errors → utils.InternalError (500)
func (us *userService) ListAllUsers(ctx context.Context) (users []usermodel.UserInterface, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, utils.InternalError("Failed to generate tracer")
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	tx, txErr := us.globalService.StartTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("user.list_all_users.tx_start_error", "error", txErr)
		return nil, utils.InternalError("Failed to start transaction")
	}
	defer func() {
		if err != nil {
			if rbErr := us.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				logger.Error("user.list_all_users.tx_rollback_error", "error", rbErr)
			}
		}
	}()

	users, err = us.repo.ListAllUsers(ctx, tx)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("user.list_all_users.read_users_error", "error", err)
		return nil, utils.MapRepositoryError(err, "Users not found")
	}

	if cmErr := us.globalService.CommitTransaction(ctx, tx); cmErr != nil {
		utils.SetSpanError(ctx, cmErr)
		logger.Error("user.list_all_users.tx_commit_error", "error", cmErr)
		return nil, utils.InternalError("Failed to commit transaction")
	}

	return
}
