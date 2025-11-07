# Plano de Refatoração — Reposicionamento da Gestão de `user_roles`

## 1. Diagnóstico
- **Arquivos analisados:**
  - `internal/core/model/user_model/*`
  - `internal/core/model/permission_model/*`
  - `internal/core/port/right/repository/{user_repository,permission_repository}/*`
  - `internal/adapter/right/mysql/{user,permission}/*`
  - `internal/core/service/{user_service,permission_service}/*`
  - `internal/adapter/left/http/handlers/user_handlers/*`
  - `internal/core/factory/*`
- **Desvios identificados:**
  1. **Adapters de permissionamento** (`internal/adapter/right/mysql/permission/*`): mantêm criação, deleção e leitura de `user_roles`, violando a Regra de Espelhamento Port ↔ Adapter (Seção 2.1) e gerando acoplamento indevido (Seção 14).
  2. **Serviço de permissionamento** (`permission_service`): centraliza regras de associação usuário-role (ex.: `AssignRoleToUserWithTx`, `GetUserRolesWithTx`), invertendo a dependência de domínio definida na Seção 1.
  3. **Modelos de domínio** (`user_model`, `permission_model`):
     - `UserInterface` não expõe métodos seguros para atualizar roles nem mantém slice inicializado, descumprindo o padrão da Seção 8.7.
     - `UserRoleInterface` em `permission_model` mantém responsabilidades de persistência e comentários incompletos, desalinhando o contrato com o domínio user.
  4. **Repositórios de usuário** (`user_repository`): `GetUserByID` não retorna role ativa; métodos relativos a roles residem no `permission_repository`, quebrando encapsulamento (Seção 7.3).
  5. **Serviço de usuário** (`user_service`): fluxos como `signin`, `create_*`, `delete_account`, `switch_user_role`, `get_user_by_id` realizam round-trips extras via permission service, aumentando risco transacional (Seções 5 e 7.1).
  6. **Handlers HTTP** (`user_handlers`): `GetUserRoles` depende diretamente do serviço de permissionamento, violando o fluxo Handler → Service → Repository (Seção 1).
- **Impacto:**
  - Latência maior e transações complexas para montar usuários completos.
  - Fronteiras de domínio violadas, dificultando evolução do user service.
  - Documentação inconsistente nos modelos (Seção 8).
- **Melhorias possíveis:**
  - Centralizar gestão de `user_roles` no domínio user.
  - Atualizar modelos para refletir invariantes e documentação adequadas.
  - Enxugar permission service para catálogo e cache.
  - Reduzir round-trips e garantir agregados consistentes.

## 2. Code Skeletons

### 2.1 Modelos (Domain)
```go
// internal/core/model/user_model/user_interface.go
package usermodel

// UserInterface represents a platform user aggregate with its active role.
//
// Responsibilities:
//   - Store identity information and current role assignment.
//   - Provide helpers to manage role collections while keeping invariants consistent.
//   - Remain decoupled from adapters/HTTP concerns (Guide Section 8.7).
type UserInterface interface {
    GetID() int64
    SetID(int64)
    GetFullName() string
    SetFullName(string)
    // … existing getters/setters …

    // GetRoles returns all roles attached to the aggregate (active + historical).
    GetRoles() []permissionmodel.UserRoleInterface
    // ReplaceRoles swaps the current role slice preserving immutability guarantees.
    ReplaceRoles([]permissionmodel.UserRoleInterface)
    // SetActiveRole sets the current active role for quick access.
    SetActiveRole(permissionmodel.UserRoleInterface)
    GetActiveRole() permissionmodel.UserRoleInterface
}

// NewUser creates a new user aggregate with initialized role slice.
func NewUser() UserInterface { /* initialize roles slice, set defaults */ }
```

```go
// internal/core/model/user_model/user_role_assignment.go
package usermodel

// UserRoleAssignment encapsulates the domain request to assign a role to a user.
//
// This DTO decouples service orchestration from persistence details.
type UserRoleAssignment struct {
    UserID       int64
    RoleID       int64
    Status       permissionmodel.UserRoleStatus
    IsActive     bool
    ExpiresAt    *time.Time
    BlockedUntil *time.Time
}
```

