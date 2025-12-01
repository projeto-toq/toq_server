package mediaprocessingentities

import "database/sql"

// AssetEntity represents records in media_assets table.
type AssetEntity struct {
	ID                uint64         `db:"id"`
	ListingIdentityID uint64         `db:"listing_identity_id"`
	AssetType         string         `db:"asset_type"`
	Sequence          uint8          `db:"sequence"`
	Status            string         `db:"status"`
	S3KeyRaw          sql.NullString `db:"s3_key_raw"`
	S3KeyProcessed    sql.NullString `db:"s3_key_processed"`
	Title             sql.NullString `db:"title"`
	Metadata          sql.NullString `db:"metadata"`
}
