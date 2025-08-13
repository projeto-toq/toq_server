package gcsport

import "context"

type GCSPortInterface interface {
	CreateUserBucket(ctx context.Context, UserID int64) (err error)
	DeleteUserBucket(ctx context.Context, UserID int64) (err error)
	ListBucketObjects(ctx context.Context, bucketName string) (objects []string, err error)
	DeleteBucketObject(ctx context.Context, bucketName, objectName string) (err error)
	GenerateV4PutObjectSignedURL(bucketName, objectName, contentType string) (string, error)
	GenerateV4GetObjectSignedURL(bucketName, objectName string) (string, error)
}
