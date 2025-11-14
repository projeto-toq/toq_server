package mysqllistingadapter

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// UpdateWarehouseAdditionalFloors replaces all warehouse additional floors for a listing version.
// Deletes existing floors and inserts new ones in a single transaction.
func (la *ListingAdapter) UpdateWarehouseAdditionalFloors(ctx context.Context, tx *sql.Tx, listingVersionID int64, floors []listingmodel.WarehouseAdditionalFloorInterface) (err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	// Remove all warehouse additional floors from listing version
	err = la.DeleteListingWarehouseAdditionalFloors(ctx, tx, listingVersionID)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			utils.SetSpanError(ctx, err)
			logger.Error("mysql.listing.update_warehouse_additional_floors.delete_error", "error", err)
			return fmt.Errorf("delete listing warehouse additional floors: %w", err)
		}
	}

	if len(floors) == 0 {
		return nil
	}

	// Insert the new warehouse additional floors
	for _, floor := range floors {
		floor.SetListingVersionID(listingVersionID)
		err = la.CreateWarehouseAdditionalFloor(ctx, tx, floor)
		if err != nil {
			utils.SetSpanError(ctx, err)
			logger.Error("mysql.listing.update_warehouse_additional_floors.create_error", "error", err)
			return fmt.Errorf("create warehouse additional floor: %w", err)
		}
	}

	return nil
}
