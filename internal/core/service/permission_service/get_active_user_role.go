package permissionservice

import (
	"context"
	"database/sql"

	permissionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/permission_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

// GetActiveUserRole returns the single active role for the given user.
// It returns (nil, nil) when the user has no active role.
//
// Português: inicia tracing, valida entrada, abre transação via global service,
// delega ao repositório e padroniza erros via WrapDomainErrorWithSource.
func (p *permissionServiceImpl) GetActiveUserRole(ctx context.Context, userID int64) (permissionmodel.UserRoleInterface, error) {
	ctx, end, _ := utils.GenerateTracer(ctx)
	defer end()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	if userID <= 0 {
		return nil, utils.BadRequest("invalid user id")
	}

	// Start transaction
	tx, err := p.globalService.StartTransaction(ctx)
	if err != nil {
		logger.Error("permission.user_role.active.tx_start_failed", "user_id", userID, "error", err)
		utils.SetSpanError(ctx, err)
		return nil, utils.InternalError("")
	}
	defer func() {
		if err != nil {
			if rbErr := p.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				logger.Error("permission.user_role.active.tx_rollback_failed", "user_id", userID, "error", rbErr)
				utils.SetSpanError(ctx, rbErr)
			}
		}
	}()

	userRole, repoErr := p.permissionRepository.GetActiveUserRoleByUserID(ctx, tx, userID)
	if repoErr != nil {
		logger.Error("permission.user_role.active.db_failed", "user_id", userID, "error", repoErr)
		utils.SetSpanError(ctx, repoErr)
		err = repoErr
		return nil, utils.InternalError("")
	}

	// Commit the transaction
	if err = p.globalService.CommitTransaction(ctx, tx); err != nil {
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
func (p *permissionServiceImpl) GetActiveUserRoleWithTx(ctx context.Context, tx *sql.Tx, userID int64) (permissionmodel.UserRoleInterface, error) {
	ctx, end, _ := utils.GenerateTracer(ctx)
	defer end()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	if userID <= 0 {
		return nil, utils.BadRequest("invalid user id")
	}

	userRole, repoErr := p.permissionRepository.GetActiveUserRoleByUserID(ctx, tx, userID)
	if repoErr != nil {
		utils.SetSpanError(ctx, repoErr)
		logger.Error("permission.user_role.active.db_failed", "user_id", userID, "error", repoErr)
		return nil, utils.InternalError("")
	}
	return userRole, nil
}
