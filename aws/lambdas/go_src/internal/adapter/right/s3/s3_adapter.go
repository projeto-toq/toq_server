package s3adapter

import (
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/projeto-toq/toq_server/aws/lambdas/go_src/internal/core/port"
)

// S3Adapter implements StoragePort using AWS SDK v2
type S3Adapter struct {
	client   *s3.Client
	uploader *manager.Uploader
}

// NewS3Adapter creates a new S3 adapter
func NewS3Adapter(client *s3.Client) port.StoragePort {
	return &S3Adapter{
		client:   client,
		uploader: manager.NewUploader(client),
	}
}

// Download retrieves the object body
func (a *S3Adapter) Download(ctx context.Context, bucket, key string) (io.ReadCloser, error) {
	resp, err := a.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to download from s3: %w", err)
	}
	return resp.Body, nil
}

// Upload uploads the object using manager.Uploader for efficient multipart/stream uploads
func (a *S3Adapter) Upload(ctx context.Context, bucket, key string, body io.Reader, contentType string) error {
	_, err := a.uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(key),
		Body:        body,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return fmt.Errorf("failed to upload to s3: %w", err)
	}
	return nil
}

// GetMetadata retrieves object metadata
func (a *S3Adapter) GetMetadata(ctx context.Context, bucket, key string) (int64, string, error) {
	resp, err := a.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return 0, "", fmt.Errorf("failed to head object: %w", err)
	}

	var size int64
	if resp.ContentLength != nil {
		size = *resp.ContentLength
	}

	var etag string
	if resp.ETag != nil {
		etag = *resp.ETag
	}

	return size, etag, nil
}
