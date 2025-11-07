package userservices

import (
	"context"
	"database/sql"

	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// GetActiveUserRole returns the single active role for the given user.
// It returns (nil, nil) when the user has no active role.
//
// Português: inicia tracing, valida entrada, abre transação via global service,
// delega ao repositório e padroniza erros via WrapDomainErrorWithSource.
func (us *userService) GetActiveUserRole(ctx context.Context, userID int64) (usermodel.UserRoleInterface, error) {
	ctx, end, _ := utils.GenerateTracer(ctx)
	defer end()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	if userID <= 0 {
		return nil, utils.BadRequest("invalid user id")
	}

	// Start transaction
	tx, err := us.globalService.StartTransaction(ctx)
	if err != nil {
		logger.Error("permission.user_role.active.tx_start_failed", "user_id", userID, "error", err)
		utils.SetSpanError(ctx, err)
		return nil, utils.InternalError("")
	}
	defer func() {
		if err != nil {
			if rbErr := us.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				logger.Error("permission.user_role.active.tx_rollback_failed", "user_id", userID, "error", rbErr)
				utils.SetSpanError(ctx, rbErr)
			}
		}
	}()

	userRole, repoErr := us.repo.GetActiveUserRoleByUserID(ctx, tx, userID)
	if repoErr != nil {
		logger.Error("permission.user_role.active.db_failed", "user_id", userID, "error", repoErr)
		utils.SetSpanError(ctx, repoErr)
		err = repoErr
		return nil, utils.InternalError("")
	}

	// Commit the transaction
	if err = us.globalService.CommitTransaction(ctx, tx); err != nil {
		logger.Error("permission.user_role.active.tx_commit_failed", "user_id", userID, "error", err)
		utils.SetSpanError(ctx, err)
		return nil, utils.InternalError("")
	}

	return userRole, nil
}

// GetActiveUserRoleWithTx returns the single active role for the given user using the provided transaction.
// It returns (nil, nil) when the user has no active role.
//
// Português: variante com transação do chamador; mantém tracing e erros padronizados.
func (us *userService) GetActiveUserRoleWithTx(ctx context.Context, tx *sql.Tx, userID int64) (usermodel.UserRoleInterface, error) {
	ctx, end, _ := utils.GenerateTracer(ctx)
	defer end()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	if userID <= 0 {
		return nil, utils.BadRequest("invalid user id")
	}

	userRole, repoErr := us.repo.GetActiveUserRoleByUserID(ctx, tx, userID)
	if repoErr != nil {
		utils.SetSpanError(ctx, repoErr)
		logger.Error("permission.user_role.active.db_failed", "user_id", userID, "error", repoErr)
		return nil, utils.InternalError("")
	}
	return userRole, nil
}
