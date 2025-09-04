package complexservices

import (
	"context"
	"database/sql"

	"log/slog"

	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (cs *complexService) GetOptions(ctx context.Context, zipCode string, number string) (propertyTypes globalmodel.PropertyType, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return propertyTypes, err
	}
	defer spanEnd()

	tx, txErr := cs.gsi.StartTransaction(ctx)
	if txErr != nil {
		slog.Error("complex.get_options.tx_start_error", "err", txErr)
		return propertyTypes, utils.InternalError("Failed to start transaction")
	}
	defer func() {
		if err != nil {
			if rbErr := cs.gsi.RollbackTransaction(ctx, tx); rbErr != nil {
				slog.Error("complex.get_options.tx_rollback_error", "err", rbErr)
			}
		}
	}()

	propertyTypes, err = cs.getOptions(ctx, tx, zipCode, number)
	if err != nil {
		return
	}

	if cmErr := cs.gsi.CommitTransaction(ctx, tx); cmErr != nil {
		slog.Error("complex.get_options.tx_commit_error", "err", cmErr)
		return propertyTypes, utils.InternalError("Failed to commit transaction")
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
			return
		}
	}
	if callhorizontal {
		complex, err = cs.complexRepository.GetHorizontalComplex(ctx, tx, zipCode)
		if err != nil {
			if err == sql.ErrNoRows {
				return propertyTypes, utils.ValidationError("area_not_covered", "Area not covered yet")
			} else {
				return
			}
		}
	}

	propertyTypes = complex.GetPropertyType()
	return

}
