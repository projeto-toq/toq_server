package mysqllistingfavoriteadapter

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// GetUserFlags returns which listing identities are favorited by the given user.
func (a *ListingFavoriteAdapter) GetUserFlags(ctx context.Context, tx *sql.Tx, listingIdentityIDs []int64, userID int64) (map[int64]bool, error) {
	ctx, spanEnd, _ := utils.GenerateTracer(ctx)
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	flags := make(map[int64]bool)
	if len(listingIdentityIDs) == 0 {
		return flags, nil
	}

	placeholders := make([]string, len(listingIdentityIDs))
	args := make([]any, 0, len(listingIdentityIDs)+1)
	args = append(args, userID)
	for i, id := range listingIdentityIDs {
		placeholders[i] = "?"
		args = append(args, id)
	}

	query := fmt.Sprintf(`SELECT listing_identity_id FROM listing_favorites WHERE user_id = ? AND listing_identity_id IN (%s)`, strings.Join(placeholders, ","))

	rows, err := a.QueryContext(ctx, tx, "select", query, args...)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.listing_favorite.flags.query_error", "user_id", userID, "err", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id int64
		if scanErr := rows.Scan(&id); scanErr != nil {
			utils.SetSpanError(ctx, scanErr)
			logger.Error("mysql.listing_favorite.flags.scan_error", "user_id", userID, "err", scanErr)
			return nil, scanErr
		}
		flags[id] = true
	}

	if err = rows.Err(); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.listing_favorite.flags.rows_error", "user_id", userID, "err", err)
		return nil, err
	}

	return flags, nil
}
