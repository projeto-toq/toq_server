package permissionservice

import (
	"context"
	"database/sql"
	"log/slog"

	permissionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/permission_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

// GetUserRoles returns all roles of a user, independent of is_active
func (p *permissionServiceImpl) GetUserRoles(ctx context.Context, userID int64) ([]permissionmodel.UserRoleInterface, error) {
	ctx, end, _ := utils.GenerateTracer(ctx)
	defer end()

	if userID <= 0 {
		return nil, utils.BadRequest("invalid user id")
	}

	// Start transaction
	tx, err := p.globalService.StartTransaction(ctx)
	if err != nil {
		slog.Error("permission.user_roles.tx_start_failed", "user_id", userID, "error", err)
		return nil, utils.WrapDomainErrorWithSource(utils.InternalError("Failed to start transaction"))
	}
	defer func() {
		if err != nil {
			if rbErr := p.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				slog.Error("permission.user_roles.tx_rollback_failed", "user_id", userID, "error", rbErr)
			}
		}
	}()

	// Busca todas as roles do usuário (ativas e inativas); a regra de negócio prevê apenas uma ativa.
	userRoles, err := p.permissionRepository.GetUserRolesByUserID(ctx, tx, userID)
	if err != nil {
		slog.Error("permission.user_roles.db_failed", "user_id", userID, "error", err)
		return nil, utils.WrapDomainErrorWithSource(utils.InternalError("Failed to get user roles"))
	}

	// Commit the transaction
	if err = p.globalService.CommitTransaction(ctx, tx); err != nil {
		slog.Error("permission.user_roles.tx_commit_failed", "user_id", userID, "error", err)
		return nil, utils.WrapDomainErrorWithSource(utils.InternalError("Failed to commit transaction"))
	}

	return userRoles, nil
}

// GetUserRolesWithTx returns all roles of a user within a provided transaction (used in flows)
func (p *permissionServiceImpl) GetUserRolesWithTx(ctx context.Context, tx *sql.Tx, userID int64) ([]permissionmodel.UserRoleInterface, error) {
	ctx, end, _ := utils.GenerateTracer(ctx)
	defer end()

	if userID <= 0 {
		return nil, utils.BadRequest("invalid user id")
	}

	// Busca todas as roles do usuário (ativas e inativas); a regra de negócio prevê apenas uma ativa.
	userRoles, err := p.permissionRepository.GetUserRolesByUserID(ctx, tx, userID)
	if err != nil {
		return nil, utils.WrapDomainErrorWithSource(utils.InternalError("Failed to get user roles"))
	}

	return userRoles, nil
}
