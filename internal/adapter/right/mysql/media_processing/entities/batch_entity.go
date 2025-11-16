package mediaprocessingentities

import (
	"database/sql"
	"time"
)

// BatchEntity espelha a tabela listing_media_batches.
type BatchEntity struct {
	ID              uint64
	ListingID       uint64
	Reference       string
	Status          string
	StatusMessage   sql.NullString
	StatusReason    sql.NullString
	StatusDetails   sql.NullString
	StatusUpdatedBy uint64
	StatusUpdatedAt time.Time
	DeletedAt       sql.NullTime
}
