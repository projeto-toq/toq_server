package s3adapter

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// GenerateV4PutObjectSignedURL gera uma URL assinada para upload (PUT) no S3
func (s *S3Adapter) GenerateV4PutObjectSignedURL(bucketName, objectName, contentType string) (string, error) {
	if s.adminClient == nil {
		return "", fmt.Errorf("admin client is not initialized")
	}

	presignClient := s3.NewPresignClient(s.adminClient)

	request, err := presignClient.PresignPutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(objectName),
		ContentType: aws.String(contentType),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = 15 * time.Minute
	})

	if err != nil {
		return "", fmt.Errorf("failed to generate S3 signed URL for PUT: %w", err)
	}

	return request.URL, nil
}

// GenerateV4GetObjectSignedURL gera uma URL assinada para download (GET) no S3
func (s *S3Adapter) GenerateV4GetObjectSignedURL(bucketName, objectName string) (string, error) {
	if s.readerClient == nil {
		return "", fmt.Errorf("reader client is not initialized")
	}

	presignClient := s3.NewPresignClient(s.readerClient)

	request, err := presignClient.PresignGetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectName),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = 60 * time.Minute
	})

	if err != nil {
		return "", fmt.Errorf("failed to generate S3 signed URL for GET: %w", err)
	}

	return request.URL, nil
}

// GeneratePhotoSignedURL gera uma URL para upload de foto específica do usuário
func (s *S3Adapter) GeneratePhotoSignedURL(bucketName string, userID int64, photoType, contentType string) (string, error) {
	objectPath := fmt.Sprintf("%d/%s", userID, photoType)
	return s.GenerateV4PutObjectSignedURL(bucketName, objectPath, contentType)
}

// GeneratePhotoDownloadURL gera uma URL para download de foto específica do usuário
func (s *S3Adapter) GeneratePhotoDownloadURL(bucketName string, userID int64, photoType string) (string, error) {
	objectPath := fmt.Sprintf("%d/%s", userID, photoType)
	return s.GenerateV4GetObjectSignedURL(bucketName, objectPath)
}
