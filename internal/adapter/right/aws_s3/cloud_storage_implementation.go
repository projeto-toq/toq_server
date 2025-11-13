package s3adapter

import (
	"fmt"

	storagemodel "github.com/projeto-toq/toq_server/internal/core/model/storage_model"
)

// GetBucketConfig retorna a configuração do bucket de usuários
// For listing-specific operations, use GetListingBucketConfig() instead
func (s *S3Adapter) GetBucketConfig() storagemodel.BucketConfig {
	return storagemodel.BucketConfig{
		Name:   s.userBucketName,
		Region: s.region,
	}
}

// GeneratePhotoUploadURL gera uma URL para upload de foto específica do usuário usando abstração
func (s *S3Adapter) GeneratePhotoUploadURL(userID int64, photoType storagemodel.PhotoType, contentType string) (string, error) {
	objectPath := fmt.Sprintf("%d/%s", userID, string(photoType))
	return s.GenerateV4PutObjectSignedURL(s.userBucketName, objectPath, contentType)
}

// GeneratePhotoDownloadURL gera uma URL para download de foto específica do usuário usando abstração
func (s *S3Adapter) GeneratePhotoDownloadURL(userID int64, photoType storagemodel.PhotoType) (string, error) {
	objectPath := fmt.Sprintf("%d/%s", userID, string(photoType))
	return s.GenerateV4GetObjectSignedURL(s.userBucketName, objectPath)
}

// GenerateDocumentUploadURL gera uma URL para upload de documento específico do usuário
func (s *S3Adapter) GenerateDocumentUploadURL(userID int64, docType storagemodel.DocumentType, contentType string) (string, error) {
	objectPath := fmt.Sprintf("%d/%s", userID, string(docType))
	return s.GenerateV4PutObjectSignedURL(s.userBucketName, objectPath, contentType)
}

// GenerateDocumentDownloadURL gera uma URL para download de documento específico do usuário
func (s *S3Adapter) GenerateDocumentDownloadURL(userID int64, docType storagemodel.DocumentType) (string, error) {
	objectPath := fmt.Sprintf("%d/%s", userID, string(docType))
	return s.GenerateV4GetObjectSignedURL(s.userBucketName, objectPath)
}

// GenerateUploadURL é um wrapper para o método existing GenerateV4PutObjectSignedURL
func (s *S3Adapter) GenerateUploadURL(bucketName, objectName, contentType string) (string, error) {
	return s.GenerateV4PutObjectSignedURL(bucketName, objectName, contentType)
}

// GenerateDownloadURL é um wrapper para o método existing GenerateV4GetObjectSignedURL
func (s *S3Adapter) GenerateDownloadURL(bucketName, objectName string) (string, error) {
	return s.GenerateV4GetObjectSignedURL(bucketName, objectName)
}
