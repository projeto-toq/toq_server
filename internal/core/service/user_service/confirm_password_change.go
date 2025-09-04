package userservices

import (
	"context"
	"database/sql"
	"log/slog"
	"strings"
	"time"

	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"

	"errors"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (us *userService) ConfirmPasswordChange(ctx context.Context, nationalID string, password string, code string) (err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return utils.InternalError("Failed to generate tracer")
	}
	defer spanEnd()

	// Start transaction
	tx, txErr := us.globalService.StartTransaction(ctx)
	if txErr != nil {
		slog.Error("user.confirm_password_change.tx_start_error", "err", txErr)
		if mp := us.globalService.GetMetrics(); mp != nil {
			mp.IncrementPasswordChangeConfirm("start_tx_error")
		}
		return utils.InternalError("Failed to start transaction")
	}
	defer func() {
		if err != nil {
			if rbErr := us.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				slog.Error("user.confirm_password_change.tx_rollback_error", "err", rbErr)
			}
		}
	}()

	err = us.confirmPasswordChange(ctx, tx, nationalID, password, code)
	if err != nil {
		if mp := us.globalService.GetMetrics(); mp != nil {
			switch err {
			case utils.ErrPasswordChangeNotPending:
				mp.IncrementPasswordChangeConfirm("not_pending")
			case utils.ErrPasswordChangeCodeInvalid:
				mp.IncrementPasswordChangeConfirm("invalid")
			case utils.ErrPasswordChangeCodeExpired:
				mp.IncrementPasswordChangeConfirm("expired")
			default:
				mp.IncrementPasswordChangeConfirm("domain_error")
			}
		}
		return
	}

	if commitErr := us.globalService.CommitTransaction(ctx, tx); commitErr != nil {
		slog.Error("user.confirm_password_change.tx_commit_error", "err", commitErr)
		if mp := us.globalService.GetMetrics(); mp != nil {
			mp.IncrementPasswordChangeConfirm("commit_error")
		}
		return utils.InternalError("Failed to commit transaction")
	}

	if mp := us.globalService.GetMetrics(); mp != nil {
		mp.IncrementPasswordChangeConfirm("success")
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
		err = utils.ErrPasswordChangeNotPending
		return
	}

	//check if the code is correct
	if !strings.EqualFold(userValidation.GetPasswordCode(), code) {
		err = utils.ErrPasswordChangeCodeInvalid
		return
	}

	//check if the validation is in time
	if userValidation.GetPasswordCodeExp().Before(now) {
		err = utils.ErrPasswordChangeCodeExpired
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
	if err != nil && errors.Is(err, sql.ErrNoRows) {
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
