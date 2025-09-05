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
		utils.SetSpanError(ctx, txErr)
		return utils.InternalError("Failed to start transaction")
	}
	defer func() {
		if err != nil {
			if rbErr := us.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				slog.Error("user.confirm_password_change.tx_rollback_error", "err", rbErr)
				utils.SetSpanError(ctx, rbErr)
			}
		}
	}()

	err = us.confirmPasswordChange(ctx, tx, nationalID, password, code)
	if err != nil {
		return
	}

	if commitErr := us.globalService.CommitTransaction(ctx, tx); commitErr != nil {
		slog.Error("user.confirm_password_change.tx_commit_error", "err", commitErr)
		utils.SetSpanError(ctx, commitErr)
		return utils.InternalError("Failed to commit transaction")
	}

	return
}

func (us *userService) confirmPasswordChange(ctx context.Context, tx *sql.Tx, nationalID string, password string, code string) (err error) {
	now := time.Now().UTC()

	user, err := us.repo.GetUserByNationalID(ctx, tx, nationalID)
	if err != nil {
		// Não revelar existência do usuário: mapear ausência como fluxo não pendente
		if errors.Is(err, sql.ErrNoRows) {
			return utils.ErrPasswordChangeNotPending
		}
		// Outros erros são infra
		utils.SetSpanError(ctx, err)
		slog.Error("user.confirm_password_change.stage_error", "stage", "get_user_by_national_id", "err", err)
		return utils.InternalError("Failed to get user")
	}

	//read the user validation
	userValidation, err := us.repo.GetUserValidations(ctx, tx, user.GetID())
	if err != nil {
		// Sem validação pendente → domínio
		if errors.Is(err, sql.ErrNoRows) {
			return utils.ErrPasswordChangeNotPending
		}
		// Infra
		utils.SetSpanError(ctx, err)
		slog.Error("user.confirm_password_change.stage_error", "stage", "get_validations", "err", err)
		return utils.InternalError("Failed to get user validations")
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
		utils.SetSpanError(ctx, err)
		slog.Error("user.confirm_password_change.stage_error", "stage", "update_validations", "err", err)
		return utils.InternalError("Failed to update validations")
	}

	//delete the temp_wrong_signin
	_, err = us.repo.DeleteWrongSignInByUserID(ctx, tx, user.GetID())
	if err != nil {
		// Ignorar ausência de registros; tratar demais erros
		if !errors.Is(err, sql.ErrNoRows) {
			utils.SetSpanError(ctx, err)
			slog.Error("user.confirm_password_change.stage_error", "stage", "delete_wrong_signin", "err", err)
			return utils.InternalError("Failed to cleanup wrong sign in attempts")
		}
	}

	user.SetLastActivityAt(now)

	err = us.repo.UpdateUserPasswordByID(ctx, tx, user)
	if err != nil {
		utils.SetSpanError(ctx, err)
		slog.Error("user.confirm_password_change.stage_error", "stage", "update_user_password", "err", err)
		return utils.InternalError("Failed to update user password")
	}

	err = us.globalService.CreateAudit(ctx, tx, globalmodel.TableUsers, "Alterada a senha do usuário")
	if err != nil {
		utils.SetSpanError(ctx, err)
		slog.Error("user.confirm_password_change.stage_error", "stage", "audit", "err", err)
		return utils.InternalError("Failed to create audit")
	}

	return
}
