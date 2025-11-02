package mysqllistingadapter

import (
	"context"
	"database/sql"
	"fmt"

	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (la *ListingAdapter) UpdateListing(ctx context.Context, tx *sql.Tx, listing listingmodel.ListingInterface) (err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `UPDATE listings SET
			user_id = ?, code = ?, version = ?, status = ?, zip_code = ?, street = ?, number = ?, complement = ?, neighborhood = ?, city = ?, state = ?, title = ?,
				type = ?, owner = ?, land_size = ?, corner = ?, non_buildable = ?, buildable = ?, delivered = ?, who_lives = ?, description = ?,
				transaction = ?, sell_net = ?, rent_net = ?, condominium = ?, annual_tax = ?, annual_ground_rent = ?, exchange = ?, exchange_perc = ?,
				installment = ?, financing = ?, visit = ?, tenant_name = ?, tenant_email = ?, tenant_phone = ?, accompanying = ?, deleted = ?
			WHERE id = ?`

	street := sql.NullString{String: listing.Street(), Valid: listing.Street() != ""}
	number := sql.NullString{String: listing.Number(), Valid: true}
	complement := sql.NullString{String: listing.Complement(), Valid: listing.Complement() != ""}
	neighborhood := sql.NullString{String: listing.Neighborhood(), Valid: listing.Neighborhood() != ""}
	city := sql.NullString{String: listing.City(), Valid: listing.City() != ""}
	state := sql.NullString{String: listing.State(), Valid: listing.State() != ""}
	title := sql.NullString{}
	if listing.HasTitle() {
		title = sql.NullString{String: listing.Title(), Valid: true}
	}
	typeValue := sql.NullInt64{Int64: int64(listing.ListingType()), Valid: true}
	owner := sql.NullInt64{}
	if listing.HasOwner() {
		owner = sql.NullInt64{Int64: int64(listing.Owner()), Valid: true}
	}
	landSize := sql.NullFloat64{}
	if listing.HasLandSize() {
		landSize = sql.NullFloat64{Float64: listing.LandSize(), Valid: true}
	}
	corner := sql.NullBool{}
	if listing.HasCorner() {
		corner = sql.NullBool{Bool: listing.Corner(), Valid: true}
	}
	nonBuildable := sql.NullFloat64{}
	if listing.HasNonBuildable() {
		nonBuildable = sql.NullFloat64{Float64: listing.NonBuildable(), Valid: true}
	}
	buildable := sql.NullFloat64{}
	if listing.HasBuildable() {
		buildable = sql.NullFloat64{Float64: listing.Buildable(), Valid: true}
	}
	delivered := sql.NullInt64{}
	if listing.HasDelivered() {
		delivered = sql.NullInt64{Int64: int64(listing.Delivered()), Valid: true}
	}
	whoLives := sql.NullInt64{}
	if listing.HasWhoLives() {
		whoLives = sql.NullInt64{Int64: int64(listing.WhoLives()), Valid: true}
	}
	description := sql.NullString{}
	if listing.HasDescription() {
		description = sql.NullString{String: listing.Description(), Valid: true}
	}
	transaction := sql.NullInt64{}
	if listing.HasTransaction() {
		transaction = sql.NullInt64{Int64: int64(listing.Transaction()), Valid: true}
	}
	sellNet := sql.NullFloat64{}
	if listing.HasSellNet() {
		sellNet = sql.NullFloat64{Float64: listing.SellNet(), Valid: true}
	}
	rentNet := sql.NullFloat64{}
	if listing.HasRentNet() {
		rentNet = sql.NullFloat64{Float64: listing.RentNet(), Valid: true}
	}
	condominium := sql.NullFloat64{}
	if listing.HasCondominium() {
		condominium = sql.NullFloat64{Float64: listing.Condominium(), Valid: true}
	}
	annualTax := sql.NullFloat64{}
	if listing.HasAnnualTax() {
		annualTax = sql.NullFloat64{Float64: listing.AnnualTax(), Valid: true}
	}
	annualGroundRent := sql.NullFloat64{}
	if listing.HasAnnualGroundRent() {
		annualGroundRent = sql.NullFloat64{Float64: listing.AnnualGroundRent(), Valid: true}
	}
	exchange := sql.NullBool{}
	if listing.HasExchange() {
		exchange = sql.NullBool{Bool: listing.Exchange(), Valid: true}
	}
	exchangePercentual := sql.NullFloat64{}
	if listing.HasExchangePercentual() {
		exchangePercentual = sql.NullFloat64{Float64: listing.ExchangePercentual(), Valid: true}
	}
	installment := sql.NullInt64{}
	if listing.HasInstallment() {
		installment = sql.NullInt64{Int64: int64(listing.Installment()), Valid: true}
	}
	financing := sql.NullBool{}
	if listing.HasFinancing() {
		financing = sql.NullBool{Bool: listing.Financing(), Valid: true}
	}
	visit := sql.NullInt64{}
	if listing.HasVisit() {
		visit = sql.NullInt64{Int64: int64(listing.Visit()), Valid: true}
	}
	tenantName := sql.NullString{}
	if listing.HasTenantName() {
		tenantName = sql.NullString{String: listing.TenantName(), Valid: true}
	}
	tenantEmail := sql.NullString{}
	if listing.HasTenantEmail() {
		tenantEmail = sql.NullString{String: listing.TenantEmail(), Valid: true}
	}
	tenantPhone := sql.NullString{}
	if listing.HasTenantPhone() {
		tenantPhone = sql.NullString{String: listing.TenantPhone(), Valid: true}
	}
	accompanying := sql.NullInt64{}
	if listing.HasAccompanying() {
		accompanying = sql.NullInt64{Int64: int64(listing.Accompanying()), Valid: true}
	}
	deletedValue := listing.Deleted()

	if _, execErr := la.ExecContext(ctx, tx, "update", query,
		listing.UserID(), listing.Code(), listing.Version(), listing.Status(), listing.ZipCode(),
		street,
		number,
		complement,
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
		annualGroundRent,
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
		listing.ID()); execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.listing.update_listing.exec_error", "error", execErr, "listing_id", listing.ID())
		return fmt.Errorf("exec update listing: %w", execErr)
	}

	err = la.UpdateExchangePlaces(ctx, tx, listing.ID(), listing.ExchangePlaces())
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.listing.update_listing.exchange_places_error", "error", err, "listing_id", listing.ID())
		return fmt.Errorf("update exchange places: %w", err)
	}
	err = la.UpdateFeatures(ctx, tx, listing.ID(), listing.Features())
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.listing.update_listing.features_error", "error", err, "listing_id", listing.ID())
		return fmt.Errorf("update features: %w", err)
	}
	err = la.UpdateGuarantees(ctx, tx, listing.ID(), listing.Guarantees())
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.listing.update_listing.guarantees_error", "error", err, "listing_id", listing.ID())
		return fmt.Errorf("update guarantees: %w", err)
	}
	err = la.UpdateFinancingBlockers(ctx, tx, listing.ID(), listing.FinancingBlockers())
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.listing.update_listing.financing_blockers_error", "error", err, "listing_id", listing.ID())
		return fmt.Errorf("update financing blockers: %w", err)
	}

	return nil
}
