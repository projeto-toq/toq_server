package mysqllistingadapter

import (
	"context"
	"database/sql"

	listingmodel "github.com/giulio-alfieri/toq_server/internal/core/model/listing_model"
	
	

"github.com/giulio-alfieri/toq_server/internal/core/utils"
"errors"
)

func (la *ListingAdapter) UpdateExchangePlaces(ctx context.Context, tx *sql.Tx, places []listingmodel.ExchangePlaceInterface) (err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	//check if there is any data to update
	if len(places) == 0 {
		return
	}

	//remove all exchange places from listing
	err = la.DeleteListingExchangePlaces(ctx, tx, places[0].ListingID())
	if err != nil {
		//check if the error is not found, because it's ok if there is no row to delete
		if errors.Is(err, sql.ErrNoRows) {
			return
		}
	}

	//insert the new exchange places
	for _, place := range places {
		err = la.CreateExchangePlace(ctx, tx, place)
		if err != nil {
			return
		}
	}

	return
}
