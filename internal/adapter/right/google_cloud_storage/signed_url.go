package gcsadapter

import (
	"fmt"
	"time"

	"cloud.google.com/go/storage"
)

// GenerateV4PutObjectSignedURL gera uma URL para upload (PUT).
func (g *GCSAdapter) GenerateV4PutObjectSignedURL(bucketName, objectName, contentType string) (string, error) {
	if g.writerClient == nil {
		return "", fmt.Errorf("writer client is not initialized")
	}

	opts := &storage.SignedURLOptions{
		GoogleAccessID: g.writerSAEmail,
		PrivateKey:     g.writerPrivateKey,
		Scheme:         storage.SigningSchemeV4,
		Method:         "PUT",
		Expires:        time.Now().Add(15 * time.Minute),
		ContentType:    contentType,
	}

	url, err := storage.SignedURL(bucketName, objectName, opts)
	if err != nil {
		return "", fmt.Errorf("failed to generate signed URL for PUT: %w", err)
	}

	return url, nil
}

// GenerateV4GetObjectSignedURL gera uma URL para download (GET).
func (g *GCSAdapter) GenerateV4GetObjectSignedURL(bucketName, objectName string) (string, error) {
	if g.readerClient == nil {
		return "", fmt.Errorf("reader client is not initialized")
	}

	opts := &storage.SignedURLOptions{
		GoogleAccessID: g.readerSAEmail,
		PrivateKey:     g.readerPrivateKey,
		Scheme:         storage.SigningSchemeV4,
		Method:         "GET",
		Expires:        time.Now().Add(60 * time.Minute),
	}

	url, err := storage.SignedURL(bucketName, objectName, opts)
	if err != nil {
		return "", fmt.Errorf("failed to generate signed URL for GET: %w", err)
	}

	return url, nil
}

// GeneratePhotoSignedURL gera uma URL para upload de foto específica do usuário.
func (g *GCSAdapter) GeneratePhotoSignedURL(bucketName string, userID int64, photoType, contentType string) (string, error) {
	objectPath := fmt.Sprintf("%d/%s", userID, photoType)
	return g.GenerateV4PutObjectSignedURL(bucketName, objectPath, contentType)
}

// GeneratePhotoDownloadURL gera uma URL para download de foto específica do usuário.
func (g *GCSAdapter) GeneratePhotoDownloadURL(bucketName string, userID int64, photoType string) (string, error) {
	objectPath := fmt.Sprintf("%d/%s", userID, photoType)
	return g.GenerateV4GetObjectSignedURL(bucketName, objectPath)
}
