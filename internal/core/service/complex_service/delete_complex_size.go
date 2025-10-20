package complexservices

import (
	"context"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// DeleteComplexSize remove um tamanho cadastrado.
func (cs *complexService) DeleteComplexSize(ctx context.Context, id int64) error {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return utils.InternalError("")
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	if err := ensurePositiveID("id", id); err != nil {
		return err
	}

	tx, txErr := cs.gsi.StartTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("complex.size.delete.tx_start_error", "err", txErr)
		return utils.InternalError("")
	}
	success := false
	defer func() {
		if !success {
			if rbErr := cs.gsi.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				logger.Error("complex.size.delete.tx_rollback_error", "err", rbErr)
			}
		}
	}()

	rows, err := cs.complexRepository.DeleteComplexSize(ctx, tx, id)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("complex.size.delete.repo_error", "err", err, "id", id)
		return utils.InternalError("")
	}

	if rows == 0 {
		return utils.NotFoundError("complex_size")
	}

	if cmErr := cs.gsi.CommitTransaction(ctx, tx); cmErr != nil {
		utils.SetSpanError(ctx, cmErr)
		logger.Error("complex.size.delete.tx_commit_error", "err", cmErr)
		return utils.InternalError("")
	}

	success = true
	return nil
}
