package permissionservice

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	permissionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/permission_model"
)

// CreateRole cria um novo role no sistema
func (p *permissionServiceImpl) CreateRole(ctx context.Context, name, slug, description string, isSystemRole bool) (permissionmodel.RoleInterface, error) {
	if strings.TrimSpace(name) == "" {
		return nil, fmt.Errorf("role name cannot be empty")
	}

	if strings.TrimSpace(slug) == "" {
		return nil, fmt.Errorf("role slug cannot be empty")
	}

	slog.Debug("Creating role", "name", name, "slug", slug, "isSystemRole", isSystemRole)

	// Verificar se o slug já existe
	existingRole, err := p.permissionRepository.GetRoleBySlug(ctx, nil, slug)
	if err == nil && existingRole != nil {
		return nil, fmt.Errorf("role with slug '%s' already exists", slug)
	}

	// Criar o novo role
	newRole := permissionmodel.NewRole()
	newRole.SetName(name)
	newRole.SetSlug(slug)
	newRole.SetDescription(description)
	newRole.SetIsSystemRole(isSystemRole)
	newRole.SetIsActive(true)

	// Salvar no banco
	err = p.permissionRepository.CreateRole(ctx, nil, newRole)
	if err != nil {
		return nil, fmt.Errorf("failed to create role: %w", err)
	}

	slog.Info("Role created successfully", "name", name, "slug", slug)
	return newRole, nil
}
