package propertycoverageservice

import (
	"context"
	"database/sql"

	propertycoveragemodel "github.com/projeto-toq/toq_server/internal/core/model/property_coverage_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// CreateComplexTower registers a tower under a vertical complex.
func (s *propertyCoverageService) CreateComplexTower(ctx context.Context, input CreateComplexTowerInput) (propertycoveragemodel.VerticalComplexTowerInterface, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, utils.InternalError("")
	}
	defer spanEnd()

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
		logger.Error("property_coverage.tower.create.tx_start_error", "err", txErr)
		return nil, utils.InternalError("")
	}

	success := false
	defer func() {
		if !success {
			if rbErr := s.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				logger.Error("property_coverage.tower.create.tx_rollback_error", "err", rbErr)
			}
		}
	}()

	if _, err := s.repository.GetManagedComplex(ctx, tx, input.VerticalComplexID, propertycoveragemodel.CoverageKindVertical); err != nil {
		if err == sql.ErrNoRows {
			return nil, utils.NotFoundError("complex")
		}
		utils.SetSpanError(ctx, err)
		logger.Error("property_coverage.tower.create.parent_error", "err", err, "complex_id", input.VerticalComplexID)
		return nil, utils.InternalError("")
	}

	tower := propertycoveragemodel.NewVerticalComplexTower()
	tower.SetVerticalComplexID(input.VerticalComplexID)
	tower.SetTower(sanitizeString(input.Tower))
	tower.SetFloors(input.Floors)
	tower.SetTotalUnits(input.TotalUnits)
	tower.SetUnitsPerFloor(input.UnitsPerFloor)

	id, err := s.repository.CreateVerticalComplexTower(ctx, tx, tower)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("property_coverage.tower.create.repo_error", "err", err)
		return nil, utils.InternalError("")
	}

	tower.SetID(id)

	if commitErr := s.globalService.CommitTransaction(ctx, tx); commitErr != nil {
		utils.SetSpanError(ctx, commitErr)
		logger.Error("property_coverage.tower.create.tx_commit_error", "err", commitErr)
		return nil, utils.InternalError("")
	}

	success = true
	return tower, nil
}
