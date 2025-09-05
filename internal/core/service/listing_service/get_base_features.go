package listingservices

import (
	"context"
	"log/slog"

	listingmodel "github.com/giulio-alfieri/toq_server/internal/core/model/listing_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (ls *listingService) GetBaseFeatures(ctx context.Context) (features []listingmodel.BaseFeatureInterface, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, utils.InternalError("")
	}
	defer spanEnd()

	tx, txErr := ls.gsi.StartTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		slog.Error("listing.get_base_features.tx_start_error", "err", txErr)
		return nil, utils.InternalError("")
	}
	defer func() {
		if err != nil {
			if rbErr := ls.gsi.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				slog.Error("listing.get_base_features.tx_rollback_error", "err", rbErr)
			}
		}
	}()

	features, err = ls.listingRepository.GetBaseFeatures(ctx, tx)
	if err != nil {
		utils.SetSpanError(ctx, err)
		return nil, utils.InternalError("")
	}

	if cmErr := ls.gsi.CommitTransaction(ctx, tx); cmErr != nil {
		utils.SetSpanError(ctx, cmErr)
		slog.Error("listing.get_base_features.tx_commit_error", "err", cmErr)
		return nil, utils.InternalError("")
	}

	return
}
