package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sfn"
	mediaprocessingmodel "github.com/projeto-toq/toq_server/internal/core/model/media_processing_model"
)

var (
	s3Client        *s3.Client
	sfnClient       *sfn.Client
	logger          *slog.Logger
	bucket          string
	stateMachineArn string
)

func init() {
	logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	bucket = os.Getenv("MEDIA_BUCKET")
	if bucket == "" {
		bucket = "toq-listing-medias"
	}
	stateMachineArn = os.Getenv("STATE_MACHINE_ARN")
	if stateMachineArn == "" {
		stateMachineArn = "arn:aws:states:us-east-1:058264253741:stateMachine:listing-media-processing-sm-staging"
	}

	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		logger.Error("Failed to load AWS config", "error", err)
		os.Exit(1)
	}

	s3Client = s3.NewFromConfig(cfg)
	sfnClient = sfn.NewFromConfig(cfg)
}

func HandleRequest(ctx context.Context, rawEvent json.RawMessage) (mediaprocessingmodel.StepFunctionPayload, error) {
	// LOG: Raw event for debugging
	logger.Info("Validate Lambda received raw event", "raw_event", string(rawEvent))

	// Try to parse as SQS Event
	var sqsEvent events.SQSEvent
	if err := json.Unmarshal(rawEvent, &sqsEvent); err == nil && len(sqsEvent.Records) > 0 && sqsEvent.Records[0].EventSource == "aws:sqs" {
		logger.Info("Detected SQS Event", "record_count", len(sqsEvent.Records))

		for _, record := range sqsEvent.Records {
			logger.Info("Processing SQS record", "message_id", record.MessageId)

			var rawPayload struct {
				JobID     uint64          `json:"jobId"`
				BatchID   uint64          `json:"batchId"`
				ListingID uint64          `json:"listingId"`
				Assets    json.RawMessage `json:"assets"`
				Retry     uint16          `json:"retry"`
			}

			if err := json.Unmarshal([]byte(record.Body), &rawPayload); err != nil {
				logger.Error("Failed to unmarshal SQS body structure", "error", err, "body", record.Body)
				continue
			}

			var assets []mediaprocessingmodel.JobAsset
			// Try []JobAsset
			if err := json.Unmarshal(rawPayload.Assets, &assets); err != nil {
				// Try []string (Legacy/Backend mismatch fix)
				var assetKeys []string
				if err2 := json.Unmarshal(rawPayload.Assets, &assetKeys); err2 == nil {
					logger.Info("Detected legacy string assets format", "count", len(assetKeys))
					for _, key := range assetKeys {
						assetType := "PHOTO"
						if strings.Contains(key, "/video/") {
							assetType = "VIDEO"
						}
						assets = append(assets, mediaprocessingmodel.JobAsset{
							Key:  key,
							Type: assetType,
						})
					}
				} else {
					logger.Error("Failed to parse assets as struct or string", "error", err)
					continue
				}
			}

			payload := mediaprocessingmodel.StepFunctionPayload{
				JobID:     rawPayload.JobID,
				BatchID:   rawPayload.BatchID,
				ListingID: rawPayload.ListingID,
				Assets:    assets,
			}

			// Start Step Function
			inputBytes, _ := json.Marshal(payload)
			inputStr := string(inputBytes)

			_, err := sfnClient.StartExecution(ctx, &sfn.StartExecutionInput{
				StateMachineArn: aws.String(stateMachineArn),
				Input:           aws.String(inputStr),
			})

			if err != nil {
				logger.Error("Failed to start Step Function", "error", err, "batch_id", payload.BatchID)
				return mediaprocessingmodel.StepFunctionPayload{}, err // Fail the lambda so SQS retries
			}

			logger.Info("Started Step Function execution", "batch_id", payload.BatchID)
		}

		return mediaprocessingmodel.StepFunctionPayload{}, nil // Return empty success
	}

	var event mediaprocessingmodel.StepFunctionPayload
	if err := json.Unmarshal(rawEvent, &event); err != nil {
		logger.Error("Failed to unmarshal event", "error", err)
		return mediaprocessingmodel.StepFunctionPayload{}, err
	}

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

	// Check for videos
	for _, asset := range validAssets {
		if strings.Contains(asset.Type, "VIDEO") {
			event.HasVideos = true
			break
		}
	}

	// LOG: Final summary
	logger.Info("Validation complete",
		"batch_id", event.BatchID,
		"valid_count", len(validAssets),
		"invalid_count", len(event.Assets)-len(validAssets),
		"has_videos", event.HasVideos,
	)

	return event, nil
}

func main() {
	lambda.Start(HandleRequest)
}
