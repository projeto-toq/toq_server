package propertycoverageservice

import (
	"context"
	"database/sql"

	propertycoveragemodel "github.com/projeto-toq/toq_server/internal/core/model/property_coverage_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// GetComplexSizeDetail fetches a single size row by identifier.
func (s *propertyCoverageService) GetComplexSizeDetail(ctx context.Context, sizeID int64) (propertycoveragemodel.VerticalComplexSizeInterface, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, utils.InternalError("")
	}
	defer spanEnd()

	if err := ensurePositiveID("id", sizeID); err != nil {
		return nil, err
	}

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	tx, txErr := s.globalService.StartTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("property_coverage.size.detail.tx_start_error", "err", txErr)
		return nil, utils.InternalError("")
	}

	success := false
	defer func() {
		if !success {
			if rbErr := s.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				logger.Error("property_coverage.size.detail.tx_rollback_error", "err", rbErr)
			}
		}
	}()

	size, err := s.repository.GetVerticalComplexSize(ctx, tx, sizeID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, utils.NotFoundError("size")
		}
		utils.SetSpanError(ctx, err)
		logger.Error("property_coverage.size.detail.repo_error", "err", err, "id", sizeID)
		return nil, utils.InternalError("")
	}

	if commitErr := s.globalService.CommitTransaction(ctx, tx); commitErr != nil {
		utils.SetSpanError(ctx, commitErr)
		logger.Error("property_coverage.size.detail.tx_commit_error", "err", commitErr)
		return nil, utils.InternalError("")
	}

	success = true
	return size, nil
}
