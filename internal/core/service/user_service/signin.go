package userservices

import (
	"context"
	"database/sql"
	"log/slog"
	"time"

	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"

	"errors"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (us *userService) SignIn(ctx context.Context, nationalID string, password string, deviceToken string) (tokens usermodel.Tokens, err error) {

	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	tx, err := us.globalService.StartTransaction(ctx)
	if err != nil {
		return
	}

	tokens, err = us.signIn(ctx, tx, nationalID, password, deviceToken)
	if err != nil {
		us.globalService.RollbackTransaction(ctx, tx)
		return
	}

	// Commit the transaction
	err = us.globalService.CommitTransaction(ctx, tx)
	if err != nil {
		us.globalService.RollbackTransaction(ctx, tx)
		return
	}

	return
}

func (us *userService) signIn(ctx context.Context, tx *sql.Tx, nationalID string, password string, deviceToken string) (tokens usermodel.Tokens, err error) {
	criptoPassword := us.encryptPassword(password)
	user, err := us.repo.GetUserByNationalID(ctx, tx, nationalID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = utils.ErrInternalServer
			return
		}
		return
	}

	// Get active roles via Permission Service
	userRoles, err := us.permissionService.GetUserRolesWithTx(ctx, tx, user.GetID())
	if err != nil {
		return
	}

	if len(userRoles) == 0 {
		err = utils.ErrInternalServer
		return
	}

	// Convert permission model to user model (get the first active role)
	// TODO: Implement proper active role selection logic based on business rules
	activeRole := usermodel.NewUserRole()
	activeRole.SetID(userRoles[0].GetID())
	activeRole.SetUserID(userRoles[0].GetUserID())
	activeRole.SetBaseRoleID(userRoles[0].GetRoleID())
	activeRole.SetActive(userRoles[0].GetIsActive())

	user.SetActiveRole(activeRole)

	if user.GetActiveRole().GetStatus() == usermodel.StatusBlocked {
		err = utils.ErrInternalServer
		return
	}

	//compare the password with cripto password
	if user.GetPassword() != criptoPassword {
		err = checkWrongSignin(ctx, tx, us, user)
		if err != nil {
			return
		}

		err = utils.ErrInternalServer
		return
	}

	_, err = us.repo.DeleteWrongSignInByUserID(ctx, tx, user.GetID())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return
		}
	}

	// Attach device token if provided (add or ignore if duplicate)
	if deviceToken != "" {
		if errAdd := us.repo.AddDeviceToken(ctx, tx, user.GetID(), deviceToken, nil); errAdd != nil {
			slog.Warn("signIn: failed to add device token", "userID", user.GetID(), "err", errAdd)
		}
	}

	//generate the tokens
	tokens, err = us.CreateTokens(ctx, tx, user, false)
	if err != nil {
		return
	}

	// Note: Last activity is now tracked automatically by AuthInterceptor → Redis → Batch worker
	// No need for direct UpdateUserLastActivity call

	return
}

func checkWrongSignin(ctx context.Context, tx *sql.Tx, us *userService, user usermodel.UserInterface) (err error) {

	wrongSignin, err := us.repo.GetWrongSigninByUserID(ctx, tx, user.GetID())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		wrongSignin = usermodel.NewWrongSignin()
	}
	wrongSignin.SetUserID(user.GetID())
	wrongSignin.SetLastAttemptAt(time.Now().UTC())
	wrongSignin.SetFailedAttempts(wrongSignin.GetFailedAttempts() + 1)
	err = us.repo.UpdateWrongSignIn(ctx, tx, wrongSignin)
	if err != nil {
		return
	}

	if wrongSignin.GetFailedAttempts() >= usermodel.MaxWrongSigninAttempts {
		user.GetActiveRole().SetStatus(usermodel.StatusBlocked)
		user.GetActiveRole().SetStatusReason("Too many wrong signins attempts")
		user.SetLastSignInAttempt(time.Now().UTC())
		err = us.repo.UpdateUserByID(ctx, tx, user)
		if err != nil {
			return
		}
	}

	return
}
