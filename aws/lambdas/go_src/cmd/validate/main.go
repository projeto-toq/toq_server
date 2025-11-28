package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sfn"
	mediaprocessingmodel "github.com/projeto-toq/toq_server/internal/core/model/media_processing_model"
)

var (
	s3Client  *s3.Client
	sfnClient *sfn.Client
	logger    *slog.Logger
	bucket    string
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
	sfnClient = sfn.NewFromConfig(cfg)
}

func HandleRequest(ctx context.Context, event mediaprocessingmodel.StepFunctionPayload) (mediaprocessingmodel.StepFunctionPayload, error) {
	logger.Info("Validate Lambda started", "batchId", event.BatchID, "listingId", event.ListingID)

	validAssets := make([]mediaprocessingmodel.JobAsset, 0, len(event.Assets))

	for _, asset := range event.Assets {
		headOutput, err := s3Client.HeadObject(ctx, &s3.HeadObjectInput{
			Bucket: &bucket,
			Key:    &asset.Key,
		})

		if err != nil {
			logger.Error("Failed to validate asset", "key", asset.Key, "error", err)
			// We skip invalid assets but log them. In a stricter mode, we might fail the batch.
			continue
		}

		if headOutput.ContentLength != nil {
			asset.Size = *headOutput.ContentLength
		}
		if headOutput.ETag != nil {
			asset.ETag = *headOutput.ETag
		}
		// SourceKey is same as Key for now, but could be different if we moved files.
		asset.SourceKey = asset.Key

		validAssets = append(validAssets, asset)
	}

	event.ValidAssets = validAssets
	logger.Info("Validation complete", "validCount", len(validAssets), "totalCount", len(event.Assets))

	return event, nil
}

func main() {
	lambda.Start(HandleRequest)
}
