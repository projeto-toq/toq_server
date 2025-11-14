package mysqllistingadapter

import (
	"context"
	"database/sql"
	"fmt"

	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// CreateWarehouseAdditionalFloor inserts a new warehouse additional floor record for a listing version.
func (la *ListingAdapter) CreateWarehouseAdditionalFloor(ctx context.Context, tx *sql.Tx, floor listingmodel.WarehouseAdditionalFloorInterface) (err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	statement := `INSERT INTO warehouse_additional_floors (listing_version_id, floor_name, floor_order, floor_height) VALUES (?, ?, ?, ?);`

	result, execErr := la.ExecContext(ctx, tx, "insert", statement,
		floor.ListingVersionID(),
		floor.FloorName(),
		floor.FloorOrder(),
		floor.FloorHeight(),
	)
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.listing.create_warehouse_additional_floor.exec_error", "error", execErr)
		return fmt.Errorf("exec create warehouse additional floor: %w", execErr)
	}

	id, lastErr := result.LastInsertId()
	if lastErr != nil {
		utils.SetSpanError(ctx, lastErr)
		logger.Error("mysql.listing.create_warehouse_additional_floor.last_insert_error", "error", lastErr)
		return fmt.Errorf("last insert id for warehouse additional floor: %w", lastErr)
	}

	floor.SetID(id)

	return nil
}
