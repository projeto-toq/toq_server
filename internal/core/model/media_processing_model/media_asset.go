package mediaprocessingmodel

import (
	"database/sql"
)

// MediaAsset represents a single media file associated with a listing.
// Audit fields removed as per project standard.
type MediaAsset struct {
	id             uint64
	listingID      uint64
	assetType      MediaAssetType
	sequence       uint8
	status         MediaAssetStatus
	s3KeyRaw       sql.NullString
	s3KeyProcessed sql.NullString
	title          sql.NullString
	metadata       sql.NullString // JSON string
}

// NewMediaAsset creates a new MediaAsset instance.
func NewMediaAsset(listingID uint64, assetType MediaAssetType, sequence uint8) MediaAsset {
	return MediaAsset{
		listingID: listingID,
		assetType: assetType,
		sequence:  sequence,
		status:    MediaAssetStatusPendingUpload,
	}
}

func (a *MediaAsset) ID() uint64                   { return a.id }
func (a *MediaAsset) SetID(id uint64)              { a.id = id }
func (a *MediaAsset) ListingID() uint64            { return a.listingID }
func (a *MediaAsset) AssetType() MediaAssetType    { return a.assetType }
func (a *MediaAsset) Sequence() uint8              { return a.sequence }
func (a *MediaAsset) Status() MediaAssetStatus     { return a.status }
func (a *MediaAsset) SetStatus(s MediaAssetStatus) { a.status = s }

func (a *MediaAsset) S3KeyRaw() string {
	if a.s3KeyRaw.Valid {
		return a.s3KeyRaw.String
	}
	return ""
}
func (a *MediaAsset) SetS3KeyRaw(key string) {
	a.s3KeyRaw = sql.NullString{String: key, Valid: key != ""}
}

func (a *MediaAsset) S3KeyProcessed() string {
	if a.s3KeyProcessed.Valid {
		return a.s3KeyProcessed.String
	}
	return ""
}
func (a *MediaAsset) SetS3KeyProcessed(key string) {
	a.s3KeyProcessed = sql.NullString{String: key, Valid: key != ""}
}

func (a *MediaAsset) Title() string {
	if a.title.Valid {
		return a.title.String
	}
	return ""
}
func (a *MediaAsset) SetTitle(title string) {
	a.title = sql.NullString{String: title, Valid: title != ""}
}

func (a *MediaAsset) Metadata() string {
	if a.metadata.Valid {
		return a.metadata.String
	}
	return ""
}
func (a *MediaAsset) SetMetadata(metadata string) {
	a.metadata = sql.NullString{String: metadata, Valid: metadata != ""}
}
