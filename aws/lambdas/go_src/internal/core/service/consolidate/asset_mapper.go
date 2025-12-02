package consolidate

import (
	"strings"

	mediaprocessingmodel "github.com/projeto-toq/toq_server/internal/core/model/media_processing_model"
)

// PayloadAccumulator keeps track of the aggregated payload for one original asset
// plus helper metadata (asset type and best processed resolution seen so far).
type PayloadAccumulator struct {
	payload           *mediaprocessingmodel.MediaProcessingJobPayload
	assetType         string
	bestProcessedRank int
}

// BranchError represents failures surfaced by parallel branches so that the
// backend can persist consistent diagnostics for each asset.
type BranchError struct {
	SourceKey    string
	ErrorCode    string
	ErrorMessage string
}

// InitializePayloads builds the accumulator map using the raw assets received
// from the backend. Validation errors detected earlier are preserved.
func InitializePayloads(assets []mediaprocessingmodel.JobAsset) map[string]*PayloadAccumulator {
	payloads := make(map[string]*PayloadAccumulator, len(assets))

	for _, asset := range assets {
		acc := &PayloadAccumulator{
			payload: &mediaprocessingmodel.MediaProcessingJobPayload{
				RawKey:  asset.Key,
				Outputs: make(map[string]string),
			},
			assetType:         strings.ToLower(asset.Type),
			bestProcessedRank: -1,
		}

		if asset.Error != "" {
			acc.payload.ErrorCode = "VALIDATION_ERROR"
			acc.payload.ErrorMessage = asset.Error
		}

		payloads[asset.Key] = acc
	}

	return payloads
}

// MapGeneratedAsset enriches the accumulator with a derived object (thumbnail or
// resized image), updating the canonical processed key based on resolution
// priority and storing thumbnails separately.
func MapGeneratedAsset(acc *PayloadAccumulator, derivative mediaprocessingmodel.JobAsset) {
	if acc == nil || acc.payload == nil {
		return
	}

	resolution := extractResolution(derivative.Key)
	outputsKey := buildOutputsKey(resolution, acc.assetType)
	acc.payload.Outputs[outputsKey] = derivative.Key

	if strings.EqualFold(resolution, "thumbnail") {
		acc.payload.ThumbnailKey = derivative.Key
	}

	rank := processedResolutionRank(resolution)
	if rank > acc.bestProcessedRank {
		acc.payload.ProcessedKey = derivative.Key
		acc.bestProcessedRank = rank
	}
}

// ApplyBranchErrors attaches errors reported by derived processing stages to the
// related payloads so the backend can expose them to clients.
func ApplyBranchErrors(payloads map[string]*PayloadAccumulator, branchErrors []BranchError) {
	if len(branchErrors) == 0 {
		return
	}

	for _, branchErr := range branchErrors {
		if branchErr.SourceKey == "" {
			continue
		}

		if acc, ok := payloads[branchErr.SourceKey]; ok && acc.payload != nil {
			acc.payload.ErrorCode = branchErr.ErrorCode
			acc.payload.ErrorMessage = branchErr.ErrorMessage
		}
	}
}

// FlattenPayloads converts the accumulator map back into a slice compatible with
// the callback response contract expected by the backend.
func FlattenPayloads(payloads map[string]*PayloadAccumulator) []mediaprocessingmodel.MediaProcessingJobPayload {
	flattened := make([]mediaprocessingmodel.MediaProcessingJobPayload, 0, len(payloads))
	for _, acc := range payloads {
		if acc == nil || acc.payload == nil {
			continue
		}
		flattened = append(flattened, *acc.payload)
	}
	return flattened
}

func buildOutputsKey(resolution string, assetType string) string {
	normalizedResolution := strings.ToLower(strings.TrimSpace(resolution))
	if normalizedResolution == "" {
		normalizedResolution = "original"
	}

	normalizedType := strings.ReplaceAll(strings.ToLower(strings.TrimSpace(assetType)), " ", "_")
	if normalizedType == "" {
		normalizedType = "unknown"
	}

	return normalizedResolution + "_" + normalizedType
}

func extractResolution(key string) string {
	lowerKey := strings.ToLower(key)
	idx := strings.Index(lowerKey, "/processed/")
	if idx == -1 {
		return ""
	}

	remainder := key[idx+len("/processed/"):]
	segments := strings.Split(remainder, "/")
	if len(segments) < 3 {
		return ""
	}

	// segments pattern: mediaType / orientation / resolution / filename
	return segments[2]
}

func processedResolutionRank(resolution string) int {
	switch strings.ToLower(strings.TrimSpace(resolution)) {
	case "original":
		return 4
	case "large":
		return 3
	case "medium":
		return 2
	case "small":
		return 1
	case "thumbnail":
		return 0
	default:
		return -1
	}
}
