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

// ConfirmEmailChange confirms a pending email change without creating or returning tokens.
// The user ID is extracted from context (SSOT).
func (us *userService) ConfirmEmailChange(ctx context.Context, code string) (err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return utils.InternalError("Failed to generate tracer")
	}
	defer spanEnd()

	// Obter o ID do usuário a partir do contexto
	userID, err := us.globalService.GetUserIDFromContext(ctx)
	if err != nil || userID == 0 {
		return utils.AuthenticationError("")
	}

	// Start transaction
	tx, txErr := us.globalService.StartTransaction(ctx)
	if txErr != nil {
		slog.Error("user.confirm_email_change.tx_start_error", "err", txErr)
		if mp := us.globalService.GetMetrics(); mp != nil {
			mp.IncrementEmailChangeConfirm("start_tx_error")
		}
		return utils.InternalError("Failed to start transaction")
	}
	defer func() {
		if err != nil {
			if rbErr := us.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				slog.Error("user.confirm_email_change.tx_rollback_error", "err", rbErr)
			}
		}
	}()

	err = us.confirmEmailChange(ctx, tx, userID, code)
	if err != nil {
		if mp := us.globalService.GetMetrics(); mp != nil {
			// map domain errors to results
			switch err {
			case utils.ErrEmailChangeNotPending:
				mp.IncrementEmailChangeConfirm("not_pending")
			case utils.ErrEmailChangeCodeInvalid:
				mp.IncrementEmailChangeConfirm("invalid")
			case utils.ErrEmailChangeCodeExpired:
				mp.IncrementEmailChangeConfirm("expired")
			case utils.ErrEmailAlreadyInUse:
				mp.IncrementEmailChangeConfirm("already_in_use")
			default:
				mp.IncrementEmailChangeConfirm("domain_error")
			}
		}
		return
	}

	// Após confirmar e-mail com sucesso, aplicar transição simples de status (na mesma transação)
	if _, changed, terr := us.applyStatusTransitionAfterContactChange(ctx, tx, true /*emailJustConfirmed*/); terr != nil {
		return terr
	} else {
		_ = changed // mudança pode ser no-op
	}

	if commitErr := us.globalService.CommitTransaction(ctx, tx); commitErr != nil {
		slog.Error("user.confirm_email_change.tx_commit_error", "err", commitErr)
		if mp := us.globalService.GetMetrics(); mp != nil {
			mp.IncrementEmailChangeConfirm("commit_error")
		}
		return utils.InternalError("Failed to commit transaction")
	}

	// Push notifications no longer dispatched automatically here.

	if mp := us.globalService.GetMetrics(); mp != nil {
		mp.IncrementEmailChangeConfirm("success")
	}

	return
}

func (us *userService) confirmEmailChange(ctx context.Context, tx *sql.Tx, userID int64, code string) (err error) {

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

	// Deve haver um código pendente e um novo e-mail definido
	if userValidation.GetEmailCode() == "" || userValidation.GetNewEmail() == "" {
		err = utils.ErrEmailChangeNotPending
		return
	}

	//check if the code is correct
	if !strings.EqualFold(userValidation.GetEmailCode(), code) {
		err = utils.ErrEmailChangeCodeInvalid
		return
	}

	//check if the validation is in time
	if userValidation.GetEmailCodeExp().Before(now) {
		err = utils.ErrEmailChangeCodeExpired
		return
	}

	// Carrega usuário com active role via Service (invariável: requer active role)
	user, err := us.GetUserByIDWithTx(ctx, tx, userID)
	if err != nil {
		return
	}

	// Verificar se o novo e-mail já está sendo utilizado por outro usuário
	if exist, verr := us.repo.ExistsEmailForAnotherUser(ctx, tx, userValidation.GetNewEmail(), userID); verr != nil {
		return verr
	} else if exist {
		return utils.ErrEmailAlreadyInUse
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

	err = us.repo.UpdateUserByID(ctx, tx, user)
	if err != nil {
		return
	}

	err = us.globalService.CreateAudit(ctx, tx, globalmodel.TableUsers, "email do usuário alterado")
	if err != nil {
		return
	}
	return
}
