package userservices

import (
	"context"
	"database/sql"
	"log/slog"

	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

// UpdateOptStatus consolidates opt-in/out behavior with audit and transactions
func (us *userService) UpdateOptStatus(ctx context.Context, opt bool) (err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return utils.InternalError("Failed to generate tracer")
	}
	defer spanEnd()

	userID, err := us.globalService.GetUserIDFromContext(ctx)
	if err != nil {
		return utils.AuthenticationError("")
	}

	tx, err := us.globalService.StartTransaction(ctx)
	if err != nil {
		slog.Error("user.update_opt_status.tx_start_error", "err", err)
		return utils.InternalError("Failed to start transaction")
	}
	defer func() {
		if err != nil {
			if rbErr := us.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				slog.Error("user.update_opt_status.tx_rollback_error", "err", rbErr)
			}
		}
	}()

	if err = us.updateOptStatus(ctx, tx, userID, opt); err != nil {
		return
	}

	if err = us.globalService.CommitTransaction(ctx, tx); err != nil {
		slog.Error("user.update_opt_status.tx_commit_error", "err", err)
		return utils.InternalError("Failed to commit transaction")
	}
	return
}

func (us *userService) updateOptStatus(ctx context.Context, tx *sql.Tx, userID int64, opt bool) (err error) {
	// Busca usuário atual para aplicar idempotência
	user, err := us.repo.GetUserByID(ctx, tx, userID)
	if err != nil {
		return
	}

	// Se já está no estado desejado, não faz nada (idempotente)
	if user.IsOptStatus() == opt {
		return nil
	}

	// Transição para opt-out: remover tokens antes de persistir
	if !opt {
		if err = us.repo.RemoveAllDeviceTokens(ctx, tx, userID); err != nil {
			return
		}
	}

	// Atualiza status e persiste
	user.SetOptStatus(opt)
	if err = us.repo.UpdateUserByID(ctx, tx, user); err != nil {
		return
	}

	// Auditoria conforme novo estado
	if opt {
		if err = us.globalService.CreateAudit(ctx, tx, globalmodel.TableUsers, "Usuário aceitou receber notificações"); err != nil {
			return
		}
	} else {
		if err = us.globalService.CreateAudit(ctx, tx, globalmodel.TableUsers, "Usuário rejeitou receber notificações"); err != nil {
			return
		}
	}
	return
}
