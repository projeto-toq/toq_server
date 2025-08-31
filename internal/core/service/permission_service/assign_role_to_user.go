package permissionservice

import (
	"context"
	"database/sql"
	"log/slog"
	"time"

	permissionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/permission_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

// AssignRoleToUser atribui um role a um usuário (sem transação - uso direto)
func (p *permissionServiceImpl) AssignRoleToUser(ctx context.Context, userID, roleID int64, expiresAt *time.Time) error {
	// Start transaction
	tx, err := p.globalService.StartTransaction(ctx)
	if err != nil {
		return utils.ErrInternalServer
	}
	defer func() {
		if rollbackErr := p.globalService.RollbackTransaction(ctx, tx); rollbackErr != nil {
			slog.Error("Failed to rollback transaction", "error", rollbackErr)
		}
	}()

	err = p.AssignRoleToUserWithTx(ctx, tx, userID, roleID, expiresAt)
	if err != nil {
		return err
	}

	// Commit the transaction
	err = p.globalService.CommitTransaction(ctx, tx)
	if err != nil {
		return utils.ErrInternalServer
	}

	return nil
}

// AssignRoleToUserWithTx atribui um role a um usuário (com transação - uso em fluxos)
func (p *permissionServiceImpl) AssignRoleToUserWithTx(ctx context.Context, tx *sql.Tx, userID, roleID int64, expiresAt *time.Time) error {
	if userID <= 0 {
		return utils.ErrBadRequest
	}

	if roleID <= 0 {
		return utils.ErrBadRequest
	}

	slog.Debug("Assigning role to user", "userID", userID, "roleID", roleID, "expiresAt", expiresAt)

	// Verificar se o role existe
	role, err := p.permissionRepository.GetRoleByID(ctx, tx, roleID)
	if err != nil {
		return utils.ErrInternalServer
	}
	if role == nil {
		return utils.ErrNotFound
	}

	// Verificar se o usuário já tem este role
	existingUserRole, err := p.permissionRepository.GetUserRoleByUserIDAndRoleID(ctx, tx, userID, roleID)
	if err == nil && existingUserRole != nil {
		return utils.ErrConflict
	}

	// Criar o novo UserRole
	userRole := permissionmodel.NewUserRole()
	userRole.SetUserID(userID)
	userRole.SetRoleID(roleID)
	userRole.SetIsActive(true)

	if expiresAt != nil {
		userRole.SetExpiresAt(expiresAt)
	}

	// Salvar no banco
	err = p.permissionRepository.CreateUserRole(ctx, tx, userRole)
	if err != nil {
		return utils.ErrInternalServer
	}

	slog.Info("Role assigned to user successfully", "userID", userID, "roleID", roleID, "roleName", role.GetName())
	return nil
}
