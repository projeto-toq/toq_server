package permissionservice

import (
	"context"
	"log/slog"

	permissionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/permission_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

// GrantPermissionToRole concede uma permissão a um role
func (p *permissionServiceImpl) GrantPermissionToRole(ctx context.Context, roleID, permissionID int64) error {
	if roleID <= 0 {
		return utils.ErrBadRequest
	}

	if permissionID <= 0 {
		return utils.ErrBadRequest
	}

	slog.Debug("Granting permission to role", "roleID", roleID, "permissionID", permissionID)

	// Start transaction
	tx, err := p.globalService.StartTransaction(ctx)
	if err != nil {
		return utils.ErrInternalServer
	}
	defer p.globalService.RollbackTransaction(ctx, tx)

	// Verificar se o role existe
	role, err := p.permissionRepository.GetRoleByID(ctx, tx, roleID)
	if err != nil {
		return utils.ErrInternalServer
	}
	if role == nil {
		return utils.ErrNotFound
	}

	// Verificar se a permissão existe
	permission, err := p.permissionRepository.GetPermissionByID(ctx, tx, permissionID)
	if err != nil {
		return utils.ErrInternalServer
	}
	if permission == nil {
		return utils.ErrNotFound
	}

	// Verificar se a relação já existe
	existingRolePermission, err := p.permissionRepository.GetRolePermissionByRoleIDAndPermissionID(ctx, tx, roleID, permissionID)
	if err == nil && existingRolePermission != nil {
		return utils.ErrConflict
	}

	// Criar a nova RolePermission
	rolePermission := permissionmodel.NewRolePermission()
	rolePermission.SetRoleID(roleID)
	rolePermission.SetPermissionID(permissionID)
	rolePermission.SetGranted(true)

	// Salvar no banco
	err = p.permissionRepository.CreateRolePermission(ctx, tx, rolePermission)
	if err != nil {
		return utils.ErrInternalServer
	}

	// Commit the transaction
	err = p.globalService.CommitTransaction(ctx, tx)
	if err != nil {
		return utils.ErrInternalServer
	}

	slog.Info("Permission granted to role successfully", "roleID", roleID, "permissionID", permissionID)
	return nil
}
