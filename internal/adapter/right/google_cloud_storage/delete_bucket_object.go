package gcsadapter

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (g *GCSAdapter) DeleteBucketObject(ctx context.Context, bucketName, objectName string) (err error) {
	if g.writerClient == nil {
		err = status.Error(codes.FailedPrecondition, "gcs writer client not initialized")
		return
	}
	bucketHandle := g.writerClient.Bucket(bucketName)

	err = bucketHandle.Object(objectName).Delete(ctx)

	return
}
