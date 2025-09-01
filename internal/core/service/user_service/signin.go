package userservices

import (
	"context"
	"database/sql"
	"log/slog"
	"time"

	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"

	"errors"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
	"golang.org/x/crypto/bcrypt"
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
	user, err := us.repo.GetUserByNationalID(ctx, tx, nationalID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			slog.Warn("Signin attempt with non-existent nationalID", "nationalID", nationalID)
			err = utils.AuthenticationError("Invalid credentials")
			return
		}
		slog.Error("Failed to get user by national ID", "nationalID", nationalID, "error", err)
		err = utils.InternalError("Failed to validate credentials")
		return
	}

	// Check if user is temporarily blocked before any password validation
	isBlocked, err := us.permissionService.IsUserTempBlocked(ctx, user.GetID())
	if err != nil {
		slog.Error("Failed to check if user is temporarily blocked", "userID", user.GetID(), "error", err)
		err = utils.InternalError("Failed to validate user status")
		return
	}

	if isBlocked {
		slog.Warn("Signin attempt from temporarily blocked user", "userID", user.GetID(), "nationalID", nationalID)
		err = utils.AuthenticationError("Account temporarily blocked due to too many failed attempts")
		return
	}

	// Get active roles via Permission Service
	userRoles, err := us.permissionService.GetUserRolesWithTx(ctx, tx, user.GetID())
	if err != nil {
		slog.Error("Failed to get user roles", "userID", user.GetID(), "error", err)
		err = utils.InternalError("Failed to validate user permissions")
		return
	}

	if len(userRoles) == 0 {
		slog.Warn("User has no active roles", "userID", user.GetID())
		err = utils.AuthenticationError("No active user roles")
		return
	}

	if len(userRoles) > 0 {
		user.SetActiveRole(userRoles[0])
	}

	// Comparar a senha fornecida com o hash armazenado (bcrypt)
	if bcrypt.CompareHashAndPassword([]byte(user.GetPassword()), []byte(password)) != nil {
		err = checkWrongSignin(ctx, tx, us, user)
		if err != nil {
			return
		}

		slog.Warn("Invalid password attempt", "userID", user.GetID(), "nationalID", nationalID)
		err = utils.AuthenticationError("Invalid credentials")
		return
	}

	_, err = us.repo.DeleteWrongSignInByUserID(ctx, tx, user.GetID())
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			slog.Error("Failed to delete wrong signin record", "userID", user.GetID(), "error", err)
			err = utils.InternalError("Failed to clear signin attempts")
			return
		}
	}

	// Clear any temporary blocks since login was successful
	isBlocked, errBlock := us.permissionService.IsUserTempBlocked(ctx, user.GetID())
	if errBlock != nil {
		slog.Warn("Failed to check temp block status after successful login", "userID", user.GetID(), "error", errBlock)
	} else if isBlocked {
		errUnblock := us.permissionService.UnblockUser(ctx, tx, user.GetID())
		if errUnblock != nil {
			slog.Error("Failed to unblock user after successful login", "userID", user.GetID(), "error", errUnblock)
		} else {
			slog.Info("User unblocked after successful login", "userID", user.GetID())
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
			wrongSignin = usermodel.NewWrongSignin()
		} else {
			slog.Error("Failed to get wrong signin record", "userID", user.GetID(), "error", err)
			return utils.InternalError("Failed to check signin attempts")
		}
	}

	wrongSignin.SetUserID(user.GetID())
	wrongSignin.SetLastAttemptAt(time.Now().UTC())
	wrongSignin.SetFailedAttempts(wrongSignin.GetFailedAttempts() + 1)

	err = us.repo.UpdateWrongSignIn(ctx, tx, wrongSignin)
	if err != nil {
		slog.Error("Failed to update wrong signin record", "userID", user.GetID(), "error", err)
		return utils.InternalError("Failed to update signin attempts")
	}

	if wrongSignin.GetFailedAttempts() >= usermodel.MaxWrongSigninAttempts {
		// Block user temporarily using permission service
		err = us.permissionService.BlockUserTemporarily(ctx, tx, user.GetID(), "Too many failed signin attempts")
		if err != nil {
			slog.Error("Failed to block user temporarily", "userID", user.GetID(), "error", err)
			return utils.InternalError("Failed to process security measures")
		}

		user.SetLastSignInAttempt(time.Now().UTC())
		err = us.repo.UpdateUserByID(ctx, tx, user)
		if err != nil {
			slog.Error("Failed to update user last signin attempt", "userID", user.GetID(), "error", err)
			return utils.InternalError("Failed to update user record")
		}

		slog.Warn("User temporarily blocked due to too many failed signin attempts", "userID", user.GetID(), "attempts", wrongSignin.GetFailedAttempts())
	}

	return
}
