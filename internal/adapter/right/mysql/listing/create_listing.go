package mysqllistingadapter

import (
	"context"
	"database/sql"
	"fmt"

	listingmodel "github.com/giulio-alfieri/toq_server/internal/core/model/listing_model"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (la *ListingAdapter) CreateListing(ctx context.Context, tx *sql.Tx, listing listingmodel.ListingInterface) (err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	sql := `INSERT INTO listings (
				user_id, code, version, status, zip_code, street, number, complement, neighborhood, city, state,
				type, owner, land_size, corner, non_buildable, buildable, delivered, who_lives, description,
				transaction, sell_net, rent_net, condominium, annual_tax, annual_ground_rent, exchange, exchange_perc,
				installment, financing, visit, tenant_name, tenant_email, tenant_phone, accompanying, deleted)
			VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)` //36 ?

	stmt, err := tx.PrepareContext(ctx, sql)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.listing.create_listing.prepare_error", "error", err)
		return fmt.Errorf("prepare create listing: %w", err)
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx,
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
		listing.TenantPhone(), listing.Accompanying(), listing.Deleted())
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.listing.create_listing.exec_error", "error", err)
		return fmt.Errorf("exec create listing: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.listing.create_listing.last_insert_error", "error", err)
		return fmt.Errorf("last insert id for create listing: %w", err)
	}

	listing.SetID(id)

	return nil
}
