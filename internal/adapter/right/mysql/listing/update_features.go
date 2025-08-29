package mysqllistingadapter

import (
	"context"
	"database/sql"

	listingmodel "github.com/giulio-alfieri/toq_server/internal/core/model/listing_model"
	
	

"github.com/giulio-alfieri/toq_server/internal/core/utils"
"errors"
)

func (la *ListingAdapter) UpdateFeatures(ctx context.Context, tx *sql.Tx, features []listingmodel.FeatureInterface) (err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	//check if there is any data to update
	if len(features) == 0 {
		return
	}

	// Remove all features from listing
	err = la.DeleteListingFeatures(ctx, tx, features[0].ListingID())
	if err != nil {
		//check if the error is not found, because it's ok if there is no row to delete
		if errors.Is(err, sql.ErrNoRows) {
			return
		}
	}

	// Insert the new features
	for _, feature := range features {
		err = la.CreateFeature(ctx, tx, feature)
		if err != nil {
			return
		}
	}

	return
}
