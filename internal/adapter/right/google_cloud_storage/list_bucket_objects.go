package gcsadapter

import (
	"context"
	"log/slog"

	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (g *GCSAdapter) ListBucketObjects(ctx context.Context, bucketName string) (objects []string, err error) {
	if g.readerClient == nil {
		err = status.Error(codes.FailedPrecondition, "gcs reader client not initialized")
		return
	}
	bucketHandle := g.adminClient.Bucket(bucketName)

	it := bucketHandle.Objects(ctx, nil)
	for {
		objAttrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			slog.Error("failed to iterate over bucket objects", "error", err)
			return nil, status.Error(codes.Internal, "failed to iterate over bucket objects")
		}
		objects = append(objects, objAttrs.Name)
	}

	return
}
