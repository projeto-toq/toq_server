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
		lv.deleted`

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

	listing = listingconverters.ListingEntityToDomain(entityListing)
	return
}
