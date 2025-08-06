package gcsadapter

import (
	"context"
	"fmt"
	"log/slog"

	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (g *GCSAdapter) DeleteUserBucket(ctx context.Context, UserID int64) (err error) {
	if g.adminClient == nil {
		err = status.Error(codes.FailedPrecondition, "gcs admin client not initialized")
		return
	}
	bucketName := fmt.Sprintf("user-%v-bucket", UserID)
	bucketHandle := g.adminClient.Bucket(bucketName)

	it := bucketHandle.Objects(ctx, nil)
	for {
		objAttrs, itErr := it.Next()
		if itErr == iterator.Done {
			break
		}
		if itErr != nil {
			err = status.Error(codes.Internal, "failed to list objects")
			return
		}
		if delErr := bucketHandle.Object(objAttrs.Name).Delete(ctx); delErr != nil {
			err = status.Error(codes.Internal, "failed to delete object")
			return
		}
	}

	err = bucketHandle.Delete(ctx)
	if err != nil {
		slog.Error("failed to delete bucket", "error", err)
		err = status.Error(codes.Internal, "failed to delete bucket")
		return
	}

	return
}
