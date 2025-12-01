package converters

import (
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	domaindto "github.com/projeto-toq/toq_server/internal/core/domain/dto"
	mediaprocessingmodel "github.com/projeto-toq/toq_server/internal/core/model/media_processing_model"
)

// DTOToRequestUploadURLsInput converts HTTP request to service input
func DTOToRequestUploadURLsInput(req dto.RequestUploadURLsRequest) domaindto.RequestUploadURLsInput {
	files := make([]domaindto.RequestUploadFile, 0, len(req.Files))
	for _, f := range req.Files {
		files = append(files, domaindto.RequestUploadFile{
			AssetType:   mediaprocessingmodel.MediaAssetType(f.AssetType),
			Sequence:    f.Sequence,
			Filename:    f.Filename,
			ContentType: f.ContentType,
			Bytes:       f.Bytes,
			Checksum:    f.Checksum,
			Title:       f.Title,
			Metadata:    f.Metadata,
		})
	}

	return domaindto.RequestUploadURLsInput{
		ListingIdentityID: int64(req.ListingIdentityID),
		Files:             files,
	}
}

// RequestUploadURLsOutputToDTO converts service output to HTTP response
func RequestUploadURLsOutputToDTO(output domaindto.RequestUploadURLsOutput) dto.RequestUploadURLsResponse {
	files := make([]dto.RequestUploadInstructionResponse, 0, len(output.Files))
	for _, f := range output.Files {
		files = append(files, dto.RequestUploadInstructionResponse{
			AssetType: f.AssetType,
			UploadURL: f.UploadURL,
			Method:    f.Method,
			Headers:   f.Headers,
			ObjectKey: f.ObjectKey,
			Sequence:  f.Sequence,
			Title:     f.Title,
		})
	}

	return dto.RequestUploadURLsResponse{
		ListingIdentityID:   uint64(output.ListingIdentityID),
		UploadURLTTLSeconds: output.UploadURLTTLSeconds,
		Files:               files,
	}
}

// DTOToListDownloadURLsInput converts HTTP request to service input
func DTOToListDownloadURLsInput(req dto.ListDownloadURLsRequest) domaindto.ListDownloadURLsInput {
	assetTypes := make([]mediaprocessingmodel.MediaAssetType, 0, len(req.AssetTypes))
	for _, at := range req.AssetTypes {
		assetTypes = append(assetTypes, mediaprocessingmodel.MediaAssetType(at))
	}

	return domaindto.ListDownloadURLsInput{
		ListingIdentityID: int64(req.ListingIdentityID),
		AssetTypes:        assetTypes,
	}
}

// ListDownloadURLsOutputToDTO converts service output to HTTP response
func ListDownloadURLsOutputToDTO(output domaindto.ListDownloadURLsOutput) dto.ListDownloadURLsResponse {
	downloads := make([]dto.DownloadEntryResponse, 0, len(output.Assets))
	for _, d := range output.Assets {
		downloads = append(downloads, dto.DownloadEntryResponse{
			AssetType: string(d.AssetType),
			Sequence:  d.Sequence,
			Status:    string(d.Status),
			Title:     d.Title,
			URL:       d.DownloadURL,
		})
	}

	return dto.ListDownloadURLsResponse{
		ListingIdentityID: uint64(output.ListingIdentityID),
		Downloads:         downloads,
	}
}
