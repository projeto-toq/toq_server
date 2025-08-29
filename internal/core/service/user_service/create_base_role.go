package userservices

import (
	"context"
	"database/sql"

	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (us *userService) CreateBaseRole(ctx context.Context, role usermodel.UserRole, name string) (err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	tx, err := us.globalService.StartTransaction(ctx)
	if err != nil {
		return
	}

	err = us.createBaseRole(ctx, tx, role, name)
	if err != nil {
		us.globalService.RollbackTransaction(ctx, tx)
		return
	}

	err = us.globalService.CommitTransaction(ctx, tx)
	if err != nil {
		us.globalService.RollbackTransaction(ctx, tx)
		return
	}

	return
}

func (us *userService) createBaseRole(ctx context.Context, tx *sql.Tx, role usermodel.UserRole, name string) (err error) {

	baseRole := usermodel.NewBaseRole()
	baseRole.SetName(name)
	baseRole.SetRole(role)

	switch role {
	case usermodel.RoleRoot:
		// Admin privileges will be handled by permission system
	case usermodel.RoleOwner:
		baseRole.SetPrivileges(us.CreateOwnerPrivileges())
	case usermodel.RoleRealtor:
		baseRole.SetPrivileges(us.CreateRealtorPrivileges())
	case usermodel.RoleAgency:
		baseRole.SetPrivileges(us.CreateAgencyPrivileges())
	}

	us.repo.CreateBaseRole(ctx, tx, baseRole)

	err = us.globalService.CreateAudit(ctx, tx, globalmodel.TableBaseRoles, "Criado novo papel base")
	if err != nil {
		return
	}

	return

}
