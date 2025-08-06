package userservices

import (
	"context"
	"database/sql"
	"strings"
	"time"

	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (us *userService) ConfirmPasswordChange(ctx context.Context, nationalID string, password string, code string) (err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	// Start transaction
	tx, err := us.globalService.StartTransaction(ctx)
	if err != nil {
		return
	}

	err = us.confirmPasswordChange(ctx, tx, nationalID, password, code)
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

func (us *userService) confirmPasswordChange(ctx context.Context, tx *sql.Tx, nationalID string, password string, code string) (err error) {
	now := time.Now().UTC()

	user, err := us.repo.GetUserByNationalID(ctx, tx, nationalID)
	if err != nil {
		return
	}

	//read the user validation
	userValidation, err := us.repo.GetUserValidations(ctx, tx, user.GetID())
	if err != nil {
		return
	}

	//check if the user is awaiting password reset
	if userValidation.GetPasswordCode() == "" {
		err = status.Error(codes.FailedPrecondition, "User is not awaiting password validation")
		return
	}

	//check if the code is correct
	if !strings.EqualFold(userValidation.GetPasswordCode(), code) {
		err = status.Error(codes.InvalidArgument, "Invalid code")
		return
	}

	//check if the validation is in time
	if userValidation.GetPasswordCodeExp().Before(now) {
		err = status.Error(codes.InvalidArgument, "Code expired")
		return
	}

	user.SetPassword(us.encryptPassword(password))

	//update the user validation
	userValidation.SetPasswordCode("")
	userValidation.SetPasswordCodeExp(time.Time{})

	err = us.repo.UpdateUserValidations(ctx, tx, userValidation)
	if err != nil {
		return
	}

	//delete the temp_wrong_signin
	_, err = us.repo.DeleteWrongSignInByUserID(ctx, tx, user.GetID())
	if err != nil && status.Code(err) != codes.NotFound {
		return
	}

	user.SetLastActivityAt(now)

	err = us.repo.UpdateUserPasswordByID(ctx, tx, user)
	if err != nil {
		return
	}

	err = us.globalService.CreateAudit(ctx, tx, globalmodel.TableUsers, "Alterada a senha do usuário")
	if err != nil {
		return
	}

	return
}
