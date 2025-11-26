package mediaprocessingentities

import "database/sql"

// JobEntity representa registros em listing_media_jobs.
type JobEntity struct {
	ID         uint64
	BatchID    uint64
	ListingID  uint64 // Populated via JOIN with listing_media_batches
	Status     string
	Provider   string
	ExternalID sql.NullString
	Payload    sql.NullString
	StartedAt  sql.NullTime
	FinishedAt sql.NullTime
}
