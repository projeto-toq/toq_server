package s3adapter

import (
	"context"
	"errors"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func (s *S3Adapter) ListBucketObjects(ctx context.Context, bucketName string) (objects []string, err error) {
	if s.readerClient == nil {
		err = errors.New("s3 reader client is nil")
		return
	}

	slog.Debug("Listing bucket objects in S3", "bucket", bucketName)

	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
	}

	paginator := s3.NewListObjectsV2Paginator(s.readerClient, input)

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			slog.Error("failed to iterate over bucket objects in S3", "error", err)
			return nil, err
		}

		for _, obj := range output.Contents {
			if obj.Key != nil {
				objects = append(objects, *obj.Key)
			}
		}
	}

	slog.Debug("Successfully listed bucket objects", "bucket", bucketName, "count", len(objects))
	return objects, nil
}
