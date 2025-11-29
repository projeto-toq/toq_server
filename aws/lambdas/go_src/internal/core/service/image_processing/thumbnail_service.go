package imageprocessing

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"log/slog"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/projeto-toq/toq_server/aws/lambdas/go_src/internal/core/port"
	mediaprocessingmodel "github.com/projeto-toq/toq_server/internal/core/model/media_processing_model"
)

// ThumbnailService handles image resizing and processing logic
type ThumbnailService struct {
	storage port.StoragePort
	bucket  string
	logger  *slog.Logger
}

// NewThumbnailService creates a new instance of ThumbnailService
func NewThumbnailService(storage port.StoragePort, bucket string, logger *slog.Logger) *ThumbnailService {
	return &ThumbnailService{
		storage: storage,
		bucket:  bucket,
		logger:  logger,
	}
}

type targetSize struct {
	Name    string
	Width   int
	Quality int
}

// ProcessAsset downloads, resizes, and uploads thumbnails for a given asset
// It generates Large (1920px), Medium (1280px), Small (640px), and Tiny (300px) versions.
func (s *ThumbnailService) ProcessAsset(ctx context.Context, listingID uint64, asset mediaprocessingmodel.JobAsset) ([]mediaprocessingmodel.JobAsset, error) {
	if !strings.HasPrefix(asset.Type, "PHOTO") {
		s.logger.Info("Skipping non-photo asset", "key", asset.Key, "type", asset.Type)
		return nil, nil
	}

	s.logger.Info("Processing photo", "key", asset.Key)

	// 1. Download
	body, err := s.storage.Download(ctx, s.bucket, asset.Key)
	if err != nil {
		return nil, fmt.Errorf("failed to download asset %s: %w", asset.Key, err)
	}
	defer body.Close()

	// 2. Decode
	img, _, err := image.Decode(body)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image %s: %w", asset.Key, err)
	}

	// 3. Define sizes
	sizes := []targetSize{
		{Name: "large", Width: 1920, Quality: 80},
		{Name: "medium", Width: 1280, Quality: 70},
		{Name: "small", Width: 640, Quality: 60},
		{Name: "tiny", Width: 300, Quality: 50},
	}

	var generatedAssets []mediaprocessingmodel.JobAsset

	// 4. Process each size
	for _, size := range sizes {
		// Resize maintaining aspect ratio
		resized := imaging.Resize(img, size.Width, 0, imaging.Lanczos)

		// Encode to JPEG
		buf := new(bytes.Buffer)
		if err := jpeg.Encode(buf, resized, &jpeg.Options{Quality: size.Quality}); err != nil {
			s.logger.Error("Failed to encode thumbnail", "key", asset.Key, "size", size.Name, "error", err)
			continue
		}

		// Generate Key
		newKey := s.generateKey(asset.Key, size.Name)

		// Upload
		if err := s.storage.Upload(ctx, s.bucket, newKey, buf, "image/jpeg"); err != nil {
			s.logger.Error("Failed to upload thumbnail", "key", newKey, "error", err)
			continue
		}

		generatedAssets = append(generatedAssets, mediaprocessingmodel.JobAsset{
			Key:  newKey,
			Type: fmt.Sprintf("PHOTO_%s", strings.ToUpper(size.Name)),
		})
	}

	return generatedAssets, nil
}

func (s *ThumbnailService) generateKey(originalKey, sizeName string) string {
	// Default fallback if structure is unexpected
	dir := filepath.Dir(originalKey)
	base := filepath.Base(originalKey)
	ext := filepath.Ext(base)
	name := strings.TrimSuffix(base, ext)

	// Try to replace /raw/ with /processed/{size}/
	if strings.Contains(originalKey, "/raw/") {
		// 1. Replace /raw/ with /processed/{size}/
		tempKey := strings.Replace(originalKey, "/raw/", fmt.Sprintf("/processed/%s/", sizeName), 1)

		// 2. Handle filename
		tempDir := filepath.Dir(tempKey)
		// Re-assemble
		return fmt.Sprintf("%s/%s_%s.jpg", tempDir, name, sizeName)
	}

	// Fallback: append size to directory and filename
	return fmt.Sprintf("%s/%s/%s_%s.jpg", dir, sizeName, name, sizeName)
}
