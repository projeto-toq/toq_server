package mediaprocessingentities

import "database/sql"

// AssetEntity representa registros de listing_media_assets.
type AssetEntity struct {
	ID             uint64
	BatchID        uint64
	ListingID      uint64
	AssetType      string
	Orientation    sql.NullString
	Filename       string
	ContentType    string
	Sequence       uint8
	SizeInBytes    int64
	Checksum       string
	RawObjectKey   string
	ProcessedKey   sql.NullString
	ThumbnailKey   sql.NullString
	Width          sql.NullInt64
	Height         sql.NullInt64
	DurationMillis sql.NullInt64
	Metadata       sql.NullString
}
