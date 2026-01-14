package s3adapter

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"

	derrors "github.com/projeto-toq/toq_server/internal/core/derrors"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

const maxDeleteObjectsBatchSize = 1000

// DeleteKeys removes multiple objects from the listing bucket using S3 DeleteObjects.
// Missing keys (NoSuchKey) are treated as successful deletions to keep the operation idempotent.
func (a *ListingMediaStorageAdapter) DeleteKeys(ctx context.Context, keys []string) error {
	ctx = utils.ContextWithLogger(ctx)
	ctx, end, _ := utils.GenerateBusinessTracer(ctx, "ListingMediaStorage.DeleteKeys")
	defer end()

	if err := a.ensureClients(); err != nil {
		utils.SetSpanError(ctx, err)
		return err
	}

	sanitized := dedupeKeys(keys)
	if len(sanitized) == 0 {
		return nil
	}

	for start := 0; start < len(sanitized); start += maxDeleteObjectsBatchSize {
		endIdx := start + maxDeleteObjectsBatchSize
		if endIdx > len(sanitized) {
			endIdx = len(sanitized)
		}

		batch := sanitized[start:endIdx]
		if err := a.deleteBatch(ctx, batch); err != nil {
			utils.SetSpanError(ctx, err)
			return err
		}
	}

	return nil
}

func (a *ListingMediaStorageAdapter) deleteBatch(ctx context.Context, keys []string) error {
	identifiers := make([]types.ObjectIdentifier, 0, len(keys))
	for _, key := range keys {
		identifiers = append(identifiers, types.ObjectIdentifier{Key: aws.String(key)})
	}

	output, err := a.adminClient.DeleteObjects(ctx, &s3.DeleteObjectsInput{
		Bucket: aws.String(a.bucket),
		Delete: &types.Delete{Objects: identifiers, Quiet: aws.Bool(true)},
	})
	if err != nil {
		logger := utils.LoggerFromContext(ctx)
		logger.Error("adapter.s3.listing.delete_objects_error", "error", err, "keys_count", len(keys))
		return derrors.Infra("failed to delete objects", err)
	}

	if len(output.Errors) == 0 {
		return nil
	}

	var failures []string
	for _, item := range output.Errors {
		code := strings.TrimSpace(aws.ToString(item.Code))
		if strings.EqualFold(code, "NoSuchKey") {
			// Idempotent delete: ignore missing keys
			continue
		}
		failures = append(failures, fmt.Sprintf("key=%s code=%s message=%s", aws.ToString(item.Key), code, aws.ToString(item.Message)))
	}

	if len(failures) == 0 {
		return nil
	}

	errMsg := fmt.Sprintf("s3 delete_objects returned errors: %s", strings.Join(failures, "; "))
	logger := utils.LoggerFromContext(ctx)
	logger.Error("adapter.s3.listing.delete_objects_partial_failure", "details", errMsg)
	return derrors.Infra(errMsg, fmt.Errorf("%s", errMsg))
}

func dedupeKeys(keys []string) []string {
	seen := make(map[string]struct{}, len(keys))
	deduped := make([]string, 0, len(keys))
	for _, key := range keys {
		trimmed := strings.TrimSpace(key)
		if trimmed == "" {
			continue
		}
		if _, ok := seen[trimmed]; ok {
			continue
		}
		seen[trimmed] = struct{}{}
		deduped = append(deduped, trimmed)
	}
	return deduped
}
