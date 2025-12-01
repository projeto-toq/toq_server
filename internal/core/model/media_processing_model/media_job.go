package mediaprocessingmodel

import "time"

// MediaProcessingJob describes an asynchronous job tracked in the database.
type MediaProcessingJob struct {
	id                uint64
	listingIdentityID uint64
	status            MediaProcessingJobStatus
	provider          MediaProcessingProvider
	externalID        string
	payload           MediaProcessingJobPayload
	retryCount        uint16
	startedAt         *time.Time
	completedAt       *time.Time
	lastError         string
	callbackBody      string
}

// MediaProcessingJobRecord rehydrates a job from persistent storage.
type MediaProcessingJobRecord struct {
	ID                uint64
	ListingIdentityID uint64
	Status            MediaProcessingJobStatus
	Provider          MediaProcessingProvider
	ExternalID        string
	Payload           MediaProcessingJobPayload
	RetryCount        uint16
	StartedAt         *time.Time
	CompletedAt       *time.Time
	LastError         string
	CallbackBody      string
}

// RestoreMediaProcessingJob rebuilds a job from a storage record.
func RestoreMediaProcessingJob(record MediaProcessingJobRecord) MediaProcessingJob {
	return MediaProcessingJob{
		id:                record.ID,
		listingIdentityID: record.ListingIdentityID,
		status:            record.Status,
		provider:          record.Provider,
		externalID:        record.ExternalID,
		payload:           record.Payload,
		retryCount:        record.RetryCount,
		startedAt:         record.StartedAt,
		completedAt:       record.CompletedAt,
		lastError:         record.LastError,
		callbackBody:      record.CallbackBody,
	}
}

func NewMediaProcessingJob(listingIdentityID uint64, provider MediaProcessingProvider) MediaProcessingJob {
	return MediaProcessingJob{
		listingIdentityID: listingIdentityID,
		provider:          provider,
		status:            MediaProcessingJobStatusPending,
	}
}

func (j *MediaProcessingJob) ID() uint64                        { return j.id }
func (j *MediaProcessingJob) SetID(id uint64)                   { j.id = id }
func (j *MediaProcessingJob) ListingIdentityID() uint64         { return j.listingIdentityID }
func (j *MediaProcessingJob) Status() MediaProcessingJobStatus  { return j.status }
func (j *MediaProcessingJob) Provider() MediaProcessingProvider { return j.provider }
func (j *MediaProcessingJob) ExternalID() string                { return j.externalID }
func (j *MediaProcessingJob) SetExternalID(externalID string)   { j.externalID = externalID }
func (j *MediaProcessingJob) Payload() MediaProcessingJobPayload {
	return j.payload
}
func (j *MediaProcessingJob) RetryCount() uint16 { return j.retryCount }
func (j *MediaProcessingJob) StartedAt() *time.Time {
	return j.startedAt
}
func (j *MediaProcessingJob) CompletedAt() *time.Time {
	return j.completedAt
}
func (j *MediaProcessingJob) LastError() string    { return j.lastError }
func (j *MediaProcessingJob) CallbackBody() string { return j.callbackBody }

func (j *MediaProcessingJob) MarkRunning(externalID string, startedAt time.Time) {
	j.status = MediaProcessingJobStatusRunning
	j.externalID = externalID
	j.startedAt = &startedAt
}

func (j *MediaProcessingJob) MarkCompleted(status MediaProcessingJobStatus, payload MediaProcessingJobPayload, completedAt time.Time) {
	j.status = status
	j.payload = payload
	j.completedAt = &completedAt
}

func (j *MediaProcessingJob) AppendError(message string) {
	j.lastError = message
}

func (j *MediaProcessingJob) SetCallbackBody(body string) {
	j.callbackBody = body
}

// JobAsset defines the contract between Backend -> SQS -> Lambdas.
type JobAsset struct {
	Key       string `json:"key"`
	Type      string `json:"type"`                // Enum: PHOTO_VERTICAL, VIDEO_HORIZONTAL, etc.
	SourceKey string `json:"sourceKey,omitempty"` // Filled by validation
	Size      int64  `json:"size,omitempty"`
	ETag      string `json:"etag,omitempty"`
	Error     string `json:"error,omitempty"`
}

// StepFunctionPayload is the unified payload for Step Functions.
type StepFunctionPayload struct {
	JobID             uint64     `json:"jobId"` // Added
	ListingIdentityID uint64     `json:"listingIdentityId"`
	Assets            []JobAsset `json:"assets"`    // Raw input
	HasVideos         bool       `json:"hasVideos"` // Flag for video processing
	Traceparent       string     `json:"traceparent"`
}

// LambdaResponse wraps the output to match Step Functions expectation ($.body).
type LambdaResponse struct {
	Body any `json:"body"`
}

// MediaProcessingJobPayload keeps information serialized from the async provider.
type MediaProcessingJobPayload struct {
	RawKey       string            `json:"rawKey"`
	ProcessedKey string            `json:"processedKey"`
	ThumbnailKey string            `json:"thumbnailKey"`
	Outputs      map[string]string `json:"outputs"`
	ErrorCode    string            `json:"errorCode"`
	ErrorMessage string            `json:"errorMessage"`
}

// MediaProcessingJobMessage is the payload sent to SQS/Step Functions.
type MediaProcessingJobMessage struct {
	JobID             uint64     `json:"jobId"`
	ListingIdentityID uint64     `json:"listingIdentityId"`
	Assets            []JobAsset `json:"assets"`
	Retry             uint16     `json:"retry"`
}

// MediaProcessingCallback represents the structure received from the async workflow.
type MediaProcessingCallback struct {
	JobID         uint64                      `json:"jobId"`
	ExternalID    string                      `json:"externalId"`
	Status        MediaProcessingJobStatus    `json:"status"`
	Provider      MediaProcessingProvider     `json:"provider"`
	Outputs       []MediaProcessingJobPayload `json:"outputs"`
	FailureReason string                      `json:"failureReason"`
	Error         any                         `json:"error"`
	RawBody       string
}
