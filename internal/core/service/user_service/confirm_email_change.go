package userservices

import (
	"context"
	"database/sql"
	"strings"
	"time"

	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	
	
"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (us *userService) ConfirmEmailChange(ctx context.Context, userID int64, code string) (tokens usermodel.Tokens, err error) {
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

	tokens, err = us.confirmEmailChange(ctx, tx, userID, code)
	if err != nil {
		us.globalService.RollbackTransaction(ctx, tx)
		return
	}

	err = us.globalService.CommitTransaction(ctx, tx)
	if err != nil {
		us.globalService.RollbackTransaction(ctx, tx)
		return
	}

	// Push notifications no longer dispatched automatically here.

	return
}

func (us *userService) confirmEmailChange(ctx context.Context, tx *sql.Tx, userID int64, code string) (tokens usermodel.Tokens, err error) {

	now := time.Now()

	//read the user validation
	userValidation, err := us.repo.GetUserValidations(ctx, tx, userID)
	if err != nil {
		return
	}

	// //verify is the user is on profile validation where email and phone should be validated
	// //in this case the phone must be validated first
	// if userValidation.GetPhoneCode() != "" {
	// 	err = utils.ErrInternalServer
	// 	return
	// }

	//check if the user is awaiting email validation
	if userValidation.GetEmailCode() == "" {
		err = utils.ErrInternalServer
		return
	}

	//check if the code is correct
	if !strings.EqualFold(userValidation.GetEmailCode(), code) {
		err = utils.ErrInternalServer
		return
	}

	//check if the validation is in time
	if userValidation.GetEmailCodeExp().Before(now) {
		err = utils.ErrInternalServer
		return
	}

	user, err := us.repo.GetUserByID(ctx, tx, userID)
	if err != nil {
		return
	}

	user.SetEmail(userValidation.GetNewEmail())
	// Device token capture removed from email confirmation.

	//update the user validation
	userValidation.SetNewEmail("")
	userValidation.SetEmailCode("")
	userValidation.SetEmailCodeExp(time.Time{})

	err = us.repo.UpdateUserValidations(ctx, tx, userValidation)
	if err != nil {
		return
	}

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

	err = us.globalService.CreateAudit(ctx, tx, globalmodel.TableUsers, "alterado e email do usuário")
	if err != nil {
		return
	}
	return
}
