package userservices

import (
	"context"
	"database/sql"
	"log/slog"
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
		return utils.AuthenticationError("")
	}

	// Start transaction
	tx, txErr := us.globalService.StartTransaction(ctx)
	if txErr != nil {
		slog.Error("user.confirm_phone_change.tx_start_error", "err", txErr)
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
		return
	}

	// Após confirmar telefone com sucesso, aplicar transição simples de status (na mesma transação)
	if _, changed, terr := us.applyStatusTransitionAfterContactChange(ctx, tx, false /*emailJustConfirmed*/); terr != nil {
		return terr
	} else {
		_ = changed
	}

	if commitErr := us.globalService.CommitTransaction(ctx, tx); commitErr != nil {
		slog.Error("user.confirm_phone_change.tx_commit_error", "err", commitErr)
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
		return
	}

	// Deve haver um código pendente e um novo telefone definido (espelha o fluxo de e-mail)
	if userValidation.GetPhoneCode() == "" || userValidation.GetNewPhone() == "" {
		err = utils.ErrPhoneChangeNotPending
		return
	}

	//check if the code is correct
	if !strings.EqualFold(userValidation.GetPhoneCode(), code) {
		err = utils.ErrPhoneChangeCodeInvalid
		return
	}

	//check if the validation is in time
	if userValidation.GetPhoneCodeExp().Before(now) {
		err = utils.ErrPhoneChangeCodeExpired
		return
	}

	//read the user to update the phone number
	// Carrega usuário com active role via Service (invariável: requer active role)
	user, err := us.GetUserByIDWithTx(ctx, tx, userID)
	if err != nil {
		return
	}

	// Uniqueness check before setting
	if exist, verr := us.repo.ExistsPhoneForAnotherUser(ctx, tx, userValidation.GetNewPhone(), userID); verr != nil {
		return verr
	} else if exist {
		return utils.ErrPhoneAlreadyInUse
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

	// Status transition is handled centrally by applyStatusTransitionAfterContactChange in the public method.
	// Não realizar transições de status aqui para evitar duplicidade e inconsistências.

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
