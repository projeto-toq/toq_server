package validate

import (
	"context"
	"encoding/json"
	"log/slog"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/projeto-toq/toq_server/aws/lambdas/go_src/internal/core/port"
	mediaprocessingmodel "github.com/projeto-toq/toq_server/internal/core/model/media_processing_model"
)

type ValidateService struct {
	storage  port.StoragePort
	workflow port.WorkflowPort
	bucket   string
	smArn    string
	logger   *slog.Logger
}

func NewValidateService(
	storage port.StoragePort,
	workflow port.WorkflowPort,
	bucket string,
	smArn string,
	logger *slog.Logger,
) *ValidateService {
	return &ValidateService{
		storage:  storage,
		workflow: workflow,
		bucket:   bucket,
		smArn:    smArn,
		logger:   logger,
	}
}

func (s *ValidateService) ProcessSQSEvent(ctx context.Context, records []events.SQSMessage) error {
	for _, record := range records {
		s.logger.Info("Processing SQS record", "message_id", record.MessageId)

		var rawPayload struct {
			JobID             uint64          `json:"jobId"`
			ListingIdentityID uint64          `json:"listingIdentityId"`
			Assets            json.RawMessage `json:"assets"`
			Retry             uint16          `json:"retry"`
		}

		if err := json.Unmarshal([]byte(record.Body), &rawPayload); err != nil {
			s.logger.Error("Failed to unmarshal SQS body structure", "error", err, "body", record.Body)
			continue
		}

		var assets []mediaprocessingmodel.JobAsset
		// Try []JobAsset
		if err := json.Unmarshal(rawPayload.Assets, &assets); err != nil {
			// Try []string (Legacy/Backend mismatch fix)
			var assetKeys []string
			if err2 := json.Unmarshal(rawPayload.Assets, &assetKeys); err2 == nil {
				s.logger.Info("Detected legacy string assets format", "count", len(assetKeys))
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
				s.logger.Error("Failed to parse assets as struct or string", "error", err)
				continue
			}
		}

		payload := mediaprocessingmodel.StepFunctionPayload{
			JobID:             rawPayload.JobID,
			ListingIdentityID: rawPayload.ListingIdentityID,
			Assets:            assets,
		}

		// Start Step Function
		inputBytes, _ := json.Marshal(payload)
		inputStr := string(inputBytes)

		if err := s.workflow.StartExecution(ctx, s.smArn, inputStr); err != nil {
			s.logger.Error("Failed to start Step Function", "error", err, "job_id", payload.JobID)
			return err // Fail the lambda so SQS retries
		}

		s.logger.Info("Started Step Function execution", "job_id", payload.JobID)
	}
	return nil
}

func (s *ValidateService) ValidateAssets(ctx context.Context, event mediaprocessingmodel.StepFunctionPayload) (mediaprocessingmodel.StepFunctionPayload, error) {
	s.logger.Info("Validate Lambda started",
		"job_id", event.JobID,
		"listing_identity_id", event.ListingIdentityID,
		"input_assets_count", len(event.Assets),
	)

	validatedAssets := make([]mediaprocessingmodel.JobAsset, 0, len(event.Assets))

	for _, asset := range event.Assets {
		s.logger.Debug("Validating asset", "key", asset.Key, "job_id", event.JobID)

		size, etag, err := s.storage.GetMetadata(ctx, s.bucket, asset.Key)
		if err != nil {
			s.logger.Error("Asset validation failed",
				"key", asset.Key,
				"job_id", event.JobID,
				"error", err,
			)
			// Mark error instead of discarding
			asset.Error = err.Error()
			asset.SourceKey = asset.Key
		} else {
			asset.Size = size
			asset.ETag = etag
			asset.SourceKey = asset.Key
			s.logger.Debug("Asset valid", "key", asset.Key, "size", asset.Size)
		}

		validatedAssets = append(validatedAssets, asset)
	}

	event.Assets = validatedAssets

	// Check for videos
	for _, asset := range validatedAssets {
		if strings.Contains(asset.Type, "VIDEO") {
			event.HasVideos = true
			break
		}
	}

	s.logger.Info("Validation complete",
		"job_id", event.JobID,
		"total_assets", len(validatedAssets),
		"has_videos", event.HasVideos,
	)

	return event, nil
}
