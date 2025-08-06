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

func (us *userService) ConfirmEmailChange(ctx context.Context, userID int64, code string, deviceToken string) (tokens usermodel.Tokens, err error) {
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

	tokens, err = us.confirmEmailChange(ctx, tx, userID, code, deviceToken)
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

func (us *userService) confirmEmailChange(ctx context.Context, tx *sql.Tx, userID int64, code string, deviceToken string) (tokens usermodel.Tokens, err error) {

	now := time.Now()

	//read the user validation
	userValidation, err := us.repo.GetUserValidations(ctx, tx, userID)
	if err != nil {
		return
	}

	//verify is the user is on profile validation where email and phone should be validated
	//in this case the phone must be validated first
	if userValidation.GetPhoneCode() != "" {
		err = status.Error(codes.FailedPrecondition, "Phone must be validated first")
		return
	}

	//check if the user is awaiting email validation
	if userValidation.GetEmailCode() == "" {
		err = status.Error(codes.FailedPrecondition, "User is not awaiting email validation")
		return
	}

	//check if the code is correct
	if !strings.EqualFold(userValidation.GetEmailCode(), code) {
		err = status.Error(codes.InvalidArgument, "Invalid code")
		return
	}

	//check if the validation is in time
	if userValidation.GetEmailCodeExp().Before(now) {
		err = status.Error(codes.InvalidArgument, "Code expired")
		return
	}

	user, err := us.repo.GetUserByID(ctx, tx, userID)
	if err != nil {
		return
	}

	user.SetEmail(userValidation.GetNewEmail())
	//TODO: Remove this hardcoded deviceToken
	_ = deviceToken
	user.SetDeviceToken("dDWfs2iRThyJvzd_dSvyah:APA91bGp1GdU1zNsTzpaNb9gJpPdPTOVvJFpL2vpT52E7wemRocGtCe8HN5rpxk_Ys5NH4qo__7CD4_TZ0ahbTk2CyRaj36gCwlV9IANjFFtiQpQEvbSenw")
	// user.SetDeviceToken(deviceToken)

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
