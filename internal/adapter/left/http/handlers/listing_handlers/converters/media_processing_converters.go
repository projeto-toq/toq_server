package converters

import (
	"encoding/json"

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

// DTOToListMediaInput converts HTTP request to service input
func DTOToListMediaInput(req dto.ListMediaRequest) domaindto.ListMediaInput {
	return domaindto.ListMediaInput{
		ListingIdentityID: req.ListingIdentityID,
		AssetType:         req.AssetType,
		Sequence:          req.Sequence,
		Page:              req.Page,
		Limit:             req.Limit,
		Sort:              req.Sort,
		Order:             req.Order,
	}
}

// ListMediaOutputToDTO converts service output to HTTP response
func ListMediaOutputToDTO(output domaindto.ListMediaOutput) dto.ListMediaResponse {
	data := make([]dto.MediaAssetResponse, 0, len(output.Assets))
	for _, a := range output.Assets {
		var metaMap map[string]string
		if metaStr := a.Metadata(); metaStr != "" {
			_ = json.Unmarshal([]byte(metaStr), &metaMap)
		}

		data = append(data, dto.MediaAssetResponse{
			ID:                a.ID(),
			ListingIdentityID: a.ListingIdentityID(),
			AssetType:         string(a.AssetType()),
			Sequence:          a.Sequence(),
			Status:            string(a.Status()),
			Title:             a.Title(),
			Metadata:          metaMap,
			S3KeyRaw:          a.S3KeyRaw(),
			S3KeyProcessed:    a.S3KeyProcessed(),
		})
	}

	return dto.ListMediaResponse{
		Data: data,
		Pagination: dto.PaginationResponse{
			Page:  output.Page,
			Limit: output.Limit,
			Total: output.TotalCount,
		},
	}
}

// DTOToGenerateDownloadURLsInput converts HTTP request to service input
func DTOToGenerateDownloadURLsInput(req dto.GenerateDownloadURLsRequest) domaindto.GenerateDownloadURLsInput {
	requests := make([]domaindto.DownloadRequestItemInput, 0, len(req.Requests))
	for _, r := range req.Requests {
		requests = append(requests, domaindto.DownloadRequestItemInput{
			AssetType:  mediaprocessingmodel.MediaAssetType(r.AssetType),
			Sequence:   r.Sequence,
			Resolution: r.Resolution,
		})
	}

	return domaindto.GenerateDownloadURLsInput{
		ListingIdentityID: req.ListingIdentityID,
		Requests:          requests,
	}
}

// GenerateDownloadURLsOutputToDTO converts service output to HTTP response
func GenerateDownloadURLsOutputToDTO(output domaindto.GenerateDownloadURLsOutput) dto.GenerateDownloadURLsResponse {
	urls := make([]dto.DownloadURLResponse, 0, len(output.Urls))
	for _, u := range output.Urls {
		urls = append(urls, dto.DownloadURLResponse{
			AssetType:  string(u.AssetType),
			Sequence:   u.Sequence,
			Resolution: u.Resolution,
			Url:        u.Url,
			ExpiresIn:  u.ExpiresIn,
		})
	}

	return dto.GenerateDownloadURLsResponse{
		ListingIdentityID: output.ListingIdentityID,
		Urls:              urls,
	}
}
