package mysqlvisitadapter

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// GetVisitWithListingByID returns a single visit enriched with listing snapshot and participant metadata.
func (a *VisitAdapter) GetVisitWithListingByID(ctx context.Context, tx *sql.Tx, visitID int64) (listingmodel.VisitWithListing, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return listingmodel.VisitWithListing{}, err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	scheduledStartExpr := "CAST(CONCAT(lv.scheduled_date, ' ', lv.scheduled_time_start) AS DATETIME)"
	scheduledEndExpr := "CAST(CONCAT(lv.scheduled_date, ' ', lv.scheduled_time_end) AS DATETIME)"
	baseSelect := fmt.Sprintf(visitWithListingSelectBase, scheduledStartExpr, scheduledEndExpr)
	query := fmt.Sprintf(`%s WHERE lv.id = ? LIMIT 1`, baseSelect)

	row := a.QueryRowContext(ctx, tx, "get_visit_with_listing", query, visitID)
	entry, scanErr := scanVisitWithListingRow(row)
	if scanErr != nil {
		if errors.Is(scanErr, sql.ErrNoRows) {
			return listingmodel.VisitWithListing{}, sql.ErrNoRows
		}
		utils.SetSpanError(ctx, scanErr)
		logger.Error("mysql.visit.get_with_listing.scan_error", "visit_id", visitID, "err", scanErr)
		return listingmodel.VisitWithListing{}, fmt.Errorf("get visit with listing: %w", scanErr)
	}

	return entry, nil
}
