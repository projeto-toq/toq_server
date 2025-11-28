package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
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

func HandleRequest(ctx context.Context, event mediaprocessingmodel.StepFunctionPayload) (mediaprocessingmodel.StepFunctionPayload, error) {
	// LOG: Start with context
	logger.Info("Validate Lambda started",
		"batch_id", event.BatchID,
		"listing_id", event.ListingID,
		"input_assets_count", len(event.Assets),
	)

	validAssets := make([]mediaprocessingmodel.JobAsset, 0, len(event.Assets))

	for _, asset := range event.Assets {
		// LOG: Checking each asset
		logger.Debug("Validating asset", "key", asset.Key, "batch_id", event.BatchID)

		headOutput, err := s3Client.HeadObject(ctx, &s3.HeadObjectInput{
			Bucket: &bucket,
			Key:    &asset.Key,
		})

		if err != nil {
			// LOG: Specific failure
			logger.Error("Asset validation failed",
				"key", asset.Key,
				"batch_id", event.BatchID,
				"error", err,
			)
			continue
		}

		if headOutput.ContentLength != nil {
			asset.Size = *headOutput.ContentLength
		}
		if headOutput.ETag != nil {
			asset.ETag = *headOutput.ETag
		}
		asset.SourceKey = asset.Key

		validAssets = append(validAssets, asset)

		// LOG: Individual success
		logger.Debug("Asset valid", "key", asset.Key, "size", asset.Size)
	}

	event.ValidAssets = validAssets

	// LOG: Final summary
	logger.Info("Validation complete",
		"batch_id", event.BatchID,
		"valid_count", len(validAssets),
		"invalid_count", len(event.Assets)-len(validAssets),
	)

	return event, nil
}

func main() {
	lambda.Start(HandleRequest)
}
