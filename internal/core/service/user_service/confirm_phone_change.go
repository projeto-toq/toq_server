package userservices

import (
	"context"
	"database/sql"
	"strings"
	"time"

	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

// ConfirmPhoneChange confirms a pending phone change without creating or returning tokens.
func (us *userService) ConfirmPhoneChange(ctx context.Context, code string) (err error) {
	// Obter o ID do usuário do contexto (SSOT)
	userID, err := us.globalService.GetUserIDFromContext(ctx)
	if err != nil || userID == 0 {
		return utils.ErrInternalServer
	}
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	// Start transaction
	tx, err := us.globalService.StartTransaction(ctx)
	if err != nil {
		if mp := us.globalService.GetMetrics(); mp != nil {
			mp.IncrementPhoneChangeConfirm("start_tx_error")
		}
		return
	}

	err = us.confirmPhoneChange(ctx, tx, userID, code)
	if err != nil {
		us.globalService.RollbackTransaction(ctx, tx)
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

	err = us.globalService.CommitTransaction(ctx, tx)
	if err != nil {
		us.globalService.RollbackTransaction(ctx, tx)
		if mp := us.globalService.GetMetrics(); mp != nil {
			mp.IncrementPhoneChangeConfirm("commit_error")
		}
		return
	}

	if mp := us.globalService.GetMetrics(); mp != nil {
		mp.IncrementPhoneChangeConfirm("success")
	}
	// Após confirmar telefone com sucesso, aplicar a transição de status adequada
	if _, _, terr := us.ApplyUserStatusTransitionAfterPhoneConfirmed(ctx); terr != nil {
		// Não falha o fluxo principal; apenas registra
		_ = terr
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

	//check if the user is awaiting phone validation
	if userValidation.GetPhoneCode() == "" {
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
	user, err := us.repo.GetUserByID(ctx, tx, userID)
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

	// Update user status if needed, but do not create tokens in this flow
	_, err = us.UpdateUserValidationByUserRole(ctx, tx, &user, userValidation)
	if err != nil {
		return
	}

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
