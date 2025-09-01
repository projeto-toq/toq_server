package permissionservice

import (
	"context"
	"database/sql"

	permissionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/permission_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

// GetUserPermissions retorna todas as permissões de um usuário
func (p *permissionServiceImpl) GetUserPermissions(ctx context.Context, userID int64) ([]permissionmodel.PermissionInterface, error) {
	if userID <= 0 {
		return nil, utils.ErrBadRequest
	}

	// Start transaction
	tx, err := p.globalService.StartTransaction(ctx)
	if err != nil {
		return nil, utils.ErrInternalServer
	}

	permissions, err := p.permissionRepository.GetUserPermissions(ctx, tx, userID)
	if err != nil {
		err = p.globalService.RollbackTransaction(ctx, tx)
		return nil, utils.ErrInternalServer
	}

	// Commit the transaction
	err = p.globalService.CommitTransaction(ctx, tx)
	if err != nil {
		err = p.globalService.RollbackTransaction(ctx, tx)
		return nil, utils.ErrInternalServer
	}

	return permissions, nil
}

// GetUserPermissionsWithTx retorna todas as permissões de um usuário (com transação - uso em fluxos)
func (p *permissionServiceImpl) GetUserPermissionsWithTx(ctx context.Context, tx *sql.Tx, userID int64) ([]permissionmodel.PermissionInterface, error) {
	if userID <= 0 {
		return nil, utils.ErrBadRequest
	}

	permissions, err := p.permissionRepository.GetUserPermissions(ctx, tx, userID)
	if err != nil {
		return nil, utils.ErrInternalServer
	}

	return permissions, nil
}
