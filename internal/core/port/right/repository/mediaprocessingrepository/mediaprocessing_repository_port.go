package mediaprocessingrepository

import (
	"context"
	"database/sql"

	mediaprocessingmodel "github.com/projeto-toq/toq_server/internal/core/model/media_processing_model"
)

// RepositoryInterface exposes the persistence contract for media processing entities.
type RepositoryInterface interface {
	CreateBatch(ctx context.Context, tx *sql.Tx, batch mediaprocessingmodel.MediaBatch) (uint64, error)
	UpdateBatchStatus(ctx context.Context, tx *sql.Tx, batchID uint64, status mediaprocessingmodel.BatchStatus, metadata mediaprocessingmodel.BatchStatusMetadata) error
	GetBatchByID(ctx context.Context, tx *sql.Tx, batchID uint64) (mediaprocessingmodel.MediaBatch, error)
	ListBatchesByListing(ctx context.Context, tx *sql.Tx, filter BatchQueryFilter) ([]mediaprocessingmodel.MediaBatch, error)

	UpsertAssets(ctx context.Context, tx *sql.Tx, assets []mediaprocessingmodel.MediaAsset) error
	ListAssetsByBatch(ctx context.Context, tx *sql.Tx, batchID uint64) ([]mediaprocessingmodel.MediaAsset, error)

	RegisterProcessingJob(ctx context.Context, tx *sql.Tx, job mediaprocessingmodel.MediaProcessingJob) (uint64, error)
	GetProcessingJobByID(ctx context.Context, tx *sql.Tx, jobID uint64) (mediaprocessingmodel.MediaProcessingJob, error)
	UpdateProcessingJob(ctx context.Context, tx *sql.Tx, jobID uint64, status mediaprocessingmodel.MediaProcessingJobStatus, output mediaprocessingmodel.MediaProcessingJobPayload) error

	SoftDeleteBatch(ctx context.Context, tx *sql.Tx, batchID uint64) error
}

// BatchQueryFilter narrows down batch lookups for handlers and workers.
type BatchQueryFilter struct {
	ListingID      uint64
	Statuses       []mediaprocessingmodel.BatchStatus
	Limit          int
	IncludeDeleted bool
}
