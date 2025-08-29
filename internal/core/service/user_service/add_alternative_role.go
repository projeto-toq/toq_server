package userservices

import (
	"context"
	"database/sql"

	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (us *userService) AddAlternativeRole(ctx context.Context, userID int64, role usermodel.UserRole, creciInfo ...string) (err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	tx, err := us.globalService.StartTransaction(ctx)
	if err != nil {
		return
	}

	_, err = us.addAlternativeRole(ctx, tx, userID, role, creciInfo...)
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

func (us *userService) addAlternativeRole(ctx context.Context, tx *sql.Tx, userID int64, role usermodel.UserRole, creciInfo ...string) (userRole usermodel.UserRoleInterface, err error) {

	//verify if the user is on active status
	user, err := us.repo.GetUserByID(ctx, tx, userID)
	if err != nil {
		return
	}

	if user.GetActiveRole().GetStatus() != usermodel.StatusActive {
		err = utils.ErrInternalServer
		return
	}

	if role == usermodel.RoleRealtor && len(creciInfo) != 3 {
		err = utils.ErrInternalServer
		return
	}

	baseRole, err := us.repo.GetBaseRoleByRole(ctx, tx, role)
	if err != nil {
		return
	}

	userRole = usermodel.NewUserRole()
	userRole.SetUserID(userID)
	userRole.SetBaseRoleID(baseRole.GetID())
	userRole.SetRole(baseRole.GetRole())
	userRole.SetActive(false)
	switch {
	case role == usermodel.RoleOwner:
		userRole.SetStatus(usermodel.StatusActive)
		userRole.SetStatusReason("")
	case role == usermodel.RoleRealtor:
		userRole.SetStatus(usermodel.StatusPendingImages)
		userRole.SetStatusReason("Awaiting creci images to verify")
		// CRECI functionality removed - keeping user role logic only
		// user.SetCreciNumber(creciInfo[0])
		// user.SetCreciState(creciInfo[1])
		// t, err1 := converters.StrngToTime(creciInfo[2])
		// if err1 != nil {
		// 	slog.Error("userservices.addAlternativeRole", "error converters.StrngToTime", err1)
		// 	err = utils.ErrInternalServer
		// 	return
		// }
		// user.SetCreciValidity(t)
	}

	err = us.repo.CreateUserRole(ctx, tx, userRole)
	if err != nil {
		return
	}

	if role == usermodel.RoleRealtor {
		err = us.repo.UpdateUserByID(ctx, tx, user)
		if err != nil {
			return
		}

		err = us.CreateUserFolder(ctx, user.GetID())
		if err != nil {
			return
		}
	}

	err = us.globalService.CreateAudit(ctx, tx, globalmodel.TableUserRoles, "Criado papel alternativo")
	if err != nil {
		return
	}

	return
}
