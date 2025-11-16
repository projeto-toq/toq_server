package converters

import (
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	mediaprocessingmodel "github.com/projeto-toq/toq_server/internal/core/model/media_processing_model"
	mediaprocessingservice "github.com/projeto-toq/toq_server/internal/core/service/media_processing_service"
)

// DTOToCreateUploadBatchInput converts HTTP request to service input
func DTOToCreateUploadBatchInput(req dto.CreateUploadBatchRequest) mediaprocessingservice.CreateUploadBatchInput {
	files := make([]mediaprocessingservice.CreateUploadBatchFile, 0, len(req.Files))
	for _, f := range req.Files {
		files = append(files, mediaprocessingservice.CreateUploadBatchFile{
			ClientID:    f.ClientID,
			AssetType:   mediaprocessingmodel.MediaAssetType(f.AssetType),
			Orientation: mediaprocessingmodel.MediaAssetOrientation(f.Orientation),
			Filename:    f.Filename,
			ContentType: f.ContentType,
			Bytes:       f.Bytes,
			Checksum:    f.Checksum,
			Title:       f.Title,
			Sequence:    f.Sequence,
			Metadata:    f.Metadata,
		})
	}

	return mediaprocessingservice.CreateUploadBatchInput{
		ListingID:      req.ListingID,
		BatchReference: req.BatchReference,
		Files:          files,
	}
}

// CreateUploadBatchOutputToDTO converts service output to HTTP response
func CreateUploadBatchOutputToDTO(output mediaprocessingservice.CreateUploadBatchOutput) dto.CreateUploadBatchResponse {
	files := make([]dto.UploadInstructionResponse, 0, len(output.Files))
	for _, f := range output.Files {
		files = append(files, dto.UploadInstructionResponse{
			ClientID:  f.ClientID,
			UploadURL: f.UploadURL,
			Method:    f.Method,
			Headers:   f.Headers,
			ObjectKey: f.ObjectKey,
			Sequence:  f.Sequence,
			Title:     f.Title,
		})
	}

	return dto.CreateUploadBatchResponse{
		ListingID:           output.ListingID,
		BatchID:             output.BatchID,
		UploadURLTTLSeconds: output.UploadURLTTLSeconds,
		Files:               files,
	}
}

// DTOToCompleteUploadBatchInput converts HTTP request to service input
func DTOToCompleteUploadBatchInput(req dto.CompleteUploadBatchRequest) mediaprocessingservice.CompleteUploadBatchInput {
	files := make([]mediaprocessingservice.CompletedUploadFile, 0, len(req.Files))
	for _, f := range req.Files {
		files = append(files, mediaprocessingservice.CompletedUploadFile{
			ClientID:  f.ClientID,
			ObjectKey: f.ObjectKey,
			Bytes:     f.Bytes,
			Checksum:  f.Checksum,
			ETag:      f.ETag,
		})
	}

	return mediaprocessingservice.CompleteUploadBatchInput{
		ListingID: req.ListingID,
		BatchID:   req.BatchID,
		Files:     files,
	}
}

// CompleteUploadBatchOutputToDTO converts service output to HTTP response
func CompleteUploadBatchOutputToDTO(output mediaprocessingservice.CompleteUploadBatchOutput) dto.CompleteUploadBatchResponse {
	return dto.CompleteUploadBatchResponse{
		ListingID:                output.ListingID,
		BatchID:                  output.BatchID,
		JobID:                    output.JobID,
		Status:                   output.Status.String(),
		EstimatedDurationSeconds: int(output.EstimatedDuration.Seconds()),
	}
}

// DTOToGetBatchStatusInput converts HTTP request to service input
func DTOToGetBatchStatusInput(req dto.GetBatchStatusRequest) mediaprocessingservice.GetBatchStatusInput {
	return mediaprocessingservice.GetBatchStatusInput{
		ListingID: req.ListingID,
		BatchID:   req.BatchID,
	}
}

// GetBatchStatusOutputToDTO converts service output to HTTP response
func GetBatchStatusOutputToDTO(output mediaprocessingservice.GetBatchStatusOutput) dto.GetBatchStatusResponse {
	assets := make([]dto.BatchAssetStatusResponse, 0, len(output.Assets))
	for _, a := range output.Assets {
		assets = append(assets, dto.BatchAssetStatusResponse{
			ClientID:     a.ClientID,
			Title:        a.Title,
			AssetType:    string(a.AssetType),
			Sequence:     a.Sequence,
			RawObjectKey: a.RawObjectKey,
			ProcessedKey: a.ProcessedKey,
			ThumbnailKey: a.ThumbnailKey,
			Metadata:     a.Metadata,
		})
	}

	return dto.GetBatchStatusResponse{
		ListingID:     output.ListingID,
		BatchID:       output.BatchID,
		Status:        output.Status.String(),
		StatusMessage: output.StatusMessage,
		Assets:        assets,
	}
}

// DTOToListDownloadURLsInput converts HTTP request to service input
func DTOToListDownloadURLsInput(req dto.ListDownloadURLsRequest) mediaprocessingservice.ListDownloadURLsInput {
	return mediaprocessingservice.ListDownloadURLsInput{
		ListingID: req.ListingID,
		BatchID:   req.BatchID,
	}
}

// ListDownloadURLsOutputToDTO converts service output to HTTP response
func ListDownloadURLsOutputToDTO(output mediaprocessingservice.ListDownloadURLsOutput) dto.ListDownloadURLsResponse {
	downloads := make([]dto.DownloadEntryResponse, 0, len(output.Downloads))
	for _, d := range output.Downloads {
		downloads = append(downloads, dto.DownloadEntryResponse{
			ClientID:   d.ClientID,
			AssetType:  string(d.AssetType),
			Title:      d.Title,
			Sequence:   d.Sequence,
			URL:        d.URL,
			ExpiresAt:  d.ExpiresAt.Format("2006-01-02T15:04:05Z07:00"),
			PreviewURL: d.PreviewURL,
			Metadata:   d.Metadata,
		})
	}

	return dto.ListDownloadURLsResponse{
		ListingID:   output.ListingID,
		BatchID:     output.BatchID,
		GeneratedAt: output.GeneratedAt.Format("2006-01-02T15:04:05Z07:00"),
		TTLSeconds:  output.TTLSeconds,
		Downloads:   downloads,
	}
}

// DTOToRetryMediaBatchInput converts HTTP request to service input
func DTOToRetryMediaBatchInput(req dto.RetryMediaBatchRequest) mediaprocessingservice.RetryMediaBatchInput {
	return mediaprocessingservice.RetryMediaBatchInput{
		ListingID: req.ListingID,
		BatchID:   req.BatchID,
		Reason:    req.Reason,
	}
}

// RetryMediaBatchOutputToDTO converts service output to HTTP response
func RetryMediaBatchOutputToDTO(output mediaprocessingservice.RetryMediaBatchOutput) dto.RetryMediaBatchResponse {
	return dto.RetryMediaBatchResponse{
		ListingID: output.ListingID,
		BatchID:   output.BatchID,
		JobID:     output.JobID,
		Status:    output.Status.String(),
	}
}
