package mysqllistingadapter

import (
	"context"
	"database/sql"
	"fmt"

	listingmodel "github.com/giulio-alfieri/toq_server/internal/core/model/listing_model"

	"errors"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (la *ListingAdapter) UpdateExchangePlaces(ctx context.Context, tx *sql.Tx, places []listingmodel.ExchangePlaceInterface) (err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	//check if there is any data to update
	if len(places) == 0 {
		return
	}

	//remove all exchange places from listing
	err = la.DeleteListingExchangePlaces(ctx, tx, places[0].ListingID())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.listing.update_exchange_places.delete_error", "error", err)
		return fmt.Errorf("delete listing exchange places: %w", err)
	}

	//insert the new exchange places
	for _, place := range places {
		err = la.CreateExchangePlace(ctx, tx, place)
		if err != nil {
			utils.SetSpanError(ctx, err)
			logger.Error("mysql.listing.update_exchange_places.create_error", "error", err)
			return fmt.Errorf("create exchange place: %w", err)
		}
	}

	return nil
}
