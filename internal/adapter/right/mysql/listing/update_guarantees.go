package mysqllistingadapter

import (
	"context"
	"database/sql"
	"fmt"

	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"

	"errors"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (la *ListingAdapter) UpdateGuarantees(ctx context.Context, tx *sql.Tx, listingVersionID int64, guarantees []listingmodel.GuaranteeInterface) (err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	// Remove all guarantees from listing
	err = la.DeleteListingGuarantees(ctx, tx, listingVersionID)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			utils.SetSpanError(ctx, err)
			logger.Error("mysql.listing.update_guarantees.delete_error", "error", err)
			return fmt.Errorf("delete listing guarantees: %w", err)
		}
	}

	if len(guarantees) == 0 {
		return nil
	}

	// Insert the new guarrantees
	for _, guarantee := range guarantees {
		guarantee.SetListingVersionID(listingVersionID)
		err = la.CreateGuarantee(ctx, tx, guarantee)
		if err != nil {
			utils.SetSpanError(ctx, err)
			logger.Error("mysql.listing.update_guarantees.create_error", "error", err)
			return fmt.Errorf("create guarantee: %w", err)
		}
	}

	return nil
}
