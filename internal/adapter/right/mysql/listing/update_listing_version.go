package mysqllistingadapter

import (
	"context"
	"database/sql"
	"fmt"

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
			status = ?, title = ?, zip_code = ?, street = ?, number = ?, complement = ?,
			neighborhood = ?, city = ?, state = ?, type = ?, owner = ?, land_size = ?,
			corner = ?, non_buildable = ?, buildable = ?, delivered = ?, who_lives = ?,
			description = ?, transaction = ?, sell_net = ?, rent_net = ?, condominium = ?,
			annual_tax = ?, monthly_tax = ?, annual_ground_rent = ?, monthly_ground_rent = ?,
			exchange = ?, exchange_perc = ?, installment = ?, financing = ?, visit = ?,
			tenant_name = ?, tenant_email = ?, tenant_phone = ?, accompanying = ?
		WHERE id = ? AND deleted = 0
	`

	var title, description, tenantName, tenantEmail, tenantPhone interface{}
	var owner, landSize, corner, nonBuildable, buildable, delivered, whoLives, transaction interface{}
	var sellNet, rentNet, condominium, annualTax, monthlyTax, annualGroundRent, monthlyGroundRent interface{}
	var exchange, exchangePerc, installment, financing, visit, accompanying interface{}

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

	_, execErr := la.ExecContext(ctx, tx, "update", query,
		uint8(version.Status()), title, version.ZipCode(), street, version.Number(), complement,
		neighborhood, city, state, uint8(version.ListingType()), owner, landSize,
		corner, nonBuildable, buildable, delivered, whoLives,
		description, transaction, sellNet, rentNet, condominium,
		annualTax, monthlyTax, annualGroundRent, monthlyGroundRent,
		exchange, exchangePerc, installment, financing, visit,
		tenantName, tenantEmail, tenantPhone, accompanying,
		version.ID(),
	)

	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.listing.update_listing_version.exec_error", "error", execErr, "version_id", version.ID())
		return fmt.Errorf("exec update listing version: %w", execErr)
	}

	return nil
}
