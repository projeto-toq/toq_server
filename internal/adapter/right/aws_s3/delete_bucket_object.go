package s3adapter

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	
	
"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (s *S3Adapter) DeleteBucketObject(ctx context.Context, bucketName, objectName string) (err error) {
	if s.adminClient == nil {
		err = utils.ErrInternalServer
		return
	}

	input := &s3.DeleteObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectName),
	}

	_, err = s.adminClient.DeleteObject(ctx, input)
	if err != nil {
		return utils.ErrInternalServer
	}

	return nil
}
