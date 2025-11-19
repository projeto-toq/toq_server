package propertycoverageservice

import (
	"context"
	"database/sql"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// DeleteComplex removes a managed coverage entry from the database.
func (s *propertyCoverageService) DeleteComplex(ctx context.Context, input DeleteComplexInput) error {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return utils.InternalError("")
	}
	defer spanEnd()

	if err := ensurePositiveID("id", input.ID); err != nil {
		return err
	}

	if err := validateCoverageKind(input.Kind); err != nil {
		return err
	}

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	tx, txErr := s.globalService.StartTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("property_coverage.delete.tx_start_error", "err", txErr)
		return utils.InternalError("")
	}

	success := false
	defer func() {
		if !success {
			if rbErr := s.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				logger.Error("property_coverage.delete.tx_rollback_error", "err", rbErr)
			}
		}
	}()

	if _, err := s.repository.GetManagedComplex(ctx, tx, input.ID, input.Kind); err != nil {
		if err == sql.ErrNoRows {
			return utils.NotFoundError("complex")
		}
		utils.SetSpanError(ctx, err)
		logger.Error("property_coverage.delete.get_error", "err", err, "id", input.ID)
		return utils.InternalError("")
	}

	rows, err := s.repository.DeleteManagedComplex(ctx, tx, input.ID, input.Kind)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("property_coverage.delete.repo_error", "err", err)
		return utils.InternalError("")
	}

	if rows == 0 {
		return utils.NotFoundError("complex")
	}

	if commitErr := s.globalService.CommitTransaction(ctx, tx); commitErr != nil {
		utils.SetSpanError(ctx, commitErr)
		logger.Error("property_coverage.delete.tx_commit_error", "err", commitErr)
		return utils.InternalError("")
	}

	success = true
	return nil
}
