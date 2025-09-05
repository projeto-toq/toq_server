package permissionservice

import (
	"context"
	"database/sql"
	"log/slog"

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

	if userID <= 0 {
		return nil, utils.BadRequest("invalid user id")
	}

	// Start transaction
	tx, err := p.globalService.StartTransaction(ctx)
	if err != nil {
		slog.Error("permission.user_role.active.tx_start_failed", "user_id", userID, "error", err)
		return nil, utils.WrapDomainErrorWithSource(utils.InternalError("Failed to start transaction"))
	}
	defer func() {
		if err != nil {
			if rbErr := p.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				slog.Error("permission.user_role.active.tx_rollback_failed", "user_id", userID, "error", rbErr)
			}
		}
	}()

	userRole, repoErr := p.permissionRepository.GetActiveUserRoleByUserID(ctx, tx, userID)
	if repoErr != nil {
		slog.Error("permission.user_role.active.db_failed", "user_id", userID, "error", repoErr)
		err = repoErr
		return nil, utils.WrapDomainErrorWithSource(utils.InternalError("Failed to get active user role"))
	}

	// Commit the transaction
	if err = p.globalService.CommitTransaction(ctx, tx); err != nil {
		slog.Error("permission.user_role.active.tx_commit_failed", "user_id", userID, "error", err)
		return nil, utils.WrapDomainErrorWithSource(utils.InternalError("Failed to commit transaction"))
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

	if userID <= 0 {
		return nil, utils.BadRequest("invalid user id")
	}

	userRole, repoErr := p.permissionRepository.GetActiveUserRoleByUserID(ctx, tx, userID)
	if repoErr != nil {
		return nil, utils.WrapDomainErrorWithSource(utils.InternalError("Failed to get active user role"))
	}
	return userRole, nil
}
