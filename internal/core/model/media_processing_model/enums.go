package mediaprocessingmodel

// BatchStatus represents the lifecycle of a media batch within TOQ Server.
type BatchStatus string

const (
	// BatchStatusPendingUpload indicates the manifest was validated and signed URLs are waiting for confirmation.
	BatchStatusPendingUpload BatchStatus = "PENDING_UPLOAD"
	// BatchStatusReceived means all uploads were confirmed with S3 HEAD/GetObjectAttributes checks.
	BatchStatusReceived BatchStatus = "RECEIVED"
	// BatchStatusProcessing states that the batch is currently handled by the async pipeline (Step Functions).
	BatchStatusProcessing BatchStatus = "PROCESSING"
	// BatchStatusReady signals that the pipeline finished successfully and downloads can be generated.
	BatchStatusReady BatchStatus = "READY"
	// BatchStatusFailed indicates that processing failed but raw media is still available for retries.
	BatchStatusFailed BatchStatus = "FAILED"
)

// String returns the textual representation used by logs and metrics.
func (s BatchStatus) String() string {
	return string(s)
}

// IsTerminal reports whether the batch reached a final state.
func (s BatchStatus) IsTerminal() bool {
	return s == BatchStatusReady || s == BatchStatusFailed
}

// MediaAssetType enumerates every supported type stored in listing_media_assets.
type MediaAssetType string

const (
	MediaAssetTypePhotoVertical   MediaAssetType = "PHOTO_VERTICAL"
	MediaAssetTypePhotoHorizontal MediaAssetType = "PHOTO_HORIZONTAL"
	MediaAssetTypeVideoVertical   MediaAssetType = "VIDEO_VERTICAL"
	MediaAssetTypeVideoHorizontal MediaAssetType = "VIDEO_HORIZONTAL"
	MediaAssetTypeThumbnail       MediaAssetType = "THUMBNAIL"
	MediaAssetTypeZip             MediaAssetType = "ZIP"
	MediaAssetTypeProjectDoc      MediaAssetType = "PROJECT_DOC"
	MediaAssetTypeProjectRender   MediaAssetType = "PROJECT_RENDER"
)

// MediaAssetOrientation stores the canonical orientation for assets that support layout decisions.
type MediaAssetOrientation string

const (
	MediaAssetOrientationVertical   MediaAssetOrientation = "VERTICAL"
	MediaAssetOrientationHorizontal MediaAssetOrientation = "HORIZONTAL"
)

// MediaProcessingProvider identifies the external system responsible for a processing job.
type MediaProcessingProvider string

const (
	MediaProcessingProviderStepFunctions             MediaProcessingProvider = "STEP_FUNCTIONS"
	MediaProcessingProviderStepFunctionsFinalization MediaProcessingProvider = "STEP_FUNCTIONS_FINALIZATION"
	MediaProcessingProviderMediaConvert              MediaProcessingProvider = "MEDIACONVERT"
)

// MediaProcessingJobStatus mirrors the async job state reported by Step Functions/MediaConvert.
type MediaProcessingJobStatus string

const (
	MediaProcessingJobStatusPending   MediaProcessingJobStatus = "PENDING"
	MediaProcessingJobStatusRunning   MediaProcessingJobStatus = "RUNNING"
	MediaProcessingJobStatusSucceeded MediaProcessingJobStatus = "SUCCEEDED"
	MediaProcessingJobStatusFailed    MediaProcessingJobStatus = "FAILED"
)

// IsTerminal reports whether the async job reached a final outcome.
func (s MediaProcessingJobStatus) IsTerminal() bool {
	return s == MediaProcessingJobStatusSucceeded || s == MediaProcessingJobStatusFailed
}

// MediaAssetStatus represents the lifecycle of a single media asset.
type MediaAssetStatus string

const (
	MediaAssetStatusPendingUpload MediaAssetStatus = "PENDING_UPLOAD"
	MediaAssetStatusProcessing    MediaAssetStatus = "PROCESSING"
	MediaAssetStatusProcessed     MediaAssetStatus = "PROCESSED"
	MediaAssetStatusFailed        MediaAssetStatus = "FAILED"
)
