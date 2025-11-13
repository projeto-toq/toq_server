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

	statement := `INSERT INTO exchange_places (listing_version_id, neighborhood, city, state) VALUES (?, ?, ?, ?);`

	result, execErr := la.ExecContext(ctx, tx, "insert", statement, place.ListingVersionID(), place.Neighborhood(), place.City(), place.State())
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.listing.create_exchange_place.exec_error", "error", execErr)
		return fmt.Errorf("exec create exchange place: %w", execErr)
	}

	id, lastErr := result.LastInsertId()
	if lastErr != nil {
		utils.SetSpanError(ctx, lastErr)
		logger.Error("mysql.listing.create_exchange_place.last_insert_error", "error", lastErr)
		return fmt.Errorf("last insert id for exchange place: %w", lastErr)
	}

	place.SetID(id)

	return nil
}
