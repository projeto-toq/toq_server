package s3adapter

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	awshttp "github.com/aws/aws-sdk-go-v2/aws/transport/http"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go"
)

// ObjectExists checks if an object exists in the configured S3 bucket using HEAD request.
func (s *S3Adapter) ObjectExists(ctx context.Context, bucketName, objectName string) (bool, error) {
	// Usa o readerClient por princípio de menor privilégio
	_, err := s.readerClient.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectName),
	})
	if err == nil {
		return true, nil
	}
	// Detecta 404 via ResponseError
	var respErr *awshttp.ResponseError
	if errors.As(err, &respErr) && respErr.Response != nil && respErr.Response.StatusCode == 404 {
		return false, nil
	}
	// Detecta códigos semânticos via smithy APIError
	var apiErr smithy.APIError
	if errors.As(err, &apiErr) {
		code := apiErr.ErrorCode()
		if code == "NotFound" || code == "NoSuchKey" { // segurança
			return false, nil
		}
	}
	// Para outros erros (permissão, rede, etc) propagar
	return false, err
}
