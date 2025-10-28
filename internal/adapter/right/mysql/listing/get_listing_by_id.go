package mysqllistingadapter

import (
	"context"
	"database/sql"
	"fmt"

	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (la *ListingAdapter) GetListingByID(ctx context.Context, tx *sql.Tx, listingID int64) (listing listingmodel.ListingInterface, err error) {
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
		annual_ground_rent,
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
	WHERE id = ? AND deleted = 0;`

	listing, err = la.GetListingByQuery(ctx, tx, query, listingID)
	if err != nil {
		return nil, fmt.Errorf("get listing by id: %w", err)
	}

	return listing, nil

}
