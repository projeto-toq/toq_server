package mysqllistingadapter

import (
	"context"
	"database/sql"
	"fmt"

	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"

	"errors"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (la *ListingAdapter) UpdateFinancingBlockers(ctx context.Context, tx *sql.Tx, blockers []listingmodel.FinancingBlockerInterface) (err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	//check if there is any data to update
	if len(blockers) == 0 {
		return
	}

	// Remove all blocker from listing
	err = la.DeleteListingFinancingBlockers(ctx, tx, blockers[0].ListingID())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.listing.update_financing_blockers.delete_error", "error", err)
		return fmt.Errorf("delete listing financing blockers: %w", err)
	}

	// Insert the new blokers
	for _, blocker := range blockers {
		err = la.CreateFinancingBlocker(ctx, tx, blocker)
		if err != nil {
			utils.SetSpanError(ctx, err)
			logger.Error("mysql.listing.update_financing_blockers.create_error", "error", err)
			return fmt.Errorf("create financing blocker: %w", err)
		}
	}

	return nil
}
