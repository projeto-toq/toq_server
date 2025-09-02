package userservices

import (
	"context"
	"database/sql"

	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	permissionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/permission_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

// CreateAgency cria a conta de imobiliária e autentica via SignIn padrão
func (us *userService) CreateAgency(ctx context.Context, agency usermodel.UserInterface, plainPassword string, deviceToken string, ipAddress string, userAgent string) (tokens usermodel.Tokens, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	tx, err := us.globalService.StartTransaction(ctx)
	if err != nil {
		return
	}

	created, err := us.createAgency(ctx, tx, agency)
	if err != nil {
		us.globalService.RollbackTransaction(ctx, tx)
		return
	}

	err = us.globalService.CommitTransaction(ctx, tx)
	if err != nil {
		us.globalService.RollbackTransaction(ctx, tx)
		return
	}

	// Autentica após commit
	tokens, err = us.SignInWithContext(ctx, created.GetNationalID(), plainPassword, deviceToken, ipAddress, userAgent)
	return tokens, err
}

// createAgency cria o usuário imobiliária e retorna o usuário criado
func (us *userService) createAgency(ctx context.Context, tx *sql.Tx, agency usermodel.UserInterface) (created usermodel.UserInterface, err error) {

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
	userRole, err := us.permissionService.AssignRoleToUserWithTx(ctx, tx, agency.GetID(), role.GetID(), nil)
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
