package imageprocessing

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"io"
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
	data, err := s.downloadBytes(ctx, bucket, key)
	if err != nil {
		return nil, err
	}

	img, err := s.decodeWithOrientation(data)
	if err != nil {
		return nil, err
	}

	generatedKeys := make([]string, 0, len(targetSizes))
	for _, size := range targetSizes {
		newKey, err := s.persistVariant(ctx, bucket, key, size, img)
		if err != nil {
			return nil, err
		}
		generatedKeys = append(generatedKeys, newKey)
	}

	return generatedKeys, nil
}

var targetSizes = []struct {
	Name  string
	Width int
}{
	{Name: "thumbnail", Width: 200},
	{Name: "small", Width: 400},
	{Name: "medium", Width: 800},
	{Name: "large", Width: 1200},
}

func (s *ThumbnailService) downloadBytes(ctx context.Context, bucket, key string) ([]byte, error) {
	reader, err := s.storage.Download(ctx, bucket, key)
	if err != nil {
		return nil, fmt.Errorf("failed to download image: %w", err)
	}
	defer reader.Close()

	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read image: %w", err)
	}
	return data, nil
}

func (s *ThumbnailService) decodeWithOrientation(data []byte) (image.Image, error) {
	img, err := imaging.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}
	orientation := readOrientation(data)
	return normalizeOrientation(img, orientation), nil
}

func (s *ThumbnailService) persistVariant(ctx context.Context, bucket, originalKey string, size struct {
	Name  string
	Width int
}, img image.Image) (string, error) {
	resizedImg := imaging.Resize(img, size.Width, 0, imaging.Lanczos)
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, resizedImg, &jpeg.Options{Quality: 85}); err != nil {
		return "", fmt.Errorf("failed to encode resized image %s: %w", size.Name, err)
	}

	newKey, err := s.generateKey(originalKey, size.Name)
	if err != nil {
		return "", fmt.Errorf("failed to generate key for %s: %w", size.Name, err)
	}

	if err := s.storage.Upload(ctx, bucket, newKey, bytes.NewReader(buf.Bytes()), "image/jpeg"); err != nil {
		return "", fmt.Errorf("failed to upload %s: %w", size.Name, err)
	}

	return newKey, nil
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
