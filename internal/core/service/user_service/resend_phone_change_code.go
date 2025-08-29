package userservices

import (
	"context"
	"database/sql"
	"time"

	
	
"github.com/giulio-alfieri/toq_server/internal/core/utils"
"errors"
)

func (us *userService) ResendPhoneChangeCode(ctx context.Context, userID int64) (code string, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	tx, err := us.globalService.StartTransaction(ctx)
	if err != nil {
		return
	}

	code, err = us.resendPhoneChangeCode(ctx, tx, userID)
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

func (us *userService) resendPhoneChangeCode(ctx context.Context, tx *sql.Tx, userID int64) (code string, err error) {
	validation, err := us.repo.GetUserValidations(ctx, tx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = utils.ErrInternalServer
		}
		return
	}

	phoneCode := validation.GetPhoneCode()
	if phoneCode == "" {
		err = utils.ErrInternalServer
		return
	}

	// Verificar se código não expirou
	if time.Now().UTC().After(validation.GetPhoneCodeExp()) {
		err = utils.ErrInternalServer
		return
	}

	return phoneCode, nil
}
