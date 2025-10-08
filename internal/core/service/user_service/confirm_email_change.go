package userservices

import (
	"context"
	"database/sql"
	"errors"
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

	ctx = utils.ContextWithLogger(ctx)

	// Obter o ID do usuário a partir do contexto
	userID, err := us.globalService.GetUserIDFromContext(ctx)
	if err != nil || userID == 0 {
		return utils.AuthenticationError("")
	}

	// Start transaction
	tx, txErr := us.globalService.StartTransaction(ctx)
	if txErr != nil {
		utils.LoggerFromContext(ctx).Error("user.confirm_email_change.tx_start_error", "err", txErr)
		utils.SetSpanError(ctx, txErr)
		return utils.InternalError("Failed to start transaction")
	}
	defer func() {
		if err != nil {
			if rbErr := us.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.LoggerFromContext(ctx).Error("user.confirm_email_change.tx_rollback_error", "err", rbErr)
				utils.SetSpanError(ctx, rbErr)
			}
		}
	}()

	err = us.confirmEmailChange(ctx, tx, userID, code)
	if err != nil {
		return
	}

	// Após confirmar e-mail com sucesso, aplicar transição simples de status (na mesma transação)
	if _, changed, terr := us.applyStatusTransitionAfterContactChange(ctx, tx, true /*emailJustConfirmed*/); terr != nil {
		return terr
	} else {
		_ = changed // mudança pode ser no-op
	}

	if commitErr := us.globalService.CommitTransaction(ctx, tx); commitErr != nil {
		utils.LoggerFromContext(ctx).Error("user.confirm_email_change.tx_commit_error", "err", commitErr)
		utils.SetSpanError(ctx, commitErr)
		return utils.InternalError("Failed to commit transaction")
	}

	// Push notifications no longer dispatched automatically here.

	return
}

func (us *userService) confirmEmailChange(ctx context.Context, tx *sql.Tx, userID int64, code string) (err error) {

	now := time.Now()

	//read the user validation
	userValidation, err := us.repo.GetUserValidations(ctx, tx, userID)
	if err != nil {
		// Se não há validação, tratar como fluxo não pendente (domínio)
		if errors.Is(err, sql.ErrNoRows) {
			return utils.ErrEmailChangeNotPending
		}
		// Outros erros são infraestrutura
		utils.SetSpanError(ctx, err)
		utils.LoggerFromContext(ctx).Error("user.confirm_email_change.stage_error", "stage", "get_validations", "err", err)
		return utils.InternalError("Failed to get user validations")
	}

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
		utils.SetSpanError(ctx, verr)
		utils.LoggerFromContext(ctx).Error("user.confirm_email_change.stage_error", "stage", "exists_email_for_another_user", "err", verr)
		return utils.InternalError("Failed to check email uniqueness")
	} else if exist {
		return utils.ErrEmailAlreadyInUse
	}

	user.SetEmail(userValidation.GetNewEmail())

	//update the user validation
	userValidation.SetNewEmail("")
	userValidation.SetEmailCode("")
	userValidation.SetEmailCodeExp(time.Time{})

	err = us.repo.UpdateUserValidations(ctx, tx, userValidation)
	if err != nil {
		utils.SetSpanError(ctx, err)
		utils.LoggerFromContext(ctx).Error("user.confirm_email_change.stage_error", "stage", "update_validations", "err", err)
		return utils.InternalError("Failed to update validations")
	}

	err = us.repo.UpdateUserByID(ctx, tx, user)
	if err != nil {
		utils.SetSpanError(ctx, err)
		utils.LoggerFromContext(ctx).Error("user.confirm_email_change.stage_error", "stage", "update_user", "err", err)
		return utils.InternalError("Failed to update user")
	}

	err = us.globalService.CreateAudit(ctx, tx, globalmodel.TableUsers, "email do usuário alterado")
	if err != nil {
		utils.SetSpanError(ctx, err)
		utils.LoggerFromContext(ctx).Error("user.confirm_email_change.stage_error", "stage", "audit", "err", err)
		return utils.InternalError("Failed to create audit")
	}
	return
}
