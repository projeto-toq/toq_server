package userservices

import (
	"context"
	"database/sql"

	auditmodel "github.com/projeto-toq/toq_server/internal/core/model/audit_model"
	globalmodel "github.com/projeto-toq/toq_server/internal/core/model/global_model"
	permissionmodel "github.com/projeto-toq/toq_server/internal/core/model/permission_model"
	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"
	auditservice "github.com/projeto-toq/toq_server/internal/core/service/audit_service"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// CreateRealtor cria a conta de corretor e autentica via SignIn padrão
func (us *userService) CreateRealtor(ctx context.Context, realtor usermodel.UserInterface, plainPassword string, deviceToken string, ipAddress string, userAgent string) (tokens usermodel.Tokens, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return tokens, utils.InternalError("Failed to generate tracer")
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	rawDeviceID, _ := ctx.Value(globalmodel.DeviceIDKey).(string)
	ctx, trimmedDeviceToken, trimmedDeviceID, derr := us.sanitizeDeviceContext(ctx, deviceToken, rawDeviceID, "user.create_realtor")
	if derr != nil {
		return tokens, derr
	}

	tx, err := us.globalService.StartTransaction(ctx)
	if err != nil {
		logger.Error("user.create_realtor.tx_start_error", "err", err)
		utils.SetSpanError(ctx, err)
		return tokens, utils.InternalError("Failed to start transaction")
	}
	defer func() {
		if err != nil {
			if rbErr := us.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				logger.Error("user.create_realtor.tx_rollback_error", "err", rbErr)
				utils.SetSpanError(ctx, rbErr)
			}
		}
	}()

	created, err := us.createRealtor(ctx, tx, realtor)
	if err != nil {
		return tokens, err
	}

	tokens, err = us.signIn(ctx, tx, created.GetNationalID(), plainPassword, trimmedDeviceToken, trimmedDeviceID, ipAddress, userAgent)
	if err != nil {
		return tokens, err
	}

	err = us.globalService.CommitTransaction(ctx, tx)
	if err != nil {
		logger.Error("user.create_realtor.tx_commit_error", "err", err)
		utils.SetSpanError(ctx, err)
		return tokens, utils.InternalError("Failed to commit transaction")
	}

	return tokens, nil
}

// createRealtor cria o usuário corretor e retorna o usuário criado
func (us *userService) createRealtor(ctx context.Context, tx *sql.Tx, realtor usermodel.UserInterface) (created usermodel.UserInterface, err error) {
	ctx = utils.ContextWithLogger(ctx)

	// Usar permission service diretamente para buscar role
	role, err := us.permissionService.GetRoleBySlugWithTx(ctx, tx, permissionmodel.RoleSlugRealtor)
	if err != nil {
		return
	}

	err = us.ValidateUserData(ctx, tx, realtor, permissionmodel.RoleSlugRealtor)
	if err != nil {
		return
	}

	err = us.repo.CreateUser(ctx, tx, realtor)
	if err != nil {
		return nil, err
	}

	// Usar permission service diretamente para atribuir role
	userRole, err := us.AssignRoleToUserWithTx(ctx, tx, realtor.GetID(), role.GetID(), nil, nil)
	if err != nil {
		return
	}

	realtor.SetActiveRole(userRole)

	err = us.CreateUserValidations(ctx, tx, realtor)
	if err != nil {
		return
	}

	err = us.cloudStorageService.CreateUserFolder(ctx, realtor.GetID())
	if err != nil {
		return
	}

	auditRecord := auditservice.BuildRecordFromContext(
		ctx,
		realtor.GetID(),
		auditmodel.AuditTarget{Type: auditmodel.TargetUser, ID: realtor.GetID()},
		auditmodel.OperationCreate,
		map[string]any{"role_slug": string(permissionmodel.RoleSlugRealtor), "origin": "realtor_signup"},
	)
	if err = us.auditService.RecordChange(ctx, tx, auditRecord); err != nil {
		return nil, err
	}
	return realtor, nil
}
