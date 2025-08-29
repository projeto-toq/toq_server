package mysqllistingadapter

import (
	"context"
	"database/sql"
	"log/slog"

	listingconverters "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql/listing/converters"
	listingentity "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql/listing/entity"
	listingmodel "github.com/giulio-alfieri/toq_server/internal/core/model/listing_model"
	
	
"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (la *ListingAdapter) GetListingByQuery(ctx context.Context, tx *sql.Tx, query string, args ...any) (listing listingmodel.ListingInterface, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	entityListing := listingentity.ListingEntity{}

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		slog.Error("Error preparing statement on msqllistingadapter/GetListingByQuery", "error", err)
		err = utils.ErrInternalServer
		return
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, args...)

	err = row.Scan(
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
		&entityListing.AnnualGroundRent,
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
	)
	if err != nil {
		if err == sql.ErrNoRows {
			err = utils.ErrInternalServer
		} else {
			slog.Error("Error scanning row on msqllistingadapter/GetListingByQuery", "error", err)
			err = utils.ErrInternalServer
		}
		return
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
