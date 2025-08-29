package listingservices

import (
	"context"

	listingmodel "github.com/giulio-alfieri/toq_server/internal/core/model/listing_model"
"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (ls *listingService) GetBaseFeatures(ctx context.Context) (features []listingmodel.BaseFeatureInterface, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	tx, err := ls.gsi.StartTransaction(ctx)
	if err != nil {
		return
	}

	features, err = ls.listingRepository.GetBaseFeatures(ctx, tx)
	if err != nil {
		ls.gsi.RollbackTransaction(ctx, tx)
		return
	}

	err = ls.gsi.CommitTransaction(ctx, tx)
	if err != nil {
		ls.gsi.RollbackTransaction(ctx, tx)
		return
	}

	return
}
