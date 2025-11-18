package mysqllistingadapter

import (
	"context"
	"database/sql"
	"fmt"

	listingentity "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/listing/entity"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// GetEntityFeaturesByListing retrieves all features associated with a specific listing version
//
// This function queries the features table to fetch all feature records linked to a given
// listing_version_id. Each feature includes a reference to base_features (feature_id) and
// a quantity value indicating how many instances of that feature exist.
//
// The function ensures:
//   - Only features linked to the specified listing version are returned
//   - Proper error handling with span marking for observability
//   - Efficient scanning of result rows into entity structs
//
// Database Schema:
//   - Table: features
//   - Columns: id (INT), listing_version_id (INT), feature_id (INT), qty (TINYINT)
//   - Foreign Key: listing_version_id → listing_versions.id (CASCADE DELETE)
//   - Foreign Key: feature_id → base_features.id (NO ACTION)
//
// Parameters:
//   - ctx: Context for tracing, cancellation, and logging
//   - tx: Database transaction (can be nil for standalone queries)
//   - listingVersionID: Unique identifier of the listing version (listing_versions.id)
//
// Returns:
//   - features: Slice of EntityFeature structs containing all matching features
//   - error: Database query/scan errors, or nil on success
//
// Notes:
//   - Returns empty slice (not error) if no features exist for the listing version
//   - Uses InstrumentedAdapter for automatic metrics and tracing
//   - Column name is 'qty' in database (not 'quantity')
func (la *ListingAdapter) GetEntityFeaturesByListing(ctx context.Context, tx *sql.Tx, listingVersionID int64) (features []listingentity.EntityFeature, err error) {
	// Initialize tracing for observability (metrics + distributed tracing)
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	// Attach logger to context to ensure request_id/trace_id propagation
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	// Query features table - note column name is 'qty' not 'quantity'
	// This matches the actual database schema in db_creation.sql
	query := `SELECT id, listing_version_id, feature_id, qty FROM features WHERE listing_version_id = ?`

	// Execute query using instrumented adapter (auto-generates metrics + tracing)
	rows, queryErr := la.QueryContext(ctx, tx, "select", query, listingVersionID)
	if queryErr != nil {
		// Mark span as error for distributed tracing analysis
		utils.SetSpanError(ctx, queryErr)
		logger.Error("mysql.listing.get_entity_features.query_error", "error", queryErr)
		return nil, fmt.Errorf("query features by listing: %w", queryErr)
	}
	defer rows.Close()

	// Scan all feature rows into entity structs
	for rows.Next() {
		feature := listingentity.EntityFeature{}
		err = rows.Scan(
			&feature.ID,
			&feature.ListingVersionID,
			&feature.FeatureID,
			&feature.Quantity,
		)
		if err != nil {
			utils.SetSpanError(ctx, err)
			logger.Error("mysql.listing.get_entity_features.scan_error", "error", err)
			return nil, fmt.Errorf("scan feature row: %w", err)
		}

		features = append(features, feature)
	}

	// Check for errors that occurred during iteration
	if err = rows.Err(); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.listing.get_entity_features.rows_error", "error", err)
		return nil, fmt.Errorf("rows iteration for features: %w", err)
	}

	return features, nil
}
