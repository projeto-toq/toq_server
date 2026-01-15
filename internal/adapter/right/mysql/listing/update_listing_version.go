package mysqllistingadapter

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// UpdateListingVersion updates an existing listing version record.
// This method updates all mutable fields of a version (except id, listing_identity_id, code, version).
func (la *ListingAdapter) UpdateListingVersion(ctx context.Context, tx *sql.Tx, version listingmodel.ListingVersionInterface) error {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `
		UPDATE listing_versions SET
			status = ?, title = ?, zip_code = ?, street = ?, number = ?, complement = ?, complex = ?,
			neighborhood = ?, city = ?, state = ?, type = ?, owner = ?, land_size = ?,
			corner = ?, non_buildable = ?, buildable = ?, delivered = ?, who_lives = ?,
			description = ?, transaction = ?, sell_net = ?, rent_net = ?, condominium = ?,
			annual_tax = ?, monthly_tax = ?, annual_ground_rent = ?, monthly_ground_rent = ?,
			exchange = ?, exchange_perc = ?, installment = ?, financing = ?, visit = ?,
			tenant_name = ?, tenant_email = ?, tenant_phone = ?, accompanying = ?,
			completion_forecast = ?, land_block = ?, land_lot = ?, land_front = ?, land_side = ?,
			land_back = ?, land_terrain_type = ?, has_kmz = ?, kmz_file = ?, building_floors = ?,
			unit_tower = ?, unit_floor = ?, unit_number = ?, warehouse_manufacturing_area = ?,
			warehouse_sector = ?, warehouse_has_primary_cabin = ?, warehouse_cabin_kva = ?,
			warehouse_ground_floor = ?, warehouse_floor_resistance = ?, warehouse_zoning = ?,
			warehouse_has_office_area = ?, warehouse_office_area = ?, store_has_mezzanine = ?,
			store_mezzanine_area = ?, price_updated_at = ?
		WHERE id = ? AND deleted = 0
	`

	var title, description, complexValue, tenantName, tenantEmail, tenantPhone interface{}
	var owner, landSize, corner, nonBuildable, buildable, delivered, whoLives, transaction interface{}
	var sellNet, rentNet, condominium, annualTax, monthlyTax, annualGroundRent, monthlyGroundRent interface{}
	var exchange, exchangePerc, installment, financing, visit, accompanying interface{}
	var completionForecast, landBlock, landLot, landFront, landSide, landBack, landTerrainType interface{}
	var hasKmz, kmzFile, buildingFloors, unitTower, unitFloor, unitNumber interface{}
	var warehouseManufacturingArea, warehouseSector, warehouseHasPrimaryCabin, warehouseCabinKva interface{}
	var warehouseGroundFloor, warehouseFloorResistance, warehouseZoning interface{}
	var warehouseHasOfficeArea, warehouseOfficeArea, storeHasMezzanine, storeMezzanineArea interface{}
	var priceUpdatedAt interface{}

	// Required string fields - always set
	street := version.Street()
	complement := version.Complement()
	neighborhood := version.Neighborhood()
	city := version.City()
	state := version.State()

	if version.HasTitle() {
		title = version.Title()
	}
	if version.HasDescription() {
		description = version.Description()
	}
	if version.HasComplex() {
		complexValue = version.Complex()
	}
	if version.HasOwner() {
		owner = uint8(version.Owner())
	}
	if version.HasLandSize() {
		landSize = version.LandSize()
	}
	if version.HasCorner() {
		corner = version.Corner()
	}
	if version.HasNonBuildable() {
		nonBuildable = version.NonBuildable()
	}
	if version.HasBuildable() {
		buildable = version.Buildable()
	}
	if version.HasDelivered() {
		delivered = uint8(version.Delivered())
	}
	if version.HasWhoLives() {
		whoLives = uint8(version.WhoLives())
	}
	if version.HasTransaction() {
		transaction = uint8(version.Transaction())
	}
	if version.HasSellNet() {
		sellNet = version.SellNet()
	}
	if version.HasRentNet() {
		rentNet = version.RentNet()
	}
	if version.HasCondominium() {
		condominium = version.Condominium()
	}
	if version.HasAnnualTax() {
		annualTax = version.AnnualTax()
	}
	if version.HasMonthlyTax() {
		monthlyTax = version.MonthlyTax()
	}
	if version.HasAnnualGroundRent() {
		annualGroundRent = version.AnnualGroundRent()
	}
	if version.HasMonthlyGroundRent() {
		monthlyGroundRent = version.MonthlyGroundRent()
	}
	if version.HasExchange() {
		exchange = version.Exchange()
	}
	if version.HasExchangePercentual() {
		exchangePerc = version.ExchangePercentual()
	}
	if version.HasInstallment() {
		installment = uint8(version.Installment())
	}
	if version.HasFinancing() {
		financing = version.Financing()
	}
	if version.HasVisit() {
		visit = uint8(version.Visit())
	}
	if version.HasTenantName() {
		tenantName = version.TenantName()
	}
	if version.HasTenantEmail() {
		tenantEmail = version.TenantEmail()
	}
	if version.HasTenantPhone() {
		tenantPhone = version.TenantPhone()
	}
	if version.HasAccompanying() {
		accompanying = uint8(version.Accompanying())
	}

	// Sanitize completion forecast to ensure MySQL DATE compatibility
	// Layer of defense: even if domain validation failed, ensure correct format here
	if version.HasCompletionForecast() {
		rawValue := version.CompletionForecast()

		// Attempt to normalize format (defensive programming)
		normalized, parseErr := utils.ParseCompletionForecast(rawValue)
		if parseErr != nil {
			// Log warning but don't fail the update (domain should have validated)
			logger.Warn("mysql.listing.update_listing_version.invalid_completion_forecast_format",
				"raw_value", rawValue, "error", parseErr)
			// Store raw value (will likely fail MySQL constraint, triggering proper error handling)
			completionForecast = rawValue
		} else {
			completionForecast = normalized
		}
	}

	if version.HasLandBlock() {
		landBlock = version.LandBlock()
	}
	if version.HasLandLot() {
		landLot = version.LandLot()
	}
	if version.HasLandFront() {
		landFront = version.LandFront()
	}
	if version.HasLandSide() {
		landSide = version.LandSide()
	}
	if version.HasLandBack() {
		landBack = version.LandBack()
	}
	if version.HasLandTerrainType() {
		landTerrainType = uint8(version.LandTerrainType())
	}
	if version.HasHasKmz() {
		hasKmz = version.HasKmz()
	}
	if version.HasKmzFile() {
		kmzFile = version.KmzFile()
	}
	if version.HasBuildingFloors() {
		buildingFloors = version.BuildingFloors()
	}
	if version.HasUnitTower() {
		unitTower = version.UnitTower()
	}
	if version.HasUnitFloor() {
		unitFloor = version.UnitFloor()
	}
	if version.HasUnitNumber() {
		unitNumber = version.UnitNumber()
	}
	if version.HasWarehouseManufacturingArea() {
		warehouseManufacturingArea = version.WarehouseManufacturingArea()
	}
	if version.HasWarehouseSector() {
		warehouseSector = uint8(version.WarehouseSector())
	}
	if version.HasWarehouseHasPrimaryCabin() {
		warehouseHasPrimaryCabin = version.WarehouseHasPrimaryCabin()
	}
	if version.HasWarehouseCabinKva() {
		warehouseCabinKva = version.WarehouseCabinKva()
	}
	if version.HasWarehouseGroundFloor() {
		warehouseGroundFloor = version.WarehouseGroundFloor()
	}
	if version.HasWarehouseFloorResistance() {
		warehouseFloorResistance = version.WarehouseFloorResistance()
	}
	if version.HasWarehouseZoning() {
		warehouseZoning = version.WarehouseZoning()
	}
	if version.HasWarehouseHasOfficeArea() {
		warehouseHasOfficeArea = version.WarehouseHasOfficeArea()
	}
	if version.HasWarehouseOfficeArea() {
		warehouseOfficeArea = version.WarehouseOfficeArea()
	}
	if version.HasStoreHasMezzanine() {
		storeHasMezzanine = version.StoreHasMezzanine()
	}
	if version.HasStoreMezzanineArea() {
		storeMezzanineArea = version.StoreMezzanineArea()
	}

	if listingWithPrice, ok := version.(listingmodel.ListingInterface); ok {
		priceUpdatedAt = listingWithPrice.PriceUpdatedAt()
	} else {
		priceUpdatedAt = time.Now().UTC()
	}

	_, execErr := la.ExecContext(ctx, tx, "update", query,
		uint8(version.Status()), title, version.ZipCode(), street, version.Number(), complement, complexValue,
		neighborhood, city, state, uint16(version.ListingType()), owner, landSize,
		corner, nonBuildable, buildable, delivered, whoLives,
		description, transaction, sellNet, rentNet, condominium,
		annualTax, monthlyTax, annualGroundRent, monthlyGroundRent,
		exchange, exchangePerc, installment, financing, visit,
		tenantName, tenantEmail, tenantPhone, accompanying,
		completionForecast, landBlock, landLot, landFront, landSide,
		landBack, landTerrainType, hasKmz, kmzFile, buildingFloors,
		unitTower, unitFloor, unitNumber, warehouseManufacturingArea,
		warehouseSector, warehouseHasPrimaryCabin, warehouseCabinKva,
		warehouseGroundFloor, warehouseFloorResistance, warehouseZoning,
		warehouseHasOfficeArea, warehouseOfficeArea, storeHasMezzanine,
		storeMezzanineArea, priceUpdatedAt,
		version.ID(),
	)

	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.listing.update_listing_version.exec_error", "error", execErr, "version_id", version.ID())
		return fmt.Errorf("exec update listing version: %w", execErr)
	}

	return nil
}
