package mysqllistingfavoriteadapter

import (
	"context"
	"database/sql"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// ListByUser returns paginated listing identity IDs favorited by the user and the total count.
func (a *ListingFavoriteAdapter) ListByUser(ctx context.Context, tx *sql.Tx, userID int64, page, limit int) ([]int64, int64, error) {
	ctx, spanEnd, _ := utils.GenerateTracer(ctx)
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}
	offset := (page - 1) * limit

	countQuery := `SELECT COUNT(*) FROM listing_favorites WHERE user_id = ?`
	var total int64
	if err := a.QueryRowContext(ctx, tx, "count", countQuery, userID).Scan(&total); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.listing_favorite.list.count_error", "user_id", userID, "err", err)
		return nil, 0, err
	}

	if total == 0 {
		return []int64{}, 0, nil
	}

	query := `SELECT listing_identity_id FROM listing_favorites WHERE user_id = ? ORDER BY id DESC LIMIT ? OFFSET ?`
	rows, err := a.QueryContext(ctx, tx, "select", query, userID, limit, offset)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.listing_favorite.list.query_error", "user_id", userID, "err", err)
		return nil, 0, err
	}
	defer rows.Close()

	ids := make([]int64, 0)
	for rows.Next() {
		var id int64
		if scanErr := rows.Scan(&id); scanErr != nil {
			utils.SetSpanError(ctx, scanErr)
			logger.Error("mysql.listing_favorite.list.scan_error", "user_id", userID, "err", scanErr)
			return nil, 0, scanErr
		}
		ids = append(ids, id)
	}

	if err = rows.Err(); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.listing_favorite.list.rows_error", "user_id", userID, "err", err)
		return nil, 0, err
	}

	return ids, total, nil
}
