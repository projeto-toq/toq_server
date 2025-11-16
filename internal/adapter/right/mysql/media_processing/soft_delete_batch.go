package mysqlmediaprocessingadapter

import (
	"context"
	"database/sql"
	"time"
)

const softDeleteBatchQuery = `
UPDATE listing_media_batches
SET deleted_at = ?
WHERE id = ? AND deleted_at IS NULL
`

// SoftDeleteBatch aplica soft delete no lote indicado.
func (a *MediaProcessingAdapter) SoftDeleteBatch(ctx context.Context, tx *sql.Tx, batchID uint64) error {
	now := time.Now()
	observer := a.ObserveOnComplete("update", softDeleteBatchQuery)
	defer observer()

	_, err := a.ExecContext(ctx, tx, "update", softDeleteBatchQuery,
		now,
		batchID,
	)
	return err
}
