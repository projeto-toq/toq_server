package userservices

import (
	"context"
	"database/sql"

	globalmodel "github.com/projeto-toq/toq_server/internal/core/model/global_model"
	permissionmodel "github.com/projeto-toq/toq_server/internal/core/model/permission_model"
	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// CreateAgency cria a conta de imobiliária e autentica via SignIn padrão
func (us *userService) CreateAgency(ctx context.Context, agency usermodel.UserInterface, plainPassword string, deviceToken string, ipAddress string, userAgent string) (tokens usermodel.Tokens, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return tokens, utils.InternalError("Failed to generate tracer")
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	tx, err := us.globalService.StartTransaction(ctx)
	if err != nil {
		logger.Error("user.create_agency.tx_start_error", "err", err)
		utils.SetSpanError(ctx, err)
		return tokens, utils.InternalError("Failed to start transaction")
	}
	defer func() {
		if err != nil {
			if rbErr := us.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				logger.Error("user.create_agency.tx_rollback_error", "err", rbErr)
				utils.SetSpanError(ctx, rbErr)
			}
		}
	}()

	created, err := us.createAgency(ctx, tx, agency)
	if err != nil {
		return tokens, err
	}

	err = us.globalService.CommitTransaction(ctx, tx)
	if err != nil {
		logger.Error("user.create_agency.tx_commit_error", "err", err)
		utils.SetSpanError(ctx, err)
		return tokens, utils.InternalError("Failed to commit transaction")
	}

	// Autentica após commit
	deviceID, _ := ctx.Value(globalmodel.DeviceIDKey).(string)
	tokens, err = us.SignInWithContext(ctx, created.GetNationalID(), plainPassword, deviceToken, deviceID, ipAddress, userAgent)
	return tokens, err
}

// createAgency cria o usuário imobiliária e retorna o usuário criado
func (us *userService) createAgency(ctx context.Context, tx *sql.Tx, agency usermodel.UserInterface) (created usermodel.UserInterface, err error) {
	ctx = utils.ContextWithLogger(ctx)

	// Usar permission service com constante tipada
	role, err := us.permissionService.GetRoleBySlugWithTx(ctx, tx, permissionmodel.RoleSlugAgency)
	if err != nil {
		return
	}

	err = us.ValidateUserData(ctx, tx, agency, permissionmodel.RoleSlugAgency)
	if err != nil {
		return
	}

	err = us.repo.CreateUser(ctx, tx, agency)
	if err != nil {
		return nil, err
	}

	// Usar permission service diretamente para atribuir role
	userRole, err := us.AssignRoleToUserWithTx(ctx, tx, agency.GetID(), role.GetID(), nil, nil)
	if err != nil {
		return
	}

	agency.SetActiveRole(userRole)

	err = us.CreateUserValidations(ctx, tx, agency)
	if err != nil {
		return
	}

	err = us.cloudStorageService.CreateUserFolder(ctx, agency.GetID())
	if err != nil {
		return
	}

	err = us.globalService.CreateAudit(ctx, tx, globalmodel.TableUsers, "Criado novo usuário com papel de Imobiliária", agency.GetID())
	if err != nil {
		return nil, err
	}
	return agency, nil
}
