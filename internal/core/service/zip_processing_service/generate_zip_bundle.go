package zipprocessingservice

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"path"
	"time"

	"github.com/projeto-toq/toq_server/internal/core/derrors"
	mediaprocessingmodel "github.com/projeto-toq/toq_server/internal/core/model/media_processing_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// GenerateZipBundle orchestrates the download, compression, and upload of assets
func (s *zipProcessingService) GenerateZipBundle(ctx context.Context, input mediaprocessingmodel.GenerateZipInput) (mediaprocessingmodel.ZipOutput, error) {
	ctx, spanEnd, _ := utils.GenerateTracer(ctx)
	defer spanEnd()

	logger := utils.LoggerFromContext(ctx)
	logger.Info("service.zip.generate_bundle.start", "batch_id", input.BatchID, "listing_id", input.ListingID)

	// 1. Create ZIP buffer
	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)

	// Helper function to add file to zip
	addFileToZip := func(key string, filename string) error {
		content, err := s.s3Adapter.DownloadFile(ctx, key)
		if err != nil {
			logger.Error("service.zip.download_error", "key", key, "err", err)
			return err
		}

		f, err := zipWriter.Create(filename)
		if err != nil {
			logger.Error("service.zip.create_entry_error", "filename", filename, "err", err)
			return err
		}

		_, err = f.Write(content)
		if err != nil {
			logger.Error("service.zip.write_entry_error", "filename", filename, "err", err)
			return err
		}
		return nil
	}

	// 2. Process Original Assets
	for _, asset := range input.ValidAssets {
		// Use SourceKey for download, and a clean filename for the zip
		// Assuming SourceKey format: {listingID}/raw/{type}/{date}/{filename}
		// We want just {filename} in the zip, or maybe organized folders?
		// Let's keep it simple for now: just the filename from the key
		filename := path.Base(asset.SourceKey)
		if err := addFileToZip(asset.SourceKey, "originals/"+filename); err != nil {
			return mediaprocessingmodel.ZipOutput{}, derrors.Infra("failed to process asset", err)
		}
	}

	// 3. Process Thumbnails
	for _, thumb := range input.Thumbnails {
		// Assuming ThumbnailKey is available in MediaAsset (it might be mapped to ProcessedKey or similar)
		// The model MediaAsset has ProcessedKey, ThumbnailKey etc.
		// In the Step Function output, thumbnails are MediaAssets.
		key := thumb.ThumbnailKey // Or SourceKey? The lambda code used 'thumbnailKey'
		if key == "" {
			key = thumb.SourceKey // Fallback
		}

		filename := path.Base(key)
		if err := addFileToZip(key, "thumbnails/"+filename); err != nil {
			return mediaprocessingmodel.ZipOutput{}, derrors.Infra("failed to process thumbnail", err)
		}
	}

	if err := zipWriter.Close(); err != nil {
		return mediaprocessingmodel.ZipOutput{}, derrors.Infra("failed to close zip", err)
	}

	// 4. Upload ZIP to correct path (Fixing the path issue)
	// Path: {listingID}/zip/complete_{batchID}.zip
	key := fmt.Sprintf("%d/zip/complete_%d_%d.zip", input.ListingID, input.BatchID, time.Now().Unix())

	if err := s.s3Adapter.UploadFile(ctx, key, buf.Bytes(), "application/zip"); err != nil {
		return mediaprocessingmodel.ZipOutput{}, derrors.Infra("failed to upload zip", err)
	}

	logger.Info("service.zip.generate_bundle.success", "key", key)

	return mediaprocessingmodel.ZipOutput{
		ZipKey:       key,
		AssetsZipped: len(input.ValidAssets) + len(input.Thumbnails),
	}, nil
}
