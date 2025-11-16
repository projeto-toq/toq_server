package mysqlmediaprocessingadapter

import (
	"context"
	"database/sql"
	"strings"

	mediaprocessingconverters "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/media_processing/converters"
	mediaprocessingentities "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/media_processing/entities"
	mediaprocessingmodel "github.com/projeto-toq/toq_server/internal/core/model/media_processing_model"
	mediaprocessingrepository "github.com/projeto-toq/toq_server/internal/core/port/right/repository/mediaprocessingrepository"
)

const listBatchesBaseQuery = `
SELECT id, listing_id, reference, status, status_message, status_reason, status_details, status_updated_by, status_updated_at,
	   deleted_at
FROM listing_media_batches
WHERE listing_id = ?
`

// ListBatchesByListing retorna lotes conforme filtro fornecido.
func (a *MediaProcessingAdapter) ListBatchesByListing(ctx context.Context, tx *sql.Tx, filter mediaprocessingrepository.BatchQueryFilter) ([]mediaprocessingmodel.MediaBatch, error) {
	queryBuilder := strings.Builder{}
	queryBuilder.WriteString(listBatchesBaseQuery)
	args := []any{filter.ListingID}

	if len(filter.Statuses) > 0 {
		placeholders := make([]string, len(filter.Statuses))
		for i, status := range filter.Statuses {
			placeholders[i] = "?"
			args = append(args, status.String())
		}
		queryBuilder.WriteString(" AND status IN (" + strings.Join(placeholders, ",") + ")")
	}

	if !filter.IncludeDeleted {
		queryBuilder.WriteString(" AND deleted_at IS NULL")
	}

	queryBuilder.WriteString(" ORDER BY id DESC")
	limit := filter.Limit
	if limit <= 0 {
		limit = 20
	}
	queryBuilder.WriteString(" LIMIT ?")
	args = append(args, limit)

	query := queryBuilder.String()
	observer := a.ObserveOnComplete("select", query)
	defer observer()

	rows, err := a.QueryContext(ctx, tx, "select", query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var batches []mediaprocessingmodel.MediaBatch
	for rows.Next() {
		var entity mediaprocessingentities.BatchEntity
		if err := rows.Scan(
			&entity.ID,
			&entity.ListingID,
			&entity.Reference,
			&entity.Status,
			&entity.StatusMessage,
			&entity.StatusReason,
			&entity.StatusDetails,
			&entity.StatusUpdatedBy,
			&entity.StatusUpdatedAt,
			&entity.DeletedAt,
		); err != nil {
			return nil, err
		}
		batches = append(batches, mediaprocessingconverters.BatchEntityToDomain(entity))
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return batches, nil
}
