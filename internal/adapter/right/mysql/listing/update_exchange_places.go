package mysqllistingadapter

import (
	"context"
	"database/sql"
	"fmt"

	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"

	"errors"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (la *ListingAdapter) UpdateExchangePlaces(ctx context.Context, tx *sql.Tx, listingID int64, places []listingmodel.ExchangePlaceInterface) (err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	//remove all exchange places from listing
	err = la.DeleteListingExchangePlaces(ctx, tx, listingID)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			utils.SetSpanError(ctx, err)
			logger.Error("mysql.listing.update_exchange_places.delete_error", "error", err)
			return fmt.Errorf("delete listing exchange places: %w", err)
		}
	}

	if len(places) == 0 {
		return nil
	}

	//insert the new exchange places
	for _, place := range places {
		place.SetListingID(listingID)
		err = la.CreateExchangePlace(ctx, tx, place)
		if err != nil {
			utils.SetSpanError(ctx, err)
			logger.Error("mysql.listing.update_exchange_places.create_error", "error", err)
			return fmt.Errorf("create exchange place: %w", err)
		}
	}

	return nil
}
