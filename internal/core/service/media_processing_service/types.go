package mediaprocessingservice

import (
	"time"

	mediaprocessingmodel "github.com/projeto-toq/toq_server/internal/core/model/media_processing_model"
)

// CreateUploadBatchInput carries the manifest required to request signed upload URLs.
type CreateUploadBatchInput struct {
	ListingID      uint64
	RequestedBy    uint64
	BatchReference string
	Files          []CreateUploadBatchFile
}

// CreateUploadBatchFile describes a single asset to be uploaded by the client.
type CreateUploadBatchFile struct {
	ClientID    string
	AssetType   mediaprocessingmodel.MediaAssetType
	Orientation mediaprocessingmodel.MediaAssetOrientation
	Filename    string
	ContentType string
	Bytes       int64
	Checksum    string
	Title       string
	Sequence    uint8
	Metadata    map[string]string
}

// CreateUploadBatchOutput returns signed URLs ready to be used by the uploader.
type CreateUploadBatchOutput struct {
	ListingID           uint64
	BatchID             uint64
	UploadURLTTLSeconds int
	Files               []UploadInstruction
}

// UploadInstruction carries the information required to perform a PUT upload to S3.
type UploadInstruction struct {
	ClientID  string
	UploadURL string
	Method    string
	Headers   map[string]string
	ObjectKey string
	Sequence  uint8
	Title     string
}

// CompleteUploadBatchInput confirms that every asset in the manifest has been uploaded successfully.
type CompleteUploadBatchInput struct {
	ListingID   uint64
	BatchID     uint64
	RequestedBy uint64
	Files       []CompletedUploadFile
}

// CompletedUploadFile links the logical asset to the physical object persisted in S3.
type CompletedUploadFile struct {
	ClientID  string
	ObjectKey string
	Bytes     int64
	Checksum  string
	ETag      string
}

// CompleteUploadBatchOutput exposes the async job metadata published to the processing queue.
type CompleteUploadBatchOutput struct {
	ListingID         uint64
	BatchID           uint64
	JobID             uint64
	Status            mediaprocessingmodel.BatchStatus
	EstimatedDuration time.Duration
}

// GetBatchStatusInput retrieves detailed status for a specific batch under a listing identity.
type GetBatchStatusInput struct {
	ListingID uint64
	BatchID   uint64
}

// GetBatchStatusOutput aggregates batch status and asset metadata for UI polling.
type GetBatchStatusOutput struct {
	ListingID     uint64
	BatchID       uint64
	Status        mediaprocessingmodel.BatchStatus
	StatusMessage string
	Assets        []BatchAssetStatus
}

// BatchAssetStatus mirrors the information frontend needs to display upload and processing progress.
type BatchAssetStatus struct {
	ClientID     string
	Title        string
	AssetType    mediaprocessingmodel.MediaAssetType
	Sequence     uint8
	RawObjectKey string
	ProcessedKey string
	ThumbnailKey string
	Metadata     map[string]string
}

// ListDownloadURLsInput requests signed GET URLs for processed assets.
type ListDownloadURLsInput struct {
	ListingID uint64
	BatchID   uint64 // optional: when zero the most recent READY batch will be used
}

// ListDownloadURLsOutput returns signed URLs for processed assets within a batch.
type ListDownloadURLsOutput struct {
	ListingID   uint64
	BatchID     uint64
	GeneratedAt time.Time
	TTLSeconds  int
	Downloads   []DownloadEntry
}

// DownloadEntry encapsulates individual signed URLs derived from processed assets.
type DownloadEntry struct {
	ClientID   string
	AssetType  mediaprocessingmodel.MediaAssetType
	Title      string
	Sequence   uint8
	URL        string
	ExpiresAt  time.Time
	PreviewURL string
	Metadata   map[string]string
}

// RetryMediaBatchInput allows the caller to re-enqueue a finished batch.
type RetryMediaBatchInput struct {
	ListingID   uint64
	BatchID     uint64
	RequestedBy uint64
	Reason      string
}

// RetryMediaBatchOutput exposes the identifier of the newly created processing job.
type RetryMediaBatchOutput struct {
	ListingID uint64
	BatchID   uint64
	JobID     uint64
	Status    mediaprocessingmodel.BatchStatus
}

// HandleProcessingCallbackInput wraps the payload received from Step Functions or Lambda callbacks.
type HandleProcessingCallbackInput struct {
	Callback      mediaprocessingmodel.MediaProcessingCallback
	ReceiptHandle string
}

// HandleProcessingCallbackOutput returns the resulting state transition after applying the callback.
type HandleProcessingCallbackOutput struct {
	ListingID uint64
	BatchID   uint64
	Status    mediaprocessingmodel.BatchStatus
}
