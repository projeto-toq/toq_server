package userservices

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"os"
	"strings"
	"time"

	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

// ConfirmPhoneChange confirms a pending phone change without creating or returning tokens.
func (us *userService) ConfirmPhoneChange(ctx context.Context, code string) (err error) {
	// Padronizar: gerar tracer antes e só então obter userID do contexto (como no fluxo de e-mail)
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return utils.InternalError("Failed to generate tracer")
	}
	defer spanEnd()

	// Obter o ID do usuário do contexto (SSOT)
	userID, err := us.globalService.GetUserIDFromContext(ctx)
	if err != nil || userID == 0 {
		// erro de autenticação already carries source here
		return utils.AuthenticationError("")
	}
	if v := os.Getenv("TOQ_DEBUG_ERROR_TRACE"); v == "true" {
		slog.Debug("user.confirm_phone_change.start", "user_id", userID)
	}

	// Start transaction
	tx, txErr := us.globalService.StartTransaction(ctx)
	if txErr != nil {
		slog.Error("user.confirm_phone_change.tx_start_error", "err", txErr)
		utils.SetSpanError(ctx, txErr)
		if mp := us.globalService.GetMetrics(); mp != nil {
			mp.IncrementPhoneChangeConfirm("start_tx_error")
		}
		return utils.InternalError("Failed to start transaction")
	}
	defer func() {
		if err != nil {
			if rbErr := us.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				slog.Error("user.confirm_phone_change.tx_rollback_error", "err", rbErr)
			}
		}
	}()

	err = us.confirmPhoneChange(ctx, tx, userID, code)
	if err != nil {
		if mp := us.globalService.GetMetrics(); mp != nil {
			switch err {
			case utils.ErrPhoneChangeNotPending:
				mp.IncrementPhoneChangeConfirm("not_pending")
			case utils.ErrPhoneChangeCodeInvalid:
				mp.IncrementPhoneChangeConfirm("invalid")
			case utils.ErrPhoneChangeCodeExpired:
				mp.IncrementPhoneChangeConfirm("expired")
			case utils.ErrPhoneAlreadyInUse:
				mp.IncrementPhoneChangeConfirm("already_in_use")
			default:
				mp.IncrementPhoneChangeConfirm("domain_error")
			}
		}
		if v := os.Getenv("TOQ_DEBUG_ERROR_TRACE"); v == "true" {
			status := 500
			if derr, ok := err.(utils.DomainError); ok {
				status = derr.Code()
			}
			slog.Debug("user.confirm_phone_change.error", "stage", "confirm_phone_change_call", "status", status)
		}
		// Se já é um DomainError, envolvemos com source do service; caso contrário, mapeamos para InternalError
		if derr, ok := err.(utils.DomainError); ok {
			return utils.WrapDomainErrorWithSource(derr)
		}
		utils.SetSpanError(ctx, err)
		return utils.InternalError("Failed to confirm phone change")
	}

	// Após confirmar telefone com sucesso, aplicar transição simples de status (na mesma transação)
	if _, changed, terr := us.applyStatusTransitionAfterContactChange(ctx, tx, false /*emailJustConfirmed*/); terr != nil {
		return terr
	} else {
		_ = changed
	}

	if commitErr := us.globalService.CommitTransaction(ctx, tx); commitErr != nil {
		slog.Error("user.confirm_phone_change.tx_commit_error", "err", commitErr)
		utils.SetSpanError(ctx, commitErr)
		if mp := us.globalService.GetMetrics(); mp != nil {
			mp.IncrementPhoneChangeConfirm("commit_error")
		}
		return utils.InternalError("Failed to commit transaction")
	}

	if mp := us.globalService.GetMetrics(); mp != nil {
		mp.IncrementPhoneChangeConfirm("success")
	}
	return
}