```go
// internal/core/model/permission_model/user_role_interface.go
package permissionmodel

// UserRoleInterface describes a user-role association enriched with role metadata.
//
// Updated to include optional blocking and lifecycle timestamps (Guide Section 8.7).
type UserRoleInterface interface {
    GetID() int64
    SetID(int64)
    GetUserID() int64
    SetUserID(int64)
    GetRoleID() int64
    SetRoleID(int64)
    GetIsActive() bool
    SetIsActive(bool)
    GetStatus() UserRoleStatus
    SetStatus(UserRoleStatus)
    GetExpiresAt() *time.Time
    SetExpiresAt(*time.Time)
    GetBlockedUntil() *time.Time
    SetBlockedUntil(*time.Time)
    GetRole() RoleInterface
    SetRole(RoleInterface)
}
```

### 2.2 Portas (Repositories)
```go
// internal/core/port/right/repository/user_repository/user_repository_interface.go
package userrepository

// UserRepoPortInterface centralizes persistence for user aggregates (Guide Section 12).
type UserRepoPortInterface interface {
    // Existing signatures …

    // GetUserWithActiveRole loads user and current active role in a single query.
    GetUserWithActiveRole(ctx context.Context, tx *sql.Tx, id int64) (usermodel.UserInterface, error)

    // ListUserRoles returns all role assignments (active + inactive).
    ListUserRoles(ctx context.Context, tx *sql.Tx, userID int64) ([]permissionmodel.UserRoleInterface, error)

    // AssignUserRole persists a new role association.
    AssignUserRole(ctx context.Context, tx *sql.Tx, assignment usermodel.UserRoleAssignment) (permissionmodel.UserRoleInterface, error)

    // UpdateUserRoleState updates status/is_active for a given role slug.
    UpdateUserRoleState(ctx context.Context, tx *sql.Tx, userID int64, roleSlug permissionmodel.RoleSlug, status permissionmodel.UserRoleStatus, isActive bool) error

    // DeactivateAllUserRoles deactivates every role for consistency checks.
    DeactivateAllUserRoles(ctx context.Context, tx *sql.Tx, userID int64) error

    // RemoveUserRole deletes a specific role assignment.
    RemoveUserRole(ctx context.Context, tx *sql.Tx, userRoleID int64) error
}
```

```go
// internal/core/port/right/repository/permission_repository/permission_repository_interface.go
package permissionrepository

// PermissionRepositoryInterface keeps only catalog operations and complex queries for permissions.
type PermissionRepositoryInterface interface {
    // Role catalog CRUD …
    // Permission catalog CRUD …
    // Role-permission bindings …

    // GetUserPermissions is retained for permission checks.
    GetUserPermissions(ctx context.Context, tx *sql.Tx, userID int64) ([]permissionmodel.PermissionInterface, error)

    // User-role mutating operations removed (managed by user repository).
}
```

### 2.3 Adapters (MySQL)
```go
// internal/adapter/right/mysql/user/get_user_with_active_role.go
package mysqluseradapter

// GetUserWithActiveRole loads the user and active role using explicit column selection.
func (ua *UserAdapter) GetUserWithActiveRole(ctx context.Context, tx *sql.Tx, id int64) (usermodel.UserInterface, error) {
    ctx, spanEnd, err := utils.GenerateTracer(ctx)
    if err != nil {
        return nil, err
    }
    defer spanEnd()

    ctx = utils.ContextWithLogger(ctx)

    const query = `SELECT /* explicit column list across users, user_roles, roles */ ...`

    row := ua.QueryRowContext(ctx, tx, "select", query, id)
    userEntity, roleEntity, scanErr := scanUserWithActiveRole(row)
    if scanErr != nil {
        utils.SetSpanError(ctx, scanErr)
        return nil, fmt.Errorf("scan user with active role: %w", scanErr)
    }

    user := userconverters.UserEntityToDomain(userEntity)
    if roleEntity != nil {
        userRole := permissionconverters.UserRoleEntityToDomain(*roleEntity)
        user.SetActiveRole(userRole)
        user.ReplaceRoles([]permissionmodel.UserRoleInterface{userRole})
    }

    return user, nil
}
```

