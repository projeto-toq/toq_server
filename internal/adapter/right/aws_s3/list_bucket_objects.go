package s3adapter

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (s *S3Adapter) ListBucketObjects(ctx context.Context, bucketName string) (objects []string, err error) {
	if s.readerClient == nil {
		err = errors.New("s3 reader client is nil")
		return
	}

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)
	logger.Debug("adapter.s3.list_objects.start", "bucket", bucketName)

	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
	}

	paginator := s3.NewListObjectsV2Paginator(s.readerClient, input)

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			utils.SetSpanError(ctx, err)
			logger.Error("adapter.s3.list_objects.iteration_error", "bucket", bucketName, "error", err)
			return nil, err
		}

		for _, obj := range output.Contents {
			if obj.Key != nil {
				objects = append(objects, *obj.Key)
			}
		}
	}

	logger.Debug("adapter.s3.list_objects.success", "bucket", bucketName, "count", len(objects))
	return objects, nil
}
