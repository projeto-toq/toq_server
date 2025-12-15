package media_processing_repository

import (
	"context"
	"database/sql"

	mediaprocessingmodel "github.com/projeto-toq/toq_server/internal/core/model/media_processing_model"
)

// RepositoryInterface exposes the persistence contract for media processing entities.
type RepositoryInterface interface {
	// New methods for granular asset management
	UpsertAsset(ctx context.Context, tx *sql.Tx, asset mediaprocessingmodel.MediaAsset) error
	GetAsset(ctx context.Context, tx *sql.Tx, listingIdentityID uint64, assetType mediaprocessingmodel.MediaAssetType, sequence uint8) (mediaprocessingmodel.MediaAsset, error)
	GetAssetBySequence(ctx context.Context, tx *sql.Tx, listingIdentityID uint64, assetType mediaprocessingmodel.MediaAssetType, sequence uint8) (mediaprocessingmodel.MediaAsset, error)
	GetAssetByID(ctx context.Context, tx *sql.Tx, assetID uint64) (mediaprocessingmodel.MediaAsset, error)
	GetAssetByRawKey(ctx context.Context, tx *sql.Tx, rawKey string) (mediaprocessingmodel.MediaAsset, error)
	ListAssets(ctx context.Context, tx *sql.Tx, listingIdentityID uint64, filter AssetFilter, pagination *Pagination) ([]mediaprocessingmodel.MediaAsset, error)
	CountAssets(ctx context.Context, tx *sql.Tx, listingIdentityID uint64, filter AssetFilter) (int64, error)
	DeleteAsset(ctx context.Context, tx *sql.Tx, listingIdentityID uint64, assetType mediaprocessingmodel.MediaAssetType, sequence uint8) error
	BulkUpdateAssetStatus(ctx context.Context, tx *sql.Tx, listingIdentityID uint64, fromStatus, toStatus mediaprocessingmodel.MediaAssetStatus) error

	RegisterProcessingJob(ctx context.Context, tx *sql.Tx, job mediaprocessingmodel.MediaProcessingJob) (uint64, error)
	GetProcessingJobByID(ctx context.Context, tx *sql.Tx, jobID uint64) (mediaprocessingmodel.MediaProcessingJob, error)
	UpdateProcessingJob(ctx context.Context, tx *sql.Tx, job mediaprocessingmodel.MediaProcessingJob) error
	GetLatestFinalizationJob(ctx context.Context, tx *sql.Tx, listingIdentityID uint64) (mediaprocessingmodel.MediaProcessingJob, error)
}

// Pagination define opções de paginação para repositórios.
type Pagination struct {
	Page  int
	Limit int
	Sort  string
	Order string
}

// AssetFilter narrows down asset lookups.
type AssetFilter struct {
	AssetTypes []mediaprocessingmodel.MediaAssetType
	Status     []mediaprocessingmodel.MediaAssetStatus
	Sequence   *uint8
}
