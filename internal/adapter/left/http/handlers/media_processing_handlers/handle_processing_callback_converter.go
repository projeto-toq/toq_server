package mediaprocessinghandlers

import (
	"strings"

	httpdto "github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	"github.com/projeto-toq/toq_server/internal/core/derrors"
	"github.com/projeto-toq/toq_server/internal/core/domain/dto"
	mediaprocessingmodel "github.com/projeto-toq/toq_server/internal/core/model/media_processing_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func toHandleProcessingCallbackInput(request httpdto.MediaProcessingCallbackRequest) (dto.HandleProcessingCallbackInput, error) {
	jobID, err := utils.ParseUintFromJSON(request.JobID)
	if err != nil {
		return dto.HandleProcessingCallbackInput{}, derrors.Validation("invalid job identifier", err.Error())
	}

	listingIdentityID, err := utils.ParseUintFromJSON(request.ListingIdentityID)
	if err != nil {
		return dto.HandleProcessingCallbackInput{}, derrors.Validation("invalid listing identifier", err.Error())
	}

	var payloadError string
	var errorCode string
	var errorMetadata map[string]string
	if request.Error != nil {
		payloadError = strings.TrimSpace(request.Error.Message)
		errorCode = strings.ToUpper(strings.TrimSpace(request.Error.Code))
		if len(request.Error.Details) > 0 {
			errorMetadata = make(map[string]string, len(request.Error.Details))
			for k, v := range request.Error.Details {
				trimmed := strings.TrimSpace(v)
				if trimmed == "" {
					continue
				}
				errorMetadata[k] = trimmed
			}
		}
	}

	input := dto.HandleProcessingCallbackInput{
		JobID:             jobID,
		ListingIdentityID: listingIdentityID,
		Provider:          strings.ToUpper(strings.TrimSpace(request.Provider)),
		Status:            strings.ToUpper(strings.TrimSpace(request.Status)),
		Results:           mapOutputsToResults(request.Outputs),
		AssetsZipped:      request.AssetsZipped,
		ZipBundles:        cloneStringSlice(request.ZipBundles),
		ZipSizeBytes:      request.ZipSizeBytes,
		UnzippedSizeBytes: request.UnzippedSizeBytes,
		Error:             payloadError,
		ErrorCode:         errorCode,
		ErrorMetadata:     errorMetadata,
		FailureReason:     request.FailureReason,
		Traceparent:       strings.TrimSpace(request.Traceparent),
		RawPayload:        string(request.RawBody),
	}

	if input.Status == "" {
		input.Status = "UNKNOWN"
	}

	return input, nil
}

func mapOutputsToResults(outputs []mediaprocessingmodel.MediaProcessingJobPayload) []dto.ProcessingResult {
	if len(outputs) == 0 {
		return nil
	}
	results := make([]dto.ProcessingResult, 0, len(outputs))
	for _, output := range outputs {
		status := "PROCESSED"
		if output.ErrorCode != "" || output.ErrorMessage != "" {
			status = "FAILED"
		}

		var metadata map[string]string
		if len(output.Outputs) > 0 {
			metadata = make(map[string]string, len(output.Outputs))
			for k, v := range output.Outputs {
				trimmedKey := strings.TrimSpace(k)
				trimmedValue := strings.TrimSpace(v)
				if trimmedKey == "" || trimmedValue == "" {
					continue
				}
				metadata[trimmedKey] = trimmedValue
			}
		}

		results = append(results, dto.ProcessingResult{
			RawKey:       output.RawKey,
			Status:       status,
			ProcessedKey: output.ProcessedKey,
			ThumbnailKey: output.ThumbnailKey,
			Metadata:     metadata,
			Error:        output.ErrorMessage,
			ErrorCode:    output.ErrorCode,
		})
	}
	return results
}

func cloneStringSlice(values []string) []string {
	if len(values) == 0 {
		return nil
	}
	cloned := make([]string, 0, len(values))
	for _, raw := range values {
		trimmed := strings.TrimSpace(raw)
		if trimmed == "" {
			continue
		}
		cloned = append(cloned, trimmed)
	}
	if len(cloned) == 0 {
		return nil
	}
	return cloned
}
