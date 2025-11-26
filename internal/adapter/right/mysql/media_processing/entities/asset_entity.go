package mediaprocessingentities

import "database/sql"

// AssetEntity representa registros de listing_media_assets.
type AssetEntity struct {
	ID             uint64
	BatchID        uint64
	AssetType      string
	Orientation    sql.NullString
	RawObjectKey   string // DB: raw_object_key
	ProcessedKey   sql.NullString
	ThumbnailKey   sql.NullString
	Checksum       string // DB: checksum
	ContentType    string
	Filename       string         // DB: filename
	SizeBytes      int64          // DB: size_bytes
	Width          sql.NullInt64  // DB: width
	Height         sql.NullInt64  // DB: height
	DurationMillis sql.NullInt64  // DB: duration_millis
	Title          sql.NullString // DB: title
	Sequence       uint8
	Metadata       sql.NullString // DB: metadata (JSON)
}
