package port

import (
	"context"
	"io"
)

// StoragePort defines the interface for object storage operations (e.g., S3)
type StoragePort interface {
	// Download retrieves an object from storage
	Download(ctx context.Context, bucket, key string) (io.ReadCloser, error)

	// Upload stores an object and returns its location/error
	Upload(ctx context.Context, bucket, key string, body io.Reader, contentType string) error
}
