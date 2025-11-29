package imageprocessing

import (
	"bytes"
	"context"
	"fmt"
	"image/jpeg"
	"path"
	"regexp"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/projeto-toq/toq_server/aws/lambdas/go_src/internal/core/port"
)

type ThumbnailService struct {
	storage port.StoragePort
}

func NewThumbnailService(storage port.StoragePort) *ThumbnailService {
	return &ThumbnailService{
		storage: storage,
	}
}

func (s *ThumbnailService) ProcessImage(ctx context.Context, bucket, key string) ([]string, error) {
	// 1. Download
	reader, err := s.storage.Download(ctx, bucket, key)
	if err != nil {
		return nil, fmt.Errorf("failed to download image: %w", err)
	}
	defer reader.Close()

	// 2. Decode with auto-orientation (fixes rotation issues)
	// imaging.Decode handles EXIF orientation automatically
	img, err := imaging.Decode(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	// 3. Define sizes
	// "tiny" renamed to "thumbnail" as requested
	sizes := []struct {
		Name  string
		Width int
	}{
		{Name: "thumbnail", Width: 200},
		{Name: "small", Width: 400},
		{Name: "medium", Width: 800},
		{Name: "large", Width: 1200},
	}

	var generatedKeys []string

	// 4. Process each size
	for _, size := range sizes {
		resizedImg := imaging.Resize(img, size.Width, 0, imaging.Lanczos)

		// Encode to JPEG
		var buf bytes.Buffer
		if err := jpeg.Encode(&buf, resizedImg, &jpeg.Options{Quality: 85}); err != nil {
			return nil, fmt.Errorf("failed to encode resized image %s: %w", size.Name, err)
		}

		// Generate new key
		newKey, err := s.generateKey(key, size.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to generate key for %s: %w", size.Name, err)
		}

		// Upload
		if err := s.storage.Upload(ctx, bucket, newKey, &buf, "image/jpeg"); err != nil {
			return nil, fmt.Errorf("failed to upload %s: %w", size.Name, err)
		}

		generatedKeys = append(generatedKeys, newKey)
	}

	return generatedKeys, nil
}

// generateKey transforms the raw key into the processed key
func (s *ThumbnailService) generateKey(originalKey, sizeName string) (string, error) {
	// Check for 'raw/' segment
	if !strings.Contains(originalKey, "raw/") {
		return "", fmt.Errorf("invalid key format: must contain 'raw/' segment")
	}

	// Split into prefix (e.g. "123/" or empty) and suffix (e.g. "photo/vertical/...")
	parts := strings.SplitN(originalKey, "raw/", 2)
	prefix := parts[0] // "123/" or ""
	suffix := parts[1] // "photo/vertical/..."

	// Remove date segment from suffix if present
	dateRegex := regexp.MustCompile(`\d{4}-\d{2}-\d{2}/`)
	cleanSuffix := dateRegex.ReplaceAllString(suffix, "")

	// Construct new path components
	dir := path.Dir(cleanSuffix)   // "photo/vertical"
	file := path.Base(cleanSuffix) // "uuid.jpg"

	// Final structure: {prefix}processed/{dir}/{size}/{file}
	// Example: 123/processed/photo/vertical/thumbnail/uuid.jpg
	newKey := fmt.Sprintf("%sprocessed/%s/%s/%s", prefix, dir, sizeName, file)

	// Clean up any double slashes just in case
	newKey = strings.ReplaceAll(newKey, "//", "/")

	return newKey, nil
}
