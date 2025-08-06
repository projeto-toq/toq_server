package userservices

import (
	"context"
	"database/sql"

	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
)

func (us *userService) AddFirstUserRole(ctx context.Context, tx *sql.Tx, user usermodel.UserInterface, role usermodel.UserRole) (err error) {

	baseRole, err := us.repo.GetBaseRoleByRole(ctx, tx, role)
	if err != nil {
		return
	}

	iRole := usermodel.NewUserRole()
	iRole.SetRole(baseRole.GetRole())
	iRole.SetActive(true)
	iRole.SetBaseRoleID(baseRole.GetID())
	iRole.SetUserID(user.GetID())

	status, reason, _, err := us.updateUserStatus(ctx, tx, iRole.GetRole(), usermodel.ActionFinishedCreated)
	if err != nil {
		return
	}
	iRole.SetStatus(status)
	iRole.SetStatusReason(reason)

	err = us.repo.CreateUserRole(ctx, tx, iRole)
	if err != nil {
		return
	}

	user.SetActiveRole(iRole)

	err = us.globalService.CreateAudit(ctx, tx, globalmodel.TableUserRoles, "Criado papel inicial")
	if err != nil {
		return
	}

	return
}