```go
// internal/adapter/right/mysql/user/assign_user_role.go
package mysqluseradapter

// AssignUserRole inserts a new user_role row using InstrumentedAdapter for observability.
func (ua *UserAdapter) AssignUserRole(ctx context.Context, tx *sql.Tx, assignment usermodel.UserRoleAssignment) (permissionmodel.UserRoleInterface, error) {
    ctx, spanEnd, err := utils.GenerateTracer(ctx)
    if err != nil {
        return nil, err
    }
    defer spanEnd()

    ctx = utils.ContextWithLogger(ctx)

    const query = `INSERT INTO user_roles (user_id, role_id, is_active, status, expires_at, blocked_until) VALUES (?, ?, ?, ?, ?, ?)`

    result, execErr := ua.ExecContext(ctx, tx, "insert", query,
        assignment.UserID,
        assignment.RoleID,
        assignment.IsActive,
        int(assignment.Status),
        assignment.ExpiresAt,
        assignment.BlockedUntil,
    )
    if execErr != nil {
        utils.SetSpanError(ctx, execErr)
        return nil, fmt.Errorf("insert user role: %w", execErr)
    }

    id, lastErr := result.LastInsertId()
    if lastErr != nil {
        utils.SetSpanError(ctx, lastErr)
        return nil, fmt.Errorf("user role last insert id: %w", lastErr)
    }

    role := permissionmodel.NewUserRole()
    role.SetID(id)
    role.SetUserID(assignment.UserID)
    role.SetRoleID(assignment.RoleID)
    role.SetStatus(assignment.Status)
    role.SetIsActive(assignment.IsActive)
    role.SetExpiresAt(assignment.ExpiresAt)
    role.SetBlockedUntil(assignment.BlockedUntil)
    return role, nil
}
```

### 2.4 Serviço de Usuário
```go
// internal/core/service/user_service/user_roles.go
package userservices

// AssignRoleToUser orchestrates the assignment ensuring cache invalidation and invariants.
func (us *userService) AssignRoleToUser(ctx context.Context, tx *sql.Tx, user usermodel.UserInterface, role permissionmodel.RoleInterface, opts AssignUserRoleOptions) (UserRoleAssignmentResult, error) {
    ctx, spanEnd, err := utils.GenerateTracer(ctx)
    if err != nil {
        return UserRoleAssignmentResult{}, utils.InternalError("Failed to generate tracer")
    }
    defer spanEnd()

    ctx = utils.ContextWithLogger(ctx)

    assignment := usermodel.UserRoleAssignment{
        UserID:       user.GetID(),
        RoleID:       role.GetID(),
        Status:       opts.Status,
        IsActive:     opts.IsActive,
        ExpiresAt:    opts.ExpiresAt,
        BlockedUntil: opts.BlockedUntil,
    }

    storedRole, err := us.repo.AssignUserRole(ctx, tx, assignment)
    if err != nil {
        utils.SetSpanError(ctx, err)
        return UserRoleAssignmentResult{}, utils.InternalError("Failed to assign role")
    }

    if assignment.IsActive {
        if err := us.repo.DeactivateAllUserRoles(ctx, tx, user.GetID()); err != nil {
            utils.SetSpanError(ctx, err)
            return UserRoleAssignmentResult{}, utils.InternalError("Failed to deactivate previous roles")
        }
        storedRole.SetIsActive(true)
        user.SetActiveRole(storedRole)
    }

    user.ReplaceRoles(append(user.GetRoles(), storedRole))

    if err := us.permissionService.InvalidateUserCache(ctx, user.GetID()); err != nil {
        utils.SetSpanError(ctx, err)
        return UserRoleAssignmentResult{}, utils.InternalError("Failed to invalidate permission cache")
    }

    return UserRoleAssignmentResult{Assignment: storedRole, User: user}, nil
}
```

### 2.5 Serviço de Permissionamento
```go
// internal/core/service/permission_service/permission_service.go
package permissionservice

// PermissionServiceInterface keeps catalog operations and cache helpers only.
type PermissionServiceInterface interface {
    HasHTTPPermission(ctx context.Context, userID int64, method, path string) (bool, error)
    ListRoles(ctx context.Context, input ListRolesInput) (ListRolesOutput, error)
    GetRoleBySlug(ctx context.Context, slug permissionmodel.RoleSlug) (permissionmodel.RoleInterface, error)
    // … demais operações de catálogo …

    InvalidateUserCache(ctx context.Context, userID int64) error
    ClearUserPermissionsCache(ctx context.Context, userID int64) error

    BlockUserTemporarily(ctx context.Context, tx *sql.Tx, userID int64, reason string) error
    UnblockUser(ctx context.Context, tx *sql.Tx, userID int64) error
    IsUserTempBlockedWithTx(ctx context.Context, tx *sql.Tx, userID int64) (bool, error)
}
```

