package dto

import mediaprocessingmodel "github.com/projeto-toq/toq_server/internal/core/model/media_processing_model"

// RequestUploadURLsInput defines the input for generating upload URLs.
type RequestUploadURLsInput struct {
	ListingIdentityID int64               `json:"listingIdentityId" validate:"required,gt=0"`
	Files             []RequestUploadFile `json:"files" validate:"required,min=1,dive"`
	RequestedBy       uint64              `json:"-"`
}

type RequestUploadFile struct {
	AssetType   mediaprocessingmodel.MediaAssetType `json:"assetType" validate:"required"`
	Sequence    uint8                               `json:"sequence" validate:"required,gt=0"`
	Filename    string                              `json:"filename" validate:"required"`
	ContentType string                              `json:"contentType" validate:"required"`
	Bytes       int64                               `json:"bytes" validate:"required,gt=0"`
	Checksum    string                              `json:"checksum" validate:"required"`
	Title       string                              `json:"title"`
	Metadata    map[string]string                   `json:"metadata"`
}

// ProcessMediaInput defines the input for triggering media processing.
type ProcessMediaInput struct {
	ListingIdentityID int64  `json:"listingIdentityId" validate:"required,gt=0"`
	RequestedBy       uint64 `json:"-"`
}

// CompleteMediaInput defines the input for finalizing media processing.
type CompleteMediaInput struct {
	ListingIdentityID int64  `json:"listingIdentityId" validate:"required,gt=0"`
	RequestedBy       uint64 `json:"-"`
}

// ListDownloadURLsInput defines filters for listing media URLs.
type ListDownloadURLsInput struct {
	ListingIdentityID int64                                 `json:"listingIdentityId" validate:"required,gt=0"`
	AssetTypes        []mediaprocessingmodel.MediaAssetType `json:"assetTypes"` // Optional filter
	RequestedBy       uint64                                `json:"-"`
}

// ListDownloadURLsOutput contains the list of assets with signed URLs.
type ListDownloadURLsOutput struct {
	ListingIdentityID int64           `json:"listingIdentityId"`
	Assets            []DownloadAsset `json:"assets"`
}

type DownloadAsset struct {
	AssetType    mediaprocessingmodel.MediaAssetType   `json:"assetType"`
	Sequence     uint8                                 `json:"sequence"`
	Status       mediaprocessingmodel.MediaAssetStatus `json:"status"`
	Title        string                                `json:"title"`
	DownloadURL  string                                `json:"downloadUrl,omitempty"` // Signed URL (Processed or Raw)
	ThumbnailURL string                                `json:"thumbnailUrl,omitempty"`
	Metadata     map[string]string                     `json:"metadata,omitempty"`
}

// UpdateMediaInput defines the input for updating a media asset.
type UpdateMediaInput struct {
	ListingIdentityID int64                               `json:"listingIdentityId" validate:"required,gt=0"`
	AssetType         mediaprocessingmodel.MediaAssetType `json:"assetType" validate:"required"`
	Sequence          uint8                               `json:"sequence" validate:"required,gt=0"`
	Title             string                              `json:"title"`
	Metadata          map[string]string                   `json:"metadata"`
	RequestedBy       uint64                              `json:"-"`
}

// DeleteMediaInput defines the input for deleting a media asset.
type DeleteMediaInput struct {
	ListingIdentityID int64                               `json:"listingIdentityId" validate:"required,gt=0"`
	AssetType         mediaprocessingmodel.MediaAssetType `json:"assetType" validate:"required"`
	Sequence          uint8                               `json:"sequence" validate:"required,gt=0"`
	RequestedBy       uint64                              `json:"-"`
}

// HandleProcessingCallbackInput defines the payload received from the processing pipeline.
type HandleProcessingCallbackInput struct {
	JobID             uint64             `json:"jobId"`
	ListingIdentityID uint64             `json:"listingIdentityId"`
	Provider          string             `json:"provider"`
	Status            string             `json:"status"` // "SUCCEEDED", "FAILED"
	Results           []ProcessingResult `json:"results"`
	Error             string             `json:"error,omitempty"`
	ErrorCode         string             `json:"errorCode,omitempty"`
	ErrorMetadata     map[string]string  `json:"errorMetadata,omitempty"`
	FailureReason     string             `json:"failureReason,omitempty"`
	RawPayload        string             `json:"-"`
}

type ProcessingResult struct {
	AssetID      uint64            `json:"assetId"`
	RawKey       string            `json:"rawKey"` // Alternative to AssetID
	Status       string            `json:"status"` // "PROCESSED", "FAILED"
	ProcessedKey string            `json:"processedKey,omitempty"`
	ThumbnailKey string            `json:"thumbnailKey,omitempty"`
	Metadata     map[string]string `json:"metadata,omitempty"`
	Error        string            `json:"error,omitempty"`
	ErrorCode    string            `json:"errorCode,omitempty"`
}

type HandleProcessingCallbackOutput struct {
	Success bool `json:"success"`
}

// ListMediaInput define a entrada para o serviço de listagem.
type ListMediaInput struct {
	ListingIdentityID uint64
	AssetType         string
	Sequence          *uint8
	Page              int
	Limit             int
	Sort              string
	Order             string
}

// ListMediaOutput define a saída do serviço de listagem.
type ListMediaOutput struct {
	Assets     []mediaprocessingmodel.MediaAsset
	TotalCount int64
	Page       int
	Limit      int
}

// GenerateDownloadURLsInput define a entrada para geração de URLs.
type GenerateDownloadURLsInput struct {
	ListingIdentityID uint64
	Requests          []DownloadRequestItemInput
}

type DownloadRequestItemInput struct {
	AssetType  mediaprocessingmodel.MediaAssetType
	Sequence   uint8
	Resolution string
}

// GenerateDownloadURLsOutput define a saída com as URLs geradas.
type GenerateDownloadURLsOutput struct {
	ListingIdentityID uint64
	Urls              []DownloadURLOutput
}

type DownloadURLOutput struct {
	AssetType  mediaprocessingmodel.MediaAssetType
	Sequence   uint8
	Resolution string
	Url        string
	ExpiresIn  int
}
