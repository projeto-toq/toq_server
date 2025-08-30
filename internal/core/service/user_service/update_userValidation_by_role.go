package userservices

import (
	"context"
	"database/sql"

	permissionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/permission_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (us *userService) UpdateUserValidationByUserRole(ctx context.Context, tx *sql.Tx, user *usermodel.UserInterface, userValidation usermodel.ValidationInterface) (generateTokens bool, err error) {
	_, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	iUser := *user
	generateTokens = false

	//verify what validation was done
	switch {
	case userValidation.GetEmailCode() == "" && userValidation.GetPhoneCode() == "": //both email and phone are validated
		// Converter RoleInterface para RoleSlug
		roleSlug := permissionmodel.RoleSlug(iUser.GetActiveRole().GetRole().GetSlug())
		_, _, _, err = us.updateUserStatus(ctx, tx, roleSlug, usermodel.ActionFinishedEmailVerified)
		if err != nil {
			return false, err
		}

		// TODO: Implementar atualização de status via permission service
		// iUser.GetActiveRole().SetStatus(status)
		// iUser.GetActiveRole().SetStatusReason(reason)
		generateTokens = true

	case userValidation.GetEmailCode() == "" && userValidation.GetPhoneCode() != "": //phone is validated but email is not
		// iUser.GetActiveRole().SetStatusReason("Pending email validation")
		// Converter RoleInterface para RoleSlug
		roleSlug := permissionmodel.RoleSlug(iUser.GetActiveRole().GetRole().GetSlug())
		_, _, _, err = us.updateUserStatus(ctx, tx, roleSlug, usermodel.ActionFinishedPhoneVerified)
		if err != nil {
			return false, err
		}
		// TODO: Implementar atualização de status via permission service
		// iUser.GetActiveRole().SetStatus(status)
		// iUser.GetActiveRole().SetStatusReason(reason)
	}

	return
}
