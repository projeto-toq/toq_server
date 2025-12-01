package mysqlmediaprocessingadapter

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	mediaprocessingconverters "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/media_processing/converters"
	mediaprocessingentities "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/media_processing/entities"
	mediaprocessingmodel "github.com/projeto-toq/toq_server/internal/core/model/media_processing_model"
	mediaprocessingrepository "github.com/projeto-toq/toq_server/internal/core/port/right/repository/mediaprocessingrepository"
)

const listAssetsBaseQuery = `
SELECT
    id, listing_identity_id, asset_type, sequence, status, s3_key_raw, s3_key_processed, title, metadata
FROM media_assets
WHERE listing_identity_id = ?
`

// ListAssets retrieves assets for a listing with optional filters.
func (a *MediaProcessingAdapter) ListAssets(ctx context.Context, tx *sql.Tx, listingIdentityID uint64, filter mediaprocessingrepository.AssetFilter) ([]mediaprocessingmodel.MediaAsset, error) {
	query := listAssetsBaseQuery
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

	query += " ORDER BY asset_type, sequence"

	rows, err := a.QueryContext(ctx, tx, "list_assets", query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var assets []mediaprocessingmodel.MediaAsset
	for rows.Next() {
		var entity mediaprocessingentities.AssetEntity
		if err := rows.Scan(
			&entity.ID,
			&entity.ListingIdentityID,
			&entity.AssetType,
			&entity.Sequence,
			&entity.Status,
			&entity.S3KeyRaw,
			&entity.S3KeyProcessed,
			&entity.Title,
			&entity.Metadata,
		); err != nil {
			return nil, err
		}
		assets = append(assets, mediaprocessingconverters.AssetEntityToDomain(entity))
	}

	return assets, nil
}
