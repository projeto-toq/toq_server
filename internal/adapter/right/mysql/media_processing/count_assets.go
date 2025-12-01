package mysqlmediaprocessingadapter

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	mediaprocessingrepository "github.com/projeto-toq/toq_server/internal/core/port/right/repository/media_processing_repository"
)

const countAssetsBaseQuery = `
SELECT COUNT(*)
FROM media_assets
WHERE listing_identity_id = ?
`

func (a *MediaProcessingAdapter) CountAssets(ctx context.Context, tx *sql.Tx, listingIdentityID uint64, filter mediaprocessingrepository.AssetFilter) (int64, error) {
	query := countAssetsBaseQuery
	args := []interface{}{listingIdentityID}

	if len(filter.AssetTypes) > 0 {
		placeholders := make([]string, len(filter.AssetTypes))
		for i, t := range filter.AssetTypes {
			placeholders[i] = "?"
			args = append(args, string(t))
		}
		query += fmt.Sprintf(" AND asset_type IN (%s)", strings.Join(placeholders, ","))
	}

	if len(filter.Status) > 0 {
		placeholders := make([]string, len(filter.Status))
		for i, s := range filter.Status {
			placeholders[i] = "?"
			args = append(args, string(s))
		}
		query += fmt.Sprintf(" AND status IN (%s)", strings.Join(placeholders, ","))
	}

	if filter.Sequence != nil {
		query += " AND sequence = ?"
		args = append(args, *filter.Sequence)
	}

	var count int64
	err := a.QueryRowContext(ctx, tx, "count_assets", query, args...).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
