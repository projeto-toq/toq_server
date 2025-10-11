package complexservices

import (
	"context"
	"database/sql"

	globalmodel "github.com/projeto-toq/toq_server/internal/core/model/global_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (cs *complexService) GetOptions(ctx context.Context, zipCode string, number string) (propertyTypes globalmodel.PropertyType, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return 0, utils.InternalError("")
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	tx, txErr := cs.gsi.StartTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("complex.get_options.tx_start_error", "err", txErr)
		return 0, utils.InternalError("")
	}
	defer func() {
		if err != nil {
			if rbErr := cs.gsi.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				logger.Error("complex.get_options.tx_rollback_error", "err", rbErr)
			}
		}
	}()

	propertyTypes, err = cs.getOptions(ctx, tx, zipCode, number)
	if err != nil {
		return 0, err
	}

	if cmErr := cs.gsi.CommitTransaction(ctx, tx); cmErr != nil {
		utils.SetSpanError(ctx, cmErr)
		logger.Error("complex.get_options.tx_commit_error", "err", cmErr)
		return 0, utils.InternalError("")
	}

	return
}

func (cs *complexService) getOptions(ctx context.Context, tx *sql.Tx, zipCode string, number string) (propertyTypes globalmodel.PropertyType, err error) {
	callhorizontal := false
	complex, err := cs.complexRepository.GetVerticalComplex(ctx, tx, zipCode, number)
	if err != nil {
		if err == sql.ErrNoRows {
			callhorizontal = true
		} else {
			utils.SetSpanError(ctx, err)
			return 0, utils.InternalError("")
		}
	}
	if callhorizontal {
		complex, err = cs.complexRepository.GetHorizontalComplex(ctx, tx, zipCode)
		if err != nil {
			if err == sql.ErrNoRows {
				return propertyTypes, utils.NotFoundError("complex")
			}
			utils.SetSpanError(ctx, err)
			return 0, utils.InternalError("")
		}
	}

	propertyTypes = complex.GetPropertyType()
	return

}
