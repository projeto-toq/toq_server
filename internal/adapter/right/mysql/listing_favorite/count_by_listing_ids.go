package mysqllistingfavoriteadapter

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// CountByListingIdentities returns favorites counts for the provided listing identity IDs.
func (a *ListingFavoriteAdapter) CountByListingIdentities(ctx context.Context, tx *sql.Tx, listingIdentityIDs []int64) (map[int64]int64, error) {
	ctx, spanEnd, _ := utils.GenerateTracer(ctx)
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	result := make(map[int64]int64)
	if len(listingIdentityIDs) == 0 {
		return result, nil
	}

	placeholders := make([]string, len(listingIdentityIDs))
	args := make([]any, 0, len(listingIdentityIDs))
	for i, id := range listingIdentityIDs {
		placeholders[i] = "?"
		args = append(args, id)
	}

	query := fmt.Sprintf(`SELECT listing_identity_id, COUNT(*) FROM listing_favorites WHERE listing_identity_id IN (%s) GROUP BY listing_identity_id`, strings.Join(placeholders, ","))

	rows, err := a.QueryContext(ctx, tx, "select", query, args...)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.listing_favorite.count.query_error", "ids", listingIdentityIDs, "err", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var listingID int64
		var count int64
		if scanErr := rows.Scan(&listingID, &count); scanErr != nil {
			utils.SetSpanError(ctx, scanErr)
			logger.Error("mysql.listing_favorite.count.scan_error", "err", scanErr)
			return nil, scanErr
		}
		result[listingID] = count
	}

	if err = rows.Err(); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.listing_favorite.count.rows_error", "err", err)
		return nil, err
	}

	return result, nil
}
