package mysqllistingadapter

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	globalmodel "github.com/projeto-toq/toq_server/internal/core/model/global_model"
	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
	listingrepository "github.com/projeto-toq/toq_server/internal/core/port/right/repository/listing_repository"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// GetListingForEndUpdate retrieves comprehensive listing version data for validation during end-update and promotion flows.
//
// This method fetches a specific listing version by its version ID and returns aggregated data including:
// - Version metadata (identity_id, user_id, code, version, status)
// - Property details (address, type, transaction, pricing)
// - Satellite entity counts (features, guarantees, exchange_places, financing_blockers)
//
// The method joins listing_versions with listing_identities to ensure both records are not soft-deleted.
// It uses the InstrumentedAdapter for query execution, providing automatic tracing and metrics.
//
// Parameters:
//   - ctx: Context with request/trace information
//   - tx: Active database transaction
//   - versionID: The ID of the listing_version record to fetch (NOT the listing_identity_id)
//
// Returns:
//   - ListingEndUpdateData: Aggregated data struct with all necessary fields for validation
//   - error: sql.ErrNoRows if version not found, or infrastructure error if query fails
//
// Business Rules in Query:
//   - Only returns versions where both lv.deleted = 0 AND li.deleted = 0
//   - Joins listing_identities to validate referential integrity
//
// Edge Cases:
//   - Satellite entities (features, guarantees, etc.) may return sql.ErrNoRows → treated as empty, not error
//   - Returns sql.ErrNoRows if version doesn't exist or is soft-deleted
func (la *ListingAdapter) GetListingForEndUpdate(ctx context.Context, tx *sql.Tx, versionID int64) (listingrepository.ListingEndUpdateData, error) {
	data := listingrepository.ListingEndUpdateData{}

	// Initialize tracer for this repository operation
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return data, err
	}
	defer spanEnd()

	// Propagate logger with trace context for structured logging
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	// Query explanation:
	// - Selects lv.listing_identity_id (NOT lv.id) to populate ListingID field
	// - Joins with listing_identities to enforce referential integrity
	// - Filters by lv.id (version ID) and ensures both version and identity are not deleted
	query := `SELECT lv.listing_identity_id, lv.user_id, lv.code, lv.version, lv.status, lv.zip_code, lv.street, lv.number, lv.complex, lv.city, lv.state,
		lv.title, lv.type, lv.owner, lv.buildable, lv.delivered, lv.who_lives, lv.description, lv.transaction, lv.visit,
		lv.accompanying, lv.annual_tax, lv.monthly_tax, lv.annual_ground_rent, lv.monthly_ground_rent, lv.exchange,
		lv.exchange_perc, lv.sell_net, lv.rent_net, lv.condominium, lv.land_size, lv.corner, lv.tenant_name, lv.tenant_phone,
		lv.tenant_email, lv.financing, lv.completion_forecast, lv.land_block, lv.land_lot, lv.land_front, lv.land_side,
		lv.land_back, lv.land_terrain_type, lv.has_kmz, lv.kmz_file, lv.building_floors, lv.unit_tower, lv.unit_floor,
		lv.unit_number, lv.warehouse_manufacturing_area, lv.warehouse_sector, lv.warehouse_has_primary_cabin,
		lv.warehouse_cabin_kva, lv.warehouse_ground_floor, lv.warehouse_floor_resistance, lv.warehouse_zoning,
		lv.warehouse_has_office_area, lv.warehouse_office_area, lv.store_has_mezzanine, lv.store_mezzanine_area
		FROM listing_versions lv
		INNER JOIN listing_identities li ON li.id = lv.listing_identity_id
		WHERE lv.id = ? AND lv.deleted = 0 AND li.deleted = 0`

	var (
		status                     uint8
		listingType                uint16
		street                     sql.NullString
		number                     sql.NullString
		complex                    sql.NullString
		city                       sql.NullString
		state                      sql.NullString
		title                      sql.NullString
		owner                      sql.NullInt16
		buildable                  sql.NullFloat64
		delivered                  sql.NullInt16
		whoLives                   sql.NullInt16
		description                sql.NullString
		transaction                sql.NullInt16
		visit                      sql.NullInt16
		accompanying               sql.NullInt16
		annualTax                  sql.NullFloat64
		monthlyTax                 sql.NullFloat64
		annualGroundRent           sql.NullFloat64
		monthlyGroundRent          sql.NullFloat64
		exchange                   sql.NullInt16
		exchangePerc               sql.NullFloat64
		saleNet                    sql.NullFloat64
		rentNet                    sql.NullFloat64
		condominium                sql.NullFloat64
		landSize                   sql.NullFloat64
		corner                     sql.NullInt16
		tenantName                 sql.NullString
		tenantPhone                sql.NullString
		tenantEmail                sql.NullString
		financing                  sql.NullInt16
		completionForecast         sql.NullString
		landBlock                  sql.NullString
		landLot                    sql.NullString
		landFront                  sql.NullFloat64
		landSide                   sql.NullFloat64
		landBack                   sql.NullFloat64
		landTerrainType            sql.NullInt16
		hasKmz                     sql.NullInt16
		kmzFile                    sql.NullString
		buildingFloors             sql.NullInt16
		unitTower                  sql.NullString
		unitFloor                  sql.NullString
		unitNumber                 sql.NullString
		warehouseManufacturingArea sql.NullFloat64
		warehouseSector            sql.NullInt16
		warehouseHasPrimaryCabin   sql.NullInt16
		warehouseCabinKva          sql.NullString
		warehouseGroundFloor       sql.NullInt16
		warehouseFloorResistance   sql.NullFloat64
		warehouseZoning            sql.NullString
		warehouseHasOfficeArea     sql.NullInt16
		warehouseOfficeArea        sql.NullFloat64
		storeHasMezzanine          sql.NullInt16
		storeMezzanineArea         sql.NullFloat64
	)

	// Use InstrumentedAdapter for query execution (automatic tracing + metrics)
	row := la.QueryRowContext(ctx, tx, "select", query, versionID)

	// Scan result into data struct
	// IMPORTANT: First column is now lv.listing_identity_id (not lv.id)
	err = row.Scan(
		&data.ListingID,
		&data.UserID,
		&data.Code,
		&data.Version,
		&status,
		&data.ZipCode,
		&street,
		&number,
		&complex,
		&city,
		&state,
		&title,
		&listingType,
		&owner,
		&buildable,
		&delivered,
		&whoLives,
		&description,
		&transaction,
		&visit,
		&accompanying,
		&annualTax,
		&monthlyTax,
		&annualGroundRent,
		&monthlyGroundRent,
		&exchange,
		&exchangePerc,
		&saleNet,
		&rentNet,
		&condominium,
		&landSize,
		&corner,
		&tenantName,
		&tenantPhone,
		&tenantEmail,
		&financing,
		&completionForecast,
		&landBlock,
		&landLot,
		&landFront,
		&landSide,
		&landBack,
		&landTerrainType,
		&hasKmz,
		&kmzFile,
		&buildingFloors,
		&unitTower,
		&unitFloor,
		&unitNumber,
		&warehouseManufacturingArea,
		&warehouseSector,
		&warehouseHasPrimaryCabin,
		&warehouseCabinKva,
		&warehouseGroundFloor,
		&warehouseFloorResistance,
		&warehouseZoning,
		&warehouseHasOfficeArea,
		&warehouseOfficeArea,
		&storeHasMezzanine,
		&storeMezzanineArea,
	)
	if err != nil {
		// Return sql.ErrNoRows directly for "not found" handling in service layer
		if errors.Is(err, sql.ErrNoRows) {
			return data, sql.ErrNoRows
		}
		// Mark span as error and log infrastructure failure
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.listing.get_listing_for_end_update.scan_error", "error", err, "version_id", versionID)
		return data, fmt.Errorf("scan listing for end update: %w", err)
	}

	data.Status = listingmodel.ListingStatus(status)
	data.ListingType = globalmodel.PropertyType(listingType)
	data.Street = street
	data.Number = number
	data.Complex = complex
	data.City = city
	data.State = state
	data.Title = title
	data.Owner = owner
	data.Buildable = buildable
	data.Delivered = delivered
	data.WhoLives = whoLives
	data.Description = description
	data.Transaction = transaction
	data.Visit = visit
	data.Accompanying = accompanying
	data.AnnualTax = annualTax
	data.MonthlyTax = monthlyTax
	data.AnnualGroundRent = annualGroundRent
	data.MonthlyGroundRent = monthlyGroundRent
	data.Exchange = exchange
	data.ExchangePercentual = exchangePerc
	data.SaleNet = saleNet
	data.RentNet = rentNet
	data.Condominium = condominium
	data.LandSize = landSize
	data.Corner = corner
	data.TenantName = tenantName
	data.TenantPhone = tenantPhone
	data.TenantEmail = tenantEmail
	data.Financing = financing
	data.CompletionForecast = completionForecast
	data.LandBlock = landBlock
	data.LandLot = landLot
	data.LandFront = landFront
	data.LandSide = landSide
	data.LandBack = landBack
	data.LandTerrainType = landTerrainType
	data.HasKmz = hasKmz
	data.KmzFile = kmzFile
	data.BuildingFloors = buildingFloors
	data.UnitTower = unitTower
	data.UnitFloor = unitFloor
	data.UnitNumber = unitNumber
	data.WarehouseManufacturingArea = warehouseManufacturingArea
	data.WarehouseSector = warehouseSector
	data.WarehouseHasPrimaryCabin = warehouseHasPrimaryCabin
	data.WarehouseCabinKva = warehouseCabinKva
	data.WarehouseGroundFloor = warehouseGroundFloor
	data.WarehouseFloorResistance = warehouseFloorResistance
	data.WarehouseZoning = warehouseZoning
	data.WarehouseHasOfficeArea = warehouseHasOfficeArea
	data.WarehouseOfficeArea = warehouseOfficeArea
	data.StoreHasMezzanine = storeHasMezzanine
	data.StoreMezzanineArea = storeMezzanineArea

	// Fetch satellite entities counts
	// Note: Satellite tables reference listing_version_id, not listing_identity_id
	// Therefore we use versionID to fetch related entities

	// Fetch features count
	features, ferr := la.GetEntityFeaturesByListing(ctx, tx, versionID)
	if ferr != nil && !errors.Is(ferr, sql.ErrNoRows) {
		utils.SetSpanError(ctx, ferr)
		logger.Error("mysql.listing.get_listing_for_end_update.features_error", "error", ferr, "version_id", versionID)
		return data, fmt.Errorf("get features for end update: %w", ferr)
	}
	data.FeaturesCount = len(features)

	// Fetch exchange places count
	exchangePlaces, eerr := la.GetEntityExchangePlacesByListing(ctx, tx, versionID)
	if eerr != nil && !errors.Is(eerr, sql.ErrNoRows) {
		utils.SetSpanError(ctx, eerr)
		logger.Error("mysql.listing.get_listing_for_end_update.exchange_places_error", "error", eerr, "version_id", versionID)
		return data, fmt.Errorf("get exchange places for end update: %w", eerr)
	}
	data.ExchangePlacesCount = len(exchangePlaces)

	// Fetch financing blockers count
	financingBlockers, berr := la.GetEntityFinancingBlockersByListing(ctx, tx, versionID)
	if berr != nil && !errors.Is(berr, sql.ErrNoRows) {
		utils.SetSpanError(ctx, berr)
		logger.Error("mysql.listing.get_listing_for_end_update.financing_blockers_error", "error", berr, "version_id", versionID)
		return data, fmt.Errorf("get financing blockers for end update: %w", berr)
	}
	data.FinancingBlockersCount = len(financingBlockers)

	// Fetch guarantees count
	guarantees, gerr := la.GetEntityGuaranteesByListing(ctx, tx, versionID)
	if gerr != nil && !errors.Is(gerr, sql.ErrNoRows) {
		utils.SetSpanError(ctx, gerr)
		logger.Error("mysql.listing.get_listing_for_end_update.guarantees_error", "error", gerr, "version_id", versionID)
		return data, fmt.Errorf("get guarantees for end update: %w", gerr)
	}
	data.GuaranteesCount = len(guarantees)

	return data, nil
}
