package propertycoverageservice

import (
	"context"
	"database/sql"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// DeleteComplexTower removes a tower row from a vertical complex.
func (s *propertyCoverageService) DeleteComplexTower(ctx context.Context, towerID int64) error {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return utils.InternalError("")
	}
	defer spanEnd()

	if err := ensurePositiveID("id", towerID); err != nil {
		return err
	}

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	tx, txErr := s.globalService.StartTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("property_coverage.tower.delete.tx_start_error", "err", txErr)
		return utils.InternalError("")
	}

	success := false
	defer func() {
		if !success {
			if rbErr := s.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				logger.Error("property_coverage.tower.delete.tx_rollback_error", "err", rbErr)
			}
		}
	}()

	if _, err := s.repository.GetVerticalComplexTower(ctx, tx, towerID); err != nil {
		if err == sql.ErrNoRows {
			return utils.NotFoundError("tower")
		}
		utils.SetSpanError(ctx, err)
		logger.Error("property_coverage.tower.delete.get_error", "err", err, "tower_id", towerID)
		return utils.InternalError("")
	}

	rows, err := s.repository.DeleteVerticalComplexTower(ctx, tx, towerID)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("property_coverage.tower.delete.repo_error", "err", err)
		return utils.InternalError("")
	}

	if rows == 0 {
		return utils.NotFoundError("tower")
	}

	if commitErr := s.globalService.CommitTransaction(ctx, tx); commitErr != nil {
		utils.SetSpanError(ctx, commitErr)
		logger.Error("property_coverage.tower.delete.tx_commit_error", "err", commitErr)
		return utils.InternalError("")
	}

	success = true
	return nil
}
