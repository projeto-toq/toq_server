package mysqllistingadapter

import (
	"context"
	"database/sql"
	"fmt"

	listingconverters "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/listing/converters"
	listingentity "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/listing/entity"
	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

const listingSelectColumns = `
		lv.id,
		lv.listing_identity_id,
		li.listing_uuid,
		li.active_version_id,
		lv.user_id,
		lv.code,
		lv.version,
		lv.status,
		lv.zip_code,
		lv.street,
		lv.number,
		lv.complement,
		lv.complex,
		lv.neighborhood,
		lv.city,
		lv.state,
		lv.title,
		lv.type,
		lv.owner,
		lv.land_size,
		lv.corner,
		lv.non_buildable,
		lv.buildable,
		lv.delivered,
		lv.who_lives,
		lv.description,
		lv.transaction,
		lv.sell_net,
		lv.rent_net,
		lv.condominium,
		lv.annual_tax,
		lv.monthly_tax,
		lv.annual_ground_rent,
		lv.monthly_ground_rent,
		lv.exchange,
		lv.exchange_perc,
		lv.installment,
		lv.financing,
		lv.visit,
		lv.tenant_name,
		lv.tenant_email,
		lv.tenant_phone,
		lv.accompanying,
		lv.deleted,
		lv.completion_forecast,
		lv.land_block,
		lv.land_lot,
		lv.land_front,
		lv.land_side,
		lv.land_back,
		lv.land_terrain_type,
		lv.has_kmz,
		lv.kmz_file,
		lv.building_floors,
		lv.unit_tower,
		lv.unit_floor,
		lv.unit_number,
		lv.warehouse_manufacturing_area,
		lv.warehouse_sector,
		lv.warehouse_has_primary_cabin,
		lv.warehouse_cabin_kva,
		lv.warehouse_ground_floor,
		lv.warehouse_floor_resistance,
		lv.warehouse_zoning,
		lv.warehouse_has_office_area,
		lv.warehouse_office_area,
		lv.store_has_mezzanine,
		lv.store_mezzanine_area`

type rowScanner interface {
	Scan(dest ...any) error
}

func scanListingEntity(row rowScanner) (listingentity.ListingEntity, error) {
	entity := listingentity.ListingEntity{}
	if err := row.Scan(
		&entity.ID,
		&entity.ListingIdentityID,
		&entity.ListingUUID,
		&entity.ActiveVersionID,
		&entity.UserID,
		&entity.Code,
		&entity.Version,
		&entity.Status,
		&entity.ZipCode,
		&entity.Street,
		&entity.Number,
		&entity.Complement,
		&entity.Complex,
		&entity.Neighborhood,
		&entity.City,
		&entity.State,
		&entity.Title,
		&entity.ListingType,
		&entity.Owner,
		&entity.LandSize,
		&entity.Corner,
		&entity.NonBuildable,
		&entity.Buildable,
		&entity.Delivered,
		&entity.WhoLives,
		&entity.Description,
		&entity.Transaction,
		&entity.SellNet,
		&entity.RentNet,
		&entity.Condominium,
		&entity.AnnualTax,
		&entity.MonthlyTax,
		&entity.AnnualGroundRent,
		&entity.MonthlyGroundRent,
		&entity.Exchange,
		&entity.ExchangePercentual,
		&entity.Installment,
		&entity.Financing,
		&entity.Visit,
		&entity.TenantName,
		&entity.TenantEmail,
		&entity.TenantPhone,
		&entity.Accompanying,
		&entity.Deleted,
		&entity.CompletionForecast,
		&entity.LandBlock,
		&entity.LandLot,
		&entity.LandFront,
		&entity.LandSide,
		&entity.LandBack,
		&entity.LandTerrainType,
		&entity.HasKmz,
		&entity.KmzFile,
		&entity.BuildingFloors,
		&entity.UnitTower,
		&entity.UnitFloor,
		&entity.UnitNumber,
		&entity.WarehouseManufacturingArea,
		&entity.WarehouseSector,
		&entity.WarehouseHasPrimaryCabin,
		&entity.WarehouseCabinKva,
		&entity.WarehouseGroundFloor,
		&entity.WarehouseFloorResistance,
		&entity.WarehouseZoning,
		&entity.WarehouseHasOfficeArea,
		&entity.WarehouseOfficeArea,
		&entity.StoreHasMezzanine,
		&entity.StoreMezzanineArea,
	); err != nil {
		return listingentity.ListingEntity{}, err
	}

	return entity, nil
}

// GetListingByQuery executes a custom query and returns a FULLY ENRICHED listing
//
// This method is the centralized helper for fetching listing versions with ALL satellite tables.
// It performs the following enrichment operations IN ORDER:
//  1. Scans base listing entity from query result
//  2. Fetches ExchangePlaces (listing_exchange_places table)
//  3. Fetches Features (listing_features table)
//  4. Fetches Guarantees (listing_guarantees table)
//  5. Fetches FinancingBlockers (listing_financing_blockers table)
//  6. Fetches WarehouseAdditionalFloors (listing_warehouse_additional_floors table)
//  7. Converts enriched entity to domain model via ListingEntityToDomain
//
// The query parameter MUST return columns matching listingSelectColumns constant.
//
// Parameters:
//   - ctx: Context for tracing, cancellation, and logging
//   - tx: Database transaction (can be nil for standalone queries)
//   - query: SQL query returning listingSelectColumns (must use SELECT <listingSelectColumns>)
//   - args: Query arguments (e.g., version_id, listing_identity_id)
//
// Returns:
//   - listing: ListingInterface with ALL satellite tables populated
//   - error: sql.ErrNoRows if no listing found, or other database errors
//
// Usage:
//
//	// Example 1: Get by version ID
//	query := `SELECT ` + listingSelectColumns + ` FROM listing_versions lv
//	          JOIN listing_identities li ON lv.listing_identity_id = li.id
//	          WHERE lv.id = ? AND lv.deleted = 0 LIMIT 1`
//	listing, err := la.GetListingByQuery(ctx, tx, query, versionID)
//
//	// Example 2: Get active version
//	query := `SELECT ` + listingSelectColumns + ` FROM listing_versions lv
//	          JOIN listing_identities li ON lv.listing_identity_id = li.id
//	          WHERE li.id = ? AND lv.id = li.active_version_id LIMIT 1`
//	listing, err := la.GetListingByQuery(ctx, tx, query, identityID)
//
// Error Handling:
//   - Returns sql.ErrNoRows if query returns 0 rows
//   - Logs ERROR and marks span on scan failures
//   - Logs ERROR and marks span if any satellite table fetch fails
//   - Propagates errors from enrichment methods (features, guarantees, etc.)
func (la *ListingAdapter) GetListingByQuery(ctx context.Context, tx *sql.Tx, query string, args ...any) (listing listingmodel.ListingInterface, err error) {
	// Initialize tracing for observability (metrics + distributed tracing)
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	// Attach logger to context to ensure request_id/trace_id propagation
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	// Execute query using instrumented adapter (auto-generates metrics + tracing)
	row := la.QueryRowContext(ctx, tx, "select", query, args...)

	// Scan base listing entity from query result
	entityListing, err := scanListingEntity(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		// Mark span as error for distributed tracing analysis
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.listing.get_listing_by_query.scan_error", "error", err)
		return nil, fmt.Errorf("scan listing by query: %w", err)
	}

	// Enrich with ExchangePlaces (listing_exchange_places table)
	// Note: Returns empty slice if no exchange places (not an error)
	entityExchangePlaces, err := la.GetEntityExchangePlacesByListing(ctx, tx, entityListing.ID)
	if err != nil {
		// Infrastructure error: log and mark span
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.listing.get_listing_by_query.exchange_places_error", "listing_version_id", entityListing.ID, "error", err)
		return
	}
	entityListing.ExchangePlaces = entityExchangePlaces

	// Enrich with Features (listing_features table)
	// Note: Returns empty slice if no features (not an error)
	entityFeatures, err := la.GetEntityFeaturesByListing(ctx, tx, entityListing.ID)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.listing.get_listing_by_query.features_error", "listing_version_id", entityListing.ID, "error", err)
		return
	}
	entityListing.Features = entityFeatures

	// Enrich with Guarantees (listing_guarantees table)
	// Note: Returns empty slice if no guarantees (not an error)
	entityGuarantees, err := la.GetEntityGuaranteesByListing(ctx, tx, entityListing.ID)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.listing.get_listing_by_query.guarantees_error", "listing_version_id", entityListing.ID, "error", err)
		return
	}
	entityListing.Guarantees = entityGuarantees

	// Enrich with FinancingBlockers (listing_financing_blockers table)
	// Note: Returns empty slice if no blockers (not an error)
	entityFinancingBlockers, err := la.GetEntityFinancingBlockersByListing(ctx, tx, entityListing.ID)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.listing.get_listing_by_query.financing_blockers_error", "listing_version_id", entityListing.ID, "error", err)
		return
	}
	entityListing.FinancingBlocker = entityFinancingBlockers

	// Enrich with WarehouseAdditionalFloors (listing_warehouse_additional_floors table)
	// Note: Returns empty slice if no additional floors (not an error)
	entityWarehouseFloors, err := la.GetEntityWarehouseAdditionalFloorsByListing(ctx, tx, entityListing.ID)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.listing.get_listing_by_query.warehouse_floors_error", "listing_version_id", entityListing.ID, "error", err)
		return
	}
	entityListing.WarehouseAdditionalFloors = entityWarehouseFloors

	// Convert enriched entity to domain model (separation of concerns)
	// Note: ListingEntityToDomain handles all sql.Null* conversions and populates all fields
	listing = listingconverters.ListingEntityToDomain(entityListing)

	return
}
