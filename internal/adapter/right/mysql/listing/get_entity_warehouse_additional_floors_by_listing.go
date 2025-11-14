package mysqllistingadapter

import (
	"context"
	"database/sql"
	"fmt"

	listingentity "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/listing/entity"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// GetEntityWarehouseAdditionalFloorsByListing retrieves all warehouse additional floors for a listing version.
func (la *ListingAdapter) GetEntityWarehouseAdditionalFloorsByListing(ctx context.Context, tx *sql.Tx, listingVersionID int64) (floors []listingentity.EntityWarehouseAdditionalFloor, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `SELECT id, listing_version_id, floor_name, floor_order, floor_height 
	          FROM warehouse_additional_floors 
	          WHERE listing_version_id = ? 
	          ORDER BY floor_order ASC;`

	rows, queryErr := la.QueryContext(ctx, tx, "select", query, listingVersionID)
	if queryErr != nil {
		utils.SetSpanError(ctx, queryErr)
		logger.Error("mysql.listing.get_entity_warehouse_additional_floors.query_error", "error", queryErr)
		return nil, fmt.Errorf("query warehouse additional floors by listing: %w", queryErr)
	}
	defer rows.Close()

	for rows.Next() {
		floor := listingentity.EntityWarehouseAdditionalFloor{}
		err = rows.Scan(
			&floor.ID,
			&floor.ListingVersionID,
			&floor.FloorName,
			&floor.FloorOrder,
			&floor.FloorHeight,
		)
		if err != nil {
			utils.SetSpanError(ctx, err)
			logger.Error("mysql.listing.get_entity_warehouse_additional_floors.scan_error", "error", err)
			return nil, fmt.Errorf("scan warehouse additional floor row: %w", err)
		}

		floors = append(floors, floor)
	}

	if err = rows.Err(); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.listing.get_entity_warehouse_additional_floors.rows_error", "error", err)
		return nil, fmt.Errorf("rows iteration for warehouse additional floors: %w", err)
	}

	return floors, nil
}
