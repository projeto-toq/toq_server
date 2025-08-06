package mysqllistingadapter

import (
	"context"
	"database/sql"

	listingmodel "github.com/giulio-alfieri/toq_server/internal/core/model/listing_model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (la *ListingAdapter) UpdateGuarantees(ctx context.Context, tx *sql.Tx, guarantees []listingmodel.GuaranteeInterface) (err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	//check if there is any data to update
	if len(guarantees) == 0 {
		return
	}

	// Remove all guarantees from listing
	err = la.DeleteListingGuarantees(ctx, tx, guarantees[0].ListingID())
	if err != nil {
		//check if the error is not found, because it's ok if there is no row to delete
		if status.Code(err) != codes.NotFound {
			return
		}
	}

	// Insert the new guarrantees
	for _, guarantee := range guarantees {
		err = la.CreateGuarantee(ctx, tx, guarantee)
		if err != nil {
			return
		}
	}

	return
}
