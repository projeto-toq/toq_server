package mediaprocessingmodel

import "time"

// MediaProcessingJob describes an asynchronous job tracked in the database.
type MediaProcessingJob struct {
	id           uint64
	batchID      uint64
	listingID    uint64
	status       MediaProcessingJobStatus
	provider     MediaProcessingProvider
	externalID   string
	payload      MediaProcessingJobPayload
	retryCount   uint16
	startedAt    *time.Time
	completedAt  *time.Time
	lastError    string
	callbackBody string
}

// MediaProcessingJobRecord rehydrates a job from persistent storage.
type MediaProcessingJobRecord struct {
	ID           uint64
	BatchID      uint64
	ListingID    uint64
	Status       MediaProcessingJobStatus
	Provider     MediaProcessingProvider
	ExternalID   string
	Payload      MediaProcessingJobPayload
	RetryCount   uint16
	StartedAt    *time.Time
	CompletedAt  *time.Time
	LastError    string
	CallbackBody string
}

// RestoreMediaProcessingJob rebuilds a job from a storage record.
func RestoreMediaProcessingJob(record MediaProcessingJobRecord) MediaProcessingJob {
	return MediaProcessingJob{
		id:           record.ID,
		batchID:      record.BatchID,
		listingID:    record.ListingID,
		status:       record.Status,
		provider:     record.Provider,
		externalID:   record.ExternalID,
		payload:      record.Payload,
		retryCount:   record.RetryCount,
		startedAt:    record.StartedAt,
		completedAt:  record.CompletedAt,
		lastError:    record.LastError,
		callbackBody: record.CallbackBody,
	}
}

func NewMediaProcessingJob(batchID, listingID uint64, provider MediaProcessingProvider) MediaProcessingJob {
	return MediaProcessingJob{
		batchID:   batchID,
		listingID: listingID,
		provider:  provider,
		status:    MediaProcessingJobStatusPending,
	}
}

func (j *MediaProcessingJob) ID() uint64                        { return j.id }
func (j *MediaProcessingJob) SetID(id uint64)                   { j.id = id }
func (j *MediaProcessingJob) BatchID() uint64                   { return j.batchID }
func (j *MediaProcessingJob) ListingID() uint64                 { return j.listingID }
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

// MediaProcessingJobPayload keeps information serialized from the async provider.
type MediaProcessingJobPayload struct {
	RawKey       string
	ProcessedKey string
	ThumbnailKey string
	Outputs      map[string]string
	ErrorCode    string
	ErrorMessage string
}

// MediaProcessingJobMessage is the payload sent to SQS/Step Functions.
type MediaProcessingJobMessage struct {
	JobID     uint64   `json:"jobId"`
	BatchID   uint64   `json:"batchId"`
	ListingID uint64   `json:"listingId"`
	Assets    []string `json:"assets"`
	Retry     uint16   `json:"retry"`
}

// MediaProcessingCallback represents the structure received from the async workflow.
type MediaProcessingCallback struct {
	JobID         uint64                      `json:"jobId"`
	ExternalID    string                      `json:"externalId"`
	Status        MediaProcessingJobStatus    `json:"status"`
	Provider      MediaProcessingProvider     `json:"provider"`
	Outputs       []MediaProcessingJobPayload `json:"outputs"`
	FailureReason string                      `json:"failureReason"`
	RawBody       string
}
