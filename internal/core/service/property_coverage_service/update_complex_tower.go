package propertycoverageservice

import (
	"context"
	"database/sql"

	propertycoveragemodel "github.com/projeto-toq/toq_server/internal/core/model/property_coverage_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// UpdateComplexTower updates tower metadata for a vertical complex.
func (s *propertyCoverageService) UpdateComplexTower(ctx context.Context, input UpdateComplexTowerInput) (propertycoveragemodel.VerticalComplexTowerInterface, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, utils.InternalError("")
	}
	defer spanEnd()

	if err := ensurePositiveID("id", input.ID); err != nil {
		return nil, err
	}

	if err := ensurePositiveID("verticalComplexId", input.VerticalComplexID); err != nil {
		return nil, err
	}

	if err := validateRequiredField("tower", input.Tower); err != nil {
		return nil, err
	}

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	tx, txErr := s.globalService.StartTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("property_coverage.tower.update.tx_start_error", "err", txErr)
		return nil, utils.InternalError("")
	}

	success := false
	defer func() {
		if !success {
			if rbErr := s.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				logger.Error("property_coverage.tower.update.tx_rollback_error", "err", rbErr)
			}
		}
	}()

	if _, err := s.repository.GetManagedComplex(ctx, tx, input.VerticalComplexID, propertycoveragemodel.CoverageKindVertical); err != nil {
		if err == sql.ErrNoRows {
			return nil, utils.NotFoundError("complex")
		}
		utils.SetSpanError(ctx, err)
		logger.Error("property_coverage.tower.update.parent_error", "err", err, "complex_id", input.VerticalComplexID)
		return nil, utils.InternalError("")
	}

	tower, err := s.repository.GetVerticalComplexTower(ctx, tx, input.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, utils.NotFoundError("tower")
		}
		utils.SetSpanError(ctx, err)
		logger.Error("property_coverage.tower.update.get_error", "err", err, "tower_id", input.ID)
		return nil, utils.InternalError("")
	}

	tower.SetVerticalComplexID(input.VerticalComplexID)
	tower.SetTower(sanitizeString(input.Tower))
	tower.SetFloors(input.Floors)
	tower.SetTotalUnits(input.TotalUnits)
	tower.SetUnitsPerFloor(input.UnitsPerFloor)

	rows, err := s.repository.UpdateVerticalComplexTower(ctx, tx, tower)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("property_coverage.tower.update.repo_error", "err", err)
		return nil, utils.InternalError("")
	}

	if rows == 0 {
		return nil, utils.NotFoundError("tower")
	}

	if commitErr := s.globalService.CommitTransaction(ctx, tx); commitErr != nil {
		utils.SetSpanError(ctx, commitErr)
		logger.Error("property_coverage.tower.update.tx_commit_error", "err", commitErr)
		return nil, utils.InternalError("")
	}

	success = true
	return tower, nil
}
