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

func (la *ListingAdapter) GetListingByQuery(ctx context.Context, tx *sql.Tx, query string, args ...any) (listing listingmodel.ListingInterface, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	entityListing := listingentity.ListingEntity{}

	row := la.QueryRowContext(ctx, tx, "select", query, args...)

	if err = row.Scan(
		&entityListing.ID,
		&entityListing.UserID,
		&entityListing.Code,
		&entityListing.Version,
		&entityListing.Status,
		&entityListing.ZipCode,
		&entityListing.Street,
		&entityListing.Number,
		&entityListing.Complement,
		&entityListing.Neighborhood,
		&entityListing.City,
		&entityListing.State,
		&entityListing.Title,
		&entityListing.ListingType,
		&entityListing.Owner,
		&entityListing.LandSize,
		&entityListing.Corner,
		&entityListing.NonBuildable,
		&entityListing.Buildable,
		&entityListing.Delivered,
		&entityListing.WhoLives,
		&entityListing.Description,
		&entityListing.Transaction,
		&entityListing.SellNet,
		&entityListing.RentNet,
		&entityListing.Condominium,
		&entityListing.AnnualTax,
		&entityListing.MonthlyTax,
		&entityListing.AnnualGroundRent,
		&entityListing.MonthlyGroundRent,
		&entityListing.Exchange,
		&entityListing.ExchangePercentual,
		&entityListing.Installment,
		&entityListing.Financing,
		&entityListing.Visit,
		&entityListing.TenantName,
		&entityListing.TenantEmail,
		&entityListing.TenantPhone,
		&entityListing.Accompanying,
		&entityListing.Deleted,
	); err != nil {
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
