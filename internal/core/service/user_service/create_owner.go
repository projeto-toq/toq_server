package userservices

import (
	"context"
	"database/sql"

	globalmodel "github.com/projeto-toq/toq_server/internal/core/model/global_model"
	permissionmodel "github.com/projeto-toq/toq_server/internal/core/model/permission_model"
	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// CreateOwner cria a conta e autentica o usuário via fluxo padrão de SignIn
func (us *userService) CreateOwner(ctx context.Context, owner usermodel.UserInterface, plainPassword string, deviceToken string, ipAddress string, userAgent string) (tokens usermodel.Tokens, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return tokens, utils.InternalError("Failed to generate tracer")
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	tx, err := us.globalService.StartTransaction(ctx)
	if err != nil {
		logger.Error("user.create_owner.tx_start_error", "err", err)
		utils.SetSpanError(ctx, err)
		return tokens, utils.InternalError("Failed to start transaction")
	}
	defer func() {
		if err != nil {
			if rbErr := us.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				logger.Error("user.create_owner.tx_rollback_error", "err", rbErr)
				utils.SetSpanError(ctx, rbErr)
			}
		}
	}()

	// Executa criação e audit em transação; retorna usuário criado
	created, err := us.createOwner(ctx, tx, owner)
	if err != nil {
		return tokens, err
	}

	err = us.globalService.CommitTransaction(ctx, tx)
	if err != nil {
		logger.Error("user.create_owner.tx_commit_error", "err", err)
		utils.SetSpanError(ctx, err)
		return tokens, utils.InternalError("Failed to commit transaction")
	}

	// Após commit, autentica usando SignIn padrão com nationalID normalizado
	deviceID, _ := ctx.Value(globalmodel.DeviceIDKey).(string)
	tokens, err = us.SignInWithContext(ctx, created.GetNationalID(), plainPassword, deviceToken, deviceID, ipAddress, userAgent)
	return tokens, err
}

// createOwner executa a criação transacional do usuário Owner e retorna o usuário criado
func (us *userService) createOwner(ctx context.Context, tx *sql.Tx, owner usermodel.UserInterface) (created usermodel.UserInterface, err error) {
	ctx = utils.ContextWithLogger(ctx)

	// Usar permission service diretamente para buscar role
	role, err := us.permissionService.GetRoleBySlugWithTx(ctx, tx, permissionmodel.RoleSlugOwner)
	if err != nil {
		return
	}

	err = us.ValidateUserData(ctx, tx, owner, permissionmodel.RoleSlugOwner)
	if err != nil {
		return
	}

	err = us.repo.CreateUser(ctx, tx, owner)
	if err != nil {
		return nil, err
	}

	userRole, err := us.permissionService.AssignRoleToUserWithTx(ctx, tx, owner.GetID(), role.GetID(), nil, nil)
	if err != nil {
		return
	}

	owner.SetActiveRole(userRole)

	err = us.CreateUserValidations(ctx, tx, owner)
	if err != nil {
		return
	}

	err = us.cloudStorageService.CreateUserFolder(ctx, owner.GetID())
	if err != nil {
		return
	}

	err = us.globalService.CreateAudit(ctx, tx, globalmodel.TableUsers, "Criado novo usuário com papel de Proprietário", owner.GetID())
	if err != nil {
		return nil, err
	}

	return owner, nil
}