func (us *userService) confirmPhoneChange(ctx context.Context, tx *sql.Tx, userID int64, code string) (err error) {
	now := time.Now().UTC()

	//read the user validation
	userValidation, err := us.repo.GetUserValidations(ctx, tx, userID)
	if err != nil {
		// mapear ausência para domínio; outras falhas são infra
		if errors.Is(err, sql.ErrNoRows) {
			if v := os.Getenv("TOQ_DEBUG_ERROR_TRACE"); v == "true" {
				slog.Debug("user.confirm_phone_change.stage", "stage", "get_validations", "outcome", "no_rows")
			}
			return utils.ErrPhoneChangeNotPending
		}
		utils.SetSpanError(ctx, err)
		slog.Error("user.confirm_phone_change.stage_error", "stage", "get_validations", "err", err)
		return utils.InternalError("Failed to get user validations")
	}
	if v := os.Getenv("TOQ_DEBUG_ERROR_TRACE"); v == "true" {
		slog.Debug("user.confirm_phone_change.stage", "stage", "get_validations", "outcome", "found")
	}

	// Deve haver um código pendente e um novo telefone definido (espelha o fluxo de e-mail)
	if userValidation.GetPhoneCode() == "" || userValidation.GetNewPhone() == "" {
		if v := os.Getenv("TOQ_DEBUG_ERROR_TRACE"); v == "true" {
			slog.Debug("user.confirm_phone_change.stage", "stage", "validate_pending", "outcome", "not_pending")
		}
		err = utils.ErrPhoneChangeNotPending
		return
	}

	//check if the code is correct
	if !strings.EqualFold(userValidation.GetPhoneCode(), code) {
		if v := os.Getenv("TOQ_DEBUG_ERROR_TRACE"); v == "true" {
			slog.Debug("user.confirm_phone_change.stage", "stage", "validate_code_match", "outcome", "mismatch")
		}
		err = utils.ErrPhoneChangeCodeInvalid
		return
	}

	//check if the validation is in time
	if userValidation.GetPhoneCodeExp().Before(now) {
		if v := os.Getenv("TOQ_DEBUG_ERROR_TRACE"); v == "true" {
			slog.Debug("user.confirm_phone_change.stage", "stage", "validate_code_exp", "outcome", "expired")
		}
		err = utils.ErrPhoneChangeCodeExpired
		return
	}

	//read the user to update the phone number
	// Carrega usuário com active role via Service (invariável: requer active role)
	user, err := us.GetUserByIDWithTx(ctx, tx, userID)
	if err != nil {
		if v := os.Getenv("TOQ_DEBUG_ERROR_TRACE"); v == "true" {
			slog.Debug("user.confirm_phone_change.stage", "stage", "get_user_by_id", "outcome", "domain_error")
		}
		// GetUserByIDWithTx já retorna DomainError
		return err
	}
	if v := os.Getenv("TOQ_DEBUG_ERROR_TRACE"); v == "true" {
		slog.Debug("user.confirm_phone_change.stage", "stage", "get_user_by_id", "outcome", "ok")
	}

	// Uniqueness check before setting
	if exist, verr := us.repo.ExistsPhoneForAnotherUser(ctx, tx, userValidation.GetNewPhone(), userID); verr != nil {
		utils.SetSpanError(ctx, verr)
		slog.Error("user.confirm_phone_change.stage_error", "stage", "exists_phone_for_another_user", "err", verr)
		return utils.InternalError("Failed to check phone uniqueness")
	} else if exist {
		if v := os.Getenv("TOQ_DEBUG_ERROR_TRACE"); v == "true" {
			slog.Debug("user.confirm_phone_change.stage", "stage", "exists_phone_for_another_user", "outcome", "exists")
		}
		return utils.ErrPhoneAlreadyInUse
	}
	if v := os.Getenv("TOQ_DEBUG_ERROR_TRACE"); v == "true" {
		slog.Debug("user.confirm_phone_change.stage", "stage", "exists_phone_for_another_user", "outcome", "not_exists")
	}
	user.SetPhoneNumber(userValidation.GetNewPhone())

	//update the user validation
	userValidation.SetNewPhone("")
	userValidation.SetPhoneCode("")
	userValidation.SetPhoneCodeExp(time.Time{})

	err = us.repo.UpdateUserValidations(ctx, tx, userValidation)
	if err != nil {
		utils.SetSpanError(ctx, err)
		slog.Error("user.confirm_phone_change.stage_error", "stage", "update_validations", "err", err)
		return utils.InternalError("Failed to update validations")
	}
	if v := os.Getenv("TOQ_DEBUG_ERROR_TRACE"); v == "true" {
		slog.Debug("user.confirm_phone_change.stage", "stage", "update_validations", "outcome", "ok")
	}

	// Status transition is handled centrally by applyStatusTransitionAfterContactChange in the public method.
	// Não realizar transições de status aqui para evitar duplicidade e inconsistências.

	err = us.repo.UpdateUserByID(ctx, tx, user)
	if err != nil {
		utils.SetSpanError(ctx, err)
		slog.Error("user.confirm_phone_change.stage_error", "stage", "update_user", "err", err)
		return utils.InternalError("Failed to update user")
	}
	if v := os.Getenv("TOQ_DEBUG_ERROR_TRACE"); v == "true" {
		slog.Debug("user.confirm_phone_change.stage", "stage", "update_user", "outcome", "ok")
	}

	err = us.globalService.CreateAudit(ctx, tx, globalmodel.TableUsers, "Alterada o telefone do usuário")
	if err != nil {
		utils.SetSpanError(ctx, err)
		slog.Error("user.confirm_phone_change.stage_error", "stage", "audit", "err", err)
		return utils.InternalError("Failed to create audit")
	}
	if v := os.Getenv("TOQ_DEBUG_ERROR_TRACE"); v == "true" {
		slog.Debug("user.confirm_phone_change.stage", "stage", "audit", "outcome", "ok")
	}

	return
}
