package complexservices

import (
	"context"
	"database/sql"

	complexmodel "github.com/projeto-toq/toq_server/internal/core/model/complex_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// UpdateComplexTower atualiza dados de uma torre existente.
func (cs *complexService) UpdateComplexTower(ctx context.Context, input UpdateComplexTowerInput) (complexmodel.ComplexTowerInterface, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, utils.InternalError("")
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	if err := ensurePositiveID("id", input.ID); err != nil {
		return nil, err
	}

	if err := ensurePositiveID("complexId", input.ComplexID); err != nil {
		return nil, err
	}

	if err := validateRequiredField("tower", input.Tower); err != nil {
		return nil, err
	}

	tx, txErr := cs.gsi.StartTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("complex.tower.update.tx_start_error", "err", txErr)
		return nil, utils.InternalError("")
	}
	success := false
	defer func() {
		if !success {
			if rbErr := cs.gsi.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				logger.Error("complex.tower.update.tx_rollback_error", "err", rbErr)
			}
		}
	}()

	if _, err = cs.complexRepository.GetComplexByID(ctx, tx, input.ComplexID); err != nil {
		if err == sql.ErrNoRows {
			return nil, utils.NotFoundError("complex")
		}
		utils.SetSpanError(ctx, err)
		logger.Error("complex.tower.update.parent_error", "err", err, "complex_id", input.ComplexID)
		return nil, utils.InternalError("")
	}

	tower, err := cs.complexRepository.GetComplexTowerByID(ctx, tx, input.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, utils.NotFoundError("complex_tower")
		}
		utils.SetSpanError(ctx, err)
		logger.Error("complex.tower.update.get_error", "err", err, "id", input.ID)
		return nil, utils.InternalError("")
	}

	tower.SetTower(sanitizeString(input.Tower))
	tower.SetFloors(input.Floors)
	tower.SetTotalUnits(input.TotalUnits)
	tower.SetUnitsPerFloor(input.UnitsPerFloor)

	rows, err := cs.complexRepository.UpdateComplexTower(ctx, tx, tower)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("complex.tower.update.repo_error", "err", err, "id", input.ID)
		return nil, utils.InternalError("")
	}

	if rows == 0 {
		return nil, utils.NotFoundError("complex_tower")
	}

	if cmErr := cs.gsi.CommitTransaction(ctx, tx); cmErr != nil {
		utils.SetSpanError(ctx, cmErr)
		logger.Error("complex.tower.update.tx_commit_error", "err", cmErr)
		return nil, utils.InternalError("")
	}

	success = true
	return tower, nil
}
