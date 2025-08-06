package globalentities

import "time"

type AuditEntity struct {
	ID         int64
	ExecutedAT time.Time
	ExecutedBY int64
	TableName  string
	TableID    int64
	Action     string
}
