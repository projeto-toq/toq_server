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

func (la *ListingAdapter) GetListingForEndUpdate(ctx context.Context, tx *sql.Tx, listingID int64) (listingrepository.ListingEndUpdateData, error) {
	data := listingrepository.ListingEndUpdateData{}

	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return data, err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `SELECT id, user_id, code, version, status, zip_code, street, number, city, state, title, type, owner,
		buildable, delivered, who_lives, description, transaction, visit, accompanying, annual_tax, monthly_tax, 
		annual_ground_rent, monthly_ground_rent, exchange,
		exchange_perc, sell_net, rent_net, condominium, land_size, corner, tenant_name, tenant_phone, tenant_email,
		financing FROM listings WHERE id = ?`

	var (
		status            uint8
		listingType       uint16
		street            sql.NullString
		number            sql.NullString
		city              sql.NullString
		state             sql.NullString
		title             sql.NullString
		owner             sql.NullInt16
		buildable         sql.NullFloat64
		delivered         sql.NullInt16
		whoLives          sql.NullInt16
		description       sql.NullString
		transaction       sql.NullInt16
		visit             sql.NullInt16
		accompanying      sql.NullInt16
		annualTax         sql.NullFloat64
		monthlyTax        sql.NullFloat64
		annualGroundRent  sql.NullFloat64
		monthlyGroundRent sql.NullFloat64
		exchange          sql.NullInt16
		exchangePerc      sql.NullFloat64
		saleNet           sql.NullFloat64
		rentNet           sql.NullFloat64
		condominium       sql.NullFloat64
		landSize          sql.NullFloat64
		corner            sql.NullInt16
		tenantName        sql.NullString
		tenantPhone       sql.NullString
		tenantEmail       sql.NullString
		financing         sql.NullInt16
	)

	row := la.QueryRowContext(ctx, tx, "select", query, listingID)
	err = row.Scan(
		&data.ListingID,
		&data.UserID,
		&data.Code,
		&data.Version,
		&status,
		&data.ZipCode,
		&street,
		&number,
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
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return data, sql.ErrNoRows
		}
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.listing.get_listing_for_end_update.scan_error", "error", err, "listing_id", listingID)
		return data, fmt.Errorf("scan listing for end update: %w", err)
	}

	data.Status = listingmodel.ListingStatus(status)
	data.ListingType = globalmodel.PropertyType(listingType)
	data.Street = street
	data.Number = number
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

	features, ferr := la.GetEntityFeaturesByListing(ctx, tx, data.ListingID)
	if ferr != nil && !errors.Is(ferr, sql.ErrNoRows) {
		utils.SetSpanError(ctx, ferr)
		logger.Error("mysql.listing.get_listing_for_end_update.features_error", "error", ferr, "listing_id", listingID)
		return data, fmt.Errorf("get features for end update: %w", ferr)
	}
	data.FeaturesCount = len(features)

	exchangePlaces, eerr := la.GetEntityExchangePlacesByListing(ctx, tx, data.ListingID)
	if eerr != nil && !errors.Is(eerr, sql.ErrNoRows) {
		utils.SetSpanError(ctx, eerr)
		logger.Error("mysql.listing.get_listing_for_end_update.exchange_places_error", "error", eerr, "listing_id", listingID)
		return data, fmt.Errorf("get exchange places for end update: %w", eerr)
	}
	data.ExchangePlacesCount = len(exchangePlaces)

	financingBlockers, berr := la.GetEntityFinancingBlockersByListing(ctx, tx, data.ListingID)
	if berr != nil && !errors.Is(berr, sql.ErrNoRows) {
		utils.SetSpanError(ctx, berr)
		logger.Error("mysql.listing.get_listing_for_end_update.financing_blockers_error", "error", berr, "listing_id", listingID)
		return data, fmt.Errorf("get financing blockers for end update: %w", berr)
	}
	data.FinancingBlockersCount = len(financingBlockers)

	guarantees, gerr := la.GetEntityGuaranteesByListing(ctx, tx, data.ListingID)
	if gerr != nil && !errors.Is(gerr, sql.ErrNoRows) {
		utils.SetSpanError(ctx, gerr)
		logger.Error("mysql.listing.get_listing_for_end_update.guarantees_error", "error", gerr, "listing_id", listingID)
		return data, fmt.Errorf("get guarantees for end update: %w", gerr)
	}
	data.GuaranteesCount = len(guarantees)

	return data, nil
}
