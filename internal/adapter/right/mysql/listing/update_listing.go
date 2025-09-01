package mysqllistingadapter

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	listingmodel "github.com/giulio-alfieri/toq_server/internal/core/model/listing_model"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (la *ListingAdapter) UpdateListing(ctx context.Context, tx *sql.Tx, listing listingmodel.ListingInterface) (err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	sql := `UPDATE listings SET
				user_id = ?, code = ?, version = ?, status = ?, zip_code = ?, street = ?, number = ?, complement = ?, neighborhood = ?, city = ?, state = ?,
				type = ?, owner = ?, land_size = ?, corner = ?, non_buildable = ?, buildable = ?, delivered = ?, who_lives = ?, description = ?,
				transaction = ?, sell_net = ?, rent_net = ?, condominium = ?, annual_tax = ?, annual_ground_rent = ?, exchange = ?, exchange_perc = ?,
				installment = ?, financing = ?, visit = ?, tenant_name = ?, tenant_email = ?, tenant_phone = ?, accompanying = ?, deleted = ?
			WHERE id = ?`

	stmt, err := tx.PrepareContext(ctx, sql)
	if err != nil {
		slog.Error("mysqllistingadapter/UpdateListing: error preparing statement", "error", err)
		err = fmt.Errorf("prepare update listing: %w", err)
		return
	}

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
		slog.Error("mysqllistingadapter/UpdateListing: error executing statement", "error", err)
		err = fmt.Errorf("exec update listing: %w", err)
		return
	}

	err = la.UpdateExchangePlaces(ctx, tx, listing.ExchangePlaces())
	if err != nil {
		return
	}
	err = la.UpdateFeatures(ctx, tx, listing.Features())
	if err != nil {
		return
	}
	err = la.UpdateGuarantees(ctx, tx, listing.Guarantees())
	if err != nil {
		return
	}
	err = la.UpdateFinancingBlockers(ctx, tx, listing.FinancingBlockers())
	if err != nil {
		return
	}

	return
}
