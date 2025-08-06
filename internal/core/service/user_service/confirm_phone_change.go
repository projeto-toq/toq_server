package userservices

import (
	"context"
	"database/sql"
	"strings"
	"time"

	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (us *userService) ConfirmPhoneChange(ctx context.Context, userID int64, code string) (tokens usermodel.Tokens, err error) {
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

	tokens, err = us.confirmPhoneChange(ctx, tx, userID, code)
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

func (us *userService) confirmPhoneChange(ctx context.Context, tx *sql.Tx, userID int64, code string) (tokens usermodel.Tokens, err error) {
	now := time.Now().UTC()

	//read the user validation
	userValidation, err := us.repo.GetUserValidations(ctx, tx, userID)
	if err != nil {
		return
	}

	//check if the user is awaiting phone validation
	if userValidation.GetPhoneCode() == "" {
		err = status.Error(codes.FailedPrecondition, "User is not awaiting phone validation")
		return
	}

	//check if the code is correct
	if !strings.EqualFold(userValidation.GetPhoneCode(), code) {
		err = status.Error(codes.InvalidArgument, "Invalid code")
		return
	}

	//check if the validation is in time
	if userValidation.GetPhoneCodeExp().Before(now) {
		err = status.Error(codes.InvalidArgument, "Code expired")
		return
	}

	//read the user to update the phone number
	user, err := us.repo.GetUserByID(ctx, tx, userID)
	if err != nil {
		return
	}

	user.SetPhoneNumber(userValidation.GetNewPhone())

	//update the user validation
	userValidation.SetNewPhone("")
	userValidation.SetPhoneCode("")
	userValidation.SetPhoneCodeExp(time.Time{})

	err = us.repo.UpdateUserValidations(ctx, tx, userValidation)
	if err != nil {
		return
	}

	//update the user Status and create tokens if needed
	mustCreateTokens, err := us.UpdateUserValidationByUserRole(ctx, tx, &user, userValidation)
	if err != nil {
		return
	}

	if mustCreateTokens {
		tokens, err = us.CreateTokens(ctx, tx, user, false)
		if err != nil {
			return
		}
	}

	err = us.repo.UpdateUserByID(ctx, tx, user)
	if err != nil {
		return
	}

	err = us.globalService.CreateAudit(ctx, tx, globalmodel.TableUsers, "Alterada o telefone do usuário")
	if err != nil {
		return
	}

	return
}
