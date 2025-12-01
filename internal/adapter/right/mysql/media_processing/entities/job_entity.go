package mediaprocessingentities

import "database/sql"

// JobEntity represents records in media_jobs.
type JobEntity struct {
	ID         uint64         `db:"id"`
	ListingID  uint64         `db:"listing_id"`
	Status     string         `db:"status"`
	Provider   string         `db:"provider"`
	ExternalID sql.NullString `db:"external_id"`
	Payload    sql.NullString `db:"payload"`
	StartedAt  sql.NullTime   `db:"started_at"`
	FinishedAt sql.NullTime   `db:"finished_at"`
}
