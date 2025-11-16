package mediaprocessingentities

import "database/sql"

// JobEntity representa registros em listing_media_jobs.
type JobEntity struct {
	ID          uint64
	BatchID     uint64
	ListingID   uint64
	Status      string
	Provider    string
	ExternalID  sql.NullString
	Payload     sql.NullString
	RetryCount  uint16
	StartedAt   sql.NullTime
	CompletedAt sql.NullTime
	LastError   sql.NullString
	CallbackRaw sql.NullString
}
