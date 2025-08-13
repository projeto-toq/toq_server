package gcsadapter

import (
	"context"
	"fmt"
	"log/slog"

	"cloud.google.com/go/storage"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (g *GCSAdapter) CreateUserBucket(ctx context.Context, UserID int64) (err error) {
	if g.adminClient == nil {
		err = status.Error(codes.FailedPrecondition, "gcs admin client not initialized")
		return
	}
	bucketName := fmt.Sprintf("user-%v-bucket", UserID)
	bucketHandle := g.adminClient.Bucket(bucketName)

	bucketAttrs := &storage.BucketAttrs{
		LocationType:             "region",
		Location:                 "southamerica-east1",
		StorageClass:             "STANDARD",
		UniformBucketLevelAccess: storage.UniformBucketLevelAccess{Enabled: true},
		PublicAccessPrevention:   storage.PublicAccessPreventionEnforced,
	}
	if err = bucketHandle.Create(ctx, g.projectID, bucketAttrs); err != nil {
		slog.Error("failed to create bucket", "error", err)
		err = status.Error(codes.Internal, "failed to create bucket")
		return
	}

	policy, err := bucketHandle.IAM().Policy(ctx)
	if err != nil {
		slog.Error("failed to get IAM policy", "error", err)
		err = status.Error(codes.Internal, "failed to get IAM policy")
		return
	}

	writer := fmt.Sprintf("serviceAccount:%s", g.writerSAEmail)
	policy.Add(writer, "roles/storage.legacyBucketWriter")
	reader := fmt.Sprintf("serviceAccount:%s", g.readerSAEmail)
	policy.Add(reader, "roles/storage.legacyBucketReader")
	if err = bucketHandle.IAM().SetPolicy(ctx, policy); err != nil {
		slog.Error("failed to set IAM policy", "error", err)
		err = status.Error(codes.Internal, "failed to set IAM policy")
		return
	}

	return
}
