package s3adapter

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *S3Adapter) DeleteBucketObject(ctx context.Context, bucketName, objectName string) (err error) {
	if s.adminClient == nil {
		err = status.Error(codes.FailedPrecondition, "s3 admin client not initialized")
		return
	}

	input := &s3.DeleteObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectName),
	}

	_, err = s.adminClient.DeleteObject(ctx, input)
	if err != nil {
		return status.Error(codes.Internal, "failed to delete object from S3")
	}

	return nil
}
