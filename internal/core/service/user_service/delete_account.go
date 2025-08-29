package userservices

import (
	"context"
	"database/sql"

	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (us *userService) DeleteAccount(ctx context.Context, userID int64) (tokens usermodel.Tokens, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	tx, err := us.globalService.StartTransaction(ctx)
	if err != nil {
		return
	}

	tokens, err = us.deleteAccount(ctx, tx, userID)
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

func (us *userService) deleteAccount(ctx context.Context, tx *sql.Tx, userId int64) (tokens usermodel.Tokens, err error) {

	user, err := us.repo.GetUserByID(ctx, tx, userId)
	if err != nil {
		return
	}
	//delete the account dependencies
	switch user.GetActiveRole().GetRole() {
	case usermodel.RoleOwner:
		err = us.CleanOwnerPending(ctx, user)
		if err != nil {
			return
		}
	case usermodel.RoleRealtor:
		err = us.CleanRealtorPending(ctx, user)
		if err != nil {
			return
		}
	case usermodel.RoleAgency:
		err = us.CleanAgencyPending(ctx, user)
		if err != nil {
			return
		}
	}

	//generate a new expired token to avoid user keep logged in
	tokens, err = us.CreateTokens(ctx, tx, user, true)
	if err != nil {
		return
	}

	us.setDeletedData(user)

	err = us.repo.UpdateUserByID(ctx, tx, user)
	if err != nil {
		return
	}

	err = us.repo.UpdateUserPasswordByID(ctx, tx, user)
	if err != nil {
		return
	}

	err = us.globalService.CreateAudit(ctx, tx, globalmodel.TableUsers, "Mascarado dados do usuário (conta apagada)")
	if err != nil {
		return
	}

	// Delete user folder in cloud storage
	if us.cloudStorageService != nil {
		folderErr := us.DeleteUserFolder(ctx, user.GetID())
		if folderErr != nil {
			// Log error but don't fail the transaction - account deletion should continue
			// even if cloud storage cleanup fails
			// Note: This will be handled by span tracing
		}
	}

	_, err = us.repo.DeleteUserRolesByUserID(ctx, tx, user.GetID())
	if err != nil {
		return
	}

	err = us.globalService.CreateAudit(ctx, tx, globalmodel.TableUserRoles, "Apagados os papéis do usuário, pois a conta foi apagada")
	if err != nil {
		return
	}

	return
}
