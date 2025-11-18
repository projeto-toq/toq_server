package mysqllistingadapter

import (
	"context"
	"database/sql"
	"fmt"

	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (la *ListingAdapter) CreateListingVersion(ctx context.Context, tx *sql.Tx, version listingmodel.ListingVersionInterface) error {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	listingIdentityID := version.ListingIdentityID()
	if listingIdentityID == 0 {
		return fmt.Errorf("listing version missing identity id")
	}

	query := `INSERT INTO listing_versions (
        listing_identity_id,
        user_id,
		code,
		version,
		status,
		zip_code,
		street,
		number,
		complement,
		complex,
		neighborhood,
		city,
		state,
		title,
		type,
		owner,
		land_size,
		corner,
		non_buildable,
		buildable,
		delivered,
		who_lives,
		description,
		transaction,
		sell_net,
		rent_net,
		condominium,
		annual_tax,
		monthly_tax,
		annual_ground_rent,
		monthly_ground_rent,
		exchange,
		exchange_perc,
		installment,
		financing,
		visit,
		tenant_name,
		tenant_email,
		tenant_phone,
		accompanying,
		deleted,
		completion_forecast,
		land_block,
		land_lot,
		land_front,
		land_side,
		land_back,
		land_terrain_type,
		has_kmz,
		kmz_file,
		building_floors,
		unit_tower,
		unit_floor,
		unit_number,
		warehouse_manufacturing_area,
		warehouse_sector,
		warehouse_has_primary_cabin,
		warehouse_cabin_kva,
		warehouse_ground_floor,
		warehouse_floor_resistance,
		warehouse_zoning,
		warehouse_has_office_area,
		warehouse_office_area,
		store_has_mezzanine,
		store_mezzanine_area
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	street := sql.NullString{String: version.Street(), Valid: version.Street() != ""}
	number := sql.NullString{String: version.Number(), Valid: true}
	complement := sql.NullString{String: version.Complement(), Valid: version.Complement() != ""}
	complexValue := sql.NullString{}
	if version.HasComplex() {
		complexValue = sql.NullString{String: version.Complex(), Valid: true}
	}
	neighborhood := sql.NullString{String: version.Neighborhood(), Valid: version.Neighborhood() != ""}
	city := sql.NullString{String: version.City(), Valid: version.City() != ""}
	state := sql.NullString{String: version.State(), Valid: version.State() != ""}

	title := sql.NullString{}
	if version.HasTitle() {
		title = sql.NullString{String: version.Title(), Valid: true}
	}

	typeValue := sql.NullInt64{Int64: int64(version.ListingType()), Valid: true}

	owner := sql.NullInt64{}
	if version.HasOwner() {
		owner = sql.NullInt64{Int64: int64(version.Owner()), Valid: true}
	}

	landSize := sql.NullFloat64{}
	if version.HasLandSize() {
		landSize = sql.NullFloat64{Float64: version.LandSize(), Valid: true}
	}

	corner := sql.NullBool{}
	if version.HasCorner() {
		corner = sql.NullBool{Bool: version.Corner(), Valid: true}
	}

	nonBuildable := sql.NullFloat64{}
	if version.HasNonBuildable() {
		nonBuildable = sql.NullFloat64{Float64: version.NonBuildable(), Valid: true}
	}

	buildable := sql.NullFloat64{}
	if version.HasBuildable() {
		buildable = sql.NullFloat64{Float64: version.Buildable(), Valid: true}
	}

	delivered := sql.NullInt64{}
	if version.HasDelivered() {
		delivered = sql.NullInt64{Int64: int64(version.Delivered()), Valid: true}
	}

	whoLives := sql.NullInt64{}
	if version.HasWhoLives() {
		whoLives = sql.NullInt64{Int64: int64(version.WhoLives()), Valid: true}
	}

	description := sql.NullString{}
	if version.HasDescription() {
		description = sql.NullString{String: version.Description(), Valid: true}
	}

	transaction := sql.NullInt64{}
	if version.HasTransaction() {
		transaction = sql.NullInt64{Int64: int64(version.Transaction()), Valid: true}
	}

	sellNet := sql.NullFloat64{}
	if version.HasSellNet() {
		sellNet = sql.NullFloat64{Float64: version.SellNet(), Valid: true}
	}

	rentNet := sql.NullFloat64{}
	if version.HasRentNet() {
		rentNet = sql.NullFloat64{Float64: version.RentNet(), Valid: true}
	}

	condominium := sql.NullFloat64{}
	if version.HasCondominium() {
		condominium = sql.NullFloat64{Float64: version.Condominium(), Valid: true}
	}

	annualTax := sql.NullFloat64{}
	if version.HasAnnualTax() {
		annualTax = sql.NullFloat64{Float64: version.AnnualTax(), Valid: true}
	}

	monthlyTax := sql.NullFloat64{}
	if version.HasMonthlyTax() {
		monthlyTax = sql.NullFloat64{Float64: version.MonthlyTax(), Valid: true}
	}

	annualGroundRent := sql.NullFloat64{}
	if version.HasAnnualGroundRent() {
		annualGroundRent = sql.NullFloat64{Float64: version.AnnualGroundRent(), Valid: true}
	}

	monthlyGroundRent := sql.NullFloat64{}
	if version.HasMonthlyGroundRent() {
		monthlyGroundRent = sql.NullFloat64{Float64: version.MonthlyGroundRent(), Valid: true}
	}

	exchange := sql.NullBool{}
	if version.HasExchange() {
		exchange = sql.NullBool{Bool: version.Exchange(), Valid: true}
	}

	exchangePercentual := sql.NullFloat64{}
	if version.HasExchangePercentual() {
		exchangePercentual = sql.NullFloat64{Float64: version.ExchangePercentual(), Valid: true}
	}

	installment := sql.NullInt64{}
	if version.HasInstallment() {
		installment = sql.NullInt64{Int64: int64(version.Installment()), Valid: true}
	}

	financing := sql.NullBool{}
	if version.HasFinancing() {
		financing = sql.NullBool{Bool: version.Financing(), Valid: true}
	}

	visit := sql.NullInt64{}
	if version.HasVisit() {
		visit = sql.NullInt64{Int64: int64(version.Visit()), Valid: true}
	}

	tenantName := sql.NullString{}
	if version.HasTenantName() {
		tenantName = sql.NullString{String: version.TenantName(), Valid: true}
	}

	tenantEmail := sql.NullString{}
	if version.HasTenantEmail() {
		tenantEmail = sql.NullString{String: version.TenantEmail(), Valid: true}
	}

	tenantPhone := sql.NullString{}
	if version.HasTenantPhone() {
		tenantPhone = sql.NullString{String: version.TenantPhone(), Valid: true}
	}

	accompanying := sql.NullInt64{}
	if version.HasAccompanying() {
		accompanying = sql.NullInt64{Int64: int64(version.Accompanying()), Valid: true}
	}

	deletedValue := version.Deleted()

	// New property-specific fields
	// Sanitize completion forecast to ensure MySQL DATE compatibility
	completionForecast := sql.NullString{}
	if version.HasCompletionForecast() {
		rawValue := version.CompletionForecast()

		// Attempt to normalize format (defensive programming)
		normalized, parseErr := utils.ParseCompletionForecast(rawValue)
		if parseErr != nil {
			// Log warning but don't fail the insert (domain should have validated)
			logger.Warn("mysql.listing.create_listing_version.invalid_completion_forecast_format",
				"raw_value", rawValue, "error", parseErr)
			// Store raw value (will likely fail MySQL constraint, triggering proper error handling)
			completionForecast = sql.NullString{String: rawValue, Valid: true}
		} else {
			completionForecast = sql.NullString{String: normalized, Valid: true}
		}
	}

	landBlock := sql.NullString{}
	if version.HasLandBlock() {
		landBlock = sql.NullString{String: version.LandBlock(), Valid: true}
	}

	landLot := sql.NullString{}
	if version.HasLandLot() {
		landLot = sql.NullString{String: version.LandLot(), Valid: true}
	}

	landFront := sql.NullFloat64{}
	if version.HasLandFront() {
		landFront = sql.NullFloat64{Float64: version.LandFront(), Valid: true}
	}

	landSide := sql.NullFloat64{}
	if version.HasLandSide() {
		landSide = sql.NullFloat64{Float64: version.LandSide(), Valid: true}
	}

	landBack := sql.NullFloat64{}
	if version.HasLandBack() {
		landBack = sql.NullFloat64{Float64: version.LandBack(), Valid: true}
	}

	landTerrainType := sql.NullInt64{}
	if version.HasLandTerrainType() {
		landTerrainType = sql.NullInt64{Int64: int64(version.LandTerrainType()), Valid: true}
	}

	hasKmz := sql.NullBool{}
	if version.HasHasKmz() {
		hasKmz = sql.NullBool{Bool: version.HasKmz(), Valid: true}
	}

	kmzFile := sql.NullString{}
	if version.HasKmzFile() {
		kmzFile = sql.NullString{String: version.KmzFile(), Valid: true}
	}

	buildingFloors := sql.NullInt64{}
	if version.HasBuildingFloors() {
		buildingFloors = sql.NullInt64{Int64: int64(version.BuildingFloors()), Valid: true}
	}

	unitTower := sql.NullString{}
	if version.HasUnitTower() {
		unitTower = sql.NullString{String: version.UnitTower(), Valid: true}
	}

	unitFloor := sql.NullString{}
	if version.HasUnitFloor() {
		unitFloor = sql.NullString{String: version.UnitFloor(), Valid: true}
	}

	unitNumber := sql.NullString{}
	if version.HasUnitNumber() {
		unitNumber = sql.NullString{String: version.UnitNumber(), Valid: true}
	}

	warehouseManufacturingArea := sql.NullFloat64{}
	if version.HasWarehouseManufacturingArea() {
		warehouseManufacturingArea = sql.NullFloat64{Float64: version.WarehouseManufacturingArea(), Valid: true}
	}

	warehouseSector := sql.NullInt64{}
	if version.HasWarehouseSector() {
		warehouseSector = sql.NullInt64{Int64: int64(version.WarehouseSector()), Valid: true}
	}

	warehouseHasPrimaryCabin := sql.NullBool{}
	if version.HasWarehouseHasPrimaryCabin() {
		warehouseHasPrimaryCabin = sql.NullBool{Bool: version.WarehouseHasPrimaryCabin(), Valid: true}
	}

	warehouseCabinKva := sql.NullString{}
	if version.HasWarehouseCabinKva() {
		warehouseCabinKva = sql.NullString{String: version.WarehouseCabinKva(), Valid: true}
	}

	warehouseGroundFloor := sql.NullInt64{}
	if version.HasWarehouseGroundFloor() {
		warehouseGroundFloor = sql.NullInt64{Int64: int64(version.WarehouseGroundFloor()), Valid: true}
	}

	warehouseFloorResistance := sql.NullFloat64{}
	if version.HasWarehouseFloorResistance() {
		warehouseFloorResistance = sql.NullFloat64{Float64: version.WarehouseFloorResistance(), Valid: true}
	}

	warehouseZoning := sql.NullString{}
	if version.HasWarehouseZoning() {
		warehouseZoning = sql.NullString{String: version.WarehouseZoning(), Valid: true}
	}

	warehouseHasOfficeArea := sql.NullBool{}
	if version.HasWarehouseHasOfficeArea() {
		warehouseHasOfficeArea = sql.NullBool{Bool: version.WarehouseHasOfficeArea(), Valid: true}
	}

	warehouseOfficeArea := sql.NullFloat64{}
	if version.HasWarehouseOfficeArea() {
		warehouseOfficeArea = sql.NullFloat64{Float64: version.WarehouseOfficeArea(), Valid: true}
	}

	storeHasMezzanine := sql.NullBool{}
	if version.HasStoreHasMezzanine() {
		storeHasMezzanine = sql.NullBool{Bool: version.StoreHasMezzanine(), Valid: true}
	}

	storeMezzanineArea := sql.NullFloat64{}
	if version.HasStoreMezzanineArea() {
		storeMezzanineArea = sql.NullFloat64{Float64: version.StoreMezzanineArea(), Valid: true}
	}

	result, execErr := la.ExecContext(ctx, tx, "insert", query,
		listingIdentityID,
		version.UserID(),
		version.Code(),
		version.Version(),
		version.Status(),
		version.ZipCode(),
		street,
		number,
		complement,
		complexValue,
		neighborhood,
		city,
		state,
		title,
		typeValue,
		owner,
		landSize,
		corner,
		nonBuildable,
		buildable,
		delivered,
		whoLives,
		description,
		transaction,
		sellNet,
		rentNet,
		condominium,
		annualTax,
		monthlyTax,
		annualGroundRent,
		monthlyGroundRent,
		exchange,
		exchangePercentual,
		installment,
		financing,
		visit,
		tenantName,
		tenantEmail,
		tenantPhone,
		accompanying,
		deletedValue,
		completionForecast,
		landBlock,
		landLot,
		landFront,
		landSide,
		landBack,
		landTerrainType,
		hasKmz,
		kmzFile,
		buildingFloors,
		unitTower,
		unitFloor,
		unitNumber,
		warehouseManufacturingArea,
		warehouseSector,
		warehouseHasPrimaryCabin,
		warehouseCabinKva,
		warehouseGroundFloor,
		warehouseFloorResistance,
		warehouseZoning,
		warehouseHasOfficeArea,
		warehouseOfficeArea,
		storeHasMezzanine,
		storeMezzanineArea,
	)
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.listing.create_listing_version.exec_error", "error", execErr)
		return fmt.Errorf("exec create listing version: %w", execErr)
	}

	id, lastErr := result.LastInsertId()
	if lastErr != nil {
		utils.SetSpanError(ctx, lastErr)
		logger.Error("mysql.listing.create_listing_version.last_insert_error", "error", lastErr)
		return fmt.Errorf("last insert id for listing version: %w", lastErr)
	}

	version.SetID(id)

	return nil
}
