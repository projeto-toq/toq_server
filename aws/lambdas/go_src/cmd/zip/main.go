package main

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	mediaprocessingmodel "github.com/projeto-toq/toq_server/internal/core/model/media_processing_model"
)

var (
	s3Client *s3.Client
	logger   *slog.Logger
	bucket   string
)

func init() {
	logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	bucket = os.Getenv("MEDIA_BUCKET")
	if bucket == "" {
		bucket = "toq-listing-medias"
	}

	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		logger.Error("Failed to load AWS config", "error", err)
		os.Exit(1)
	}

	s3Client = s3.NewFromConfig(cfg)
}

type ZipOutput struct {
	ZipKey       string `json:"zipKey"`
	AssetsZipped int    `json:"assetsZipped"`
}

func HandleRequest(ctx context.Context, event mediaprocessingmodel.StepFunctionPayload) (mediaprocessingmodel.LambdaResponse, error) {
	logger.Info("Zip Lambda started", "batchId", event.BatchID)

	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)

	zippedCount := 0

	for _, asset := range event.ValidAssets {
		logger.Info("Zipping asset", "key", asset.Key)

		resp, err := s3Client.GetObject(ctx, &s3.GetObjectInput{
			Bucket: &bucket,
			Key:    &asset.Key,
		})
		if err != nil {
			logger.Error("Failed to download asset for zip", "key", asset.Key, "error", err)
			continue
		}

		f, err := zipWriter.Create(asset.Key)
		if err != nil {
			logger.Error("Failed to create zip entry", "key", asset.Key, "error", err)
			resp.Body.Close()
			continue
		}

		_, err = io.Copy(f, resp.Body)
		resp.Body.Close()
		if err != nil {
			logger.Error("Failed to write zip entry", "key", asset.Key, "error", err)
			continue
		}

		zippedCount++
	}

	if err := zipWriter.Close(); err != nil {
		return mediaprocessingmodel.LambdaResponse{}, fmt.Errorf("failed to close zip writer: %w", err)
	}

	zipKey := fmt.Sprintf("%d/zip/complete_%d_%d.zip", event.ListingID, event.BatchID, time.Now().Unix())

	_, err := s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      &bucket,
		Key:         &zipKey,
		Body:        bytes.NewReader(buf.Bytes()),
		ContentType: aws.String("application/zip"),
	})
	if err != nil {
		return mediaprocessingmodel.LambdaResponse{}, fmt.Errorf("failed to upload zip: %w", err)
	}

	return mediaprocessingmodel.LambdaResponse{
		Body: ZipOutput{
			ZipKey:       zipKey,
			AssetsZipped: zippedCount,
		},
	}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
