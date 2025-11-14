package mysqllistingadapter

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// CloneListingVersionSatellites copies all satellite entities (features, exchange_places,
// financing_blockers, guarantees) from sourceVersionID to targetVersionID.
// This is used when creating a draft version from an active version.
func (la *ListingAdapter) CloneListingVersionSatellites(ctx context.Context, tx *sql.Tx, sourceVersionID, targetVersionID int64) error {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	// Clone features
	featuresQuery := `
		INSERT INTO features (listing_version_id, feature_id, qty)
		SELECT ?, feature_id, qty
		FROM features
		WHERE listing_version_id = ?
	`
	if _, err := la.ExecContext(ctx, tx, "insert", featuresQuery, targetVersionID, sourceVersionID); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.listing.clone_satellites.features_error", "error", err, "source", sourceVersionID, "target", targetVersionID)
		return fmt.Errorf("clone features: %w", err)
	}

	// Clone exchange_places
	exchangeQuery := `
		INSERT INTO exchange_places (listing_version_id, neighborhood, city, state)
		SELECT ?, neighborhood, city, state
		FROM exchange_places
		WHERE listing_version_id = ?
	`
	if _, err := la.ExecContext(ctx, tx, "insert", exchangeQuery, targetVersionID, sourceVersionID); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.listing.clone_satellites.exchange_places_error", "error", err, "source", sourceVersionID, "target", targetVersionID)
		return fmt.Errorf("clone exchange places: %w", err)
	}

	// Clone financing_blockers
	blockersQuery := `
		INSERT INTO financing_blockers (listing_version_id, blocker)
		SELECT ?, blocker
		FROM financing_blockers
		WHERE listing_version_id = ?
	`
	if _, err := la.ExecContext(ctx, tx, "insert", blockersQuery, targetVersionID, sourceVersionID); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.listing.clone_satellites.financing_blockers_error", "error", err, "source", sourceVersionID, "target", targetVersionID)
		return fmt.Errorf("clone financing blockers: %w", err)
	}

	// Clone guarantees
	guaranteesQuery := `
		INSERT INTO guarantees (listing_version_id, priority, guarantee)
		SELECT ?, priority, guarantee
		FROM guarantees
		WHERE listing_version_id = ?
	`
	if _, err := la.ExecContext(ctx, tx, "insert", guaranteesQuery, targetVersionID, sourceVersionID); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.listing.clone_satellites.guarantees_error", "error", err, "source", sourceVersionID, "target", targetVersionID)
		return fmt.Errorf("clone guarantees: %w", err)
	}

	// Clone warehouse_additional_floors
	warehouseFloorsQuery := `
		INSERT INTO warehouse_additional_floors (listing_version_id, floor_name, floor_order, floor_height)
		SELECT ?, floor_name, floor_order, floor_height
		FROM warehouse_additional_floors
		WHERE listing_version_id = ?
	`
	if _, err := la.ExecContext(ctx, tx, "insert", warehouseFloorsQuery, targetVersionID, sourceVersionID); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.listing.clone_satellites.warehouse_floors_error", "error", err, "source", sourceVersionID, "target", targetVersionID)
		return fmt.Errorf("clone warehouse additional floors: %w", err)
	}

	logger.Info("listing.clone_satellites.success", "source_version_id", sourceVersionID, "target_version_id", targetVersionID)
	return nil
}
