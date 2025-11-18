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

func (la *ListingAdapter) GetListingByQuery(ctx context.Context, tx *sql.Tx, query string, args ...any) (listing listingmodel.ListingInterface, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	row := la.QueryRowContext(ctx, tx, "select", query, args...)
	entityListing, err := scanListingEntity(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.listing.get_listing_by_query.scan_error", "error", err)
		return nil, fmt.Errorf("scan listing by query: %w", err)
	}

	entityExchangePlaces, err := la.GetEntityExchangePlacesByListing(ctx, tx, entityListing.ID)
	if err != nil {
		return
	}
	entityListing.ExchangePlaces = entityExchangePlaces

	entityFeatures, err := la.GetEntityFeaturesByListing(ctx, tx, entityListing.ID)
	if err != nil {
		return
	}
	entityListing.Features = entityFeatures

	entityGuarantees, err := la.GetEntityGuaranteesByListing(ctx, tx, entityListing.ID)
	if err != nil {
		return
	}
	entityListing.Guarantees = entityGuarantees

	entityFinancingBlockers, err := la.GetEntityFinancingBlockersByListing(ctx, tx, entityListing.ID)
	if err != nil {
		return
	}
	entityListing.FinancingBlocker = entityFinancingBlockers

	entityWarehouseFloors, err := la.GetEntityWarehouseAdditionalFloorsByListing(ctx, tx, entityListing.ID)
	if err != nil {
		return
	}
	entityListing.WarehouseAdditionalFloors = entityWarehouseFloors

	listing = listingconverters.ListingEntityToDomain(entityListing)
	return
}
