package complexservices

import (
	"context"
	"database/sql"

	complexmodel "github.com/projeto-toq/toq_server/internal/core/model/complex_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// CreateComplexTower cria uma nova torre associada a um empreendimento.
func (cs *complexService) CreateComplexTower(ctx context.Context, input CreateComplexTowerInput) (complexmodel.ComplexTowerInterface, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, utils.InternalError("")
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	if err := ensurePositiveID("complexId", input.ComplexID); err != nil {
		return nil, err
	}

	if err := validateRequiredField("tower", input.Tower); err != nil {
		return nil, err
	}

	tx, txErr := cs.gsi.StartTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("complex.tower.create.tx_start_error", "err", txErr)
		return nil, utils.InternalError("")
	}
	success := false
	defer func() {
		if !success {
			if rbErr := cs.gsi.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				logger.Error("complex.tower.create.tx_rollback_error", "err", rbErr)
			}
		}
	}()

	if _, err = cs.complexRepository.GetComplexByID(ctx, tx, input.ComplexID); err != nil {
		if err == sql.ErrNoRows {
			return nil, utils.NotFoundError("complex")
		}
		utils.SetSpanError(ctx, err)
		logger.Error("complex.tower.create.parent_error", "err", err, "complex_id", input.ComplexID)
		return nil, utils.InternalError("")
	}

	tower := complexmodel.NewComplexTower()
	tower.SetComplexID(input.ComplexID)
	tower.SetTower(sanitizeString(input.Tower))
	tower.SetFloors(input.Floors)
	tower.SetTotalUnits(input.TotalUnits)
	tower.SetUnitsPerFloor(input.UnitsPerFloor)

	id, err := cs.complexRepository.CreateComplexTower(ctx, tx, tower)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("complex.tower.create.repo_error", "err", err, "complex_id", input.ComplexID)
		return nil, utils.InternalError("")
	}

	tower.SetID(id)

	if cmErr := cs.gsi.CommitTransaction(ctx, tx); cmErr != nil {
		utils.SetSpanError(ctx, cmErr)
		logger.Error("complex.tower.create.tx_commit_error", "err", cmErr)
		return nil, utils.InternalError("")
	}

	success = true
	return tower, nil
}
