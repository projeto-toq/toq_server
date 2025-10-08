package mysqllistingadapter

import (
	"context"
	"database/sql"
	"fmt"

	listingmodel "github.com/giulio-alfieri/toq_server/internal/core/model/listing_model"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
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
				user_id = ?, code = ?, version = ?, status = ?, zip_code = ?, street = ?, number = ?, complement = ?, neighborhood = ?, city = ?, state = ?,
				type = ?, owner = ?, land_size = ?, corner = ?, non_buildable = ?, buildable = ?, delivered = ?, who_lives = ?, description = ?,
				transaction = ?, sell_net = ?, rent_net = ?, condominium = ?, annual_tax = ?, annual_ground_rent = ?, exchange = ?, exchange_perc = ?,
				installment = ?, financing = ?, visit = ?, tenant_name = ?, tenant_email = ?, tenant_phone = ?, accompanying = ?, deleted = ?
			WHERE id = ?`

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.listing.update_listing.prepare_error", "error", err, "listing_id", listing.ID())
		return fmt.Errorf("prepare update listing: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx,
		listing.UserID(), listing.Code(), listing.Version(), listing.Status(), listing.ZipCode(),
		listing.ToSQLNullString(listing.Street()),
		listing.ToSQLNullString(listing.Number()),
		listing.ToSQLNullString(listing.Complement()),
		listing.ToSQLNullString(listing.Neighborhood()),
		listing.ToSQLNullString(listing.City()),
		listing.ToSQLNullString(listing.State()),
		listing.ToSQLNullInt(uint8(listing.ListingType())),
		listing.ToSQLNullInt(uint8(listing.Owner())),
		listing.LandSize(), listing.Corner(), listing.NonBuildable(), listing.Buildable(), listing.Delivered(), listing.WhoLives(), listing.Description(),
		listing.Transaction(), listing.SellNet(), listing.RentNet(), listing.Condominium(), listing.AnnualTax(), listing.AnnualGroundRent(), listing.Exchange(),
		listing.ExchangePercentual(), listing.Installment(), listing.Financing(), listing.Visit(), listing.TenantName(), listing.TenantEmail(),
		listing.TenantPhone(), listing.Accompanying(), listing.Deleted(),
		listing.ID())
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.listing.update_listing.exec_error", "error", err, "listing_id", listing.ID())
		return fmt.Errorf("exec update listing: %w", err)
	}

	err = la.UpdateExchangePlaces(ctx, tx, listing.ExchangePlaces())
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.listing.update_listing.exchange_places_error", "error", err, "listing_id", listing.ID())
		return fmt.Errorf("update exchange places: %w", err)
	}
	err = la.UpdateFeatures(ctx, tx, listing.Features())
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.listing.update_listing.features_error", "error", err, "listing_id", listing.ID())
		return fmt.Errorf("update features: %w", err)
	}
	err = la.UpdateGuarantees(ctx, tx, listing.Guarantees())
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.listing.update_listing.guarantees_error", "error", err, "listing_id", listing.ID())
		return fmt.Errorf("update guarantees: %w", err)
	}
	err = la.UpdateFinancingBlockers(ctx, tx, listing.FinancingBlockers())
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.listing.update_listing.financing_blockers_error", "error", err, "listing_id", listing.ID())
		return fmt.Errorf("update financing blockers: %w", err)
	}

	return nil
}
