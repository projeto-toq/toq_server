package s3adapter

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func (s *S3Adapter) DeleteBucketObject(ctx context.Context, bucketName, objectName string) (err error) {
	if s.adminClient == nil {
		err = errors.New("s3 admin client is nil")
		return
	}

	input := &s3.DeleteObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectName),
	}

	_, err = s.adminClient.DeleteObject(ctx, input)
	if err != nil {
		return err
	}

	return nil
}
