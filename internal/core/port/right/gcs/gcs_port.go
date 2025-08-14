package gcsport

import "context"

type GCSPortInterface interface {
	CreateUserFolder(ctx context.Context, UserID int64) (err error)
	DeleteUserFolder(ctx context.Context, UserID int64) (err error)
	ListBucketObjects(ctx context.Context, bucketName string) (objects []string, err error)
	DeleteBucketObject(ctx context.Context, bucketName, objectName string) (err error)
	GenerateV4PutObjectSignedURL(bucketName, objectName, contentType string) (string, error)
	GenerateV4GetObjectSignedURL(bucketName, objectName string) (string, error)
	GeneratePhotoSignedURL(bucketName string, userID int64, photoType, contentType string) (string, error)
	GeneratePhotoDownloadURL(bucketName string, userID int64, photoType string) (string, error)
}
