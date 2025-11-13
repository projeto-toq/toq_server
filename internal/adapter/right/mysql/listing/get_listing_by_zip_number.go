package mysqllistingadapter

import (
	"context"
	"database/sql"
	"fmt"

	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (la *ListingAdapter) GetListingByZipNumber(ctx context.Context, tx *sql.Tx, zip string, number string) (listing listingmodel.ListingInterface, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)

	query := `SELECT
		id,
		user_id,
		code,
		version,
		status,
		zip_code,
		street,
		number,
		complement,
		neighborhood,
		city,
		state,
		title,
		type,
		owner,
		land_size,
		corner,
		non_buildable,
		buildable,
		delivered,
		who_lives,
		description,
		transaction,
		sell_net,
		rent_net,
		condominium,
		annual_tax,
		monthly_tax,
		annual_ground_rent,
		monthly_ground_rent,
		exchange,
		exchange_perc,
		installment,
		financing,
		visit,
		tenant_name,
		tenant_email,
		tenant_phone,
		accompanying,
		deleted
	FROM listings
	WHERE zip_code = ? AND number = ? AND deleted = 0;`

	listing, err = la.GetListingByQuery(ctx, tx, query, zip, number)
	if err != nil {
		return nil, fmt.Errorf("get listing by zip number: %w", err)
	}

	return listing, nil

}
