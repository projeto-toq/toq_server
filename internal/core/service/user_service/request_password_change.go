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

func (us *userService) RequestPasswordChange(ctx context.Context, nationalID string) (err error) {

	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	tx, err := us.globalService.StartTransaction(ctx)
	if err != nil {
		return
	}

	err = us.requestPasswordChange(ctx, tx, nationalID)
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

func (us *userService) requestPasswordChange(ctx context.Context, tx *sql.Tx, nationaID string) (err error) {

	user, err := us.repo.GetUserByNationalID(ctx, tx, nationaID)
	if err != nil {
		return
	}

	var validation usermodel.ValidationInterface

	//set the user validation as pending for password
	validation, err = us.repo.GetUserValidations(ctx, tx, user.GetID())
	if err != nil {
		if status.Code(err) != codes.NotFound {
			return
		}
		validation = usermodel.NewValidation()
	}

	validation.SetUserID(user.GetID())
	validation.SetPasswordCode(us.random6Digits())
	validation.SetPasswordCodeExp(time.Now().UTC().Add(usermodel.ValidationCodeExpiration))

	err = us.repo.UpdateUserValidations(ctx, tx, validation)
	if err != nil {
		return
	}
	err = us.globalService.SendNotification(ctx, user, globalmodel.NotificationPasswordChange, validation.GetPasswordCode())
	if err != nil {
		return
	}

	//update the lastactivity
	err = us.repo.UpdateUserLastActivity(ctx, tx, user.GetID())
	if err != nil {
		return
	}

	return
}
