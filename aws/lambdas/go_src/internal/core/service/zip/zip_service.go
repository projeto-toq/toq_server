package zip

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/projeto-toq/toq_server/aws/lambdas/go_src/internal/core/port"
)

type ZipService struct {
	storage port.StoragePort
}

func NewZipService(storage port.StoragePort) *ZipService {
	return &ZipService{
		storage: storage,
	}
}

// CreateZip creates a zip file from the given keys and uploads it to the destination key
func (s *ZipService) CreateZip(ctx context.Context, bucket string, sourceKeys []string, destinationKey string) error {
	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)

	for _, key := range sourceKeys {
		// 1. Download file
		reader, err := s.storage.Download(ctx, bucket, key)
		if err != nil {
			// Log error but maybe continue? For now, fail.
			return fmt.Errorf("failed to download %s: %w", key, err)
		}

		// 2. Determine internal path in zip
		internalPath := s.cleanPath(key)

		// 3. Create entry in zip
		header := &zip.FileHeader{
			Name:   internalPath,
			Method: zip.Deflate,
		}
		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			reader.Close()
			return fmt.Errorf("failed to create zip header for %s: %w", key, err)
		}

		// 4. Copy content
		_, err = io.Copy(writer, reader)
		reader.Close()
		if err != nil {
			return fmt.Errorf("failed to write %s to zip: %w", key, err)
		}
	}

	if err := zipWriter.Close(); err != nil {
		return fmt.Errorf("failed to close zip writer: %w", err)
	}

	// 5. Upload zip
	if err := s.storage.Upload(ctx, bucket, destinationKey, buf, "application/zip"); err != nil {
		return fmt.Errorf("failed to upload zip to %s: %w", destinationKey, err)
	}

	return nil
}

// cleanPath removes 'processed/' prefix and date segments to create a clean internal zip path
func (s *ZipService) cleanPath(key string) string {
	// Find where 'processed/' starts
	idx := strings.Index(key, "processed/")
	if idx != -1 {
		// Take everything after 'processed/'
		key = key[idx+len("processed/"):]
	}

	// Remove date segment (YYYY-MM-DD/)
	dateRegex := regexp.MustCompile(`\d{4}-\d{2}-\d{2}/`)
	key = dateRegex.ReplaceAllString(key, "")

	return key
}
