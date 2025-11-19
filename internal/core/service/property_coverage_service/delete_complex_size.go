package propertycoverageservice

import (
	"context"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// DeleteComplexSize removes a size entry by identifier.
func (s *propertyCoverageService) DeleteComplexSize(ctx context.Context, sizeID int64) error {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return utils.InternalError("")
	}
	defer spanEnd()

	if err := ensurePositiveID("id", sizeID); err != nil {
		return err
	}

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	tx, txErr := s.globalService.StartTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("property_coverage.size.delete.tx_start_error", "err", txErr)
		return utils.InternalError("")
	}

	success := false
	defer func() {
		if !success {
			if rbErr := s.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				logger.Error("property_coverage.size.delete.tx_rollback_error", "err", rbErr)
			}
		}
	}()

	rows, err := s.repository.DeleteVerticalComplexSize(ctx, tx, sizeID)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("property_coverage.size.delete.repo_error", "err", err, "id", sizeID)
		return utils.InternalError("")
	}

	if rows == 0 {
		return utils.NotFoundError("size")
	}

	if commitErr := s.globalService.CommitTransaction(ctx, tx); commitErr != nil {
		utils.SetSpanError(ctx, commitErr)
		logger.Error("property_coverage.size.delete.tx_commit_error", "err", commitErr)
		return utils.InternalError("")
	}

	success = true
	return nil
}
