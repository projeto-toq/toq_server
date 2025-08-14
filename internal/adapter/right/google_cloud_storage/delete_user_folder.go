package gcsadapter

import (
	"context"
	"fmt"
	"log/slog"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (g *GCSAdapter) DeleteUserFolder(ctx context.Context, UserID int64) (err error) {
	if g.writerClient == nil {
		err = status.Error(codes.FailedPrecondition, "gcs writer client not initialized")
		return
	}

	bucketHandle := g.writerClient.Bucket(UsersBucketName)
	prefix := fmt.Sprintf("%d/", UserID)

	it := bucketHandle.Objects(ctx, &storage.Query{Prefix: prefix})
	for {
		objAttrs, itErr := it.Next()
		if itErr == iterator.Done {
			break
		}
		if itErr != nil {
			slog.Error("failed to iterate over user objects", "userID", UserID, "error", itErr)
			err = status.Error(codes.Internal, "failed to list user objects")
			return
		}
		if delErr := bucketHandle.Object(objAttrs.Name).Delete(ctx); delErr != nil {
			slog.Error("failed to delete user object", "userID", UserID, "object", objAttrs.Name, "error", delErr)
			err = status.Error(codes.Internal, "failed to delete user object")
			return
		}
	}

	slog.Info("user folder deleted successfully", "userID", UserID, "bucket", UsersBucketName)
	return
}
