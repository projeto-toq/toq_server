package userservices

import (
	"context"
	"database/sql"
	"time"

	
	
"github.com/giulio-alfieri/toq_server/internal/core/utils"
"errors"
)

func (us *userService) ResendEmailChangeCode(ctx context.Context, userID int64) (code string, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	tx, err := us.globalService.StartTransaction(ctx)
	if err != nil {
		return
	}

	code, err = us.resendEmailChangeCode(ctx, tx, userID)
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

func (us *userService) resendEmailChangeCode(ctx context.Context, tx *sql.Tx, userID int64) (code string, err error) {
	validation, err := us.repo.GetUserValidations(ctx, tx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = utils.ErrInternalServer
		}
		return
	}

	emailCode := validation.GetEmailCode()
	if emailCode == "" {
		err = utils.ErrInternalServer
		return
	}

	// Verificar se código não expirou
	if time.Now().UTC().After(validation.GetEmailCodeExp()) {
		err = utils.ErrInternalServer
		return
	}

	return emailCode, nil
}
