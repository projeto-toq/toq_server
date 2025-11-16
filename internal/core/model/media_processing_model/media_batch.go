package mediaprocessingmodel

import "time"

// BatchStatusMetadata stores contextual information attached to a batch status transition.
type BatchStatusMetadata struct {
	Message   string
	Reason    string
	Details   map[string]string
	UpdatedBy uint64
	UpdatedAt time.Time
}

// MediaBatch represents a logical bundle of assets uploaded for a listing.
type MediaBatch struct {
	id             uint64
	listingID      uint64
	reference      string
	status         BatchStatus
	statusMetadata BatchStatusMetadata
	assets         []MediaAsset
	deletedAt      *time.Time
}

// MediaBatchRecord rehydrates a media batch from persistent storage.
type MediaBatchRecord struct {
	ID             uint64
	ListingID      uint64
	Reference      string
	Status         BatchStatus
	StatusMetadata BatchStatusMetadata
	Assets         []MediaAsset
	DeletedAt      *time.Time
}

// RestoreMediaBatch rebuilds a MediaBatch from a storage record.
func RestoreMediaBatch(record MediaBatchRecord) MediaBatch {
	batch := MediaBatch{
		id:             record.ID,
		listingID:      record.ListingID,
		reference:      record.Reference,
		status:         record.Status,
		statusMetadata: record.StatusMetadata,
		assets:         record.Assets,
		deletedAt:      record.DeletedAt,
	}

	if batch.statusMetadata.Details == nil {
		batch.statusMetadata.Details = map[string]string{}
	}

	if batch.assets == nil {
		batch.assets = []MediaAsset{}
	}

	return batch
}

// NewMediaBatch builds a new in-memory representation with sane defaults.
func NewMediaBatch(listingID uint64, reference string, createdBy uint64) MediaBatch {
	return MediaBatch{
		listingID: listingID,
		reference: reference,
		status:    BatchStatusPendingUpload,
		statusMetadata: BatchStatusMetadata{
			Message:   "batch_created",
			UpdatedBy: createdBy,
			UpdatedAt: time.Now(),
		},
	}
}

func (b *MediaBatch) ID() uint64 {
	return b.id
}

func (b *MediaBatch) SetID(id uint64) {
	b.id = id
}

func (b *MediaBatch) ListingID() uint64 {
	return b.listingID
}

func (b *MediaBatch) Reference() string {
	return b.reference
}

func (b *MediaBatch) Status() BatchStatus {
	return b.status
}

func (b *MediaBatch) StatusMetadata() BatchStatusMetadata {
	return b.statusMetadata
}

func (b *MediaBatch) Assets() []MediaAsset {
	return b.assets
}

func (b *MediaBatch) DeletedAt() *time.Time {
	return b.deletedAt
}

func (b *MediaBatch) WithAssets(assets []MediaAsset) {
	b.assets = assets
}

func (b *MediaBatch) UpdateStatus(status BatchStatus, metadata BatchStatusMetadata) {
	b.status = status
	b.statusMetadata = metadata
}

func (b *MediaBatch) MarkDeleted(deletedAt time.Time) {
	b.deletedAt = &deletedAt
}
