package userservices

import (
	"context"
	"database/sql"
	"time"

	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (us *userService) RequestEmailChange(ctx context.Context, userID int64, newEmail string) (err error) {

	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	tx, err := us.globalService.StartTransaction(ctx)
	if err != nil {
		return
	}

	err = us.requestEmailChange(ctx, tx, userID, newEmail)
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

func (us *userService) requestEmailChange(ctx context.Context, tx *sql.Tx, id int64, email string) (err error) {

	user, err := us.repo.GetUserByID(ctx, tx, id)
	if err != nil {
		return
	}

	var validation usermodel.ValidationInterface

	//set the user validation as pending for email
	validation, err = us.repo.GetUserValidations(ctx, tx, user.GetID())
	if err != nil {
		if status.Code(err) != codes.NotFound {
			return
		}
		validation = usermodel.NewValidation()
	}
	validation.SetUserID(user.GetID())
	validation.SetEmailCode(us.random6Digits())
	validation.SetEmailCodeExp(time.Now().UTC().Add(usermodel.ValidationCodeExpiration))
	validation.SetNewEmail(email)

	err = us.repo.UpdateUserValidations(ctx, tx, validation)
	if err != nil {
		return
	}
	err = us.globalService.SendNotification(ctx, user, globalmodel.NotificationEmailChange, validation.GetEmailCode())

	if err != nil {
		return
	}

	err = us.repo.UpdateUserLastActivity(ctx, tx, user.GetID())
	if err != nil {
		return
	}

	return
}
