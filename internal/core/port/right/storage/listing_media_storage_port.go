package storageport

import (
	"context"
	"time"

	mediaprocessingmodel "github.com/projeto-toq/toq_server/internal/core/model/media_processing_model"
)

// ListingMediaStoragePort exposes the contract with the storage provider used for raw/processed media.
type ListingMediaStoragePort interface {
	GenerateRawUploadURL(ctx context.Context, listingID uint64, asset mediaprocessingmodel.MediaAsset, contentType, checksum string) (SignedURL, error)
	GenerateProcessedDownloadURL(ctx context.Context, listingID uint64, asset mediaprocessingmodel.MediaAsset, resolution string) (SignedURL, error)
	GenerateDownloadURL(ctx context.Context, key string) (SignedURL, error)
	ValidateObjectChecksum(ctx context.Context, bucketKey string, expectedChecksum string) (StorageObjectMetadata, error)
	DeleteObject(ctx context.Context, bucketKey string) error
	DeleteKeys(ctx context.Context, keys []string) error
	DownloadFile(ctx context.Context, key string) ([]byte, error)
	UploadFile(ctx context.Context, key string, content []byte, contentType string) error
}

// SignedURL stores the URL plus mandatory headers for clients.
type SignedURL struct {
	URL       string
	Method    string
	Headers   map[string]string
	ExpiresIn time.Duration
	ObjectKey string
}

// StorageObjectMetadata captures attributes read from S3/compatible providers.
type StorageObjectMetadata struct {
	SizeInBytes  int64
	Checksum     string
	ETag         string
	ContentType  string
	LastModified time.Time
}
