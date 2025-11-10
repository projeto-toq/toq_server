package userservices

import (
	"context"
	"database/sql"
	"time"

	globalmodel "github.com/projeto-toq/toq_server/internal/core/model/global_model"
	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// AssignRoleOptions permite personalizar campos do user_role criado pelo serviço.
type AssignRoleOptions struct {
	IsActive *bool
	Status   *globalmodel.UserRoleStatus
}

// AssignRoleToUser atribui um role a um usuário (sem transação - uso direto)
func (us *userService) AssignRoleToUser(ctx context.Context, userID, roleID int64, expiresAt *time.Time, opts *AssignRoleOptions) (usermodel.UserRoleInterface, error) {
	ctx, end, _ := utils.GenerateTracer(ctx)
	defer end()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	// Start transaction
	tx, err := us.globalService.StartTransaction(ctx)
	if err != nil {
		logger.Error("permission.role.assign.tx_start_failed", "user_id", userID, "role_id", roleID, "error", err)
		utils.SetSpanError(ctx, err)
		return nil, utils.InternalError("")
	}
	defer func() {
		if err != nil {
			if rollbackErr := us.globalService.RollbackTransaction(ctx, tx); rollbackErr != nil {
				logger.Error("permission.role.assign.tx_rollback_failed", "user_id", userID, "role_id", roleID, "error", rollbackErr)
				utils.SetSpanError(ctx, rollbackErr)
			}
		}
	}()

	userRole, err := us.AssignRoleToUserWithTx(ctx, tx, userID, roleID, expiresAt, opts)
	if err != nil {
		return nil, err
	}

	// Commit the transaction
	err = us.globalService.CommitTransaction(ctx, tx)
	if err != nil {
		logger.Error("permission.role.assign.tx_commit_failed", "user_id", userID, "role_id", roleID, "error", err)
		utils.SetSpanError(ctx, err)
		return nil, utils.InternalError("")
	}

	return userRole, nil
}

// AssignRoleToUserWithTx atribui um role a um usuário (com transação - uso em fluxos)
//
// Esta função orquestra a atribuição de um role com as seguintes etapas:
//  1. Valida parâmetros de entrada (userID, roleID)
//  2. Verifica existência do role no banco
//  3. Verifica se o usuário já possui o role (evita duplicação)
//  4. Cria o registro UserRole com status e ativação configuráveis
//  5. Persiste no banco via repositório
//  6. Invalida cache de permissões do usuário (best-effort)
//
// A invalidação de cache ocorre APÓS a persistência e NÃO bloqueia a operação
// mesmo em caso de falha, pois o cache será eventualmente consistente.
//
// Parameters:
//   - ctx: Context for tracing, cancellation, and logging
//   - tx: Database transaction (must not be nil)
//   - userID: ID do usuário que receberá o role (must be > 0)
//   - roleID: ID do role a ser atribuído (must be > 0)
//   - expiresAt: Data de expiração do role (nil = sem expiração)
//   - opts: Opções para customizar is_active e status do UserRole (nil = defaults)
//
// Returns:
//   - userRole: UserRoleInterface com ID populado
//   - error: Domain error (400/404/409) ou infrastructure error (500)
//
// Business Rules:
//   - UserID e RoleID devem ser > 0
//   - Role deve existir no banco (404 se não encontrado)
//   - Usuário NÃO pode ter o mesmo role duplicado (409 se já existe)
//   - Status padrão: StatusPendingBoth (se não especificado em opts)
//   - IsActive padrão: true (se não especificado em opts)
//
// Side Effects:
//   - Cria registro em user_roles table
//   - Invalida cache de permissões do usuário (best-effort)
//   - Registra log Info com detalhes da atribuição
//
// Example:
//
//	opts := &AssignRoleOptions{
//	    IsActive: utils.BoolPtr(true),
//	    Status: utils.Ptr(globalmodel.StatusActive),
//	}
//	userRole, err := us.AssignRoleToUserWithTx(ctx, tx, userID, roleID, nil, opts)
//	if err != nil {
//	    return nil, err
//	}
func (us *userService) AssignRoleToUserWithTx(ctx context.Context, tx *sql.Tx, userID, roleID int64, expiresAt *time.Time, opts *AssignRoleOptions) (usermodel.UserRoleInterface, error) {
	// Initialize tracing for observability
	ctx, end, _ := utils.GenerateTracer(ctx)
	defer end()

	// Ensure logger propagation with request_id and trace_id
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	// Validate input parameters (business rules)
	if userID <= 0 {
		return nil, utils.BadRequest("invalid user id")
	}

	if roleID <= 0 {
		return nil, utils.BadRequest("invalid role id")
	}

	logger.Debug("permission.role.assign.request", "user_id", userID, "role_id", roleID, "expires_at", expiresAt)

	// Verify role exists (infrastructure + domain validation)
	role, err := us.permissionService.GetRoleByIDWithTx(ctx, tx, roleID)
	if err != nil {
		logger.Error("permission.role.assign.db_failed", "stage", "get_role", "user_id", userID, "role_id", roleID, "error", err)
		utils.SetSpanError(ctx, err)
		return nil, utils.InternalError("")
	}
	if role == nil {
		return nil, utils.NotFoundError("role")
	}

	// Check for duplicate role assignment (business rule: unique constraint)
	existingUserRole, err := us.repo.GetUserRoleByUserIDAndRoleID(ctx, tx, userID, roleID)
	if err != nil {
		logger.Error("permission.role.assign.db_failed", "stage", "get_user_role", "user_id", userID, "role_id", roleID, "error", err)
		utils.SetSpanError(ctx, err)
		return nil, utils.InternalError("")
	}
	if existingUserRole != nil {
		return nil, utils.ConflictError("role already assigned to user")
	}

	// Create new UserRole domain entity with defaults or custom options
	userRole := usermodel.NewUserRole()
	userRole.SetUserID(userID)
	userRole.SetRoleID(roleID)

	// Apply is_active option (default: true)
	isActive := true
	if opts != nil && opts.IsActive != nil {
		isActive = *opts.IsActive
	}
	userRole.SetIsActive(isActive)

	// Apply status option (default: StatusPendingBoth)
	status := globalmodel.StatusPendingBoth
	if opts != nil && opts.Status != nil {
		status = *opts.Status
	}
	userRole.SetStatus(status)

	// Set expiration if provided
	if expiresAt != nil {
		userRole.SetExpiresAt(expiresAt)
	}

	// Persist to database
	userRole, err = us.repo.CreateUserRole(ctx, tx, userRole)
	if err != nil {
		logger.Error("permission.role.assign.db_failed", "stage", "create_user_role", "user_id", userID, "role_id", roleID, "error", err)
		utils.SetSpanError(ctx, err)
		return nil, utils.InternalError("")
	}

	// Log success (domain event)
	logger.Info("permission.role.assigned", "user_id", userID, "role_id", roleID, "role_name", role.GetName(), "is_active", isActive, "status", status.String())

	// Invalidate user permissions cache (best-effort, post-commit operation)
	// Failure does not block the operation as the cache will be eventually consistent
	us.permissionService.InvalidateUserCacheSafe(ctx, userID, "assign_role")

	return userRole, nil
}
