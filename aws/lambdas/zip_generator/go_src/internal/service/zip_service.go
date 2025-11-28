package service

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"path"
	"time"

	"zip_generator/internal/adapter"
	"zip_generator/internal/model"
)

type ZipService struct {
	s3     *adapter.S3Adapter
	logger *slog.Logger
}

func NewZipService(s3 *adapter.S3Adapter, logger *slog.Logger) *ZipService {
	return &ZipService{s3: s3, logger: logger}
}

func (s *ZipService) GenerateZip(ctx context.Context, payload model.StepFunctionPayload) (model.ZipOutput, error) {
	s.logger.Info("Starting ZIP generation", "batch_id", payload.BatchID)

	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)

	// Helper to add files
	addFile := func(key, filename string) error {
		if key == "" {
			return nil
		}
		content, err := s.s3.Download(ctx, key)
		if err != nil {
			s.logger.Error("Failed to download asset", "key", key, "error", err)
			return err
		}

		f, err := zipWriter.Create(filename)
		if err != nil {
			return err
		}
		_, err = f.Write(content)
		return err
	}

	count := 0

	// 1. Originals
	for _, asset := range payload.ValidAssets {
		filename := "originals/" + path.Base(asset.SourceKey)
		if err := addFile(asset.SourceKey, filename); err != nil {
			return model.ZipOutput{}, err
		}
		count++
	}

	// 2. Thumbnails (extracted from ParallelResults)
	for _, result := range payload.ParallelResults {
		for _, thumb := range result.Body.Thumbnails {
			key := thumb.ThumbnailKey
			if key == "" {
				key = thumb.SourceKey
			}
			filename := "thumbnails/" + path.Base(key)
			if err := addFile(key, filename); err != nil {
				return model.ZipOutput{}, err
			}
			count++
		}
	}

	if err := zipWriter.Close(); err != nil {
		return model.ZipOutput{}, err
	}

	// Upload ZIP
	key := fmt.Sprintf("%d/zip/complete_%d_%d.zip", payload.ListingID, payload.BatchID, time.Now().Unix())
	if err := s.s3.Upload(ctx, key, buf.Bytes(), "application/zip"); err != nil {
		return model.ZipOutput{}, err
	}

	return model.ZipOutput{
		ZipKey:       key,
		AssetsZipped: count,
	}, nil
}
