package permissionservice

import (
	"context"
	"database/sql"

	permissionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/permission_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

// GetUserRoles retorna todos os roles ativos de um usuário
func (p *permissionServiceImpl) GetUserRoles(ctx context.Context, userID int64) ([]permissionmodel.UserRoleInterface, error) {
	if userID <= 0 {
		return nil, utils.ErrBadRequest
	}

	// Start transaction
	tx, err := p.globalService.StartTransaction(ctx)
	if err != nil {
		return nil, utils.ErrInternalServer
	}

	userRoles, err := p.permissionRepository.GetActiveUserRolesByUserID(ctx, tx, userID)
	if err != nil {
		err = p.globalService.RollbackTransaction(ctx, tx)
		return nil, utils.ErrInternalServer
	}

	// Commit the transaction
	err = p.globalService.CommitTransaction(ctx, tx)
	if err != nil {
		p.globalService.RollbackTransaction(ctx, tx)
		return nil, utils.ErrInternalServer
	}

	return userRoles, nil
}

// GetUserRolesWithTx retorna todos os roles ativos de um usuário (com transação - uso em fluxos)
func (p *permissionServiceImpl) GetUserRolesWithTx(ctx context.Context, tx *sql.Tx, userID int64) ([]permissionmodel.UserRoleInterface, error) {
	if userID <= 0 {
		return nil, utils.ErrBadRequest
	}

	userRoles, err := p.permissionRepository.GetActiveUserRolesByUserID(ctx, tx, userID)
	if err != nil {
		return nil, utils.ErrInternalServer
	}

	return userRoles, nil
}
