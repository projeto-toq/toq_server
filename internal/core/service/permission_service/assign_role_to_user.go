package permissionservice

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	permissionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/permission_model"
)

// AssignRoleToUser atribui um role a um usuário
func (p *permissionServiceImpl) AssignRoleToUser(ctx context.Context, userID, roleID int64, expiresAt *time.Time) error {
	if userID <= 0 {
		return fmt.Errorf("invalid user ID: %d", userID)
	}

	if roleID <= 0 {
		return fmt.Errorf("invalid role ID: %d", roleID)
	}

	slog.Debug("Assigning role to user", "userID", userID, "roleID", roleID, "expiresAt", expiresAt)

	// Verificar se o role existe
	role, err := p.permissionRepository.GetRoleByID(ctx, nil, roleID)
	if err != nil {
		return fmt.Errorf("failed to verify role existence: %w", err)
	}
	if role == nil {
		return fmt.Errorf("role with ID %d does not exist", roleID)
	}

	// Verificar se o usuário já tem este role
	existingUserRole, err := p.permissionRepository.GetUserRoleByUserIDAndRoleID(ctx, nil, userID, roleID)
	if err == nil && existingUserRole != nil {
		return fmt.Errorf("user %d already has role %d", userID, roleID)
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
	err = p.permissionRepository.CreateUserRole(ctx, nil, userRole)
	if err != nil {
		return fmt.Errorf("failed to assign role to user: %w", err)
	}

	slog.Info("Role assigned to user successfully", "userID", userID, "roleID", roleID, "roleName", role.GetName())
	return nil
}
