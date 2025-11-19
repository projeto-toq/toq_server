package propertycoverageservice

import (
	"context"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// DeleteComplexZipCode removes a horizontal zip code entry.
func (s *propertyCoverageService) DeleteComplexZipCode(ctx context.Context, zipCodeID int64) error {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return utils.InternalError("")
	}
	defer spanEnd()

	if err := ensurePositiveID("id", zipCodeID); err != nil {
		return err
	}

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	tx, txErr := s.globalService.StartTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("property_coverage.zip.delete.tx_start_error", "err", txErr)
		return utils.InternalError("")
	}

	success := false
	defer func() {
		if !success {
			if rbErr := s.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				logger.Error("property_coverage.zip.delete.tx_rollback_error", "err", rbErr)
			}
		}
	}()

	rows, err := s.repository.DeleteHorizontalComplexZipCode(ctx, tx, zipCodeID)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("property_coverage.zip.delete.repo_error", "err", err, "id", zipCodeID)
		return utils.InternalError("")
	}

	if rows == 0 {
		return utils.NotFoundError("zipCode")
	}

	if commitErr := s.globalService.CommitTransaction(ctx, tx); commitErr != nil {
		utils.SetSpanError(ctx, commitErr)
		logger.Error("property_coverage.zip.delete.tx_commit_error", "err", commitErr)
		return utils.InternalError("")
	}

	success = true
	return nil
}
