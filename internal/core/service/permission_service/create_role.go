package permissionservice

import (
	"context"
	"log/slog"
	"strings"

	permissionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/permission_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

// CreateRole cria um novo role no sistema
func (p *permissionServiceImpl) CreateRole(ctx context.Context, name string, slug permissionmodel.RoleSlug, description string, isSystemRole bool) (permissionmodel.RoleInterface, error) {
	// Tracing da operação
	ctx, end, _ := utils.GenerateTracer(ctx)
	defer end()

	if strings.TrimSpace(name) == "" {
		return nil, utils.ValidationError("name", "role name cannot be empty")
	}

	if !slug.IsValid() {
		return nil, utils.ValidationError("slug", "invalid role slug")
	}

	slog.Debug("permission.role.create.start", "name", name, "slug", slug, "is_system_role", isSystemRole)

	// Verificar se o slug já existe (read-only tx)
	tx, err := p.globalService.StartReadOnlyTransaction(ctx)
	if err != nil {
		slog.Error("permission.role.tx_ro_start_failed", "slug", slug, "error", err)
		return nil, utils.InternalError("")
	}
	existingRole, getErr := p.permissionRepository.GetRoleBySlug(ctx, tx, slug.String())
	if getErr != nil {
		slog.Error("permission.role.get_by_slug_failed", "slug", slug, "error", getErr)
		_ = p.globalService.RollbackTransaction(ctx, tx)
		return nil, utils.InternalError("")
	}
	if cerr := p.globalService.CommitTransaction(ctx, tx); cerr != nil {
		slog.Error("permission.role.tx_ro_commit_failed", "slug", slug, "error", cerr)
		return nil, utils.InternalError("")
	}
	if existingRole != nil {
		return nil, utils.ConflictError("role slug already exists")
	}

	// Criar o novo role
	newRole := permissionmodel.NewRole()
	newRole.SetName(name)
	newRole.SetSlug(slug.String())
	newRole.SetDescription(description)
	newRole.SetIsSystemRole(isSystemRole)
	newRole.SetIsActive(true)

	// Salvar no banco (tx de escrita)
	wtx, werr := p.globalService.StartTransaction(ctx)
	if werr != nil {
		slog.Error("permission.role.tx_start_failed", "slug", slug, "error", werr)
		return nil, utils.InternalError("")
	}
	if err = p.permissionRepository.CreateRole(ctx, wtx, newRole); err != nil {
		slog.Error("permission.role.create_failed", "slug", slug, "error", err)
		return nil, utils.InternalError("")
	}
	if cerr := p.globalService.CommitTransaction(ctx, wtx); cerr != nil {
		slog.Error("permission.role.tx_commit_failed", "slug", slug, "error", cerr)
		return nil, utils.InternalError("")
	}

	slog.Info("permission.role.created", "role_id", newRole.GetID(), "slug", slug, "is_system_role", isSystemRole)
	return newRole, nil
}
