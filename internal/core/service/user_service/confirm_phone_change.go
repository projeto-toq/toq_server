package userservices

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"strings"
	"time"

	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"

	derrors "github.com/giulio-alfieri/toq_server/internal/core/derrors"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

// ConfirmPhoneChange confirms a pending phone change without creating or returning tokens.
func (us *userService) ConfirmPhoneChange(ctx context.Context, code string) (err error) {
	// Padronizar: gerar tracer antes e só então obter userID do contexto (como no fluxo de e-mail)
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return derrors.Infra("Failed to generate tracer", err)
	}
	defer spanEnd()

	// Obter o ID do usuário do contexto (SSOT)
	userID, err := us.globalService.GetUserIDFromContext(ctx)
	if err != nil || userID == 0 {
		return derrors.Auth("Authentication required")
	}
	// start debug log removed; keep error logs only

	// Start transaction
	tx, txErr := us.globalService.StartTransaction(ctx)
	if txErr != nil {
		slog.Error("user.confirm_phone_change.tx_start_error", "err", txErr)
		utils.SetSpanError(ctx, txErr)
		return derrors.Infra("Failed to start transaction", txErr)
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
		// Retornar o próprio erro de domínio (sentinela) ou erro de infra (já logado acima)
		// Em caso de erro inesperado (genérico), classificar como Infra
		if _, ok := derrors.AsKind(err); !ok {
			utils.SetSpanError(ctx, err)
			return derrors.Infra("Failed to confirm phone change", err)
		}
		return err
	}

	// Após confirmar telefone com sucesso, aplicar transição simples de status (na mesma transação)
	if _, changed, terr := us.applyStatusTransitionAfterContactChange(ctx, tx, false /*emailJustConfirmed*/); terr != nil {
		// Classificar transição: domínio vs infra
		if _, ok := derrors.AsKind(terr); ok {
			// erro de domínio da transição – propagar
			return terr
		}
		// infra: logar e mapear
		slog.Error("user.confirm_phone_change.apply_status_transition_error", "err", terr)
		utils.SetSpanError(ctx, terr)
		return derrors.Infra("Failed to apply status transition", terr)
	} else {
		_ = changed
	}

	if commitErr := us.globalService.CommitTransaction(ctx, tx); commitErr != nil {
		slog.Error("user.confirm_phone_change.tx_commit_error", "err", commitErr)
		utils.SetSpanError(ctx, commitErr)
		return derrors.Infra("Failed to commit transaction", commitErr)
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
			return derrors.ErrPhoneChangeNotPending
		}
		utils.SetSpanError(ctx, err)
		slog.Error("user.confirm_phone_change.stage_error", "stage", "get_validations", "err", err)
		return derrors.Infra("Failed to get user validations", err)
	}

	// Deve haver um código pendente e um novo telefone definido (espelha o fluxo de e-mail)
	if userValidation.GetPhoneCode() == "" || userValidation.GetNewPhone() == "" {
		err = derrors.ErrPhoneChangeNotPending
		return
	}

	//check if the code is correct
	if !strings.EqualFold(userValidation.GetPhoneCode(), code) {
		err = derrors.ErrPhoneChangeCodeInvalid
		return
	}

	//check if the validation is in time
	if userValidation.GetPhoneCodeExp().Before(now) {
		err = derrors.ErrPhoneChangeCodeExpired
		return
	}

	//read the user to update the phone number
	// Carrega usuário com active role via Service (invariável: requer active role)
	user, err := us.GetUserByIDWithTx(ctx, tx, userID)
	if err != nil {
		// GetUserByIDWithTx já retorna DomainError
		return err
	}

	// Uniqueness check before setting
	if exist, verr := us.repo.ExistsPhoneForAnotherUser(ctx, tx, userValidation.GetNewPhone(), userID); verr != nil {
		utils.SetSpanError(ctx, verr)
		slog.Error("user.confirm_phone_change.stage_error", "stage", "exists_phone_for_another_user", "err", verr)
		return derrors.Infra("Failed to check phone uniqueness", verr)
	} else if exist {
		return derrors.ErrPhoneAlreadyInUse
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
		return derrors.Infra("Failed to update validations", err)
	}

	// Status transition is handled centrally by applyStatusTransitionAfterContactChange in the public method.
	// Não realizar transições de status aqui para evitar duplicidade e inconsistências.

	err = us.repo.UpdateUserByID(ctx, tx, user)
	if err != nil {
		utils.SetSpanError(ctx, err)
		slog.Error("user.confirm_phone_change.stage_error", "stage", "update_user", "err", err)
		return derrors.Infra("Failed to update user", err)
	}

	err = us.globalService.CreateAudit(ctx, tx, globalmodel.TableUsers, "Alterada o telefone do usuário")
	if err != nil {
		utils.SetSpanError(ctx, err)
		slog.Error("user.confirm_phone_change.stage_error", "stage", "audit", "err", err)
		return derrors.Infra("Failed to create audit", err)
	}

	return
}
