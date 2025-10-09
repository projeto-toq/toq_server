package mysqllistingadapter

import (
	"context"
	"database/sql"
	"fmt"

	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (la *ListingAdapter) CreateExchangePlace(ctx context.Context, tx *sql.Tx, place listingmodel.ExchangePlaceInterface) (err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	sql := `INSERT INTO exchange_places (listing_id, neighborhood, city, state) VALUES (?, ?, ?, ?);`

	stmt, err := tx.PrepareContext(ctx, sql)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.listing.create_exchange_place.prepare_error", "error", err)
		return fmt.Errorf("prepare create exchange place: %w", err)
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, place.ListingID(), place.Neighborhood(), place.City(), place.State())
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.listing.create_exchange_place.exec_error", "error", err)
		return fmt.Errorf("exec create exchange place: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.listing.create_exchange_place.last_insert_error", "error", err)
		return fmt.Errorf("last insert id for exchange place: %w", err)
	}

	place.SetID(id)

	return nil
}
