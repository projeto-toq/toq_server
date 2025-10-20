package complexservices

import (
	"context"
	"database/sql"

	complexmodel "github.com/projeto-toq/toq_server/internal/core/model/complex_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// GetComplexSizeDetail returns a single complex size by its identifier.
func (cs *complexService) GetComplexSizeDetail(ctx context.Context, id int64) (complexmodel.ComplexSizeInterface, error) {
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
		logger.Error("complex.size.detail.tx_start_error", "err", txErr)
		return nil, utils.InternalError("")
	}
	success := false
	defer func() {
		if !success {
			if rbErr := cs.gsi.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				logger.Error("complex.size.detail.tx_rollback_error", "err", rbErr)
			}
		}
	}()

	size, err := cs.complexRepository.GetComplexSizeByID(ctx, tx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, utils.NotFoundError("complex_size")
		}
		utils.SetSpanError(ctx, err)
		logger.Error("complex.size.detail.repo_error", "err", err, "id", id)
		return nil, utils.InternalError("")
	}

	if cmErr := cs.gsi.CommitTransaction(ctx, tx); cmErr != nil {
		utils.SetSpanError(ctx, cmErr)
		logger.Error("complex.size.detail.tx_commit_error", "err", cmErr)
		return nil, utils.InternalError("")
	}

	success = true
	return size, nil
}
