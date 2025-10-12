package mysqllistingadapter

import (
	"context"
	"database/sql"
	"fmt"

	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"

	"errors"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (la *ListingAdapter) UpdateFeatures(ctx context.Context, tx *sql.Tx, listingID int64, features []listingmodel.FeatureInterface) (err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	// Remove all features from listing
	err = la.DeleteListingFeatures(ctx, tx, listingID)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			utils.SetSpanError(ctx, err)
			logger.Error("mysql.listing.update_features.delete_error", "error", err)
			return fmt.Errorf("delete listing features: %w", err)
		}
	}

	if len(features) == 0 {
		return nil
	}

	// Insert the new features
	for _, feature := range features {
		feature.SetListingID(listingID)
		err = la.CreateFeature(ctx, tx, feature)
		if err != nil {
			utils.SetSpanError(ctx, err)
			logger.Error("mysql.listing.update_features.create_error", "error", err)
			return fmt.Errorf("create feature: %w", err)
		}
	}

	return nil
}
