package storageport

import (
	"context"

	storagemodel "github.com/giulio-alfieri/toq_server/internal/core/model/storage_model"
)

// CloudStoragePortInterface define a abstração para serviços de armazenamento na nuvem
// Substitui GCSPortInterface seguindo os princípios da arquitetura hexagonal
type CloudStoragePortInterface interface {
	// Folder Operations
	CreateUserFolder(ctx context.Context, userID int64) error
	DeleteUserFolder(ctx context.Context, userID int64) error

	// Object Operations
	ListBucketObjects(ctx context.Context, bucketName string) ([]string, error)
	DeleteBucketObject(ctx context.Context, bucketName, objectName string) error
	// Object existence check
	ObjectExists(ctx context.Context, bucketName, objectName string) (bool, error)

	// Generic Signed URLs
	GenerateUploadURL(bucketName, objectName, contentType string) (string, error)
	GenerateDownloadURL(bucketName, objectName string) (string, error)

	// Domain-specific methods with proper abstraction
	GeneratePhotoUploadURL(userID int64, photoType storagemodel.PhotoType, contentType string) (string, error)
	GeneratePhotoDownloadURL(userID int64, photoType storagemodel.PhotoType) (string, error)
	GenerateDocumentUploadURL(userID int64, docType storagemodel.DocumentType, contentType string) (string, error)
	GenerateDocumentDownloadURL(userID int64, docType storagemodel.DocumentType) (string, error)

	// Configuration methods
	GetBucketConfig() storagemodel.BucketConfig
}