### 2.6 Handlers HTTP
```go
// internal/adapter/left/http/handlers/user_handlers/get_user_roles.go
package userhandlers

// GetUserRoles lists every role assignment of the authenticated user.
// @Summary     List user roles
// @Description Returns all role assignments (active and inactive) linked to the authenticated user.
// @Tags        Users
// @Produce     json
// @Security    BearerAuth
// @Success     200 {object} dto.GetUserRolesResponse
// @Failure     401 {object} dto.ErrorResponse
// @Failure     500 {object} dto.ErrorResponse
// @Router      /users/me/roles [get]
func (uh *UserHandler) GetUserRoles(c *gin.Context) {
    ctx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)
    info, ok := middlewares.GetUserInfoFromContext(c)
    if !ok {
        httperrors.SendHTTPError(c, http.StatusUnauthorized, "AUTH_CONTEXT_MISSING", "User context not found")
        return
    }

    roles, err := uh.userService.ListUserRoles(ctx, info.ID)
    if err != nil {
        httperrors.SendHTTPErrorObj(c, err)
        return
    }

    response := dto.NewGetUserRolesResponse(roles)
    c.JSON(http.StatusOK, response)
}
```

## 3. Estrutura de Diretórios (final)
```
internal/core/model/user_model/
  user_interface.go        (atualizado)
  user_domain.go           (atualizado)
  user_role_assignment.go  (novo)
  …
internal/core/model/permission_model/
  user_role_interface.go   (atualizado)
  user_role_domain.go      (atualizado)
  …

internal/core/port/right/repository/
  user_repository/
    user_repository_interface.go (atualizado)
  permission_repository/
    permission_repository_interface.go (atualizado)

internal/adapter/right/mysql/user/
  user_adapter.go
  get_user_by_id.go
  get_user_with_active_role.go   (novo)
  assign_user_role.go            (novo)
  list_user_roles.go             (novo)
  update_user_role_state.go      (atualizado)
  deactivate_all_user_roles.go   (atualizado)
  converters/
    user_entity_to_domain.go
    user_role_entity_to_domain.go (migrado do permission)
    user_role_domain_to_entity.go (novo)
  entities/
    user_entity.go
    user_role_entity.go          (migrado do permission)
    user_role_with_role_entity.go (novo)

internal/adapter/right/mysql/permission/
  (mantém apenas operações de roles/permissions — arquivos de user_roles removidos)

internal/core/service/user_service/
  user_service.go          (ajuste de dependências)
  user_roles.go            (novo)
  get_user_by_id.go        (atualizado)
  signin.go, create_owner.go, … (atualizados para novo helper)
  models.go                (novo helpers internos, se necessário)

internal/core/service/permission_service/
  permission_service.go    (interface enxuta)
  assign_role_to_user.go, get_user_roles.go, … (removidos ou reduzidos)

internal/adapter/left/http/handlers/user_handlers/
  get_user_roles.go        (atualizado)

docs/refactors/
  user_roles_refactor_plan.md (este plano)
```

## 4. Ordem de Execução
  1. **Modelos:** atualizar `user_model` e `permission_model` com novos métodos, documentação e structs auxiliares.
  2. **Portas:** ajustar interfaces de `user_repository` e `permission_repository` refletindo responsabilidades.
  3. **Adapters:** migrar entidades/converters de `user_roles` para o adapter de user, criar queries combinadas e remover arquivos equivalentes do adapter de permission.
4. **Serviço de Usuário:** implementar helpers (`AssignRoleToUser`, `ListUserRoles`, `hydrateActiveRole`) e atualizar chamadas existentes.
5. **Serviço de Permissionamento:** remover métodos de gestão de `user_roles`, manter apenas catálogo/cache, adaptar usos remanescentes.
6. **Handlers/DTOs:** direcionar `GetUserRoles` para o serviço de usuário, revisar DTOs se necessário.
7. **Injeção de dependências:** atualizar factories e construções de serviços para refletir novas interfaces.
8. **Observabilidade:** garantir que novos fluxos mantenham tracing/logging conforme Seções 5, 7 e 9; revisar invalidação de cache.
9. **Documentação:** revisar Godoc/Swagger conforme Seção 8 e atualizar checklist.
10. **Validação manual:** rodar `go vet`/lint localmente, validar queries manualmente e preparar PR com riscos conhecidos.

## 5. Checklist de Conformidade
- [ ] Arquitetura hexagonal (Seção 1)
- [ ] Regra de Espelhamento Port ↔ Adapter (Seção 2.1)
- [ ] InstrumentedAdapter em repos (Seção 7.3)
- [ ] Transações via globalService (Seção 7.1)
- [ ] Tracing/Logging/Erros (Seções 5, 7, 9)
- [ ] Documentação (Seção 8) — especialmente modelos atualizados
- [ ] Sem anti-padrões (Seção 14)

## Observações Finais
- Seguir restrições: sem alterações em testes, scripts SQL ou swagger.json/yaml.
- Executar `make swagger` somente após ajustar comentários em handlers.
- Dividir PRs por etapas (modelos + portas, adapters, serviços/handlers) para facilitar revisão.
