package validate

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"path"
	"regexp"
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

		executionArn, err := s.workflow.StartExecution(ctx, s.smArn, inputStr)
		if err != nil {
			s.logger.Error("Failed to start Step Function", "error", err, "job_id", payload.JobID)
			return err // Fail the lambda so SQS retries
		}

		s.logger.Info("Started Step Function execution", "job_id", payload.JobID, "execution_arn", executionArn)
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
	var firstVideoKey string

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

		if firstVideoKey == "" && strings.Contains(strings.ToUpper(asset.Type), "VIDEO") {
			firstVideoKey = asset.Key
		}
	}

	event.Assets = validatedAssets

	// Check for videos and prepare MediaConvert paths
	if firstVideoKey != "" {
		event.HasVideos = true

		videoInput, videoOutput := deriveVideoPaths(s.bucket, firstVideoKey)
		if videoInput != "" && videoOutput != "" {
			event.VideoInput = videoInput
			event.VideoOutputPath = videoOutput
			s.logger.Info("Video processing paths prepared",
				"job_id", event.JobID,
				"video_input", videoInput,
				"video_output_path", videoOutput,
			)
		} else {
			s.logger.Warn("Could not derive video paths from asset key", "job_id", event.JobID, "asset_key", firstVideoKey)
		}
	}

	s.logger.Info("Validation complete",
		"job_id", event.JobID,
		"total_assets", len(validatedAssets),
		"has_videos", event.HasVideos,
	)

	return event, nil
}

var dateSegmentRegex = regexp.MustCompile(`\d{4}-\d{2}-\d{2}/`)

func deriveVideoPaths(bucket, rawKey string) (string, string) {
	if rawKey == "" || bucket == "" {
		return "", ""
	}

	input := fmt.Sprintf("s3://%s/%s", bucket, rawKey)

	parts := strings.SplitN(rawKey, "raw/", 2)
	if len(parts) != 2 {
		return input, ""
	}

	prefix := parts[0]
	suffix := dateSegmentRegex.ReplaceAllString(parts[1], "")
	mediaDir := path.Dir(suffix) // e.g., video/horizontal

	if strings.TrimSpace(mediaDir) == "" {
		return input, ""
	}

	output := fmt.Sprintf("s3://%s/%sprocessed/%s/original/", bucket, prefix, mediaDir)

	return input, output
}
