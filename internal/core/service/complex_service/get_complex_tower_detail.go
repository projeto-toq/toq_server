package complexservices

import (
	"context"
	"database/sql"

	complexmodel "github.com/projeto-toq/toq_server/internal/core/model/complex_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// GetComplexTowerDetail returns a single complex tower by its identifier.
func (cs *complexService) GetComplexTowerDetail(ctx context.Context, id int64) (complexmodel.ComplexTowerInterface, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, utils.InternalError("")
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	if err := ensurePositiveID("id", id); err != nil {
		return nil, err
	}

	tx, txErr := cs.gsi.StartTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("complex.tower.detail.tx_start_error", "err", txErr)
		return nil, utils.InternalError("")
	}
	success := false
	defer func() {
		if !success {
			if rbErr := cs.gsi.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				logger.Error("complex.tower.detail.tx_rollback_error", "err", rbErr)
			}
		}
	}()

	tower, err := cs.complexRepository.GetComplexTowerByID(ctx, tx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, utils.NotFoundError("complex_tower")
		}
		utils.SetSpanError(ctx, err)
		logger.Error("complex.tower.detail.repo_error", "err", err, "id", id)
		return nil, utils.InternalError("")
	}

	if cmErr := cs.gsi.CommitTransaction(ctx, tx); cmErr != nil {
		utils.SetSpanError(ctx, cmErr)
		logger.Error("complex.tower.detail.tx_commit_error", "err", cmErr)
		return nil, utils.InternalError("")
	}

	success = true
	return tower, nil
}
